// Package observeclock provides a clock implementation that allows observing
// Schedule, Timer and Ticker method calls.
//
// The package is intended as an alternative to mockclock for simple cases when
// observing an event creation is enough to reach the desired state.
package observeclock

import (
	"sync"
	"time"

	"go.pact.im/x/clock"
)

var _ interface {
	clock.Scheduler
	clock.TimerScheduler
	clock.TickerScheduler
} = (*Clock)(nil)

// Clock allows observing creation of new events on the underlying clock.Clock
// instance.
type Clock struct {
	*clock.Clock

	mu sync.Mutex
	xs []chan struct{}
}

// NewClock returns a new Clock that observes the given clock.Clock.
func NewClock(c *clock.Clock) *Clock {
	return &Clock{
		Clock: c,
	}
}

// New returns a new Clock that observes the given clock.Scheduler including
// optional interface method calls.
func New(s clock.Scheduler) *Clock {
	return NewClock(clock.NewClock(s))
}

// Schedule implements the clock.Scheduler interface.
func (c *Clock) Schedule(d time.Duration, f func(time.Time)) clock.Event {
	t := c.Clock.Schedule(d, f)
	c.event()
	return t
}

// Timer implements the clock.TimerScheduler interface.
func (c *Clock) Timer(d time.Duration) clock.Timer {
	t := c.Clock.Timer(d)
	c.event()
	return t
}

// Ticker implements the clock.TickerScheduler interface.
func (c *Clock) Ticker(d time.Duration) clock.Ticker {
	t := c.Clock.Ticker(d)
	c.event()
	return t
}

// Observe returns a channel that is closed on Schedule, Timer and Ticker calls.
func (c *Clock) Observe() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	x := make(chan struct{})
	c.xs = append(c.xs, x)
	return x
}

// event triggers an observable event.
func (c *Clock) event() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, x := range c.xs {
		c.xs[i] = nil
		close(x)
	}
	c.xs = c.xs[:0]
}
