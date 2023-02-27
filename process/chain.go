package process

import (
	"context"
)

// Chain returns a [Runner] instance that starts and runs processes by nesting
// them in callbacks. If no processes are given, it returns [Nop] instance.
func Chain(deps ...Runner) Runner {
	switch len(deps) {
	case 0:
		return Nop()
	case 1:
		return deps[0]
	}
	return &chainRunner{deps}
}

type chainRunner struct {
	deps []Runner
}

func (r *chainRunner) Run(ctx context.Context, callback Callback) error {
	s := chainState{deps: r.deps}
	return s.Run(ctx, callback)
}

type chainState struct {
	index int
	deps  []Runner // len(deps) >= 2
	main  Callback
}

func (r *chainState) Run(ctx context.Context, callback Callback) error {
	switch r.index {
	case 0:
		r.main = callback
		callback = r.next
	case len(r.deps) - 1:
		callback = r.main
	default:
		callback = r.next
	}
	i := r.index
	r.index++
	return r.deps[i].Run(ctx, callback)
}

func (r *chainState) next(ctx context.Context) error {
	return r.Run(ctx, nil)
}
