package cache

import "sync"

// KV
type KV interface {
	// basic feature
	Interface
	// put a [k: v] to cache
	Put(k string, v interface{}) error
	// get object from store by k
	Get(k string) (interface{}, error)
	// delete the element from the store
	Delete(k string) error
	// watch the object
	Watch(string) <-chan Event
	// list all key from the store
	ListKeys() ([]string, error)
}

// k/v cache
type kvCache struct {
	mu    sync.RWMutex
	store map[string]interface{}
}
