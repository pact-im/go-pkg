package mockclock

import (
	"go.pact.im/x/clock"
)

// Clock is a union of all interfaces implemented by a full-featured clock
// implementation. It is satisfied by clock.Clock that provides shims for
// missing functionality.
//
// It is an input for mock types generator and must not be used outside of this
// package.
type Clock interface {
	clock.Scheduler
	clock.NowScheduler
	clock.TimerScheduler
	clock.TickerScheduler

	private()
}
