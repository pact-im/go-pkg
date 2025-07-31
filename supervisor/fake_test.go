package supervisor

import (
	"context"
	"sync"
	"testing"

	"go.pact.im/x/process"
)

// fakeRunnable is a fake process.Runnable interface implementation.
type fakeRunnable struct{}

// newFakeRunnable returns a new fakeRunnable instance.
func newFakeRunnable() *fakeRunnable {
	return (*fakeRunnable)(nil)
}

// Run implements the process.Runnable interface.
func (*fakeRunnable) Run(ctx context.Context, callback process.Callback) error {
	return callback(ctx)
}

// observeRunnable allows simulating cases in tests where Run blocks
// indefenitely.
type observeRunnable struct {
	proc process.Runnable
	runc chan chan<- struct{}
}

// newRunObserver returns a new observeRunnable instance that wraps proc
// process.Runnable.
func newObserveRunnable(proc process.Runnable) *observeRunnable {
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

// Run implements the process.Runnable interface.
func (r *observeRunnable) Run(ctx context.Context, callback process.Callback) error {
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
	return r.proc.Run(ctx, callback)
}

func TestObserveRunnable(_ *testing.T) {
	r := newObserveRunnable(newFakeRunnable())
	observer := r.Observe()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = r.Run(context.Background(), func(_ context.Context) error {
			return nil
		})
	}()

	close(<-observer)
	wg.Wait()
}
