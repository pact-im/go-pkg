package httptrack

import (
	"net"
	"net/http"
	"sync"
)

// ConnTracker is a connection tracker for HTTP connections. It can be
// used to wait for all existing HTTP connections to finish after performing
// HTTP server shutdown.
//
// Note that [http.Server.Close] closes all connections but does not wait for
// per-connection goroutines to return. If the connection was active, it is up
// to the running [http.Handler] to handle connection error and eventually
// return. ConnTracker allows waiting for all in-flight handlers.
type ConnTracker struct {
	wg sync.WaitGroup
}

// Track updates the connection counter based on the connection’s state.
// It should be assigned to the ConnState field of an [http.Server].
func (c *ConnTracker) Track(_ net.Conn, state http.ConnState) {
	switch state {
	// Newly accepted connection from http.Server.Serve method (same
	// goroutine as the caller).
	case http.StateNew:
		c.wg.Add(1)
	// Hijacked (e.g. WebSockets) and closed connections from
	// http.conn.serve method (separate goroutine per connection).
	case http.StateHijacked, http.StateClosed:
		c.wg.Done()
	}
}

// Wait blocks until all tracked HTTP connections have completed. It should be
// called after [http.Server.Serve] has returned.
func (c *ConnTracker) Wait() {
	c.wg.Wait()
}
