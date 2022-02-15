package flaky

import (
	"context"

	"go.pact.im/x/clock"
)

// WatchdogExecutor is an executor that keeps executing an operation on schedule
// using the underlying executor until an error is returned or context expires.
// It allow monitoring a resource that should not fail.
type WatchdogExecutor struct {
	exec *ScheduleExecutor
}

// Watchdog executes an operation on the given schedule using the executor until
// an error is returned or the context expires. It allow monitoring a resource
// that should not fail.
//
// Note that it returns a nil error iff context expires. Otherwise an error from
// the executor indicates an operation failure.
//
// This design allows interrupting execution using the context but the tradeoff
// is that an error is discarded. That is mostly noticeable when the executor
// retries operations and a context expires after a failure. In that case the
// error would not be propagated to the watchdog user.
//
// Unlike a simple loop that executes the given executor and waits a certain
// amount of time, watchdog enforces the use of a configurable schedule.
//
func Watchdog(e Executor, s Schedule) *WatchdogExecutor {
	return &WatchdogExecutor{
		exec: WithSchedule(e, s),
	}
}

// WithClock returns a copy of the executor that uses the given clock.
func (w *WatchdogExecutor) WithClock(c *clock.Clock) *WatchdogExecutor {
	return &WatchdogExecutor{
		exec: w.exec.WithClock(c),
	}
}

// Execute implements the Executor interface.
func (w *WatchdogExecutor) Execute(ctx context.Context, f Op) error {
	for {
		err := w.exec.Execute(ctx, f)
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
