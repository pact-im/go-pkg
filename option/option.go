// Package option implements optional values to represent the lack of value
// without pointers.
package option

// Of represents an optional value that may be nil.
type Of[T any] struct {
	isSet bool
	value T
}

// Value returns an option with the given value.
func Value[T any](v T) Of[T] {
	return Of[T]{true, v}
}

// Nil returns nil option with type T.
func Nil[T any]() Of[T] {
	return Of[T]{}
}

// Unwrap returns the underlying value and a boolean flag indicating whether it
// is set.
func (v Of[T]) Unwrap() (T, bool) {
	return v.value, v.isSet
}

// UnwrapOrZero returns the option value or its zero value if it is not set.
func UnwrapOrZero[T any](opt Of[T]) T {
	v, ok := opt.Unwrap()
	if !ok {
		var zero T
		return zero
	}
	return v
}

// IsNil returns true if the value is nil.
func IsNil[T any](opt Of[T]) bool {
	return !opt.isSet
}

// Map applies function f to the underlying option value if it is set.
func Map[T, U any](opt Of[T], f func(T) U) Of[U] {
	v, ok := opt.Unwrap()
	if !ok {
		return Nil[U]()
	}
	return Value(f(v))
}
