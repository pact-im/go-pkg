package process

import (
	"context"
	"errors"
	"testing"
)

func TestLeafCallbackReturns(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		<-ctx.Done()
		return oops
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}

func TestLeafCallbackError(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		return oops
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}

func TestLeafRunReturns(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		return nil
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return oops
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}

func TestLeafRunError(t *testing.T) {
	oops := errors.New("oops")
	leaf := Leaf(func(ctx context.Context) error {
		return oops
	}, nil)
	err := leaf.Run(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}

func TestStartStopLeafProcess(t *testing.T) {
	started, stopped := make(chan struct{}), make(chan struct{})
	p := NewProcess(
		context.Background(),
		Leaf(func(ctx context.Context) error {
			close(started)
			<-ctx.Done()
			<-stopped
			return nil
		}, func(ctx context.Context) error {
			close(stopped)
			return nil
		}),
	)
	proc := StartStop(p.Start, p.Stop)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		<-started
		return nil
	})
	if err != nil {
		t.FailNow()
	}
}

func TestStartStopErrorOnStart(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return oops },
		func(ctx context.Context) error { panic("unreachable") },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}

func TestStartStopErrorOnCallback(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return nil },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return oops
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}

func TestStartStopErrorOnCallbackAndStop(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return errors.New("ignored") },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return oops
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}

func TestStartStopErrorOnStop(t *testing.T) {
	oops := errors.New("oops")
	proc := StartStop(
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return oops },
	)
	err := proc.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(err, oops) {
		t.FailNow()
	}
}
