package httpprocess

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"

	"go.uber.org/goleak"

	"go.pact.im/x/netchan"
	"go.pact.im/x/process"
)

func TestServer(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	lis := netchan.NewListener()
	srv := &http.Server{}

	p := process.NewProcess(ctx, Server(srv, lis))
	if err := p.Start(ctx); err != nil {
		t.Fatalf("start server: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return lis.Dial(ctx)
			},
		},
	}
	if err := httpGet(client, "http://example.com"); err != nil {
		t.Fatalf("http get: %v", err)
	}

	if err := p.Stop(ctx); err != nil {
		t.Fatalf("stop server: %v", err)
	}
}

func httpGet(client *http.Client, url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
