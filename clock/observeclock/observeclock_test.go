package observeclock

import (
	"testing"
	"time"

	"go.pact.im/x/clock"
	"go.pact.im/x/clock/fakeclock"
)

func TestObserve(t *testing.T) {
	c := New(clock.NewClock(fakeclock.Go()))
	observer := c.Observe()
	if isClosed(observer) {
		t.Fatal("observation without an event")
	}

	timer := c.Timer(time.Second)
	if !isClosed(observer) {
		t.Fatal("no observation on event")
	}
	timer.Stop()

	observer = c.Observe()
	if isClosed(observer) {
		t.Fatal("observation without an event")
	}

	ticker := c.Ticker(time.Second)
	if !isClosed(observer) {
		t.Fatal("no observation on event")
	}
	ticker.Stop()
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
	}
	return false
}
