package httptrack

import (
	"net"
	"net/http"
)

// Tracker is an interface for tracking HTTP connection state transitions.
// It matches the signature required by [http.Server.ConnState], allowing
// implementations to observe and respond to connection lifecycle events.
type Tracker interface {
	Track(net.Conn, http.ConnState)
}

// TrackerFunc is an adapter to allow the use of ordinary functions as
// [Tracker] implementations. If f is a function with the appropriate
// signature, TrackerFunc(f) is a [Tracker] that calls f.
type TrackerFunc func(net.Conn, http.ConnState)

// Track calls f(conn, state).
func (f TrackerFunc) Track(conn net.Conn, state http.ConnState) {
	f(conn, state)
}

// Compose returns a TrackerFunc that invokes multiple [Tracker] hooks. Each call
// to Track on the returned TrackerFunc will call Track on all provided Tracker
// implementations, in the order they are passed.
//
// When using a [ConnTracker] for graceful shutdown, it should be passed
// as the last argument to Compose. This ensures that other trackers have
// completed their processing before the server is allowed to shut down.
func Compose(hooks ...Tracker) TrackerFunc {
	return func(conn net.Conn, state http.ConnState) {
		for _, h := range hooks {
			h.Track(conn, state)
		}
	}
}

// Wrap configures an [http.Server] to use the given [Tracker] via server’s
// ConnState hook. If ConnState is already set, it is preserved and runs before
// the provided [Tracker], and Wrap returns the old value.
func Wrap(s *http.Server, h Tracker) TrackerFunc {
	oldConnState := TrackerFunc(s.ConnState)
	if oldConnState != nil {
		s.ConnState = Compose(oldConnState, h)
	} else {
		s.ConnState = h.Track
	}
	return oldConnState
}
