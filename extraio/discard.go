package extraio

import (
	"io"
)

// DiscardReader is an io.Reader that discard all reads from the underlying reader.
type DiscardReader struct {
	R io.Reader // underlying reader
}

// NewDiscardReader returns a new reader that discard all reads from r.
func NewDiscardReader(r io.Reader) *DiscardReader {
	return &DiscardReader{r}
}

// Read implements io.Reader interface. It reads from the underlying io.Reader.
func (d *DiscardReader) Read(p []byte) (int, error) {
	_, err := io.Copy(io.Discard, d.R)
	if err == nil {
		err = io.EOF
	}
	return 0, err
}
