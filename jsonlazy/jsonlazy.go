// Package jsonlazy implements delayed JSON encoding for values that do not
// implement json.Marshaler interface.
package jsonlazy

import (
	"encoding/json"
	"sync"
)

var (
	_ json.Marshaler = (*Marshaler)(nil)
	_ json.Marshaler = (*OnceMarshaler)(nil)
)

// Marshaler implements a delayed JSON encoding for values.
type Marshaler struct {
	v interface{}
}

// NewMarshaler returns a new Marshaler for the given value.
func NewMarshaler(v interface{}) *Marshaler {
	return &Marshaler{v}
}

// MarshalJSON implemnets the json.Marshaler interface.
func (m *Marshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.v)
}

// Once is a shorthand for jsonlazy.Once(m).
func (m *Marshaler) Once() *OnceMarshaler {
	return Once(m)
}

// OnceMarshaler is a json.Marshaler that caches the result of MarshalJSON call.
type OnceMarshaler struct {
	m json.Marshaler

	once sync.Once
	data []byte
	err  error
}

// Once returns a new marshaler that caches the result of MarshalJSON call. That
// is, it calls MarshalJSON at most once.
func Once(m json.Marshaler) *OnceMarshaler {
	return &OnceMarshaler{
		m: m,
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (m *OnceMarshaler) MarshalJSON() ([]byte, error) {
	m.once.Do(func() {
		m.data, m.err = m.m.MarshalJSON()
	})
	return m.data, m.err
}
