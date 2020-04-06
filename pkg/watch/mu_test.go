package watch

import (
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"sync"
	"testing"

	"github.com/sxllwx/vulcanus/pkg/log"
)

func init() {
	go http.ListenAndServe(":6060", nil)
}

func TestBroadcaster(t *testing.T) {
	table := []Event{
		{Type: Create, Object: "create pod-a message"},
		{Type: Update, Object: "update pod-a message"},
		{Type: Delete, Object: "delete pod-a message"},
	}

	// The broadcaster we're testing
	m := New(0, WaitIfChannelFull)

	// Add a bunch of watchers
	const testWatchers = 2
	wg := sync.WaitGroup{}
	wg.Add(testWatchers)
	for i := 0; i < testWatchers; i++ {
		// Verify that each watcher gets the events in the correct order

		w, err := m.Watch()
		if err != nil {
			t.Fatal(err)
		}
		go func(watcher int, w Watcher) {
			tableLine := 0
			for {
				event, ok := <-w.ResultChan()
				if !ok {
					break
				}
				if e, a := table[tableLine], event; !reflect.DeepEqual(e, a) {
					t.Fatalf("Watcher %v, line %v: Expected (%v, %#v), got (%v, %#v)",
						watcher, tableLine, e.Type, e.Object, a.Type, a.Object)
				} else {
					log.Infof("watcher (%d) Got (%v, %#v)", watcher, event.Type, event.Object)
				}
				tableLine++
			}
			wg.Done()
		}(i, w)
	}

	for _, item := range table {
		if err := m.Action(item.Type, item.Object); err != nil {
			t.Fatal(err)
		}
	}

	if err := m.Shutdown(); err != nil {
		t.Fatal(err)
	}

	if err := m.Action(Create, "abc"); err != ErrBroadCasterStopped {
		t.Fatal("shutdown not work")
	}

	wg.Wait()
}

func TestBroadcasterWatcherClose(t *testing.T) {
	m := New(0, WaitIfChannelFull)
	w, err := m.Watch()
	if err != nil {
		t.Fatal(err)
	}
	w2, err := m.Watch()
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Stop(); err != nil {
		t.Fatal(err)
	}
	if err := m.Shutdown(); err != nil {
		t.Fatal(err)
	}

	if _, open := <-w.ResultChan(); open {
		t.Errorf("Stop didn't work?")
	}
	if _, open := <-w2.ResultChan(); open {
		t.Errorf("Shutdown didn't work?")
	}
	// Extra stops don't hurt things
	if err := w.Stop(); err != nil {
		t.Fatal(err)
	}
	if err := w2.Stop(); err != nil {
		t.Fatal(err)
	}
}

//func TestBroadcasterWatcherStopDeadlock(t *testing.T) {
//	done := make(chan bool)
//	m := New(0, WaitIfChannelFull)
//	w0, err := m.Watch()
//	if err != nil {
//		t.Fatal(err)
//	}
//	w1, err := m.Watch()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	go func(w0, w1 *watcherImpl) {
//		// We know Broadcaster is in the distribute loop once one watcher receives
//		// an event. Stop the other watcher while distribute is trying to
//		// send to it.
//		select {
//		case <-w0.ResultChan():
//			if err := w1.Stop(); err != nil {
//				t.Fatal(err)
//			}
//		case <-w1.ResultChan():
//			if err := w0.Stop(); err != nil {
//				t.Fatal(err)
//			}
//		}
//
//		close(done)
//		fmt.Println("called close done")
//	}(w0, w1)
//	if err := m.Action(Event{Create, "a"}); err != nil {
//		t.Fatal(err)
//	}
//
//	select {
//	case <-time.After(time.Second * 30):
//		t.Error("timeout: deadlocked")
//	case <-done:
//	}
//	m.Shutdown()
//}

func TestBroadcasterDropIfChannelFull(t *testing.T) {
	m := New(1, DropIfChannelFull)

	event1 := Event{Type: Create, Object: "hello world 1"}
	event2 := Event{Type: Update, Object: "hello world 2"}

	// Add a couple watchers
	var (
		watches = make([]Watcher, 2)
		err     error
	)
	for i := range watches {
		watches[i], err = m.Watch()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Send a couple events before closing the watch channel.
	if err := m.Action(event1.Type, event1.Object); err != nil {
		t.Fatal(err)
	}
	if err := m.Action(event2.Type, event2.Object); err != nil {
		t.Fatal(err)
	}
	if err := m.Shutdown(); err != nil {
		t.Fatal(err)
	}

	// Pull events from the queue.
	wg := sync.WaitGroup{}
	wg.Add(len(watches))
	for i := range watches {
		go func(watcher int, w Watcher) {
			defer wg.Done()
			for e := range w.ResultChan() {
				t.Logf("Watcher %d Got (%v, %#v)", watcher, e.Type, e.Object)
			}
		}(i, watches[i])
	}
	wg.Wait()
}

func TestBroadCasterDie(t *testing.T) {
	m := New(1, DropIfChannelFull)

	w, err := m.Watch()
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Action(Create, "a"); err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for e := range w.ResultChan() {
			t.Logf("got message: %v", e)
		}
		wg.Done()
	}()

	if err := m.Shutdown(); err != nil {
		t.Fatal(err)
	}

	if m.Action(Create, "asd") != ErrBroadCasterStopped {
		t.Fatal("shutdown now work")
	}

	if err := w.Stop(); err != nil {
		t.Fatal("shutdown now work", err)
	}

	wg.Wait()
}
