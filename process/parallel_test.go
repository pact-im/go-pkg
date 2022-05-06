package process

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
)

func TestParallel(t *testing.T) {
	for i := 0; i < 5; i++ {
		testParallel(t, i)
	}
}

func testParallel(t *testing.T, count int) {
	expected := make([]int, count)
	for i := range expected {
		expected[i] = i
	}

	values := make([]int, count)

	deps := make([]Runnable, count)
	for i := range deps {
		i := i
		deps[i] = RunnableFunc(func(ctx context.Context, callback func(ctx context.Context) error) error {
			values[i] = i
			return callback(ctx)
		})
	}

	par := Parallel(deps...)
	err := par.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, values)
}
