package phcformat

import (
	"strings"
	"testing"

	"go.pact.im/x/phcformat/option"
)

func TestMustParse(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Panics bool
	}{{
		"PanicsOnEmptyInput",
		"",
		true,
	}, {
		"AcceptsGoodInput",
		"$",
		false,
	}}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			defer func() {
				e := recover()
				switch {
				case e == nil && tc.Panics:
					t.Fatal("expected panic")
				case e != nil && !tc.Panics:
					t.Fatal("unexpected panic")
				}
			}()
			_ = MustParse(tc.Input)
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string
		Hash  Hash
		Good  bool
	}{{
		"RejectsEmptyString",
		"",
		Hash{},
		false,
	}, {
		"RejectsWithoutLeadingSep",
		"algo",
		Hash{},
		false,
	}, {
		"AcceptsEmptyID",
		"$",
		Hash{},
		true,
	}, {
		"AcceptsIDLength32",
		"$" + strings.Repeat("z", 32),
		Hash{ID: strings.Repeat("z", 32)},
		true,
	}, {
		"RejectsIDLength33",
		"$" + strings.Repeat("z", 33),
		Hash{},
		false,
	}, {
		"RejectsInvalidID",
		"$!",
		Hash{},
		false,
	}, {
		"AcceptsEmptyVersion",
		"$$v=",
		Hash{Version: option.Value("")},
		true,
	}, {
		"AcceptsNumericVersion",
		"$$v=42",
		Hash{Version: option.Value("42")},
		true,
	}, {
		"AcceptsVersionAndParams",
		"$$v=$k=v",
		Hash{Version: option.Value(""), Params: option.Value("k=v")},
		true,
	}, {
		"AcceptsParamV",
		"$$v=param",
		Hash{Params: option.Value("v=param")},
		true,
	}, {
		"AcceptsMultipleParamsWithV",
		"$$v=,k=v",
		Hash{Params: option.Value("v=,k=v")},
		true,
	}, {
		"RejectsInvalidVersion",
		"$$v=!",
		Hash{},
		false,
	}, {
		"AcceptsEmptyParamName",
		"$$=value",
		Hash{Params: option.Value("=value")},
		true,
	}, {
		"AcceptsEmptyParamValue",
		"$$name=",
		Hash{Params: option.Value("name=")},
		true,
	}, {
		"AcceptsEmptyParam",
		"$$=",
		Hash{Params: option.Value("=")},
		true,
	}, {
		"AcceptsMultipleEmptyParams",
		"$$=,name=,=value",
		Hash{Params: option.Value("=,name=,=value")},
		true,
	}, {
		"RejectsInvalidParamName",
		"$$!=",
		Hash{},
		false,
	}, {
		"RejectsInvalidParamValue",
		"$$=!",
		Hash{},
		false,
	}, {
		"RejectsParamTrailingComma",
		"$$k=v,",
		Hash{},
		false,
	}, {
		"AcceptsSalt",
		"$$salt",
		Hash{Salt: option.Value("salt")},
		true,
	}, {
		"AcceptsParamsAndSalt",
		"$$k=v$salt",
		Hash{Params: option.Value("k=v"), Salt: option.Value("salt")},
		true,
	}, {
		"AcceptsSaltBase64",
		"$$gZiV/M1gPc22ElAH/Jh1Hw",
		Hash{Salt: option.Value("gZiV/M1gPc22ElAH/Jh1Hw")},
		true,
	}, {
		"AcceptsSaltAndHash",
		"$$salt$out",
		Hash{Salt: option.Value("salt"), Output: option.Value("out")},
		true,
	}, {
		"AcceptsEmptySaltAndHash",
		"$$$",
		Hash{Salt: option.Value(""), Output: option.Value("")},
		true,
	}, {
		"AcceptsParamsAndSaltAndHash",
		"$$k=v$salt$out",
		Hash{Params: option.Value("k=v"), Salt: option.Value("salt"), Output: option.Value("out")},
		true,
	}, {
		"RejectsParamsAndInvalidSalt",
		"$$k=v$!",
		Hash{},
		false,
	}, {
		"RejectsSaltAndInvalidHash",
		"$$salt$!",
		Hash{},
		false,
	}}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testParse(t, tc.Input, tc.Hash, tc.Good)
		})
	}
}

func testParse(t *testing.T, input string, expected Hash, good bool) {
	if good {
		expected.Raw = input
	}
	h, ok := Parse(input)
	switch {
	case ok && !good:
		t.Fatal("expected parse error")
	case !ok && good:
		t.Fatal("unexpected parse error")
	}
	if h != expected {
		t.Fatal("invalid parsed hash")
	}
}
