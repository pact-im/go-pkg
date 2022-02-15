package flaky

import (
	"context"
	"errors"
)

// PermanentError indicates that a permanent error has been encountered during a
// retry attempt.
type PermanentError struct {
	err    error
	unwrap bool
}

// Internal marks an error as a permanent internal error. Executors unwrap
// internal errors prior to returning them.
func Internal(err error) *PermanentError {
	return &PermanentError{err, true}
}

// Permanent marks an error as permanent. That is, returning permanent error
// from retry attempt would stop the retry process.
func Permanent(err error) *PermanentError {
	return &PermanentError{err, false}
}

// IsPermanentError returns whether err’s error chain contains PermanentError.
func IsPermanentError(err error) bool {
	_, ok := AsPermanentError(err)
	return ok
}

// AsPermanentError attempts to extract PermanentError from err’s error chain.
func AsPermanentError(err error) (*PermanentError, bool) {
	var e *PermanentError
	ok := errors.As(err, &e)
	return e, ok
}

// Error implements the error interface.
func (e *PermanentError) Error() string {
	return "permanent: " + e.err.Error()
}

// Unwrap unwraps the underlying error.
func (e *PermanentError) Unwrap() error {
	return e.err
}

// Internal returns whether the permanent error is internal. It is guaranteed
// that Executor’s Execute method does not return internal errors.
func (e *PermanentError) Internal() bool {
	return e.unwrap
}

// unwrapInternal unwraps the underlying error of an internal PermanentError
// from err’s error chain.
func unwrapInternal(err error) error {
	p, ok := AsPermanentError(err)
	if !ok {
		return err
	}
	if !p.Internal() {
		return err
	}
	return p.Unwrap()
}

type untilPermanentExecutor struct{}

// UntilPermanent returns a new executor that retries operations until either a
// permanent error or context expiration.
func UntilPermanent() Executor {
	return (*untilPermanentExecutor)(nil)
}

// Execute implements the Executor interface.
func (*untilPermanentExecutor) Execute(ctx context.Context, f Op) error {
	var err error
	for {
		err = f(ctx)
		if err == nil {
			return nil
		}
		if IsPermanentError(err) {
			return unwrapInternal(err)
		}
		if ctx.Err() != nil {
			break
		}
	}
	return err
}
