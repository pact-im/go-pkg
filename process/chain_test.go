package process

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
)

func TestChain(t *testing.T) {
	for i := 0; i < 5; i++ {
		testChain(t, i)
	}
}

func testChain(t *testing.T, count int) {
	expected := make([]int, count)
	for i := range expected {
		expected[i] = i
	}

	values := make([]int, 0, count)

	deps := make([]Runnable, count)
	for i := range deps {
		i := i
		deps[i] = RunnableFunc(func(ctx context.Context, callback func(ctx context.Context) error) error {
			values = append(values, i)
			return callback(ctx)
		})
	}

	seq := Chain(deps...)
	err := seq.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, values)
}
