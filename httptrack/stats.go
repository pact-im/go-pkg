package httptrack

import (
	"maps"
	"net"
	"net/http"
	"sync"
)

// counterSize is the count of [http.ConnState] states that are used for
// counter-based statistics.
const counterSize = http.StateClosed + 1 // all states

// gaugeSize is the count of [http.ConnState] states that are used for
// gauge-based statistics.
const gaugeSize = http.StateIdle + 1 // non-terminal states

// Stats holds a snapshot of HTTP connection state statistics.
type Stats struct {
	counter [counterSize]uint64
	gauge   [gaugeSize]uint
}

// Accepted returns the number of connections currently in the [http.StateNew]
// state.
//
// A non-zero value typically indicates that there are accepted connections
// which have not yet sent a complete HTTP request. This may occur naturally
// under normal traffic, but a persistently high number of connections in this
// state can suggest:
//
//   - Slow or high-latency clients
//   - Clients intentionally delaying requests (e.g. denial-of-service attack)
//   - Server under heavy load or accepting connections faster than it can
//     process them
func (s *Stats) Accepted() uint {
	return s.gauge[http.StateNew]
}

// Active returns the number of connections currently in the [http.StateActive]
// state.
//
// In HTTP/1.x, this usually means that a request is currently being processed.
//
// In HTTP/2, this state indicates that the connection has at least one open
// stream. However, the connection may briefly enter the Active state after
// being accepted as part of reading the HTTP/2 connection preface, before
// transitioning to Idle state. Therefore, a connection in the Active state
// is not guaranteed to be actively handling a request.
func (s *Stats) Active() uint {
	return s.gauge[http.StateActive]
}

// Idle returns the number of connections currently in the [http.StateIdle]
// state.
//
// In HTTP/1.x, Idle connections are typically those that have completed
// a request and are being kept alive. A non-zero Idle count is normal in
// servers that support HTTP keep-alive.
//
// In HTTP/2, this state indicates that the connection has no open streams.
// A newly established HTTP/2 connection enters the Idle state after reading
// the connection preface.
func (s *Stats) Idle() uint {
	return s.gauge[http.StateIdle]
}

// AcceptedTotal returns the total number of connections ever seen in the
// [http.StateNew] state.
//
// This is a monotonically increasing counter that reflects the cumulative
// number of accepted connections. It includes all connections, regardless of
// whether they were later closed, hijacked, or reused.
func (s *Stats) AcceptedTotal() uint64 {
	return s.counter[http.StateNew]
}

// ActiveTotal returns the total number of times connections have transitioned
// into the [http.StateActive] state.
//
// This counter increases every time a connection becomes active. It does not
// represent the number of requests, since a single connection may enter Active
// state for multiple requests (e.g. HTTP/2 streams).
func (s *Stats) ActiveTotal() uint64 {
	return s.counter[http.StateActive]
}

// IdleTotal returns the total number of times connections have entered the
// [http.StateIdle] state.
func (s *Stats) IdleTotal() uint64 {
	return s.counter[http.StateIdle]
}

// HijackedTotal returns the number of connections that have been hijacked via
// [http.Hijacker].
//
// A hijacked connection is taken over by the application (e.g. for WebSocket
// upgrade) and is no longer managed by the HTTP server.
func (s *Stats) HijackedTotal() uint64 {
	return s.counter[http.StateHijacked]
}

// ClosedTotal returns the number of connections that have been closed.
//
// A non-zero Closed count indicates connections that have been fully closed
// by either the server or the client, including normal connection termination
// and error scenarios.
//
// High numbers in Closed are typical over time, but a sudden spike may
// indicate:
//
//   - Mass client disconnects
//   - Server-side timeouts or errors
//   - Deployment cycles or load balancer resets
func (s *Stats) ClosedTotal() uint64 {
	return s.counter[http.StateClosed]
}

// StatsTracker is a [Tracker] implementation that keeps counts of HTTP
// connections in each state. It can be used for monitoring connections
// and their lifecycle.
type StatsTracker struct {
	mu    sync.Mutex
	conns map[net.Conn]http.ConnState
	stats Stats
}

// Track records the transition of a connection’s state. It updates internal
// statistics and tracks the current state of each connection.
//
// The dynamic type of the provided net.Conn must be comparable (i.e. usable as
// a map key). Most standard connection types (such as tls.Conn pointer) meet
// this requirement. Passing a non-comparable type will cause a runtime panic.
func (s *StatsTracker) Track(conn net.Conn, state http.ConnState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.counter[state]++

	if oldState, ok := s.conns[conn]; ok {
		s.stats.gauge[oldState]--
	}

	switch state {
	case http.StateHijacked, http.StateClosed:
		delete(s.conns, conn)
	default: // non-terminal states
		if s.conns == nil {
			s.conns = make(map[net.Conn]http.ConnState)
		}
		s.conns[conn] = state
		s.stats.gauge[state]++
	}
}

// Stats returns a snapshot of connection statistics.
func (s *StatsTracker) Stats() Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stats
}

// Connections returns a snapshot of connection states.
func (s *StatsTracker) Connections() map[net.Conn]http.ConnState {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Note that we return the original map instead of a cloned map. That
	// is, as a side effect of this method, the internal map is shrinked.
	// See https://go.dev/issue/20135 and https://antonz.org/go-map-shrink
	var conns map[net.Conn]http.ConnState
	s.conns, conns = maps.Clone(s.conns), s.conns
	return conns
}
