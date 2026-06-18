package logs

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		input string
		want  []slog.Record
	}{
		{
			name:  "plain message",
			level: slog.LevelInfo,
			input: "hello",
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "hello"},
			},
		},
		{
			name:  "trailing newline stripped",
			level: slog.LevelInfo,
			input: "hello\n",
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "hello"},
			},
		},
		{
			name:  "warn level",
			level: slog.LevelWarn,
			input: "warning\n",
			want: []slog.Record{
				{Level: slog.LevelWarn, Message: "warning"},
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
			l := New(handler).WithTime(nil).WithCapturePC(false)
			w := l.Writer(context.Background(), tt.level)

			n, err := w.Write([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n != len(tt.input) {
				t.Fatalf("wrote %d bytes, want %d", n, len(tt.input))
			}
			if !slices.EqualFunc(records, tt.want, slogRecordEqual) {
				t.Fatalf("unexpected records")
			}
		})
	}
}

func TestWriter_Write_Disabled(t *testing.T) {
	handler := slog.DiscardHandler
	l := New(handler).WithTime(nil).WithCapturePC(false)
	w := l.Writer(context.Background(), slog.LevelInfo)

	buf := []byte("hello\n")
	n, err := w.Write(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("wrote %d bytes, want %d", n, len(buf))
	}
}

func TestNewLogLogger(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []slog.Record
	}{
		{
			name:  "print",
			input: "hello",
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "hello"},
			},
		},
		{
			name:  "println",
			input: "world",
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "world"},
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
			l := New(handler).WithTime(nil).WithCapturePC(false)
			ll := NewLogLogger(context.Background(), slog.LevelInfo, l)

			ll.Println(tt.input)

			if !slices.EqualFunc(records, tt.want, slogRecordEqual) {
				t.Fatalf("unexpected records")
			}
		})
	}
}

func TestWriter_Write_Error(t *testing.T) {
	wantErr := errors.New("handle failed")
	handler := HandlerFunc(func(_ context.Context, _ slog.Record) error {
		return wantErr
	})
	l := New(handler).WithTime(nil).WithCapturePC(false)
	w := l.Writer(context.Background(), slog.LevelInfo)

	_, err := w.Write([]byte("hello"))
	if err != wantErr {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
}

func TestWriter_Write_Multiple(t *testing.T) {
	var records []slog.Record
	handler := HandlerFunc(func(_ context.Context, r slog.Record) error {
		records = append(records, r)
		return nil
	})
	l := New(handler).WithTime(nil).WithCapturePC(false)
	w := l.Writer(context.Background(), slog.LevelInfo)

	_, _ = w.Write([]byte("first\n"))
	_, _ = w.Write([]byte("second\n"))

	want := []slog.Record{
		{Level: slog.LevelInfo, Message: "first"},
		{Level: slog.LevelInfo, Message: "second"},
	}
	if !slices.EqualFunc(records, want, slogRecordEqual) {
		t.Fatalf("unexpected records")
	}
}
