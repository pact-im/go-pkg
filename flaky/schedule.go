package flaky

import (
	"context"
	"errors"
	"time"

	"go.pact.im/x/clock"
)

var (
	// ErrNoNextSchedule is an error that is returned by ScheduleExecutor
	// if there is no next scheduled time or it is before current time.
	ErrNoNextSchedule = errors.New("flaky: no next scheduled time")

	// ErrScheduleDeadline is an error that is returned by ScheduleExecutor
	// if scheduled time exceeds the context deadline.
	ErrScheduleDeadline = errors.New("flaky: scheduled time exceeds context deadline")
)

// Schedule defines the schedule to use for WithSchedule.
type Schedule interface {
	// Next returns the next schedule time relative to now. If there is no
	// schedule past the given time, it returns zero time or a value before
	// now.
	Next(now time.Time) time.Time
}

type untilSchedule struct {
	sched Schedule
	until time.Time
}

// Until returns a schedule that stops at the given time. It allows defining the
// schedule for operations that should fail after the given time.
func Until(s Schedule, t time.Time) Schedule {
	return &untilSchedule{
		sched: s,
		until: t,
	}
}

// Next implements the Schedule interface.
func (t *untilSchedule) Next(now time.Time) time.Time {
	next := t.sched.Next(now)
	if t.until.Before(next) {
		return time.Time{}
	}
	return next
}

type midnightSchedule struct {
	d   time.Duration
	loc *time.Location
}

// Midnight returns a repeated Schedule for the duration d since midnight in the
// given in location. If the location is nil, it defaults to UTC.
//
// Passing a non-positive duration is equivalent to not using a schedule.
//
func Midnight(d time.Duration, loc *time.Location) Schedule {
	if loc == nil {
		loc = time.UTC
	}
	return &midnightSchedule{d, loc}
}

// Next implements the Schedule interface.
func (t *midnightSchedule) Next(now time.Time) time.Time {
	if t.d < 0 {
		return now
	}

	nowLoc := now.Location()
	now = now.In(t.loc)

	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	next := midnight.Add(t.d)
	if next.Before(now) {
		// Go does not use leap seconds so the day is always 24 hours.
		// See https://github.com/golang/go/issues/15247
		next = next.Add(24 * time.Hour)
	}

	return next.In(nowLoc)
}

type sleepSchedule struct {
	d time.Duration
}

// Sleep returns a Schedule that delays execution for a fixed interval d.
//
// Passing a non-positive duration is equivalent to not using a schedule.
//
func Sleep(d time.Duration) Schedule {
	return sleepSchedule{d}
}

// Next implements the Schedule interface.
func (t sleepSchedule) Next(now time.Time) time.Time {
	if t.d < 0 {
		return now
	}
	return now.Add(t.d)
}

// ScheduleExecutor is an executor that restricts operation execution to the
// specified schedule. It allows executing operations that may only succeed
// at the given time.
type ScheduleExecutor struct {
	clock *clock.Clock
	sched Schedule
	exec  Executor
}

// WithSchedule restricts the executor to wait for the given scheduled time
// before executing an operation.
func WithSchedule(e Executor, s Schedule) *ScheduleExecutor {
	return &ScheduleExecutor{
		clock: clock.System(),
		sched: s,
		exec:  e,
	}
}

// WithClock returns a copy of the executor that uses the given clock.
func (s *ScheduleExecutor) WithClock(c *clock.Clock) *ScheduleExecutor {
	if c == nil {
		c = clock.System()
	}
	return &ScheduleExecutor{
		clock: c,
		sched: s.sched,
		exec:  s.exec,
	}
}

// Execute implements the Executor interface.
func (s *ScheduleExecutor) Execute(ctx context.Context, f Op) error {
	now := s.clock.Now()
	next := s.sched.Next(now)

	if now.Equal(next) {
		return s.exec.Execute(ctx, f)
	}

	if next.Before(now) {
		return ErrNoNextSchedule
	}

	d := next.Sub(now)
	if !withinDeadline(ctx, d) {
		return ErrScheduleDeadline
	}

	timer := s.clock.Timer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C():
	}

	return s.exec.Execute(ctx, f)
}
