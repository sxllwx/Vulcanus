package store

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrStoreAlreadyStopped = errors.New("the store already stopped")
	ErrNoMoreElements      = errors.New("the store no more elements")
)

// Simple
// this is simple store
type Simple interface {

	// put the element to store
	Put(interface{}) error
	// get the  first element from store
	Get() (interface{}, error)

	// get the size of elements in store
	// if current store elements > size, get the size elements from store
	// if current store elements <= current, get all  elements in store
	Batch(size int) ([]interface{}, error)

	// close the store, get the rest elements in store
	Close() ([]interface{}, error)

	Done() <-chan struct{}
}

type options struct {
	ctx context.Context
}

type Option func(*options)

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

type simpleImpl struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu    sync.Mutex // protect the store
	store []interface{}
}

func (s *simpleImpl) running() bool {

	select {
	case <-s.ctx.Done():
		return false
	default:
		return true
	}

}

func (s *simpleImpl) Put(o interface{}) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running() {
		return ErrStoreAlreadyStopped
	}
	s.store = append(s.store, o)
	return nil
}

func (s *simpleImpl) Get() (interface{}, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running() {
		return nil, ErrStoreAlreadyStopped
	}

	if len(s.store) == 0 {
		return nil, ErrNoMoreElements
	}

	out := s.store[0]
	s.store = s.store[1:]

	return out, nil
}

func (s *simpleImpl) Batch(size int) ([]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running() {
		return nil, ErrStoreAlreadyStopped
	}

	if len(s.store) == 0 {
		return nil, ErrNoMoreElements
	}

	switch {
	case len(s.store) <= size:
		out := s.store
		s.store = nil
		return out, nil

	case len(s.store) > size:
		out := s.store[:size]
		s.store = s.store[size:]
		return out, nil
	}
	panic("unreachable")
}

func (s *simpleImpl) Close() ([]interface{}, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running() {
		return nil, ErrStoreAlreadyStopped
	}

	out := s.store

	// cancel ctx
	s.cancel()

	return out, nil
}

func (s *simpleImpl) Done() <-chan struct{} {
	return s.ctx.Done()
}

func NewSimpleStore(ctx context.Context) Simple {

	ctx, cancel := context.WithCancel(ctx)
	out := &simpleImpl{
		ctx:    ctx,
		cancel: cancel,
		store:  []interface{}{},
	}

	return out
}
