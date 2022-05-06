package process

import (
	"context"
	"sync"

	"go.pact.im/x/task"
)

// Parallel returns a Runnable instance that starts and runs processes in
// parallel. If no processes are given, it returns Nop instance.
//
// The resulting Runnable calls callback after all process dependencies are
// successfully started. If any dependecy fails to start, processes that have
// already started are gracefully stopped. If any dependency fails before the
// main callback returns, the context passed to callback is canceled and all
// processes are gracefully stopped (unless the parent context has expired).
//
// The callbacks of dependencies return after the callback of the resulting
// dependent process. Run returns callback error if it is not nil, otherwise it
// returns combined errors from dependencies.
func Parallel(deps ...Runnable) Runnable {
	switch len(deps) {
	case 0:
		return Nop()
	case 1:
		return deps[0]
	}
	return &groupRunnable{
		deps: deps,
		exec: task.ParallelExecutor(),
	}
}

// Sequential returns a Runnable instance that provides the same guarantees as
// returned by Parallel function, but starts and stops processes in sequential
// order.
func Sequential(deps ...Runnable) Runnable {
	switch len(deps) {
	case 0:
		return Nop()
	case 1:
		return deps[0]
	}
	return &groupRunnable{
		deps: deps,
		exec: task.SequentialExecutor(),
	}
}

type groupRunnable struct {
	deps []Runnable
	exec task.Executor
}

func (r *groupRunnable) Run(ctx context.Context, callback func(ctx context.Context) error) error {
	var once sync.Once
	var wg sync.WaitGroup
	wg.Add(1)

	// fgctx is passed to callback and cancel is used from child
	// process below to cancel callback invocation after startup.
	fgctx, cancel := context.WithCancel(ctx)
	defer cancel()

	child := func(ctx context.Context, callback func(ctx context.Context) error) error {
		err := callback(ctx)

		// Propagate process shutdown to main callback and wait
		// for it to return before exiting.
		cancel()
		wg.Wait()

		return err
	}

	n := len(r.deps)
	procs := make([]*Process, n)
	tasksArena := make([]task.Task, 2*n)
	startTasks := tasksArena[0*n : 1*n]
	stopTasks := tasksArena[1*n : 2*n]
	for i, dep := range r.deps {
		p := NewProcess(ctx, Chain(dep, RunnableFunc(child)))
		procs[i] = p
		startTasks[i] = func(ctx context.Context) error {
			err := p.Start(ctx)
			if err == nil {
				return nil
			}

			// If startup failed, we do not invoke callback.
			// Unblock callbacks for process dependencies
			// that have already started.
			once.Do(wg.Done)

			return err
		}
		stopTasks[i] = func(ctx context.Context) error {
			// We get either ErrProcessInvalidState or p.Err
			// from Stop so it is safe to ignore error here.
			_ = p.Stop(ctx)
			return p.Err()
		}
	}

	startError := r.exec.Execute(ctx, task.CancelOnError(), startTasks...)

	var callbackError error
	if startError == nil {
		callbackError = callback(fgctx)

		// Main callback has returned, unblock callbacks for
		// dependencies.
		wg.Done()
	}

	stopError := r.exec.Execute(ctx, task.NeverCancel(), stopTasks...)

	if callbackError != nil {
		return callbackError
	}

	return stopError
}
