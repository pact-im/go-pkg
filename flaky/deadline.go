package flaky

import (
	"context"
	"time"
)

// withinDeadline returns whether the duration d is within context’s deadline.
func withinDeadline(ctx context.Context, d time.Duration) bool {
	deadline, ok := ctx.Deadline()
	if !ok {
		return true
	}

	// TODO(tie): allow passing clock through context.
	until := time.Until(deadline)
	return d < until
}
