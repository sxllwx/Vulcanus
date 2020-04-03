package io

import (
	"context"
	"io"
	"net"
	"sync/atomic"
	"time"
)

// default Measurable Suite
// basic model of the metric of everyone call be measurable
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

func (r *defaultMeasurableSuite) Cost() time.Duration {
	return time.Since(r.start)
}

func (r *defaultMeasurableSuite) TotalBytes() uint64 {
	return atomic.LoadUint64(&r.totalBytes)
}

func (r *defaultMeasurableSuite) BytesPerSecond() uint64 {
	return atomic.LoadUint64(&r.bps)
}

func (r *defaultMeasurableSuite) AverageBytesPerSecond() float64 {
	return float64(atomic.LoadUint64(&r.totalBytes)) / time.Since(r.start).Seconds()
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

// loop
// calculate the metric
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

type readWriteCloser struct {
	rm measurable
	wm measurable
	io.ReadWriteCloser
}

// Wrapper  io.ReadWriteCloser to MeasurableReadWriteCloser
func DecorateReadWriteCloser(rwc io.ReadWriteCloser) MeasurableReadWriteCloser {
	return &readWriteCloser{
		rm:              newMeasurableSuite(),
		wm:              newMeasurableSuite(),
		ReadWriteCloser: rwc,
	}
}

func (rwc *readWriteCloser) ReadMetric() measurable {
	return rwc.rm
}

func (rwc *readWriteCloser) WriteMetric() measurable {
	return rwc.wm
}

func (rwc *readWriteCloser) Close() error {
	rwc.rm.stop()
	rwc.wm.stop()
	return rwc.ReadWriteCloser.Close()
}

func (rwc *readWriteCloser) Read(b []byte) (int, error) {

	n, err := rwc.ReadWriteCloser.Read(b)
	if err != nil {
		return 0, err
	}

	rwc.rm.addTotal(uint64(n))
	return n, nil
}

func (rwc *readWriteCloser) Write(b []byte) (int, error) {

	n, err := rwc.ReadWriteCloser.Write(b)
	if err != nil {
		return 0, err
	}

	rwc.wm.addTotal(uint64(n))
	return n, nil
}

type conn struct {
	rm measurable
	wm measurable
	net.Conn
}

// wrapper net.Conn to  MeasurableConn
func DecorateConn(c net.Conn) MeasurableConn {
	return &conn{
		rm:   newMeasurableSuite(),
		wm:   newMeasurableSuite(),
		Conn: c,
	}
}

func (c *conn) ReadMetric() measurable {
	return c.rm
}

func (c *conn) WriteMetric() measurable {
	return c.wm
}

func (c *conn) Read(b []byte) (int, error) {

	n, err := c.Conn.Read(b)
	if err != nil {
		return 0, err
	}

	c.rm.addTotal(uint64(n))
	return n, nil
}

func (c *conn) Write(b []byte) (int, error) {

	n, err := c.Conn.Write(b)
	if err != nil {
		return 0, err
	}

	c.wm.addTotal(uint64(n))
	return n, nil
}

func (c *conn) Close() error {

	c.wm.stop()
	c.rm.stop()

	return c.Conn.Close()
}
