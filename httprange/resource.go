package httprange

import (
	"context"
	"errors"
	"io"
)

var _ interface {
	io.ReaderAt
	io.Reader
	io.Seeker
	io.Closer
} = (*BytesResourceReader)(nil)

// BytesResource represents a resource with a known length that supports byte
// range requests.
type BytesResource struct {
	// Length is the complete length of the resource in bytes.
	Length int64

	// Ranger is the Ranger instance for performing range requests.
	Ranger Ranger
}

// Reader returns a [BytesResourceReader] for reading the resource.
//
// The provided context controls the lifetime of the reader; canceling it will
// abort any in-flight range requests.
//
// The returned reader must be closed after use by calling its Close method.
func (r *BytesResource) Reader(ctx context.Context) *BytesResourceReader {
	ctx, cancel := context.WithCancelCause(ctx)
	reader := &BytesReader{
		Context: ctx,
		Ranger:  r.Ranger,
	}
	return &BytesResourceReader{
		Reader: reader,
		Cancel: cancel,
		Length: r.Length,
	}
}

// BytesResourceReader provides sequential and random access to a [BytesResource].
// It implements [io.ReaderAt], [io.Reader], [io.Seeker], and [io.Closer].
type BytesResourceReader struct {
	// Reader is the underlying resource reader.
	Reader io.ReaderAt

	// Cancel is a function that cancels reads.
	Cancel context.CancelCauseFunc

	// Length is the complete length of the resource in bytes.
	Length int64

	// Offset is the current offset for reading and seeking.
	Offset int64 // mutable
}

// ReadAt reads len(p) bytes from offset off into p.
func (r *BytesResourceReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("httprange: negative read offset")
	}
	if off >= r.Length {
		return 0, io.EOF
	}
	p, eof := r.slice(p, off)
	n, err := r.Reader.ReadAt(p, off)
	switch {
	case eof && err == nil:
		err = io.EOF
	case !eof && err == io.EOF:
		err = io.ErrUnexpectedEOF
	}
	return n, err
}

func (r *BytesResourceReader) slice(p []byte, off int64) ([]byte, bool) {
	n := r.Length - off
	if int64(len(p)) > n {
		return p[:n], true
	}
	return p, int64(len(p)) == n
}

// Read reads up to len(p) bytes into p.
func (r *BytesResourceReader) Read(p []byte) (int, error) {
	n, err := r.ReadAt(p, r.Offset)
	r.Offset += int64(n)
	return n, err
}

// Seek sets the offset for the next read.
func (r *BytesResourceReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		// OK
	case io.SeekCurrent:
		offset += r.Offset
	case io.SeekEnd:
		offset += r.Length
	default:
		return 0, errors.New("httprange: invalid seek whence")
	}
	if offset < 0 || offset > r.Length {
		return 0, errors.New("httprange: invalid seek offset")
	}
	r.Offset = offset
	return offset, nil
}

// Close releases resources associated with the reader and cancels any pending
// range requests with [BytesResourceReaderClosedError] cause.
// After Close returns, all subsequent operations on the reader will fail.
func (r *BytesResourceReader) Close() error {
	r.Cancel(errBytesResourceReaderClosed)
	return nil
}
