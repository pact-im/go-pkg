package zapjournal

import (
	"testing"
)

func TestAppendVarName(t *testing.T) {
	testCases := []struct {
		parts  []string
		expect string
	}{
		{nil, defaultPrefix},
		{[]string{""}, defaultPrefix},
		{[]string{"", ""}, defaultPrefix},
		{[]string{"", "foo"}, "FOO"},
		{[]string{"foo", ""}, "FOO"},
		{[]string{"foo", "bar"}, "FOO_BAR"},
		{[]string{"foo", "bar"}, "FOO_BAR"},
		{[]string{"foo", "_bar"}, "FOO_BAR"},
		{[]string{"foo_", "bar"}, "FOO_BAR"},
		{[]string{"foo", "bar-baz"}, "FOO_BAR_BAZ"},
		{[]string{"foo-bar", "baz"}, "FOO_BAR_BAZ"},
		{[]string{"__foo__", "__bar__"}, "FOO_BAR"},
		{[]string{"__foo___bar__"}, "FOO_BAR"},
		{[]string{"__foo@bar__"}, "FOO_BAR"},
		{[]string{"foo\xffbar"}, "FOOBAR"},
		{[]string{"camelCase"}, "CAMEL_CASE"},
		{[]string{"UPPER_SNAKE"}, "UPPER_SNAKE"},
		{[]string{"Title"}, "TITLE"},
		{[]string{"netIPAddr"}, "NET_IP_ADDR"},
		{[]string{"96neko"}, "96_NEKO"},
		{[]string{"42"}, "42"},
	}
	for _, tc := range testCases {
		out := string(appendVarName(nil, tc.parts...))
		if out != tc.expect {
			t.Fatalf("expected %q, got %q", tc.expect, out)
		}
	}
}
