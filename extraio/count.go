package extraio

import (
	"errors"
	"io"
	"math"
)

// ErrCounterOverflow is an error that CountReader returns when the counter
// overflows 64-bit unsigned integer.
var ErrCounterOverflow = errors.New("extraio: counter overflow")

// CountReader counts bytes read from the underlying io.Reader. It is implicitly
// limited to at most math.MaxUint64 bytes and returns ErrCounterOverflow error
// if the limit is exceeded. In practice, that would require counting thousands
// of petabytes to reach this limitation.
type CountReader struct {
	reader io.Reader
	count  uint64 // mutable
}

// NewCountReader returns a new reader that counts bytes read from r.
func NewCountReader(r io.Reader) *CountReader {
	return &CountReader{reader: r}
}

// Count returns the count of bytes read.
func (r *CountReader) Count() uint64 {
	return r.count
}

// Read implements the io.Reader interface. It reads from the underlying
// io.Reader and increments read bytes counter.
func (r *CountReader) Read(p []byte) (int, error) {
	switch k := math.MaxUint64 - r.count; {
	case k == 0:
		return 0, ErrCounterOverflow
	case k < uint64(len(p)):
		p = p[:k]
	}
	n, err := r.reader.Read(p)
	if n > 0 {
		r.count += uint64(n)
	}
	return n, err
}
