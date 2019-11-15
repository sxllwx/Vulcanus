package simple

// the simple queue
// producer and consumer are both thread-safe
type Interface interface {

	// insert the element to the header of queue
	Insert(interface{}) error
	// append the element to the tail of queue
	Append(interface{}) error
	// pop the element from the queue
	Pop() (interface{}, error)
	// close the queue, and get the reset element from the queue
	Close() ([]interface{}, error)
	// the queue status
	Done() <-chan struct{}
}
