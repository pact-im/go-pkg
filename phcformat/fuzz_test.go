package phcformat

import (
	"testing"
)

func FuzzParse(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		if _, ok := Parse(s); !ok {
			t.SkipNow()
		}
	})
}
