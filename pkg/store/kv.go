package store

import (
	"time"
)

type KV interface {

	// put the k/v to store
	// put will update the exist object
	Put(string, []byte) error
	// get the k from the store
	Get(string) ([]byte, error)
	// delete the element from the store
	Delete(string) error
	// list all key from the store
	ListKeys() ([]string, error)
	// bind the ttl for a key
	TTL(string, time.Duration) error

	NiceContainer
}
