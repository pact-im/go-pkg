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

		loggerZap := zap.New(New(context.Background(), handler))

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
		requireNoError(t, err)

		got := b.String()
		assertNotContains(t, got, `skip debug`)
		assertContains(t, got, `keep info`)
	})

	t.Run("named logger and fields", func(t *testing.T) {
		var b bytes.Buffer
		handler := slog.NewTextHandler(&b,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)

		loggerZap := zap.New(New(context.Background(), handler)).WithOptions(zap.WithCaller(false))
		fullTime := testTime.Add(123 * time.Nanosecond)

		loggerZap = loggerZap.Named("example").With(
			zap.String("base_key", "base_value"),
			zap.Int("base_count", 7),
		)

		loggerZap.Info("fields",
			zap.String("key", "value"),
			zap.Int("count", 42),
			zap.Float64("float64", 123.456789),
			zap.Float32("float32", 12.5),
			zap.Bool("flag", true),
			zap.Duration("timeout", 2*time.Second),
			zap.Time("created_at", testTime),
			zap.Int32("int32", 32),
			zap.Int16("int16", 16),
			zap.Int8("int8", 8),
			zap.Uint64("uint64", 64),
			zap.Uint32("uint32", 32),
			zap.Uint16("uint16", 16),
			zap.Uint8("uint8", 8),
			zap.Stringer("stringer", testStringer("stringer-value")),
			zap.Error(errors.New("boom")),
			zap.Binary("binary", []byte("abc")),
			zap.ByteString("byte_string", []byte("hello")),
			zap.Complex64("complex64", complex64(1+2i)),
			zap.Complex128("complex128", complex128(3+4i)),
			zap.Array("array", testArray{1, "two"}),
			zap.Array("array_objects", objectsArray{
				{id: 1, name: "value1"},
				{id: 2, name: "value2"},
			}),
			zap.Array("array_nested_objects", objectsArrayNested{
				{name: "value1", nested: "nested1"},
				{name: "value2", nested: "nested2"},
			}),
			zap.Object("object", testObject{id: 1, name: "obj"}),
			zap.Inline(testObject{id: 7, name: "inline"}),
			zap.Reflect("reflect", map[string]any{"answer": 42}),
			zap.Uintptr("pointer", uintptr(123)),
			zap.Field{Key: "full_time", Type: zapcore.TimeFullType, Interface: fullTime},
			zap.Skip(),
			zap.Namespace("ns"),
			zap.String("nested", "value"),
			zap.Object("deep", testObject{id: 2, name: "nested-obj"}),
		)

		err := loggerZap.Sync()
		requireNoError(t, err)

		got := b.String()
		assertContains(t, got, `msg=fields`)
		assertContains(t, got, `logger_name=example`)
		assertContains(t, got, `base_key=base_value`)
		assertContains(t, got, `base_count=7`)
		assertContains(t, got, `key=value`)
		assertContains(t, got, `count=42`)
		assertContains(t, got, `float64=123.456789`)
		assertContains(t, got, `float32=12.5`)
		assertContains(t, got, `flag=true`)
		assertContains(t, got, `timeout=2s`)
		assertContains(t, got, `created_at=2026-05-27T12:34:56.000Z`)
		assertContains(t, got, `int32=32`)
		assertContains(t, got, `int16=16`)
		assertContains(t, got, `int8=8`)
		assertContains(t, got, `uint64=64`)
		assertContains(t, got, `uint32=32`)
		assertContains(t, got, `uint16=16`)
		assertContains(t, got, `uint8=8`)
		assertContains(t, got, `stringer=stringer-value`)
		assertContains(t, got, `error=boom`)
		assertContains(t, got, `binary="abc"`)
		assertContains(t, got, `byte_string=hello`)
		assertContains(t, got, `complex64=(1+2i)`)
		assertContains(t, got, `complex128=(3+4i)`)
		assertContains(t, got, `array.0=1`)
		assertContains(t, got, `array.1=two`)
		assertContains(t, got, `array_objects.0.id=1`)
		assertContains(t, got, `array_objects.0.name=value1`)
		assertContains(t, got, `array_objects.1.id=2`)
		assertContains(t, got, `array_objects.1.name=value2`)
		assertContains(t, got, `array_nested_objects.0.name=value1`)
		assertContains(t, got, `array_nested_objects.0.ns.nested=nested1`)
		assertContains(t, got, `array_nested_objects.1.name=value2`)
		assertContains(t, got, `array_nested_objects.1.ns.nested=nested2`)
		assertContains(t, got, `object.id=1`)
		assertContains(t, got, `object.name=obj`)
		assertContains(t, got, `id=7`)
		assertContains(t, got, `name=inline`)
		assertContains(t, got, `reflect=map[answer:42]`)
		assertContains(t, got, `pointer=123`)
		assertContains(t, got, `full_time=2026-05-27T12:34:56.000Z`)
		assertContains(t, got, `ns.nested=value`)
		assertContains(t, got, `ns.deep.id=2`)
		assertContains(t, got, `ns.deep.name=nested-obj`)
		assertNotContains(t, got, `skip`)
	})

	t.Run("named logger and fields json", func(t *testing.T) {
		var b bytes.Buffer
		handler := slog.NewJSONHandler(&b,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)

		loggerZap := zap.New(New(context.Background(), handler)).WithOptions(zap.WithCaller(false))
		fullTime := testTime.Add(123 * time.Nanosecond)

		loggerZap = loggerZap.Named("example").With(
			zap.String("base_key", "base_value"),
			zap.Int("base_count", 7),
		)

		loggerZap.Info("fields",
			zap.String("key", "value"),
			zap.Int("count", 42),
			zap.Float64("float64", 123.456789),
			zap.Float32("float32", 12.5),
			zap.Bool("flag", true),
			zap.Duration("timeout", 2*time.Second),
			zap.Time("created_at", testTime),
			zap.Int32("int32", 32),
			zap.Int16("int16", 16),
			zap.Int8("int8", 8),
			zap.Uint64("uint64", 64),
			zap.Uint32("uint32", 32),
			zap.Uint16("uint16", 16),
			zap.Uint8("uint8", 8),
			zap.Stringer("stringer", testStringer("stringer-value")),
			zap.Error(errors.New("boom")),
			zap.Binary("binary", []byte("abc")),
			zap.ByteString("byte_string", []byte("hello")),
			zap.Complex64("complex64", complex64(1+2i)),
			zap.Complex128("complex128", complex128(3+4i)),
			zap.Array("array", testArray{1, "two"}),
			zap.Array("array_objects", objectsArray{
				{id: 1, name: "value1"},
				{id: 2, name: "value2"},
			}),
			zap.Array("array_nested_objects", objectsArrayNested{
				{name: "value1", nested: "nested1"},
				{name: "value2", nested: "nested2"},
			}),
			zap.Object("object", testObject{id: 1, name: "obj"}),
			zap.Inline(testObject{id: 7, name: "inline"}),
			zap.Reflect("reflect", map[string]any{"answer": 42}),
			zap.Uintptr("pointer", uintptr(123)),
			zap.Field{Key: "full_time", Type: zapcore.TimeFullType, Interface: fullTime},
			zap.Skip(),
			zap.Namespace("ns"),
			zap.String("nested", "value"),
			zap.Object("deep", testObject{id: 2, name: "nested-obj"}),
		)

		err := loggerZap.Sync()
		requireNoError(t, err)

		got := b.String()
		assertContains(t, got, `"msg":"fields"`)
		assertContains(t, got, `"logger_name":"example"`)
		assertContains(t, got, `"base_key":"base_value"`)
		assertContains(t, got, `"base_count":7`)
		assertContains(t, got, `"key":"value"`)
		assertContains(t, got, `"count":42`)
		assertContains(t, got, `"float64":123.456789`)
		assertContains(t, got, `"float32":12.5`)
		assertContains(t, got, `"flag":true`)
		assertContains(t, got, `"timeout":2000000000`)
		assertContains(t, got, `"created_at":"2026-05-27T12:34:56Z"`)
		assertContains(t, got, `"int32":32`)
		assertContains(t, got, `"int16":16`)
		assertContains(t, got, `"int8":8`)
		assertContains(t, got, `"uint64":64`)
		assertContains(t, got, `"uint32":32`)
		assertContains(t, got, `"uint16":16`)
		assertContains(t, got, `"uint8":8`)
		assertContains(t, got, `"stringer":"stringer-value"`)
		assertContains(t, got, `"error":"boom"`)
		assertContains(t, got, `"binary":"YWJj"`)
		assertContains(t, got, `"byte_string":"hello"`)
		assertContains(t, got, `"complex64":`)
		assertContains(t, got, `"complex128":`)
		assertContains(t, got, `"array":{"0":1,"1":"two"}`)
		assertContains(t, got, `"array_objects":{"0":{"id":1,"name":"value1"},"1":{"id":2,"name":"value2"}}`)
		assertContains(t, got, `"array_nested_objects":{"0":{"name":"value1","ns":{"nested":"nested1"}},"1":{"name":"value2","ns":{"nested":"nested2"}}}`)
		assertContains(t, got, `"object":{"id":1,"name":"obj"}`)
		assertContains(t, got, `"id":7`)
		assertContains(t, got, `"name":"inline"`)
		assertContains(t, got, `"reflect":{"answer":42}`)
		assertContains(t, got, `"pointer":123`)
		assertContains(t, got, `"full_time":"2026-05-27T12:34:56.000000123Z"`)
		assertContains(t, got, `"ns":{"nested":"value","deep":{"id":2,"name":"nested-obj"}}`)
		assertNotContains(t, got, `"skip"`)
	})
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
