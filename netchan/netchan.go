// Package netchan provides an in-memory implementation of the [net.Listener]
// interface, allowing programs to simulate network connections using channels
// and pipes instead of real network sockets.
package netchan

import (
	"context"
	"errors"
	"net"
	"sync"
)

var _ net.Listener = (*Listener)(nil)

// ErrClosed indicates that the [Listener] has been closed and no further
// operations can be performed.
var ErrClosed = errors.New("netchan: closed")

// chanAddr implements the [net.chanAddr] interface representing the network
// address of an in-memory [Listener]. Since connections are in-memory, chanAddr
// returns constant values indicating the connection uses channels.
type chanAddr struct{}

// Network returns the name of the network, which is always "chan".
func (chanAddr) Network() string { return "chan" }

// String returns a string representation of the Addr, which is always "chan".
func (chanAddr) String() string { return "chan" }

// Listener implements the [net.Listener] interface using channels for
// accepting connections. Incoming connections are sent on an internal
// channel and accepted by Accept calls. The Listener can be closed to
// unblock Accept calls and signal shutdown.
type Listener struct {
	addr net.Addr
	conn chan net.Conn
	done chan struct{}
	once sync.Once
}

// NewListener returns a new Listener instance with the given network address.
// If the provided addr is nil, the default address is used.
func NewListener(addr net.Addr) *Listener {
	if addr == nil {
		addr = chanAddr{}
	}
	return &Listener{
		addr: addr,
		conn: make(chan net.Conn),
		done: make(chan struct{}),
	}
}

// Accept waits for and returns the next connection sent to the listener’s
// connection channel. If the listener is closed, Accept returns [ErrClosed].
func (l *Listener) Accept() (net.Conn, error) {
	select {
	// Accept should return an error after Close.
	case <-l.done:
		return nil, ErrClosed
	default:
	}
accept:
	select {
	case conn, ok := <-l.conn:
		// Handle unlikely case of close(l.C()).
		if !ok {
			return nil, ErrClosed
		}
		// Allow sending nil connection to wait for Accept call.
		if conn == nil {
			goto accept
		}
		return conn, nil
	// Unblock Accept when Close is called.
	case <-l.done:
		return nil, ErrClosed
	}
}

// Addr implements the [net.Listener] interface.
func (l *Listener) Addr() net.Addr {
	return chanAddr{}
}

// Close closes the listener and unblocks all Accept calls.
func (l *Listener) Close() error {
	l.once.Do(func() { close(l.done) })
	return nil
}

// C returns the internal channel used to send [net.Conn] connections to the
// Accept callers. The returned channel should not be closed; use Close method
// to unblock Accept instead.
//
// Sending a nil connection on this channel is allowed and can be used
// to wait for the Accept caller to be ready.
func (l *Listener) C() chan net.Conn {
	return l.conn
}

// Done returns a channel that is closed when the Listener is closed.
func (l *Listener) Done() <-chan struct{} {
	return l.done
}

// Dial creates a new in-memory connection pair using net.Pipe and sends the
// server side connection to the Listener’s internal connection channel.
//
// It returns the client side connection, which can be used to communicate
// with the accepted connection on the Listener side.
//
// Dial blocks until the server side connection is accepted by a call to Accept,
// or until the Listener is closed or the provided context is canceled.
//
// The second and third string parameters are ignored. They exist for
// convenience to match [net.Dialer.DialContext] function signature.
func (l *Listener) Dial(ctx context.Context, _, _ string) (net.Conn, error) {
	clientConn, serverConn := net.Pipe()
	select {
	case l.conn <- serverConn:
		return clientConn, nil
	case <-l.done:
		return nil, ErrClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
