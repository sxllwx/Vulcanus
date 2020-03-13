package cache

import (
	"sync"
	"time"

	"github.com/sxllwx/vulcanus/pkg/store"
)

type entry struct {
	data  []byte
	timer *time.Timer
}

type kvStore struct {
	store.LifeCycle

	lock  sync.RWMutex
	store map[string]*entry
}

func (s *kvStore) Put(k string, v []byte) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.Alive() {
		return store.ErrContainerAlreadyClosed
	}

	s.store[k] = &entry{
		data: v,
	}
	return nil
}

func (s *kvStore) Get(k string) ([]byte, error) {

	s.lock.RLock()
	defer s.lock.RUnlock()

	if !s.Alive() {
		return nil, store.ErrContainerAlreadyClosed
	}
	out, ok := s.store[k]
	if !ok {
		return nil, store.ErrNotFound
	}
	return out.data, nil
}

func (s *kvStore) Delete(k string) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.Alive() {
		return store.ErrContainerAlreadyClosed
	}

	delete(s.store, k)
	return nil
}

func (s *kvStore) ListKeys() ([]string, error) {

	s.lock.RLock()
	defer s.lock.RUnlock()

	var out []string
	for k, _ := range s.store {
		out = append(out, k)
	}

	return out, nil
}

func (s *kvStore) TTL(k string, ttl time.Duration) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	out, ok := s.store[k]
	if !ok {
		return store.ErrNotFound
	}

	// stop old timer
	if out.timer != nil {
		out.timer.Stop()
		out.timer = nil
	}

	// set new timer
	out.timer = time.AfterFunc(ttl, func() {

		s.lock.Lock()
		defer s.lock.Unlock()

		entity, ok := s.store[k]
		if !ok {
			return
		}

		if entity.timer == nil {
			return
		}

		entity.timer.Stop()
		delete(s.store, k)
	})

	return nil
}

func (s *kvStore) Rest() ([]interface{}, error) {

	if s.Alive() {
		return nil, store.ErrRestShouldNotBeCall
	}

	// TODO
	return nil, nil
}

func (s *kvStore) Len() (int, error) {

	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.store), nil
}
