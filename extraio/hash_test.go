package extraio

import (
	"crypto/sha256"
	"io"
	"testing"
	"testing/iotest"

	"gotest.tools/v3/assert"
)

func TestHashReader(t *testing.T) {
	data := []byte("test")
	hash := sha256ToSlice(sha256.Sum256(data))

	h := NewHashReader(sha256.New())

	_, err := h.Hash().Write(data)
	assert.NilError(t, err)

	err = iotest.TestReader(h, hash)
	assert.NilError(t, err)

	h.Reset()

	err = iotest.TestReader(h, hash)
	assert.NilError(t, err)

	h.Hash().Reset()

	_, err = h.Read(nil)
	assert.ErrorIs(t, err, io.EOF)

	h.Reset()

	hash = sha256ToSlice(sha256.Sum256(nil))
	err = iotest.TestReader(h, hash)
	assert.NilError(t, err)
}

func sha256ToSlice(b [sha256.Size]byte) []byte {
	return b[:]
}
