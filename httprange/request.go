package httprange

import (
	"context"
	"maps"
	"net/http"
)

// HTTPRequestBuilder is a factory for [http.Request] instances.
type HTTPRequestBuilder interface {
	// Build creates a new HTTP request for the given context.
	Build(context.Context) (*http.Request, error)
}

// HTTPRequestBuilderFunc is a function that implements [HTTPRequestBuilder]
// interface.
type HTTPRequestBuilderFunc func(context.Context) (*http.Request, error)

// Build implements the [HTTPRequestBuilder] interface.
func (f HTTPRequestBuilderFunc) Build(ctx context.Context) (*http.Request, error) {
	return f(ctx)
}

// HTTPRequest is an [HTTPRequestBuilder] implementation that wraps an existing
// [http.Request].
type HTTPRequest http.Request

// Build returns a shallow copy of the request with the given context.
// It implements the [HTTPRequestBuilder] interface.
func (r *HTTPRequest) Build(ctx context.Context) (*http.Request, error) {
	return (*http.Request)(r).WithContext(ctx), nil
}

func setHeader(target, source http.Header) http.Header {
	if target == nil {
		return source
	}
	target = target.Clone()
	maps.Copy(target, source)
	return target
}
