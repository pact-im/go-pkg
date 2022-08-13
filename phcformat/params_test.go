package phcformat

import (
	"testing"
)

func TestParamsIterator(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string
		State []ParamsIterator
	}{{
		"EmptyInput",
		"",
		nil,
	}, {
		"InvalidInput",
		"bad",
		[]ParamsIterator{{
			After: "bad",
		}},
	}, {
		"CommaInput",
		",",
		[]ParamsIterator{{
			After: ",",
		}},
	}, {
		"EmptySingleParam",
		"=",
		[]ParamsIterator{{
			Name:  "",
			Value: "",
			Valid: true,
		}},
	}, {
		"SingleParam",
		"k=v",
		[]ParamsIterator{{
			Name:  "k",
			Value: "v",
			Valid: true,
		}},
	}, {
		"EmptyMultipleParams",
		"=,=",
		[]ParamsIterator{{
			Name:  "",
			Value: "",
			After: "=",
			Valid: true,
		}, {
			Name:  "",
			Value: "",
			Valid: true,
		}},
	}, {
		"MultipleParams",
		"a=b,c=d",
		[]ParamsIterator{{
			Name:  "a",
			Value: "b",
			After: "c=d",
			Valid: true,
		}, {
			Name:  "c",
			Value: "d",
			Valid: true,
		}},
	}, {
		"SingleParamTrailingComma",
		"k=v,",
		[]ParamsIterator{{
			Name:  "k",
			Value: "v",
			After: ",",
			Valid: true,
		}, {
			After: ",",
		}},
	}, {
		"MultipleParamsTrailingComma",
		"a=b,c=d,",
		[]ParamsIterator{{
			Name:  "a",
			Value: "b",
			After: "c=d,",
			Valid: true,
		}, {
			Name:  "c",
			Value: "d",
			After: ",",
			Valid: true,
		}, {
			After: ",",
		}},
	}}
	compareState := func(t *testing.T, it ParamsIterator, i int, states []ParamsIterator) {
		var expected ParamsIterator
		if i < len(states) {
			expected = states[i]
		}
		if it != expected {
			t.Fatalf("invalid iterator state %d: %#v != %#v", i, it, expected)
		}
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			it := IterParams(tc.Input)
			for i := 0; ; i++ {
				compareState(t, it, i, tc.State)
				if !it.Valid {
					break
				}
				it = it.Next()
			}
		})
	}
}
