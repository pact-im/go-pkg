package zapslog

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testTime := time.Date(2026, time.May, 27, 12, 34, 56, 0, time.UTC)

	t.Run("levels", func(t *testing.T) {
		var b bytes.Buffer
		handler := slog.NewTextHandler(&b,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)

		loggerZap := New(handler)

		loggerZap.Debug("debug level")
		loggerZap.Info("info level")
		loggerZap.Warn("warn level")
		loggerZap.Error("error level")

		err := loggerZap.Sync()
		require.NoError(t, err)

		got := b.String()
		assert.Contains(t, got, `level=DEBUG msg="debug level"`)
		assert.Contains(t, got, `level=INFO msg="info level"`)
		assert.Contains(t, got, `level=WARN msg="warn level"`)
		assert.Contains(t, got, `level=ERROR msg="error level"`)
	})

	t.Run("named logger and fields", func(t *testing.T) {
		var b bytes.Buffer
		handler := slog.NewTextHandler(&b,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)

		loggerZap := New(handler).WithOptions(zap.WithCaller(false))

		loggerZap = loggerZap.Named("example").With(
			zap.String("base_key", "base_value"),
			zap.Int("base_count", 7),
		)

		loggerZap.Info("fields",
			zap.String("key", "value"),
			zap.Int("count", 42),
			zap.Bool("flag", true),
			zap.Duration("timeout", 2*time.Second),
			zap.Time("created_at", testTime),
			zap.Stringer("stringer", testStringer("stringer-value")),
			zap.Error(errors.New("boom")),
		)

		err := loggerZap.Sync()
		require.NoError(t, err)

		got := b.String()
		assert.Contains(t, got, `msg=fields`)
		assert.Contains(t, got, `name=example`)
		assert.Contains(t, got, `base_key=base_value`)
		assert.Contains(t, got, `base_count=7`)
		assert.Contains(t, got, `key=value`)
		assert.Contains(t, got, `count=42`)
		assert.Contains(t, got, `flag=true`)
		assert.Contains(t, got, `timeout=2s`)
		assert.Contains(t, got, `created_at=2026-05-27T12:34:56.000Z`)
		assert.Contains(t, got, `stringer=stringer-value`)
		assert.Contains(t, got, `error=boom`)
	})

	t.Run("level filtering", func(t *testing.T) {
		var b bytes.Buffer
		handler := slog.NewTextHandler(&b,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		)

		loggerZap := New(handler).WithOptions(zap.WithCaller(false))

		loggerZap.Debug("skip debug")
		loggerZap.Info("keep info")

		err := loggerZap.Sync()
		require.NoError(t, err)

		got := b.String()
		assert.NotContains(t, got, `skip debug`)
		assert.Contains(t, got, `keep info`)
	})
}

func BenchmarkNew(b *testing.B) {
	loggerZap := New(noopSlogHandler{})

	b.ResetTimer()

	for b.Loop() {
		loggerZap.Info("hello world")
	}
}

func BenchmarkNewNoCaller(b *testing.B) {
	loggerZap := New(noopSlogHandler{}).WithOptions(zap.WithCaller(false))

	b.ResetTimer()

	for b.Loop() {
		loggerZap.Info("hello world")
	}
}

func BenchmarkNewFields(b *testing.B) {
	loggerZap := New(noopSlogHandler{}).WithOptions(zap.WithCaller(false))

	b.ResetTimer()

	for b.Loop() {
		loggerZap.Info("hello world",
			zap.String("key", "value"),
			zap.Int("count", 42),
			zap.Bool("flag", true),
		)
	}
}

func BenchmarkZap(b *testing.B) {
	loggerZap, err := zap.NewProduction(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewNopCore()
	}))
	require.NoError(b, err)

	b.ResetTimer()

	for b.Loop() {
		loggerZap.Info("hello world")
	}
}

func BenchmarkSlog(b *testing.B) {
	loggerSlog := slog.New(noopSlogHandler{})

	b.ResetTimer()

	for b.Loop() {
		loggerSlog.Info("hello world")
	}
}

type noopSlogHandler struct{}

func (noopSlogHandler) Enabled(context.Context, slog.Level) bool  { return true }
func (noopSlogHandler) Handle(context.Context, slog.Record) error { return nil }
func (h noopSlogHandler) WithAttrs([]slog.Attr) slog.Handler      { return h }
func (h noopSlogHandler) WithGroup(string) slog.Handler           { return h }

type testStringer string

func (s testStringer) String() string {
	return string(s)
}
