package logs

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// Logger is an alternative implementation of [slog.Logger] that allows
// configuring time source and program counter for log records.
//
// Unlike [slog.Logger], Logger does not provide convenience methods like
// Info, Debug, Warn, or Error. Instead, all logging is done through the
// Log method with explicit context, level, and attributes.
type Logger struct {
	// handler is the underlying log handler.
	handler slog.Handler

	// timeNow is a function that returns the current time.
	// If timeNow is nil, [slog.Record.Time] will not be set.
	timeNow func() time.Time

	// capturePC indicates whether the logger should capture program
	// counter.
	capturePC bool

	// skipPC is the number of stack frames to skip for program counter.
	// If skipPC is negative, [slog.Record.PC] will not be set.
	skipPC int
}

// New creates a new [Logger] with the given non-nil [slog.Handler].
func New(h slog.Handler) *Logger {
	return &Logger{
		handler:   h,
		timeNow:   time.Now,
		capturePC: true,
	}
}

// Handler returns the underlying [slog.Handler].
func (l *Logger) Handler() slog.Handler {
	return l.handler
}

// WithHandler returns a [Logger] that uses the given [slog.Handler].
func (l *Logger) WithHandler(h slog.Handler) *Logger {
	c := *l
	c.handler = h
	return &c
}

// WithAttrs returns a [Logger] that includes the given attributes in each
// output operation.
func (l *Logger) WithAttrs(attrs ...slog.Attr) *Logger {
	if len(attrs) == 0 {
		return l
	}
	return l.WithHandler(l.handler.WithAttrs(attrs))
}

// WithGroup returns a [Logger] that starts a group. The keys of all attributes
// added to the [Logger] will be qualified by the given name.
//
// If name is empty, WithGroup returns the receiver.
func (l *Logger) WithGroup(name string) *Logger {
	if name == "" {
		return l
	}
	return l.WithHandler(l.handler.WithGroup(name))
}

// WithTime returns a [Logger] that uses the given time function.
// If now is nil, [slog.Record.Time] will not be set.
func (l *Logger) WithTime(now func() time.Time) *Logger {
	if now == nil && l.timeNow == nil {
		return l
	}
	c := *l
	c.timeNow = now
	return &c
}

// WithCapturePC returns a [Logger] with program counter capture toggled.
// If v is false, [slog.Record.PC] will not be set regardless of WithSkipPC.
func (l *Logger) WithCapturePC(v bool) *Logger {
	if l.capturePC == v {
		return l
	}
	c := *l
	c.capturePC = v
	return &c
}

// WithSkipPC returns a [Logger] that skips additional stack frames.
// If n is non-positive, WithSkipPC returns the receiver.
func (l *Logger) WithSkipPC(n int) *Logger {
	if n <= 0 || l.skipPC < 0 {
		return l
	}
	c := *l
	c.skipPC += n
	if c.skipPC < 0 {
		c.skipPC = -1
	}
	return &c
}

// Enabled reports whether l emits log records at the given context and level.
func (l *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.handler.Enabled(ctx, level)
}

// Log emits a log record with the given level, message, and attributes.
//
// By default, [slog.Record.PC] will be set to the caller’s program counter
// and [slog.Record.Time] will be set to the current time. Use WithCapturePC,
// WithSkipPC and WithTime methods to change this behavior.
func (l *Logger) Log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if !l.Enabled(ctx, level) {
		return
	}
	_ = l.output(ctx, level, msg, attrs)
}

func (l *Logger) output(
	ctx context.Context,
	level slog.Level,
	msg string,
	attrs []slog.Attr,
) error {
	var t time.Time
	if l.timeNow != nil {
		t = l.timeNow()
	}

	var pc uintptr
	if l.capturePC && l.skipPC >= 0 {
		// Skip 4 stack frames to log actual caller:
		//   - runtime.Callers
		//   - callerPC
		//   - this function
		//   - this function’s caller (Logger.Log or Writer.Write)
		skipPC := l.skipPC + 4
		if skipPC >= 0 {
			pc = callerPC(skipPC)
		}
	}

	r := slog.NewRecord(t, level, msg, pc)
	r.AddAttrs(attrs...)
	return l.handler.Handle(ctx, r)
}

// callerPC returns the program counter at the given stack depth.
func callerPC(depth int) uintptr {
	var pcs [1]uintptr
	_ = runtime.Callers(depth, pcs[:])
	return pcs[0]
}
