package phcformat

import (
	"reflect"
	"strings"
	"testing"
)

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
		Hash{Version: String("")},
		true,
	}, {
		"AcceptsNumericVersion",
		"$$v=42",
		Hash{Version: String("42")},
		true,
	}, {
		"AcceptsVersionAndParams",
		"$$v=$k=v",
		Hash{Version: String(""), Params: String("k=v")},
		true,
	}, {
		"AcceptsParamV",
		"$$v=param",
		Hash{Params: String("v=param")},
		true,
	}, {
		"AcceptsMultipleParamsWithV",
		"$$v=,k=v",
		Hash{Params: String("v=,k=v")},
		true,
	}, {
		"RejectsInvalidVersion",
		"$$v=!",
		Hash{},
		false,
	}, {
		"AcceptsEmptyParamName",
		"$$=value",
		Hash{Params: String("=value")},
		true,
	}, {
		"AcceptsEmptyParamValue",
		"$$name=",
		Hash{Params: String("name=")},
		true,
	}, {
		"AcceptsEmptyParam",
		"$$=",
		Hash{Params: String("=")},
		true,
	}, {
		"AcceptsMultipleEmptyParams",
		"$$=,name=,=value",
		Hash{Params: String("=,name=,=value")},
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
		Hash{Salt: String("salt")},
		true,
	}, {
		"AcceptsParamsAndSalt",
		"$$k=v$salt",
		Hash{Params: String("k=v"), Salt: String("salt")},
		true,
	}, {
		"AcceptsSaltBase64",
		"$$gZiV/M1gPc22ElAH/Jh1Hw",
		Hash{Salt: String("gZiV/M1gPc22ElAH/Jh1Hw")},
		true,
	}, {
		"AcceptsSaltAndHash",
		"$$salt$out",
		Hash{Salt: String("salt"), Output: String("out")},
		true,
	}, {
		"AcceptsEmptySaltAndHash",
		"$$$",
		Hash{Salt: String(""), Output: String("")},
		true,
	}, {
		"AcceptsParamsAndSaltAndHash",
		"$$k=v$salt$out",
		Hash{Params: String("k=v"), Salt: String("salt"), Output: String("out")},
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
		if tc.Good {
			tc.Hash.Raw = tc.Input
		}
		t.Run(tc.Name, func(t *testing.T) {
			h, ok := Parse(tc.Input)
			switch {
			case ok && !tc.Good:
				t.Fatal("expected parse error")
			case !ok && tc.Good:
				t.Fatal("unexpected parse error")
			}
			if !reflect.DeepEqual(h, tc.Hash) {
				t.Fatal("invalid parsed hash")
			}
		})
	}
}
