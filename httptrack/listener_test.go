package httptrack

import (
	"context"
	"net"

	"go.pact.im/x/netchan"
)

// testListener wraps netchan.Listener with additional methods for convenience.
type testListener struct {
	*netchan.Listener
}

// newTestListener returns a new testListener instance.
func newTestListener() testListener {
	return testListener{
		Listener: netchan.NewListener(),
	}
}

// Dial is a convenience method that matches [net.Dialer.DialContext] function
// signature.
func (l testListener) Dial(_ context.Context, _, _ string) (net.Conn, error) {
	return l.Pipe(), nil
}

// Pipe creates a pipe and sends the server side to the Accept caller. It
// returns client side of the pipe.
func (l testListener) Pipe() net.Conn {
	server, client := net.Pipe()
	l.C() <- server
	return client
}

// Wait waits for Accept caller to be ready. Note that our Accept implementation
// retries reading from the channel on nil values.
func (l *testListener) Wait() {
	l.C() <- nil
}
