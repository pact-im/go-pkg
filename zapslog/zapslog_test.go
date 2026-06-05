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
)

type handler struct {
	level   slog.Level
	records []slog.Record
}

func (h *handler) Enabled(_ context.Context, lev slog.Level) bool {
	return lev >= h.level
}

func (h *handler) Handle(_ context.Context, rec slog.Record) error {
	h.records = append(h.records, rec)
	return nil
}

func (h *handler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *handler) WithGroup(string) slog.Handler {
	return h
}

func TestNew(t *testing.T) {
	t.Run("levels", func(t *testing.T) {
		handler := &handler{
			level: slog.LevelDebug,
		}

		loggerZap := zap.New(New(context.Background(), handler))

		loggerZap.Debug("debug level")
		loggerZap.Info("info level")
		loggerZap.Warn("warn level")
		loggerZap.Error("error level")

		err := loggerZap.Sync()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(handler.records) != 4 {
			t.Fatalf("expected 4 records, got %d", len(handler.records))
		}

		if handler.records[0].Level != slog.LevelDebug || handler.records[0].Message != "debug level" {
			t.Fatalf("unexpected first record: level=%v message=%q", handler.records[0].Level, handler.records[0].Message)
		}
		if handler.records[1].Level != slog.LevelInfo || handler.records[1].Message != "info level" {
			t.Fatalf("unexpected second record: level=%v message=%q", handler.records[1].Level, handler.records[1].Message)
		}
		if handler.records[2].Level != slog.LevelWarn || handler.records[2].Message != "warn level" {
			t.Fatalf("unexpected third record: level=%v message=%q", handler.records[2].Level, handler.records[2].Message)
		}
		if handler.records[3].Level != slog.LevelError || handler.records[3].Message != "error level" {
			t.Fatalf("unexpected fourth record: level=%v message=%q", handler.records[3].Level, handler.records[3].Message)
		}
	})

	t.Run("level filtering", func(t *testing.T) {
		handler := &handler{
			level: slog.LevelInfo,
		}

		loggerZap := zap.New(New(context.Background(), handler)).WithOptions(zap.WithCaller(false))

		loggerZap.Debug("skip debug")
		loggerZap.Info("keep info")

		err := loggerZap.Sync()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(handler.records) != 1 {
			t.Fatalf("expected 1 record, got %d", len(handler.records))
		}

		if handler.records[0].Level != slog.LevelInfo || handler.records[0].Message != "keep info" {
			t.Fatalf("unexpected record: level=%v message=%q", handler.records[0].Level, handler.records[0].Message)
		}
	})

	t.Run("with namespace", func(t *testing.T) {
		handler := &handler{
			level: slog.LevelDebug,
		}

		loggerZap := zap.New(New(context.Background(), handler)).WithOptions(zap.WithCaller(false))

		newCore := loggerZap.With(zap.Namespace("foo"))
		newCore = newCore.With(zap.String("key", "hello"))
		newCore.Info("msg", zap.Int("int", 232))

		if err := newCore.Sync(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(handler.records) != 1 {
			t.Fatalf("expected 1 record, got %d", len(handler.records))
		}

		record := handler.records[0]
		if record.Level != slog.LevelInfo || record.Message != "msg" {
			t.Fatalf("unexpected record: level=%v message=%q", record.Level, record.Message)
		}

		var attrs []slog.Attr
		record.Attrs(func(attr slog.Attr) bool {
			attrs = append(attrs, attr)
			return true
		})

		if len(attrs) != 1 {
			t.Fatalf("expected 1 slog.Attr, got %d", len(attrs))
		}

		want := slog.GroupAttrs("foo", slog.String("key", "hello"), slog.Int("int", 232))
		if got := attrs[0]; !slogAttrEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
}

func TestEncodeFields(t *testing.T) {
	testTime := time.Date(2026, time.May, 27, 12, 34, 56, 0, time.UTC)
	fullTime := testTime.Add(123 * time.Nanosecond)
	tests := []struct {
		name         string
		fields       []zap.Field
		wantContains []slog.Attr
	}{
		{
			name:         "string",
			fields:       []zap.Field{zap.String("key", "value")},
			wantContains: []slog.Attr{slog.String("key", "value")},
		},
		{
			name:         "int",
			fields:       []zap.Field{zap.Int("count", 42)},
			wantContains: []slog.Attr{slog.Int("count", 42)},
		},
		{
			name:         "float64",
			fields:       []zap.Field{zap.Float64("float64", 123.456789)},
			wantContains: []slog.Attr{slog.Float64("float64", 123.456789)},
		},
		{
			name:         "float32",
			fields:       []zap.Field{zap.Float32("float32", 12.5)},
			wantContains: []slog.Attr{slog.Float64("float32", 12.5)},
		},
		{
			name:         "bool",
			fields:       []zap.Field{zap.Bool("flag", true)},
			wantContains: []slog.Attr{slog.Bool("flag", true)},
		},
		{
			name:         "duration",
			fields:       []zap.Field{zap.Duration("timeout", 2*time.Second)},
			wantContains: []slog.Attr{slog.Duration("timeout", 2*time.Second)},
		},
		{
			name:         "time",
			fields:       []zap.Field{zap.Time("created_at", testTime)},
			wantContains: []slog.Attr{slog.Time("created_at", testTime)},
		},
		{
			name:         "int32",
			fields:       []zap.Field{zap.Int32("int32", 32)},
			wantContains: []slog.Attr{slog.Int64("int32", 32)},
		},
		{
			name:         "int16",
			fields:       []zap.Field{zap.Int16("int16", 16)},
			wantContains: []slog.Attr{slog.Int64("int16", 16)},
		},
		{
			name:         "int8",
			fields:       []zap.Field{zap.Int8("int8", 8)},
			wantContains: []slog.Attr{slog.Int64("int8", 8)},
		},
		{
			name:         "uint64",
			fields:       []zap.Field{zap.Uint64("uint64", 64)},
			wantContains: []slog.Attr{slog.Uint64("uint64", 64)},
		},
		{
			name:         "uint32",
			fields:       []zap.Field{zap.Uint32("uint32", 32)},
			wantContains: []slog.Attr{slog.Uint64("uint32", 32)},
		},
		{
			name:         "uint16",
			fields:       []zap.Field{zap.Uint16("uint16", 16)},
			wantContains: []slog.Attr{slog.Uint64("uint16", 16)},
		},
		{
			name:         "uint8",
			fields:       []zap.Field{zap.Uint8("uint8", 8)},
			wantContains: []slog.Attr{slog.Uint64("uint8", 8)},
		},
		{
			name:         "uintptr",
			fields:       []zap.Field{zap.Uintptr("pointer", uintptr(123))},
			wantContains: []slog.Attr{slog.Uint64("pointer", 123)},
		},
		{
			name:         "time full",
			fields:       []zap.Field{{Key: "full_time", Type: zapcore.TimeFullType, Interface: fullTime}},
			wantContains: []slog.Attr{slog.Time("full_time", fullTime)},
		},
		{
			name:         "stringer",
			fields:       []zap.Field{zap.Stringer("stringer", testStringer("stringer-value"))},
			wantContains: []slog.Attr{slog.String("stringer", "stringer-value")},
		},
		{
			name:         "error",
			fields:       []zap.Field{zap.Error(errors.New("boom"))},
			wantContains: []slog.Attr{slog.String("error", "boom")},
		},
		{
			name:         "binary",
			fields:       []zap.Field{zap.Binary("binary", []byte("abc"))},
			wantContains: []slog.Attr{slog.Any("binary", []byte("abc"))},
		},
		{
			name:         "byte string",
			fields:       []zap.Field{zap.ByteString("byte_string", []byte("hello"))},
			wantContains: []slog.Attr{slog.String("byte_string", "hello")},
		},
		{
			name:         "complex64",
			fields:       []zap.Field{zap.Complex64("complex64", complex64(1+2i))},
			wantContains: []slog.Attr{slog.Any("complex64", complex64(1+2i))},
		},
		{
			name:         "complex128",
			fields:       []zap.Field{zap.Complex128("complex128", complex128(3+4i))},
			wantContains: []slog.Attr{slog.Any("complex128", complex128(3+4i))},
		},
		{
			name:         "reflect",
			fields:       []zap.Field{zap.Reflect("reflect", []byte("hello"))},
			wantContains: []slog.Attr{slog.Any("reflect", []byte("hello"))},
		},
		{
			name:   "array",
			fields: []zap.Field{zap.Array("array", testArray{1, "two"})},
			wantContains: []slog.Attr{
				slog.GroupAttrs("array",
					slog.Int64("0", 1),
					slog.String("1", "two"),
				),
			},
		},
		{
			name: "array objects",
			fields: []zap.Field{zap.Array("array_objects", objectsArray{
				{id: 1, name: "value1"},
				{id: 2, name: "value2"},
			})},
			wantContains: []slog.Attr{
				slog.GroupAttrs("array_objects",
					slog.GroupAttrs("0",
						slog.Int("id", 1),
						slog.String("name", "value1"),
					),
					slog.GroupAttrs("1",
						slog.Int("id", 2),
						slog.String("name", "value2"),
					),
				),
			},
		},
		{
			name: "array nested objects",
			fields: []zap.Field{zap.Array("array_nested_objects", objectsArrayNested{
				{name: "value1", nested: "nested1"},
				{name: "value2", nested: "nested2"},
			})},
			wantContains: []slog.Attr{
				slog.GroupAttrs("array_nested_objects",
					slog.GroupAttrs("0",
						slog.String("name", "value1"),
						slog.GroupAttrs("ns", slog.String("nested", "nested1")),
					),
					slog.GroupAttrs("1",
						slog.String("name", "value2"),
						slog.GroupAttrs("ns", slog.String("nested", "nested2")),
					),
				),
			},
		},
		{
			name:   "object",
			fields: []zap.Field{zap.Object("object", testObject{id: 1, name: "obj"})},
			wantContains: []slog.Attr{
				slog.GroupAttrs("object",
					slog.Int("id", 1),
					slog.String("name", "obj"),
				),
			},
		},
		{
			name:   "inline",
			fields: []zap.Field{zap.Inline(testObject{id: 7, name: "inline"})},
			wantContains: []slog.Attr{
				slog.Int("id", 7),
				slog.String("name", "inline"),
			},
		},
		{
			name: "namespace",
			fields: []zap.Field{
				zap.Namespace("ns"),
				zap.String("nested", "value"),
				zap.Object("deep", testObject{id: 2, name: "nested-obj"}),
			},
			wantContains: []slog.Attr{
				slog.GroupAttrs("ns",
					slog.String("nested", "value"),
					slog.GroupAttrs("deep",
						slog.Int("id", 2),
						slog.String("name", "nested-obj"),
					),
				),
			},
		},
		{
			name:         "skip",
			fields:       []zap.Field{zap.Skip()},
			wantContains: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := newObjectEncoder(len(tt.fields))
			for _, field := range tt.fields {
				field.AddTo(enc)
			}
			got := enc.Attrs()

			if len(got) != len(tt.wantContains) {
				t.Fatalf("expected %d slog.Attrs, got %d", len(tt.wantContains), len(got))
			}

			for _, want := range tt.wantContains {
				found := false
				for _, attr := range got {
					if slogAttrEqual(attr, want) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected %v to contain %v", got, want)
				}
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	loggerZap := zap.New(New(context.Background(), noopSlogHandler{}))

	b.ResetTimer()

	for b.Loop() {
		loggerZap.Info("hello world")
	}
}

func BenchmarkNewNoCaller(b *testing.B) {
	loggerZap := zap.New(New(context.Background(), noopSlogHandler{})).WithOptions(zap.WithCaller(false))

	b.ResetTimer()

	for b.Loop() {
		loggerZap.Info("hello world")
	}
}

func BenchmarkNewFields(b *testing.B) {
	loggerZap := zap.New(New(context.Background(), noopSlogHandler{})).WithOptions(zap.WithCaller(false))

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
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}

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

type testObject struct {
	id   int
	name string
}

func (o testObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", o.id)
	enc.AddString("name", o.name)

	return nil
}

type testArray []any

func (a testArray) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, item := range a {
		switch v := item.(type) {
		case int:
			enc.AppendInt(v)
		case string:
			enc.AppendString(v)
		default:
			if err := enc.AppendReflected(v); err != nil {
				return err
			}
		}
	}

	return nil
}

type objectsArray []testObject

func (a objectsArray) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, item := range a {
		if err := enc.AppendObject(item); err != nil {
			return err
		}
	}

	return nil
}

type testNestedObject struct {
	name   string
	nested string
}

func (o testNestedObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", o.name)
	enc.OpenNamespace("ns")
	enc.AddString("nested", o.nested)

	return nil
}

type objectsArrayNested []testNestedObject

func (a objectsArrayNested) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, item := range a {
		if err := enc.AppendObject(item); err != nil {
			return err
		}
	}

	return nil
}

func slogAttrEqual(a, b slog.Attr) bool {
	return a.Key == b.Key && slogValueEqual(a.Value, b.Value)
}

func slogValueEqual(a, b slog.Value) bool {
	aSlice, aOK := a.Any().([]byte)
	bSlice, bOK := b.Any().([]byte)
	if aOK || bOK {
		return aOK && bOK && bytes.Equal(aSlice, bSlice)
	}
	return a.Equal(b)
}
