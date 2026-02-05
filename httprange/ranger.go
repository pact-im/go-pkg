package httprange

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"golang.org/x/net/http/httpguts"

	"go.pact.im/x/httpclient"
)

const (
	noneUnit                = "none"
	httpHeaderRange         = "Range"
	httpHeaderAcceptRanges  = "Accept-Ranges"
	httpHeaderContentRange  = "Content-Range"
	httpHeaderContentType   = "Content-Type"
	mimeMultipartByteranges = "multipart/byteranges"
)

// Specifier represents a Range header value. It consists of a unit and a range
// set (e.g. "bytes=0-99,100-200").
type Specifier string

// Unit returns the unit of the range.
func (s Specifier) Unit() string {
	unit, _, _ := strings.Cut(string(s), "=")
	return unit
}

// Range represents a single range response.
type Range struct {
	// Resp is the range specification for the content.
	Resp string
	// Reader is the reader for the range content.
	Reader io.Reader
}

// Ranger performs range requests and returns the requested ranges as a sequence.
type Ranger interface {
	// Range returns a non-empty sequence of ranges for the given specifier.
	// Iteration yields Range values until an error. The [Range.Reader] is
	// valid only until the next iteration.
	Range(ctx context.Context, spec Specifier) iter.Seq2[Range, error]
}

// HTTPRanger implements the [Ranger] interface for HTTP resources.
type HTTPRanger struct {
	// Request is a builder for HTTP GET requests.
	Request HTTPRequestBuilder

	// Client is the HTTP client instance to use.
	Client httpclient.Client
}

// Range performs a range request for the given specifier. It sends an HTTP
// request with a Range header and returns an iterator over the ranges in the
// response (a single range or a multipart/byteranges parts).
func (r *HTTPRanger) Range(ctx context.Context, spec Specifier) iter.Seq2[Range, error] {
	return func(yield func(Range, error) bool) {
		req, err := r.Request.Build(ctx)
		if err != nil {
			_ = yield(Range{}, err)
			return
		}

		req.Header = setHeader(req.Header, http.Header{
			httpHeaderRange: {string(spec)},
		})

		resp, err := r.Client.Do(req)
		if err != nil {
			_ = yield(Range{}, err)
			return
		}

		closeBody := func() { _ = resp.Body.Close() }
		stop := context.AfterFunc(ctx, closeBody)
		defer func() {
			if stop() {
				closeBody()
			}
		}()

		if err := parseHTTPResponse(
			resp,
			spec.Unit(),
			func(rr Range) bool {
				return yield(rr, nil)
			},
		); err != nil {
			_ = yield(Range{}, &HTTPResponseError{
				Response: resp,
				cause:    err,
			})
		}
	}
}

func parseHTTPResponse(resp *http.Response, unit string, yield func(Range) bool) error {
	// Note that we do not support Accept-Ranges in a trailer section.
	if err := checkAcceptRangesHeader(resp.Header, unit); err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusPartialContent:
		// OK
	case http.StatusRequestedRangeNotSatisfiable:
		if resp.Header.Get(httpHeaderContentRange) == "" {
			return &UnsatisfiedRangeError{}
		}
		r, err := parseContentRangeFromHeader(resp.Header, unit)
		if err != nil {
			return errors.Join(
				&UnsatisfiedRangeError{},
				err,
			)
		}
		return &UnsatisfiedRangeError{Resp: r}
	default:
		return &UnexpectedStatusCodeError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
		}
	}

	if contentType := resp.Header.Get(httpHeaderContentType); contentType != "" {
		typ, params, err := mime.ParseMediaType(contentType)
		// NB we ignore errors in other content types.
		if typ == mimeMultipartByteranges {
			if err != nil {
				return fmt.Errorf(
					"httprange: parse Content-Type header: %w",
					err,
				)
			}
			reader := multipart.NewReader(resp.Body, params["boundary"])
			return readMultipartRanges(reader, unit, yield)
		}
	}

	r, err := parseContentRangeFromHeader(resp.Header, unit)
	if err != nil {
		return err
	}
	_ = yield(Range{
		Resp:   r,
		Reader: resp.Body,
	})
	return nil
}

func readMultipartRanges(r *multipart.Reader, unit string, yield func(Range) bool) error {
	for empty := true; ; empty = false {
		part, err := r.NextPart()
		if err == io.EOF {
			if empty {
				return errors.New("httprange: empty multipart ranges")
			}
			return nil
		}
		if err != nil {
			return err
		}

		resp, err := parseContentRangeFromHeader(http.Header(part.Header), unit)
		if err != nil {
			return err
		}
		if !yield(Range{Resp: resp, Reader: part}) {
			return nil
		}
	}
}

func checkAcceptRangesHeader(h http.Header, unit string) error {
	return checkUnitIsAccepted(h.Values(httpHeaderAcceptRanges), unit)
}

func checkUnitIsAccepted(values []string, unit string) error {
	if len(values) == 0 {
		return nil // assume it is accepted
	}
	for _, v := range values {
		if equalFoldASCII(v, noneUnit) {
			return errRangesNotSupported
		}
	}
	if !httpguts.HeaderValuesContainsToken(values, unit) {
		return &UnacceptedUnitError{
			AcceptRanges: values,
		}
	}
	return nil
}

func equalFoldASCII(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range len(a) {
		if lowerASCII(a[i]) != lowerASCII(b[i]) {
			return false
		}
	}
	return true
}

// lowerASCII returns the ASCII lowercase version of b.
func lowerASCII(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
