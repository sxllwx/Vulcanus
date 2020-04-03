package broadcast

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

var (
	ErrBroadCasterAlreadyStop = errors.New("broadcaster already stopped")
	ErrWatcherAlreadyStop     = errors.New("watcher already stopped")
)

type FullChannelBehavior int

const (
	WaitIfChannelFull FullChannelBehavior = iota
	DropIfChannelFull
)

type EventType string

const (
	Create    EventType = "create"
	Update    EventType = "update"
	Delete    EventType = "delete"
	opWatcher EventType = ""
)

type Event struct {
	Type   EventType
	Object interface{}
}

type broadCaster struct {

	// lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// protect blow message
	mu       sync.Mutex
	nextID   uint64
	watchers map[uint64]*watcherImpl

	sink chan Event

	b FullChannelBehavior
}

// New a broadcast
func New(parent context.Context, b FullChannelBehavior) *broadCaster {

	ctx, cancel := context.WithCancel(parent)
	out := &broadCaster{
		ctx:      ctx,
		cancel:   cancel,
		watchers: make(map[uint64]*watcherImpl),
		sink:     make(chan Event, 16),
		b:        b,
	}

	return out
}

func (b *broadCaster) loop() {

	for {
		select {
		case e, ok := <-b.sink:
			if !ok {
				// broadcaster stop
				return
			}
			go b.distribute(e)
		}
	}
}

func (b *broadCaster) Watch() (*watcherImpl, error) {

	wg := sync.WaitGroup{}
	wg.Add(1)
	var out *watcherImpl
	e := Event{
		Type: opWatcher,
		Object: func() {
			defer wg.Done()
			out = b.addWatcher()
		},
	}

	wg.Wait()
	return out
}

func (b *broadCaster) distribute(e Event) {

	b.mu.Lock()
	defer b.mu.Unlock()

	if e.Type == opWatcher {
		// handle watcher operation
		e.Object.(func())()
		return
	}

	// this is a event
	for _, w := range b.watchers {
		go w.action(e)
	}
}

func (b *broadCaster) Close() error {
	b.cancel()
	return nil
}

// this method should be protect by mu
func (b *broadCaster) addWatcher() *watcherImpl {

	// new watcher
	currentId := b.nextID
	w := newWatcher(currentId, b)
	b.watchers[currentId] = w

	b.nextID++
	return w
}

// this method should be protect by mu
func (b *broadCaster) deleteWatcher(id uint64) {
	_, ok := b.watchers[id]
	if !ok {
		return
	}
	delete(b.watchers, id)
}

type watcherImpl struct {
	ctx    context.Context
	cancel context.CancelFunc

	id         uint64
	resultChan chan Event

	b *broadCaster
}

func newWatcher(id uint64, b *broadCaster) *watcherImpl {

	ctx, cancel := context.WithCancel(b.ctx)
	return &watcherImpl{
		ctx:        ctx,
		cancel:     cancel,
		id:         id,
		resultChan: make(chan Event),
		b:          b,
	}
}

// broadcaster should promise not send event to a stopped watcher
func (w *watcherImpl) action(e Event) {

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
