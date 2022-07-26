package tests

import (
	"testing"
	"time"
)

type failTestWriter struct {
	t *testing.T
}

func (w failTestWriter) Write(p []byte) (int, error) {
	w.t.Fatal(string(p))
	return len(p), nil
}

func (w failTestWriter) Sync() error {
	return nil
}

type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	return c.now
}

func (*fakeClock) NewTicker(time.Duration) *time.Ticker {
	panic("not reachable")
}
