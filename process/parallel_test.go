package process

import (
	"context"
	"testing"
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

	deps := make([]Runner, count)
	for i := range deps {
		i := i
		deps[i] = RunnerFunc(func(ctx context.Context, callback Callback) error {
			values[i] = i
			return callback(ctx)
		})
	}

	par := Parallel(deps...)
	err := par.Run(context.Background(), func(_ context.Context) error {
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
