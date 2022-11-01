package phcformat

import (
	"strings"
	"testing"

	"go.pact.im/x/option"
	"go.pact.im/x/phcformat/encode"
)

func TestAppend(t *testing.T) {
	testCases := []struct {
		Name    string
		HashID  encode.Appender
		Version encode.Appender
		Params  encode.Appender
		Salt    encode.Appender
		Output  encode.Appender
		Expect  string
		Parsed  option.Of[Hash]
	}{{
		Name:   "HashID",
		HashID: encode.String("name"),
		Expect: "$name",
		Parsed: option.Value(Hash{
			ID: "name",
		}),
	}, {
		Name:   "EmptyHashID",
		Expect: "$",
		Parsed: option.Value(Hash{}),
	}, {
		Name:    "Version",
		Version: encode.NewUint(42),
		Expect:  "$$v=42",
		Parsed: option.Value(Hash{
			Version: option.Value("42"),
		}),
	}, {
		Name:    "EmptyVersion",
		Version: encode.NewEmpty(),
		Expect:  "$$v=",
		Parsed: option.Value(Hash{
			Version: option.Value(""),
		}),
	}, {
		Name:   "AmbiguousParamV",
		Params: encode.String("v="),
		Expect: "$$v=",
		Parsed: option.Value(Hash{
			Version: option.Value(""),
		}),
	}, {
		Name:   "Params",
		Params: encode.String("k=v"),
		Expect: "$$k=v",
		Parsed: option.Value(Hash{
			Params: option.Value("k=v"),
		}),
	}, {
		Name:   "Salt",
		Salt:   encode.NewBase64([]byte{0xB1, 0xA9, 0x6D}),
		Expect: "$$salt",
		Parsed: option.Value(Hash{
			Salt: option.Value("salt"),
		}),
	}, {
		Name:   "EmptySalt",
		Salt:   encode.NewEmpty(),
		Expect: "$$",
		Parsed: option.Value(Hash{
			Salt: option.Value(""),
		}),
	}, {
		Name:   "OutputWithoutSalt",
		Output: encode.NewBase64([]byte{0x85, 0xAB, 0x21}),
		Expect: "$$hash",
		Parsed: option.Value(Hash{
			Salt: option.Value("hash"),
		}),
	}, {
		Name:   "SaltAndOutput",
		Salt:   encode.NewBase64([]byte{0xB1, 0xA9, 0x6D}),
		Output: encode.NewBase64([]byte{0x85, 0xAB, 0x21}),
		Expect: "$$salt$hash",
		Parsed: option.Value(Hash{
			Salt:   option.Value("salt"),
			Output: option.Value("hash"),
		}),
	}}
	orEmpty := func(v encode.Appender) encode.Appender {
		if v == nil {
			return encode.NewEmpty()
		}
		return v
	}
	orNil := func(v encode.Appender) option.Of[encode.Appender] {
		if v == nil {
			return encode.Nil()
		}
		return option.Value(v)
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			sentinel := "x"
			buf := Append([]byte(sentinel), orEmpty(tc.HashID), orNil(tc.Version), orNil(tc.Params), orNil(tc.Salt), orNil(tc.Output))
			if string(buf) != sentinel+tc.Expect {
				t.Fatal("invalid appended hash")
			}
			h, ok := tc.Parsed.Unwrap()
			testParse(t, strings.TrimPrefix(string(buf), sentinel), h, ok)
		})
	}
}
