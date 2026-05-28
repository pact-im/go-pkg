// Package zapslog provides a zapcore.Core implementation that forwards logs to
// slog.Handler.
package zapslog

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a zap.Logger backed by the provided slog.Handler.
func New(handler slog.Handler) *zap.Logger {
	return zap.New(&zapSlogCore{handler: handler})
}

// zapSlogCore implements zapcore.Core and forwards log records to slog.Handler.
type zapSlogCore struct {
	handler slog.Handler
}

// Enabled reports whether the underlying slog handler accepts the given level.
func (c *zapSlogCore) Enabled(level zapcore.Level) bool {
	return c.handler.Enabled(context.Background(), zapCoreLevelToSlogLevel(level))
}

// fieldToAttr converts a zap field into a slog attribute.
func fieldToAttr(field zapcore.Field) slog.Attr {
	switch field.Type {
	case zapcore.StringType:
		return slog.String(field.Key, field.String)
	case zapcore.Int64Type:
		return slog.Int64(field.Key, field.Integer)
	case zapcore.Int32Type:
		return slog.Int(field.Key, int(field.Integer))
	case zapcore.Uint64Type:
		return slog.Uint64(field.Key, uint64(field.Integer))
	case zapcore.Float64Type:
		return slog.Float64(field.Key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.BoolType:
		return slog.Bool(field.Key, field.Integer == 1)
	case zapcore.TimeType:
		if field.Interface != nil {
			loc, ok := field.Interface.(*time.Location)
			if ok {
				return slog.Time(field.Key, time.Unix(0, field.Integer).In(loc))
			}
		}

		return slog.Time(field.Key, time.Unix(0, field.Integer))
	case zapcore.DurationType:
		return slog.Duration(field.Key, time.Duration(field.Integer))
	case zapcore.StringerType:
		if value, ok := field.Interface.(fmt.Stringer); ok {
			return slog.String(field.Key, value.String())
		}

		return slog.Any(field.Key, field.Interface)
	default:
		return slog.Any(field.Key, field.Interface)
	}
}

// fieldToAttrs converts zap fields into slog attributes.
func fieldToAttrs(fields []zapcore.Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, fieldToAttr(field))
	}

	return attrs
}

// With implements the [zapcore.Core] interface.
func (c *zapSlogCore) With(fields []zapcore.Field) zapcore.Core {
	handler := c.handler.WithAttrs(fieldToAttrs(fields))

	return &zapSlogCore{handler: handler}
}

// Check implements the [zapcore.Core] interface.
func (c *zapSlogCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}

	return ce
}

// Write implements the [zapcore.Core] interface.
func (c *zapSlogCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// https://pkg.go.dev/log/slog#hdr-Writing_a_handler
	record := slog.NewRecord(entry.Time, zapCoreLevelToSlogLevel(entry.Level), entry.Message, entry.Caller.PC)

	if entry.LoggerName != "" {
		record.AddAttrs(slog.String("name", entry.LoggerName))
	}

	for _, field := range fields {
		record.AddAttrs(fieldToAttr(field))
	}

	if entry.Stack != "" {
		record.AddAttrs(slog.String("stack", entry.Stack))
	}

	err := c.handler.Handle(context.Background(), record)
	if err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	return nil
}

// Sync implements the [zapcore.Core] interface.
func (c *zapSlogCore) Sync() error {
	return nil
}

// zapCoreLevelToSlogLevel converts a zapcore.Level to a slog.Level.
// Unsupported levels are converted to slog.LevelDebug.
func zapCoreLevelToSlogLevel(level zapcore.Level) slog.Level {
	switch level {
	case zapcore.DebugLevel:
		return slog.LevelDebug
	case zapcore.InfoLevel:
		return slog.LevelInfo
	case zapcore.WarnLevel:
		return slog.LevelWarn
	case zapcore.ErrorLevel:
		return slog.LevelError
	case zapcore.DPanicLevel:
		return slog.LevelError
	case zapcore.PanicLevel:
		return slog.LevelError
	case zapcore.FatalLevel:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
