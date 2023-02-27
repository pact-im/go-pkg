package process

import (
	"context"
	"testing"
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

	deps := make([]Runner, count)
	for i := range deps {
		i := i
		deps[i] = RunnerFunc(func(ctx context.Context, callback Callback) error {
			values = append(values, i)
			return callback(ctx)
		})
	}

	seq := Chain(deps...)
	err := seq.Run(context.Background(), func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.FailNow()
	}
	if len(expected) != len(values) {
		t.FailNow()
	}
	for i := range expected {
		if expected[i] != values[i] {
			t.FailNow()
		}
	}
}
