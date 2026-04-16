// Package pgxslog implements [tracelog.Logger] adapter for [slog.Handler].
package pgxslog

import (
	"cmp"
	"context"
	"log/slog"
	"runtime"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
)

var _ tracelog.Logger = (*Logger)(nil)

// Logger is an adapter from [slog.Handler] to [tracelog.Logger].
type Logger struct {
	handler slog.Handler
	timeNow func() time.Time
}

// NewLogger return a new [Logger] instance.
func NewLogger(handler slog.Handler, timeNow func() time.Time) *Logger {
	return &Logger{
		handler: handler,
		timeNow: timeNow,
	}
}

// Log implements the [tracelog.Logger] interface.
func (l *Logger) Log(
	ctx context.Context,
	traceLogLevel tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	level := translateTraceLogLevel(traceLogLevel)
	if !l.handler.Enabled(ctx, level) {
		return
	}

	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}
	slices.SortFunc(attrs, func(a, b slog.Attr) int {
		return cmp.Compare(a.Key, b.Key)
	})

	// Skip 4 stack frames to log pgx.QueryTracer caller:
	//   - runtime.Callers
	//   - this function
	//   - this function’s caller (TraceLog.log helper)
	//   - QueryTracer hook
	var pcs [1]uintptr
	runtime.Callers(4, pcs[:])

	r := slog.NewRecord(l.now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = l.handler.Handle(ctx, r)
}

func (l *Logger) now() time.Time {
	if l.timeNow == nil {
		return time.Now()
	}
	return l.timeNow()
}

func translateTraceLogLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace:
		return slog.LevelDebug - 1
	case tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	case tracelog.LogLevelError:
		return slog.LevelError
	}
	return 0
}
