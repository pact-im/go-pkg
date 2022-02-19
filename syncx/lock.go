package syncx

import (
	"context"
)

// Lock is a mutually exclusive lock that supports cancelling acquire operation.
type Lock chan struct{}

// NewLock returns a new Lock that allows at most one goroutine to acquire it.
func NewLock() Lock {
	return make(Lock, 1)
}

// Acquire locks c. If the lock is already in use, the calling goroutine blocks
// until either the lock is available or the context expires. It returns nil
// if the lock was acquired and a non-nil context error otherwise.
func (c Lock) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c <- struct{}{}:
		return nil
	}
}

// Release unlocks c. Unlike sync.Mutex, it is valid to release a Lock without
// a corresponding Acquire. In that case the next Acquire call will unlock the
// Release.
func (c Lock) Release() {
	<-c
}
