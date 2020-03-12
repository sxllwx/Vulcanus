package store

// Queue
// queue is FIFO
type Queue interface {
	NiceContainer

	// enqueue, add the element to the tail
	// if the queue is full,
	// the actually behavior will be decided by impl
	EnQueue(interface{}) error

	// pop the element from the header
	// if the queue is empty,
	// the actually behavior will be decided by impl
	DeQueue() (interface{}, error)
}

// DoubleEndQueue
// can be get element from tail
// can be push element to header
type DoubleEndQueue interface {
	Queue
	// insert element to queue header
	EnQueueToFront(interface{}) error
	//  pop the element from the queue tail
	DeQueueFromTail() (interface{}, error)
}
