package process

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"

	"go.pact.im/x/clock"
	"go.pact.im/x/clock/fakeclock"
)

func TestManagerStartConflict(t *testing.T) {
	ctx := context.Background()

	var pk struct{}
	var tab mapTable[struct{}, *observeRunnable]
	m := NewManager[struct{}, *observeRunnable](&tab, Options{
		Clock: clock.NewClock(fakeclock.Go()),
	})

	err := m.Run(ctx, func(ctx context.Context) error {
		r := newObserveRunnable(newFakeRunnable())
		observer := r.Observe()
		tab.m.Store(pk, r)

		g, ctx := errgroup.WithContext(ctx)
		defer func() { _ = g.Wait() }()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Wait for ongoing Run method call. Unblock it once we
		// get a result from the second concurrent Start call.
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case unblock := <-observer:
				defer close(unblock)
			}

			_, err := m.Start(ctx, pk)
			if errors.Is(err, ErrProcessExists) {
				return nil
			}
			return fmt.Errorf("unexpected error: %w", err)
		})

		_, err := m.Start(ctx, pk)

		// Wait for the result of the second Start call and only
		// then check that our first call succeeded. That makes
		// it a bit easier to debug should there be an error.
		assert.NilError(t, g.Wait())
		assert.NilError(t, err)
		return nil
	})
	assert.NilError(t, err)
}

func TestManagerStartStop(t *testing.T) {
	ctx := context.Background()

	var pk struct{}
	var tab mapTable[struct{}, *fakeRunnable]
	m := NewManager[struct{}, *fakeRunnable](&tab, Options{
		Clock: clock.NewClock(fakeclock.Go()),
	})

	err := m.Run(ctx, func(ctx context.Context) error {
		tab.m.Store(pk, newFakeRunnable())

		if _, err := m.Start(ctx, pk); err != nil {
			return err
		}
		if _, err := m.Get(ctx, pk); err != nil {
			return err
		}
		if err := m.Stop(ctx, pk); err != nil {
			return err
		}

		return nil
	})
	assert.NilError(t, err)
}
