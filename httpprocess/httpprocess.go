// Package httpprocess provides [process.Runnable] wrapper for [http.Server].
package httpprocess

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"

	"go.pact.im/x/httptrack"
	"go.pact.im/x/process"
)

// Server returns a [process.Runnable] instance for the given HTTP server and
// network listener.
func Server(srv *http.Server, lis net.Listener) process.Runnable {
	var connTracker httptrack.ConnTracker
	oldConnState := httptrack.Wrap(srv, &connTracker)
	wrappedListener := &nilCloserListener{Listener: lis}
	return process.Leaf(
		func(ctx context.Context) error {
			oldBaseContext := srv.BaseContext
			srv.BaseContext = func(_ net.Listener) context.Context {
				return ctx
			}

			err := srv.Serve(wrappedListener)
			if errors.Is(err, http.ErrServerClosed) {
				err = nil
			}
			if err == nil {
				err = wrappedListener.CloseError()
			}

			// Wait until all connections are actually closed (i.e. until all
			// per-connection goroutines return).
			//
			// Otherwise, process.Leaf would cancel the context on return, and
			// the cancellation would propagate to in-flight HTTP handlers. We
			// avoid that because our graceful shutdown model assumes the
			// context only expires during forced shutdown.
			//
			// Note that this can still happen if the shutdown function returns,
			// but by that point, graceful shutdown would have already failed.
			connTracker.Wait()

			srv.BaseContext = oldBaseContext
			srv.ConnState = oldConnState

			return err
		},
		func(ctx context.Context) error {
			// Shutdown and Close forward error from net.Listener.Close.
			// We already capture this error via nilCloserListener.
			_ = srv.Shutdown(ctx)
			_ = srv.Close()
			return nil
		},
	)
}

// nilCloserListener is a wrapper around a [net.Listener] that ensures the
// underlying listener is closed exactly once and returns nil error.
type nilCloserListener struct {
	net.Listener
	once sync.Once
	err  error
}

// Close closes the underlying [net.Listener] exactly once. It always returns
// nil error.
func (l *nilCloserListener) Close() error {
	l.once.Do(func() { l.err = l.Listener.Close() })
	return nil
}

// CloseError returns Close error from the underlying [net.Listener].
func (l *nilCloserListener) CloseError() error {
	return l.err
}
