// Package grpcprocess provides [process.Runnable] wrapper for [grpc.Server].
package grpcprocess

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"go.pact.im/x/process"
)

// Server returns a [process.Runnable] instance for the given gRPC server and
// network listener.
func Server(srv *grpc.Server, lis net.Listener) process.Runnable {
	return process.Leaf(
		func(_ context.Context) error {
			// TODO: consider injecting base context into gRPC’s
			// internal/transport/http2_server.go NewServerTransport.
			return srv.Serve(lis)
		},
		func(ctx context.Context) error {
			done := make(chan struct{})
			go func() {
				srv.GracefulStop()
				close(done)
			}()
			select {
			case <-ctx.Done():
				// Note that Stop interrupts GracefulStop when
				// called concurrently.
				srv.Stop()
				<-done
			case <-done:
			}
			return nil
		},
	)
}
