package flaky

import (
	"math/rand"
	"time"
)

// JitterProvider provides a jitter for executing operations.
type JitterProvider interface {
	Jitter() Jitter
}

// Jitter is a function that, given a closed interval v, returns a random
// value on that interval.
type Jitter func(v JitterInterval) time.Duration

// RandomJitter returns a Jitter function for function f that defaults to
// math/rand.Int63n if it is nil.
func RandomJitter(f func(n int64) int64) Jitter {
	if f == nil {
		f = rand.Int63n
	}
	return func(v JitterInterval) time.Duration {
		return randomJitter(v, f)
	}
}

func randomJitter(v JitterInterval, f func(n int64) int64) time.Duration {
	l := int64(v.L)
	if l == 0 {
		return 0
	}

	neg := l < 0
	if neg {
		l = -l
	}

	n := f(l) + 1
	if neg {
		n, l = -n, -l
	}

	if q := v.Q; q != 0 {
		n -= l / q
	}

	return time.Duration(n)
}

// JitterInterval describes a closed jitter interval. The final jitter value
// for random number N is computed as N - L/Q.
type JitterInterval struct {
	// L is the length of the jitter interval.
	L time.Duration

	// Q specifies 1/q quotient of the interval length to subtract from the
	// random number.
	Q int64
}

// Jitter implements the JitterProvider interface.
func (j Jitter) Jitter() Jitter {
	return j
}

type jitterSchedule struct {
	jitter   Jitter
	interval JitterInterval

	sched Schedule
}

// Schedule returns a new Schedule with jitter for the given interval.
func (j Jitter) Schedule(s Schedule, v JitterInterval) Schedule {
	return &jitterSchedule{
		jitter:   j,
		sched:    s,
		interval: v,
	}
}

// Next implements the Schedule interface.
func (j *jitterSchedule) Next(now time.Time) time.Time {
	next := j.sched.Next(now)
	if next.IsZero() || next.Before(now) {
		return next
	}

	d := j.jitter(j.interval)

	next = next.Add(d)
	if next.Before(now) {
		return now
	}
	return next
}

// Backoff returns a new Backoff with jitter for the given interval.
func (j Jitter) Backoff(b Backoff, v JitterInterval) Backoff {
	return func(n uint) (time.Duration, bool) {
		d, ok := b(n)
		if !ok {
			return d, ok
		}
		return d + j(v), ok
	}
}
