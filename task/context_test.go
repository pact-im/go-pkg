package task

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestContextifyCancel(t *testing.T) {
	oops := errors.New("oops")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stop := make(chan struct{})
	c := Contextify(func() error {
		<-stop
		return oops
	}, func() {
		close(stop)
	})

	err := c.Run(ctx)
	assert.ErrorIs(t, err, oops)
}

func TestContextifyReturn(t *testing.T) {
	c := Contextify(func() error {
		return nil
	}, func() {
		panic("unreachable")
	})
	err := c.Run(context.Background())
	assert.NilError(t, err)
}
