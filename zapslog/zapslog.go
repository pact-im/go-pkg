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
	enc := newEncoder(len(fields))
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

// encoder implements [zapcore.ObjectEncoder].
// Zap writes fields into an instance of this type, and it builds a tree of slog attributes.
// root keeps the full result, cur points at the namespace/object currently
// receiving fields.
type encoder struct {
	root *node
	cur  *node
}

type node struct {
	entries []objectEntry
}

type objectEntry struct {
	attr  slog.Attr
	child *groupEntry
}

type groupEntry struct {
	key  string
	node *node
}

func newEncoder(capEncoder int) *encoder {
	root := &node{entries: make([]objectEntry, 0, capEncoder)}
	return &encoder{
		root: root,
		cur:  root,
	}
}

// Attrs materializes the accumulated object tree into the final slog attrs.
func (e *encoder) Attrs() []slog.Attr {
	return e.root.attrs()
}

func (e *encoder) addAttr(attr slog.Attr) {
	e.cur.entries = append(e.cur.entries, objectEntry{attr: attr})
}

// attrs materializes child groups lazily so namespaces can keep collecting
// fields until the entire object is finished.
func (n *node) attrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, len(n.entries))
	for _, entry := range n.entries {
		if entry.child == nil {
			attrs = append(attrs, entry.attr)
			continue
		}

		attrs = append(attrs, slog.Attr{
			Key:   entry.child.key,
			Value: slog.GroupValue(entry.child.node.attrs()...),
		})
	}

	return attrs
}

func (e *encoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	arr := &arrayEncoder{elems: make([]any, 0)}
	err := marshaler.MarshalLogArray(arr)
	e.addAttr(slog.Any(key, arr.elems))
	return err
}

func (e *encoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	obj := newEncoder(0)
	err := marshaler.MarshalLogObject(obj)
	e.addAttr(slog.Attr{Key: key, Value: slog.GroupValue(obj.Attrs()...)})
	return err
}

func (e *encoder) AddBinary(key string, value []byte) { e.addAttr(slog.Any(key, value)) }
func (e *encoder) AddByteString(key string, value []byte) {
	e.addAttr(slog.String(key, string(value)))
}
func (e *encoder) AddBool(key string, value bool) { e.addAttr(slog.Bool(key, value)) }
func (e *encoder) AddComplex128(key string, value complex128) {
	e.addAttr(slog.String(key, fmt.Sprint(value)))
}

func (e *encoder) AddComplex64(key string, value complex64) {
	e.addAttr(slog.String(key, fmt.Sprint(value)))
}

func (e *encoder) AddDuration(key string, value time.Duration) {
	e.addAttr(slog.Duration(key, value))
}
func (e *encoder) AddFloat64(key string, value float64) { e.addAttr(slog.Float64(key, value)) }
func (e *encoder) AddFloat32(key string, value float32) {
	e.addAttr(slog.Float64(key, float64(value)))
}
func (e *encoder) AddInt(key string, value int)        { e.addAttr(slog.Int(key, value)) }
func (e *encoder) AddInt64(key string, value int64)    { e.addAttr(slog.Int64(key, value)) }
func (e *encoder) AddInt32(key string, value int32)    { e.addAttr(slog.Int64(key, int64(value))) }
func (e *encoder) AddInt16(key string, value int16)    { e.addAttr(slog.Int64(key, int64(value))) }
func (e *encoder) AddInt8(key string, value int8)      { e.addAttr(slog.Int64(key, int64(value))) }
func (e *encoder) AddString(key, value string)         { e.addAttr(slog.String(key, value)) }
func (e *encoder) AddTime(key string, value time.Time) { e.addAttr(slog.Time(key, value)) }
func (e *encoder) AddUint(key string, value uint)      { e.addAttr(slog.Uint64(key, uint64(value))) }
func (e *encoder) AddUint64(key string, value uint64)  { e.addAttr(slog.Uint64(key, value)) }
func (e *encoder) AddUint32(key string, value uint32) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *encoder) AddUint16(key string, value uint16) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}
func (e *encoder) AddUint8(key string, value uint8) { e.addAttr(slog.Uint64(key, uint64(value))) }
func (e *encoder) AddUintptr(key string, value uintptr) {
	e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *encoder) AddReflected(key string, value any) error {
	e.addAttr(slog.Any(key, value))
	return nil
}

// OpenNamespace switches subsequent writes into a child node, matching zap's
// namespace semantics where all following fields belong to that group.
func (e *encoder) OpenNamespace(key string) {
	ns := &node{}
	e.cur.entries = append(e.cur.entries, objectEntry{
		child: &groupEntry{
			key:  key,
			node: ns,
		},
	})
	e.cur = ns
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
	obj := newEncoder(0)
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
