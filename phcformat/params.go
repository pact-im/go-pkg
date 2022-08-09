package phcformat

import "strings"

// ParamsIterator iterates over comma-separated key=value parameter pairs. Note
// that iterator does not validate characters in parameter’s key and value.
type ParamsIterator struct {
	// Param is the current parameter.
	Param HashParam
	// After is the string with remaining parameters.
	After string
	// Valid indicates that the iterator is valid.
	Valid bool
}

// IterParams returns ParamsIterator for an optional string.
func IterParams(s OptionalString) ParamsIterator {
	if !s.IsSet {
		return ParamsIterator{}
	}
	it := ParamsIterator{After: s.Value}
	return it.Next()
}

// Collect collects all parameters by repeatedly calling Next.
func (it ParamsIterator) Collect() []HashParam {
	var params []HashParam
	for ; it.Valid; it = it.Next() {
		params = append(params, it.Param)
	}
	return params
}

// Next advances to the next parameter in the sequence.
func (it ParamsIterator) Next() ParamsIterator {
	it.Param, it.After, it.Valid = nextParam(it.After)
	return it
}

// nextParam returns the next parameter in s and the remaining string.
func nextParam(s string) (HashParam, string, bool) {
	var name, value string

	i := strings.IndexByte(s, '=')
	if i < 0 {
		return HashParam{}, "", false
	}
	name, s = s[:i], s[i+1:]

	j := strings.IndexByte(s, ',')
	if j < 0 {
		return HashParam{name, s}, "", true
	}
	value, s = s[:j], s[j+1:]

	return HashParam{name, value}, s, true
}
