package process

import (
	"context"
	"sync"
)

// Runnable defines a blocking long-running process that that should be
// considered running and ready for use when the Run’s method callback is
// called.
//
// Runnable performs graceful shutdown when the callback returns. In this case,
// the context is not canceled and the process may access external resources,
// e.g. perform some network requests and persist state to database. Otherwise
// a context cancellation may be used to force shutdown.
type Runnable interface {
	// Run executes the process. The given callback f is called when the
	// process has been initialized and began execution. The process is
	// interrupted if the callback returns or the given context expires.
	// On termination the context passed to the callback is canceled.
	//
	// To make debugging easier, the given f function runs in the same
	// goroutine where Run was invoked. Note though that this behavior
	// currently is not supported by all existing Run implementations.
	Run(ctx context.Context, f func(ctx context.Context) error) error
}

// RunnableFunc is a function that implements the Runnable interface.
type RunnableFunc func(ctx context.Context, f func(ctx context.Context) error) error

// Run implements the Runnable interface.
func (r RunnableFunc) Run(ctx context.Context, f func(ctx context.Context) error) error {
	return r(ctx, f)
}

// Leaf converts a “leaf” function to a runnable processs function that accepts
// callback. Note that the resulting Runnable does not support graceful shutdown
// and is terminated by context cancellation. It returns an error from the first
// function to return. If the first error is nil, the returned error is nil even
// if the second error is non-nil. The order in which run and callback functions
// return is not determenistic when the parent context is canceled.
func Leaf(run func(ctx context.Context) error) RunnableFunc {
	return func(ctx context.Context, f func(ctx context.Context) error) error {
		ctx, cancel := context.WithCancel(ctx)

		var once sync.Once
		var err error

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			runError := run(ctx)
			once.Do(func() { err = runError })
			cancel()
			wg.Done()
		}()

		callbackError := f(ctx)
		once.Do(func() { err = callbackError })
		cancel()
		wg.Wait()

		return err
	}
}

// StartStop returns a Runnable instance for the pair of start/stop functions.
//
// The resulting Runnable returns either start error or the first non-nil error
// from callback and stop functions. If both callback and stop return a non-nil
// error, the latter is ignored.
func StartStop(start, stop func(ctx context.Context) error) RunnableFunc {
	return func(ctx context.Context, f func(ctx context.Context) error) error {
		if err := start(ctx); err != nil {
			return err
		}
		callbackError := f(ctx)
		if err := stop(ctx); err != nil && callbackError == nil {
			return err
		}
		return callbackError
	}
}
