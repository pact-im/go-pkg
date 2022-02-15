package extraio

// ByteReader is an io.Reader that reads the same byte indefinitely.
type ByteReader byte

// Read implements io.Reader interface.
func (b ByteReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(b)
	}
	return len(p), nil
}
