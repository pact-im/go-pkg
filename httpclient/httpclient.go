// Package httpclient defines a mockable http.Client interface.
package httpclient

import (
	"net/http"
)

// Client is an interface describing the methods required from HTTP client.
// It is satisfied by *http.Client type.
type Client interface {
	// Do sends an HTTP request and returns an HTTP response.
	Do(req *http.Request) (*http.Response, error)
}

// DefaultClient returns http.DefaultClient.
func DefaultClient() Client {
	return http.DefaultClient
}
