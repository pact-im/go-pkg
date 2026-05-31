package flaky

import (
	"testing"
	"time"
)

func TestRandomJitter(t *testing.T) {
	j := RandomJitter(func(n int64) int64 {
		if n <= 0 {
			t.Fatalf("Int63n must not be called with n <= 0, got %d", n)
		}
		return n - 1
	})
	testCases := []struct {
		Interval JitterInterval
		Output   time.Duration
	}{
		{
			Interval: JitterInterval{
				L: 0,
			},
			Output: 0,
		},
		{
			Interval: JitterInterval{
				L: time.Second,
			},
			Output: time.Second,
		},
		{
			Interval: JitterInterval{
				L: -time.Second,
			},
			Output: -time.Second,
		},
		{
			Interval: JitterInterval{
				L: time.Second,
				Q: 2,
			},
			Output: 500 * time.Millisecond,
		},
		{
			Interval: JitterInterval{
				L: time.Second,
				Q: -4,
			},
			Output: 1250 * time.Millisecond,
		},
		{
			Interval: JitterInterval{
				L: -time.Hour,
				Q: -6,
			},
			Output: -70 * time.Minute,
		},
	}
	for _, tc := range testCases {
		d := j(tc.Interval)
		if d != tc.Output {
			t.Errorf("expected %v, got %v", tc.Output, d)
		}
	}
}
