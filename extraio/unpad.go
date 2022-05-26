package extraio

import (
	"io"
)

// UnpadReader is an io.Reader that unpads padding from PadReader. It validates
// the padding on EOF and returns an error if it is invalid.
type UnpadReader struct {
	reader    TailReader
	blockSize uint8

	incomplete int    // mutable
	lastBlock  []byte // mutable
}

// NewUnpadReader returns a new reader that unpads r using the given block size.
// If blockSize is zero, UnpadReader is a no-op, i.e. it does not attempt to
// remove padding from the underlying reader.
func NewUnpadReader(r io.Reader, blockSize uint8) *UnpadReader {
	var buf []byte
	if blockSize != 0 {
		buf = make([]byte, 0, blockSize)
	}
	return &UnpadReader{
		reader: TailReader{
			reader: r,
			buf:    buf,
		},
		blockSize: blockSize,
	}
}

// Reset resets the reader’s state.
func (r *UnpadReader) Reset() {
	r.incomplete = 0
	r.lastBlock = nil
	r.reader.Reset()
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (r *UnpadReader) Read(p []byte) (int, error) {
	if r.lastBlock != nil {
		return r.unpad(p)
	}

	n, err := r.reader.Read(p)
	if r.blockSize == 0 {
		return n, err
	}
	if n > 0 {
		bs := int(r.blockSize)
		r.incomplete += n % bs
		r.incomplete %= bs
	}
	if err != io.EOF {
		return n, err
	}

	// Check that all read bytes are divisible into blocks (i.e. we are not
	// at the incomplete block).
	if r.incomplete != 0 {
		return n, io.ErrUnexpectedEOF
	}

	// We call Tail after checking for incomplete blocks in previously read
	// bytes since it may linearize the underlying bufffer before returning
	// it. This saves us a few CPU cycles needed to rotate the ring buffer
	// that would otherwise be unused.
	//
	// Note that 0 <= len(lastBlock) <= r.blockSize <= math.MaxUint8.
	lastBlock := r.reader.Tail()

	// Check that the last block is also complete. If not, we have an
	// unexpected EOF.
	if uint8(len(lastBlock)) != r.blockSize {
		return n, io.ErrUnexpectedEOF
	}

	// Check that padding fill byte is within the block size.
	fillByte := lastBlock[len(lastBlock)-1]
	if fillByte > r.blockSize || fillByte == 0 {
		return n, io.ErrUnexpectedEOF
	}

	// Check that padding is filled with same bytes.
	payload, ok := unpadPayload(lastBlock, fillByte)
	if !ok {
		return n, io.ErrUnexpectedEOF
	}

	r.lastBlock = payload
	nn, err := r.unpad(p[n:])
	return n + nn, err
}

// unpad writes remaining payload to p.
func (r *UnpadReader) unpad(p []byte) (int, error) {
	n := copy(p, r.lastBlock)
	r.lastBlock = r.lastBlock[n:]
	if len(r.lastBlock) == 0 {
		r.Reset()
		return n, io.EOF
	}
	return n, nil
}

// unpadPayload validates that buf is padded with fillByte bytes and returns the
// unpadded payload.
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
