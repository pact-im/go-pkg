package flaky

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestRandomJitter(t *testing.T) {
	j := RandomJitter(func(n int64) int64 {
		assert.Assert(t, n > 0, "Int63n must not be called with n <= 0")
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
		assert.Check(t, is.Equal(tc.Output, d))
	}
}
