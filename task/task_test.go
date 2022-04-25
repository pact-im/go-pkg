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

func TestErrorGroup(t *testing.T) {
	oops := errors.New("oops")
	g := ErrorGroup(
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

func TestExitGroup(t *testing.T) {
	g := ExitGroup(
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
