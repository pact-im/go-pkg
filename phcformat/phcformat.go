// Package phcformat implements PHC string format parser and encoder.
//
// See https://github.com/P-H-C/phc-string-format
package phcformat

import (
	"go.pact.im/x/option"
)

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
	Version option.Of[string]

	// Params is a comma-separated sequence of "key=value" parameter pairs.
	Params option.Of[string]

	// Salt is an encoded salt string.
	Salt option.Of[string]

	// Output is the base64-encoded function output.
	Output option.Of[string]

	// Raw is the unparsed hash in PHC string format.
	Raw string
}

// String implements the fmt.Stringer interface. It returns the hash in PHC
// string format.
func (h Hash) String() string {
	return string(h.Raw)
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
