package cache

import (
	"container/list"
	"context"
	cache2 "github.com/sxllwx/vulcanus/pkg/feature/cache"
	"sync"

	"github.com/sxllwx/vulcanus/pkg/feature/cachere/cache"
)

type stack struct {

	// manage the  stack lifecycle
	cache.LifeCycle

	lock  sync.RWMutex
	store *list.List
}

func (s *stack) Len() (int, error) {

	s.lock.RLock()
	out := s.store.Len()
	s.lock.RUnlock()
	return out, nil
}

func (s *stack) Rest() ([]interface{}, error) {

	if s.Alive() {
		return nil, cache2.ErrRestShouldNotBeCall
	}

	// already stopped
	out := make([]interface{}, 0, 8) // !!!!
	for i := s.store.Front(); i != nil; i = i.Next() {
		out = append(out, i.Value)
	}
	return out, nil
}

func (s *stack) Push(e interface{}) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	return s.pushElement(e)
}

// please call me with a lock
func (s *stack) pushElement(e interface{}) error {

	if !s.Alive() {
		return cache2.ErrContainerAlreadyClosed
	}
	s.store.PushBack(e)
	return nil
}

// please call me with a lock
func (s *stack) popElement() (interface{}, error) {

	if !s.Alive() {
		return nil, cache2.ErrContainerAlreadyClosed
	}

	out := s.store.Back()
	if out != nil {
		return s.store.Remove(out), nil
	}
	return nil, cache2.ErrContainerEmpty
}

func (s *stack) Pop() (interface{}, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	return s.popElement()
}

func NewStack(ctx context.Context) cache2.Stack {

	return &stack{
		LifeCycle: cache.NewLifeCycle(ctx),
		store:     list.New(),
	}
}

type blockStack struct {
	stack

	cond         *sync.Cond
	burstChecker cache.BurstChecker
}

func NewBlockStack(ctx context.Context, checker cache.BurstChecker) cache2.Stack {

	out := &blockStack{
		stack: stack{
			LifeCycle: cache.NewLifeCycle(ctx),
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

	// burst check
	for s.burstChecker(s.store.Len()) && s.Alive() {
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

	for s.store.Len() == 0 && s.Alive() {
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

func (s *blockStack) Close() error {

	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.stack.Close()
	s.cond.Broadcast()
	return nil
}
