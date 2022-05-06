package task

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestContextifyCancel(t *testing.T) {
	oops := errors.New("oops")

	stop := make(chan struct{})
	c := Contextify(func() error {
		<-stop
		return oops
	}, func() {
		close(stop)
	})

	canceledContext, cancel := context.WithCancel(context.Background())
	cancel()
	err := c.Run(canceledContext)
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
