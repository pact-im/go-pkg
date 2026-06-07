package zapslog

import (
	"context"
	"log/slog"
	"slices"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCore(t *testing.T) {
	tests := []struct {
		name  string
		limit slog.Leveler
		run   func(*zap.Logger) error
		want  []slog.Record
	}{
		{
			name: "sync",
			run:  (*zap.Logger).Sync,
			want: nil,
		},
		{
			name: "logger name",
			run: func(l *zap.Logger) error {
				l.Named("zapslog_test").Info("msg")
				return nil
			},
			want: []slog.Record{
				slogRecordAddAttrs(
					slog.Record{Message: "msg", Level: slog.LevelInfo},
					slog.String("logger_name", "zapslog_test"),
				),
			},
		},
		{
			name:  "log levels",
			limit: slog.LevelInfo,
			run: func(l *zap.Logger) error {
				l.Debug("debug level")
				l.Info("info level")
				l.Warn("warn level")
				l.Error("error level")
				return nil
			},
			want: []slog.Record{
				{Message: "info level", Level: slog.LevelInfo},
				{Message: "warn level", Level: slog.LevelWarn},
				{Message: "error level", Level: slog.LevelError},
			},
		},
		{
			name: "with fields",
			run: func(l *zap.Logger) error {
				l.With(zap.Namespace("foo")).
					With(zap.String("key", "hello")).
					Info("msg", zap.Int("int", 232))
				return nil
			},
			want: []slog.Record{
				slogRecordAddAttrs(
					slog.Record{Message: "msg", Level: slog.LevelInfo},
					slog.GroupAttrs("foo",
						slog.String("key", "hello"),
						slog.Int("int", 232),
					),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var records []slog.Record
			handler := handlerFunc(func(_ context.Context, r slog.Record) error {
				records = append(records, r)
				return nil
			})
			core := New(
				context.Background(),
				&limitedHandler{
					handler: handler,
					level:   tt.limit,
				},
			)
			logger := zap.New(
				core,
				zap.WithClock(fakeClock{}),
				zap.ErrorOutput(zapcore.AddSync(t.Output())),
			)

			if err := tt.run(logger); err != nil {
				t.Fatal(err)
			}
			if !slices.EqualFunc(records, tt.want, slogRecordEqual) {
				t.Fatal("unexpected records")
			}
		})
	}
}

type limitedHandler struct {
	handler slog.Handler
	level   slog.Leveler
}

func (h *limitedHandler) Enabled(_ context.Context, l slog.Level) bool {
	return h.level == nil || l >= h.level.Level()
}

func (h *limitedHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *limitedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hc := *h
	hc.handler = h.handler.WithAttrs(attrs)
	return &hc
}

func (h *limitedHandler) WithGroup(group string) slog.Handler {
	hc := *h
	hc.handler = h.handler.WithGroup(group)
	return &hc
}

// handlerFunc is a function that implements the [slog.Handler] interface.
type handlerFunc func(context.Context, slog.Record) error

func (f handlerFunc) Enabled(_ context.Context, _ slog.Level) bool    { return true }
func (f handlerFunc) Handle(ctx context.Context, r slog.Record) error { return f(ctx, r) }
func (f handlerFunc) WithAttrs(_ []slog.Attr) slog.Handler            { return f }
func (f handlerFunc) WithGroup(_ string) slog.Handler                 { return f }

// fakeClock is a fake [zapcore.Clock] implementation that returns zero value of
// [time.Time]. NewTicker method panics, as it is not used in the test suite.
type fakeClock struct{}

func (fakeClock) Now() time.Time                         { return time.Time{} }
func (fakeClock) NewTicker(_ time.Duration) *time.Ticker { panic("not implemented") }

// slogRecordAddAttrs is a convenience function that adds attributes to an
// [slog.Record] and returns the updated record.
func slogRecordAddAttrs(r slog.Record, attrs ...slog.Attr) slog.Record {
	r.AddAttrs(attrs...)
	return r
}

// slogRecordEqual returns true if two [slog.Record] values are equal, false
// otherwise.
func slogRecordEqual(a, b slog.Record) bool {
	return a.Time.Equal(b.Time) &&
		a.Message == b.Message &&
		a.Level == b.Level &&
		a.PC == b.PC &&
		slices.EqualFunc(
			slices.Collect(a.Attrs),
			slices.Collect(b.Attrs),
			slogAttrEqual,
		)
}
