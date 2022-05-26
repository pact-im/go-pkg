package extraio

import (
	"errors"
	"io"
)

// ErrExceededReadLimit is an error that HardLimitedReader returns when it
// exceeds read limit.
var ErrExceededReadLimit = errors.New("extraio: exceeded read limit")

// HardLimitedReader reads at most n bytes from the underlying reader and
// returns ErrExceededReadLimit if io.EOF is not reached once the limit is
// exceeded.
type HardLimitedReader struct {
	reader    io.Reader
	limit     uint64
	readCount uint64 // mutable
}

// HardLimitReader returns a Reader that reads from r but stops with an error
// after n bytes.
func HardLimitReader(r io.Reader, n uint64) *HardLimitedReader {
	return &HardLimitedReader{
		reader: r,
		limit:  n,
	}
}

// Reset resets the reader’s state.
func (r *HardLimitedReader) Reset() {
	r.readCount = 0
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (r *HardLimitedReader) Read(p []byte) (int, error) {
	if r.readCount == r.limit {
		return 0, ErrExceededReadLimit
	}

	// Do not read more than the remaining limit. This also guarantees that
	// n will not overflow or exceed limit on addition to readCount.
	limit := r.limit - r.readCount
	if uint64(len(p)) > limit {
		p = p[:limit]
	}

	n, err := r.reader.Read(p)
	if n > 0 {
		r.readCount += uint64(n)
	}
	if err == io.EOF {
		return n, err
	}
	return n, err
}
