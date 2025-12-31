package supervisor

import (
	"context"
	"errors"
	"slices"
	"testing"
	"testing/synctest"

	"go.pact.im/x/flaky"
	"go.pact.im/x/process"
)

func TestSupervisorInterruptsRunnerAfterCallback(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorInterruptsRunnerAfterCallback)
}

func synctestTestSupervisorInterruptsRunnerAfterCallback(t *testing.T) {
	ready := make(chan struct{})
	s := NewSupervisor(process.Nop(), flaky.UntilPermanent(), Hook{
		Pre: func(_ context.Context, _ *Supervisor) error {
			close(ready)
			return nil
		},
	})
	err := s.Run(context.Background(), func(_ context.Context) error {
		<-ready
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
}

func TestSupervisorInterruptRunnerFromPreHook(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorInterruptRunnerFromPreHook)
}

func synctestTestSupervisorInterruptRunnerFromPreHook(t *testing.T) {
	s := NewSupervisor(process.Nop(), flaky.UntilPermanent(), Hook{
		Pre: func(_ context.Context, s *Supervisor) error {
			s.Interrupt()
			return nil
		},
	})
	err := s.Run(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
}

func TestSupervisorInterruptRunnerFromExecutor(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorInterruptRunnerFromExecutor)
}

func synctestTestSupervisorInterruptRunnerFromExecutor(t *testing.T) {
	var interrupt func()
	exec := func(ctx context.Context, f flaky.Op) error {
		if ctx.Err() != nil {
			return nil
		}
		interrupt()
		return f(ctx)
	}
	s := NewSupervisor(process.Nop(), flaky.ExecutorFunc(exec), Hook{})
	interrupt = s.Interrupt
	err := s.Run(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
}

func TestSupervisorPreHookError(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorPreHookError)
}

func synctestTestSupervisorPreHookError(t *testing.T) {
	sentinel := errors.New("sentinel")
	s := NewSupervisor(process.Nop(), flaky.Once(), Hook{
		Pre: func(_ context.Context, _ *Supervisor) error {
			return sentinel
		},
		Post: func(_ context.Context, _ *Supervisor) error {
			return errors.New("ignored")
		},
	})
	err := s.Run(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	if err != sentinel {
		t.Errorf("Run() error = %v, want %v", err, sentinel)
	}
}

func TestSupervisorCallbackError(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorCallbackError)
}

func synctestTestSupervisorCallbackError(t *testing.T) {
	sentinel := errors.New("sentinel")
	ready := make(chan struct{})
	runner := func(ctx context.Context, callback process.Callback) error {
		close(ready)
		return callback(ctx)
	}
	s := NewSupervisor(process.RunnerFunc(runner), flaky.UntilPermanent(), Hook{})
	err := s.Run(context.Background(), func(_ context.Context) error {
		<-ready
		return sentinel
	})
	if err != sentinel {
		t.Errorf("Run() error = %v, want %v", err, sentinel)
	}
}

func TestSupervisorBothError(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorBothError)
}

func synctestTestSupervisorBothError(t *testing.T) {
	sentinel1 := errors.New("sentinel1")
	sentinel2 := errors.New("sentinel2")
	ready := make(chan struct{})
	s := NewSupervisor(process.Nop(), flaky.UntilPermanent(), Hook{
		Pre: func(_ context.Context, _ *Supervisor) error {
			close(ready)
			return nil
		},
		Post: func(_ context.Context, _ *Supervisor) error {
			return sentinel2
		},
	})
	err := s.Run(context.Background(), func(_ context.Context) error {
		<-ready
		return sentinel1
	})
	errs := err.(interface{ Unwrap() []error }).Unwrap()
	if !slices.Equal(errs, []error{sentinel1, sentinel2}) {
		t.Errorf("Run() error = %v", err)
	}
}

func TestSupervisorInterruptsBeforeRun(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorInterruptsBeforeRun)
}

func synctestTestSupervisorInterruptsBeforeRun(t *testing.T) {
	ready := make(chan struct{})
	s := NewSupervisor(process.Nop(), flaky.UntilPermanent(), Hook{
		Pre: func(_ context.Context, _ *Supervisor) error {
			close(ready)
			return nil
		},
	})
	s.Interrupt() // no-op
	err := s.Run(context.Background(), func(_ context.Context) error {
		<-ready
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
}

func TestSupervisorRecursiveRunError(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorRecursiveRunError)
}

func synctestTestSupervisorRecursiveRunError(t *testing.T) {
	var innerError error
	exec := func(ctx context.Context, _ flaky.Op) error {
		<-ctx.Done()
		return nil
	}
	s := NewSupervisor(process.Nop(), flaky.ExecutorFunc(exec), Hook{})
	err := s.Run(context.Background(), func(ctx context.Context) error {
		innerError = s.Run(ctx, func(_ context.Context) error {
			panic("unreachable")
		})
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
	if want := errRecursiveOrConcurrentRun; innerError != want {
		t.Errorf("nested Run() error = %v, want %v", err, want)
	}
}

func TestSupervisorExecutorInterruptPanic(t *testing.T) {
	synctest.Test(t, synctestTestSupervisorExecutorInterruptPanic)
}

func synctestTestSupervisorExecutorInterruptPanic(t *testing.T) {
	var panicked bool
	runner := func(_ context.Context, _ process.Callback) error {
		return nil
	}
	exec := func(ctx context.Context, f flaky.Op) error {
		var err error
		for !panicked {
			func() {
				defer func() { panicked = recover() != nil }()
				err = f(ctx)
			}()
		}
		return err
	}
	s := NewSupervisor(process.RunnerFunc(runner), flaky.ExecutorFunc(exec), Hook{})
	err := s.Run(context.Background(), func(ctx context.Context) error {
		s.Interrupt()
		<-ctx.Done()
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
	if !panicked {
		t.Error("should panic when executor does not respect interrupt")
	}
}
