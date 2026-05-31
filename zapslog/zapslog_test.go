package zapslog

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		requireNoError(t, err)

		got := b.String()
		assertContains(t, got, `level=DEBUG msg="debug level"`)
		assertContains(t, got, `level=INFO msg="info level"`)
		assertContains(t, got, `level=WARN msg="warn level"`)
		assertContains(t, got, `level=ERROR msg="error level"`)
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
		requireNoError(t, err)

		got := b.String()
		assertContains(t, got, `msg=fields`)
		assertContains(t, got, `name=example`)
		assertContains(t, got, `base_key=base_value`)
		assertContains(t, got, `base_count=7`)
		assertContains(t, got, `key=value`)
		assertContains(t, got, `count=42`)
		assertContains(t, got, `flag=true`)
		assertContains(t, got, `timeout=2s`)
		assertContains(t, got, `created_at=2026-05-27T12:34:56.000Z`)
		assertContains(t, got, `stringer=stringer-value`)
		assertContains(t, got, `error=boom`)
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
		requireNoError(t, err)

		got := b.String()
		assertNotContains(t, got, `skip debug`)
		assertContains(t, got, `keep info`)
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
	requireNoError(b, err)

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

func requireNoError(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatalf("unexpected error: %v", err)
	}
}

func assertContains(tb testing.TB, s, substr string) {
	tb.Helper()

	if !strings.Contains(s, substr) {
		tb.Fatalf("expected %q to contain %q", s, substr)
	}
}

func assertNotContains(tb testing.TB, s, substr string) {
	tb.Helper()

	if strings.Contains(s, substr) {
		tb.Fatalf("expected %q not to contain %q", s, substr)
	}
}
