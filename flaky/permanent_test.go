package flaky

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestUntilPermanent(t *testing.T) {
	ctx := context.Background()
	oops := errors.New("oops")
	exec := UntilPermanent()
	t.Run("Nil", func(t *testing.T) {
		assert.NilError(t, exec.Execute(ctx, func(_ context.Context) error {
			return nil
		}))
	})
	t.Run("Internal", func(t *testing.T) {
		err := exec.Execute(ctx, func(_ context.Context) error {
			return Internal(oops)
		})
		assert.ErrorIs(t, oops, err)
	})
	t.Run("Permanent", func(t *testing.T) {
		err := exec.Execute(ctx, func(_ context.Context) error {
			return Permanent(oops)
		})
		assert.Assert(t, IsPermanentError(err))
		assert.ErrorIs(t, err, oops)
	})
}
