package encode

import (
	"encoding/base64"
	"testing"

	"go.pact.im/x/phcformat/option"
)

func TestNil(t *testing.T) {
	if !option.IsNil(Nil()) {
		t.Fail()
	}
}

func TestEmpty(t *testing.T) {
	if NewEmpty().Append(nil) != nil {
		t.Fail()
	}
}

func TestOption(t *testing.T) {
	if NewOption(option.Nil[Appender]()).Append(nil) != nil {
		t.Fail()
	}
	if string(NewOption(option.Value(String("foo"))).Append(nil)) != "foo" {
		t.Fail()
	}
}

func TestConcat(t *testing.T) {
	if string(NewConcat(String("foo"), String("bar")).Append(nil)) != "foobar" {
		t.Fail()
	}
}

func TestList(t *testing.T) {
	if NewList(String(","), []Appender(nil)...).Append(nil) != nil {
		t.Fail()
	}
	if string(NewList(String(","), String("foo")).Append(nil)) != "foo" {
		t.Fail()
	}
	if string(NewList(String(","), String("foo"), String("bar")).Append(nil)) != "foo,bar" {
		t.Fail()
	}
}

func TestKV(t *testing.T) {
	if string(NewKV(String("="), String("foo"), String("bar")).Append(nil)) != "foo=bar" {
		t.Fail()
	}
}

func TestByte(t *testing.T) {
	if string(NewByte('x').Append(nil)) != "x" {
		t.Fail()
	}
}

func TestString(t *testing.T) {
	if string(NewString("foo").Append(nil)) != "foo" {
		t.Fail()
	}
}

func TestBytes(t *testing.T) {
	if string(NewBytes([]byte("foo")).Append(nil)) != "foo" {
		t.Fail()
	}
}

func TestUint(t *testing.T) {
	if string(NewUint(42).Append(nil)) != "42" {
		t.Fail()
	}
}

func TestBase64(t *testing.T) {
	if string(NewBase64("base64").Append(nil)) != "YmFzZTY0" {
		t.Fail()
	}
	if string(NewBase64("base64").Append(make([]byte, 0, base64.RawStdEncoding.EncodedLen(len("base64"))))) != "YmFzZTY0" {
		t.Fail()
	}
}
