package deque

import (
	"io"
)

// Interface of double-end-queue
type Interface interface {
	// Serialize interface to persistence storage
	Serializer

	// current length of queue
	Len() (int, error)

	// push object to queue tail
	Push(interface{}) error
	// get object to queue header
	Shift() (interface{}, error)
	// get object from queue tail
	Pop() (interface{}, error)
	// ack the object
	Done(interface{}) error
}

type Serializer interface {
	Decoder
	Encoder
}

type Decoder interface {
	// decode from disk
	Decode(io.ReadCloser) error
}

type Encoder interface {
	// persistence to disk
	Encode(io.WriteCloser) error
}
