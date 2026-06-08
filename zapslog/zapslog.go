// Package zapslog provides a [zapcore.Core] implementation that forwards logs
// to [slog.Handler].
package zapslog

import (
	"context"
	"log/slog"

	"go.uber.org/zap/zapcore"
)

const (
	// loggerNameKey is the key used by the [Core] for the [zap.Entry]’s
	// LoggerName field.
	loggerNameKey = "logger_name"
	// stackKey is the key used by the [Core] for the [zap.Entry]’s
	// Stack field.
	stackKey = "stack"
)

// Core implements zapcore.Core and forwards log records to slog.Handler.
type Core struct {
	ctx     context.Context
	handler slog.Handler
	fields  []zapcore.Field
}

// New creates a [Core] backed by the provided context and [slog.Handler].
func New(ctx context.Context, handler slog.Handler) *Core {
	return &Core{
		ctx:     ctx,
		handler: handler,
	}
}

func (c *Core) appendFields(newFields []zapcore.Field) []zapcore.Field {
	allFields := newFields
	if n := len(c.fields); n != 0 {
		allFields = make([]zapcore.Field, 0, n+len(newFields))
		allFields = append(allFields, c.fields...)
		allFields = append(allFields, newFields...)
	}
	return allFields
}

// Enabled implements the [zapcore.Core] interface.
func (c *Core) Enabled(level zapcore.Level) bool {
	return c.handler.Enabled(c.ctx, convertLogLevel(level))
}

// With implements the [zapcore.Core] interface.
func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	cc := *c
	cc.fields = c.appendFields(fields)
	return &cc
}

// Check implements the [zapcore.Core] interface.
func (c *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

// Write implements the [zapcore.Core] interface.
func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	record := slog.NewRecord(
		entry.Time,
		convertLogLevel(entry.Level),
		entry.Message,
		entry.Caller.PC,
	)

	if entry.LoggerName != "" {
		record.AddAttrs(slog.String(loggerNameKey, entry.LoggerName))
	}

	record.AddAttrs(encodeFields(c.appendFields(fields))...)

	if entry.Stack != "" {
		record.AddAttrs(slog.String(stackKey, entry.Stack))
	}

	return c.handler.Handle(c.ctx, record)
}

// Sync implements the [zapcore.Core] interface.
func (c *Core) Sync() error {
	return nil
}
