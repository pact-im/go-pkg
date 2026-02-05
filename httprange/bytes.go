package httprange

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"math"
	"strconv"
)

const bytesUnit = "bytes"

// BytesSpecifier constructs a byte range specifier from a sequence of positions.
//
// The positions are interpreted as follows:
//   - An empty sequence produces the unsatisfiable range "-0".
//   - A negative integer N produces a suffix-range "-N".
//   - A non-negative integer N followed by a smaller or
//     negative integer produces an open range "N-".
//   - Two adjacent non-negative integers M, N where M ≤ N
//     produce a closed range "M-N".
//
// It preserves input order even if it produces overlapping ranges or ranges are
// not strictly in ascending order.
func BytesSpecifier(positions ...int64) Specifier {
	if len(positions) == 0 {
		return bytesUnit + "=-0" // unsatisfiable
	}
	buf := []byte(bytesUnit + "=")
	for i := 0; i < len(positions); i++ {
		if i != 0 {
			buf = append(buf, ',')
		}

		pos := positions[i]

		buf = strconv.AppendInt(buf, pos, 10)

		if pos < 0 {
			continue
		}

		buf = append(buf, '-')

		nextIndex := i + 1
		if nextIndex == len(positions) {
			break
		}

		next := positions[nextIndex]
		if next >= pos {
			i = nextIndex
			buf = strconv.AppendInt(buf, next, 10)
			if i+1 == len(positions) {
				break
			}
		}
	}
	return Specifier(buf)
}

// BytesReader implements [io.ReaderAt] for reading byte ranges from a [Ranger].
type BytesReader struct {
	// Context is the context to use for Ranger.Range calls.
	Context context.Context

	// Ranger is the Ranger instance for performing range requests.
	Ranger Ranger
}

// ReadAt reads len(p) bytes starting at byte offset off. It performs a byte
// range request for the range [off, off+len(p)-1] and reads the response into
// p. It returns [io.EOF] if the range cannot be satisfied.
func (r *BytesReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("httprange: negative read offset")
	}
	if len(p) == 0 {
		return 0, nil
	}
	if off > math.MaxInt64-int64(len(p))+1 {
		return 0, errors.New("httprange: read offset overflow")
	}

	first, last := off, off+int64(len(p))-1
	spec := BytesSpecifier(first, last)

	ctx, cancel := context.WithCancel(r.Context)
	next, stop := iter.Pull2(r.Ranger.Range(ctx, spec))
	defer stop()
	defer cancel()

	rr, err, ok := next()
	if !ok {
		return 0, errors.New("rangeio: empty sequence")
	}
	if err != nil {
		var e *UnsatisfiedRangeError
		if errors.As(err, &e) {
			// Assume EOF if complete length is unknown.
			if e.Resp == "" {
				return 0, io.EOF
			}

			ur, ok := parseUnsatisfiedRangeResp(e.Resp)
			if !ok {
				return 0, fmt.Errorf(
					"httprange: invalid unsatisfied range %q",
					e.Resp,
				)
			}

			if first < ur.CompleteLength {
				return 0, fmt.Errorf(
					"httprange: unexpected unsatisfied range %q (first byte %d is satisfiable)",
					e.Resp, first,
				)
			}

			return 0, io.EOF
		}
		return 0, err
	}

	br, ok := parseBytesRangeResp(rr.Resp)
	if !ok {
		return 0, fmt.Errorf(
			"httprange: invalid bytes range %q",
			rr.Resp,
		)
	}
	if br.First != first {
		return 0, fmt.Errorf(
			"httprange: unexpected first byte position %d (expected %d)",
			br.First, first,
		)
	}
	if br.Last > last || br.Last < last && br.CompleteLength > 0 && br.Last != br.CompleteLength-1 {
		return 0, fmt.Errorf(
			"httprange: unexpected last byte position %d (expected %d or less at EOF)",
			br.Last, last,
		)
	}

	rangeLength := int(br.Last - br.First + 1) // 0 < rangeLength ≤ len(p)
	atEOF := br.CompleteLength > 0 && br.Last == br.CompleteLength-1

	var n, nn int
	for n < rangeLength && err == nil {
		nn, err = rr.Reader.Read(p[n:])
		n += nn
	}
	if n == rangeLength && err == nil {
		var peek [1]byte // for EOF
		nn, err = rr.Reader.Read(peek[:])
		n += nn
	}
	if n != rangeLength {
		return 0, fmt.Errorf(
			"httprange: invalid range length %d (read %d bytes)",
			rangeLength, n,
		)
	}

	switch {
	case atEOF && err == nil:
		err = io.EOF
	case !atEOF && err == io.EOF:
		err = nil
	}

	cancel() // unblock reading response body
	if _, _, ok := next(); ok {
		return 0, errors.New("httprange: unexpected multiple ranges")
	}

	return n, err
}
