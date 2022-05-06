package process

import (
	"context"
	"sync"
	"testing"
)

// fakeRunnable is a fake Runnable interface implementation.
type fakeRunnable struct{}

// newFakeRunnable returns a new fakeRunnable instance.
func newFakeRunnable() *fakeRunnable {
	return (*fakeRunnable)(nil)
}

// Run implements the Runnable interface.
func (r *fakeRunnable) Run(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

// observeRunnable allows simulating cases in tests where Run blocks indefenitely.
type observeRunnable struct {
	proc Runnable
	runc chan chan<- struct{}
}

// newRunObserver returns a new observeRunnable instance that wraps proc Runnable.
func newObserveRunnable(proc Runnable) *observeRunnable {
	return &observeRunnable{
		proc: proc,
		runc: make(chan chan<- struct{}),
	}
}

// Observe returns an unbuffered channel where a channel that should be closed
// to unblock Run method gets sent for each Run call.
func (r *observeRunnable) Observe() <-chan chan<- struct{} {
	return r.runc
}

// Run implements the Runnable interface.
func (r *observeRunnable) Run(ctx context.Context, f func(ctx context.Context) error) error {
	unblock := make(chan struct{})
	select {
	case <-ctx.Done():
		return ctx.Err()
	case r.runc <- unblock:
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-unblock:
	}
	return r.proc.Run(ctx, f)
}

func TestObserveRunnable(t *testing.T) {
	r := newObserveRunnable(newFakeRunnable())
	observer := r.Observe()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = r.Run(context.Background(), func(ctx context.Context) error {
			return nil
		})
	}()

	close(<-observer)
	wg.Wait()
}
