package extraio

import (
	"hash"
	"io"
)

// HashReader wraps hash.Hash as an io.Reader. It assumes that the underlying
// hash state does not change after the first Read call.
type HashReader struct {
	h hash.Hash
	p []byte
}

// NewStrippedHashReader returns a new io.Reader that reads at most n byte of
// the given hash.
func NewStrippedHashReader(h hash.Hash, n int64) io.Reader {
	return io.LimitReader(NewHashReader(h), n)
}

// NewHashReader returns a new reader that reads from the given hash.
func NewHashReader(h hash.Hash) *HashReader {
	return &HashReader{h: h}
}

// Read implements the io.Reader interface. It computes the hash on the first
// call and advances through the hash buffer on subsequent calls to Read.
func (h *HashReader) Read(p []byte) (int, error) {
	if h.p == nil {
		h.p = h.h.Sum(nil)
	}

	n := copy(p, h.p)
	h.p = h.p[n:]
	if len(h.p) == 0 {
		return n, io.EOF
	}
	return n, nil
}
