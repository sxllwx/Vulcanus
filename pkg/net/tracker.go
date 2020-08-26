package net

import (
	"net"
	"sync"
	"time"

	"github.com/sxllwx/vulcanus/pkg/log"
)

type ConnTracker struct {

	// config
	rwTimeout time.Duration
	ticker    *time.Ticker

	// tracked conn
	mu      sync.RWMutex
	connMap map[*statusConn]struct{}
}

func NewConnTracker(cleanInterval time.Duration, rwTimeout time.Duration) *ConnTracker {

	return &ConnTracker{
		rwTimeout: rwTimeout,
		ticker:    time.NewTicker(cleanInterval),
		connMap:   make(map[*statusConn]struct{}, 8),
	}
}

// track conn
func (c *ConnTracker) track(conn *statusConn) {

	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.connMap[conn]
	if ok {
		return
	}

	c.connMap[conn] = struct{}{}
}

// untracked conn
func (c *ConnTracker) unTrack(conn *statusConn) {

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
	return c.trackNetConn(conn), nil
}

func (c *ConnTracker) trackNetConn(conn net.Conn) *statusConn {

	sc := &statusConn{
		Conn:    conn,
		tracker: c,
	}
	c.track(sc)
	return sc
}

func (c *ConnTracker) Track(conn net.Conn) net.Conn {
	return c.trackNetConn(conn)
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
	for conn := range c.connMap {
		if err := conn.SetDeadline(conn.latestRWTime.Add(c.rwTimeout)); err != nil {
			log.Warnf("set conn deadline %v", err)
		}
	}
}
