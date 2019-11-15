package queue

import (
	"context"
	"sync"
)

// the impl of Interface
type impl struct {
	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.Mutex
	store []interface{}

	// the send condition
	sendCond *sync.Cond
	// the receive condition
	rcvCond *sync.Cond

	// check is overflow
	isOverFlow func([]interface{}) bool
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

// NewQueue
// please set isOverflow be simple, because it will hold the lock
func New(opts ...Option) Interface {

	// default setting
	o := &options{
		limitChecker: notLimit,
		ctx:          context.Background(),
	}

	// apply the opts to options
	for _, f := range opts {
		f(o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	out := &impl{
		ctx:        ctx,
		cancel:     cancel,
		store:      []interface{}{},
		isOverFlow: o.limitChecker,
	}

	out.sendCond = sync.NewCond(&out.mutex)
	out.rcvCond = sync.NewCond(&out.mutex)

	return out
}

func (q *impl) Append(o interface{}) error {

	// hold the lock
	q.mutex.Lock()

	// check is queue already overflow
	for q.isOverFlow(q.store) {

		// my be the queue already stopped when wait
		if !q.running() {
			// avoid dead lock
			q.mutex.Unlock()
			return ErrAlreadyStopped
		}

		// wait consumer get the element
		q.sendCond.Wait()
	}

	// double check the queue status
	// maybe the wait() and Close() be called concurrent
	if !q.running() {
		// avoid dead lock
		q.mutex.Unlock()
		return ErrAlreadyStopped
	}

	// yeah, the queue still running
	q.store = append(q.store, o)
	q.mutex.Unlock()

	// notify the consumer try to get element
	q.rcvCond.Broadcast()
	return nil
}

func (q *impl) Insert(o interface{}) error {

	// hold the lock
	q.mutex.Lock()

	// check is queue already overflow
	for q.isOverFlow(q.store) {

		// my be the queue already stopped when wait
		if !q.running() {
			// avoid dead lock
			q.mutex.Unlock()
			return ErrAlreadyStopped
		}

		// wait consumer get the element
		q.sendCond.Wait()
	}

	// double check the queue status
	// maybe the wait() and Close() be called concurrent
	if !q.running() {
		// avoid dead lock
		q.mutex.Unlock()
		return ErrAlreadyStopped
	}

	// yeah, the queue still running
	q.store = append([]interface{}{o}, q.store...)
	q.mutex.Unlock()

	// notify the consumer try to get element
	q.rcvCond.Broadcast()
	return nil
}

func (q *impl) Get() (interface{}, error) {

	q.mutex.Lock()

	for len(q.store) == 0 {
		// no elements
		if !q.running() {
			// avoid dead lock
			q.mutex.Unlock()
			return nil, ErrAlreadyStopped
		}
		q.rcvCond.Wait()
	}

	// double check the queue status
	// maybe the wait() and Close() be called concurrent
	if !q.running() {
		// avoid dead lock
		q.mutex.Unlock()
		return nil, ErrAlreadyStopped
	}

	out := q.store[0]
	q.store = q.store[1:]

	// notify the producer to put the elements
	q.sendCond.Broadcast()
	return out, nil
}

func (q *impl) running() bool {

	select {
	case <-q.ctx.Done():
		return false
	default:
		// the queue still running
	}
	return true
}

func (q *impl) Close() ([]interface{}, error) {

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if !q.running() {
		return nil, ErrAlreadyStopped
	}

	out := q.store
	q.store = nil
	q.cancel()

	// end the put goroutines
	q.sendCond.Broadcast()

	// end the get goroutines
	q.rcvCond.Broadcast()
	return out, nil
}

func (q *impl) Done() <-chan struct{} {
	return q.ctx.Done()
}
