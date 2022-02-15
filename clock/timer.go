package clock

import "time"

// TimerScheduler is the interface implemented by a clock that provides an
// optimized implementation of ScheduleTimer.
type TimerScheduler interface {
	Scheduler

	// Timer creates a new Timer that will send the current time on
	// its channel after at least duration d.
	Timer(d time.Duration) Timer
}

// The Timer type represents a single event. When the Timer expires, the
// current time will be sent on C.
type Timer interface {
	// C returns the channel on which the current time is sent when the
	// timer expires.
	C() <-chan time.Time

	// Stop prevents the Timer from firing. It returns true if the call stops the
	// timer, false if the timer has already expired or been stopped. Stop does not
	// close the channel, to prevent a read from the channel succeeding incorrectly.
	//
	// To ensure the channel is empty after a call to Stop, check the return value
	// and drain the channel. For example, assuming the program has not received
	// from t.C already:
	//
	// 	if !t.Stop() {
	// 		<-t.C()
	// 	}
	//
	// This cannot be done concurrent to other receives from the Timer’s channel or
	// other calls to the Timer’s Stop method.
	Stop() bool

	// Reset changes the timer to expire after duration d. Reset should be
	// invoked only on stopped or expired timers with drained channels.
	//
	// If a program has already received a value from t.C, the timer is known to
	// have expired and the channel drained, so t.Reset can be used directly. If a
	// program has not yet received a value from t.C, however, the timer must be
	// stopped and—if Stop reports that the timer expired before being stopped—the
	// channel explicitly drained:
	//
	// 	select {
	// 	case <-t.C():
	// 		// Timer expired.
	// 	case <-done:
	// 		if !t.Stop() {
	// 			<-t.C()
	// 		}
	// 	}
	// 	t.Reset(d)
	//
	// Reset should always be invoked on stopped or expired channels, as
	// described above. This should not be done concurrent to other receives
	// from the Timer’s channel.
	Reset(d time.Duration)
}

// Timer creates a new Timer that will send the current time on its channel
// after at least duration d.
//
// Unless the underlying Scheduler implements TimerScheduler, Timer calls
// Schedule and sends the current time on the channel when the Event expires.
func (c *Clock) Timer(d time.Duration) Timer {
	s := c.sched()
	if s, ok := s.(TimerScheduler); ok {
		return s.Timer(d)
	}

	ch := make(chan time.Time, 1)
	event := s.Schedule(d, func(now time.Time) {
		// Do a non-blocking send of current time on ch. Default case
		// is unreachable if the timer is used correctly. See also the
		// docs for Timer’s Reset method.
		select {
		case ch <- now:
		default:
		}
	})
	return &eventTimer{event, ch}
}

// eventTimer is a Timer adapter for Event.
type eventTimer struct {
	e Event
	c chan time.Time
}

// C implements the Timer interface.
func (t *eventTimer) C() <-chan time.Time {
	return t.c
}

// Stop implements the Timer interface.
func (t *eventTimer) Stop() bool {
	return t.e.Stop()
}

// Reset implements the Timer interface.
func (t *eventTimer) Reset(d time.Duration) {
	_ = t.e.Reset(d)
}
