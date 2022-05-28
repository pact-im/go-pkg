package extraio

// ByteReader is an io.Reader that reads the same byte indefinitely.
type ByteReader byte

// Read implements the io.Reader interface. It fills p with the b byte and
// returns len(p).
func (b ByteReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(b)
	}
	return len(p), nil
}
