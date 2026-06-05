package extraio

import (
	"io"
)

var (
	_ io.Reader     = ReaderFunc(nil)
	_ io.Writer     = WriterFunc(nil)
	_ io.WriterTo   = WriterToFunc(nil)
	_ io.ReaderFrom = ReaderFromFunc(nil)
)

// ReaderFunc is an adapter to allow the use of ordinary function as io.Reader.
// If f is a function with appropriate signature, ReaderFunc(f) is an io.Reader
// that calls f.
type ReaderFunc func(p []byte) (n int, err error)

// Read implements the io.Reader interface. It calls f(p).
func (f ReaderFunc) Read(p []byte) (n int, err error) {
	return f(p)
}

// WriterFunc is an adapter to allow the use of ordinary function as io.Writer.
// If f is a function with appropriate signature, WriterFunc(f) is an io.Writer
// that calls f.
type WriterFunc func(p []byte) (n int, err error)

// Write implements the io.Writer interface. It calls f(p).
func (f WriterFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

// WriterToFunc is an adapter to allow the use of ordinary function as io.WriterTo.
// If f is a function with appropriate signature, WriterToFunc(f) is an io.WriterTo
// that calls f.
type WriterToFunc func(w io.Writer) (n int64, err error)

// WriteTo implements the io.WriterTo interface. It calls f(w).
func (f WriterToFunc) WriteTo(w io.Writer) (n int64, err error) {
	return f(w)
}

// ReaderFromFunc is an adapter to allow the use of ordinary function as io.ReaderFrom.
// If f is a function with appropriate signature, ReaderFromFunc(f) is an io.ReaderFrom
// that calls f.
type ReaderFromFunc func(r io.Reader) (n int64, err error)

// ReadFrom implements the io.ReaderFrom interface. It calls f(r).
func (f ReaderFromFunc) ReadFrom(r io.Reader) (n int64, err error) {
	return f(r)
}
