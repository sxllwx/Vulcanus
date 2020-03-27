package net

import (
	"context"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type measurable interface {
	Total() uint64
	BPS() uint64
	Cost() time.Duration
	addTotal(uint64)
	stop()
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

func (r *defaultMeasurableSuite) addTotal(t uint64) {
	atomic.AddUint64(&r.totalBytes, t)
}

func (r *defaultMeasurableSuite) stop() {
	r.cancel()
}

func newMeasurableSuite() *defaultMeasurableSuite {

	ctx, cancel := context.WithCancel(context.Background())
	out := &defaultMeasurableSuite{
		ctx:    ctx,
		cancel: cancel,
		ticker: time.NewTicker(time.Second),
		start:  time.Now(),
	}

	go out.loop()
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

type MeasurableConn interface {
	measurable
	net.Conn
}

type readCloser struct {
	measurable
	io.ReadCloser
}

func DecorateReadCloser(rc io.ReadCloser) MeasurableReadCloser {
	return &readCloser{
		measurable: newMeasurableSuite(),
		ReadCloser: rc,
	}
}

func (rc *readCloser) Close() error {
	rc.stop()
	return rc.ReadCloser.Close()
}
func (rc *readCloser) Read(b []byte) (int, error) {

	n, err := rc.ReadCloser.Read(b)
	if err != nil {
		return 0, err
	}

	rc.measurable.addTotal(uint64(n))
	return n, nil
}

func (rc *readCloser) FD() io.ReadCloser {
	return rc.ReadCloser
}

type writeCloser struct {
	measurable
	io.WriteCloser
}

func DecorateWriteCloser(wc io.WriteCloser) MeasurableWriteCloser {
	return &writeCloser{
		measurable:  newMeasurableSuite(),
		WriteCloser: wc,
	}
}

func (wc *writeCloser) FD() io.WriteCloser {
	return wc.WriteCloser
}

func (wc *writeCloser) Write(b []byte) (int, error) {

	n, err := wc.WriteCloser.Write(b)
	if err != nil {
		return 0, err
	}

	wc.measurable.addTotal(uint64(n))
	return n, nil
}

func (wc *writeCloser) Close() error {
	wc.stop()
	return wc.WriteCloser.Close()
}
