package phcformat

import (
	"reflect"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		Name    string
		HashID  string
		Version OptionalString
		Salt    HashSalt
		Output  []byte
		Params  []HashParam
		Hash    Hash
		Good    bool
	}{{
		Name:   "AcceptsEmptyHashID",
		HashID: "",
		Hash: Hash{
			ID:  "",
			Raw: "$",
		},
		Good: true,
	}, {
		Name:   "RejectsInvalidHashID",
		HashID: "!",
	}, {
		Name:   "RejectsInvalidHashIDLen32",
		HashID: strings.Repeat("z", 32),
		Hash: Hash{
			ID:  strings.Repeat("z", 32),
			Raw: "$" + strings.Repeat("z", 32),
		},
		Good: true,
	}, {
		Name:   "RejectsInvalidHashIDLen33",
		HashID: strings.Repeat("z", 33),
	}, {
		Name:    "AcceptsEmptyVersion",
		Version: String(""),
		Hash: Hash{
			Version: String(""),
			Raw:     "$$v=",
		},
		Good: true,
	}, {
		Name:    "AcceptsNumericVersion",
		Version: String("42"),
		Hash: Hash{
			Version: String("42"),
			Raw:     "$$v=42",
		},
		Good: true,
	}, {
		Name:    "RejectsInvalidVersion",
		Version: String("!"),
	}, {
		Name: "AcceptsParam",
		Params: []HashParam{
			{"k", "v"},
		},
		Hash: Hash{
			Params: String("k=v"),
			Raw:    "$$k=v",
		},
		Good: true,
	}, {
		Name: "AcceptsParamV",
		Params: []HashParam{
			{"v", "param"},
		},
		Hash: Hash{
			Params: String("v=param"),
			Raw:    "$$v=param",
		},
		Good: true,
	}, {
		Name: "RejectsAmbiguousEmptyParamV",
		Params: []HashParam{
			{"v", ""},
		},
	}, {
		Name: "RejectsAmbiguousNumericParamV",
		Params: []HashParam{
			{"v", "42"},
		},
	}, {
		Name:    "AcceptsVersionAndParamV",
		Version: String("42"),
		Params: []HashParam{
			{"v", "42"},
		},
		Hash: Hash{
			Version: String("42"),
			Params:  String("v=42"),
			Raw:     "$$v=42$v=42",
		},
		Good: true,
	}, {
		Name:   "RejectsInvalidParamName",
		Params: []HashParam{{"!", ""}},
	}, {
		Name:   "RejectsInvalidParamValue",
		Params: []HashParam{{"", "!"}},
	}, {
		Name: "AcceptsMultipleParams",
		Params: []HashParam{
			{"a", "b"},
			{"c", "d"},
		},
		Hash: Hash{
			Params: String("a=b,c=d"),
			Raw:    "$$a=b,c=d",
		},
		Good: true,
	}, {
		Name: "AcceptsSaltStringEncoded",
		Salt: HashSalt{Format: HashSaltFormatEncoded, String: "salt"},
		Hash: Hash{
			Salt: String("salt"),
			Raw:  "$$salt",
		},
		Good: true,
	}, {
		Name: "AcceptsSaltBytesEncoded",
		Salt: HashSalt{Format: HashSaltFormatEncoded, Bytes: []byte("salt")},
		Hash: Hash{
			Salt: String("salt"),
			Raw:  "$$salt",
		},
		Good: true,
	}, {
		Name: "AcceptsSaltStringBase64",
		Salt: HashSalt{Format: HashSaltFormatBase64, String: "salt"},
		Hash: Hash{
			Salt: String("c2FsdA"),
			Raw:  "$$c2FsdA",
		},
		Good: true,
	}, {
		Name: "AcceptsSaltBytesBase64",
		Salt: HashSalt{Format: HashSaltFormatBase64, Bytes: []byte("salt")},
		Hash: Hash{
			Salt: String("c2FsdA"),
			Raw:  "$$c2FsdA",
		},
		Good: true,
	}, {
		Name: "RejectsInvalidSaltStringEncoded",
		Salt: HashSalt{Format: HashSaltFormatEncoded, String: "!"},
	}, {
		Name: "RejectsInvalidSaltBytesEncoded",
		Salt: HashSalt{Format: HashSaltFormatEncoded, Bytes: []byte("!")},
	}, {
		Name:   "RejectsOutputWithoutSalt",
		Output: []byte{},
	}, {
		Name:   "AcceptsEmptyOutputAndSalt",
		Salt:   HashSalt{Format: HashSaltFormatEncoded, Bytes: []byte("")},
		Output: []byte{},
		Hash: Hash{
			Salt:   String(""),
			Output: String(""),
			Raw:    "$$$",
		},
		Good: true,
	}}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, ok := Encode(tc.HashID, tc.Version, tc.Salt, tc.Output, tc.Params...)
			switch {
			case ok && !tc.Good:
				t.Fatal("expected encode error")
			case !ok && tc.Good:
				t.Fatal("unexpected encode error")
			}
			if !reflect.DeepEqual(h, tc.Hash) {
				t.Fatal("invalid encoded hash")
			}
		})
	}
}
