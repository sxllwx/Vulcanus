package watch

import (
	"sync"

	"github.com/pkg/errors"
)

type FullChannelBehavior int

const (
	WaitIfChannelFull FullChannelBehavior = iota
	DropIfChannelFull
)

var (
	ErrBroadCasterStopped = errors.New("broadcaster already stopped")
)

type EventType string

const (
	Create           EventType = "create"
	Update           EventType = "update"
	Delete           EventType = "delete"
	opWatcher        EventType = "op" // special event type, add | delete watcher
	closeBroadCaster EventType = "close-broadcaster"
)

type SynchronizedBroadCaster interface {
	Action(EventType, interface{}) error
	Shutdown() error
	Watch() (Watcher, error)
}

type Watcher interface {
	ResultChan() <-chan Event
	Stop() error
}

// event
type Event struct {
	Type EventType
	// if event type is opWatcher, Obj should be func(){}
	Object interface{}
}

type broadCaster struct {
	nextID   uint64
	watchers map[uint64]*watcherImpl

	watcherBuffSize int

	// event sink,
	once sync.Once
	stop chan struct{}

	sink chan Event

	b FullChannelBehavior
}

// New a watch
func New(watchBufSize int, b FullChannelBehavior) SynchronizedBroadCaster {

	out := &broadCaster{
		watchers:        make(map[uint64]*watcherImpl),
		sink:            make(chan Event),
		watcherBuffSize: watchBufSize,
		b:               b,
		stop:            make(chan struct{}),
	}

	go out.loop()
	return out
}

func (b *broadCaster) loop() {

	for {
		select {
		case e := <-b.sink:

			if e.Type == closeBroadCaster {
				e.Object.(func())()
				// the broadcaster closed
				return
			}

			if e.Type == opWatcher {
				// handle watcher operation
				e.Object.(func())()
				continue
			}

			for _, w := range b.watchers {
				b.sendTo(e, w)
			}
		}
	}
}

func (b *broadCaster) Done() <-chan struct{} {
	return b.stop
}

func (b *broadCaster) Action(eventType EventType, object interface{}) error {

	e := Event{
		Type:   eventType,
		Object: object,
	}

	select {
	case <-b.stop:
		return ErrBroadCasterStopped
	case b.sink <- e:
		return nil
	}
}

func (b *broadCaster) Watch() (Watcher, error) {
	return b.newWatcher()
}

func (b *broadCaster) newWatcher() (Watcher, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var out *watcherImpl
	if err := b.Action(opWatcher, func() {
		out = b.addWatcher()
		wg.Done()
	}); err != nil {
		return nil, err
	}
	wg.Wait()
	return out, nil
}
func (b *broadCaster) stopWatcher(watcher *watcherImpl) error {

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := b.Action(opWatcher, func() {
		b.deleteWatcher(watcher.id)
		wg.Done()
	}); err != nil {
		return err
	}
	wg.Wait()
	return nil
}

// broadcaster promise never send event to a stopped watcher
func (b *broadCaster) sendTo(e Event, w *watcherImpl) {

	switch w.b.b {
	case DropIfChannelFull:
		select {
		case w.resultChan <- e: // send success
		default: // watcher be blocked
		}
	case WaitIfChannelFull:
		select {
		case w.resultChan <- e: // send success
		}
	}
}

func (b *broadCaster) addWatcher() *watcherImpl {

	// new watcher
	currentId := b.nextID
	w := newWatcher(currentId, b)
	b.watchers[currentId] = w

	b.nextID++
	return w
}

func (b *broadCaster) deleteWatcher(id uint64) {

	w, ok := b.watchers[id]
	if !ok {
		// already remove
		return
	}
	delete(b.watchers, id)

	// broadcaster stop watcher's resultChan
	close(w.resultChan)
}

func (b *broadCaster) Shutdown() error {

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := b.Action(closeBroadCaster, func() {
		// remove current watchers
		for i := range b.watchers {
			b.deleteWatcher(i)
		}
		close(b.stop)
		wg.Done()
	}); err != nil {
		return err
	}
	wg.Wait()
	return nil
}

type watcherImpl struct {
	id         uint64
	resultChan chan Event

	stopOnce sync.Once
	stop     chan struct{}

	b *broadCaster
}

func newWatcher(id uint64, b *broadCaster) *watcherImpl {

	return &watcherImpl{
		id:         id,
		resultChan: make(chan Event, b.watcherBuffSize),
		b:          b,
	}
}

func (w *watcherImpl) Stop() error {

	select {
	case <-w.b.stop:
		// broadcaster already stop, watcher already clean
		return nil
	default:
	}
	return w.b.stopWatcher(w)
}

func (w *watcherImpl) ResultChan() <-chan Event {
	return w.resultChan
}
