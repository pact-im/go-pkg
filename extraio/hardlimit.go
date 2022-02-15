package extraio

import (
	"errors"
	"io"
)

// ErrExceededReadLimit is an error that HardLimitedReader returns when it
// exceeds read limit.
var ErrExceededReadLimit = errors.New("extraio: exceeded read limit")

// HardLimitedReader reads from R but limits the amount of data returned to just
// N bytes.
type HardLimitedReader struct {
	R io.Reader // underlying reader
	N uint64    // read limit

	readCount uint64
}

// HardLimitReader returns a Reader that reads from r but stops with an error
// after n bytes.
func HardLimitReader(r io.Reader, n uint64) *HardLimitedReader {
	return &HardLimitedReader{
		R: r,
		N: n,
	}
}

// Reset resets the reader’s state.
func (r *HardLimitedReader) Reset() {
	r.readCount = 0
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (r *HardLimitedReader) Read(p []byte) (int, error) {
	n, err := r.R.Read(p)
	if n <= 0 {
		return n, err
	}

	nn := uint64(n)
	if r.N-r.readCount < nn {
		return n, ErrExceededReadLimit
	}
	r.readCount += nn

	return n, err
}
