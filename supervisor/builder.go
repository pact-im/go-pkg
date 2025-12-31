package supervisor

import (
	"context"
	"sync"

	"go.pact.im/x/option"
	"go.pact.im/x/process"
)

// BuilderFunc is a factory function that constructs [process.Runner] instances.
type BuilderFunc[T process.Runner] func(context.Context) (T, error)

// Identity returns a [BuilderFunc] that always returns the given runner.
func Identity[T process.Runner](runner T) BuilderFunc[T] {
	return func(_ context.Context) (T, error) {
		return runner, nil
	}
}

// Run implements the [process.Runner] interface. It calls the factory function
// to create a runner, then delegates to that runner’s Run method.
func (f BuilderFunc[T]) Run(ctx context.Context, callback process.Callback) error {
	runner, err := f(ctx)
	if err != nil {
		return err
	}
	return runner.Run(ctx, callback)
}

// Builder is a reusable [process.Runner] implementation that uses [BuilderFunc]
// on each [Builder.Run] invocation to create a fresh [process.Runner] instance
// for implementations that would otherwise be single-use (do not allow calling
// Run more than once per instance).
type Builder[T process.Runner] struct {
	build BuilderFunc[T]

	runnerMu sync.RWMutex
	runner   option.Of[T]
}

// NewBuilder returns a new [Builder] instance for the given factory function.
// The function will be called each time [Builder.Run] method is invoked to
// create a fresh runner instance.
func NewBuilder[T process.Runner](f BuilderFunc[T]) *Builder[T] {
	return &Builder[T]{
		build: f,
	}
}

// Runner returns the runner currently executing via [Builder.Run] It returns
// false when no execution is active.
func (b *Builder[T]) Runner() (T, bool) {
	b.runnerMu.RLock()
	defer b.runnerMu.RUnlock()
	return b.runner.Unwrap()
}

// setRunner updates the currently execution runner instance.
func (b *Builder[T]) setRunner(v option.Of[T]) {
	b.runnerMu.Lock()
	defer b.runnerMu.Unlock()
	b.runner = v
}

// Run implements the [process.Runner] interface. It uses [BuildFunc] to create
// a runner, executes it, and discards the instance when complete.
//
// Use [Builder.Runner] to get the currently executing runner instance.
func (b *Builder[T]) Run(ctx context.Context, callback process.Callback) error {
	runner, err := b.build(ctx)
	if err != nil {
		return err
	}

	b.setRunner(option.Value(runner))
	defer b.setRunner(option.Nil[T]())

	return runner.Run(ctx, callback)
}
