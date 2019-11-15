package queue

import (
	"errors"
)

var (
	ErrAlreadyStopped = errors.New("queue already stop")
)

// this is a thread-safe queue interface
// multi producer can safe pull the object to store
// multi consumer can safe get the object from store
//  p1  p2  p3 ....
//  |   |   |
//   \  |  /
//    \ | /
//   Interface
//      |
//     /|\
//    / | \
//   /  |  \
//  c1  c2 c3  ....
type Interface interface {

	// insert the element to tail
	// error always not nil, when queue already stopped
	Append(interface{}) error

	// insert the element to header
	// error always not nil, when queue already stopped
	Insert(interface{}) error

	// NOTICE:
	// Get will got the object from header
	// if there no element in queue, it'will block
	// error always not nil, when queue already stopped
	Get() (interface{}, error)

	// Close the queue and got the rest object
	// error always not nil, when queue already stopped
	Close() ([]interface{}, error)

	// the consumer and producer should put/get element on the queue is running
	Done() <-chan struct{}
}

// no limit the queue len
var notLimit = func(objectList []interface{}) bool {
	return false
}
