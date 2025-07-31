package zapjournal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"
)

type jsonEncoder struct {
	bufp *[]byte
	jsonState
}

type jsonState struct {
	depth int
	prev  bool
}

func (e *jsonEncoder) buf() []byte {
	return *e.bufp
}

func (e *jsonEncoder) setBuf(buf []byte) {
	*e.bufp = buf
}

func (e *jsonEncoder) nested() jsonState {
	old := e.jsonState
	e.jsonState = jsonState{}
	return old
}

func (e *jsonEncoder) appendKey(name string) {
	e.appendSep()
	e.appendQuotedString(name)
	e.appendColon()
}

func (e *jsonEncoder) appendQuotedString(s string) {
	e.setBuf(strconv.AppendQuote(e.buf(), s))
}

func (e *jsonEncoder) appendInt(v int64) {
	e.appendQuote()
	e.setBuf(strconv.AppendInt(e.buf(), v, 10))
	e.appendQuote()
}

func (e *jsonEncoder) appendUint(v uint64) {
	e.appendQuote()
	e.setBuf(strconv.AppendUint(e.buf(), v, 10))
	e.appendQuote()
}

func (e *jsonEncoder) appendUintptr(v uintptr) {
	e.appendQuote()
	e.setBuf(append(e.buf(), "0x"...))
	e.setBuf(strconv.AppendUint(e.buf(), uint64(v), 16))
	e.appendQuote()
}

func (e *jsonEncoder) appendFloat(v float64, bitSize int) {
	e.appendQuote()
	e.setBuf(strconv.AppendFloat(e.buf(), v, 'g', -1, bitSize))
	e.appendQuote()
}

func (e *jsonEncoder) appendComplex(v complex128, bitSize int) {
	e.appendQuote()
	e.setBuf(strconvAppendComplex(e.buf(), v, 'g', -1, bitSize))
	e.appendQuote()
}

func (e *jsonEncoder) appendMarshalJSON(v any) error {
	buf := bytes.NewBuffer(e.buf())
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	e.setBuf(buf.Bytes())
	return err
}

func (e *jsonEncoder) appendByte(c byte) {
	e.setBuf(append(e.buf(), c))
}

func (e *jsonEncoder) appendColon() {
	e.appendByte(':')
}

func (e *jsonEncoder) appendQuote() {
	e.appendByte('"')
}

func (e *jsonEncoder) appendOpenArray() {
	e.appendByte('[')
}

func (e *jsonEncoder) appendCloseArray() {
	e.appendByte(']')
}

func (e *jsonEncoder) appendOpenObject() {
	e.appendByte('{')
}

func (e *jsonEncoder) appendCloseObject() {
	e.appendByte('}')
}

func (e *jsonEncoder) appendSep() {
	if !e.prev {
		e.prev = true
		return
	}
	e.setBuf(append(e.buf(), ','))
}

func (e *jsonEncoder) complete() {
	for ; e.depth > 0; e.depth-- {
		e.appendCloseObject()
	}
	e.prev = false
}

//
// Binary
//

func (e *jsonEncoder) AddBinary(k string, v []byte) {
	e.appendKey(k)
	e.encodeBinary(v)
}

func (e *jsonEncoder) encodeBinary(v []byte) {
	e.setBuf(base64Append(base64.StdEncoding, e.buf(), v))
}

//
// Duration
//

func (e *jsonEncoder) AddDuration(k string, v time.Duration) {
	e.appendKey(k)
	e.encodeDuration(v)
}

func (e *jsonEncoder) AppendDuration(v time.Duration) {
	e.appendSep()
	e.encodeDuration(v)
}

func (e *jsonEncoder) encodeDuration(v time.Duration) {
	e.appendQuote()
	e.setBuf(append(e.buf(), v.String()...))
	e.appendQuote()
}

//
// Time
//

func (e *jsonEncoder) AddTime(k string, v time.Time) {
	e.appendKey(k)
	e.encodeTime(v)
}

func (e *jsonEncoder) AppendTime(v time.Time) {
	e.appendSep()
	e.encodeTime(v)
}

func (e *jsonEncoder) encodeTime(v time.Time) {
	e.appendQuote()
	e.setBuf(v.AppendFormat(e.buf(), time.RFC3339Nano))
	e.appendQuote()
}

//
// Bool
//

func (e *jsonEncoder) AddBool(k string, v bool) {
	e.appendKey(k)
	e.encodeBool(v)
}

func (e *jsonEncoder) AppendBool(v bool) {
	e.appendSep()
	e.encodeBool(v)
}

func (e *jsonEncoder) encodeBool(v bool) {
	e.setBuf(strconv.AppendBool(e.buf(), v))
}

//
// ByteString
//

func (e *jsonEncoder) AddByteString(k string, v []byte) {
	e.appendKey(k)
	e.encodeByteString(v)
}

func (e *jsonEncoder) AppendByteString(v []byte) {
	e.appendSep()
	e.encodeByteString(v)
}

func (e *jsonEncoder) encodeByteString(v []byte) {
	e.appendQuotedString(string(v))
}

//
// String
//

func (e *jsonEncoder) AddString(k, v string) {
	e.appendKey(k)
	e.encodeString(v)
}

func (e *jsonEncoder) AppendString(v string) {
	e.appendSep()
	e.encodeString(v)
}

func (e *jsonEncoder) encodeString(v string) {
	e.appendQuotedString(v)
}

// Complex128

func (e *jsonEncoder) AddComplex128(k string, v complex128) {
	e.appendKey(k)
	e.encodeComplex128(v)
}

func (e *jsonEncoder) AppendComplex128(v complex128) {
	e.appendSep()
	e.encodeComplex128(v)
}

func (e *jsonEncoder) encodeComplex128(v complex128) {
	e.appendComplex(v, 128)
}

//
// Complex64
//

func (e *jsonEncoder) AddComplex64(k string, v complex64) {
	e.appendKey(k)
	e.encodeComplex64(v)
}

func (e *jsonEncoder) AppendComplex64(v complex64) {
	e.appendSep()
	e.encodeComplex64(v)
}

func (e *jsonEncoder) encodeComplex64(v complex64) {
	e.appendComplex(complex128(v), 64)
}

//
// Float64
//

func (e *jsonEncoder) AddFloat64(k string, v float64) {
	e.appendKey(k)
	e.encodeFloat64(v)
}

func (e *jsonEncoder) AppendFloat64(v float64) {
	e.appendSep()
	e.encodeFloat64(v)
}

func (e *jsonEncoder) encodeFloat64(v float64) {
	e.appendFloat(v, 64)
}

//
// Float32
//

func (e *jsonEncoder) AddFloat32(k string, v float32) {
	e.appendKey(k)
	e.encodeFloat32(v)
}

func (e *jsonEncoder) AppendFloat32(v float32) {
	e.appendSep()
	e.encodeFloat32(v)
}

func (e *jsonEncoder) encodeFloat32(v float32) {
	e.appendFloat(float64(v), 32)
}

//
// Int
//

func (e *jsonEncoder) AddInt(k string, v int) {
	e.appendKey(k)
	e.encodeInt(v)
}

func (e *jsonEncoder) AppendInt(v int) {
	e.appendSep()
	e.encodeInt(v)
}

func (e *jsonEncoder) encodeInt(v int) {
	e.appendInt(int64(v))
}

//
// Int64
//

func (e *jsonEncoder) AddInt64(k string, v int64) {
	e.appendKey(k)
	e.encodeInt64(v)
}

func (e *jsonEncoder) AppendInt64(v int64) {
	e.appendSep()
	e.encodeInt64(v)
}

func (e *jsonEncoder) encodeInt64(v int64) {
	e.appendInt(v)
}

//
// Int32
//

func (e *jsonEncoder) AddInt32(k string, v int32) {
	e.appendKey(k)
	e.encodeInt32(v)
}

func (e *jsonEncoder) AppendInt32(v int32) {
	e.appendSep()
	e.encodeInt32(v)
}

func (e *jsonEncoder) encodeInt32(v int32) {
	e.appendInt(int64(v))
}

//
// Int16
//

func (e *jsonEncoder) AddInt16(k string, v int16) {
	e.appendKey(k)
	e.encodeInt16(v)
}

func (e *jsonEncoder) AppendInt16(v int16) {
	e.appendSep()
	e.encodeInt16(v)
}

func (e *jsonEncoder) encodeInt16(v int16) {
	e.appendInt(int64(v))
}

//
// Int8
//

func (e *jsonEncoder) AddInt8(k string, v int8) {
	e.appendKey(k)
	e.encodeInt8(v)
}

func (e *jsonEncoder) AppendInt8(v int8) {
	e.appendSep()
	e.encodeInt8(v)
}

func (e *jsonEncoder) encodeInt8(v int8) {
	e.appendInt(int64(v))
}

//
// Uint
//

func (e *jsonEncoder) AddUint(k string, v uint) {
	e.appendKey(k)
	e.encodeUint(v)
}

func (e *jsonEncoder) AppendUint(v uint) {
	e.appendSep()
	e.encodeUint(v)
}

func (e *jsonEncoder) encodeUint(v uint) {
	e.appendUint(uint64(v))
}

//
// Uint64
//

func (e *jsonEncoder) AddUint64(k string, v uint64) {
	e.appendKey(k)
	e.encodeUint64(v)
}

func (e *jsonEncoder) AppendUint64(v uint64) {
	e.appendSep()
	e.encodeUint64(v)
}

func (e *jsonEncoder) encodeUint64(v uint64) {
	e.appendUint(v)
}

//
// Uint32
//

func (e *jsonEncoder) AddUint32(k string, v uint32) {
	e.appendKey(k)
	e.encodeUint32(v)
}

func (e *jsonEncoder) AppendUint32(v uint32) {
	e.appendSep()
	e.encodeUint32(v)
}

func (e *jsonEncoder) encodeUint32(v uint32) {
	e.appendUint(uint64(v))
}

//
// Uint16
//

func (e *jsonEncoder) AddUint16(k string, v uint16) {
	e.appendKey(k)
	e.encodeUint16(v)
}

func (e *jsonEncoder) AppendUint16(v uint16) {
	e.appendSep()
	e.encodeUint16(v)
}

func (e *jsonEncoder) encodeUint16(v uint16) {
	e.appendUint(uint64(v))
}

//
// Uint8
//

func (e *jsonEncoder) AddUint8(k string, v uint8) {
	e.appendKey(k)
	e.encodeUint8(v)
}

func (e *jsonEncoder) AppendUint8(v uint8) {
	e.appendSep()
	e.encodeUint8(v)
}

func (e *jsonEncoder) encodeUint8(v uint8) {
	e.appendUint(uint64(v))
}

//
// Uintptr
//

func (e *jsonEncoder) AddUintptr(k string, v uintptr) {
	e.appendKey(k)
	e.encodeUintptr(v)
}

func (e *jsonEncoder) AppendUintptr(v uintptr) {
	e.appendSep()
	e.encodeUintptr(v)
}

func (e *jsonEncoder) encodeUintptr(v uintptr) {
	e.appendUintptr(v)
}

//
// Reflected
//

func (e *jsonEncoder) AddReflected(k string, v any) error {
	e.appendKey(k)
	return e.encodeReflected(v)
}

func (e *jsonEncoder) AppendReflected(v any) error {
	e.appendSep()
	return e.encodeReflected(v)
}

func (e *jsonEncoder) encodeReflected(v any) error {
	return e.appendMarshalJSON(v)
}

//
// Array
//

func (e *jsonEncoder) AddArray(k string, v zapcore.ArrayMarshaler) error {
	e.appendKey(k)
	return e.encodeArray(v)
}

func (e *jsonEncoder) AppendArray(v zapcore.ArrayMarshaler) error {
	e.appendSep()
	return e.encodeArray(v)
}

func (e *jsonEncoder) encodeArray(v zapcore.ArrayMarshaler) error {
	e.appendOpenArray()
	old := e.nested()
	err := v.MarshalLogArray(e)
	e.complete()
	e.jsonState = old
	e.appendCloseArray()
	return err
}

//
// Object
//

func (e *jsonEncoder) AddObject(k string, v zapcore.ObjectMarshaler) error {
	e.appendKey(k)
	return e.encodeObject(v)
}

func (e *jsonEncoder) AppendObject(v zapcore.ObjectMarshaler) error {
	e.appendSep()
	return e.encodeObject(v)
}

func (e *jsonEncoder) encodeObject(v zapcore.ObjectMarshaler) error {
	e.appendOpenObject()
	old := e.nested()
	err := v.MarshalLogObject(e)
	e.complete()
	e.jsonState = old
	e.appendCloseObject()
	return err
}

//
// Namespace
//

func (e *jsonEncoder) OpenNamespace(k string) {
	e.appendKey(k)
	e.appendOpenObject()
	e.prev = false
	e.depth++
}
