package cache

import (
	"container/list"
	"sync"

	"github.com/sxllwx/vulcanus/pkg/store"
)

type stack struct {

	// manage the  stack lifecycle
	store.LifeCycle

	lock sync.RWMutex
	wg   sync.WaitGroup // add the wg for every write operation

	store *list.List
}

func (s *stack) Safe() {}

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
	out := make([]interface{}, 8)
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

	if !s.Alive() {
		return store.ErrContainerAlreadyStopped
	}

	s.store.PushBack(e)
	return nil
}

func (s *stack) Pop() (interface{}, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.Alive() {
		return nil, store.ErrContainerAlreadyStopped
	}

	return s.store.Remove(s.store.Back()), nil
}
