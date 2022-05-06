package process

import (
	"context"
	"sync"
)

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
func Leaf(runInForeground, gracefulStop func(ctx context.Context) error) Runnable {
	return &leafRunnable{runInForeground, gracefulStop}
}

type leafRunnable struct {
	runInForeground func(ctx context.Context) error
	gracefulStop    func(ctx context.Context) error
}

func (r *leafRunnable) Run(ctx context.Context, callback func(ctx context.Context) error) error {
	bgctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	var runError error
	go func() {
		defer wg.Done()
		runError = r.runInForeground(bgctx)
		cancel() // cancel callback
	}()

	callbackError := callback(bgctx)

	var stopError error
	if r.gracefulStop != nil {
		stopError = r.gracefulStop(ctx)
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

// StartStop returns a Runnable instance for the pair of start/stop functions.
// The stop function should perform a graceful shutdown until a context expires,
// then proceed with a forced shutdown.
//
// The resulting Runnable returns either start error or the first non-nil error
// from callback and stop functions. If both callback and stop return a non-nil
// error, the latter is ignored.
func StartStop(startInBackground, gracefulStop func(ctx context.Context) error) Runnable {
	return &startStopRunnable{startInBackground, gracefulStop}
}

type startStopRunnable struct {
	startInBackground func(ctx context.Context) error
	gracefulStop      func(ctx context.Context) error
}

func (r *startStopRunnable) Run(ctx context.Context, callback func(ctx context.Context) error) error {
	if err := r.startInBackground(ctx); err != nil {
		return err
	}
	callbackError := callback(ctx)
	if err := r.gracefulStop(ctx); err != nil && callbackError == nil {
		return err
	}
	return callbackError
}
