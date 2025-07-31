// Package fakeclock provides support for testing users of a clock.
package fakeclock

import (
	"math"
	"sync"
	"time"

	"go.pact.im/x/clock"
)

var _ interface {
	clock.Scheduler
	clock.NowScheduler
	clock.TimerScheduler
	clock.TickerScheduler
} = (*Clock)(nil)

// moment represents a scheduled event.
type moment interface {
	// next returns the duration until the next occurrence of this event, or
	// false if it is the last event.
	next(now time.Time) (time.Duration, bool)
}

// Clock is a fake clock.Scheduler interface implementation that is safe for
// concurrent use by multiple goroutines.
//
// Use Next, Set, Add and AddDate methods to change clock time. Advancing the
// time triggers scheduled events, timers and tickers.
//
// Note that the order in which events scheduled for the same time are triggered
// is undefined, but it is guaranteed that all events that are not after the new
// current time are triggered on clock time change (even if old time is equal to
// the next time value).
//
// The zero Clock defaults to zero time and is ready for use.
type Clock struct {
	mu    sync.Mutex
	now   time.Time
	sched map[moment]time.Time
}

// Unix returns a clock set to the Unix epoch time. That is, it is set to
// 1970-01-01 00:00:00 UTC.
func Unix() *Clock {
	return Time(time.Unix(0, 0))
}

// Go returns a clock set to the Go initial release date. That is, it is set to
// 2009-11-10 23:00:00 UTC.
func Go() *Clock {
	return Time(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
}

// Y2038 returns a clock set to d duration after the Year 2038 problem time.
// That is, it is set to the given duration after 2038-01-19 03:14:07 UTC, the
// latest time that can be properly encoded as a 32-bit integer that is a number
// of seconds after the Unix epoch.
func Y2038(d time.Duration) *Clock {
	t := time.Unix(math.MaxInt32, 0)
	return Time(t.Add(d))
}

// Time returns a clock set to the given now time.
func Time(now time.Time) *Clock {
	return &Clock{
		now:   now,
		sched: map[moment]time.Time{},
	}
}

// Now returns the current clock time.
func (c *Clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

// Next advances the time to the next timer or ticker event and returns the new
// current time. If there are events in the past, the time is not changed. It
// returns false if there are no scheduled events. Otherwise it returns true and
// runs at least one scheduled event.
func (c *Clock) Next() (time.Time, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ok := len(c.sched) > 0
	c.now = c.next(c.now)
	c.advance(c.now)
	return c.now, ok
}

// Set sets the given time to be the current clock time. It is possible to set
// t that is before the current clock time.
func (c *Clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = t
	c.advance(c.now)
}

// Add adds the given duration to the current time and returns the resulting
// clock time. It is possible to add a negative duration.
//
// It is safe for concurrent use and is a shorthand for
//
//	now := c.Now().Add(d)
//	c.Set(now)
func (c *Clock) Add(d time.Duration) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.Add(d)
	c.advance(c.now)
	return c.now
}

// AddDate adds the duration corresponding to the given number of years, months
// and days relative to the current time and returns the resulting clock time.
// It is possible to add negative values.
//
// It is safe for concurrent use and is a shorthand for
//
//	now := c.Now().AddDate(years, months, days)
//	c.Set(now)
func (c *Clock) AddDate(years, months, days int) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.AddDate(years, months, days)
	c.advance(c.now)
	return c.now
}

// Schedule implements the clock.Scheduler interface.
func (c *Clock) Schedule(d time.Duration, f func(now time.Time)) clock.Event {
	t := &event{
		c: c,
		f: f,
	}
	c.reset(t, d, nil)
	return t
}

// Timer implements the clock.TimerScheduler interface.
func (c *Clock) Timer(d time.Duration) clock.Timer {
	t := &timer{
		c:  c,
		ch: make(chan time.Time, 1),
	}
	c.reset(t, d, nil)
	return t
}

// Ticker implements the clock.TickerScheduler interface. Note that the returned
// ticker does not adjust the time interval or drop ticks to make up for slow
// receivers.
func (c *Clock) Ticker(d time.Duration) clock.Ticker {
	if d <= 0 {
		panic("non-positive interval for Ticker")
	}
	t := &ticker{
		c:  c,
		ch: make(chan time.Time, 1),
	}
	c.reset(t, d, &t.d)
	return t
}

// stop removes the given moment from the set of scheduled events. It returns
// true if stop prevented the event from firing. Note that stop acquires the
// underlying lock.
func (c *Clock) stop(m moment) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.sched[m]
	delete(c.sched, m)
	return ok
}

// reset resets the given moment to run d duration after the current time.
// Note that reset acquires the underlying lock. Reset returns true if it
// rescheduled an event.
func (c *Clock) reset(m moment, d time.Duration, dp *time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if dp != nil {
		*dp = d
	}
	return c.schedule(m, c.now.Add(d))
}

// schedule schedules the given moment to run on next clock advance. It returns
// true if an event was rescheduled.
func (c *Clock) schedule(m moment, when time.Time) bool {
	if c.sched == nil {
		c.sched = map[moment]time.Time{}
	}
	_, ok := c.sched[m]
	c.sched[m] = when
	return ok
}

// next returns the time of the next scheduled event. It returns the given now
// time if there are no events or some of them are in the past.
func (c *Clock) next(now time.Time) time.Time {
	if len(c.sched) == 0 {
		return now
	}
	var next time.Time
	for _, t := range c.sched {
		next = t
		break
	}
	for _, t := range c.sched {
		if !t.After(now) {
			return now
		}
		if t.After(next) {
			continue
		}
		next = t
	}
	return next
}

// advance runs the scheduled events for the current clock time.
func (c *Clock) advance(now time.Time) {
	for m, t := range c.sched {
		if t.After(now) {
			continue
		}
		next, ok := m.next(now)
		if !ok {
			delete(c.sched, m)
			continue
		}
		_ = c.schedule(m, now.Add(next))
	}
}

// event implements the moment and clock.Event interfaces.
type event struct {
	c *Clock
	f func(time.Time)
}

// Stop implements the clock.Event interface.
func (t *event) Stop() bool {
	return t.c.stop(t)
}

// Reset implements the clock.Event interface.
func (t *event) Reset(d time.Duration) bool {
	return t.c.reset(t, d, nil)
}

// next implements the moment interface.
func (t *event) next(now time.Time) (time.Duration, bool) {
	go t.f(now)
	return 0, false
}

// timer implements the moment and clock.Timer interfaces.
type timer struct {
	c  *Clock
	ch chan time.Time
}

// C implements the clock.Timer interface.
func (t *timer) C() <-chan time.Time {
	return t.ch
}

// Stop implements the clock.Timer interface.
func (t *timer) Stop() bool {
	return t.c.stop(t)
}

// Reset implements the clock.Timer interface.
func (t *timer) Reset(d time.Duration) {
	_ = t.c.reset(t, d, nil)
}

// next implements the moment interface.
func (t *timer) next(now time.Time) (time.Duration, bool) {
	select {
	case t.ch <- now:
	default:
	}
	return 0, false
}

// ticker implements the moment and clock.Ticker interfaces.
type ticker struct {
	c  *Clock
	ch chan time.Time
	d  time.Duration
}

// C implements the clock.Ticker interface.
func (t *ticker) C() <-chan time.Time {
	return t.ch
}

// Stop implements the clock.Ticker interface.
func (t *ticker) Stop() {
	_ = t.c.stop(t)
}

// Reset implements the clock.Ticker interface.
func (t *ticker) Reset(d time.Duration) {
	if d <= 0 {
		panic("non-positive interval for Ticker.Reset")
	}
	_ = t.c.reset(t, d, &t.d)
}

// next implements the moment interface.
func (t *ticker) next(now time.Time) (time.Duration, bool) {
	select {
	case t.ch <- now:
	default:
	}
	return t.d, true
}
