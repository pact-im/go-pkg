package extraio

import (
	"errors"
	"io"
)

// PadReader is an io.Reader that always adds non-zero padding on EOF. The
// the value of the padding bytes is equal to the length of the padding.
type PadReader struct {
	reader    io.Reader
	blockSize uint8

	readCount uint64
	padding   uint8
	fillByte  byte
}

// NewPadReader returns a new reader that pads r with the given block size.
// If blockSize is zero, PadReader is a no-op, i.e. it does not attempt to
// add padding to the underlying reader.
func NewPadReader(r io.Reader, blockSize uint8) *PadReader {
	return &PadReader{
		reader:    r,
		blockSize: blockSize,
	}
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (r *PadReader) Read(p []byte) (int, error) {
	if r.fillByte != 0 {
		return r.pad(p)
	}

	n, err := r.reader.Read(p)
	if r.blockSize == 0 {
		return n, err
	}
	if n > 0 {
		r.readCount += uint64(n)
	}
	if !errors.Is(err, io.EOF) {
		return n, err
	}

	r.padding = r.blockSize - uint8(r.readCount%uint64(r.blockSize))
	r.fillByte = byte(r.padding)

	nn, err := r.pad(p[n:])
	return n + nn, err
}

// pad writes padding to p.
func (r *PadReader) pad(p []byte) (int, error) {
	if r.padding == 0 {
		return 0, io.EOF
	}

	var err error

	n := len(p)
	if int(r.padding) <= n {
		n = int(r.padding)
		r.padding = 0
		err = io.EOF
	} else {
		r.padding -= uint8(n)
	}

	for i := 0; i < n; i++ {
		p[i] = r.fillByte
	}

	return n, err
}
