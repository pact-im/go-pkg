package phcformat

import (
	"bytes"
	"strings"
)

// ParamsIterator iterates over comma-separated key=value parameter pairs. Note
// that iterator does not validate characters in parameter’s key and value.
type ParamsIterator struct {
	// Name is the name of the current parameter.
	Name string
	// Value is the value of the current parameter.
	Value string
	// After is the string with remaining parameters.
	After string
	// Valid indicates that the iterator is valid. On parse error, After
	// will contain unparsed bytes.
	Valid bool
}

// IterParams returns a ParamsIterator for the given string.
func IterParams(s string) ParamsIterator {
	if s == "" {
		return ParamsIterator{}
	}
	it := ParamsIterator{After: s}
	return it.Next()
}

// Next advances to the next parameter in the sequence.
func (it ParamsIterator) Next() ParamsIterator {
	it.Name, it.Value, it.After, it.Valid = nextParam(it.After)
	return it
}

// nextParam returns the next parameter in s and the remaining string.
func nextParam(s string) (name, value, after string, ok bool) {
	i := strings.IndexByte(s, '=')
	if i < 0 {
		return "", "", s, false
	}
	name, s = s[:i], s[i+1:]

	j := bytes.IndexByte([]byte(s), ',')
	if j < 0 {
		return name, s, "", true
	}

	// Consume parameter but make the next iteration invalid if we have a
	// trailing comma.
	off := 1
	if j+1 == len(s) {
		off = 0
	}

	value, s = s[:j], s[j+off:]

	return name, value, s, true
}
