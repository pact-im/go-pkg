package zapslog

import (
	"log/slog"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"
)

// encodeFields delegates field encoding to zap itself via Field.AddTo and
// captures the result as slog attributes.
func encodeFields(fields []zapcore.Field) []slog.Attr {
	var enc objectEncoder
	for _, field := range fields {
		field.AddTo(&enc)
	}
	return enc.attrs()
}

type namespace struct {
	name    string
	entries []slog.Attr
	child   *namespace
}

// attrs materializes child namespace lazily so namespaces can keep collecting
// fields until the entire object is finished.
func (n *namespace) attrs() []slog.Attr {
	if n.child == nil {
		return n.entries
	}
	attrs := append(make([]slog.Attr, 0, len(n.entries)+1), n.entries...)
	return append(attrs, slog.GroupAttrs(n.child.name, n.child.attrs()...))
}

// objectEncoder implements [zapcore.ObjectEncoder].
// Zap writes fields into an instance of this type, and it builds a tree of slog attributes.
// root keeps the full result, current points at the namespace currently
// receiving fields.
type objectEncoder struct {
	root    namespace // name is unused
	current *namespace
}

// attrs materializes the accumulated object tree into the final slog attrs.
func (e *objectEncoder) attrs() []slog.Attr {
	return e.root.attrs()
}

func (e *objectEncoder) init() {
	if e.current == nil {
		e.current = &e.root
	}
}

func (e *objectEncoder) addAttr(attr slog.Attr) {
	e.init()
	e.current.entries = append(e.current.entries, attr)
}

// OpenNamespace switches subsequent writes into a child namespace, matching zap’s
// namespace semantics where all following fields belong to that group.
func (e *objectEncoder) OpenNamespace(key string) {
	e.init()
	ns := &namespace{name: key}
	e.current.child = ns
	e.current = ns
}

func (e *objectEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	var arr arrayEncoder
	err := marshaler.MarshalLogArray(&arr)
	e.addAttr(slog.GroupAttrs(key, arr.attrs...))
	return err
}

func (e *objectEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	var obj objectEncoder
	err := marshaler.MarshalLogObject(&obj)
	e.addAttr(slog.GroupAttrs(key, obj.attrs()...))
	return err
}

func (e *objectEncoder) AddBinary(key string, value []byte) {
	e.addAttr(slog.Any(key, value))
}

func (e *objectEncoder) AddByteString(key string, value []byte) {
	e.addAttr(slog.String(key, string(value)))
}

func (e *objectEncoder) AddBool(key string, value bool) {
	e.addAttr(slog.Bool(key, value))
}

func (e *objectEncoder) AddComplex128(key string, value complex128) {
	e.addAttr(slog.Any(key, value))
}

func (e *objectEncoder) AddComplex64(key string, value complex64) {
	e.addAttr(slog.Any(key, value))
}

func (e *objectEncoder) AddDuration(key string, value time.Duration) {
	e.addAttr(slog.Duration(key, value))
}

func (e *objectEncoder) AddFloat64(key string, value float64) {
	e.addAttr(slog.Float64(key, value))
}

func (e *objectEncoder) AddFloat32(key string, value float32) {
	e.addAttr(slog.Float64(key, float64(value)))
}

func (e *objectEncoder) AddInt(key string, value int) {
	e.addAttr(slog.Int(key, value))
}

func (e *objectEncoder) AddInt64(key string, value int64) {
	e.addAttr(slog.Int64(key, value))
}

func (e *objectEncoder) AddInt32(key string, value int32) {
	e.addAttr(slog.Int64(key, int64(value)))
}

func (e *objectEncoder) AddInt16(key string, value int16) {
	e.addAttr(slog.Int64(key, int64(value)))
}

func (e *objectEncoder) AddInt8(key string, value int8) {
	e.addAttr(slog.Int64(key, int64(value)))
}

func (e *objectEncoder) AddString(key, value string) {
	e.addAttr(slog.String(key, value))
}

func (e *objectEncoder) AddTime(key string, value time.Time) {
	e.addAttr(slog.Time(key, value))
}

func (e *objectEncoder) AddUint(key string, value uint) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *objectEncoder) AddUint64(key string, value uint64) {
	e.addAttr(slog.Uint64(key, value))
}

func (e *objectEncoder) AddUint32(key string, value uint32) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *objectEncoder) AddUint16(key string, value uint16) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *objectEncoder) AddUint8(key string, value uint8) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *objectEncoder) AddUintptr(key string, value uintptr) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *objectEncoder) AddReflected(key string, value any) error {
	e.addAttr(slog.Any(key, value))
	return nil
}

// arrayEncoder represents array elements as indexed slog attrs and implements
// [zapcore.ArrayEncoder].
type arrayEncoder struct {
	attrs []slog.Attr
}

// appendValue appends array element as an [slog.Attr] keyed by its position.
func (e *arrayEncoder) appendValue(v slog.Value) {
	e.attrs = append(e.attrs, slog.Attr{
		Key:   strconv.Itoa(len(e.attrs)),
		Value: v,
	})
}

func (e *arrayEncoder) AppendArray(marshaler zapcore.ArrayMarshaler) error {
	var arr arrayEncoder
	err := marshaler.MarshalLogArray(&arr)
	e.appendValue(slog.GroupValue(arr.attrs...))
	return err
}

func (e *arrayEncoder) AppendObject(marshaler zapcore.ObjectMarshaler) error {
	var obj objectEncoder
	err := marshaler.MarshalLogObject(&obj)
	e.appendValue(slog.GroupValue(obj.attrs()...))
	return err
}

func (e *arrayEncoder) AppendReflected(value any) error {
	e.appendValue(slog.AnyValue(value))
	return nil
}

func (e *arrayEncoder) AppendBool(value bool) {
	e.appendValue(slog.BoolValue(value))
}

func (e *arrayEncoder) AppendByteString(value []byte) {
	e.appendValue(slog.StringValue(string(value)))
}

func (e *arrayEncoder) AppendComplex128(value complex128) {
	e.appendValue(slog.AnyValue(value))
}

func (e *arrayEncoder) AppendComplex64(value complex64) {
	e.appendValue(slog.AnyValue(value))
}

func (e *arrayEncoder) AppendDuration(value time.Duration) {
	e.appendValue(slog.DurationValue(value))
}

func (e *arrayEncoder) AppendFloat64(value float64) {
	e.appendValue(slog.Float64Value(value))
}

func (e *arrayEncoder) AppendFloat32(value float32) {
	e.appendValue(slog.Float64Value(float64(value)))
}

func (e *arrayEncoder) AppendInt(value int) {
	e.appendValue(slog.Int64Value(int64(value)))
}

func (e *arrayEncoder) AppendInt64(value int64) {
	e.appendValue(slog.Int64Value(value))
}

func (e *arrayEncoder) AppendInt32(value int32) {
	e.appendValue(slog.Int64Value(int64(value)))
}

func (e *arrayEncoder) AppendInt16(value int16) {
	e.appendValue(slog.Int64Value(int64(value)))
}

func (e *arrayEncoder) AppendInt8(value int8) {
	e.appendValue(slog.Int64Value(int64(value)))
}

func (e *arrayEncoder) AppendString(value string) {
	e.appendValue(slog.StringValue(value))
}

func (e *arrayEncoder) AppendTime(value time.Time) {
	e.appendValue(slog.TimeValue(value))
}

func (e *arrayEncoder) AppendUint(value uint) {
	e.appendValue(slog.Uint64Value(uint64(value)))
}

func (e *arrayEncoder) AppendUint64(value uint64) {
	e.appendValue(slog.Uint64Value(value))
}

func (e *arrayEncoder) AppendUint32(value uint32) {
	e.appendValue(slog.Uint64Value(uint64(value)))
}

func (e *arrayEncoder) AppendUint16(value uint16) {
	e.appendValue(slog.Uint64Value(uint64(value)))
}

func (e *arrayEncoder) AppendUint8(value uint8) {
	e.appendValue(slog.Uint64Value(uint64(value)))
}

func (e *arrayEncoder) AppendUintptr(value uintptr) {
	e.appendValue(slog.Uint64Value(uint64(value)))
}
