package clock_test

import (
	"testing"
	"time"

	"go.pact.im/x/clock"
	"go.pact.im/x/clock/fakeclock"
)

// bareScheduler hides additional functionality of the underlying Scheduler.
// leaving only the Schedule method.
type bareScheduler struct {
	s clock.Scheduler
}

// newBareScheduler returns a new bareScheduler instance.
func newBareScheduler(s clock.Scheduler) *bareScheduler {
	return &bareScheduler{s}
}

// Schedule implements the Scheduler interface.
func (s *bareScheduler) Schedule(d time.Duration, f func(t time.Time)) clock.Event {
	return s.s.Schedule(d, f)
}

// newTestClock returns a new Clock that uses the fakeclock reduced to signle
// Schedule method.
func newTestClock() (*clock.Clock, *fakeclock.Clock) {
	s := fakeclock.Go()
	c := clock.NewClock(newBareScheduler(s))
	return c, s
}

func TestClockSchedule(_ *testing.T) {
	const after = time.Second

	c, s := newTestClock()

	ch := make(chan time.Time, 1)
	c.Schedule(after, func(now time.Time) {
		ch <- now
	})
	s.Add(after)
	<-ch
}
