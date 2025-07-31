package flaky

import (
	"context"
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"github.com/robfig/cron/v3"

	"go.pact.im/x/clock"
	"go.pact.im/x/clock/fakeclock"
	"go.pact.im/x/clock/observeclock"
)

func TestScheduleExecutorCancel(t *testing.T) {
	sched, err := cron.ParseStandard("20 4 * * *") // At 04:20.
	assert.NilError(t, err)

	fakeClock := fakeclock.Unix()
	observeClock := observeclock.New(fakeClock)
	observer := observeClock.Observe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-observer
		cancel()
	}()

	executor := WithSchedule(Once(), sched).WithClock(clock.NewClock(observeClock))

	err = executor.Execute(ctx, func(_ context.Context) error {
		panic("canceled operation should not be executed")
	})
	assert.Equal(t, ctx.Err(), err)
}

func TestScheduleExecutorRetry(t *testing.T) {
	const backoff = time.Second

	sched, err := cron.ParseStandard("20 4 * * *") // At 04:20.
	assert.NilError(t, err)

	fakeClock := fakeclock.Time(time.Date(2038, time.January, 19, 3, 14, 7, 0, time.UTC)) // 03:14
	observeClock := observeclock.New(fakeClock)
	observer := observeClock.Observe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-observer
		observer = observeClock.Observe()

		fakeClock.Add(time.Hour + 6*time.Minute) // 4:20

		<-observer
		fakeClock.Add(backoff)
	}()

	c := clock.NewClock(observeClock)

	inner := Retry(Constant(backoff)).WithClock(c)
	executor := WithSchedule(inner, sched).WithClock(c)

	var init bool
	err = executor.Execute(ctx, func(_ context.Context) error {
		var err error
		if !init {
			init = true
			err = errors.New("retry")
		}
		return err
	})
	assert.NilError(t, err)
}
