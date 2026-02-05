package httprange

import (
	"bytes"
	"math/rand/v2"
	"testing"
	"testing/iotest"
)

func TestBytesResourceReader_TestReader(t *testing.T) {
	r := rand.NewChaCha8([32]byte{})
	randomBytes := func(n int) []byte {
		buf := make([]byte, n)
		_, _ = r.Read(buf)
		return buf
	}

	tests := []struct {
		name    string
		content []byte
	}{
		{
			name:    "empty",
			content: []byte{},
		},
		{
			name:    "512",
			content: randomBytes(512),
		},
		{
			name:    "4k",
			content: randomBytes(4096),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := BytesResourceReader{
				Reader: bytes.NewReader(tt.content),
				Cancel: func(error) {},
				Length: int64(len(tt.content)),
			}
			if err := iotest.TestReader(&reader, tt.content); err != nil {
				t.Fatal(err)
			}
		})
	}
}
