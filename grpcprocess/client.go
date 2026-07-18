package grpcprocess

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.pact.im/x/process"
)

// ClientConn is a process-managed gRPC client connection.
type ClientConn interface {
	grpc.ClientConnInterface
	process.Runner
}

// ClientFunc creates a gRPC client connection using the given context.
type ClientFunc func(context.Context) (*grpc.ClientConn, error)

// Client returns a [ClientConn] that creates its underlying [grpc.ClientConn]
// using connect at application runtime.
func Client(connect ClientFunc) ClientConn {
	return &client{connect: connect}
}

type client struct {
	connect ClientFunc
	conn    atomic.Pointer[grpc.ClientConn]
}

// Run implements the [process.Runner] interface.
func (c *client) Run(ctx context.Context, callback process.Callback) error {
	conn, err := c.connect(ctx)
	if err != nil {
		return err
	}

	c.conn.Store(conn)

	callbackError := callback(ctx)

	c.conn.CompareAndSwap(conn, nil)

	_ = conn.Close()

	return callbackError
}

// Invoke implements [grpc.ClientConnInterface].
func (c *client) Invoke(
	ctx context.Context,
	method string,
	args, reply any,
	options ...grpc.CallOption,
) error {
	conn, err := c.load()
	if err != nil {
		return err
	}
	return conn.Invoke(ctx, method, args, reply, options...)
}

// NewStream implements [grpc.ClientConnInterface].
func (c *client) NewStream(
	ctx context.Context,
	desc *grpc.StreamDesc,
	method string,
	options ...grpc.CallOption,
) (grpc.ClientStream, error) {
	conn, err := c.load()
	if err != nil {
		return nil, err
	}
	return conn.NewStream(ctx, desc, method, options...)
}

func (c *client) load() (*grpc.ClientConn, error) {
	conn := c.conn.Load()
	if conn == nil {
		return nil, status.Error(
			codes.Canceled,
			"grpcprocess: client connection is not running",
		)
	}
	return conn, nil
}
