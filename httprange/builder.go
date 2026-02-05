package httprange

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"go.pact.im/x/httpclient"
)

// Builder constructs a [BytesResource] for an HTTP resource that supports byte
// range requests.
type Builder struct {
	// Metadata extracts representation metadata for the resource.
	Metadata HTTPMetadataExtractor

	// Validator builds a validator for conditional requests.
	Validator HTTPValidatorBuilder

	// Request builds the HTTP GET request for range requests.
	Request HTTPRequestBuilder

	// Client is the HTTP client instance to use.
	Client httpclient.Client
}

// BuildFromRawURL is a convenience function that builds a [BytesResource] from
// the parsed URL. If client is nil, the default HTTP client is used.
func BuildFromRawURL(ctx context.Context, rawURL string, client httpclient.Client) (*BytesResource, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return NewBuilderFromURL(u, client).Build(ctx)
}

// NewBuilderFromURL returns a new [Builder] for the given URL and HTTP client.
// If client is nil, the default HTTP client is used.
func NewBuilderFromURL(u *url.URL, client httpclient.Client) *Builder {
	if client == nil {
		client = httpclient.DefaultClient()
	}
	return &Builder{
		Metadata: &HTTPResponseMetadataExtractor{
			Request: &HTTPRequest{
				Method: http.MethodHead,
				URL:    u,
			},
			Client: client,
		},
		Validator: &HTTPStrongValidatorBuilder{
			UseLastModified: false,
		},
		Request: &HTTPRequest{
			URL: u,
		},
		Client: client,
	}
}

// Build creates a [BytesResource] for the configured HTTP resource.
func (b *Builder) Build(ctx context.Context) (*BytesResource, error) {
	mp, err := b.Metadata.Extract(ctx)
	if err != nil {
		return nil, err
	}

	m := mp.Provide()

	if m.ContentLength < 0 {
		return nil, errors.New("httprange: unknown resource content length")
	}

	if err := checkUnitIsAccepted(m.AcceptRanges, bytesUnit); err != nil {
		return nil, err
	}

	client, err := b.Validator.Build(ctx, mp, b.Client)
	if err != nil {
		return nil, err
	}

	return &BytesResource{
		Length: m.ContentLength,
		Ranger: &HTTPRanger{
			Request: b.Request,
			Client:  client,
		},
	}, nil
}
