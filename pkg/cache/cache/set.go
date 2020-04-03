package cache

import (
	"context"
	"sync"

	"github.com/sxllwx/vulcanus/pkg/cache"
)

type set struct {
	cache.LifeCycle

	lock  sync.RWMutex
	store map[interface{}]struct{}
}

func (s *set) Wait() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return nil
}

func (s *set) Rest() ([]interface{}, error) {

	if s.Alive() {
		return nil, cache.ErrRestShouldNotBeCall
	}

	var out []interface{}
	for k, _ := range s.store {
		out = append(out, k)
	}
	return out, nil
}

func (s *set) Len() (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.store), nil
}

// please call me with lock
func (s *set) addElement(e interface{}) error {

	if !s.Alive() {
		return cache.ErrContainerAlreadyClosed
	}

	s.store[e] = struct{}{}
	return nil
}

// please call me with lock
func (s *set) getElement() (interface{}, error) {

	if !s.Alive() {
		return nil, cache.ErrContainerAlreadyClosed
	}

	if len(s.store) == 0 {
		return nil, cache.ErrContainerEmpty
	}

	for k, _ := range s.store {
		delete(s.store, k)
		return k, nil
	}
	panic("unreachable")
}

func (s *set) Put(e interface{}) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	return s.addElement(e)
}

func (s *set) Get() (interface{}, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	return s.getElement()
}

func (s *set) List() ([]interface{}, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.Alive() {
		return nil, cache.ErrContainerAlreadyClosed
	}

	var out []interface{}
	for k, _ := range s.store {
		out = append(out, k)
	}

	return out, nil
}

func (s *set) Close() error {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.LifeCycle.Close()
	return nil
}

func NewSet(parent context.Context) cache.Set {

	return &set{
		LifeCycle: cache.NewLifeCycle(parent),
		store:     map[interface{}]struct{}{},
	}
}

func NewBlockSet(parent context.Context, check cache.BurstChecker) cache.Set {

	out := &blockSet{

		set: set{

			LifeCycle: cache.NewLifeCycle(parent),
			store:     map[interface{}]struct{}{},
		},
		checker: check,
	}

	out.cond = sync.NewCond(&out.lock)
	return out
}

type blockSet struct {
	set

	cond    *sync.Cond
	checker cache.BurstChecker
}

func (s *blockSet) Put(e interface{}) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	for s.checker(len(s.store)) && s.Alive() {
		s.cond.Wait()
	}

	if err := s.addElement(e); err != nil {
		return err
	}
	s.cond.Broadcast()
	return nil
}

func (s *blockSet) Get() (interface{}, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	for len(s.store) == 0 && s.Alive() {
		s.cond.Wait()
	}

	e, err := s.getElement()
	if err != nil {
		return nil, err
	}

	s.cond.Broadcast()
	return e, nil
}

func (s *blockSet) Close() error {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.LifeCycle.Close()
	s.cond.Broadcast()
	return nil
}
