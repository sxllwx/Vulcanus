package queue

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrQueueAlreadyStopped = errors.New("queue already stop")
	ErrNoMoreElement       = errors.New("no more elements")
)

// this is a thread-safe queue interface like the golang chan
// this queue should be used the event so important, it should be consume as quickly as possible
// multi producer can safe pull the object to store
// multi consumer can safe get the object from store
//  p1  p2  p3 ....
//   \  |  /
//    \ | /
//     \|/
//      |
//   LikeChanQueue
//      |
//     /|\
//    / | \
//   /  |  \
//  c1  c2 c3  ....
type LikeChanQueue interface {

	// append the element to tail
	// block when the queue is full
	// error always not nil, when queue already stopped, golang chan send to a close chan, will panic
	Append(interface{}) error

	// insert the element to header
	// block when the queue is full
	// error always not nil, when queue already stopped, golang chan send to a close chan, will panic
	Insert(interface{}) error

	// Pop will got the object from header
	// if there no element in queue, and queue running, it'will block
	// if there no element in queue, and queue stopped, it'will return nil, and ErrNoMoreElement
	// if there are elements in queue, and queue stopped, it'will return first element, and nil
	Pop() (interface{}, error)

	// Close the queue
	// error always not nil, when queue already stopped
	Close() error

	// the consumer and producer should put/get element on the queue is running
	Done() <-chan struct{}
}

// no limit the queue len
var notLimit = func(objectList []interface{}) bool {
	return false
}

type options struct {
	limitChecker func([]interface{}) bool
	ctx          context.Context
}

type Option func(*options)

func WithLimiter(f func([]interface{}) bool) Option {
	return func(o *options) {
		o.limitChecker = f
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// the likeChanQueueImpl of QueueInterface
type likeChanQueueImpl struct {
	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.Mutex
	store []interface{}

	// the send condition
	sendCond *sync.Cond
	// the receive condition
	rcvCond *sync.Cond

	// check is overflow
	overFlow func([]interface{}) bool
}

// NewQueue
// please set isOverflow be simple, because it will hold the lock
func NewLikeChainQueue(opts ...Option) LikeChanQueue {

	// default setting
	o := &options{
		limitChecker: notLimit,
		ctx:          context.Background(),
	}

	// apply the opts to queueOptions
	for _, f := range opts {
		f(o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	out := &likeChanQueueImpl{
		ctx:      ctx,
		cancel:   cancel,
		store:    []interface{}{},
		overFlow: o.limitChecker,
	}

	out.sendCond = sync.NewCond(&out.mutex)
	out.rcvCond = sync.NewCond(&out.mutex)

	return out
}

func (q *likeChanQueueImpl) Append(o interface{}) error {

	// hold the lock
	q.mutex.Lock()

	// check is queue already overflow
	for q.overFlow(q.store) {

		// my be the queue already stopped when wait
		if !q.running() {
			// avoid dead lock
			q.mutex.Unlock()
			return ErrQueueAlreadyStopped
		}

		// wait consumer get the element
		q.sendCond.Wait()
	}

	// double check the queue status
	// maybe the wait() and Close() be called concurrent
	if !q.running() {
		// avoid dead lock
		q.mutex.Unlock()
		return ErrQueueAlreadyStopped
	}

	// yeah, the queue still running
	q.store = append(q.store, o)
	q.mutex.Unlock()

	// notify the consumer try to get element
	q.rcvCond.Broadcast()
	return nil
}

func (q *likeChanQueueImpl) Insert(o interface{}) error {

	// hold the lock
	q.mutex.Lock()

	// check is queue already overflow
	for q.overFlow(q.store) {

		// my be the queue already stopped when wait
		if !q.running() {
			// avoid dead lock
			q.mutex.Unlock()
			return ErrQueueAlreadyStopped
		}

		// wait consumer get the element
		q.sendCond.Wait()
	}

	// double check the queue status
	// maybe the wait() and Close() be called concurrent
	if !q.running() {
		// avoid dead lock
		q.mutex.Unlock()
		return ErrQueueAlreadyStopped
	}

	// yeah, the queue still running
	q.store = append([]interface{}{o}, q.store...)
	q.mutex.Unlock()

	// notify the consumer try to get element
	q.rcvCond.Broadcast()
	return nil
}

func (q *likeChanQueueImpl) Pop() (interface{}, error) {

	q.mutex.Lock()

	for len(q.store) == 0 {
		if !q.running() {
			// no elements && the queue stopped
			// avoid dead lock
			q.mutex.Unlock()
			return nil, ErrNoMoreElement
		}
		q.rcvCond.Wait()
	}

	// not check the q status
	// direct get the element from the queue
	out := q.store[0]
	q.store = q.store[1:]
	q.mutex.Unlock()

	// notify the producer to put the elements
	q.sendCond.Broadcast()
	return out, nil
}

func (q *likeChanQueueImpl) running() bool {

	select {
	case <-q.ctx.Done():
		return false
	default:
		// the queue still running
	}
	return true
}

func (q *likeChanQueueImpl) Close() error {

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if !q.running() {
		return ErrQueueAlreadyStopped
	}
	q.cancel()

	// end the put goroutines
	q.sendCond.Broadcast()

	// end the get goroutines
	q.rcvCond.Broadcast()
	return nil
}

func (q *likeChanQueueImpl) Done() <-chan struct{} {
	return q.ctx.Done()
}
