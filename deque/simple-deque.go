package deque

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type empty struct{}
type set map[interface{}]empty
type StopErr string

func (s StopErr) Error() string {
	return "double end queue already stop"
}

func (s set) add(o interface{}) {
	s[o] = empty{}
}

func (s set) has(o interface{}) bool {
	_, ok := s[o]
	return ok
}

func (s set) delete(o interface{}) {
	if _, ok := s[o]; !ok {
		return
	}
	delete(s, o)
}

// simpleDeque provider a simple double end queue with follow features
// * Fair: item can be process in order
// * Stingy: a item just will be handle by on consumer
// * No-Repeat
type simpleDeque struct {

	// protect queue
	cond sync.Cond

	queue []interface{}

	// need to process
	// distinct the object
	dirty set

	// processing
	processing set

	stop <-chan struct{}
}

func New(stop <-chan struct{}) Interface {

	q := &simpleDeque{
		stop: stop,
	}
	q.cond.L = &sync.Mutex{}
	return q
}

func (q *simpleDeque) Len() (int, error) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if err := q.check(); err != nil {
		return 0, err
	}
	return len(q.queue), nil
}

type Data struct {
	Q []interface{} `json:"Q"`
	D []interface{} `json:"D"`
	P []interface{} `json:"P"`
}

func (q *simpleDeque) Decode(readCloser io.ReadCloser) error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	var d Data
	if err := json.NewDecoder(readCloser).Decode(&d); err != nil {
		return fmt.Errorf("decode: %v", err)
	}

	q.queue = append(q.queue, d.Q)

	for _, d := range d.D {
		q.dirty[d] = empty{}
	}
	for _, p := range d.P {
		q.processing[p] = empty{}
	}

	return nil
}

func (q *simpleDeque) Encode(writeCloser io.WriteCloser) error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	d := Data{
		Q: q.queue,
	}

	for o, _ := range q.dirty {
		d.D = append(d.D, o)
	}
	for o, _ := range q.processing {
		d.P = append(d.P, o)
	}

	return json.NewEncoder(writeCloser).Encode(d)
}

func (q *simpleDeque) Push(o interface{}) error {
	return q.insert(o, insertToTail)
}

func (q *simpleDeque) Shift() (interface{}, error) {
	return q.out(outFromHeader)
}

func (q *simpleDeque) Pop() (interface{}, error) {
	return q.out(outFromTail)
}

func (q *simpleDeque) Done(o interface{}) error {
	return q.ack(o)
}

func (q *simpleDeque) check() error {

	select {
	case <-q.stop:
		return StopErr("")
	default:
	}
	return nil
}

func (q *simpleDeque) insert(o interface{}, specDirectionInsertFunc func([]interface{}, interface{}) []interface{}) error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if err := q.check(); err != nil {
		return err
	}

	// check dup
	if q.dirty.has(o) || q.processing.has(o) {
		return nil
	}
	q.dirty.add(o)

	q.queue = specDirectionInsertFunc(q.queue, o)
	q.cond.Signal()

	return nil
}

func (q *simpleDeque) out(specDirectionOutFunc func([]interface{}) (interface{}, []interface{})) (interface{}, error) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if err := q.check(); err != nil {
		return nil, err
	}

	for len(q.queue) == 0 {
		q.cond.Wait()
	}

	var o interface{}
	o, q.queue = specDirectionOutFunc(q.queue)

	q.processing.add(o)
	q.dirty.delete(o)
	return o, nil
}

func (q *simpleDeque) ack(o interface{}) error {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if err := q.check(); err != nil {
		return err
	}

	q.processing.delete(o)
	return nil
}

func insertToTail(queue []interface{}, o interface{}) []interface{} {
	return append(queue, o)
}

func insertToHeader(queue []interface{}, o interface{}) []interface{} {
	return append(append([]interface{}{}, o), queue...)
}

func outFromHeader(queue []interface{}) (interface{}, []interface{}) {
	return queue[0], queue[1:]
}
func outFromTail(queue []interface{}) (interface{}, []interface{}) {
	return queue[len(queue)-1], queue[0 : len(queue)-1]
}
