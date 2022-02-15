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
		if uint64(len(tc.Data)) <= tc.Limit {
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.Data, data)
		} else {
			assert.ErrorIs(t, err, ErrExceededReadLimit)
		}
	}
}
