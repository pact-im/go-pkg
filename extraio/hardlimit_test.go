package extraio

import (
	"bytes"
	"errors"
	"io"
	"testing"
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
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(tc.Data, data) {
				t.Fatalf("expected %v, got %v", tc.Data, data)
			}
		default:
			if !errors.Is(err, ErrExceededReadLimit) {
				t.Fatalf("expected ErrExceededReadLimit, got %v", err)
			}
		}
	}
}
