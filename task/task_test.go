package task

import (
	"context"
	"errors"
	"testing"
)

func TestNamed(t *testing.T) {
	oops := errors.New("oops")
	task := Named("test", func(_ context.Context) error {
		return oops
	})
	err := task.Run(context.Background())
	if err.Error() != "test: oops" {
		t.Fatalf("expected 'test: oops', got %v", err)
	}
}

func TestParallelCancelOnError(t *testing.T) {
	oops := errors.New("oops")
	g := Parallel(CancelOnError(),
		func(_ context.Context) error {
			return oops
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		},
	)
	err := g.Run(context.Background())
	if !errors.Is(err, oops) {
		t.Fatalf("expected error %v, got %v", oops, err)
	}
}

func TestParallelCancelOnReturn(t *testing.T) {
	g := Parallel(CancelOnReturn(),
		func(_ context.Context) error {
			return nil
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		},
	)
	err := g.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
