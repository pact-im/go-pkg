// Package clock defines basic interfaces to a clock. A clock can be provided by
// the host operating system but also by other packages.
//
// A clock provides a capability to schedule an arbitrary event to run after at
// least the given duration. On expiry, the event calls a function in its own
// goroutine with the current time argument. All other functions are implemented
// in terms of scheduled events or using an optimized Scheduler implementations.
//
// A time.Now function signature should be accepted instead of Clock instance
// if a limited capability is required that allows getting the current time but
// not scheduling new events. It is possible to implement Schedule, Timer and
// Ticker in terms of each other interchangeably and hence there is no option to
// limit the scope for these operations.
package clock

import "time"

// A Scheduler allows scheduling events to run after duration elapses.
//
// The Scheduler interface is the minimum implementation required of the clock.
// A clock may implement additional interfaces, such as TimerScheduler, to provide
// additional or optimized functionality.
//
// It is a low-level interface provided by clock implementations. Users should
// accept Clock instance instead of Scheduler.
type Scheduler interface {
	// Schedule schedules an event that calls f in its own goroutine when
	// the given duration elapses. It returns an Event that can be canceled
	// using the Stop method.
	Schedule(d time.Duration, f func(now time.Time)) Event
}

// The Event type represents a single event. When the Event expires, the
// function is called in its own goroutine.
type Event interface {
	// Stop prevents the Event from firing. It returns true if the call
	// stops the timer, false if the timer has been stopped or already
	// expired and the function has been started in its own goroutine.
	//
	// Stop does not wait for function to complete before returning.
	Stop() bool

	// Reset changes the timer to expire after duration d. It returns true
	// if the timer had been active, false if the timer had expired or been
	// stopped.
	//
	// It either reschedules when f will run, in which case Reset returns
	// true, or schedules f to run again, in which case it returns false.
	//
	// Note that Reset does not guarantee that the subsequent goroutine
	// running the function does not run concurrently with the prior one.
	Reset(d time.Duration) bool
}

// A Clock provides a functionality for measuring and displaying time.
type Clock struct {
	Scheduler
}

// NewClock returns a new Clock that uses the given Scheduler to schedule events
// and measure time.
func NewClock(s Scheduler) *Clock {
	return &Clock{
		Scheduler: s,
	}
}

// Schedule implements the Scheduler interface. It uses the clock provided by
// the host operating system if c is nil or the underlying Scheduler is not set.
func (c *Clock) Schedule(d time.Duration, f func(now time.Time)) Event {
	return c.sched().Schedule(d, f)
}

// sched returns the Scheduler implementation for this clock.
func (c *Clock) sched() Scheduler {
	if c == nil {
		return newRuntimeClock()
	}
	s := c.Scheduler
	if s == nil {
		return newRuntimeClock()
	}
	return s
}
