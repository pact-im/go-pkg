package process

import (
	"context"
	"fmt"
)

// Callback is a type alias for the function passed to Runnable’s Run method.
// See Runnable documentation for more details.
type Callback = func(context.Context) error

// Runnable defines a blocking long-running process that that should be
// considered running and ready for use when the Run’s method callback is
// called.
//
// Runnable performs graceful shutdown when the callback returns. In this case,
// the context is not canceled and the process may access external resources,
// e.g. perform some network requests and persist state to database. Otherwise
// a context cancellation is used to force shutdown.
type Runnable interface {
	// Run executes the process. Callback function is called when the
	// process has been initialized and began execution. The process is
	// interrupted if the callback returns or the given context expires.
	// On termination the context passed to the callback is canceled.
	//
	// To make debugging easier, callback runs in the same goroutine where
	// Run was invoked. Note though that this behavior currently is not
	// supported by all existing Run implementations.
	//
	// It is safe to assume that callback is called at most once.
	Run(ctx context.Context, callback Callback) error
}

// RunnableFunc is a function that implements the Runnable interface.
type RunnableFunc func(ctx context.Context, callback Callback) error

// Run implements the Runnable interface.
func (f RunnableFunc) Run(ctx context.Context, callback Callback) error {
	return f(ctx, callback)
}

// nopRunnable is a no-op Runnable implementation.
type nopRunnable struct{}

// Nop returns a Runnable instance that performs no operations and returns when
// the callback does.
func Nop() Runnable {
	return (*nopRunnable)(nil)
}

// Run implements the Runnable interface.
func (*nopRunnable) Run(ctx context.Context, callback Callback) error {
	return callback(ctx)
}

// namedRunnable is a Runnable implementation that prefixes an error with the
// set string.
type namedRunnable struct {
	proc Runnable
	name string
}

// Named returns a process that returns an error prefixed with name on failure.
func Named(name string, p Runnable) Runnable {
	return &namedRunnable{
		proc: p,
		name: name,
	}
}

// Run implements the Runnable interface.
func (p *namedRunnable) Run(ctx context.Context, callback Callback) error {
	err := p.proc.Run(ctx, callback)
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", p.name, err)
}
