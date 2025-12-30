// Package flaky implements mechanisms for executing flaky operations.
package flaky

import (
	"context"
)

// Op is an operation that is likely to fail.
type Op func(ctx context.Context) error

// Executor is an executor for flaky operations.
type Executor interface {
	// Execute executes a flaky operation f using the execution policy.
	Execute(ctx context.Context, f Op) error
}

// ExecutorFunc is a function that implements the [Executor] interface.
type ExecutorFunc func(ctx context.Context, f Op) error

// Execute implements the [Executor] interface.
func (f ExecutorFunc) Execute(ctx context.Context, op Op) error {
	return f(ctx, op)
}
