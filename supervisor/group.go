package supervisor

import (
	"context"
	"sync"

	"go.pact.im/x/process"
)

// Group manages a collection of [process.Runner] instances that run
// concurrently in the background and can be stopped together.
//
// Use Go method to add runners, Interrupt to signal them to shutdown gracefully,
// and Wait to wait for them to complete.
type Group struct {
	mu   sync.Mutex
	ctx  context.Context
	done chan struct{}
	wg   sync.WaitGroup
}

// NewGroup returns a new [Group] instance that runs processes in the background
// under the given context.
func NewGroup(ctx context.Context) *Group {
	return &Group{ctx: ctx}
}

// GroupBuilder is a [BuilderFunc] function for the [Group] type.
func GroupBuilder(ctx context.Context) (*Group, error) {
	return NewGroup(ctx), nil
}

// Run implements the [process.Runner] interface. It runs the callback, then
// interrupts all runners in the group and waits for them to complete.
func (g *Group) Run(ctx context.Context, callback process.Callback) error {
	defer g.Wait()
	defer g.Interrupt()
	return callback(ctx)
}

// Interrupt signals all runners started before the next [Group.Wait] call to
// stop.
func (g *Group) Interrupt() {
	g.mu.Lock()
	defer g.mu.Unlock()

	select {
	case <-g.done:
		return
	default:
	}

	if g.done == nil {
		g.done = make(chan struct{})
	}
	close(g.done)
}

// Wait blocks until all currently running processes exit.
func (g *Group) Wait() {
	g.wg.Wait()
	g.mu.Lock()
	g.done = nil
	g.mu.Unlock()
}

// Go starts a runner in the background, calling atExit when it finishes.
func (g *Group) Go(p process.Runner, atExit func(error)) {
	g.mu.Lock()
	if g.done == nil {
		g.done = make(chan struct{})
	}
	done := g.done
	g.mu.Unlock()

	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	if atExit == nil {
		atExit = func(_ error) {}
	}

	g.wg.Go(func() {
		atExit(p.Run(ctx, func(ctx context.Context) error {
			select {
			case <-ctx.Done():
			case <-done:
			}
			return nil
		}))
	})
}
