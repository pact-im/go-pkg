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
