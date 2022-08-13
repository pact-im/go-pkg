// Package encode provides value encoders for [go.pact.im/x/phcformat] package.
package encode

import (
	"encoding/base64"
	"strconv"

	"go.pact.im/x/phcformat/option"
)

// Appender represents an encodable value that uses append-style API.
type Appender interface {
	// Append appends the encoded value to the dst and returns the resulting
	// slice.
	Append(dst []byte) []byte
}

// Nil returns a nil Appender option.
func Nil() option.Of[Appender] {
	return option.Nil[Appender]()
}

// StringOrBytes is a union of string and byte slice types.
type StringOrBytes interface {
	~string | ~[]byte
}

// Empty is an Appender that does nothing.
type Empty struct{}

// NewEmpty returns a new Empty instance.
func NewEmpty() Empty {
	return Empty{}
}

// Append implements the Appender interface.
func (v Empty) Append(dst []byte) []byte {
	return dst
}

// Option is an Appender that appends optional value if it is set.
type Option[T Appender] struct {
	Option option.Of[T]
}

// NewOption returns a new Option instance.
func NewOption[T Appender](opt option.Of[T]) Option[T] {
	return Option[T]{
		Option: opt,
	}
}

// Append implements the Appender interface.
func (v Option[T]) Append(dst []byte) []byte {
	value, ok := v.Option.Unwrap()
	if !ok {
		return dst
	}
	return value.Append(dst)
}

// Concat is an Appender that concatenates two values.
type Concat[T, U Appender] struct {
	// A is the first value to append.
	A T
	// B is the second value to append.
	B U
}

// NewConcat returns a new Concat instance.
func NewConcat[T, U Appender](a T, b U) Concat[T, U] {
	return Concat[T, U]{
		A: a,
		B: b,
	}
}

// Append implements the Appender interface.
func (v Concat[T, U]) Append(dst []byte) []byte {
	dst = v.A.Append(dst)
	dst = v.B.Append(dst)
	return dst
}

// List is an Appender that appends a list of values separated with the given
// separator.
type List[SeparatorAppender, ElementAppender Appender] struct {
	// Separator is a separator that is appended between elements.
	Separator SeparatorAppender
	// Elements is a list of elements to append.
	Elements []ElementAppender
}

// NewList returns a new List instance.
func NewList[SeparatorAppender, ElementAppender Appender](sep SeparatorAppender, elements ...ElementAppender) List[SeparatorAppender, ElementAppender] {
	return List[SeparatorAppender, ElementAppender]{
		Separator: sep,
		Elements:  elements,
	}
}

// Append implements the Appender interface.
func (v List[ElementAppender, SeparatorAppender]) Append(dst []byte) []byte {
	if len(v.Elements) == 0 {
		return dst
	}
	dst = v.Elements[0].Append(dst)
	for _, e := range v.Elements[1:] {
		dst = v.Separator.Append(dst)
		dst = e.Append(dst)
	}
	return dst
}

// KV is an Appender that appends a key-value pair separated with the given
// separator.
type KV[KeyAppender, SepAppender, ValAppender Appender] struct {
	// Key is the key in the key-value pair.
	Key KeyAppender
	// Sep is a separator that is appended between key and value.
	Sep SepAppender
	// Val is the value in the key-value pair.
	Val ValAppender
}

// NewKV returns a new KV instance.
func NewKV[KeyAppender, SepAppender, ValAppender Appender](sep SepAppender, k KeyAppender, v ValAppender) KV[KeyAppender, SepAppender, ValAppender] {
	return KV[KeyAppender, SepAppender, ValAppender]{
		Key: k,
		Sep: sep,
		Val: v,
	}
}

// Append implements the Appender interface.
func (v KV[KeyAppender, SepAppender, ValAppender]) Append(dst []byte) []byte {
	dst = v.Key.Append(dst)
	dst = v.Sep.Append(dst)
	dst = v.Val.Append(dst)
	return dst
}

// Byte is an Appender that appends a single byte.
type Byte byte

// NewByte returns a new Byte instance.
func NewByte(c byte) Byte {
	return Byte(c)
}

// Append implements the Appender interface.
func (v Byte) Append(dst []byte) []byte {
	return append(dst, byte(v))
}

// String is an Appender that appends string.
type String string

// NewString returns a new String instance.
func NewString(s string) String {
	return String(s)
}

// Append implements the Appender interface.
func (v String) Append(dst []byte) []byte {
	return append(dst, v...)
}

// Bytes is an Appender that appends byte slice.
type Bytes []byte

// NewBytes returns a new Bytes instance.
func NewBytes(buf []byte) Bytes {
	return Bytes(buf)
}

// Append implements the Appender interface.
func (v Bytes) Append(dst []byte) []byte {
	return append(dst, v...)
}

// Uint is an Appender that encodes uint as a decimal number.
type Uint uint

// NewUint returns a new Uint instance.
func NewUint(n uint) Uint {
	return Uint(n)
}

// Append implements the Appender interface.
func (v Uint) Append(dst []byte) []byte {
	return strconv.AppendUint(dst, uint64(v), 10)
}

// Base64 is an Appender that encodes string or byte slice using
// base64.RawStdEncoding.
type Base64[T StringOrBytes] struct {
	// Data is the unencoded data.
	Data T
}

// NewBase64 returns a new Base64 instance.
func NewBase64[T StringOrBytes](data T) Base64[T] {
	return Base64[T]{
		Data: data,
	}
}

// Append implements the Appender interface.
func (v Base64[T]) Append(dst []byte) []byte {
	return base64Append(base64.RawStdEncoding, dst, []byte(v.Data))
}

// base64Append is an append-style function for base64 encoding.
//
// See also https://go.dev/issue/19366
func base64Append(e *base64.Encoding, dst, src []byte) []byte {
	k := len(dst)
	n := e.EncodedLen(len(src))
	if cap(dst)-k < n {
		dst = append(dst, make([]byte, n)...)
	} else {
		dst = dst[:k+n]
	}
	e.Encode(dst[k:], src)
	return dst
}
