package extraio

import (
	"io"
)

// TailReader buffers the last n bytes read from the underlying io.Reader. That
// is, its Read method does not return last n bytes read. Use Tail method to get
// the underlying buffer on EOF.
//
// Note that its Read method may return zero byte count with a nil error if read
// bytes fit into the underlying buffer. Some io.Reader client implementations
// return io.ErrNoProgress error when many calls to Read have failed to return
// any data or error.
type TailReader struct {
	reader io.Reader

	cur int    // mutable
	buf []byte // mutable
}

// NewTailReader returns a new reader that buffers last n bytes read from r.
func NewTailReader(r io.Reader, n uint) *TailReader {
	var buf []byte
	if n != 0 {
		buf = make([]byte, 0, n)
	}
	return &TailReader{
		reader: r,
		buf:    buf,
	}
}

// Reset resets the reader’s state.
func (r *TailReader) Reset() {
	r.buf = r.buf[:0]
}

// Length returns the length of the underlying buffer. It is faster than calling
// len(r.Tail()) since the underlying buffer may not be contiguous and Tail has
// to linearize it to return a contiguous slice of bytes.
func (r *TailReader) Length() int {
	return len(r.buf)
}

// Tail returns the buffer with last read bytes. The underlying buffer may not
// be contiguous and Tail linearizes the contents before returning a contiguous
// byte slice.
func (r *TailReader) Tail() []byte {
	r.linearize()
	return r.buf
}

// Read implements the io.Reader interface. It reads from the underlying
// io.Reader but keeps the last n read bytes in the internal ring buffer.
func (r *TailReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n <= 0 || r.buf == nil {
		return n, err
	}

	size := cap(r.buf)
	used := len(r.buf)
	free := size - used

	switch {
	// Case 1. Data fits into the free buffer space.
	case n <= free:
		r.buf = append(r.buf, p[:n]...)
		n = 0
	// Case 2. Data exceeds buffer capacity.
	case n >= size:
		if free == 0 {
			rotate(p[:n], -size)
			head := r.buf[r.cur:]
			tail := r.buf[:r.cur]
			swap(p[len(head):], tail)
			swap(p, head)
		} else {
			r.buf = r.buf[:size]
			swap(r.buf, p[n-size:n])
			n -= free
			rotate(p[:n], -used)
		}
	// Case 3. Not enough free space but data fits into the buffer.
	default:
		if free == 0 {
			k := swap(r.buf[r.cur:], p[:n])
			r.cur = (r.cur + k) % size
			k = swap(r.buf[r.cur:], p[k:n])
			r.cur = (r.cur + k) % size
		} else {
			r.buf = append(r.buf, p[:free]...)
			copy(p, p[free:n])
			n -= free
			swap(r.buf, p[:n])
			r.cur = n
		}
	}
	return n, err
}

// linearize makes the underlying ring buffer contiguous.
func (r *TailReader) linearize() {
	rotate(r.buf, r.cur)
	r.cur = 0
}

// swap swaps elements between dst and src slices. It returns the number of
// swapped elements, that is, min(len(dst), len(src)).
func swap(dst, src []byte) int {
	n := min(len(dst), len(src))
	for i := 0; i < n; i++ {
		dst[i], src[i] = src[i], dst[i]
	}
	return n
}

// rotate performs in-place left rotation of the slice by n positions from the
// start for positive n. If n is negative, the slice is right rotated by the
// absolute value of n.
//
// Examples:
//
//   rotate([]byte("lohel"), 2)
//       => []byte("hello")
//
//   rotate([]byte("45123"), 2)
//       => []byte("12345")
//
//   rotate([]byte("34512"), -2)
//       => []byte("12345")
//
func rotate(s []byte, n int) {
	k := len(s)
	if k == 0 || k == 1 {
		return
	}

	n %= k
	if n == 0 || n == k {
		return
	}

	if n < 0 {
		n += k
	}

	reverse(s[:n])
	reverse(s[n:])
	reverse(s)
}

// reverse performs in-place reversal of the slice.
func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// min returns the minimum of the two integers.
func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}
