package httprange

import (
	"fmt"
	"net/http"
)

// HTTPResponseError wraps an error that occurred while processing an HTTP response.
type HTTPResponseError struct {
	Response *http.Response
	cause    error
}

// Error returns the error message. It implements the [error] interface.
func (e *HTTPResponseError) Error() string {
	return e.cause.Error()
}

// Unwrap returns the underlying error.
func (e *HTTPResponseError) Unwrap() error {
	return e.cause
}

// UnexpectedStatusCodeError indicates that the server returned an unexpected
// HTTP status code.
type UnexpectedStatusCodeError struct {
	Status     string
	StatusCode int
}

// Error returns the error message. It implements the [error] interface.
func (e *UnexpectedStatusCodeError) Error() string {
	text := e.Status
	if text == "" {
		text = http.StatusText(e.StatusCode)
	}
	return fmt.Sprintf(
		"httprange: unexpected status code: %d %s",
		e.StatusCode, text,
	)
}

// UnsatisfiedRangeError indicates that the requested range cannot be satisfied.
// This corresponds to the HTTP status code 416 (Range Not Satisfiable).
type UnsatisfiedRangeError struct {
	// Resp contains unsatisfied range, if provided by the server.
	Resp string
}

// Error returns the error message. It implements the [error] interface.
func (e *UnsatisfiedRangeError) Error() string {
	const errorMessage = "httprange: range is not satisfiable"
	if e.Resp != "" {
		return fmt.Sprintf(errorMessage+" (%q)", e.Resp)
	}
	return errorMessage
}

var errRangesNotSupported *RangesNotSupportedError

// RangesNotSupportedError indicates that the server does not support range
// requests (response contains "Accept-Ranges: none" header).
type RangesNotSupportedError struct{}

// Error returns the error message. It implements the [error] interface.
func (*RangesNotSupportedError) Error() string {
	return "httprange: ranges are not supported"
}

// UnacceptedUnitError indicates that the server does not accept the requested
// range unit.
type UnacceptedUnitError struct {
	// AcceptRanges contains advertised Accept-Ranges header values.
	AcceptRanges []string
}

// Error returns the error message. It implements the [error] interface.
func (*UnacceptedUnitError) Error() string {
	return "httprange: unit is not acceptable"
}

var errBytesResourceReaderClosed *BytesResourceReaderClosedError

// BytesResourceReaderClosedError indicates that [BytesResourceReader] was closed.
type BytesResourceReaderClosedError struct{}

// Error returns the error message. It implements the [error] interface.
func (*BytesResourceReaderClosedError) Error() string {
	return "httprange: bytes resource reader was closed"
}

var errNoApplicableValidator *NoApplicableValidatorError

// NoApplicableValidatorError indicates that the server did not respond with
// a validator that can be used for future requests.
type NoApplicableValidatorError struct{}

// Error returns the error message. It implements the [error] interface.
func (*NoApplicableValidatorError) Error() string {
	return "httprange: no applicable validator in response"
}
