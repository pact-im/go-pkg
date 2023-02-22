package clock_test

import (
	"sync"
	"testing"
	"time"

	"go.pact.im/x/clock"
)

// instantScheduler is a Scheduler implementation that immediately fires
// scheduled events and advances current time by the given event delay.
type instantScheduler struct {
	mu  sync.Mutex
	now time.Time
}

// Schedule implements the Scheduler interface.
func (c *instantScheduler) Schedule(d time.Duration, f func(t time.Time)) clock.Event {
	c.mu.Lock()
	next := c.now.Add(d)
	c.now = next
	c.mu.Unlock()
	go f(next)
	return &instantEvent{c, f}
}

// instantEvent is the Event implementation that always returns false on Stop.
type instantEvent struct {
	c *instantScheduler
	f func(t time.Time)
}

// Stop implements the Event interface.
func (e *instantEvent) Stop() bool {
	return false
}

// Reset implements the Event interface.
func (e *instantEvent) Reset(d time.Duration) bool {
	e.c.Schedule(d, e.f)
	return false
}

func TestClockNow(t *testing.T) {
	now := time.Unix(0, 0)
	c := clock.NewClock(&instantScheduler{
		now: now,
	})
	if !now.Equal(c.Now()) {
		t.FailNow()
	}
}
