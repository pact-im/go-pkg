package extraio

import (
	"io"
)

// TailReader buffers the last n bytes read from the underlying io.Reader.
type TailReader struct {
	r io.Reader // underlying reader
	n uint64

	buf []byte
}

// NewTailReader returns a new reader that buffers last n bytes read from r.
func NewTailReader(r io.Reader, n uint64) *TailReader {
	return &TailReader{
		r: r,
		n: n,
	}
}

// Reset resets the reader’s state.
func (r *TailReader) Reset() {
	r.buf = nil
}

// Tail returns the last n bytes read.
func (r *TailReader) Tail() []byte {
	return r.buf
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (r *TailReader) Read(p []byte) (int, error) {
	// TODO(tie) this Read call must be optimized. It’s possible to perform
	// this operation with zero allocations (assuming that the internal buffer
	// is full). Haven’t profiled the program yet so not sure how it’d impact
	// the performance and GC pressure though.

	n, err := r.r.Read(p)
	if n <= 0 || r.n == 0 {
		return n, err
	}

	buf := r.buf
	buf = append(buf, p[:n]...)
	if uint64(len(buf)) <= r.n {
		r.buf = buf
		return 0, err
	}

	i := uint64(len(buf)) - r.n

	r.buf = make([]byte, r.n)
	copy(r.buf, buf[i:])

	n = copy(p, buf[:i])
	return n, err
}
