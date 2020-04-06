package cache

import (
	"container/list"
	"context"
	cache2 "github.com/sxllwx/vulcanus/pkg/feature/cache"
	"sync"

	"github.com/sxllwx/vulcanus/pkg/feature/cachere/cache"
)

type queue struct {

	// manage the  stack lifecycle
	cache.LifeCycle

	lock  sync.RWMutex
	store *list.List
}

func (q *queue) Len() (int, error) {

	q.lock.RLock()
	out := q.store.Len()
	q.lock.RUnlock()
	return out, nil
}

func (q *queue) Rest() ([]interface{}, error) {

	if q.Alive() {
		return nil, cache2.ErrRestShouldNotBeCall
	}

	// already stopped
	out := make([]interface{}, 0, 8) // !!!!
	for i := q.store.Front(); i != nil; i = i.Next() {
		out = append(out, i.Value)
	}
	return out, nil
}

// please call me with lock
func (q *queue) addElementToHeader(e interface{}) error {

	if !q.Alive() {
		return cache2.ErrContainerAlreadyClosed
	}
	q.store.PushFront(e)
	return nil
}

// please call me with lock
func (q *queue) addElementToTail(e interface{}) error {

	if !q.Alive() {
		return cache2.ErrContainerAlreadyClosed
	}
	q.store.PushBack(e)
	return nil
}

// please call me with lock
func (q *queue) getElementFromHeader() (interface{}, error) {

	if !q.Alive() {
		return nil, cache2.ErrContainerAlreadyClosed
	}

	e := q.store.Front()
	if e == nil {
		return nil, cache2.ErrContainerEmpty
	}
	return q.store.Remove(e), nil
}

// please call me with lock
func (q *queue) getElementFromTail() (interface{}, error) {

	if !q.Alive() {
		return nil, cache2.ErrContainerAlreadyClosed
	}

	e := q.store.Back()
	if e == nil {
		return nil, cache2.ErrContainerEmpty
	}
	return q.store.Remove(e), nil
}

func (q *queue) EnQueue(e interface{}) error {

	q.lock.Lock()
	defer q.lock.Unlock()

	return q.addElementToTail(e)
}

func (q *queue) EnQueueToFront(e interface{}) error {

	q.lock.Lock()
	defer q.lock.Unlock()

	return q.addElementToHeader(e)

}

func (q *queue) DeQueue() (interface{}, error) {

	q.lock.Lock()
	defer q.lock.Unlock()

	return q.getElementFromHeader()
}

func (q *queue) DeQueueFromTail() (interface{}, error) {

	q.lock.Lock()
	defer q.lock.Unlock()

	return q.getElementFromTail()
}

func NewDeQueue(ctx context.Context) cache2.DoubleEndQueue {

	return &queue{
		LifeCycle: cache.NewLifeCycle(ctx),
		store:     list.New(),
	}
}

type blockQueue struct {
	queue

	cond         *sync.Cond
	burstChecker cache.BurstChecker
}

func NewBlockDeQueue(ctx context.Context, checker cache.BurstChecker) cache2.DoubleEndQueue {

	out := &blockQueue{
		queue: queue{
			LifeCycle: cache.NewLifeCycle(ctx),
			store:     list.New(),
		},
		burstChecker: checker,
	}

	out.cond = sync.NewCond(&out.lock)
	return out
}

func (q *blockQueue) Close() error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.queue.Close()
	q.cond.Broadcast()

	return nil
}

func (q *blockQueue) EnQueue(e interface{}) error {

	q.lock.Lock()
	defer q.lock.Unlock()

	for q.burstChecker(q.store.Len()) && q.Alive() {
		q.cond.Wait()
	}

	if err := q.addElementToTail(e); err != nil {
		return err
	}

	q.cond.Broadcast()
	return nil
}

func (q *blockQueue) EnQueueToFront(e interface{}) error {

	q.lock.Lock()
	defer q.lock.Unlock()

	for q.burstChecker(q.store.Len()) && q.Alive() {
		q.cond.Wait()
	}

	if err := q.addElementToHeader(e); err != nil {
		return err
	}

	q.cond.Broadcast()
	return nil

}

func (q *blockQueue) DeQueue() (interface{}, error) {

	q.lock.Lock()
	defer q.lock.Unlock()

	for q.store.Len() == 0 && q.Alive() {
		q.cond.Wait()
	}

	out, err := q.getElementFromHeader()
	if err != nil {
		return nil, err
	}

	q.cond.Broadcast()
	return out, nil

}

func (q *blockQueue) DeQueueFromTail() (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for q.store.Len() == 0 && q.Alive() {
		q.cond.Wait()
	}

	out, err := q.getElementFromTail()
	if err != nil {
		return nil, err
	}

	q.cond.Broadcast()
	return out, nil
}
