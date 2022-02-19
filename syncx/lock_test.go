package syncx

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
)

func TestLock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	lock := NewLock()

	assert.Assert(t, lock.Acquire(ctx))

	cancel()
	assert.Assert(t, lock.Acquire(ctx) == ctx.Err())

	lock.Release()
}
