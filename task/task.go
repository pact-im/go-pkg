// Package task provides an alternative to errgroup package with builtin context
// cancellation support.
//
// This packages uses go.uber.org/multierr to combine errors from failed tasks.
package task

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/multierr"
)

// Task is a function that performs some work in foreground.
type Task func(ctx context.Context) error

// Run executes the task.
func (t Task) Run(ctx context.Context) error {
	return t(ctx)
}

// Named returns a task that returns an error prefixed with name on failure.
func Named(name string, run Task) Task {
	return func(ctx context.Context) error {
		err := run(ctx)
		if err == nil {
			return nil
		}
		return fmt.Errorf("%s: %w", name, err)
	}
}

// ErrorGroup returns a function that executes a group of tasks until the first
// non-nil error.
//
// If a task exits, a context is canceled iff error is not nil. This is the same
// behavior as the errgroup package.
//
// The resulting function returns errors for all failed tasks.
func ErrorGroup(tasks ...Task) Task {
	const cancelOnNilError = false
	return group(cancelOnNilError, tasks...)
}

// ExitGroup returns a function that executes a group of tasks until any task
// exits.
//
// If a task exits, a context is canceled and all other tasks are terminated.
//
// The resulting function returns errors for all failed tasks.
func ExitGroup(tasks ...Task) Task {
	const cancelOnNilError = true
	return group(cancelOnNilError, tasks...)
}

func group(cancelOnNilError bool, tasks ...Task) Task {
	return func(ctx context.Context) error {
		var once sync.Once
		var errs []error

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(len(tasks))
		for i, task := range tasks {
			i, task := i, task
			go func() {
				defer wg.Done()
				if err := task.Run(ctx); err != nil {
					once.Do(func() { errs = make([]error, len(tasks)) })
					errs[i] = err
				} else if !cancelOnNilError {
					return
				}
				cancel()
			}()
		}
		wg.Wait()

		return multierr.Combine(errs...)
	}
}
