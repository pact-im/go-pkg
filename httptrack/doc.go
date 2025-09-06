// Package httptrack provides tools for tracking HTTP connection state
// transitions.
//
// It integrates with [net/http.Server] via the ConnState hook to support:
//   - Graceful shutdowns (via [ConnTracker])
//   - Connection statistics (via [StatsTracker])
//
// Trackers implement the [Tracker] interface and can be composed using
// [Compose] function.
package httptrack
