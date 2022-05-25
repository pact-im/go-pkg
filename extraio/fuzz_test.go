//go:build go1.18 || gofuzzbeta
// +build go1.18 gofuzzbeta

package extraio

import (
	"bytes"
	"io"
	"testing"
)

func FuzzPadding(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte, blockSize uint8) {
		buf, err := io.ReadAll(
			NewUnpadReader(
				NewPadReader(
					bytes.NewReader(data),
					blockSize,
				),
				blockSize,
			),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, buf) {
			t.Fatal()
		}
	})
}

func FuzzUnpadReader(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte, n uint8) {
		runTestUnpadReader(t, data, n)
	})
}

func FuzzTailReader(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte, n uint) {
		runTestTailReader(t, data, n)
	})
}
