// Package process provides process management and supervision implementation.
package process

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/atomic"
)

// State represents the current state of the process.
type State int

const (
	// StateInitial is the initial state of the process. If Stop is called
	// in initial state, it prevents subsequent Start calls from succeeding,
	// thus entering StateStopped. Otherwise the transition on Start call
	// is to the StateStarting.
	StateInitial State = iota
	// StateStarting is the starting state that process enters when Start
	// is called from initial state. It transitions to either StateRunning
	// on success or StateStopped on failure (or premature shutdown observed
	// during startup).
	StateStarting
	// StateRunning is the state process enters after a successful startup.
	// The only possible transition is to the StateStopped if either Stop
	// is called or a process terminates.
	StateRunning
	// StateStopped is the final state of the process. There are no
	// transitions from this state.
	StateStopped
)

// Process represents a stateful process that is running in the background.
// It exposes Start and Stop methods that use an underlying state machine to
// prevent operations in invalid states. Process is safe for concurrent use.
//
// Unlike some implementations of the underlying Runnable interface that allow
// multiple consecutive Run invocations on the same instance, a Process may not
// be reset and started after being stopped.
type Process struct {
	proc   Runnable
	parent context.Context

	stateMu sync.Mutex
	state   State

	cancel context.CancelFunc

	stop chan struct{}
	done chan struct{}
	err  atomic.Error
}

// NewProcess returns a new stateful process instance for the given Runnable
// type parameter that would run with the ctx context.
func NewProcess(ctx context.Context, proc Runnable) *Process {
	return &Process{
		proc:   proc,
		parent: ctx,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
}

// Done returns a channel that is closed when process terminates.
func (p *Process) Done() <-chan struct{} {
	return p.done
}

// Err returns the error from running the process.
func (p *Process) Err() error {
	return p.err.Load()
}

// State returns the current process state.
func (p *Process) State() State {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()
	return p.state
}

// Start starts the process or cancels the underlying process context on error.
//
// The startup deadline may be set using the given ctx context. Note that the
// context would not be used by process directly so the associated values are
// not propagated.
//
// It returns ErrProcessInvalidState if the process is not in the initial state.
func (p *Process) Start(ctx context.Context) error {
	var bgctx context.Context
	var cancel context.CancelFunc
	if !p.transition(StateStarting, func() {
		bgctx, cancel = context.WithCancel(p.parent)
		p.parent, p.cancel = nil, cancel
	}) {
		return ErrProcessInvalidState
	}

	init := make(chan struct{})
	go func() {
		err := p.proc.Run(bgctx, func(bgctx context.Context) error {
			_ = p.transition(StateRunning, nil)

			close(init)

			select {
			case <-bgctx.Done():
			case <-p.stop:
			}
			return nil
		})
		_ = p.transition(StateStopped, func() {
			p.cancel = nil
		})
		cancel()
		p.err.Store(err)
		close(p.done)
	}()

	select {
	case <-ctx.Done():
		// Propagate context cancellation to cancel startup.
		cancel()
		<-p.Done()
		return ctx.Err()
	case <-p.Done():
		// Note that a non-nil error is still an error on start if the
		// process terminates without initialization.
		return fmt.Errorf("run: %w", p.Err())
	case <-init:
		// OK
	}
	return nil
}

// Stop stops the process by returning from the Run method callback and waiting
// for the termination.
//
// The shutdown deadline may be set using the given ctx context. If the deadline
// is exceeded, underlying context is canceled, signaling a forced shutdown to
// the process.
//
// It returns ErrProcessInvalidState if the process was not started or has
// already been stopped.
func (p *Process) Stop(ctx context.Context) error {
	var initial bool
	var cancel context.CancelFunc
	if !p.transition(StateStopped, func() {
		switch p.state {
		case StateInitial:
			close(p.done)
			initial = true
			p.parent = nil
		case StateStarting:
			p.cancel()
			cancel = p.cancel
			p.cancel = nil
		case StateRunning:
			close(p.stop)
			cancel = p.cancel
			p.cancel = nil
		}
	}) {
		return ErrProcessInvalidState
	}
	if initial {
		// We prevented a start. Do nothing.
		return nil
	}
	select {
	case <-ctx.Done():
		// Propagate context cancellation to force shutdown.
		cancel()
		<-p.Done()
		return ctx.Err()
	case <-p.Done():
		if err := p.Err(); err != nil {
			return fmt.Errorf("run: %w", err)
		}
		return nil
	}
}

// transition advances to the next process state. It returns false if there is
// not transition from the current to the given next state.
func (p *Process) transition(next State, advance func()) bool {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()

	switch p.state {
	case StateInitial:
		if next != StateStarting && next != StateStopped {
			return false
		}
	case StateStarting:
		if next != StateRunning && next != StateStopped {
			return false
		}
	case StateRunning:
		if next != StateStopped {
			return false
		}
	case StateStopped:
		return false
	}

	if advance != nil {
		advance()
	}

	p.state = next
	return true
}
