package extraio

import (
	"io"
	"math"
)

// CountReader counts bytes read from the underlying io.Reader.
type CountReader struct {
	reader   io.Reader
	count    uint64 // mutable
	overflow bool   // mutable
}

// NewCountReader returns a new reader that counts bytes read from r.
func NewCountReader(r io.Reader) *CountReader {
	return &CountReader{reader: r}
}

// Count returns the count of bytes read. It returns false if the count of read
// bytes cannot be represented as 64 bit unsigned integer. In practice, that
// would require counting thousands of petabytes to reach this limitation.
func (r *CountReader) Count() (uint64, bool) {
	return r.count, !r.overflow
}

// Read implements the io.Reader interface. It reads from the underlying
// io.Reader and increments read bytes counter.
func (r *CountReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 && !r.overflow {
		nn := uint64(n)
		r.count += nn
		if r.count < nn {
			r.overflow = true
			r.count = math.MaxUint64
		}
	}
	return n, err
}
