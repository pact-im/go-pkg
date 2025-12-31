// Package supervisor provides [process.Runner] supervision implementation.
package supervisor

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"go.pact.im/x/flaky"
	"go.pact.im/x/process"
)

// errInterrupt is an internal error used to distinguish between supervisor
// interrupts and other errors.
var errInterrupt = errors.New("supervisor: interrupt")

// errRecursiveOrConcurrentRun is an error that [Supervisor] returns on
// recursive or concurrent Run call.
var errRecursiveOrConcurrentRun = errors.New("supervisor: recursive or concurrent Supervisor.Run calls are not allowed")

// Supervisor runs a [process.Runner] alongside a control callback, managing
// their concurrent execution and coordinated shutdown. It wraps the runner
// with [flaky.Executor] retry logic and pre/post execution hooks.
type Supervisor struct {
	runner process.Runner
	exec   flaky.Executor
	hook   Hook

	intr atomic.Pointer[supervisorInterrupter]
}

// Hook is a set of hooks for supervisor’s runner.
type Hook struct {
	// Pre is a function that is called on runner’s callback. A non-nil
	// error is immediately returned from the callback. In that case, an
	// error from PostHook is ignored.
	// Defaults to a function that returns nil error.
	Pre func(context.Context, *Supervisor) error

	// Post is a function that is called before runner’s callback returns.
	// The result is returned from the callback.
	// Defaults to a function that returns nil error.
	Post func(context.Context, *Supervisor) error
}

// NewSupervisor returns a new [Supervisor] instance for the given runner.
func NewSupervisor(runner process.Runner, exec flaky.Executor, hook Hook) *Supervisor {
	return &Supervisor{
		runner: runner,
		exec:   exec,
		hook:   hook,
	}
}

// Interrupt signals the supervisor to stop the current execution. This method
// is non-blocking and returns immediately; it does not wait for execution to
// complete.
//
// If no execution is active ([Supervisor.Run] is not being called), this method
// does nothing. Otherwise, it signals the process to stop, but the caller must
// separately wait for Run to return if synchronization is needed.
//
// It is safe to call from multiple goroutines concurrently.
func (s *Supervisor) Interrupt() {
	intr := s.intr.Load()
	if intr == nil {
		return
	}
	intr.Interrupt()
}

// Run executes the supervisor’s managed process concurrently with the provided
// callback function without waiting for the initial process startup.
//
// Use [Hook.Pre] and [Hook.Post] to run code in runner’s callback.
//
// It returns an error combining any errors from the callback and runner execution.
//
// Execution timeline:
//   - T0: Start runner under executor in background goroutine.
//   - T0: Start callback in current goroutine.
//   - T1: Callback returns → interrupt executor.
//   - T1: Executor returns → cancel callback context.
//   - T2: Wait for executor to complete cleanup → return combined errors.
//
// Only one active Run invocation is permitted per Supervisor instance.
// Concurrent or recursive calls will return an error.
func (s *Supervisor) Run(ctx context.Context, callback process.Callback) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	intr := &supervisorInterrupter{
		done:   make(chan struct{}),
		cancel: cancel,
	}
	if !s.intr.CompareAndSwap(nil, intr) {
		return errRecursiveOrConcurrentRun
	}
	defer s.intr.Store(nil)

	var wg sync.WaitGroup
	var executeError error
	wg.Go(func() {
		defer cancel() // cancel callback
		executeError = s.execute(ctx, intr)
	})

	callbackError := callback(ctx)

	intr.Interrupt()
	wg.Wait()

	// Do not wrap errors in errors.joinError unless necessary.
	switch {
	case callbackError == nil:
		return executeError
	case executeError == nil:
		return callbackError
	}
	return errors.Join(callbackError, executeError)
}

// execute runs the supervised process under the executor with pre/post Run hooks.
func (s *Supervisor) execute(ctx context.Context, intr *supervisorInterrupter) error {
	err := s.exec.Execute(ctx, func(ctx context.Context) error {
		if intr.shouldStopBeforeRunner() {
			return flaky.Internal(errInterrupt)
		}

		err := s.runner.Run(ctx, func(ctx context.Context) error {
			if err := s.pre(ctx); err != nil {
				_ = s.post(ctx)
				return err
			}
			select {
			case <-ctx.Done():
			case <-intr.done:
			}
			return s.post(ctx)
		})

		intr.afterRunner()

		return err
	})
	if errors.Is(err, errInterrupt) {
		err = nil
	}
	return err
}

// pre runs pre hook of the supervisor.
func (s *Supervisor) pre(ctx context.Context) error {
	if s.hook.Pre == nil {
		return nil
	}
	return s.hook.Pre(ctx, s)
}

// pre runs post hook of the supervisor.
func (s *Supervisor) post(ctx context.Context) error {
	if s.hook.Post == nil {
		return nil
	}
	return s.hook.Post(ctx, s)
}
