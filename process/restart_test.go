package process

import (
	"context"
	"testing"

	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"

	"go.pact.im/x/clock"
	"go.pact.im/x/clock/fakeclock"
	"go.pact.im/x/clock/observeclock"
)

func TestManagerRestartInitial(t *testing.T) {
	t.Run("Timeout", func(t *testing.T) {
		testManagerRestartInitial(t, true)
	})
	t.Run("Alive", func(t *testing.T) {
		testManagerRestartInitial(t, false)
	})
}

func testManagerRestartInitial(t *testing.T, timeout bool) {
	ctx := context.Background()
	initc := make(chan struct{})
	stopc := make(chan struct{})

	var pk struct{}
	var tab mapTable[struct{}, *observeRunnable]

	fakeClock := fakeclock.Go()
	observeClock := observeclock.New(fakeClock)
	clockObserver := observeClock.Observe()

	// Create the initial process.
	r := newObserveRunnable(newFakeRunnable())
	tab.m.Store(pk, r)
	runObserver := r.Observe()

	m := NewManager[struct{}, *observeRunnable](&tab, Options{
		Clock: clock.NewClock(observeClock),
	})

	// Run process manager in the background.
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return m.Run(ctx, func(ctx context.Context) error {
			t.Log("manager is running")
			close(initc)
			select {
			case <-ctx.Done():
				t.Log("received cancellation")
			case <-stopc:
				t.Log("graceful shutdown")
			}
			return nil
		})
	})
	defer func() {
		cancel()
		_ = g.Wait()
	}()

	// Wait until we create the timer for the initial processes. Continue
	// startup in background. Note that the process startup never succeeds
	// if timeout is true since we do not unblock Run method.
	t.Log("waiting for clock observation")
	select {
	case <-ctx.Done():
		assert.NilError(t, ctx.Err())
	case <-clockObserver:
		t.Log("observed clock")
		if timeout {
			fakeClock.Add(restartInitialWait)
		}
	}
	t.Log("waiting for run observation")
	select {
	case <-ctx.Done():
		assert.NilError(t, ctx.Err())
	case unblock := <-runObserver:
		t.Log("observed run")
		if !timeout {
			close(unblock)
		}
	}

	// Note that we may not have reached Run callback at this point. It is
	// valid to call Start here since it requires start field to be true,
	// which it is since we’ve just observed another process startup.
	_, err := m.Start(ctx, pk)
	assert.ErrorIs(t, err, ErrProcessExists)

	// Wait for manager (and thus our process) to complete initialization
	// before checking the state.
	if !timeout {
		t.Log("waiting for process to reach running state")
		select {
		case <-ctx.Done():
			assert.NilError(t, ctx.Err())
		case <-initc:
		}
	}
	if _, err := m.Get(ctx, pk); timeout {
		assert.ErrorIs(t, err, ErrProcessNotRunning)
	} else {
		assert.NilError(t, err)
	}

	// Stop manager and wait for the result. Force shutdown if we are
	// simulating process that is stuck in starting state.
	t.Log("shutting down")
	if !timeout {
		close(stopc)
	} else {
		cancel()
	}
	assert.NilError(t, g.Wait())
}

func TestManagerRestartTimeoutShutdown(t *testing.T) {
	ctx := context.Background()

	var pk struct{}
	var tab mapTable[struct{}, *observeRunnable]

	fakeClock := fakeclock.Go()
	observeClock := observeclock.New(fakeClock)
	clockObserver := observeClock.Observe()

	m := NewManager[struct{}, *observeRunnable](&tab, Options{
		Clock: clock.NewClock(observeClock),
	})

	err := m.Run(ctx, func(ctx context.Context) error {
		// Create a new process.
		r := newObserveRunnable(newFakeRunnable())
		tab.m.Store(pk, r)
		runObserver := r.Observe()

		// Wait for the restart loop and start it.
		select {
		case <-ctx.Done():
			assert.NilError(t, ctx.Err())
		case <-clockObserver:
			fakeClock.Add(restartLoopInterval)
		}

		// Wait for process to begin startup.
		select {
		case <-ctx.Done():
			assert.NilError(t, ctx.Err())
		case <-runObserver:
		}

		_, err := m.Get(ctx, pk)
		assert.ErrorIs(t, err, ErrProcessNotRunning)

		// Now, shutdown. A process in starting state should be stopped.
		return nil
	})
	assert.NilError(t, err)
}
