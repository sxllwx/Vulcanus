package storage

import (
	"context"
)

type Interface interface {
	// put the k/v
	Put(key string, data []byte) error
	// delete the k
	Delete(key string) error
	// get the spec key value
	Get(key string) ([]byte, error)
	// watch the spec key value
	Watch(key string) (<-chan *Message, error)
	// get the object list
	// key start with prefix
	List(prefix string) (map[string][]byte, error)

	Reset() error
	// Close the storage
	Close() error
	// lock the database
	Locker(k string) (Locker, error)
}

type Locker interface {
	Lock(context.Context) error
	Unlock(context.Context) error
}

type EventType string

const (
	Create EventType = "create"
	Delete EventType = "delete"
	Update EventType = "update"
)

type Message struct {
	EventType EventType
	Data      []byte
}
