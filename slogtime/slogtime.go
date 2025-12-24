// Package slogtime provides an alternative [slog.Logger] implementation that
// allows using custom [time.Now] function for log records.
//
// It also has a smaller API surface that requires using explicit [slog.Level],
// [context.Context] and [slog.Attr]s.
package slogtime

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// Logger is an alternative implementation of [slog.Logger] that allows
// customizing the time source for log records.
//
// Unlike [slog.Logger], Logger does not provide convenience methods like
// Info, Debug, Warn, or Error. Instead, all logging is done through the
// Log method with explicit level, context, and attributes.
type Logger struct {
	handler slog.Handler
	now     func() time.Time
}

// New creates a new [Logger] that uses the given handler and time function.
// If now is nil, [time.Now] will be used when logging.
func New(h slog.Handler, now func() time.Time) *Logger {
	if now == nil {
		now = time.Now
	}
	return &Logger{
		handler: h,
		now:     now,
	}
}

// clone returns a copy of the Logger.
func (l *Logger) clone() *Logger {
	c := *l
	return &c
}

// Handler returns the underlying [slog.Handler].
func (l *Logger) Handler() slog.Handler {
	return l.handler
}

// Enabled reports whether l emits log records at the given context and level.
func (l *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.handler.Enabled(ctx, level)
}

// With returns a [Logger] that includes the given attributes in each output
// operation. Arguments are converted to attributes as if by [Logger.Log].
func (l *Logger) With(attrs ...slog.Attr) *Logger {
	if len(attrs) == 0 {
		return l
	}
	c := l.clone()
	c.handler = l.handler.WithAttrs(attrs)
	return c
}

// WithGroup returns a [Logger] that starts a group. The keys of all attributes
// added to the [Logger] will be qualified by the given name.
func (l *Logger) WithGroup(name string) *Logger {
	if name == "" {
		return l
	}
	c := l.clone()
	c.handler = l.handler.WithGroup(name)
	return c
}

// WithTime returns a [Logger] that uses the given time function.
// If now is nil, [time.Now] will be used.
func (l *Logger) WithTime(now func() time.Time) *Logger {
	if now == nil {
		now = time.Now
	}
	c := l.clone()
	c.now = now
	return c
}

// Log emits a log record with the given level, message, and attributes.
// The record includes a source location obtained from the call stack.
func (l *Logger) Log(
	ctx context.Context,
	level slog.Level,
	msg string,
	attrs ...slog.Attr,
) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(l.now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = l.handler.Handle(ctx, r)
}
