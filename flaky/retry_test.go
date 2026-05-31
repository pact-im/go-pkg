package flaky

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"go.pact.im/x/clock"
	"go.pact.im/x/clock/mockclock"
)

func TestRetry(t *testing.T) {
	t.Run("Backoff", func(t *testing.T) {
		opError := errors.New("op")
		backoff := []time.Duration{
			time.Nanosecond,
			time.Second,
			time.Minute,
			time.Hour,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		executor := Retry(Constant(backoff...)).WithClock(func() *clock.Clock {
			ctrl := gomock.NewController(t)
			m := mockclock.NewMockClock(ctrl)
			m.EXPECT().Timer(backoff[0]).Return(func() clock.Timer {
				timer := mockclock.NewMockTimer(ctrl)
				defer timer.EXPECT().Stop().Return(false)

				var waitC chan time.Time
				travelC := make(chan time.Time)
				close(travelC)

				timer.EXPECT().C().Return(waitC).After(
					timer.EXPECT().C().Return(travelC).Times(len(backoff) - 1),
				)

				for _, d := range backoff[1 : len(backoff)-1] {
					timer.EXPECT().Reset(d).Return()
				}

				last := backoff[len(backoff)-1]
				timer.EXPECT().Reset(last).Do(func(time.Duration) {
					cancel()
				}).Return()

				return timer
			}())
			return clock.NewClock(m)
		}())

		var n int
		err := executor.Execute(ctx, func(_ context.Context) error {
			n++
			return opError
		})
		if err != opError {
			t.Fatalf("expected opError, got %v", err)
		}
		if n != len(backoff) {
			t.Fatalf("expected n==%d, got %d", len(backoff), n)
		}
	})

	t.Run("Limit", func(t *testing.T) {
		opError := errors.New("op")
		backoff := time.Second

		ctx := context.Background()

		executor := Retry(Limit(1, Constant(backoff))).WithClock(func() *clock.Clock {
			ctrl := gomock.NewController(t)
			m := mockclock.NewMockClock(ctrl)
			m.EXPECT().Timer(backoff).Return(func() clock.Timer {
				timer := mockclock.NewMockTimer(ctrl)
				defer timer.EXPECT().Stop().Return(false)

				var waitC chan time.Time

				travelC := make(chan time.Time)
				close(travelC)

				timer.EXPECT().C().Return(waitC).AnyTimes().After(
					timer.EXPECT().C().Return(travelC),
				)
				timer.EXPECT().Reset(backoff).Return().AnyTimes()

				return timer
			}())
			return clock.NewClock(m)
		}())

		var n int
		err := executor.Execute(ctx, func(_ context.Context) error {
			n++
			return opError
		})
		if err != opError {
			t.Fatalf("expected opError, got %v", err)
		}
		if n != 2 {
			t.Fatalf("expected n==2, got %d", n)
		}
	})
}
