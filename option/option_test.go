package option

import (
	"strconv"
	"testing"
)

func TestValue(t *testing.T) {
	const value = "v"
	opt := Value(value)
	v, ok := opt.Unwrap()
	if !ok || v != value {
		t.Fail()
	}
}

func TestNil(t *testing.T) {
	opt := Nil[any]()
	v, ok := opt.Unwrap()
	if ok || v != nil {
		t.Fail()
	}
}

func TestUnwrapOrZero(t *testing.T) {
	const value = "v"
	if UnwrapOrZero(Value(value)) != value {
		t.Fail()
	}
	if UnwrapOrZero(Nil[string]()) != "" {
		t.Fail()
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(Nil[any]()) {
		t.Fail()
	}
	if IsNil(Value[any](nil)) {
		t.Fail()
	}
}

func TestMap(t *testing.T) {
	if UnwrapOrZero(Map(Nil[int](), strconv.Itoa)) != "" {
		t.Fail()
	}
	if UnwrapOrZero(Map(Value(42), strconv.Itoa)) != "42" {
		t.Fail()
	}
}
