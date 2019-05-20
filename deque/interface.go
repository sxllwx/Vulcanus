package deque

// Interface of double-end-queue
type Interface interface {
	// Serialize interface to persistence storage
	Serializer

	// current length of queue
	Len() (int, error)

	Insert(interface{}, InsertDirection) error
	Out(OutDirection) (interface{}, error)

	// get object but return too late
	// then put back to deque
	Revert(interface{}) error

	// get all queue element from queue
	Empty() ([]interface{}, error)

	// ack the object
	Done(interface{}) error
}

type Serializer interface {
	Decoder
	Encoder
}

type Decoder interface {
	// decode from disk
	Decode([]byte) error
}

type Encoder interface {
	// persistence to disk
	Encode() ([]byte, error)
}

type InsertDirection func([]interface{}, interface{}) []interface{}

type OutDirection func([]interface{}) (interface{}, []interface{})

var InsertToTail = InsertDirection(func(queue []interface{}, o interface{}) []interface{} {
	return append(queue, o)
})

var InsertToHeader = InsertDirection(func(queue []interface{}, o interface{}) []interface{} {
	return append(append([]interface{}{}, o), queue...)
})

var OutFromHeader = OutDirection(func(queue []interface{}) (interface{}, []interface{}) {
	return queue[0], queue[1:]
})

var OutFromTail = OutDirection(func(queue []interface{}) (interface{}, []interface{}) {
	return queue[len(queue)-1], queue[:len(queue)-1]
})
