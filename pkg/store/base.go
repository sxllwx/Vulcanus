package store

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

// NiceContainer
// when the container be closed, then notify the done
// when the container be closed, the goroutines can wait the container stop
type NiceContainer interface {
	io.Closer
	Done() <-chan struct{}

	// every goroutine can got the rest element from the stopped store
	Rest() ([]interface{}, error)

	// read current length of the container
	Len() (int, error)
}

// BlockContainer
// when the container is full, the current goroutine insert operation will block until the space free
// when the container is empty, the current goroutine pop operation will block until the element was be added
type BlockContainer interface {
	BlockRequest() // just a flag method
}

type BurstChecker func(int) bool

type LifeCycle struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// if the parent be closed, the LifeCycle ctx.Done() also be closed
func NewLifeCycle(parent context.Context) LifeCycle {

	ctx, cancel := context.WithCancel(parent)
	return LifeCycle{
		ctx:    ctx,
		cancel: cancel,
	}

}

func (lc *LifeCycle) Done() <-chan struct{} {
	return lc.ctx.Done()
}

func (lc *LifeCycle) Close() error {
	lc.cancel()
	return nil
}

func (lc *LifeCycle) Alive() bool {

	select {
	case <-lc.Done():
		return false
	default:
		return true
	}
}
