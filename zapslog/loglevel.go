package zapslog

import (
	"log/slog"

	"go.uber.org/zap/zapcore"
)

const (
	// LevelDPanic is the log level for [zapcore.DPanicLevel] entries.
	LevelDPanic = slog.LevelError + 1

	// LevelPanic is the log level for [zapcore.PanicLevel] entries.
	LevelPanic = slog.LevelError + 2

	// LevelFatal is the log level for [zapcore.FatalLevel] entries.
	LevelFatal = slog.LevelError + 3

	// LevelInvalid is the log level for entries with invalid
	// [zapcore.Level] values.
	LevelInvalid = slog.LevelError + 4
)

// convertLogLevel converts a [zapcore.Level] to an [slog.Level]. Unsupported
// levels are mapped to [slog.LevelError]+i, where i is the number of levels
// above [zapcore.ErrorLevel]. Unknown levels are mapped to [slog.LevelError]+4.
func convertLogLevel(level zapcore.Level) slog.Level {
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
		return LevelDPanic
	case zapcore.PanicLevel:
		return LevelPanic
	case zapcore.FatalLevel:
		return LevelFatal
	default: // including zapcore.InvalidLevel
		return LevelInvalid
	}
}
