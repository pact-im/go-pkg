package extraio

import (
	"crypto/sha256"
	"errors"
	"io"
	"testing"
	"testing/iotest"
)

func TestHashReader(t *testing.T) {
	data := []byte("test")
	hash := sha256ToSlice(sha256.Sum256(data))

	h := NewHashReader(sha256.New())

	_, err := h.Hash().Write(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = iotest.TestReader(h, hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h.Reset()

	err = iotest.TestReader(h, hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h.Hash().Reset()

	_, err = h.Read(nil)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected io.EOF, got %v", err)
	}

	h.Reset()

	hash = sha256ToSlice(sha256.Sum256(nil))
	err = iotest.TestReader(h, hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func sha256ToSlice(b [sha256.Size]byte) []byte {
	return b[:]
}
