// Package phcformat implements PHC string format parser and encoder.
//
// See https://github.com/P-H-C/phc-string-format
package phcformat

import "encoding/base64"

// Hash represents a hash function in PHC string format.
type Hash struct {
	// ID is the symbolic name for the function. The function symbolic name
	// is a sequence of characters in: [a-z0-9-] (lowercase letters, digits,
	// and the minus sign). No other characters are allowed. A name must not
	// exceed 32 characters in length. Note that it is allowed to be empty.
	ID string

	// Version is the algorithm version. The function version is a sequence
	// of characters in: [0-9]. Note that it is allowed to be empty.
	// prohibited)
	Version OptionalString

	// Params is a comma-separated sequence of "key=value" parameter pairs.
	Params OptionalString

	// Salt is an encoded salt string.
	Salt OptionalString

	// Output is the base64-encoded function output.
	Output OptionalString

	// Raw is the unparsed hash in PHC string format.
	Raw string
}

// String implements the fmt.Stringer interface. It returns the hash in PHC
// string format.
func (h Hash) String() string {
	return h.Raw
}

// HashSaltType represents the underlying type of the HashSalt.
type HashSaltType byte

const (
	// HashSaltTypeUnset indicates that the salt is not set.
	HashSaltTypeUnset HashSaltType = iota
	// HashSaltTypeString indicates that the salt is a string.
	HashSaltTypeString
	// HashSaltTypeBytes indicates that the salt is a bytes slice.
	HashSaltTypeBytes
)

// HashSaltFormat represents a supported salt encoding.
type HashSaltFormat byte

const (
	// HashSaltFormatEncoded indicates that the salt is already in encoded
	// form. Encode will only validate that the characters are valid and
	// append it without modifications.
	HashSaltFormatEncoded HashSaltFormat = iota
	// HashSaltFormatBase64 indicates that the salt should be encoded using
	// standard base64 encoding without padding characters.
	HashSaltFormatBase64
)

// HashSalt is an unencoded hash function salt used in Encode.
type HashSalt struct {
	// Format is an encoding format to use for salt. When using base64
	// format, it is recommended to set Bytes field instead of String to
	// avoid allocations.
	Format HashSaltFormat
	// String is an unencoded salt represented as string. If String is
	// empty, Bytes field is used.
	String string
	// Bytes is an unencoded salt represented as byte slice. If the slice is
	// nil, salt is not set.
	Bytes []byte
}

// Type returns the underlying salt type.
func (h HashSalt) Type() HashSaltType {
	if h.String != "" {
		return HashSaltTypeString
	}
	if h.Bytes != nil {
		return HashSaltTypeBytes
	}
	return HashSaltTypeUnset
}

// Len returns the length of the underlying unencoded string or byte slice.
func (h HashSalt) Len() int {
	n := len(h.String)
	if n != 0 {
		return n
	}
	return len(h.Bytes)
}

// EncodedLen returns the length of the encoded underlying string or byte slice.
func (h HashSalt) EncodedLen() int {
	n := h.Len()
	if h.Format == HashSaltFormatBase64 {
		n = base64.RawStdEncoding.EncodedLen(n)
	}
	return n
}

// Encode copies the encoded salt in the specified format to dst.
func (h HashSalt) Encode(dst []byte) {
	switch h.Type() {
	case HashSaltTypeString:
		switch h.Format {
		case HashSaltFormatEncoded:
			copy(dst, h.String)
		case HashSaltFormatBase64:
			base64.RawStdEncoding.Encode(dst, []byte(h.String))
		}
	case HashSaltTypeBytes:
		switch h.Format {
		case HashSaltFormatEncoded:
			copy(dst, h.Bytes)
		case HashSaltFormatBase64:
			base64.RawStdEncoding.Encode(dst, h.Bytes)
		}
	}
}

// Valid checks whether the salt contains valid characters for the specified
// encoding.
func (h HashSalt) Valid() bool {
	switch h.Type() {
	case HashSaltTypeString:
		if h.Format != HashSaltFormatEncoded {
			break
		}
		for i := 0; i < len(h.String); i++ {
			if validSalt(h.String[i]) {
				continue
			}
			return false
		}
	case HashSaltTypeBytes:
		if h.Format != HashSaltFormatEncoded {
			break
		}
		for i := 0; i < len(h.Bytes); i++ {
			if validSalt(h.Bytes[i]) {
				continue
			}
			return false
		}
	}
	return true
}

// IsSet returns true if the salt is set.
func (h HashSalt) IsSet() bool {
	return h.Type() != HashSaltTypeUnset
}

// HashParam is a hash function parameter key and value pair.
type HashParam struct {
	// Name is the name of the parameter.
	Name string
	// Value is the value of the parameter.
	Value string
}

func validID(c byte) bool {
	return 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '-'
}

func validVersion(c byte) bool {
	return '0' <= c && c <= '9'
}

func validParamName(c byte) bool {
	return 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '-'
}

func validParamValue(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' || c == '/' || c == '+' || c == '.' || c == '-'
}

func validSalt(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' || c == '/' || c == '+' || c == '.' || c == '-'
}

func validOutput(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '+' || c == '/'
}
