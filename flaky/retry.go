package flaky

import (
	"context"
	"math"

	"go.pact.im/x/clock"
)

// RetryExecutor is an executor that retries executing an operation using the
// provided backoff. It allows executing operations that may succeed after a few
// attempts.
type RetryExecutor struct {
	clock   *clock.Clock
	backoff BackoffProvider
}

// Retry returns a new executor that uses the backoff provider to retry
// operation execution.
//
// It attempts to execute f until a backoff function returns false, an attempt
// returns permanent error or a context expires. On failure it returns the last
// error encountered.
//
// It requests a new backoff function from provider for each Execute invocation.
//
func Retry(b BackoffProvider) *RetryExecutor {
	return &RetryExecutor{
		clock:   clock.System(),
		backoff: b,
	}
}

// WithClock returns a copy of the executor that uses the given clock.
func (r *RetryExecutor) WithClock(c *clock.Clock) *RetryExecutor {
	if c == nil {
		c = clock.System()
	}
	return &RetryExecutor{
		clock:   c,
		backoff: r.backoff,
	}
}

// Execute implements the Executor interface.
func (r *RetryExecutor) Execute(ctx context.Context, f Op) error {
	backoff := r.backoff.Backoff()

	var err error
	var timer clock.Timer
	for n := uint(0); n < math.MaxUint; n++ {
		err = f(ctx)
		if err == nil {
			return nil
		}
		if IsPermanentError(err) {
			return unwrapInternal(err)
		}

		if ctx.Err() != nil {
			break
		}

		d, ok := backoff(n)
		if !ok {
			break
		}

		if !withinDeadline(ctx, d) {
			break
		}

		if timer == nil {
			timer = r.clock.Timer(d)
			defer timer.Stop()
		} else {
			timer.Reset(d)
		}
		select {
		case <-ctx.Done():
			return err
		case <-timer.C():
			continue
		}
	}
	return err
}
