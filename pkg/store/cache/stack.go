package cache

import (
	"container/list"
	"context"
	"sync"

	"github.com/sxllwx/vulcanus/pkg/store"
)

type stack struct {

	// manage the  stack lifecycle
	store.LifeCycle

	lock  sync.RWMutex
	wg    sync.WaitGroup // add the wg for every write operation
	store *list.List
}

func (s *stack) Close() error {
	return s.LifeCycle.Close()
}

func (s *stack) Done() <-chan struct{} {
	return s.LifeCycle.Done()
}

func (s *stack) Wait() error {
	s.wg.Wait()
	return nil
}

func (s *stack) Rest() ([]interface{}, error) {

	if s.Alive() {
		return nil, store.ErrRestShouldNotBeCall
	}

	// already stopped
	out := make([]interface{}, 0, 8) // !!!!
	for i := s.store.Front(); i != nil; i = i.Next() {
		out = append(out, i.Value)
	}
	return out, nil
}

func (s *stack) Len() (int, error) {

	s.lock.RLock()
	out := s.store.Len()
	s.lock.RUnlock()
	return out, nil
}

func (s *stack) Push(e interface{}) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.wg.Add(1)
	defer s.wg.Done()

	return s.pushElement(e)
}

// please call me with a lock
func (s *stack) pushElement(e interface{}) error {

	if !s.Alive() {
		return store.ErrContainerAlreadyStopped
	}
	s.store.PushBack(e)
	return nil
}

// please call me with a lock
func (s *stack) popElement() (interface{}, error) {

	if !s.Alive() {
		return nil, store.ErrContainerAlreadyStopped
	}

	out := s.store.Back()
	if out != nil {
		s.store.Remove(out)
		return out.Value, nil
	}
	return nil, store.ErrContainerEmpty
}

func (s *stack) Pop() (interface{}, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.wg.Add(1)
	defer s.wg.Done()

	return s.popElement()
}

func NewStack(ctx context.Context) store.Stack {

	return &stack{
		LifeCycle: store.NewLifeCycle(ctx),
		store:     list.New(),
	}
}

type blockStack struct {
	cond *sync.Cond
	stack

	burstChecker store.BurstChecker
}

func NewBlockStack(ctx context.Context, checker store.BurstChecker) store.Stack {

	out := &blockStack{
		stack: stack{
			LifeCycle: store.NewLifeCycle(ctx),
			store:     list.New(),
		},
		burstChecker: checker,
	}
	out.cond = sync.NewCond(&out.lock)
	return out
}

func (s *blockStack) Push(e interface{}) error {

	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.wg.Add(1)
	defer s.wg.Done()

	// burst check
	for s.burstChecker(s.store.Len()) {
		// full, wait empty
		s.cond.Wait()
	}

	if err := s.pushElement(e); err != nil {
		return err
	}

	// notify there are a new element
	s.cond.Broadcast()
	return nil
}

func (s *blockStack) Pop() (interface{}, error) {

	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.wg.Add(1)
	defer s.wg.Done()

	for s.store.Len() == 0 {
		// empty
		s.cond.Wait()
	}

	e, err := s.popElement()
	if err != nil {
		return nil, err
	}

	s.cond.Broadcast()
	return e, nil
}
