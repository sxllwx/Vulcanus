package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/juju/errors"
	"github.com/sxllwx/vulcanus/pkg/cache/store"
)

var (
	ErrConsumerAlreadyStopped = errors.New("the consumer already stopped")
	ErrConsumerLowEnergy      = errors.New("the consumer low energy")
)

type RateLimiterBatchConsumer interface {
	Run()
	MaxHandlingCapOneBatch() int
	Done() <-chan struct{}
	Close() ([]interface{}, error)
}

type supportedStore interface {
	Put(interface{}) error
	Get() (interface{}, error)
	Batch(size int) ([]interface{}, error)
	Done() <-chan struct{}
}

type rateLimiterBatchConsumerImpl struct {

	// manage the consumer lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	mu                     sync.RWMutex
	maxHandlingCapOneBatch int
	dirty                  []interface{} // the handling
	wg                     sync.WaitGroup

	// the store
	store supportedStore

	// the consumer handle logic
	HandleInterval time.Duration
	HandleFunc     func([]interface{}) error
	BurstCheckFunc func([]interface{}) (bool, error)
}

func (c *rateLimiterBatchConsumerImpl) running() bool {

	select {
	case <-c.Done():
		return false
	default:
		return true
	}
}

func (c *rateLimiterBatchConsumerImpl) Run() {

	ticker := time.NewTicker(c.HandleInterval)
	defer ticker.Stop()

	for {
		select {
		// the store already be stopped
		case <-c.store.Done():
			return
		// the consumer be stopped
		case <-c.Done():
			return

		case <-ticker.C:
			c.wg.Add(1)

			go func() {

				defer c.wg.Done()
				out, err := c.getElements()
				if err != nil {
					fmt.Println("find thresh-hold", err)
					return
				}

				// no elements wait for handle
				if len(out) == 0 {
					return
				}

				if err := c.HandleFunc(out); err != nil {
					c.markAsDirty(out)
					fmt.Println("handle elements", err)
				}
			}()

		}
	}
}

func (c *rateLimiterBatchConsumerImpl) getElements() ([]interface{}, error) {

	// force find thresh hold
	bucket, err := c.findThreshHold()
	if err != nil {
		c.markAsDirty(bucket)
		return nil, errors.Annotate(err, "find thresh hold")
	}

	if len(bucket) == 0 {
		return nil, nil
	}

	return bucket, nil
}

func (c *rateLimiterBatchConsumerImpl) updateMaxHandlingCapOneBatch(l int) {

	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxHandlingCapOneBatch = l
}

func (c *rateLimiterBatchConsumerImpl) findThreshHold() ([]interface{}, error) {

	var (
		// init wind
		wind = 1
		// bucket used to check size
		bucket []interface{}
		// packaged is part of got elements, which not overflow
		packaged []interface{}
		// tail is part of  got elements, which already overflow
		rest []interface{}
	)

	// STEP.1 get win size elements
	// STEP.2 append to bucket, then burst-check
	//         not burst: {put wind <<= 1 && put got elements to packaged} goto STEP.1
	//         burst: the rest element is got elements
	for {

		got, err := c.store.Batch(wind)
		if err != nil && err != store.ErrNoMoreElements {
			// the store was stopped!!
			// return the bucket to caller, avoid message lose
			return bucket, errors.Annotatef(err, "get %d elements from store", wind)
		}

		if err == store.ErrNoMoreElements {
			// the store empty
			break
		}

		bucket = append(bucket, got...)
		burst, err := c.BurstCheckFunc(bucket)
		if err != nil {
			// the handleThreshHoldCheckFunc have bug
			// return the bucket to caller, avoid message lose
			return bucket, errors.Annotate(err, "run handle thresh-hold check func")
		}

		if !burst {
			wind <<= 1
			// not burst
			packaged = append(packaged, got...)
			continue
		}
		// the got can't be consume by packaged
		rest = got
		break
	}

	var tailThreshHold = 0
	err := c.evalTailThreshHold(&tailThreshHold, packaged, rest)
	if err != nil {
		// return all bucket info
		return bucket, errors.Annotate(err, "eval the tail thresh-hold")
	}

	packaged = append(packaged, rest[:tailThreshHold]...)

	// revert the elements to store
	for i, e := range rest[tailThreshHold:] {

		err := c.store.Put(e)
		if err != nil {
			return append(packaged, rest[i:]), errors.Annotate(err, "revert element")
		}
	}

	return packaged, nil
}

func (c *rateLimiterBatchConsumerImpl) evalTailThreshHold(
	tailThreshHold *int,
	packaged []interface{},
	rest []interface{},
) error {

	var (
		// init wind
		wind = 1
		// the bucket should be eval burst
		bucket = packaged
	)

	// STEP.1 get the wind size elements from rest,
	// STEP.2 append the elements to bucket
	// STEP.3 burst check
	//          not burst {put wind <<=1 && put got elements to bucket && increase the tailThreshHold} goto STEP.1
	//          not pass goto STEP.4
	// STEP.4 recurse call this method
	for {

		var got []interface{}

		switch {
		case len(rest) == 0:
			// the rest empty
			return nil
		case wind >= len(rest):
			// the rest is already not enough
			got, rest = rest, nil
		default:
			// rich rest,direct slice it
			got, rest = rest[:wind], rest[wind:]
		}

		bucket = append(bucket, got...)
		burst, err := c.BurstCheckFunc(bucket)
		if err != nil {
			return errors.Annotate(err, "run handle thresh-hold check func")
		}

		if !burst {
			// not bust
			// larger wind
			// add got elements to packaged
			// fix tailThreshHold
			wind <<= 1
			*tailThreshHold += len(got)
			packaged = append(packaged, got...)
			continue
		}

		// already burst
		if wind == 1 {
			// can't hold any elements
			return nil
		}

		// the got elements can't be consume
		err = c.evalTailThreshHold(tailThreshHold, packaged, got)
		if err != nil {
			return errors.Annotate(err, "eval tail thresh hold")
		}
		return nil
	}
}

func (c *rateLimiterBatchConsumerImpl) markAsDirty(bucket []interface{}) {

	if len(bucket) == 0 {
		return
	}

	c.mu.Lock()
	c.dirty = append(c.dirty, bucket...)
	c.mu.Unlock()
}

func (c *rateLimiterBatchConsumerImpl) MaxHandlingCapOneBatch() int {

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.maxHandlingCapOneBatch
}

func (c *rateLimiterBatchConsumerImpl) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *rateLimiterBatchConsumerImpl) Close() ([]interface{}, error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running() {
		return nil, ErrConsumerAlreadyStopped
	}

	c.cancel()
	c.wg.Wait()
	return c.dirty, nil
}

type options struct {
	ctx   context.Context
	store supportedStore
}

type Option func(*options)

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithStore(store supportedStore) Option {
	return func(o *options) {
		o.store = store
	}
}

func NewLimiterBatchConsumer(opts ...Option) *rateLimiterBatchConsumerImpl {

	// default setting
	o := &options{
		ctx: context.Background(),
	}

	// apply the opts to queueOptions
	for _, f := range opts {
		f(o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	out := &rateLimiterBatchConsumerImpl{
		ctx:    ctx,
		cancel: cancel,
		store:  o.store,
	}

	return out
}
