// Package process provides process management and supervision implementation.
package process

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/atomic"
)

// ProcessState represents the current state of the process.
type ProcessState int

const (
	ProcessStateInitial ProcessState = iota
	ProcessStateStarting
	ProcessStateRunning
	ProcessStateStopped
)

// Process represents a stateful process that is running in the background.
// It exposes Start and Stop methods that use an underlying state machine to
// prevent operations in invalid states. Process is safe for concurrent use.
type Process[P Runnable] struct {
	proc   P
	parent context.Context

	stateMu sync.Mutex
	state   ProcessState

	cancel context.CancelFunc

	stop chan struct{}
	done chan struct{}
	err  atomic.Error

	// stopped is used by Manager to remove the process instance from
	// internal map at most once.
	stopped atomic.Bool
}

// NewProcess returns a new stateful process instance for the given Runnable
// type parameter that would run with the ctx context.
func NewProcess[P Runnable](ctx context.Context, proc P) *Process[P] {
	return &Process[P]{
		proc:   proc,
		parent: ctx,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
}

// Done returns a channel that is closed when process terminates.
func (p *Process[P]) Done() <-chan struct{} {
	return p.done
}

// Err returns the error from running the process.
func (p *Process[P]) Err() error {
	return p.err.Load()
}

// State returns the current process state.
func (p *Process[P]) State() ProcessState {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()
	return p.state
}

// Start starts the process or cancels the underlying process context on error.
// The startup deadline may be set using the given ctx context. It would not
// be used by process.
//
// It returns ErrProcessInvalidState if the process is not in the initial state.
func (p *Process[P]) Start(ctx context.Context) error {
	var bgctx context.Context
	var cancel context.CancelFunc
	if !p.transition(ProcessStateStarting, func() {
		bgctx, p.cancel = context.WithCancel(p.parent)
		p.parent, cancel = nil, p.cancel
	}) {
		return ErrProcessInvalidState
	}

	init := make(chan struct{})
	go func() {
		err := p.proc.Run(bgctx, func(bgctx context.Context) error {
			_ = p.transition(ProcessStateRunning, nil)

			close(init)

			select {
			case <-bgctx.Done():
			case <-p.stop:
			}
			return nil
		})
		_ = p.transition(ProcessStateStopped, func() {
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

// Stop stops the process by canceling the underlying context and waiting for
// the termination.
//
// It returns ErrProcessInvalidState if the process was not started or has
// already been stopped.
func (p *Process[P]) Stop(ctx context.Context) error {
	var initial bool
	var cancel context.CancelFunc
	if !p.transition(ProcessStateStopped, func() {
		switch p.state {
		case ProcessStateInitial:
			initial = true
		case ProcessStateStarting:
			p.cancel()
			cancel = p.cancel
			p.cancel = nil
		case ProcessStateRunning:
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
func (p *Process[P]) transition(next ProcessState, advance func()) bool {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()

	switch p.state {
	case ProcessStateInitial:
		if next != ProcessStateStarting && next != ProcessStateStopped {
			return false
		}
	case ProcessStateStarting:
		if next != ProcessStateRunning && next != ProcessStateStopped {
			return false
		}
	case ProcessStateRunning:
		if next != ProcessStateStopped {
			return false
		}
	case ProcessStateStopped:
		return false
	}

	if advance != nil {
		advance()
	}

	p.state = next
	return true
}
