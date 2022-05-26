package extraio

// ReaderFunc is an adapter to allow the use of ordinary function as io.Reader.
// If f is a function with appropriate signature, ReaderFunc(f) is an io.Reader
// that calls f.
type ReaderFunc func(p []byte) (n int, err error)

// Read implements io.Reader interface. It calls f(p).
func (f ReaderFunc) Read(p []byte) (n int, err error) {
	return f(p)
}
