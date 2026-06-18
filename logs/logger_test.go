package logs

import (
	"context"
	"log/slog"
	"slices"
	"testing"
	"time"
)

// slogRecordEqual returns true if two [slog.Record] values are equal.
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

// slogAttrEqual returns true if two [slog.Attr] values are equal.
func slogAttrEqual(a, b slog.Attr) bool {
	return a.Key == b.Key && slogValueEqual(
		a.Value.Resolve(),
		b.Value.Resolve(),
	)
}

// slogValueEqual returns true if two [slog.Value] values are equal.
// For group values, it compares the inner attributes recursively.
func slogValueEqual(a, b slog.Value) bool {
	if a.Kind() == slog.KindGroup && b.Kind() == slog.KindGroup {
		return slices.EqualFunc(a.Group(), b.Group(), slogAttrEqual)
	}
	return a.Equal(b)
}

func TestNew(t *testing.T) {
	h := slog.DiscardHandler
	l := New(h)

	if l.Handler() != h {
		t.Fatal("expected handler to be set")
	}
	if l.timeNow == nil {
		t.Fatal("expected timeNow to be set")
	}
	if !l.capturePC {
		t.Fatal("expected capturePC to be true")
	}
}

func TestLogger_Enabled(t *testing.T) {
	tests := []struct {
		name    string
		handler slog.Handler
		level   slog.Level
		want    bool
	}{
		{
			name:    "discard always disabled",
			handler: slog.DiscardHandler,
			level:   slog.LevelInfo,
			want:    false,
		},
		{
			name:    "info enabled",
			handler: HandlerFunc(nil),
			level:   slog.LevelInfo,
			want:    true,
		},
		{
			name:    "debug enabled",
			handler: HandlerFunc(nil),
			level:   slog.LevelDebug,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.handler)
			got := l.Enabled(context.Background(), tt.level)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogger_Log(t *testing.T) {
	tests := []struct {
		name string
		run  func(*Logger)
		want []slog.Record
	}{
		{
			name: "debug",
			run: func(l *Logger) {
				l.Log(context.Background(), slog.LevelDebug, "debug msg")
			},
			want: []slog.Record{
				{Level: slog.LevelDebug, Message: "debug msg"},
			},
		},
		{
			name: "attrs",
			run: func(l *Logger) {
				l.Log(context.Background(), slog.LevelWarn, "warn", slog.String("k", "v"))
			},
			want: []slog.Record{
				slogRecordAddAttrs(
					slog.Record{Level: slog.LevelWarn, Message: "warn"},
					slog.String("k", "v"),
				),
			},
		},
		{
			name: "with empty attrs",
			run: func(l *Logger) {
				l.WithAttrs().Log(context.Background(), slog.LevelInfo, "hello")
			},
			want: []slog.Record{
				{Message: "hello", Level: slog.LevelInfo},
			},
		},
		{
			name: "with single attr",
			run: func(l *Logger) {
				l.WithAttrs(slog.String("key", "value")).
					Log(context.Background(), slog.LevelInfo, "hello")
			},
			want: []slog.Record{
				slogRecordAddAttrs(
					slog.Record{Message: "hello", Level: slog.LevelInfo},
					slog.String("key", "value"),
				),
			},
		},
		{
			name: "with multiple attrs",
			run: func(l *Logger) {
				l.WithAttrs(slog.String("a", "1"), slog.Int("b", 2)).
					Log(context.Background(), slog.LevelInfo, "hello")
			},
			want: []slog.Record{
				slogRecordAddAttrs(
					slog.Record{Message: "hello", Level: slog.LevelInfo},
					slog.String("a", "1"),
					slog.Int("b", 2),
				),
			},
		},
		{
			name: "with empty group",
			run: func(l *Logger) {
				l.WithGroup("").
					Log(context.Background(), slog.LevelInfo, "hello", slog.String("key", "value"))
			},
			want: []slog.Record{
				slogRecordAddAttrs(
					slog.Record{Message: "hello", Level: slog.LevelInfo},
					slog.String("key", "value"),
				),
			},
		},
		{
			name: "with named group",
			run: func(l *Logger) {
				l.WithGroup("ns").
					Log(context.Background(), slog.LevelInfo, "hello", slog.String("key", "value"))
			},
			want: []slog.Record{
				slogRecordAddAttrs(
					slog.Record{Message: "hello", Level: slog.LevelInfo},
					slog.GroupAttrs("ns", slog.String("key", "value")),
				),
			},
		},
		{
			name: "with custom time",
			run: func(l *Logger) {
				t := time.Date(2026, time.June, 18, 12, 34, 56, 789, time.UTC)
				l.WithTime(func() time.Time { return t }).
					Log(context.Background(), slog.LevelInfo, "hello")
			},
			want: []slog.Record{
				{
					Time:    time.Date(2026, time.June, 18, 12, 34, 56, 789, time.UTC),
					Level:   slog.LevelInfo,
					Message: "hello",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var records []slog.Record
			handler := HandlerFunc(func(_ context.Context, r slog.Record) error {
				records = append(records, r)
				return nil
			})
			l := New(handler).
				WithTime(nil).
				WithCapturePC(false)

			tt.run(l)

			if !slices.EqualFunc(records, tt.want, slogRecordEqual) {
				t.Fatalf("unexpected records:\n got: %#v\nwant: %#v", records, tt.want)
			}
		})
	}
}

func TestLogger_WithCapturePC(t *testing.T) {
	tests := []struct {
		name     string
		expectPC bool
	}{
		{name: "true", expectPC: true},
		{name: "false", expectPC: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var records []slog.Record
			handler := HandlerFunc(func(_ context.Context, r slog.Record) error {
				records = append(records, r)
				return nil
			})
			l := New(handler).
				WithTime(nil).
				WithCapturePC(tt.expectPC).
				WithSkipPC(1)

			l.Log(context.Background(), slog.LevelInfo, "hello")

			if len(records) != 1 {
				t.Fatalf("got %d records, want 1", len(records))
			}

			var want uintptr
			if tt.expectPC {
				// skip [runtime.Callers, callerPC, test func]
				want = callerPC(3)
			}
			if got := records[0].PC; got != want {
				t.Fatalf("unexpected PC:\n got: %#x\nwant: %#x", got, want)
			}
		})
	}
}
