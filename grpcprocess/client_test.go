package grpcprocess

import (
	"context"
	"net"
	"testing"
	"testing/synctest"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"go.pact.im/x/netchan"
	"go.pact.im/x/process"
)

func TestClient(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx := context.Background()
		lis := netchan.NewListener()
		srv := grpc.NewServer()
		healthpb.RegisterHealthServer(srv, health.NewServer())

		serverProcess := process.NewProcess(ctx, Server(srv, lis))
		if err := serverProcess.Start(ctx); err != nil {
			t.Fatalf("start server: %v", err)
		}

		client := Client(func(context.Context) (*grpc.ClientConn, error) {
			return grpc.NewClient("passthrough://netchan",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
					return lis.Dial(ctx)
				}),
			)
		})
		healthClient := healthpb.NewHealthClient(client)

		if _, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{}); status.Code(err) != codes.Canceled {
			t.Fatalf("healthcheck before start: got %v, want code %v", err, codes.Canceled)
		}

		clientProcess := process.NewProcess(ctx, client)
		if err := clientProcess.Start(ctx); err != nil {
			t.Fatalf("start client: %v", err)
		}
		if _, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{}); err != nil {
			t.Fatalf("healthcheck: %v", err)
		}

		if err := clientProcess.Stop(ctx); err != nil {
			t.Fatalf("stop client: %v", err)
		}
		if _, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{}); status.Code(err) != codes.Canceled {
			t.Fatalf("healthcheck after stop: got %v, want code %v", err, codes.Canceled)
		}

		if err := serverProcess.Stop(ctx); err != nil {
			t.Fatalf("stop server: %v", err)
		}
	})
}
