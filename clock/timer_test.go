package clock_test

import (
	"testing"
	"time"
)

func TestClockTimer(_ *testing.T) {
	const after = time.Second

	c, s := newTestClock()

	timer := c.Timer(after)
	s.Add(after)
	<-timer.C()
	timer.Reset(after)
	timer.Stop()
}
