package task

import (
	"context"
)

var (
	sequentialExecutor = executor{Sequential}
	parallelExecutor   = executor{Parallel}
)

// Executor abstracts the execution order for task groups.
type Executor interface {
	Execute(ctx context.Context, cond CancelCondition, tasks ...Task) error
}

// executor implements the Executor interface using the underlying task group
// builder.
type executor struct {
	group func(cond CancelCondition, tasks ...Task) Task
}

// SequentialExecutor returns an Executor instance for sequential execution.
func SequentialExecutor() Executor {
	return &sequentialExecutor
}

// ParallelExecutor returns an Executor instance for parallel execution.
func ParallelExecutor() Executor {
	return &parallelExecutor
}

// Execute implements the Executor interface.
func (e *executor) Execute(ctx context.Context, cond CancelCondition, tasks ...Task) error {
	return e.group(cond, tasks...).Run(ctx)
}
