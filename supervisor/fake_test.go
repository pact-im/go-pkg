package supervisor

import (
	"context"
	"sync"
	"testing"

	"go.pact.im/x/process"
)

// fakeRunner is a fake [process.Runner] interface implementation.
type fakeRunner struct{}

// newFakeRunner returns a new fakeRunner instance.
func newFakeRunner() *fakeRunner {
	return (*fakeRunner)(nil)
}

// Run implements the [process.Runner] interface.
func (*fakeRunner) Run(ctx context.Context, callback process.Callback) error {
	return callback(ctx)
}

// observeRunner allows simulating cases in tests where Run blocks
// indefenitely.
type observeRunner struct {
	runner process.Runner
	runc   chan chan<- struct{}
}

// newRunObserver returns a new observeRunner instance that wraps the given
// [process.Runner].
func newObserveRunner(runner process.Runner) *observeRunner {
	return &observeRunner{
		runner: runner,
		runc:   make(chan chan<- struct{}),
	}
}

// Observe returns an unbuffered channel where a channel that should be closed
// to unblock Run method gets sent for each Run call.
func (r *observeRunner) Observe() <-chan chan<- struct{} {
	return r.runc
}

// Run implements the [process.Runner] interface.
func (r *observeRunner) Run(ctx context.Context, callback process.Callback) error {
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
	return r.runner.Run(ctx, callback)
}

func TestObserveRunner(_ *testing.T) {
	r := newObserveRunner(newFakeRunner())
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
