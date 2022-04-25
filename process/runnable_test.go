package process

import (
	"context"
	"errors"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
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

func TestLeafCallbackReturns(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		<-ctx.Done()
		return oops
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.ErrorIs(t, err, oops)
}

func TestLeafCallbackError(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		return oops
	})
	assert.ErrorIs(t, err, oops)
}

func TestLeafRunReturns(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		return nil
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return oops
	})
	assert.ErrorIs(t, err, oops)
}

func TestLeafRunError(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		return oops
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	assert.ErrorIs(t, err, oops)
}

func TestStartStopLeafProcess(t *testing.T) {
	started, stopped := make(chan struct{}), make(chan struct{})
	p := NewProcess(
		context.Background(),
		Leaf(func(ctx context.Context) error {
			close(started)
			<-ctx.Done()
			<-stopped
			return nil
		}, func(ctx context.Context) error {
			close(stopped)
			return nil
		}),
	)
	proc := StartStop(p.Start, p.Stop)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		<-started
		return nil
	})
	assert.NilError(t, err)
}

func TestStartStopErrorOnStart(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return oops },
		func(ctx context.Context) error { panic("unreachable") },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.ErrorIs(t, err, oops)
}

func TestStartStopErrorOnCallback(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return nil },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return oops
	})
	assert.ErrorIs(t, err, oops)
}

func TestStartStopErrorOnCallbackAndStop(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return errors.New("ignored") },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return oops
	})
	assert.ErrorIs(t, err, oops)
}

func TestStartStopErrorOnStop(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return oops },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.ErrorIs(t, err, oops)
}
