package clock

import "time"

var _ interface {
	Scheduler
	NowScheduler
	TimerScheduler
	TickerScheduler
} = (*runtimeClock)(nil)

// systemClock is the Clock instance with runtimeClock Scheduler implementation.
var systemClock = Clock{newRuntimeClock()}

// System returns a clock provided by the host operating system.
func System() *Clock {
	return &systemClock
}

// runtimeClock is a clock provided by the host operating system and Go runtime.
type runtimeClock struct{}

// newRuntimeClock returns a new runtimeClock instance that satisfies Scheduler.
func newRuntimeClock() Scheduler {
	return (*runtimeClock)(nil)
}

// Scheduler implements the Scheduler interface.
func (c *runtimeClock) Schedule(d time.Duration, f func(now time.Time)) Event {
	return time.AfterFunc(d, func() {
		f(time.Now())
	})
}

// Now implements the NowScheduler interface.
func (c *runtimeClock) Now() time.Time {
	return time.Now()
}

// Timer implements the TimerScheduler interface.
func (c *runtimeClock) Timer(d time.Duration) Timer {
	return runtimeTimer{time.NewTimer(d)}
}

// Ticker implements the TickerScheduler interface.
func (c *runtimeClock) Ticker(d time.Duration) Ticker {
	return runtimeTicker{time.NewTicker(d)}
}

// runtimeTimer is a Timer adapter for time.Timer type.
type runtimeTimer struct {
	*time.Timer
}

// C implements the Timer interface.
func (r runtimeTimer) C() <-chan time.Time {
	return r.Timer.C
}

// Reset implements the Timer interface.
func (r runtimeTimer) Reset(d time.Duration) {
	_ = r.Timer.Reset(d)
}

// runtimeTicker is a Ticker adapter for time.Ticker type.
type runtimeTicker struct {
	*time.Ticker
}

// C implements the Ticker interface.
func (r runtimeTicker) C() <-chan time.Time {
	return r.Ticker.C
}
