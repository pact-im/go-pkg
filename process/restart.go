package process

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	restartInitialWait  = time.Second
	restartLoopInterval = time.Minute
	restartLoopWait     = 0
)

// restartInitial restores the last state from the underlying table.
func (m *Manager[K, P]) restartInitial(ctx context.Context) {
	_ = m.restart(ctx, restartInitialWait) // TODO: log error
}

// spawnRestartLoop spawns a restartLoop using the given context and returns
// a function that stops restart loop and waits completion.
func (m *Manager[K, P]) spawnRestartLoop(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		m.restartLoop(ctx)
	}()

	return func() {
		cancel()
		wg.Wait()
	}
}

// restartLoop runs a loop that calls restart.
func (m *Manager[K, P]) restartLoop(ctx context.Context) {
	const interval = restartLoopInterval

	timer := m.clock.Timer(interval)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C():
		}
		_ = m.restart(ctx, restartLoopWait) // TODO: log error
		timer.Reset(interval)
	}
}

// restart starts processes from the table that are not currently running.
//
// If wait duration is not zero, it waits until background startProcess calls
// complete. This allows ensuring that we restore at least partial state in
// restoreInitial before invoking Manager’s Run callback.
func (m *Manager[K, P]) restart(ctx context.Context, wait time.Duration) error {
	iter, err := m.table.Iter(ctx)
	if err != nil {
		return fmt.Errorf("create iterator: %w", err)
	}
	defer func() {
		if err := iter.Close(); err != nil {
			_ = err // TODO: log error
		}
	}()

	var wg, wg2 sync.WaitGroup
	var timer, stopc chan struct{}

	for iter.Next() {
		pk, r, err := iter.Get(ctx)
		if err != nil {
			_ = err // TODO: log error
			continue
		}
		done := m.restartProcessInBackground(ctx, pk, r)
		if wait > 0 {
			// Note that we cannot use timer.C here since there are multiple
			// consumers (each process we start) and timer sends the value
			// on channel only once (and never closes it).
			if timer == nil {
				timer = make(chan struct{})
				stopc = make(chan struct{})

				wg2.Add(1)
				go func() {
					defer wg2.Done()
					t := m.clock.Timer(wait)
					defer t.Stop()
					select {
					case <-ctx.Done():
						return
					case <-t.C():
						close(timer)
					case <-stopc:
					}
				}()
				defer wg2.Wait()
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
				case <-done:
				case <-timer:
				}
			}()
		}
	}
	if err := iter.Err(); err != nil {
		_ = err // TODO: log error
	}

	// Do not wait for timer if all goroutines under wg have finished. Note
	// that this must be run after the loop since we cannot use wg.Wait
	// concurrently with wg.Add.
	if stopc != nil {
		wg.Wait()
		close(stopc)
	}

	return nil
}

// restartProcessInBackground starts the process for the given key in the
// background. It returns a channel that is closed when startProcess call
// completes.
func (m *Manager[K, P]) restartProcessInBackground(ctx context.Context, pk K, r P) <-chan struct{} {
	done := make(chan struct{})
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer close(done)
		_, _ = m.startProcess(ctx, pk, r) // TODO: log error
	}()
	return done
}
