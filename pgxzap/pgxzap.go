// Package pgxzap implements [tracelog.Logger] adapter for [zap.Logger].
package pgxzap

import (
	"cmp"
	"context"
	"slices"

	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ tracelog.Logger = (*Logger)(nil)

// Logger is an adapter from [zap.Logger] to [tracelog.Logger].
type Logger struct {
	logger *zap.Logger
}

// NewLogger returns a new [Logger] instance.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger.WithOptions(zap.AddCallerSkip(3)),
	}
}

// Log implements the [tracelog.Logger] interface.
func (l *Logger) Log(
	_ context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	zapLevel := translateTraceLogLevel(level)
	if !l.logger.Core().Enabled(zapLevel) {
		return
	}

	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	slices.SortFunc(fields, func(a, b zap.Field) int {
		return cmp.Compare(a.Key, b.Key)
	})

	l.logger.Log(zapLevel, msg, fields...)
}

func translateTraceLogLevel(level tracelog.LogLevel) zapcore.Level {
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		return zapcore.DebugLevel
	case tracelog.LogLevelInfo:
		return zapcore.InfoLevel
	case tracelog.LogLevelWarn:
		return zapcore.WarnLevel
	case tracelog.LogLevelError:
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}
