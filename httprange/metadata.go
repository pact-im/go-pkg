package httprange

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.pact.im/x/httpclient"
)

const (
	httpHeaderETag         = "Etag"
	httpHeaderLastModified = "Last-Modified"
	httpHeaderDate         = "Date"
	weakETagPrefix         = "W/"
)

// HTTPMetadataProvider provides representation metadata for a resource.
type HTTPMetadataProvider interface {
	Provide() *HTTPMetadata
}

// HTTPMetadata contains representation metadata for a resource.
type HTTPMetadata struct {
	// ContentLength is the complete length of the representation in bytes.
	ContentLength int64

	// ETag is the current representation entity tag.
	ETag string

	// LastModified is the last modification time.
	LastModified time.Time

	// Date is the server date when the metadata was retrieved.
	Date time.Time

	// AcceptRanges contains advertised Accept-Ranges header values.
	AcceptRanges []string
}

// Provide implements the [HTTPMetadataProvider] interface.
func (m *HTTPMetadata) Provide() *HTTPMetadata {
	return m
}

// HTTPMetadataExtractor extracts representation metadata for a resource.
type HTTPMetadataExtractor interface {
	// Extract retrieves metadata for the resource. It returns the metadata
	// or an error if the extraction fails.
	Extract(context.Context) (HTTPMetadataProvider, error)
}

// HTTPResponseMetadataExtractor implements [HTTPMetadataExtractor] by
// performing an HTTP HEAD request to retrieve representation metadata.
type HTTPResponseMetadataExtractor struct {
	// Request is a builder for HTTP HEAD request that is used to retrieve
	// representation metadata.
	Request HTTPRequestBuilder

	// Client is the HTTP client instance to use.
	Client httpclient.Client
}

// Extract performs a HEAD request and extracts metadata from the response.
func (e *HTTPResponseMetadataExtractor) Extract(ctx context.Context) (HTTPMetadataProvider, error) {
	req, err := e.Request.Build(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}

	closeBody := func() { _ = resp.Body.Close() }
	stop := context.AfterFunc(ctx, closeBody)
	defer func() {
		if stop() {
			closeBody()
		}
	}()

	m, err := extractHTTPResponseMetadata(resp)
	if err != nil {
		return nil, &HTTPResponseError{
			Response: resp,
			cause:    err,
		}
	}
	return m, nil
}

func extractHTTPResponseMetadata(resp *http.Response) (*HTTPMetadata, error) {
	if resp.StatusCode/100 != 2 {
		return nil, &UnexpectedStatusCodeError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
		}
	}

	eTag, err := parseETagOrZero(resp.Header)
	if err != nil {
		return nil, err
	}

	lastModified, err := parseTimeHeaderOrZero(resp.Header, httpHeaderLastModified)
	if err != nil {
		return nil, err
	}

	date, err := parseTimeHeaderOrZero(resp.Header, httpHeaderDate)
	if err != nil {
		return nil, err
	}

	return &HTTPMetadata{
		ContentLength: resp.ContentLength,
		ETag:          eTag,
		LastModified:  lastModified,
		Date:          date,
		AcceptRanges:  resp.Header.Values(httpHeaderAcceptRanges),
	}, nil
}

func parseETagOrZero(h http.Header) (string, error) {
	value := h.Get(httpHeaderETag)
	if value == "" {
		return "", nil
	}
	if !isValidETag(value) {
		return "", fmt.Errorf(
			"httprange: invalid ETag header value %q",
			value,
		)
	}
	return value, nil
}

func isValidETag(s string) bool {
	s = strings.TrimPrefix(s, weakETagPrefix)
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return false
	}
	for i := 1; i < len(s)-1; i++ {
		// VCHAR except double quotes, plus obs-text.
		if c := s[i]; c <= 0x20 || c == '"' || c == 0x7F {
			return false
		}
	}
	return true
}

func parseTimeHeaderOrZero(h http.Header, key string) (time.Time, error) {
	value := h.Get(key)
	if value == "" {
		return time.Time{}, nil
	}
	t, err := http.ParseTime(value)
	if err != nil {
		return time.Time{}, fmt.Errorf("httprange: parse %q header: %w", key, err)
	}
	return t, nil
}
