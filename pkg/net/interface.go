package net

import (
	"context"
	"io"
	"sync/atomic"
	"time"
)

type measurable interface {
	Total() uint64
	BPS() uint64
	Cost() time.Duration
}

type defaultMeasurableSuite struct {

	// lifecycle manager
	ctx    context.Context
	cancel context.CancelFunc

	// metric value
	totalBytes uint64
	bps        uint64

	start time.Time

	ticker *time.Ticker
}

func (r *defaultMeasurableSuite) loop() {

	for {

		oldTotalBytes := atomic.LoadUint64(&r.totalBytes)

		select {
		case <-r.ctx.Done():
			r.ticker.Stop()
			return
		case <-r.ticker.C:
		}

		// store the new rate
		atomic.StoreUint64(&r.bps, atomic.LoadUint64(&r.totalBytes)-oldTotalBytes)
	}

}

func (r *defaultMeasurableSuite) Cost() time.Duration {
	return time.Since(r.start)
}

func (r *defaultMeasurableSuite) Total() uint64 {
	return atomic.LoadUint64(&r.totalBytes)
}

func (r *defaultMeasurableSuite) BPS() uint64 {
	return atomic.LoadUint64(&r.bps)
}

func newMeasurableSuite() *defaultMeasurableSuite {

	ctx, cancel := context.WithCancel(context.Background())
	out := &defaultMeasurableSuite{
		ctx:    ctx,
		cancel: cancel,
		ticker: time.NewTicker(time.Second),
		start:  time.Now(),
	}

	out.loop()
	return out
}

type MeasurableWriteCloser interface {
	measurable
	io.WriteCloser
}

type MeasurableReadCloser interface {
	measurable
	io.ReadCloser
}

type readCloser struct {
	*defaultMeasurableSuite
	io.ReadCloser
}

func DecorateReadCloser(rc io.ReadCloser) MeasurableReadCloser {
	return &readCloser{
		defaultMeasurableSuite: newMeasurableSuite(),
		ReadCloser:             rc,
	}
}

func (rc *readCloser) Close() error {
	rc.cancel()
	return rc.ReadCloser.Close()
}
func (rc *readCloser) Read(b []byte) (int, error) {

	n, err := rc.ReadCloser.Read(b)
	if err != nil {
		return 0, err
	}

	atomic.AddUint64(&rc.totalBytes, uint64(n))
	return n, nil
}

type writeCloser struct {
	*defaultMeasurableSuite
	io.WriteCloser
}

func DecorateWriteCloser(wc io.WriteCloser) MeasurableWriteCloser {
	return &writeCloser{
		defaultMeasurableSuite: newMeasurableSuite(),
		WriteCloser:            wc,
	}
}

func (wc *writeCloser) Write(b []byte) (int, error) {

	n, err := wc.WriteCloser.Write(b)
	if err != nil {
		return 0, err
	}

	atomic.AddUint64(&wc.totalBytes, uint64(n))
	return n, nil
}

func (wc *writeCloser) Close() error {
	wc.cancel()
	return wc.WriteCloser.Close()
}
