package httptrack

import (
	"errors"
	"net/http"
	"sync"
	"testing"

	"go.uber.org/goleak"
)

type testStats struct {
	AcceptedTotal uint64
	ActiveTotal   uint64
	IdleTotal     uint64
	HijackedTotal uint64
	ClosedTotal   uint64
	Accepted      uint
	Active        uint
	Idle          uint
}

func statsEqual(a Stats, b testStats) bool {
	return b == testStats{
		AcceptedTotal: a.AcceptedTotal(),
		ActiveTotal:   a.ActiveTotal(),
		IdleTotal:     a.IdleTotal(),
		HijackedTotal: a.HijackedTotal(),
		ClosedTotal:   a.ClosedTotal(),
		Accepted:      a.Accepted(),
		Active:        a.Active(),
		Idle:          a.Idle(),
	}
}

func TestStatsTrackerAccepted(t *testing.T) {
	defer goleak.VerifyNone(t)

	var connTracker ConnTracker
	var statsTracker StatsTracker

	listener := newTestListener()

	server := &http.Server{
		ConnState: Compose(
			&statsTracker,
			&connTracker,
		),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	defer func() {
		if err := server.Close(); err != nil {
			panic(err)
		}
		wg.Wait()
		connTracker.Wait()
	}()

	stats := statsTracker.Stats()
	if !statsEqual(stats, testStats{}) {
		t.Errorf("expected empty stats, but got %+v", stats)
	}

	conn := listener.Pipe()
	// StateNew is set from the Accept loop,
	// so wait for the next Accept call.
	listener.Wait()

	stats = statsTracker.Stats()
	if !statsEqual(stats, testStats{
		Accepted:      1,
		AcceptedTotal: 1,
	}) {
		t.Errorf("expected http.StateNew connection, but got stats %+v", stats)
	}

	// Close the connection and wait for http.StateClosed transition.
	if err := conn.Close(); err != nil {
		panic(err)
	}
	connTracker.Wait()

	stats = statsTracker.Stats()
	if !statsEqual(stats, testStats{
		AcceptedTotal: 1,
		ClosedTotal:   1,
	}) {
		t.Errorf("expected http.StateClosed connection, but got %+v", stats)
	}
}

func TestStatsTrackerActive(t *testing.T) {
	defer goleak.VerifyNone(t)

	var connTracker ConnTracker
	var statsTracker StatsTracker

	ch := make(chan struct{})

	listener := newTestListener()

	server := &http.Server{
		ConnState: Compose(
			&statsTracker,
			&connTracker,
		),
		Handler: http.HandlerFunc(
			func(http.ResponseWriter, *http.Request) {
				ch <- struct{}{}
				ch <- struct{}{}
			},
		),
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: listener.Dial,
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	defer func() {
		if err := server.Close(); err != nil {
			panic(err)
		}
		wg.Wait()
		connTracker.Wait()
	}()

	go func() {
		_, err := client.Get("http://example.com")
		if err != nil {
			panic(err)
		}
		ch <- struct{}{}
	}()

	// Wait for the handler and check that connection is in Active state.
	<-ch

	stats := statsTracker.Stats()
	if !statsEqual(stats, testStats{
		AcceptedTotal: 1,
		ActiveTotal:   1,
		Active:        1,
	}) {
		t.Errorf("expected http.StateActive connection, but got stats %+v", stats)
	}

	// Wait for handler to return.
	<-ch

	// Wait for client to receive response.
	<-ch

	// Disable keep-alives and wait for the connection to close.
	server.SetKeepAlivesEnabled(false)
	connTracker.Wait()
	server.SetKeepAlivesEnabled(true)

	stats = statsTracker.Stats()
	if !statsEqual(stats, testStats{
		AcceptedTotal: 1,
		ActiveTotal:   1,
		IdleTotal:     1,
		ClosedTotal:   1,
	}) {
		t.Errorf("expected http.StateClosed connection, but got %+v", stats)
	}
}
