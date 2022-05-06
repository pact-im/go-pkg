package process

import (
	"context"
)

// Chain returns a Runnable instance that starts and runs processes by nesting
// them in callbacks. If no processes are given, it returns Nop instance.
func Chain(deps ...Runnable) Runnable {
	switch len(deps) {
	case 0:
		return Nop()
	case 1:
		return deps[0]
	}
	return &chainRunnable{deps}
}

type chainRunnable struct {
	deps []Runnable
}

func (r *chainRunnable) Run(ctx context.Context, callback func(ctx context.Context) error) error {
	s := chainState{deps: r.deps}
	return s.Run(ctx, callback)
}

type chainState struct {
	index int
	deps  []Runnable // len(deps) >= 2
	main  func(ctx context.Context) error
}

func (r *chainState) Run(ctx context.Context, callback func(ctx context.Context) error) error {
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
