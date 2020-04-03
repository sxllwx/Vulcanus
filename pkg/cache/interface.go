package cache

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

var (
	ErrContainerEmpty         = errors.New("container is empty")
	ErrContainerAlreadyClosed = errors.New("container already be closed")
	ErrRestShouldNotBeCall    = errors.New("the Rest method should be call after Close")
	ErrNotFound               = errors.New("the element not be found in container")
)

// Interface
// when the store be closed, then notify the done
// when the store be closed, the goroutines can wait the container stop
type Interface interface {

	// can be closed
	io.Closer

	// stop notify chan
	Done() <-chan struct{}

	// every goroutine can got the rest element from the stopped store
	Rest() ([]interface{}, error)

	// read current length of the container
	Len() (int, error)
}

type lifeCycle struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// if the parent be closed, the lifeCycle ctx.Done() also be closed
func newLifeCycle(parent context.Context) lifeCycle {

	ctx, cancel := context.WithCancel(parent)
	return lifeCycle{
		ctx:    ctx,
		cancel: cancel,
	}

}

func (lc *lifeCycle) Done() <-chan struct{} {
	return lc.ctx.Done()
}

func (lc *lifeCycle) Close() error {
	lc.cancel()
	return nil
}

// if the cache is alive
func (lc *lifeCycle) ok() bool {
	select {
	case <-lc.Done():
		return false
	default:
		return true
	}
}

type eventType string

const (
	// the event-type
	Update eventType = "update"
	Create eventType = "create"
	Delete eventType = "delete"
)

type Event struct {
	Object interface{}
	Event  eventType
}
