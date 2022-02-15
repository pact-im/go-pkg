package clock

import (
	"sync"
	"time"
)

// TickerScheduler is the interface implemented by a clock that provides an
// optimized implementation of Ticker.
type TickerScheduler interface {
	Scheduler

	// Ticker returns a new Ticker containing a channel that will send the time
	// on the channel after each tick. The period of the ticks is specified by the
	// duration argument. The ticker will adjust the time interval or drop ticks to
	// make up for slow receivers. The duration d must be greater than zero; if
	// not, Ticker will panic. Stop the ticker to release associated resources.
	Ticker(d time.Duration) Ticker
}

// A Ticker holds a channel that delivers “ticks” of a clock at intervals.
type Ticker interface {
	// C returns the channel on which the ticks are delivered.
	C() <-chan time.Time
	// Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does
	// not close the channel, to prevent a concurrent goroutine reading from the
	// channel from seeing an erroneous “tick”.
	Stop()
	// Reset stops a ticker and resets its period to the specified duration. The
	// next tick will arrive after the new period elapses. The duration d must be
	// greater than zero; if not, Reset will panic.
	Reset(d time.Duration)
}

// Ticker returns a new Ticker containing a channel that will send the current
// time on the channel after each tick. The period of the ticks is specified by
// the duration argument. The ticker will adjust the time interval or drop ticks
// to make up for slow receivers. The duration d must be greater than zero; if
// not, Ticker will panic. Stop the ticker to release associated resources.
//
// Unless the underlying Scheduler implements TickerScheduler, Ticker calls
// Timer and resets the returned timer on each tick.
func (c *Clock) Ticker(d time.Duration) Ticker {
	if d <= 0 {
		panic("non-positive interval for Ticker")
	}

	s := c.sched()

	if s, ok := s.(TickerScheduler); ok {
		return s.Ticker(d)
	}
	return newEventTicker(c, d)
}

// eventTicker implements Ticker interface using a Timer. It intercepts timer’s
// channel and then resets it to continue ticking.
type eventTicker struct {
	t Timer
	c chan time.Time

	mu sync.Mutex
	wg sync.WaitGroup

	resetc  chan time.Duration
	stopc   chan struct{}
	stopped bool

	resetTimer bool
}

func newEventTicker(c *Clock, d time.Duration) *eventTicker {
	t := &eventTicker{
		t:      c.Timer(d),
		c:      make(chan time.Time, 1),
		resetc: make(chan time.Duration),
		stopc:  make(chan struct{}),
	}
	t.start(d)
	return t
}

// C implements the Ticker interface.
func (t *eventTicker) C() <-chan time.Time {
	return t.c
}

// Stop implements the Ticker interface.
func (t *eventTicker) Stop() {
	t.stop()
}

// Reset implements the Ticker interface.
func (t *eventTicker) Reset(d time.Duration) {
	if d <= 0 {
		panic("non-positive interval for Ticker.Reset")
	}
	t.reset(d)
}

func (t *eventTicker) start(d time.Duration) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.run(d)
	}()
}

func (t *eventTicker) run(d time.Duration) {
	timer := t.t
	timerC := timer.C()
	if t.resetTimer {
		timer.Reset(d)
	}
	t.resetTimer = true

	for {
		select {
		case <-t.stopc:
			t.stopTimer()
			return
		case d = <-t.resetc:
			t.stopTimer()
			timer.Reset(d)
		case now := <-timerC:
			timer.Reset(d)
			select {
			case t.c <- now:
			default:
			}
		}
	}
}

func (t *eventTicker) reset(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		t.stopped = false
		t.start(d)
		return
	}
	t.resetc <- d
}

func (t *eventTicker) stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		return
	}
	t.stopped = true
	t.stopc <- struct{}{}
	t.wg.Wait()
}

func (t *eventTicker) stopTimer() {
	if t.t.Stop() {
		return
	}
	<-t.t.C()
}
