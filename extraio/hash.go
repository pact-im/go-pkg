package extraio

import (
	"hash"
	"io"
)

// HashReader wraps hash.Hash as an io.Reader. It assumes that the underlying
// hash state does not change after the first Read call.
type HashReader struct {
	h hash.Hash

	buf []byte // mutable
	off int    // mutable
}

// NewStrippedHashReader returns a new io.Reader that reads at most n bytes of
// the given hash function.
func NewStrippedHashReader(h hash.Hash, n int64) io.Reader {
	return io.LimitReader(NewHashReader(h), n)
}

// NewHashReader returns a new reader that reads from the given hash function.
func NewHashReader(h hash.Hash) *HashReader {
	buf := make([]byte, 0, h.Size())
	return &HashReader{
		h:   h,
		buf: buf,
	}
}

// Hash returns the underlying hash function.
func (r *HashReader) Hash() hash.Hash {
	return r.h
}

// Reset resets the reader’s state. It does not reset the state of the
// underlying hash function.
func (r *HashReader) Reset() {
	r.buf = r.buf[:0]
	r.off = 0
}

// Read implements the io.Reader interface. It computes the hash on the first
// call and advances through the hash buffer on subsequent calls to Read.
func (r *HashReader) Read(p []byte) (int, error) {
	if len(r.buf) == 0 {
		r.buf = r.h.Sum(r.buf)
	}

	size := cap(r.buf)
	if size == r.off {
		return 0, io.EOF
	}

	n := copy(p, r.buf[r.off:])
	r.off += n

	if size == r.off {
		return n, io.EOF
	}

	return n, nil
}
