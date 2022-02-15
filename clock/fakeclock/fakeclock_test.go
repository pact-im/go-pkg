package fakeclock

import (
	"testing"
	"time"
)

func TestClockNow(t *testing.T) {
	now := time.Date(2049, 5, 6, 23, 55, 11, 1034, time.UTC)
	sim := Time(now)
	sim.AddDate(1, 0, 0)
	if sim.Now().Year() != 2050 {
		t.Fatal("unexpected year after AddDate")
	}
	sim.Set(now)
	if !sim.Now().Equal(now) {
		t.Fatal("failed to go back in time")
	}
	sim.Add(time.Second * 10)
	if sim.Now().Second() != 21 {
		t.Fatal("10 second shift failed")
	}
}

func TestClockNext(t *testing.T) {
	sim := Y2038(time.Nanosecond)
	now := sim.Now()
	next, ok := sim.Next()
	if !next.Equal(now) {
		t.Fatalf("call to Next changed current time from %v to %v when there are no events", now, next)
	}
	if ok {
		t.Fatal("call to Next returned true when there were no events")
	}
}

func TestEvent(t *testing.T) {
	const interval = time.Second

	sim := Go()
	now := sim.Now()

	c := make(chan time.Time, 1)
	e := sim.Schedule(interval, func(t time.Time) {
		c <- t
	})
	expectNoC(t, c)

	sim.Add(interval)
	awaitC(t, c)
	expectC(t, c, now.Add(interval))

	if e.Stop() {
		t.Fatal("prevented an expired event from firing")
	}
}

func TestEventReset(t *testing.T) {
	const (
		initial = time.Second
		resetTo = time.Second
	)

	sim := Go()

	c := make(chan time.Time, 1)
	e := sim.Schedule(initial, func(t time.Time) {
		c <- t
	})
	expectNoC(t, c)

	if !e.Stop() {
		t.Fatal("failed to prevent an event from firing")
	}
	if e.Reset(initial) {
		t.Fatal("rescheduled a stopped event")
	}
	if !e.Reset(initial + resetTo) {
		t.Fatal("failed to reschedule an event")
	}

	sim.Add(initial)
	expectNoC(t, c)

	sim.Add(resetTo)
	awaitC(t, c)
	expectC(t, c, sim.Now())
}

func TestEventStop(t *testing.T) {
	const interval = time.Second

	sim := Go()

	c := make(chan time.Time, 1)
	e := sim.Schedule(interval, func(t time.Time) {
		c <- t
	})
	expectNoC(t, c)

	if !e.Stop() {
		t.Fatal("failed to prevent an event from firing")
	}
}

func TestTimerStep(t *testing.T) {
	const (
		step = time.Second
		n    = 3
	)

	sim := Unix()

	timer := sim.Timer(n * step)

	c := timer.C()
	expectNoC(t, c)

	for i := 1; i < n; i++ {
		sim.Add(step)
		expectNoC(t, c)
	}
	sim.Add(step)
	expectC(t, c, sim.Now())
}

func TestTimer(t *testing.T) {
	sim := Go()

	timer := sim.Timer(0)

	c := timer.C()
	expectNoC(t, c)

	sim.Add(0)
	expectC(t, c, sim.Now())

	timer.Reset(time.Nanosecond)
	expectNoC(t, c)

	sim.Add(time.Nanosecond)
	expectC(t, c, sim.Now())
}

func TestTimerNext(t *testing.T) {
	sim := Go()

	timer := sim.Timer(-time.Second)

	c := timer.C()
	expectNoC(t, c)

	now := sim.Now()
	next, ok := sim.Next()
	if !next.Equal(now) {
		t.Fatalf("call to Next changed current time from %v to %v when there are events in the past", now, next)
	}
	if !ok {
		t.Fatal("call to Next returned false when there are scheduled events")
	}
	expectC(t, c, now)

	timer.Reset(time.Second)
	expectNoC(t, c)
	expectNext(t, sim, c, time.Second)
}

func TestTimerNonBlocking(t *testing.T) {
	const interval = time.Second

	sim := Go()

	timer := sim.Timer(interval)

	c := timer.C()
	expectNoC(t, c)

	sim.Add(interval)
	past := sim.Now()

	timer.Reset(interval)

	sim.Add(interval)
	expectC(t, c, past)
}

func TestTimerReset(t *testing.T) {
	sim := Go()

	timer := sim.Timer(time.Second)

	c := timer.C()
	expectNoC(t, c)

	if !timer.Stop() {
		t.Fatal("expected Stop to return false if it prevented the timer from firing")
	}

	timer.Reset(time.Second)
	expectNoC(t, c)

	sim.Add(time.Second)
	expectC(t, c, sim.Now())

	if timer.Stop() {
		t.Fatal("expected Stop to return true if it did not prevent the timer from firing")
	}
}

func TestTicker(t *testing.T) {
	const (
		interval = time.Second
		count    = 3
	)

	sim := Go()
	end := sim.Now().Add(count * interval)

	ticker := sim.Ticker(interval)

	c := ticker.C()
	expectNoC(t, c)

	for i := 0; i < count; i++ {
		expectNext(t, sim, c, interval)
	}
	if now := sim.Now(); !now.Equal(end) {
		t.Fatalf("expected %v time, got %v", end, now)
	}
}

func TestTickerNonBlocking(t *testing.T) {
	const interval = time.Second

	sim := Go()

	ticker := sim.Ticker(interval)

	c := ticker.C()
	expectNoC(t, c)

	sim.Add(interval)
	past := sim.Now()

	sim.Add(interval)
	expectC(t, c, past)
}

func TestTickerStop(t *testing.T) {
	const interval = time.Nanosecond

	sim := Go()

	ticker := sim.Ticker(interval)

	c := ticker.C()
	expectNoC(t, c)

	sim.Add(interval)
	expectC(t, c, sim.Now())

	ticker.Stop()
	expectNoC(t, c)

	sim.Add(interval)
	expectNoC(t, c)

	ticker.Reset(interval)
	expectNoC(t, c)
	sim.Add(interval)
	expectC(t, c, sim.Now())
}

func TestTickerPanics(t *testing.T) {
	sim := Go()

	t.Run("NewTicker", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Clock.Ticker(-1) should have panicked")
			}
		}()
		sim.Ticker(-1)
	})
	t.Run("TickerReset", func(t *testing.T) {
		ticker := sim.Ticker(time.Second)
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Ticker.Reset(0) should have panicked")
			}
		}()
		ticker.Reset(0)
	})
}

func expectNoC(t *testing.T, c <-chan time.Time) {
	t.Helper()

	if len(c) == 0 {
		return
	}
	actual := <-c
	t.Fatalf("expected empty time channel, but got %q", actual)
}

func expectC(t *testing.T, c <-chan time.Time, expected time.Time) {
	t.Helper()

	if len(c) != 1 {
		t.Fatal("time channel is empty")
	}
	actual := <-c
	if !actual.Equal(expected) {
		t.Fatalf("expected %q, but got %q", expected, actual)
	}
}

func awaitC(t *testing.T, c chan time.Time) {
	t.Helper()

	now := <-c
	select {
	case c <- now:
	default:
		t.Fatal("failed to send value back to the channel")
	}
}

func expectNext(t *testing.T, sim *Clock, c <-chan time.Time, d time.Duration) {
	t.Helper()

	now := sim.Now()
	nextNow := now.Add(d)
	next, ok := sim.Next()
	if !next.Equal(nextNow) {
		if !next.Equal(now) {
			t.Fatalf("call to Next changed current time from %v to %v instead of %v when there are no events in the past", now, sim.Now(), nextNow)
		} else {
			t.Fatalf("call to Next did not change current time from %v to %v when there are no events in the past", now, nextNow)
		}
	}
	if !ok {
		t.Fatal("call to Next returned false when there are scheduled events")
	}
	expectC(t, c, nextNow)
}
