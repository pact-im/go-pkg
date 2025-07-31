package flaky

import (
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"

	"go.pact.im/x/clock"
	"go.pact.im/x/clock/fakeclock"
	"go.pact.im/x/clock/mockclock"
	"go.pact.im/x/clock/observeclock"
)

func TestDebounceExecutor(t *testing.T) {
	const wait = time.Second

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var n uint64
	op := func(_ context.Context) error {
		n++
		return nil
	}

	fakeClock := fakeclock.Unix()
	observeClock := observeclock.New(fakeClock)
	observer := observeClock.Observe()
	go func() {
		<-observer
		fakeClock.Add(wait)
	}()

	d := new(DebounceExecutor).WithClock(clock.NewClock(observeClock)).WithWait(wait)
	assert.Assert(t, nil == d.Execute(ctx, op))
	assert.Assert(t, n == 1)
}

func TestDebounceExecutorCancel(t *testing.T) {
	const wait = time.Second

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var n uint64
	op := func(_ context.Context) error {
		n++
		return nil
	}

	fakeClock := fakeclock.Unix()
	observeClock := observeclock.New(fakeClock)
	observer := observeClock.Observe()
	go func() {
		<-observer
		cancel()
	}()

	debouncer := Debounce(wait).WithClock(clock.NewClock(observeClock))

	err := debouncer.Execute(ctx, op)
	assert.Assert(t, ctx.Err() == err)
	assert.Assert(t, n == 0)
}

func TestDebounceExecutorConcurrentExecute(t *testing.T) {
	const wait = time.Second

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var n uint64
	op := Op(func(_ context.Context) error {
		n++
		return nil
	})

	fakeClock := fakeclock.Unix()
	observeClock := observeclock.New(fakeClock)
	observer := observeClock.Observe()
	go func() {
		<-observer
		fakeClock.Add(wait)
	}()

	debouncer := Debounce(wait).WithClock(clock.NewClock(observeClock))

	err := debouncer.Execute(ctx, func(ctx context.Context) error {
		err := debouncer.Execute(ctx, op)
		assert.ErrorIs(t, err, ErrDebounced)
		return op(ctx)
	})
	assert.NilError(t, err)
	assert.Assert(t, n == 1)
}

func TestDebounceExecutorConcurrentCancel(t *testing.T) {
	const wait = time.Second

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var n uint64
	op := Op(func(_ context.Context) error {
		n++
		return nil
	})

	ctrl := gomock.NewController(t)
	mockClock := mockclock.NewMockClock(ctrl)

	debouncer := Debounce(wait).WithClock(clock.NewClock(mockClock))

	mockClock.EXPECT().Timer(wait).Times(1).DoAndReturn(func(_ time.Duration) *mockclock.MockTimer {
		// We want to deterministically test the cancellation for
		// acquiring or stealing a lock.
		//
		// At this point we are essentially blocking clock’s Timer
		// method, so the debouncer implementation will not be able to
		// pass the lock to another Execute invocation, and at the same
		// time we are already holding a lock.
		//
		// That is, it is guaranteed that we select done channel between
		// a read from done and a write to either lock or steal channels.
		//
		cancel()
		err := debouncer.Execute(ctx, op)
		assert.Assert(t, ctx.Err() == err)

		// Create a timer that never expires. Now that we’ve canceled
		// the context and there are no other ongoing Execute calls, it
		// is guaranteed that we will return from the main Execute call.
		timer := mockclock.NewMockTimer(ctrl)
		timer.EXPECT().C().Times(1).Return(nil)
		timer.EXPECT().Stop().Times(1).Return(true)
		return timer
	})

	err := debouncer.Execute(ctx, op)
	assert.Assert(t, ctx.Err() == err)
	assert.Assert(t, n == 0)
}

func TestDebounceExecutorConcurrentSteal(t *testing.T) {
	const wait = time.Second

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var n uint64
	op := Op(func(_ context.Context) error {
		n++
		return nil
	})

	ctrl := gomock.NewController(t)
	mockClock := mockclock.NewMockClock(ctrl)

	debouncer := Debounce(wait).WithClock(clock.NewClock(mockClock))

	rendezvous, stole := make(chan struct{}), make(chan error)
	go func() {
		<-rendezvous
		stole <- debouncer.Execute(ctx, op)
	}()

	timer := mockclock.NewMockTimer(ctrl)
	timerC := timer.EXPECT().C().Times(1).DoAndReturn(func() chan time.Time {
		// We want to deterministically test the Execute method stealing a
		// lock from an ongoing Execute call.
		//
		// At this point we are ready to handle steals and return a nil
		// timer channel. We close rendezvous channel to start goroutine
		// that would steal the lock from us.
		//
		close(rendezvous)
		return nil
	})
	timerStop := timer.EXPECT().Stop().Times(1).After(timerC).Return(true)
	resetTimer := timer.EXPECT().Reset(wait).Times(1).After(timerStop).Return()
	timer.EXPECT().C().Times(1).After(resetTimer).DoAndReturn(func() <-chan time.Time {
		instant := make(chan time.Time, 1)
		instant <- time.Unix(0, 0)
		return instant
	})

	mockClock.EXPECT().Timer(wait).Times(1).Return(timer)

	assert.Assert(t, ErrDebounced == debouncer.Execute(ctx, op))
	assert.Assert(t, nil == <-stole)
	assert.Assert(t, n == 1)
}

func TestDebounceExecutorTimerStop(t *testing.T) {
	const wait = time.Second

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var n uint64
	op := func(_ context.Context) error {
		n++
		return nil
	}

	ctrl := gomock.NewController(t)
	mockClock := mockclock.NewMockClock(ctrl)

	debouncer := Debounce(wait).WithClock(clock.NewClock(mockClock))

	// We want to deterministically test the Execute method stopping the
	// timer on cancellation when there was a write to timer’s channel
	// to ensure that there are no races that would otherwise deadlock
	// the debouncer.
	//
	// Once we are waiting for timer expiration, we cancel the operation and
	// simulate a timer expiration after we are reading from timer’s channel
	// but before Stop method returns.
	//
	// Note that we return false from Stop since it did not prevent a write
	// to timer’s channel (that is, a write to instant below).
	//
	instant := make(chan time.Time, 1)
	timer := mockclock.NewMockTimer(ctrl)
	timerC := timer.EXPECT().C().Times(1).DoAndReturn(func() <-chan time.Time {
		cancel()
		return instant
	})
	timerStop := timer.EXPECT().Stop().Times(1).After(timerC).DoAndReturn(func() bool {
		instant <- time.Unix(0, 0)
		return false
	})
	timer.EXPECT().C().Times(1).After(timerStop).Return(instant)

	mockClock.EXPECT().Timer(wait).Times(1).Return(timer)

	err := debouncer.Execute(ctx, op)
	assert.Assert(t, ctx.Err() == err)
	assert.Assert(t, n == 0)
}
