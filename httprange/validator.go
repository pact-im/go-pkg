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
	httpHeaderIfMatch           = "If-Match"
	httpHeaderIfUnmodifiedSince = "If-Unmodified-Since"
)

// HTTPValidatorBuilder builds HTTP clients that wraps [httpclient.Client] with
// resource state validation.
type HTTPValidatorBuilder interface {
	Build(context.Context, HTTPMetadataProvider, httpclient.Client) (httpclient.Client, error)
}

// HTTPStrongValidatorBuilder implements [HTTPValidatorBuilder] using strong
// validators. It prefers strong ETags over Last-Modified dates for validation.
type HTTPStrongValidatorBuilder struct {
	// UseLastModified controls whether to use Last-Modified dates when
	// no strong ETag is available. When false, returns an error if no
	// strong ETag is available.
	UseLastModified bool
}

// Build creates a validator client that adds either If-Match or If-Unmodified-Since
// headers to requests, depending on the available metadata.
//
// It returns an error if no applicable validator is available.
func (b *HTTPStrongValidatorBuilder) Build(_ context.Context, mp HTTPMetadataProvider, client httpclient.Client) (httpclient.Client, error) {
	switch m := mp.Provide(); {
	case m.ETag != "" && !strings.HasPrefix(m.ETag, weakETagPrefix):
		return &HTTPStrongValidator{
			Client: client,
			Precondition: http.Header{
				httpHeaderIfMatch: {m.ETag},
			},
			ETag: m.ETag,
		}, nil
	case b.UseLastModified && !m.LastModified.IsZero() && m.Date.After(m.LastModified):
		lastModified := m.LastModified.UTC().Format(http.TimeFormat)
		return &HTTPStrongValidator{
			Client: client,
			Precondition: http.Header{
				httpHeaderIfUnmodifiedSince: {lastModified},
			},
			LastModified: m.LastModified,
		}, nil
	default:
		return nil, errNoApplicableValidator
	}
}

// HTTPStrongValidator is an [httpclient.Client] that adds precondition headers
// to requests and validates that the resource hasn’t changed in responses.
type HTTPStrongValidator struct {
	// Client is the HTTP client instance to use.
	Client httpclient.Client

	// Precondition contains headers to add to requests.
	Precondition http.Header

	// ETag is the expected entity tag in responses.
	ETag string

	// LastModified is the expected Last-Modified time in responses.
	LastModified time.Time
}

// Do executes a conditional HTTP request and checks the response to ensure the
// resource hasn’t changed.
func (v *HTTPStrongValidator) Do(req *http.Request) (*http.Response, error) {
	req.Header = setHeader(req.Header, v.Precondition)

	resp, err := v.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := v.checkValidator(resp); err != nil {
		_ = resp.Body.Close()
		return nil, &HTTPResponseError{
			Response: resp,
			cause:    err,
		}
	}

	return resp, nil
}

func (v *HTTPStrongValidator) checkValidator(resp *http.Response) error {
	switch {
	case v.ETag != "":
		eTag, err := parseETagOrZero(resp.Header)
		if err != nil {
			return err
		}
		if v.ETag != eTag {
			return fmt.Errorf(
				"httprange: ETag validator changed from %q to %q",
				v.ETag, eTag,
			)
		}
	case !v.LastModified.IsZero():
		value := resp.Header.Get(httpHeaderLastModified)
		t, err := http.ParseTime(value)
		if err != nil {
			return err
		}
		if !t.Equal(v.LastModified) {
			return fmt.Errorf(
				"httprange: Last-Modified validator changed from %q to %q",
				v.LastModified, t,
			)
		}
	}
	return nil
}
