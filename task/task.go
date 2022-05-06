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

// Sequential returns a function that executes a group of tasks sequentially
// using the given cancellation condition. The resulting Task returns combined
// errors for all failed subtasks.
func Sequential(cond CancelCondition, tasks ...Task) Task {
	return func(ctx context.Context) error {
		var errs []error

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for i, task := range tasks {
			err := task.Run(ctx)
			if err != nil {
				if errs == nil {
					errs = make([]error, len(tasks))
				}
				errs[i] = err
			}
			if cond(err) {
				cancel()
			}
		}

		return multierr.Combine(errs...)
	}
}

// Parallel returns a function that executes a group of tasks in parallel using
// the given cancellation condition. The resulting Task returns combined errors
// for all failed subtasks.
func Parallel(cond CancelCondition, tasks ...Task) Task {
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
				err := task.Run(ctx)
				if err != nil {
					once.Do(func() { errs = make([]error, len(tasks)) })
					errs[i] = err
				}
				if cond(err) {
					cancel()
				}
			}()
		}
		wg.Wait()

		return multierr.Combine(errs...)
	}
}
