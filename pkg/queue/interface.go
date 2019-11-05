package queue

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrAlreadyStopped = errors.New("queue already stop")
)

// this is a thread-safe queue interface
// multi producer can safe pull the object to store
// multi consumer can safe get the object from store
//  p1  p2  p3 ....
//  |   |   |
//   \  |  /
//    \ | /
//   Interface
//      |
//     /|\
//    / | \
//   /  |  \
//  c1  c2 c3  ....
type Interface interface {

	// insert the element to tail
	// error always not nil, when queue already stopped
	Append(interface{}) error

	// insert the element to header
	// error always not nil, when queue already stopped
	Insert(interface{}) error

	// NOTICE:
	// GET will got the object from header
	// if there no element in queue, it'will block
	// error always not nil, when queue already stopped
	GET() (interface{}, error)

	// Close the queue and got the rest object
	// Shuld be call the  by consumer
	// error always not nil, when queue already stopped
	Close() ([]interface{}, error)

	// the consumer and producer should put/get element on the queue is running
	Done() <-chan struct{}
}

// the impl of Interface
type impl struct {
	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.Mutex

	store []interface{}

	isOverflow func([]interface{}) bool

	receiveCond *sync.Cond
	sendCond    *sync.Cond
}

// NewQueue
// please set isOverflow be simple, because it will hold the lock
func NewQueue(isOverflow func([]interface{}) bool) Interface {

	ctx, cancel := context.WithCancel(context.Background())

	out := &impl{
		ctx:        ctx,
		cancel:     cancel,
		store:      []interface{}{},
		isOverflow: isOverflow,
	}
	out.sendCond = sync.NewCond(&out.mutex)
	out.receiveCond = sync.NewCond(&out.mutex)
	return out
}

func (q *impl) Append(o interface{}) error {

	// check queue status
	if err := q.check(); err != nil {
		return err
	}

	q.mutex.Lock()
	for q.isOverflow(q.store) {
		// overflow, wait the element be consumed
		q.sendCond.Wait()
	}

	// double check the status
	if err := q.check(); err != nil {
		q.mutex.Unlock()
		return err
	}

	q.store = append(q.store, o)
	q.mutex.Unlock()

	// notify the consumer get the object
	q.receiveCond.Broadcast()
	return nil
}

func (q *impl) Insert(o interface{}) error {

	// check queue status
	if err := q.check(); err != nil {
		return err
	}

	q.mutex.Lock()
	for q.isOverflow(q.store) {
		// overflow, wait the element be consumed
		q.sendCond.Wait()
	}

	// double check the status
	if err := q.check(); err != nil {
		q.mutex.Unlock()
		return err
	}

	q.store = append([]interface{}{o}, q.store...)
	q.mutex.Unlock()

	// notify the consumer get the object
	q.receiveCond.Broadcast()
	return nil
}

func (q *impl) GET() (interface{}, error) {

	// check queue status
	if err := q.check(); err != nil {
		return nil, err
	}

	q.mutex.Lock()
	for len(q.store) == 0 {
		// empty queue
		q.receiveCond.Wait()
	}

	// when the goroutine running, maybe, the queue already close
	if err := q.check(); err != nil {
		q.mutex.Unlock()
		return nil, err
	}

	out := q.store[0]
	q.store = q.store[1:]
	q.mutex.Unlock()

	// notify the consumer to put the object
	q.sendCond.Broadcast()

	return out, nil
}

func (q *impl) check() error {

	select {
	case <-q.ctx.Done():
		// must unlock to avoid lock
		return ErrAlreadyStopped
	default:
		// the queue still running
	}
	return nil
}

func (q *impl) Close() ([]interface{}, error) {

	if err := q.check(); err != nil {
		return nil, err
	}

	q.mutex.Lock()
	if err := q.check(); err != nil {
		q.mutex.Unlock()
		return nil, err
	}

	out := q.store
	q.store = nil
	q.mutex.Unlock()

	// notify blocking producer && consumer, the queue already close
	q.sendCond.Broadcast()
	q.receiveCond.Broadcast()

	q.cancel()
	return out, nil
}

func (q *impl) Done() <-chan struct{} {
	return q.ctx.Done()
}
