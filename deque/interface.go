package deque

// Interface of double-end-queue
type Interface interface {
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
