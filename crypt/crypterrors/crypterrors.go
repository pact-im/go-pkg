// Package crypterrors defines errors that can be returned from crypt.Crypter
// implementations.
package crypterrors

import (
	"fmt"
)

// MalformedHashError is an error that is returned if the crypt.Crypter
// implementation receives a malformed PHC formatted string.
type MalformedHashError struct {
	// Hash is the malformed PHC formatted hash string.
	Hash string
}

// Error implements the error interface.
func (e *MalformedHashError) Error() string {
	const m = "malformed hash"
	if e == nil {
		return m
	}
	return fmt.Sprintf(m+" %q", e.Hash)
}

// UnsupportedHashError is an error that is returned if the crypt.Crypter
// implementation does not support the requested hash function.
type UnsupportedHashError struct {
	// HashID is the ID of the unsupported hash function.
	HashID string
}

// Error implements the error interface.
func (e *UnsupportedHashError) Error() string {
	const m = "unsupported hash"
	if e == nil {
		return m
	}
	return fmt.Sprintf(m+" %q", e.HashID)
}

// UnsupportedVersionError is an error that is returned if the crypt.Crypter
// implementation does not support the requested hash function version.
type UnsupportedVersionError struct {
	// Parsed is the parsed unsupported version.
	Parsed uint
	// Suggested is the suggested supported version.
	Suggested uint
}

// Error implements the error interface.
func (e *UnsupportedVersionError) Error() string {
	const m = "unsupported version"
	if e == nil {
		return m
	}
	return fmt.Sprintf(m+" %d (suggested version is %d)", e.Parsed, e.Suggested)
}

// UnsupportedParameterError is an error that is returned if the crypt.Crypter
// implementation does not support the given parameter for a hash function.
type UnsupportedParameterError struct {
	// Name is the name of the parameter.
	Name string
	// Unimplemented indicates that the parameter is defined but not
	// implemented.
	Unimplemented bool
}

// Error implements the error interface.
func (e *UnsupportedParameterError) Error() string {
	const m = "unsupported parameter"
	if e == nil {
		return m
	}
	var suffix string
	if e.Unimplemented {
		suffix = " (not implemented)"
	}
	return fmt.Sprintf(m+" %q%s", e.Name, suffix)
}

// InvalidParameterValueError is an error that is returned if the crypt.Crypter
// implementation receives an invalid hash function parameter value.
type InvalidParameterValueError struct {
	// Name is the name of the parameter.
	Name string
	// Value is the value of the parameter.
	Value string
	// Expected is a free-form expected format of the value.
	Expected string
}

// Error implements the error interface.
func (e *InvalidParameterValueError) Error() string {
	const m = "invalid value for parameter"
	if e == nil {
		return m
	}
	return fmt.Sprintf(m+" %q (given %q, expected %s)", e.Name, e.Value, e.Expected)
}

// MalformedParametersError is an error that is returned if the crypt.Crypter
// implementation receives a malformed hash function parameters.
type MalformedParametersError struct {
	// Unparsed is the unparsed parameters string.
	Unparsed string
}

// Error implements the error interface.
func (e *MalformedParametersError) Error() string {
	const m = "malformed parameters"
	if e == nil {
		return m
	}
	return fmt.Sprintf(m+" (unparsed %q)", e.Unparsed)
}

// MissingRequiredParametersError is an error that is returned if the
// crypt.Crypter implementation does not receive one of the required hash
// function parameters.
type MissingRequiredParametersError struct {
	// Required is the free-form description of the required parameters.
	Required string
}

// Error implements the error interface.
func (e *MissingRequiredParametersError) Error() string {
	const m = "missing required parameters"
	if e == nil {
		return m
	}
	return fmt.Sprintf(m+" (need %s)", e.Required)
}

// InvalidOutputLengthError is an error that is returned if the crypt.Crypter
// implementation receives a hash output with invalid length for the password
// hashing algorithm.
type InvalidOutputLengthError struct {
	// Length is the given function output length.
	Length int
	// Expected is the free-form description of the expected length.
	Expected string
}

// Error implements the error interface.
func (e *InvalidOutputLengthError) Error() string {
	const m = "invalid output length"
	if e == nil {
		return m
	}
	return fmt.Sprintf(m+" %d (expected %s)", e.Length, e.Expected)
}
