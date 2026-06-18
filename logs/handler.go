package logs

import (
	"context"
	"log/slog"
	"slices"
	"sync/atomic"
)

// HandlerFunc is a function that implements the [slog.Handler] interface.
type HandlerFunc func(context.Context, slog.Record) error

// Enabled reports whether the handler handles records at the given level.
// It always returns true.
func (f HandlerFunc) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle calls f to handle the log record.
func (f HandlerFunc) Handle(ctx context.Context, r slog.Record) error {
	return f(ctx, r)
}

// WithAttrs returns a new handler that merges the given attributes with
// any attributes from the record when handling.
func (f HandlerFunc) WithAttrs(attrs []slog.Attr) slog.Handler {
	return HandlerFunc(func(ctx context.Context, r slog.Record) error {
		return f(ctx, slogRecordAddAttrs(
			slog.NewRecord(r.Time, r.Level, r.Message, r.PC),
			slices.AppendSeq(slices.Clip(attrs), r.Attrs)...,
		))
	})
}

// WithGroup returns a new handler that wraps the record’s attributes
// in a group with the given name.
func (f HandlerFunc) WithGroup(name string) slog.Handler {
	return HandlerFunc(func(ctx context.Context, r slog.Record) error {
		return f(ctx, slogRecordAddAttrs(
			slog.NewRecord(r.Time, r.Level, r.Message, r.PC),
			slog.GroupAttrs(name, slices.Collect(r.Attrs)...),
		))
	})
}

// slogRecordAddAttrs is a convenience function that adds attributes to an
// [slog.Record] and returns the updated record.
func slogRecordAddAttrs(r slog.Record, attrs ...slog.Attr) slog.Record {
	r.AddAttrs(attrs...)
	return r
}

// Hooks contains hooks for [WrapHandler] that wrap [slog.Handler] methods.
type Hooks struct {
	// Enabled wraps the handler’s Enabled method to determines if a log
	// level is enabled.
	Enabled func(ctx context.Context, l slog.Level, next slog.Handler) bool
	// Handle wraps the handler’s Handle method to intercept records before
	// they reach the next handler.
	Handle func(ctx context.Context, r slog.Record, next slog.Handler) error
}

// WrapHandler is an [slog.Handler] that uses [Hooks] to modify the behavior of
// an underlying handler.
type WrapHandler struct {
	next  slog.Handler
	hooks Hooks
}

// Wrap creates an [slog.Handler] that wraps a handler with specific hooks.
func Wrap(next slog.Handler, hooks Hooks) *WrapHandler {
	return &WrapHandler{
		next:  next,
		hooks: hooks,
	}
}

// Enabled implements the [slog.Handler] interface. It calls [Hooks.Enabled] if
// provided, otherwise it uses the underlying handler.
func (h *WrapHandler) Enabled(ctx context.Context, l slog.Level) bool {
	if h.hooks.Enabled == nil {
		return h.next.Enabled(ctx, l)
	}
	return h.hooks.Enabled(ctx, l, h.next)
}

// Handle implements the [slog.Handler] interface. It calls [Hooks.Handle] if
// provided, otherwise it uses the underlying handler.
func (h *WrapHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.hooks.Handle == nil {
		return h.next.Handle(ctx, r)
	}
	return h.hooks.Handle(ctx, r, h.next)
}

// WithAttrs implements the [slog.Handler] interface.
func (h *WrapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hc := *h
	hc.next = h.next.WithAttrs(attrs)
	return &hc
}

// WithGroup implements the [slog.Handler] interface.
func (h *WrapHandler) WithGroup(name string) slog.Handler {
	hc := *h
	hc.next = h.next.WithGroup(name)
	return &hc
}

// ExpiredContextFilter drops log records once the context expires to reduce
// noise from [context.Canceled] and similar errors.
//
// It suppresses error logs that occur after a request or an operation has
// already been aborted, e.g. when HTTP client cancels the request.
type ExpiredContextFilter struct {
	dropped atomic.Uint64
}

// Hooks returns the filter’s hooks.
func (f *ExpiredContextFilter) Hooks() Hooks {
	return Hooks{
		Enabled: f.Enabled,
		Handle:  f.Handle,
	}
}

// Enabled returns false if the context has expired, otherwise it delegates to
// the next handler.
func (f *ExpiredContextFilter) Enabled(ctx context.Context, l slog.Level, next slog.Handler) bool {
	return next.Enabled(ctx, l) && !f.expired(ctx)
}

// Handle drops log record if the context has expired, otherwise it delegates to
// the next handler.
func (f *ExpiredContextFilter) Handle(ctx context.Context, r slog.Record, next slog.Handler) error {
	if f.expired(ctx) {
		return nil
	}
	return next.Handle(ctx, r)
}

// Dropped returns the number of log records that have been dropped due to
// expired contexts. Note that the counter wraps around on overflow.
func (f *ExpiredContextFilter) Dropped() uint64 {
	return f.dropped.Load()
}

// expired returns true if the context is canceled or timed out.
func (f *ExpiredContextFilter) expired(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		f.dropped.Add(1)
		return true
	default:
		return false
	}
}

// LevelLimitFilter drops log records with a level less than the specified
// [slog.Leveler].
type LevelLimitFilter struct {
	Level slog.Leveler
}

// Limit returns a new [slog.Handler] that wraps the provided handler h
// and only allows records with a level greater than or equal to l.
func Limit(h slog.Handler, l slog.Leveler) slog.Handler {
	return Wrap(h, (&LevelLimitFilter{l}).Hooks())
}

// Hooks returns the filter’s hooks.
func (f *LevelLimitFilter) Hooks() Hooks {
	return Hooks{
		Enabled: f.Enabled,
		Handle:  f.Handle,
	}
}

// Enabled returns false if the log level is below the limit, otherwise it
// delegates to the next handler.
func (f *LevelLimitFilter) Enabled(ctx context.Context, l slog.Level, next slog.Handler) bool {
	return f.enabled(l) && next.Enabled(ctx, l)
}

// Handle drops log record if the log level is below the limit, otherwise it
// delegates to the next handler.
func (f *LevelLimitFilter) Handle(ctx context.Context, r slog.Record, next slog.Handler) error {
	if !f.enabled(r.Level) {
		return nil
	}
	return next.Handle(ctx, r)
}

// enabled returns false if the log level is below the limit.
func (f *LevelLimitFilter) enabled(l slog.Level) bool {
	return f.Level == nil || l >= f.Level.Level()
}
