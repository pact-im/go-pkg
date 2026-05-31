package task

import (
	"context"
	"errors"
	"testing"
)

func TestContextifyCancel(t *testing.T) {
	oops := errors.New("oops")

	stop := make(chan struct{})
	c := Contextify(func() error {
		<-stop
		return oops
	}, func() {
		close(stop)
	})

	canceledContext, cancel := context.WithCancel(context.Background())
	cancel()
	err := c.Run(canceledContext)
	if !errors.Is(err, oops) {
		t.Fatalf("expected error %v, got %v", oops, err)
	}
}

func TestContextifyReturn(t *testing.T) {
	c := Contextify(func() error {
		return nil
	}, func() {
		panic("unreachable")
	})
	err := c.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
