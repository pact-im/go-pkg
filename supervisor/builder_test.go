package supervisor

import (
	"context"
	"errors"
	"testing"

	"go.pact.im/x/process"
)

func TestBuilderFuncRun(t *testing.T) {
	builder := Identity(process.Nop())

	called := false
	err := builder.Run(context.Background(), func(_ context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
	if !called {
		t.Error("callback was not called")
	}
}

func TestBuilderFuncError(t *testing.T) {
	sentinel := errors.New("sentinel")
	builder := BuilderFunc[process.Runner](
		func(context.Context) (process.Runner, error) {
			return nil, sentinel
		},
	)
	err := builder.Run(context.Background(), nil)
	if want := sentinel; err != want {
		t.Errorf("Run() error = %v, want %v", err, want)
	}
}

func TestBuilderRun(t *testing.T) {
	runner := process.Nop()
	builder := NewBuilder(Identity(runner))
	_ = builder.Run(context.Background(), func(_ context.Context) error {
		if r, ok := builder.Runner(); !ok || r != runner {
			t.Error("runner should be set in Run")
		}
		return nil
	})
	if _, ok := builder.Runner(); ok {
		t.Error("runner should not be set after Run")
	}
}

func TestBuilderError(t *testing.T) {
	sentinel := errors.New("sentinel")
	builder := NewBuilder(BuilderFunc[process.Runner](
		func(context.Context) (process.Runner, error) {
			return nil, sentinel
		},
	))
	err := builder.Run(context.Background(), nil)
	if want := sentinel; err != want {
		t.Errorf("Run() error = %v, want %v", err, want)
	}
}
