package httptrack_test

import (
	"fmt"
	"net/http"

	"go.pact.im/x/httptrack"
)

func ExampleWrap() {
	server := &http.Server{}

	var connTracker httptrack.ConnTracker
	var statsTracker httptrack.StatsTracker

	httptrack.Wrap(server, httptrack.Compose(
		&statsTracker,
		&connTracker,
	))
}

func ExampleConnTracker() {
	server := &http.Server{}

	var connTracker httptrack.ConnTracker
	server.ConnState = connTracker.Track

	// On shutdown:
	_ = server.Close()
	connTracker.Wait() // blocks until all connections are done
}

func ExampleStatsTracker() {
	server := &http.Server{}

	var statsTracker httptrack.StatsTracker
	server.ConnState = statsTracker.Track

	stats := statsTracker.Stats()
	fmt.Println(
		"Active:", stats.Active(),
		"Idle:", stats.Idle(),
		"Closed:", stats.ClosedTotal(),
	)
}
