package flaky

import (
	"math"
	"time"
)

// BackoffProvider provides a potentially stateful backoff function for
// executing an operation.
type BackoffProvider interface {
	Backoff() Backoff
}

// Backoff is a function that, given the current count of retries, returns
// backoff delay before the next attempt, or false if an executor should stop
// with the last error.
type Backoff func(n uint) (time.Duration, bool)

// Backoff implements the BackoffProvider interface.
func (b Backoff) Backoff() Backoff {
	return b
}

// Limit limits the number of retries for the given backoff. That is, it stops
// the backoff once retries counter reaches the given limit. Passing a limit of
// zero retries is equivalent to no backoff.
func Limit(limit uint, backoff Backoff) Backoff {
	return func(n uint) (time.Duration, bool) {
		if n >= limit {
			return 0, false
		}
		return backoff(n)
	}
}

// Constant returns a constant backoff that uses nth duration for the arguments.
// It returns the last element if the attempt count is greater than the number
// of arguments. Passing no arguments is equivalent to no backoff.
func Constant(xs ...time.Duration) Backoff {
	return func(n uint) (time.Duration, bool) {
		k := uint(len(xs))
		if 0 == k {
			return 0, false
		}
		if n >= k {
			return xs[k-1], true
		}
		return xs[n], true
	}
}

// Exp2 returns an exponential backoff that returns base-2 exponential of
// failed attempts count times the given duration unit.
//
// Note that the backoff stops after wait duration reaches max positive value of
// time.Duration or once n is equal to math.MaxUint.
func Exp2(unit time.Duration) Backoff {
	return func(n uint) (time.Duration, bool) {
		if n == math.MaxUint {
			return 0, false
		}
		x := unit << (n + 1)
		if x <= 0 {
			return 0, false
		}
		return x, true
	}
}
