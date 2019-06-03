package deque

import (
	"sync"
)

type StopErr string

func (s StopErr) Error() string {
	return "double end queue already stop"
}

type processQueue []interface{}

func (s *processQueue) add(o interface{}) {

	if s.has(o) {
		return
	}
	*s = append(*s, o)
}

func (s *processQueue) delete(o interface{}) {
	for i, object := range *s {
		if object == o {
			*s = append((*s)[:i], (*s)[i:])
		}
	}
}

func (s processQueue) has(o interface{}) bool {
	for _, object := range s {
		if object == o {
			return true
		}
	}
	return false
}

// simpleDeque provider a simple double end queue with follow features
// * Fair: item can be process in order
// * Stingy: a item just will be handle by on consumer
// * No-Repeat
// * Persistence: every operation will trigger flush operation
type simpleDeque struct {

	// 1. protect queue
	// 2. per op per writer
	cond sync.Cond

	queue []interface{}

	processing processQueue

	persistenceOperation func()

	running bool
}

func New(persistenceOperation func()) Interface {

	q := &simpleDeque{
		running:              true,
		persistenceOperation: persistenceOperation,
		processing:           processQueue{},
	}
	q.cond.L = &sync.Mutex{}
	return q
}

func (q *simpleDeque) Stop() {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.running = false
	q.persistenceOperation()

	q.cond.Broadcast()
}

func (q *simpleDeque) Len() (int, error) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if err := q.check(); err != nil {
		return 0, err
	}
	return len(q.queue), nil
}

func (q *simpleDeque) Revert(o interface{}) error {

	// put back to queue
	if err := q.Insert(o, InsertToHeader); err != nil {
		return err
	}

	// delete object from process
	if err := q.Done(o); err != nil {
		return err
	}
	return nil
}

func (q *simpleDeque) check() error {

	if q.running {
		return nil
	}
	return StopErr("")
}

func (q *simpleDeque) Insert(o interface{}, insertDirection InsertDirection) error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if err := q.check(); err != nil {
		return err
	}

	// check dup
	q.queue = insertDirection(q.queue, o)
	q.cond.Broadcast()
	return nil
}

// Out
// ** This method will block until Insert a new object
func (q *simpleDeque) Out(outDirection OutDirection) (interface{}, error) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.queue) == 0 && q.running {
		q.cond.Wait()
	}

	if err := q.check(); err != nil {
		return nil, err
	}

	var o interface{}
	o, q.queue = outDirection(q.queue)
	q.processing.add(o)
	q.cond.Broadcast()

	return o, nil
}

func (q *simpleDeque) Done(o interface{}) error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if err := q.check(); err != nil {
		return err
	}

	q.processing.delete(o)
	return nil
}
