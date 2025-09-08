package grpcprocess

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.uber.org/goleak"

	"go.pact.im/x/netchan"
	"go.pact.im/x/process"
)

func TestServer(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	lis := netchan.NewListener()
	srv := grpc.NewServer()

	p := process.NewProcess(ctx, Server(srv, lis))
	if err := p.Start(ctx); err != nil {
		t.Fatalf("start server: %v", err)
	}

	cc, err := grpc.DialContext(ctx, "fake",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.Dial(ctx)
		}),
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
