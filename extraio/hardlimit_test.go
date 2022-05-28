package extraio

import (
	"bytes"
	"io"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHardLimitReader(t *testing.T) {
	testCases := []struct {
		Data  []byte
		Limit uint64
	}{
		{
			Data:  []byte{},
			Limit: 0,
		},
		{
			Data:  []byte{},
			Limit: 1,
		},
		{
			Data:  []byte{3, 3, 3},
			Limit: 0,
		},
		{
			Data:  []byte{1},
			Limit: 1,
		},
		{
			Data:  []byte{5, 5, 5, 5, 5},
			Limit: 3,
		},
		{
			Data:  []byte{3, 3, 3},
			Limit: 7,
		},
	}
	for _, tc := range testCases {
		data, err := io.ReadAll(HardLimitReader(bytes.NewReader(tc.Data), tc.Limit))
		switch {
		// Special case: read returns (n, io.EOF) and together with
		// previous reads we reach read limit but stop reading due to
		// the explicit EOF.
		case len(data) == len(tc.Data) && err == nil:
			fallthrough
		// Must succeed if data does not exceed the limit.
		case uint64(len(tc.Data)) < tc.Limit:
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.Data, data)
		default:
			assert.ErrorIs(t, err, ErrExceededReadLimit)
		}
	}
}
