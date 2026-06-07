package zapslog

import (
	"bytes"
	"log/slog"
	"slices"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestEncodeFields(t *testing.T) {
	tests := []struct {
		name   string
		fields []zap.Field
		want   []slog.Attr
	}{
		{
			name:   "skip",
			fields: []zap.Field{zap.Skip()},
			want:   nil,
		},
		{
			name:   "int",
			fields: []zap.Field{zap.Int("count", 42)},
			want:   []slog.Attr{slog.Int("count", 42)},
		},
		{
			name:   "int64",
			fields: []zap.Field{zap.Int64("int64", 64)},
			want:   []slog.Attr{slog.Int64("int64", 64)},
		},
		{
			name:   "int32",
			fields: []zap.Field{zap.Int32("int32", 32)},
			want:   []slog.Attr{slog.Int64("int32", 32)},
		},
		{
			name:   "int16",
			fields: []zap.Field{zap.Int16("int16", 16)},
			want:   []slog.Attr{slog.Int64("int16", 16)},
		},
		{
			name:   "int8",
			fields: []zap.Field{zap.Int8("int8", 8)},
			want:   []slog.Attr{slog.Int64("int8", 8)},
		},
		{
			name:   "uint",
			fields: []zap.Field{zap.Uint("uint", 42)},
			want:   []slog.Attr{slog.Uint64("uint", 42)},
		},
		{
			name:   "uint64",
			fields: []zap.Field{zap.Uint64("uint64", 64)},
			want:   []slog.Attr{slog.Uint64("uint64", 64)},
		},
		{
			name:   "uint32",
			fields: []zap.Field{zap.Uint32("uint32", 32)},
			want:   []slog.Attr{slog.Uint64("uint32", 32)},
		},
		{
			name:   "uint16",
			fields: []zap.Field{zap.Uint16("uint16", 16)},
			want:   []slog.Attr{slog.Uint64("uint16", 16)},
		},
		{
			name:   "uint8",
			fields: []zap.Field{zap.Uint8("uint8", 8)},
			want:   []slog.Attr{slog.Uint64("uint8", 8)},
		},
		{
			name:   "uintptr",
			fields: []zap.Field{zap.Uintptr("pointer", uintptr(123))},
			want:   []slog.Attr{slog.Uint64("pointer", 123)},
		},
		{
			name:   "float64",
			fields: []zap.Field{zap.Float64("float64", 123.456789)},
			want:   []slog.Attr{slog.Float64("float64", 123.456789)},
		},
		{
			name:   "float32",
			fields: []zap.Field{zap.Float32("float32", 12.5)},
			want:   []slog.Attr{slog.Float64("float32", 12.5)},
		},
		{
			name:   "bool",
			fields: []zap.Field{zap.Bool("flag", true)},
			want:   []slog.Attr{slog.Bool("flag", true)},
		},
		{
			name:   "time",
			fields: []zap.Field{zap.Time("created_at", time.Date(2026, time.May, 27, 12, 34, 56, 789, time.UTC))},
			want:   []slog.Attr{slog.Time("created_at", time.Date(2026, time.May, 27, 12, 34, 56, 789, time.UTC))},
		},
		{
			name:   "duration",
			fields: []zap.Field{zap.Duration("timeout", 2*time.Second)},
			want:   []slog.Attr{slog.Duration("timeout", 2*time.Second)},
		},
		{
			name:   "string",
			fields: []zap.Field{zap.String("key", "value")},
			want:   []slog.Attr{slog.String("key", "value")},
		},
		{
			name: "stringer",
			fields: []zap.Field{zap.Stringer("stringer", stringerFunc(func() string {
				return "stringer-value"
			}))},
			want: []slog.Attr{slog.String("stringer", "stringer-value")},
		},
		{
			name: "error",
			fields: []zap.Field{zap.Error(errorFunc(func() string {
				return "boom"
			}))},
			want: []slog.Attr{slog.String("error", "boom")},
		},
		{
			name:   "binary",
			fields: []zap.Field{zap.Binary("binary", []byte("abc"))},
			want:   []slog.Attr{slog.Any("binary", []byte("abc"))},
		},
		{
			name:   "byte string",
			fields: []zap.Field{zap.ByteString("byte_string", []byte("hello"))},
			want:   []slog.Attr{slog.String("byte_string", "hello")},
		},
		{
			name:   "complex64",
			fields: []zap.Field{zap.Complex64("complex64", complex64(1+2i))},
			want:   []slog.Attr{slog.Any("complex64", complex64(1+2i))},
		},
		{
			name:   "complex128",
			fields: []zap.Field{zap.Complex128("complex128", complex128(3+4i))},
			want:   []slog.Attr{slog.Any("complex128", complex128(3+4i))},
		},
		{
			name:   "reflect",
			fields: []zap.Field{zap.Reflect("reflect", []byte("hello"))},
			want:   []slog.Attr{slog.Any("reflect", []byte("hello"))},
		},
		{
			name: "array",
			fields: []zap.Field{
				zap.Array("array", arrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
					enc.AppendInt64(1)
					enc.AppendString("two")
					return nil
				})),
			},
			want: []slog.Attr{
				slog.GroupAttrs("array",
					slog.Int64("0", 1),
					slog.String("1", "two"),
				),
			},
		},
		{
			name: "array objects",
			fields: []zap.Field{
				zap.Array("array_objects", arrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
					if err := enc.AppendObject(objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
						enc.AddInt("id", 1)
						enc.AddString("name", "value1")
						return nil
					})); err != nil {
						return err
					}
					if err := enc.AppendObject(objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
						enc.AddInt("id", 2)
						enc.AddString("name", "value2")
						return nil
					})); err != nil {
						return err
					}
					return nil
				})),
			},
			want: []slog.Attr{
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
			name: "array nested arrays",
			fields: []zap.Field{
				zap.Array("array_nested_arrays", arrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
					if err := enc.AppendArray(arrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
						enc.AppendInt(1)
						enc.AppendInt(2)
						return nil
					})); err != nil {
						return err
					}
					return nil
				})),
			},
			want: []slog.Attr{
				slog.GroupAttrs("array_nested_arrays",
					slog.GroupAttrs("0",
						slog.Int("0", 1),
						slog.Int("1", 2),
					),
				),
			},
		},
		{
			name: "array nested objects",
			fields: []zap.Field{
				zap.Array("array_nested_objects", arrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
					if err := enc.AppendObject(objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
						enc.AddString("name", "value1")
						if err := enc.AddObject("ns", objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
							enc.AddString("nested", "nested1")
							return nil
						})); err != nil {
							return err
						}
						return nil
					})); err != nil {
						return err
					}
					if err := enc.AppendObject(objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
						enc.AddString("name", "value2")
						if err := enc.AddObject("ns", objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
							enc.AddString("nested", "nested2")
							return nil
						})); err != nil {
							return err
						}
						return nil
					})); err != nil {
						return err
					}
					return nil
				})),
			},
			want: []slog.Attr{
				slog.GroupAttrs("array_nested_objects",
					slog.GroupAttrs("0",
						slog.String("name", "value1"),
						slog.GroupAttrs("ns",
							slog.String("nested", "nested1"),
						),
					),
					slog.GroupAttrs("1",
						slog.String("name", "value2"),
						slog.GroupAttrs("ns",
							slog.String("nested", "nested2"),
						),
					),
				),
			},
		},
		{
			name: "object",
			fields: []zap.Field{
				zap.Object("object", objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddInt("id", 1)
					enc.AddString("name", "obj")
					return nil
				})),
			},
			want: []slog.Attr{
				slog.GroupAttrs("object",
					slog.Int("id", 1),
					slog.String("name", "obj"),
				),
			},
		},
		{
			name: "inline",
			fields: []zap.Field{
				zap.Inline(objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddUint("id", 7)
					enc.AddString("name", "inline")
					return nil
				})),
			},
			want: []slog.Attr{
				slog.Uint64("id", 7),
				slog.String("name", "inline"),
			},
		},
		{
			name: "namespace",
			fields: []zap.Field{
				zap.Namespace("ns"),
				zap.String("nested", "value"),
				zap.Object("deep", objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddInt("id", 2)
					enc.AddString("name", "nested-obj")
					return nil
				})),
			},
			want: []slog.Attr{
				slog.GroupAttrs("ns",
					slog.String("nested", "value"),
					slog.GroupAttrs("deep",
						slog.Int("id", 2),
						slog.String("name", "nested-obj"),
					),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := encodeFields(tt.fields)
			if !slices.EqualFunc(attrs, tt.want, slogAttrEqual) {
				t.Fatalf("encodeFields() mismatch:\n got: %#v\nwant: %#v", attrs, tt.want)
			}
		})
	}
}

// stringerFunc is a function that implements the [fmt.Stringer] interface.
type stringerFunc func() string

func (f stringerFunc) String() string {
	return f()
}

// errorFunc is a function that implements the [error] interface.
type errorFunc func() string

func (f errorFunc) Error() string {
	return f()
}

// objectMarshalerFunc is a function that implements the
// [zapcore.ObjectMarshaler] interface.
type objectMarshalerFunc func(enc zapcore.ObjectEncoder) error

func (f objectMarshalerFunc) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return f(enc)
}

// arrayMarshalerFunc is a function that implements the [zapcore.ArrayMarshaler]
// interface.
type arrayMarshalerFunc func(enc zapcore.ArrayEncoder) error

func (f arrayMarshalerFunc) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	return f(enc)
}

// slogAttrEqual returns true if two [slog.Attr] values are equal, false
// otherwise.
func slogAttrEqual(a, b slog.Attr) bool {
	return a.Key == b.Key && slogValueEqual(a.Value, b.Value)
}

// slogValueEqual returns true if two [slog.Value] values are equal, false
// otherwise.
func slogValueEqual(a, b slog.Value) bool {
	aSlice, aOK := a.Any().([]byte)
	bSlice, bOK := b.Any().([]byte)
	if aOK || bOK {
		return aOK && bOK && bytes.Equal(aSlice, bSlice)
	}
	return a.Equal(b)
}
