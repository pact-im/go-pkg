package flaky

import (
	"context"
	"errors"
	"testing"
)

func TestUntilPermanent(t *testing.T) {
	ctx := context.Background()
	oops := errors.New("oops")
	exec := UntilPermanent()
	t.Run("Nil", func(t *testing.T) {
		err := exec.Execute(ctx, func(_ context.Context) error {
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("Internal", func(t *testing.T) {
		err := exec.Execute(ctx, func(_ context.Context) error {
			return Internal(oops)
		})
		if !errors.Is(err, oops) {
			t.Fatalf("expected error %v, got %v", oops, err)
		}
	})
	t.Run("Permanent", func(t *testing.T) {
		err := exec.Execute(ctx, func(_ context.Context) error {
			return Permanent(oops)
		})
		if !IsPermanentError(err) {
			t.Fatalf("expected permanent error, got %v", err)
		}
		if !errors.Is(err, oops) {
			t.Fatalf("expected error %v, got %v", oops, err)
		}
	})
}
