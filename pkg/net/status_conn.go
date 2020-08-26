package net

import (
	"net"
	"sync"
	"time"
)

// statusConn
type statusConn struct {
	net.Conn // underlay conn

	// *** timesetting ***//
	mu           sync.RWMutex
	latestRWTime time.Time
	// *** timesetting ***//

	tracker *ConnTracker
}

func (c *statusConn) Read(b []byte) (int, error) {

	n, err := c.Conn.Read(b)

	c.mu.Lock()
	c.latestRWTime = time.Now()
	c.mu.Unlock()
	return n, err
}

func (c *statusConn) Write(b []byte) (int, error) {

	n, err := c.Conn.Write(b)

	c.mu.Lock()
	c.latestRWTime = time.Now()
	c.mu.Unlock()
	return n, err
}

func (c *statusConn) Close() (err error) {
	c.tracker.unTrack(c.Conn)
	return c.Conn.Close()
}
