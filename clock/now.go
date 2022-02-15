package clock

import "time"

// NowScheduler is the interface implemented by a clock that provides an
// optimized implementation of Now.
type NowScheduler interface {
	Scheduler

	// Now returns the current local time.
	Now() time.Time
}

// Now returns the current local time.
//
// Unless the underlying Scheduler implements Timer, Now calls Timer with zero
// duration and uses C on the returned timer to wait for the current time.
func (c *Clock) Now() time.Time {
	s := c.sched()
	if s, ok := s.(NowScheduler); ok {
		return s.Now()
	}

	t := c.Timer(0)
	return <-t.C()
}
