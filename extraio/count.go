package extraio

import (
	"io"
)

// CountReader counts bytes read from the underlying io.Reader.
type CountReader struct {
	r io.Reader // underlying reader
	n uint64
}

// NewCountReader returns a new reader that counts bytes read from r.
func NewCountReader(r io.Reader) *CountReader {
	return &CountReader{r: r}
}

// Count returns the count of bytes read.
func (r *CountReader) Count() uint64 {
	return r.n
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (r *CountReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if n > 0 {
		r.n += uint64(n)
	}
	return n, err
}
