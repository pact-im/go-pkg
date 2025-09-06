package httptrack

import (
	"context"
	"errors"
	"net"
	"sync"
)

type testAddress struct{}

func (testAddress) Network() string { return "test" }
func (testAddress) String() string  { return "test" }

type testListener struct {
	conn chan net.Conn
	once sync.Once
}

func newTestListener() *testListener {
	return &testListener{
		conn: make(chan net.Conn),
	}
}

func (*testListener) Addr() net.Addr {
	return testAddress{}
}

func (l *testListener) Accept() (net.Conn, error) {
	for {
		c, ok := <-l.conn
		if !ok {
			return nil, errors.New("testListener closed")
		}
		if c == nil {
			// We’ve unblocked testListener.Wait, retry.
			continue
		}
		return c, nil
	}
}

func (l *testListener) Close() error {
	l.once.Do(func() { close(l.conn) })
	return nil
}

func (l *testListener) Dial(_ context.Context, _, _ string) (net.Conn, error) {
	return l.Pipe(), nil
}

// Pipe creates a pipe and sends the server side to the Accept caller. It
// returns client side of the pipe.
func (l *testListener) Pipe() net.Conn {
	server, client := net.Pipe()
	l.conn <- server
	return client
}

// Wait waits for Accept call to be ready. Note that our Accept implementation
// retries reading from the channel on nil values.
func (l *testListener) Wait() {
	l.conn <- nil
}
