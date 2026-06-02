// Package zapslog provides a zapcore.Core implementation that forwards logs to
// slog.Handler.
package zapslog

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.uber.org/zap/zapcore"
)

// New creates a Core backed by the provided context and slog.Handler.
func New(ctx context.Context, handler slog.Handler) *Core {
	return &Core{
		ctx:     ctx,
		handler: handler,
	}
}

// Core implements zapcore.Core and forwards log records to slog.Handler.
type Core struct {
	ctx     context.Context
	handler slog.Handler
}

// Enabled reports whether the underlying slog handler accepts the given level.
func (c *Core) Enabled(level zapcore.Level) bool {
	return c.handler.Enabled(c.ctx, zapCoreLevelToSlogLevel(level))
}

// encodeFields delegates field encoding to zap itself via Field.AddTo and
// captures the result as slog attributes.
func encodeFields(fields []zapcore.Field) []slog.Attr {
	enc := newObjectEncoder(len(fields))
	for _, field := range fields {
		field.AddTo(enc)
	}

	return enc.Attrs()
}

// With implements the [zapcore.Core] interface.
func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	handler := c.handler.WithAttrs(encodeFields(fields))

	return &Core{
		ctx:     c.ctx,
		handler: handler,
	}
}

// Check implements the [zapcore.Core] interface.
func (c *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}

	return ce
}

// Write implements the [zapcore.Core] interface.
func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// https://pkg.go.dev/log/slog#hdr-Writing_a_handler
	record := slog.NewRecord(entry.Time, zapCoreLevelToSlogLevel(entry.Level), entry.Message, entry.Caller.PC)

	if entry.LoggerName != "" {
		record.AddAttrs(slog.String("name", entry.LoggerName))
	}

	record.AddAttrs(encodeFields(fields)...)

	if entry.Stack != "" {
		record.AddAttrs(slog.String("stack", entry.Stack))
	}

	err := c.handler.Handle(c.ctx, record)
	if err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	return nil
}

// Sync implements the [zapcore.Core] interface.
func (c *Core) Sync() error {
	return nil
}

// zapCoreLevelToSlogLevel converts a zapcore.Level to a slog.Level.
// Unsupported levels are converted to slog.LevelDebug.
func zapCoreLevelToSlogLevel(level zapcore.Level) slog.Level {
	switch level {
	case zapcore.DebugLevel:
		return slog.LevelDebug
	case zapcore.InfoLevel:
		return slog.LevelInfo
	case zapcore.WarnLevel:
		return slog.LevelWarn
	case zapcore.ErrorLevel:
		return slog.LevelError
	case zapcore.DPanicLevel:
		return slog.LevelError
	case zapcore.PanicLevel:
		return slog.LevelError
	case zapcore.FatalLevel:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// objectEncoder implements [zapcore.ObjectEncoder].
// Zap writes fields into an instance of this type, and it builds a tree of slog attributes.
// root keeps the full result, current points at the namespace currently
// receiving fields.
type objectEncoder struct {
	root, current *namespace
}

type namespace struct {
	name    string
	entries []slog.Attr
	child   *namespace
}

func newObjectEncoder(capacity int) *objectEncoder {
	root := &namespace{entries: make([]slog.Attr, 0, capacity)}
	return &objectEncoder{
		root:    root,
		current: root,
	}
}

// Attrs materializes the accumulated object tree into the final slog attrs.
func (e *objectEncoder) Attrs() []slog.Attr {
	return e.root.attrs()
}

func (e *objectEncoder) addAttr(attr slog.Attr) {
	e.current.entries = append(e.current.entries, attr)
}

// attrs materializes child namespace lazily so namespaces can keep collecting
// fields until the entire object is finished.
func (n *namespace) attrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, len(n.entries))
	attrs = append(attrs, n.entries...)
	if n.child != nil {
		attrs = append(attrs, slog.Attr{
			Key:   n.child.name,
			Value: slog.GroupValue(n.child.attrs()...),
		})
	}

	return attrs
}

func (e *objectEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	arr := &arrayEncoder{elems: make([]any, 0)}
	err := marshaler.MarshalLogArray(arr)
	e.addAttr(slog.Any(key, arr.elems))
	return err
}

func (e *objectEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	obj := newObjectEncoder(0)
	err := marshaler.MarshalLogObject(obj)
	e.addAttr(slog.Attr{Key: key, Value: slog.GroupValue(obj.Attrs()...)})
	return err
}

func (e *objectEncoder) AddBinary(key string, value []byte) { e.addAttr(slog.Any(key, value)) }
func (e *objectEncoder) AddByteString(key string, value []byte) {
	e.addAttr(slog.String(key, string(value)))
}
func (e *objectEncoder) AddBool(key string, value bool) { e.addAttr(slog.Bool(key, value)) }
func (e *objectEncoder) AddComplex128(key string, value complex128) {
	e.addAttr(slog.String(key, fmt.Sprint(value)))
}

func (e *objectEncoder) AddComplex64(key string, value complex64) {
	e.addAttr(slog.String(key, fmt.Sprint(value)))
}

func (e *objectEncoder) AddDuration(key string, value time.Duration) {
	e.addAttr(slog.Duration(key, value))
}
func (e *objectEncoder) AddFloat64(key string, value float64) { e.addAttr(slog.Float64(key, value)) }
func (e *objectEncoder) AddFloat32(key string, value float32) {
	e.addAttr(slog.Float64(key, float64(value)))
}
func (e *objectEncoder) AddInt(key string, value int)        { e.addAttr(slog.Int(key, value)) }
func (e *objectEncoder) AddInt64(key string, value int64)    { e.addAttr(slog.Int64(key, value)) }
func (e *objectEncoder) AddInt32(key string, value int32)    { e.addAttr(slog.Int64(key, int64(value))) }
func (e *objectEncoder) AddInt16(key string, value int16)    { e.addAttr(slog.Int64(key, int64(value))) }
func (e *objectEncoder) AddInt8(key string, value int8)      { e.addAttr(slog.Int64(key, int64(value))) }
func (e *objectEncoder) AddString(key, value string)         { e.addAttr(slog.String(key, value)) }
func (e *objectEncoder) AddTime(key string, value time.Time) { e.addAttr(slog.Time(key, value)) }
func (e *objectEncoder) AddUint(key string, value uint)      { e.addAttr(slog.Uint64(key, uint64(value))) }
func (e *objectEncoder) AddUint64(key string, value uint64)  { e.addAttr(slog.Uint64(key, value)) }
func (e *objectEncoder) AddUint32(key string, value uint32) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *objectEncoder) AddUint16(key string, value uint16) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}
func (e *objectEncoder) AddUint8(key string, value uint8) { e.addAttr(slog.Uint64(key, uint64(value))) }
func (e *objectEncoder) AddUintptr(key string, value uintptr) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *objectEncoder) AddReflected(key string, value any) error {
	e.addAttr(slog.Any(key, value))
	return nil
}

// OpenNamespace switches subsequent writes into a child namespace, matching zap's
// namespace semantics where all following fields belong to that group.
func (e *objectEncoder) OpenNamespace(key string) {
	ns := &namespace{name: key}
	e.current.child = ns
	e.current = ns
}

// arrayEncoder collects unnamed array elements and implements
// [zapcore.ArrayEncoder].
type arrayEncoder struct {
	elems []any
}

func (e *arrayEncoder) AppendArray(marshaler zapcore.ArrayMarshaler) error {
	arr := &arrayEncoder{elems: make([]any, 0)}
	err := marshaler.MarshalLogArray(arr)
	e.elems = append(e.elems, arr.elems)
	return err
}

func (e *arrayEncoder) AppendObject(marshaler zapcore.ObjectMarshaler) error {
	obj := newObjectEncoder(0)
	err := marshaler.MarshalLogObject(obj)
	e.elems = append(e.elems, attrsToMap(obj.Attrs()))
	return err
}

func (e *arrayEncoder) AppendReflected(value any) error {
	e.elems = append(e.elems, value)
	return nil
}

func (e *arrayEncoder) AppendBool(value bool)         { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendByteString(value []byte) { e.elems = append(e.elems, string(value)) }
func (e *arrayEncoder) AppendComplex128(value complex128) {
	e.elems = append(e.elems, fmt.Sprint(value))
}
func (e *arrayEncoder) AppendComplex64(value complex64)    { e.elems = append(e.elems, fmt.Sprint(value)) }
func (e *arrayEncoder) AppendDuration(value time.Duration) { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendFloat64(value float64)        { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendFloat32(value float32)        { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendInt(value int)                { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendInt64(value int64)            { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendInt32(value int32)            { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendInt16(value int16)            { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendInt8(value int8)              { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendString(value string)          { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendTime(value time.Time)         { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendUint(value uint)              { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendUint64(value uint64)          { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendUint32(value uint32)          { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendUint16(value uint16)          { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendUint8(value uint8)            { e.elems = append(e.elems, value) }
func (e *arrayEncoder) AppendUintptr(value uintptr)        { e.elems = append(e.elems, value) }

// attrsToMap turns an object's named attrs into a single value that can live
// inside []any, for example as an element of a slog array.
func attrsToMap(attrs []slog.Attr) map[string]any {
	fields := make(map[string]any, len(attrs))
	for _, attr := range attrs {
		if attr.Key == "" {
			continue
		}
		fields[attr.Key] = attr.Value.Any()
	}

	return fields
}
