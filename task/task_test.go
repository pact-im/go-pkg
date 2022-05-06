package task

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNamed(t *testing.T) {
	oops := errors.New("oops")
	task := Named("test", func(ctx context.Context) error {
		return oops
	})
	err := task.Run(context.Background())
	assert.Equal(t, err.Error(), "test: oops")
}

func TestParallelCancelOnError(t *testing.T) {
	oops := errors.New("oops")
	g := Parallel(CancelOnError(),
		func(ctx context.Context) error {
			return oops
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		},
	)
	err := g.Run(context.Background())
	assert.ErrorIs(t, err, oops)
}

func TestParallelCancelOnReturn(t *testing.T) {
	g := Parallel(CancelOnReturn(),
		func(ctx context.Context) error {
			return nil
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		},
	)
	err := g.Run(context.Background())
	assert.NilError(t, err)
}
