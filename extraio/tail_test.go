package extraio

import (
	"bytes"
	"io"
	"testing"

	"gotest.tools/v3/assert"
)

func TestTailReader(t *testing.T) {
	testCases := []struct {
		Data []byte
		Size uint
	}{
		{
			Data: nil,
			Size: 0,
		},
		{
			Data: []byte{1},
			Size: 0,
		},
		{
			Data: []byte{},
			Size: 1,
		},
		{
			Data: []byte{1},
			Size: 1,
		},
		{
			Data: []byte{1, 2},
			Size: 1,
		},
		{
			Data: []byte{1, 2, 3, 4, 5},
			Size: 3,
		},
		{
			Data: []byte{1, 2, 3, 4, 5},
			Size: 7,
		},
	}
	for _, tc := range testCases {
		runTestTailReader(t, tc.Data, tc.Size)
	}
}

func runTestTailReader(t testing.TB, data []byte, n uint) {
	var head, tail []byte
	switch {
	case n == 0:
		head, tail = data, nil
	case n >= uint(len(data)):
		head, tail = nil, data
	default:
		i := uint(len(data)) - n
		head, tail = data[:i], data[i:]
	}

	tr := NewTailReader(bytes.NewReader(data), n)

	buf, err := io.ReadAll(tr)
	assert.NilError(t, err)

	assert.Check(t, bytes.Equal(head, buf))
	assert.Check(t, bytes.Equal(tail, tr.Tail()))
}
