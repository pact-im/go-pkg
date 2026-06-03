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
	t.Run("levels", func(t *testing.T) {
		var b bytes.Buffer
		handler := slog.NewTextHandler(&b,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)

		loggerZap := zap.New(New(context.Background(), handler))

		loggerZap.Debug("debug level")
		loggerZap.Info("info level")
		loggerZap.Warn("warn level")
		loggerZap.Error("error level")

		err := loggerZap.Sync()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := b.String()
		if !strings.Contains(got, `level=DEBUG msg="debug level"`) {
			t.Fatalf("expected %q to contain %q", got, `level=DEBUG msg="debug level"`)
		}
		if !strings.Contains(got, `level=INFO msg="info level"`) {
			t.Fatalf("expected %q to contain %q", got, `level=INFO msg="info level"`)
		}
		if !strings.Contains(got, `level=WARN msg="warn level"`) {
			t.Fatalf("expected %q to contain %q", got, `level=WARN msg="warn level"`)
		}
		if !strings.Contains(got, `level=ERROR msg="error level"`) {
			t.Fatalf("expected %q to contain %q", got, `level=ERROR msg="error level"`)
		}
	})

	t.Run("level filtering", func(t *testing.T) {
		var b bytes.Buffer
		handler := slog.NewTextHandler(&b,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		)

		loggerZap := zap.New(New(context.Background(), handler)).WithOptions(zap.WithCaller(false))

		loggerZap.Debug("skip debug")
		loggerZap.Info("keep info")

		err := loggerZap.Sync()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := b.String()
		if strings.Contains(got, `skip debug`) {
			t.Fatalf("expected %q not to contain %q", got, `skip debug`)
		}
		if !strings.Contains(got, `keep info`) {
			t.Fatalf("expected %q to contain %q", got, `keep info`)
		}
	})
}

func TestEncodeFields(t *testing.T) {
	testTime := time.Date(2026, time.May, 27, 12, 34, 56, 0, time.UTC)
	fullTime := testTime.Add(123 * time.Nanosecond)
	tests := []struct {
		name         string
		fields       []zap.Field
		wantContains []string
		notContains  []string
	}{
		{
			name:         "string",
			fields:       []zap.Field{zap.String("key", "value")},
			wantContains: []string{`"key":"value"`},
		},
		{
			name:         "int",
			fields:       []zap.Field{zap.Int("count", 42)},
			wantContains: []string{`"count":42`},
		},
		{
			name:         "float64",
			fields:       []zap.Field{zap.Float64("float64", 123.456789)},
			wantContains: []string{`"float64":123.456789`},
		},
		{
			name:         "float32",
			fields:       []zap.Field{zap.Float32("float32", 12.5)},
			wantContains: []string{`"float32":12.5`},
		},
		{
			name:         "bool",
			fields:       []zap.Field{zap.Bool("flag", true)},
			wantContains: []string{`"flag":true`},
		},
		{
			name:         "duration",
			fields:       []zap.Field{zap.Duration("timeout", 2*time.Second)},
			wantContains: []string{`"timeout":2000000000`},
		},
		{
			name:         "time",
			fields:       []zap.Field{zap.Time("created_at", testTime)},
			wantContains: []string{`"created_at":"2026-05-27T12:34:56Z"`},
		},
		{
			name:         "int32",
			fields:       []zap.Field{zap.Int32("int32", 32)},
			wantContains: []string{`"int32":32`},
		},
		{
			name:         "int16",
			fields:       []zap.Field{zap.Int16("int16", 16)},
			wantContains: []string{`"int16":16`},
		},
		{
			name:         "int8",
			fields:       []zap.Field{zap.Int8("int8", 8)},
			wantContains: []string{`"int8":8`},
		},
		{
			name:         "uint64",
			fields:       []zap.Field{zap.Uint64("uint64", 64)},
			wantContains: []string{`"uint64":64`},
		},
		{
			name:         "uint32",
			fields:       []zap.Field{zap.Uint32("uint32", 32)},
			wantContains: []string{`"uint32":32`},
		},
		{
			name:         "uint16",
			fields:       []zap.Field{zap.Uint16("uint16", 16)},
			wantContains: []string{`"uint16":16`},
		},
		{
			name:         "uint8",
			fields:       []zap.Field{zap.Uint8("uint8", 8)},
			wantContains: []string{`"uint8":8`},
		},
		{
			name:         "uintptr",
			fields:       []zap.Field{zap.Uintptr("pointer", uintptr(123))},
			wantContains: []string{`"pointer":123`},
		},
		{
			name:         "time full",
			fields:       []zap.Field{{Key: "full_time", Type: zapcore.TimeFullType, Interface: fullTime}},
			wantContains: []string{`"full_time":"2026-05-27T12:34:56.000000123Z"`},
		},
		{
			name:         "stringer",
			fields:       []zap.Field{zap.Stringer("stringer", testStringer("stringer-value"))},
			wantContains: []string{`"stringer":"stringer-value"`},
		},
		{
			name:         "error",
			fields:       []zap.Field{zap.Error(errors.New("boom"))},
			wantContains: []string{`"error":"boom"`},
		},
		{
			name:         "binary",
			fields:       []zap.Field{zap.Binary("binary", []byte("abc"))},
			wantContains: []string{`"binary":"YWJj"`},
		},
		{
			name:         "byte string",
			fields:       []zap.Field{zap.ByteString("byte_string", []byte("hello"))},
			wantContains: []string{`"byte_string":"hello"`},
		},
		{
			name:         "complex64",
			fields:       []zap.Field{zap.Complex64("complex64", complex64(1+2i))},
			wantContains: []string{`"complex64":"!ERROR:json: unsupported type: complex64"`},
		},
		{
			name:         "complex128",
			fields:       []zap.Field{zap.Complex128("complex128", complex128(3+4i))},
			wantContains: []string{`"complex128":"!ERROR:json: unsupported type: complex128"`},
		},
		{
			name:         "reflect",
			fields:       []zap.Field{zap.Reflect("reflect", map[string]any{"answer": 42})},
			wantContains: []string{`"reflect":{"answer":42}`},
		},
		{
			name:         "array",
			fields:       []zap.Field{zap.Array("array", testArray{1, "two"})},
			wantContains: []string{`"array":{"0":1,"1":"two"}`},
		},
		{
			name: "array objects",
			fields: []zap.Field{zap.Array("array_objects", objectsArray{
				{id: 1, name: "value1"},
				{id: 2, name: "value2"},
			})},
			wantContains: []string{`"array_objects":{"0":{"id":1,"name":"value1"},"1":{"id":2,"name":"value2"}}`},
		},
		{
			name: "array nested objects",
			fields: []zap.Field{zap.Array("array_nested_objects", objectsArrayNested{
				{name: "value1", nested: "nested1"},
				{name: "value2", nested: "nested2"},
			})},
			wantContains: []string{`"array_nested_objects":{"0":{"name":"value1","ns":{"nested":"nested1"}},"1":{"name":"value2","ns":{"nested":"nested2"}}}`},
		},
		{
			name:         "object",
			fields:       []zap.Field{zap.Object("object", testObject{id: 1, name: "obj"})},
			wantContains: []string{`"object":{"id":1,"name":"obj"}`},
		},
		{
			name:         "inline",
			fields:       []zap.Field{zap.Inline(testObject{id: 7, name: "inline"})},
			wantContains: []string{`"id":7`, `"name":"inline"`},
		},
		{
			name: "namespace",
			fields: []zap.Field{
				zap.Namespace("ns"),
				zap.String("nested", "value"),
				zap.Object("deep", testObject{id: 2, name: "nested-obj"}),
			},
			wantContains: []string{`"ns":{"nested":"value","deep":{"id":2,"name":"nested-obj"}}`},
		},
		{
			name:         "skip",
			fields:       []zap.Field{zap.Skip()},
			wantContains: []string{`"msg":"fields"`},
			notContains:  []string{`"skip"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			handler := slog.NewJSONHandler(&b,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			)
			loggerZap := zap.New(New(context.Background(), handler)).WithOptions(zap.WithCaller(false))

			loggerZap.Info("fields", tt.fields...)

			if err := loggerZap.Sync(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := b.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Fatalf("expected %q to contain %q", got, want)
				}
			}

			for _, notWant := range tt.notContains {
				if strings.Contains(got, notWant) {
					t.Fatalf("expected %q not to contain %q", got, notWant)
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
