package grpcprocess

import (
	"context"
	"net"
	"sync"
	"testing"

	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.pact.im/x/process"
)

// fakeAddr is a fake net.Addr implementation.
type fakeAddr struct{}

// Network implements the net.Addr interface.
func (*fakeAddr) Network() string { return "fake" }

// String implements the net.Addr interface.
func (*fakeAddr) String() string { return "fake" }

// fakeListener is a fake net.Listener implementation for in-memory connections.
type fakeListener struct {
	conn chan net.Conn
	done chan struct{}
	once sync.Once
}

// newFakeListener returns a new fakeListener instance.
func newFakeListener() *fakeListener {
	return &fakeListener{
		conn: make(chan net.Conn),
		done: make(chan struct{}),
	}
}

// Dial returns creates a new connection and waits for the server to accept it.
func (l *fakeListener) Dial(ctx context.Context, _ string) (net.Conn, error) {
	clientConn, serverConn := net.Pipe()
	select {
	case l.conn <- serverConn:
		return clientConn, nil
	case <-l.done:
		return nil, net.ErrClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Accept implements the net.Listener interface.
func (l *fakeListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.conn:
		return conn, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

// Addr implements the net.Listener interface.
func (l *fakeListener) Addr() net.Addr {
	return (*fakeAddr)(nil)
}

// Close implements the net.Listener interface.
func (l *fakeListener) Close() error {
	l.once.Do(func() { close(l.done) })
	return nil
}

func TestServer(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	lis := newFakeListener()
	srv := grpc.NewServer()

	p := process.NewProcess(ctx, Server(srv, lis))
	if err := p.Start(ctx); err != nil {
		t.Fatalf("start server: %v", err)
	}

	cc, err := grpc.DialContext(ctx, "fake",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(lis.Dial),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	if err := cc.Close(); err != nil {
		t.Fatalf("close client: %v", err)
	}

	if err := p.Stop(ctx); err != nil {
		t.Fatalf("stop server: %v", err)
	}
}
