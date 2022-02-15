package monikernames

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestNamer(t *testing.T) {
	name, err := New().Build().Name(context.Background())
	assert.NilError(t, err)
	assert.Assert(t, is.Regexp(`^\w+ \w+$`, name))
}
