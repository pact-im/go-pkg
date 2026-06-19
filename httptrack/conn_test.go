package httptrack

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"
	"testing/synctest"
)

func TestConnTracker(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var connTracker ConnTracker

		ch := make(chan struct{})

		listener := newTestListener()

		handler := func(http.ResponseWriter, *http.Request) {
			ch <- struct{}{}
			ch <- struct{}{}
		}

		server := &http.Server{
			ConnState: connTracker.Track,
			Handler:   http.HandlerFunc(handler),
		}

		client := &http.Client{
			Transport: &http.Transport{
				DialContext: listener.Dial,
			},
		}

		var wg sync.WaitGroup
		wg.Go(func() {
			err := server.Serve(listener)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		})
		wg.Go(func() {
			_, err := client.Get("http://example.com")
			if err != nil && !errors.Is(err, io.EOF) {
				panic(err)
			}
		})

		// Wait for handler to receive a request.
		<-ch

		// Shutdown the server while there is an in-flight request.
		if err := server.Close(); err != nil {
			t.Fatal(err)
		}

		// Wait for both client and server to return.
		wg.Wait()

		// Unblock handler.
		<-ch

		// Wait for connection http.StateClosed transition.
		connTracker.Wait()
	})
}
