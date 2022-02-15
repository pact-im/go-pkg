package flaky

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.pact.im/x/clock"
)

// ErrDebounced is an error that DebounceExecutor returns when an Execute call
// is superseded by another operation.
var ErrDebounced = errors.New("flaky: debounced")

// DebounceExecutor is an executor that debounces an operation. That is, it
// performs an operation once it stops being called for the specified debounce
// duration.
//
// Keep in mind that using DebounceExecutor introduces latency to the operation
// for at least the specified debounce duration. In most cases it should be used
// when the operation is expensive to execute or under high event rate or load.
// As a side effect, DebounceExecutor guarantees that at most one operation is
// executing at a time.
//
type DebounceExecutor struct {
	once sync.Once
	debounceState

	c *clock.Clock
	w time.Duration

	timer clock.Timer
}

// debounceState is the debounce state for DebounceExecutor.
type debounceState struct {
	lock  chan struct{}
	exec  chan struct{}
	steal chan struct{}
	next  chan struct{}
}

// Debounce returns a new DebounceExecutor instance for the given wait duration.
func Debounce(t time.Duration) *DebounceExecutor {
	return &DebounceExecutor{
		w: t,
	}
}

// clone returns a copy of the executor.
func (d *DebounceExecutor) clone() *DebounceExecutor {
	d.init()
	return &DebounceExecutor{
		debounceState: d.debounceState,
		c:             d.c,
		w:             d.w,
		timer:         d.timer,
	}
}

// WithClock returns a copy of the executor that uses the given clock.
func (d *DebounceExecutor) WithClock(c *clock.Clock) *DebounceExecutor {
	d = d.clone()
	d.timer = nil
	d.c = c
	d.init()
	return d
}

// WithWait returns a copy of the executor that uses the given wait duration.
//
// It shares the underlying state and can be used to debounce operation with
// non-default wait duration.
func (d *DebounceExecutor) WithWait(t time.Duration) *DebounceExecutor {
	d = d.clone()
	d.w = t
	d.init()
	return d
}

// Execute calls the function f if and only if Execute is not called again
// during the debounce interval. Context expiration cancels an operation. It
// returns ErrDebounced error if the given operation f was superseded by another
// Execute call. If Execute is called concurrently, the last invocation wins.
// When another operation is already executing, it returns ErrDebounced.
//
// Callers should handle ErrDebounced error to avoid breaking assumptions when
// running under another DebounceExecutor.
//
// As a side effect, it guarantees that at most one f is executing at a time.
// That is, it is safe to avoid locking if the state is mutated exclusively
// under the DebounceExecutor.
func (d *DebounceExecutor) Execute(ctx context.Context, f Op) error {
	d.init()

	// Acquire a lock or steal it from an ongoing debounced Execute call.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case v := <-d.exec:
		d.exec <- v
		return ErrDebounced
	case d.lock <- struct{}{}:
	case d.steal <- struct{}{}:
		// Wait until the Execute call we are stealing from stops the
		// timer and passes the lock to us.
		<-d.next
	}

	// Reset timer for the debounce duration.
	if d.timer == nil {
		d.timer = d.c.Timer(d.w)
	} else {
		d.timer.Reset(d.w)
	}

	// Wait until the debounce timer expiration, another Execute invocation
	// that steals our lock, or the action cancellation.
	var steal bool
	select {
	case <-d.timer.C():
		d.exec <- struct{}{}
		err := f(ctx)
		<-d.exec
		<-d.lock
		return unwrapInternal(err)
	case <-d.steal:
		steal = true
	case <-ctx.Done():
	}
	if !d.timer.Stop() {
		<-d.timer.C()
	}
	if steal {
		d.next <- struct{}{}
		return ErrDebounced
	}

	<-d.lock
	return ctx.Err()
}

// init initializes the internal debouncer state.
func (d *DebounceExecutor) init() {
	d.once.Do(func() {
		if d.c == nil {
			d.c = clock.System()
		}
		if d.lock == nil {
			d.lock = make(chan struct{}, 1)
		}
		if d.exec == nil {
			d.exec = make(chan struct{}, 1)
		}
		if d.steal == nil {
			d.steal = make(chan struct{})
		}
		if d.next == nil {
			d.next = make(chan struct{})
		}
	})
}
