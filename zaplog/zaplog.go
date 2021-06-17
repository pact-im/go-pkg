// Package zaplog provides a constructor for zap.Logger with sensible defaults.
package zaplog

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a new zap.Logger that writes to w.
func New(w io.Writer) *zap.Logger {
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		}),
		zapcore.AddSync(w),
		zap.DebugLevel,
	), zap.WithCaller(true))
}

// Tee returns log’s clone that duplicates log entries into another core.
func Tee(log *zap.Logger, core zapcore.Core) *zap.Logger {
	return log.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, core)
	}))
}
