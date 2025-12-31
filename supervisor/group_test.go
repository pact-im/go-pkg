package supervisor

import (
	"context"
	"testing"
	"testing/synctest"

	"go.pact.im/x/process"
)

func TestGroupRunEmpty(t *testing.T) {
	synctest.Test(t, synctestTestGroupRunEmpty)
}

func synctestTestGroupRunEmpty(t *testing.T) {
	var g Group

	called := false
	err := g.Run(context.Background(), func(_ context.Context) error {
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

func TestGroupGoInRun(t *testing.T) {
	synctest.Test(t, synctestTestGroupGoInRun)
}

func synctestTestGroupGoInRun(t *testing.T) {
	var g Group

	err := g.Run(context.Background(), func(_ context.Context) error {
		g.Go(process.Nop(), nil)
		return nil
	})
	if err != nil {
		t.Errorf("Run() error = %v, want nil", err)
	}
}

func TestGroupGoAndInterrupt(t *testing.T) {
	synctest.Test(t, synctestTestGroupGoAndInterrupt)
}

func synctestTestGroupGoAndInterrupt(t *testing.T) {
	ctx := context.Background()
	g := NewGroup(ctx)

	var started, stopped bool
	runner := process.RunnerFunc(func(ctx context.Context, callback process.Callback) error {
		started = true
		return callback(ctx)
	})

	g.Go(runner, func(_ error) {
		stopped = true
	})

	synctest.Wait()

	if !started {
		t.Error("runner not started")
	}
	if stopped {
		t.Error("runner stopped too early")
	}

	g.Interrupt()
	g.Interrupt() // also test interrupting twice
	g.Wait()

	if !stopped {
		t.Error("runner not stopped after interrupt")
	}
}
