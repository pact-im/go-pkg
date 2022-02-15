package extraio

import (
	"errors"
	"io"
)

// UnpadReader is an io.Reader that unpads padding from PadReader. It validates
// the padding on EOF and returns an error if it is invalid.
type UnpadReader struct {
	reader    TailReader
	blockSize uint8
	readCount uint64
	tail      []byte
}

// NewUnpadReader returns a new reader that unpads r using the given block size.
// If blockSize is zero, UnpadReader is a no-op, i.e. it does not attempt to
// remove padding from the underlying reader.
func NewUnpadReader(r io.Reader, blockSize uint8) *UnpadReader {
	return &UnpadReader{
		reader: TailReader{
			r: r,
			n: uint64(blockSize),
		},
		blockSize: blockSize,
	}
}

// Reset resets the reader’s state.
func (r *UnpadReader) Reset() {
	r.readCount = 0
	r.tail = nil
	r.reader.Reset()
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (r *UnpadReader) Read(p []byte) (int, error) {
	if r.tail != nil {
		return r.unpad(p)
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

	tail := r.reader.Tail()
	blockSize := uint64(r.blockSize)

	// Check that stream is divisible into blocks and we have at least one block.
	if r.readCount%blockSize != 0 || uint64(len(tail))%blockSize != 0 || len(tail) == 0 {
		return 0, io.ErrUnexpectedEOF
	}

	// Check that padding is within block size.
	fillByte := tail[len(tail)-1]
	if fillByte > r.blockSize || fillByte == 0 {
		return 0, io.ErrUnexpectedEOF
	}

	// Check that padding is filled with same bytes.
	payload, ok := unpadPayload(tail, fillByte)
	if !ok {
		return 0, io.ErrUnexpectedEOF
	}

	r.tail = payload
	nn, err := r.unpad(p[n:])
	return n + nn, err
}

// unpad writes remaining payload to p.
func (r *UnpadReader) unpad(p []byte) (int, error) {
	n := copy(p, r.tail)
	r.tail = r.tail[n:]
	if len(r.tail) == 0 {
		r.Reset()
		return n, io.EOF
	}
	return n, nil
}

// unpadPayload validates that buf’s padding is consists of fillByte bytes an
// returns the unpadded payload.
func unpadPayload(buf []byte, fillByte byte) ([]byte, bool) {
	n := int(fillByte)
	if len(buf) < n {
		return nil, false
	}
	offset := len(buf) - n

	padding := buf[offset:]
	for _, c := range padding {
		if c == fillByte {
			continue
		}
		return nil, false
	}

	return buf[:offset], true
}
