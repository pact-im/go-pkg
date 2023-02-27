package process

import (
	"context"
	"fmt"
)

// Callback is a type alias for the function passed to Runner’s Run method. See
// [Runner] documentation for more details.
type Callback = func(context.Context) error

// Runner defines an entrypoint for the process that should be considered
// running and ready for use when the Run’s method callback is called.
//
// Runner performs graceful shutdown when the callback returns. In this case,
// the context is not canceled and the process may access external resources,
// e.g. perform some network requests and persist state to database. Otherwise
// a context cancellation is used to force shutdown.
type Runner interface {
	// Run executes the process. Callback function is called when the
	// process has been initialized and began execution. The process is
	// interrupted if the callback returns or the given context expires.
	// On termination the context passed to the callback is canceled.
	//
	// To make debugging easier, callback runs in the same goroutine where
	// Run was invoked. Note though that this behavior currently is not
	// supported by all existing Run implementations.
	//
	// It is safe to assume that callback is called at most once per Run
	// invocation.
	Run(ctx context.Context, callback Callback) error
}

// RunnerFunc is a function that implements the Runner interface.
type RunnerFunc func(context.Context, Callback) error

// Run implements the Runner interface.
func (f RunnerFunc) Run(ctx context.Context, callback Callback) error {
	return f(ctx, callback)
}

// nopRunner is a no-op Runner implementation.
type nopRunner struct{}

// Nop returns a Runner instance that performs no operations and returns when
// the callback does.
func Nop() Runner {
	return (*nopRunner)(nil)
}

// Run implements the Runner interface.
func (*nopRunner) Run(ctx context.Context, callback Callback) error {
	return callback(ctx)
}

// prefixedErrorRunner is a Runner implementation that prefixes an error with
// the set string on failure.
type prefixedErrorRunner struct {
	runner Runner
	prefix string
}

// PrefixedError returns a process that returns an error prefixed with the given
// prefix string on failure.
func PrefixedError(prefix string, runner Runner) Runner {
	return &prefixedErrorRunner{
		runner: runner,
		prefix: prefix,
	}
}

// Run implements the Runner interface.
func (p *prefixedErrorRunner) Run(ctx context.Context, callback Callback) error {
	err := p.runner.Run(ctx, callback)
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", p.prefix, err)
}
