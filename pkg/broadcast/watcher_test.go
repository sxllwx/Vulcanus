package broadcast

import (
	"fmt"
	"github.com/sxllwx/vulcanus/pkg/log"
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"sync"
	"testing"
	"time"
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
		go func(watcher int, w *watcherImpl) {
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
		if err := m.Action(item); err != nil {
			t.Fatal(err)
		}
	}

	m.Shutdown()

	if err := m.Action(Event{
		Type:   Create,
		Object: "abc",
	}); err != ErrBroadCasterAlreadyStop {
		t.Fatal("not expect err")
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
	m.Shutdown()

	if _, open := <-w.ResultChan(); open {
		t.Errorf("Stop didn't work?")
	}
	if _, open := <-w2.ResultChan(); open {
		t.Errorf("Shutdown didn't work?")
	}
	// Extra stops don't hurt things
	w.Stop()
	w2.Stop()
}

func TestBroadcasterWatcherStopDeadlock(t *testing.T) {
	done := make(chan bool)
	m := New(0, WaitIfChannelFull)
	w0, err := m.Watch()
	if err != nil {
		t.Fatal(err)
	}
	w1, err := m.Watch()
	if err != nil {
		t.Fatal(err)
	}

	go func(w0, w1 *watcherImpl) {
		// We know Broadcaster is in the distribute loop once one watcher receives
		// an event. Stop the other watcher while distribute is trying to
		// send to it.
		select {
		case <-w0.ResultChan():
			if err := w1.Stop(); err != nil {
				t.Fatal(err)
			}
		case <-w1.ResultChan():
			if err := w0.Stop(); err != nil {
				t.Fatal(err)
			}
		}

		close(done)
		fmt.Println("called close done")
	}(w0, w1)
	if err := m.Action(Event{Create, "a"}); err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.After(time.Second * 30):
		t.Error("timeout: deadlocked")
	case <-done:
	}
	m.Shutdown()
}

func TestBroadcasterDropIfChannelFull(t *testing.T) {
	m := New(1, DropIfChannelFull)

	event1 := Event{Type: Create, Object: "hello world 1"}
	event2 := Event{Type: Create, Object: "hello world 2"}

	// Add a couple watchers
	watches := make([]*watcherImpl, 2)
	var err error
	for i := range watches {
		watches[i], err = m.Watch()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Send a couple events before closing the broadcast channel.
	fmt.Printf("Sending event 1\n")
	if err := m.Action(event1); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Sending event 2\n")
	if err := m.Action(event2); err != nil {
		t.Fatal(err)
	}
	m.Shutdown()

	// Pull events from the queue.
	wg := sync.WaitGroup{}
	wg.Add(len(watches))
	for i := range watches {
		// Verify that each watcher only gets the first event because its watch
		// queue of length one was full from the first one.
		go func(watcher int, w *watcherImpl) {
			defer wg.Done()
			e1, ok := <-w.ResultChan()
			if !ok {
				t.Errorf("Watcher %v failed to retrieve first event.", watcher)
			}
			if e, a := event1, e1; !reflect.DeepEqual(e, a) {
				t.Errorf("Watcher %v: Expected (%v, %#v), got (%v, %#v)",
					watcher, e.Type, e.Object, a.Type, a.Object)
			}
			t.Logf("Watcher %d Got (%v, %#v)", watcher, e1.Type, e1.Object)
		}(i, watches[i])
	}
	wg.Wait()
}
