package queue

import (
	"sync"
)

import (
	"github.com/pkg/errors"
)

var ErrQueueAlreadyStopped = errors.New("queue already stop")

type DoubleEndQueue interface {
	InsertToHeader(interface{}) error
	InsertToTail(interface{}) error
	PopFromHeader() (interface{}, error)
	PopFromTail() (interface{}, error)
	Stop()
}

type deque struct {

	// 1. protect queue
	// 2. notify consumer to get object
	cond sync.Cond

	store []interface{}

	doneOnce sync.Once // protect just stop once, avoid close channel twice
	done     chan struct{}
}

func NewDoubleEndQueue() DoubleEndQueue {

	q := &deque{
		cond: *sync.NewCond(&sync.Mutex{}),
		done: make(chan struct{}),
	}
	return q
}

func (q *deque) InsertToHeader(o interface{}) error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	select {
	case <-q.done:
		return ErrQueueAlreadyStopped
	default:
	}

	q.store = append(append([]interface{}{}, o), q.store...)
	q.cond.Signal() // send signal
	return nil
}

func (q *deque) InsertToTail(o interface{}) error {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	select {
	case <-q.done:
		return ErrQueueAlreadyStopped
	default:
	}

	q.store = append(q.store, o)
	q.cond.Signal()
	return nil
}

func (q *deque) PopFromHeader() (interface{}, error) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	select {
	case <-q.done:
		return nil, ErrQueueAlreadyStopped
	default:
	}

	// if q.store is empty, wait signal
	for len(q.store) == 0 {
		q.cond.Wait()
	}

	o := q.store[0]
	q.store = append([]interface{}{}, q.store[1:]...)
	return o, nil
}

func (q *deque) PopFromTail() (interface{}, error) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	select {
	case <-q.done:
		return nil, ErrQueueAlreadyStopped
	default:
	}

	o := q.store[len(q.store)-1]
	q.store = append([]interface{}{}, q.store[:len(q.store)-1])
	return o, nil
}

func (q *deque) Stop() {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.doneOnce.Do(func() {
		close(q.done)
	})
}
