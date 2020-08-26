package net

import (
	"net"
	"sync"
	"time"

	"github.com/sxllwx/vulcanus/pkg/log"
)

type ConnTracker struct {
	rwTimeout time.Duration

	ticker *time.Ticker

	mu sync.RWMutex

	connMap map[net.Conn]*statusConn
}

func NewConnTracker(cleanInterval time.Duration, rwTimeout time.Duration) *ConnTracker {

	return &ConnTracker{
		rwTimeout: 0,
		ticker:    time.NewTicker(cleanInterval),
		connMap:   make(map[net.Conn]*statusConn, 8),
	}
}

// track cnn
func (c *ConnTracker) track(conn net.Conn) net.Conn {

	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.connMap[conn]
	if ok {
		return nil
	}

	sc := &statusConn{
		Conn:    conn,
		tracker: c,
	}
	c.connMap[conn] = sc
	return sc
}

// untrack conn
func (c *ConnTracker) unTrack(conn net.Conn) {

	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.connMap[conn]
	if !ok {
		return
	}
	delete(c.connMap, conn)
}

func (c *ConnTracker) Dial(network string, addr string, dialFunc func(string, string) (net.Conn, error)) (net.Conn, error) {
	conn, err := dialFunc(network, addr)
	if err != nil {
		return nil, err
	}

	c.track(conn)
	return conn, nil
}

func (c *ConnTracker) Track(conn net.Conn) net.Conn {
	return c.track(conn)
}

func (c *ConnTracker) Start() {

	for {
		_, ok := <-c.ticker.C
		if !ok {
			// already stopped
			return
		}
		c.resetDeadline()
	}

}

func (c *ConnTracker) Stop() {
	c.ticker.Stop()
}

func (c *ConnTracker) resetDeadline() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, conn := range c.connMap {
		if err := conn.SetDeadline(conn.latestRWTime.Add(c.rwTimeout)); err != nil {
			log.Warnf("set conn deadline %v", err)
		}
	}
}
