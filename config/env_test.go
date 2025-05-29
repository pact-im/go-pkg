package config

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestEnv(t *testing.T) {
	testCases := []struct {
		EnvFile string
		EnvVars []string
	}{
		{
			strings.Join([]string{
				"",
				"FOO=BAR",
				"#",
			}, "\n"),
			[]string{
				"FOO=BAR",
			},
		},
		{
			strings.Join([]string{
				" #FOO=BAR",
				"# FOO=BAR",
				"#FOO=BAR",
				"FOO=BAR #",
				"FOO=BAR#",
			}, "\n"),
			[]string{
				"FOO=BAR",
				"FOO=BAR",
			},
		},
	}
	for _, tc := range testCases {
		env, err := parseEnv(strings.NewReader(tc.EnvFile))
		assert.Check(t, err)
		assert.Check(t, is.DeepEqual(tc.EnvVars, env))
	}
}
