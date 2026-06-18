// Package logs provides [Logger], an alternative [slog] frontend with
// configurable time source and program counter capture.
//
// Logger has no convenience methods. Call [Logger.Log] directly with a level,
// context, and attributes. Use [Logger.WithTime] to change the time source,
// and [Logger.WithCapturePC] / [Logger.WithSkipPC] to control PC capture.
package logs
