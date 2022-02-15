package flaky

import (
	"context"
)

type onceExecutor struct{}

// Once returns a new executor that does not retry operations. That is,
// it executes the operation exactly once and returns the result.
func Once() Executor {
	return (*onceExecutor)(nil)
}

// Execute implements the Executor interface.
func (*onceExecutor) Execute(ctx context.Context, f Op) error {
	return unwrapInternal(f(ctx))
}
