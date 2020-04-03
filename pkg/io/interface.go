package io

import (
	"io"
	"net"
	"time"
)

// internal measurable interface
type measurable interface {

	// real time metric
	BytesPerSecond() uint64

	// bytes number until now Read/Write
	TotalBytes() uint64
	// time cost until now
	Cost() time.Duration
	// Read/Write bytes per second
	AverageBytesPerSecond() float64

	// calculate the metric
	addTotal(uint64)
	stop()
}

type MeasurableReadWriteCloser interface {
	ReadMetric() measurable
	WriteMetric() measurable
	io.ReadWriteCloser
}

type MeasurableConn interface {
	ReadMetric() measurable
	WriteMetric() measurable
	net.Conn
}
