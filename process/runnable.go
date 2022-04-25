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
// a context cancellation is used to force shutdown.
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

// Leaf converts a “leaf” function to a runnable process function that accepts
// callback. It accepts an optional gracefulStop function to perform graceful
// shutdown. If the function is nil, the process will be terminated by context
// cancellation instead.
//
// The resulting Runnable returns first non-nil error from functions in the
// following order: callback, gracefulStop, runInForeground. That is, if both
// callback and gracefulStop return a non-nil error, the latter is ignored.
//
// Example (HTTP):
//
//  var lis net.Listener
//  var srv *http.Server
//
//  process.Leaf(
//    func(_ context.Context) error {
//      err := srv.Serve(lis)
//      if errors.Is(err, http.ErrServerClosed) {
//        return nil
//      }
//      return err
//    },
//    func(ctx context.Context) error {
//      err := srv.Shutdown(ctx)
//      if err != nil {
//        return srv.Close()
//      }
//      return nil
//    },
//  )
//
// Example (gRPC):
//
//  var lis net.Listener
//  var srv *grpc.Server
//
//  process.Leaf(
//    func(ctx context.Context) error {
//      return srv.Serve(lis)
//    },
//    func(_ context.Context) error {
//      done := make(chan struct{})
//      go func() {
//        srv.GracefulStop()
//        close(done)
//      }()
//      select {
//      case <-ctx.Done():
//        srv.Stop()
//        <-done
//      case <-done:
//      }
//      return nil
//    },
//  )
//
func Leaf(runInForeground, gracefulStop func(ctx context.Context) error) RunnableFunc {
	return func(ctx context.Context, callback func(ctx context.Context) error) error {
		bgctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(1)

		var runError error
		go func() {
			defer wg.Done()
			runError = runInForeground(bgctx)
			cancel() // cancel callback
		}()

		callbackError := callback(bgctx)

		var stopError error
		if gracefulStop != nil {
			stopError = gracefulStop(ctx)
		}

		cancel()
		wg.Wait()

		switch {
		case callbackError != nil:
			return callbackError
		case stopError != nil:
			return stopError
		}
		return runError
	}
}

// StartStop returns a Runnable instance for the pair of start/stop functions.
// The stop function should perform a graceful shutdown until a context expires,
// then proceed with a forced shutdown.
//
// The resulting Runnable returns either start error or the first non-nil error
// from callback and stop functions. If both callback and stop return a non-nil
// error, the latter is ignored.
func StartStop(startInBackground, gracefulStop func(ctx context.Context) error) RunnableFunc {
	return func(ctx context.Context, callback func(ctx context.Context) error) error {
		if err := startInBackground(ctx); err != nil {
			return err
		}
		callbackError := callback(ctx)
		if err := gracefulStop(ctx); err != nil && callbackError == nil {
			return err
		}
		return callbackError
	}
}
