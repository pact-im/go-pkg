package process

import (
	"context"
)

// Runnable defines a blocking long-running process that that should be
// considered running and ready for use when the Run’s method callback is
// called.
type Runnable interface {
	// Run executes the process. The given callback f is called when the
	// process has been initialized and began execution. The process is
	// interrupted if the callback returns or the given context expires.
	// On termination the context passed to the callback is canceled.
	//
	// To make debugging easier, the given f function runs in the same
	// goroutine where Run was invoked. Note though that this behavior
	// currently is not supported by all existing Run implementations.
	//
	Run(ctx context.Context, f func(ctx context.Context) error) error
}

// RunnableFunc is a function that implements the Runnable interface.
type RunnableFunc func(ctx context.Context, f func(ctx context.Context) error) error

// Run implements the Runnable interface.
func (r RunnableFunc) Run(ctx context.Context, f func(ctx context.Context) error) error {
	return r(ctx, f)
}
