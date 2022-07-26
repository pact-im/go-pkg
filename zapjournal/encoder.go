package zapjournal

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

// hdrLen is the size of the length header for a multi-line variable. The header
// is a single 64-bit little endian integer.
const hdrLen = 8

// Severity levels used in PRIORITY field for syslog compatibility.
const (
	priEmerg = string('0' + iota)
	priAlert
	priCrit
	priErr
	priWarning
	priNotice
	priInfo
	priDebug
)

// varsEncoder is a partial zapcore.Encoder implementation that encodes entries
// to the systemd-journald wire format. Top-level fields are set as variables
// with the configured prefix with open namespaces appended, while composite
// values (arrays and objects) and reflected values are encoded as JSON.
//
// Note that, in JSON mode, numeric values (float, complex, int) are encoded as
// strings. This addresses the limitations of JSON encoding, i.e. the lack of
// complex numbers, NaN and Inf, and integer interoperability. This does not
// apply to the reflected values though, that use standard library for encoding.
type varsEncoder struct {
	// prefix is appended to variable names.
	prefix string

	// buf holds the encoded variables.
	buf []byte

	// hdr, if non zero, is an index to buf that is used in endVar to write
	// the length of the data appended since beginVar.
	hdr int

	// json encoder is used to encode composite values (objects and arrays).
	json jsonEncoder
}

// beginSinglelineVar starts a single-line variable.
func (e *varsEncoder) beginSinglelineVar(name string) {
	e.beginVar(name, false)
}

// beginMultilineVar starts a multi-line variable.
func (e *varsEncoder) beginMultilineVar(name string) {
	e.beginVar(name, true)
}

// beginVar starts a single-line or multi-line variable. It should not be used
// directly. Instead, call the convenience functions beginSinglelineVar and
// beginMultilineVar. Each call should have a corresponding endVar. Nested
// variables are not allowed.
func (e *varsEncoder) beginVar(name string, multiline bool) {
	e.buf = appendVarName(e.buf, e.prefix, name)
	if !multiline {
		e.buf = append(e.buf, '=')
		return
	}
	e.buf = append(e.buf, '\n')
	e.hdr = len(e.buf)
	e.buf = append(e.buf, make([]byte, hdrLen)...)
}

// endVar ends a variable. Both single-line and multi-line variables are
// supported.
func (e *varsEncoder) endVar() {
	if e.hdr != 0 {
		size := uint64(len(e.buf) - e.hdr - hdrLen)
		binary.LittleEndian.PutUint64(e.buf[e.hdr:], size)
		e.hdr = 0
	}
	e.buf = append(e.buf, '\n')
}

func (e *varsEncoder) OpenNamespace(key string) {
	e.prefix += "_" + key
}

func (e *varsEncoder) AddBinary(key string, value []byte) {
	e.beginMultilineVar(key)
	e.buf = append(e.buf, value...)
	e.endVar()
}

func (e *varsEncoder) AddByteString(key string, value []byte) {
	multiline := bytes.Contains(value, []byte{'\n'})
	e.beginVar(key, multiline)
	e.buf = append(e.buf, value...)
	e.endVar()
}

func (e *varsEncoder) AddString(key, value string) {
	multiline := strings.Contains(value, "\n")
	e.beginVar(key, multiline)
	e.buf = append(e.buf, value...)
	e.endVar()
}

func (e *varsEncoder) AddBool(key string, value bool) {
	e.beginSinglelineVar(key)
	e.buf = strconv.AppendBool(e.buf, value)
	e.endVar()
}

func (e *varsEncoder) AddComplex128(key string, value complex128) {
	e.addComplex(key, value, 128)
}

func (e *varsEncoder) AddComplex64(key string, value complex64) {
	e.addComplex(key, complex128(value), 64)
}

func (e *varsEncoder) addComplex(key string, value complex128, bitSize int) {
	e.beginSinglelineVar(key)
	e.buf = strconvAppendComplex(e.buf, value, 'g', -1, bitSize)
	e.endVar()
}

func (e *varsEncoder) AddDuration(key string, value time.Duration) {
	e.AddString(key, value.String())
}

func (e *varsEncoder) AddFloat64(key string, value float64) {
	e.addFloat(key, value, 64)
}

func (e *varsEncoder) AddFloat32(key string, value float32) {
	e.addFloat(key, float64(value), 32)
}

func (e *varsEncoder) addFloat(key string, value float64, bitSize int) {
	e.beginSinglelineVar(key)
	e.buf = strconv.AppendFloat(e.buf, value, 'g', -1, bitSize)
	e.endVar()
}

func (e *varsEncoder) AddInt(key string, value int) {
	e.AddInt64(key, int64(value))
}

func (e *varsEncoder) AddInt64(key string, value int64) {
	e.beginSinglelineVar(key)
	e.buf = strconv.AppendInt(e.buf, value, 10)
	e.endVar()
}

func (e *varsEncoder) AddInt32(key string, value int32) {
	e.AddInt64(key, int64(value))
}

func (e *varsEncoder) AddInt16(key string, value int16) {
	e.AddInt64(key, int64(value))
}

func (e *varsEncoder) AddInt8(key string, value int8) {
	e.AddInt64(key, int64(value))
}

func (e *varsEncoder) AddUint(key string, value uint) {
	e.AddUint64(key, uint64(value))
}

func (e *varsEncoder) AddUint64(key string, value uint64) {
	e.beginSinglelineVar(key)
	e.buf = strconv.AppendUint(e.buf, value, 10)
	e.endVar()
}

func (e *varsEncoder) AddUint32(key string, value uint32) {
	e.AddUint64(key, uint64(value))
}

func (e *varsEncoder) AddUint16(key string, value uint16) {
	e.AddUint64(key, uint64(value))
}

func (e *varsEncoder) AddUint8(key string, value uint8) {
	e.AddUint64(key, uint64(value))
}

func (e *varsEncoder) AddUintptr(key string, value uintptr) {
	e.beginSinglelineVar(key)
	e.buf = strconvAppendUintptr(e.buf, value)
	e.endVar()
}

func (e *varsEncoder) AddTime(key string, value time.Time) {
	e.beginSinglelineVar(key)
	e.buf = value.AppendFormat(e.buf, time.RFC3339Nano)
	e.endVar()
}

func (e *varsEncoder) AddReflected(key string, value any) error {
	e.beginMultilineVar(key)
	buf := bytes.NewBuffer(e.buf)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(value)
	e.buf = buf.Bytes()
	// Remove newline character that is always added by Encode.
	// See https://stackoverflow.com/a/36320146
	e.buf = e.buf[:len(e.buf)-1]
	e.endVar()
	return err
}

func (e *varsEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	e.beginMultilineVar(key)
	e.json.appendOpenArray()
	err := marshaler.MarshalLogArray(&e.json)
	e.json.complete()
	e.json.appendCloseArray()
	e.endVar()
	return err
}

func (e *varsEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	e.beginMultilineVar(key)
	e.json.appendOpenObject()
	err := marshaler.MarshalLogObject(&e.json)
	e.json.complete()
	e.json.appendCloseObject()
	e.endVar()
	return err
}

func (e *varsEncoder) encodeEntryVars(ent zapcore.Entry) {
	// Temporarily unset prefix to add entry variables.
	tmp := e.prefix
	e.prefix = ""

	// https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html#MESSAGE=
	e.AddString("MESSAGE", ent.Message)

	// https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html#CODE_FILE=
	if caller := ent.Caller; caller.Defined {
		e.AddString("CODE_FILE", caller.File)
		e.AddInt("CODE_LINE", caller.Line)
		e.AddString("CODE_FUNC", caller.Function)
	}

	// https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html#PRIORITY=
	pri := priNotice
	switch ent.Level {
	case zapcore.DebugLevel:
		pri = priDebug
	case zapcore.InfoLevel:
		pri = priInfo
	case zapcore.WarnLevel:
		pri = priWarning
	case zapcore.ErrorLevel:
		pri = priErr
	case zapcore.DPanicLevel:
		pri = priCrit
	case zapcore.PanicLevel:
		pri = priAlert
	case zapcore.FatalLevel:
		pri = priEmerg
	}
	e.AddString("PRIORITY", pri)

	// Also add textual representation of the log level.
	e.AddString("LOG_LEVEL", ent.Level.String())

	if !ent.Time.IsZero() {
		e.AddTime("TIMESTAMP", ent.Time)
	}

	if ent.LoggerName != "" {
		e.AddString("LOG_NAME", ent.LoggerName)
	}

	if ent.Stack != "" {
		e.AddString("STACK", ent.Stack)
	}

	e.prefix = tmp
}

// encodeEntry encodes the entry using the underlying buffer as the initial
// state. Note that it copies the internal buffer and does not mutate it. The
// resulting encoder is owned by the caller.
func (e *varsEncoder) encodeEntry(ent zapcore.Entry, fields []zapcore.Field) *varsEncoder {
	e = cloneVarsEncoder(e)
	addFields(e, fields)
	e.encodeEntryVars(ent)
	return e
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
