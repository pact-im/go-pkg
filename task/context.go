package task

import (
	"context"
	"sync"
)

// Contextify adds context cancellation support to the task using the given
// shutdown function.
func Contextify(task func() error, shutdown func()) Task {
	return func(ctx context.Context) error {
		done := make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			select {
			case <-done:
			case <-ctx.Done():
				shutdown()
			}
		}()

		err := task()

		close(done)
		wg.Wait()

		return err
	}
}
