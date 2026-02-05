package httprange

import (
	"net/http"
	"testing"
)

func TestParseContentRange(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantUnit string
		wantResp string
		wantOK   bool
	}{
		{
			name:     "bytes range",
			input:    "bytes 0-99/*",
			wantUnit: "bytes",
			wantResp: "0-99/*",
			wantOK:   true,
		},
		{
			name:     "custom range",
			input:    "custom-unit-name123 0.1+0.42",
			wantUnit: "custom-unit-name123",
			wantResp: "0.1+0.42",
			wantOK:   true,
		},
		{
			name:   "missing space",
			input:  "bytes0-99/100",
			wantOK: false,
		},
		{
			name:   "empty string",
			input:  "",
			wantOK: false,
		},
		{
			name:   "space only",
			input:  " ",
			wantOK: false,
		},
		{
			name:   "unit only",
			input:  "bytes",
			wantOK: false,
		},
		{
			name:     "unit with space only",
			input:    "bytes ",
			wantUnit: "bytes",
			wantResp: "",
			wantOK:   true,
		},
		{
			name:   "space then range only",
			input:  " 0-99/100",
			wantOK: false,
		},
		{
			name:   "invalid token character (parenthesis)",
			input:  "unit(test) 0-99/100",
			wantOK: false,
		},
		{
			name:   "invalid token character (comma)",
			input:  "unit,test 0-99/100",
			wantOK: false,
		},
		{
			name:   "invalid unit character (colon)",
			input:  "unit:test 0-99/100",
			wantOK: false,
		},
		{
			name:   "range with non-ASCII character",
			input:  "unit 0-99/©",
			wantOK: false,
		},
		{
			name:   "range with non-ASCII byte",
			input:  "bytes 0-99/\x80",
			wantOK: false,
		},
		{
			name:   "range with null character",
			input:  "bytes 0-99/\x00100",
			wantOK: false,
		},
		{
			name:     "range with newline",
			input:    "unit 0-99/\n100",
			wantUnit: "unit",
			wantResp: "0-99/\n100",
			wantOK:   true,
		},
		{
			name:     "multiple spaces before range",
			input:    "unit  0-99/100",
			wantUnit: "unit",
			wantResp: " 0-99/100",
			wantOK:   true,
		},
		{
			name:   "tab instead of space",
			input:  "bytes\t0-99/100",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseContentRange(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("parseContentRange() ok = %t, want %t", ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if got.Unit != tt.wantUnit {
				t.Errorf("parseContentRange() Unit = %q, want %q", got.Unit, tt.wantUnit)
			}
			if got.Resp != tt.wantResp {
				t.Errorf("parseContentRange() Resp = %q, want %q", got.Resp, tt.wantResp)
			}
		})
	}
}

func TestParseContentRangeFromHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		unit    string
		want    string
		wantErr bool
	}{
		{
			name: "valid range",
			headers: http.Header{
				httpHeaderContentRange: {"bytes 0-99/100"},
			},
			unit:    "bytes",
			want:    "0-99/100",
			wantErr: false,
		},
		{
			name: "missing header",
			headers: http.Header{
				"Other-Header": {"value"},
			},
			unit:    "bytes",
			wantErr: true,
		},
		{
			name: "empty header value",
			headers: http.Header{
				httpHeaderContentRange: {""},
			},
			unit:    "bytes",
			wantErr: true,
		},
		{
			name: "unit mismatch",
			headers: http.Header{
				httpHeaderContentRange: {"seconds 0-59/60"},
			},
			unit:    "bytes",
			wantErr: true,
		},
		{
			name: "multiple values takes first",
			headers: http.Header{
				httpHeaderContentRange: {
					"bytes 0-99/100",
					"bytes 100-199/200",
				},
			},
			unit:    "bytes",
			want:    "0-99/100",
			wantErr: false,
		},
		{
			name: "case-sensitive unit mismatch",
			headers: http.Header{
				httpHeaderContentRange: {"Bytes 0-99/100"},
			},
			unit:    "bytes",
			want:    "0-99/100",
			wantErr: false,
		},
		{
			name: "header with extra spaces",
			headers: http.Header{
				httpHeaderContentRange: {"  bytes  0-99/100  "},
			},
			unit:    "bytes",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseContentRangeFromHeader(tt.headers, tt.unit)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("parseContentRangeFromHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseUnsatisfiedRangeResp(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   int64
		wantOK bool
	}{
		{
			name:   "valid with length",
			input:  "*/100",
			want:   100,
			wantOK: true,
		},
		{
			name:   "valid zero length",
			input:  "*/0",
			want:   0,
			wantOK: true,
		},
		{
			name:   "missing asterisk",
			input:  "/100",
			wantOK: false,
		},
		{
			name:   "missing slash",
			input:  "*100",
			wantOK: false,
		},
		{
			name:   "empty after slash",
			input:  "*/",
			wantOK: false,
		},
		{
			name:   "not digits after slash",
			input:  "*/abc",
			wantOK: false,
		},
		{
			name:   "mixed digits and letters",
			input:  "*/123abc",
			wantOK: false,
		},
		{
			name:   "negative number",
			input:  "*/-100",
			wantOK: false,
		},
		{
			name:   "empty string",
			input:  "",
			wantOK: false,
		},
		{
			name:   "just asterisk",
			input:  "*",
			wantOK: false,
		},
		{
			name:   "just slash",
			input:  "/",
			wantOK: false,
		},
		{
			name:   "extra characters before asterisk",
			input:  "abc*/100",
			wantOK: false,
		},
		{
			name:   "extra characters after length",
			input:  "*/100extra",
			wantOK: false,
		},
		{
			name:   "space before asterisk",
			input:  " */100",
			wantOK: false,
		},
		{
			name:   "space after asterisk",
			input:  "* /100",
			wantOK: false,
		},
		{
			name:   "space before length",
			input:  "*/ 100",
			wantOK: false,
		},
		{
			name:   "space after length",
			input:  "*/100 ",
			wantOK: false,
		},
		{
			name:   "max int64",
			input:  "*/9223372036854775807",
			want:   9223372036854775807,
			wantOK: true,
		},
		{
			name:   "overflow int64",
			input:  "*/9223372036854775808",
			wantOK: false,
		},
		{
			name:   "leading zeros",
			input:  "*/00100",
			want:   100,
			wantOK: true,
		},
		{
			name:   "all zeros",
			input:  "*/000",
			want:   0,
			wantOK: true,
		},
		{
			name:   "only asterisk and slash",
			input:  "*/",
			wantOK: false,
		},
		{
			name:   "asterisk slash asterisk",
			input:  "*/*",
			wantOK: false,
		},
		{
			name:   "with decimal point",
			input:  "*/100.5",
			wantOK: false,
		},
		{
			name:   "with underscore",
			input:  "*/100_000",
			wantOK: false,
		},
		{
			name:   "with plus sign",
			input:  "*/+100",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseUnsatisfiedRangeResp(tt.input)
			if ok != tt.wantOK {
				t.Fatalf(
					"parseUnsatisfiedRangeResp(%q) ok = %t, want %t",
					tt.input, ok, tt.wantOK,
				)
			}
			if !ok {
				return
			}
			if got.CompleteLength != tt.want {
				t.Errorf(
					"parseUnsatisfiedRangeResp(%q) CompleteLength = %d, want %d",
					tt.input, got.CompleteLength, tt.want,
				)
			}
		})
	}
}

func TestParseBytesRangeResp(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   bytesRangeResp
		wantOK bool
	}{
		{
			name:   "valid range with known length",
			input:  "0-99/100",
			want:   bytesRangeResp{First: 0, Last: 99, CompleteLength: 100},
			wantOK: true,
		},
		{
			name:   "valid single byte range",
			input:  "50-50/100",
			want:   bytesRangeResp{First: 50, Last: 50, CompleteLength: 100},
			wantOK: true,
		},
		{
			name:   "valid at start",
			input:  "0-0/100",
			want:   bytesRangeResp{First: 0, Last: 0, CompleteLength: 100},
			wantOK: true,
		},
		{
			name:   "valid at end",
			input:  "99-99/100",
			want:   bytesRangeResp{First: 99, Last: 99, CompleteLength: 100},
			wantOK: true,
		},
		{
			name:   "valid with leading zeros",
			input:  "000-099/100",
			want:   bytesRangeResp{First: 0, Last: 99, CompleteLength: 100},
			wantOK: true,
		},
		{
			name:   "valid unknown length",
			input:  "0-99/*",
			want:   bytesRangeResp{First: 0, Last: 99},
			wantOK: true,
		},
		{
			name:   "valid single byte unknown length",
			input:  "50-50/*",
			want:   bytesRangeResp{First: 50, Last: 50},
			wantOK: true,
		},
		{
			name:   "invalid zero length content",
			input:  "0-0/0",
			wantOK: false,
		},
		{
			name:   "first > last",
			input:  "100-50/200",
			wantOK: false,
		},
		{
			name:   "last ≥ complete length",
			input:  "0-100/100",
			wantOK: false,
		},
		{
			name:   "first ≥ complete length",
			input:  "100-150/100",
			wantOK: false,
		},
		{
			name:   "both > complete length",
			input:  "100-150/50",
			wantOK: false,
		},
		{
			name:   "missing dash",
			input:  "0100/200",
			wantOK: false,
		},
		{
			name:   "missing slash",
			input:  "0-100",
			wantOK: false,
		},
		{
			name:   "empty first",
			input:  "-100/200",
			wantOK: false,
		},
		{
			name:   "empty last",
			input:  "0-/200",
			wantOK: false,
		},
		{
			name:   "non-digit first",
			input:  "abc-100/200",
			wantOK: false,
		},
		{
			name:   "non-digit last",
			input:  "0-abc/200",
			wantOK: false,
		},
		{
			name:   "non-digit length",
			input:  "0-100/abc",
			wantOK: false,
		},
		{
			name:   "negative first",
			input:  "-10-100/200",
			wantOK: false,
		},
		{
			name:   "negative last",
			input:  "0--10/200",
			wantOK: false,
		},
		{
			name:   "negative length",
			input:  "0-100/-200",
			wantOK: false,
		},
		{
			name:   "empty string",
			input:  "",
			wantOK: false,
		},
		{
			name:   "empty all",
			input:  "-/",
			wantOK: false,
		},
		{
			name:   "only asterisk",
			input:  "*",
			wantOK: false,
		},
		{
			name:   "multiple dashes",
			input:  "0-100-200/300",
			wantOK: false,
		},
		{
			name:   "multiple slashes",
			input:  "0-100/200/300",
			wantOK: false,
		},
		{
			name:   "space before first",
			input:  " 0-100/200",
			wantOK: false,
		},
		{
			name:   "space after last",
			input:  "0-100 /200",
			wantOK: false,
		},
		{
			name:   "space before slash",
			input:  "0-100 /200",
			wantOK: false,
		},
		{
			name:   "space after slash",
			input:  "0-100/ 200",
			wantOK: false,
		},
		{
			name:   "space before asterisk",
			input:  "0-100/ *",
			wantOK: false,
		},
		{
			name:   "plus sign in first",
			input:  "+0-100/200",
			wantOK: false,
		},
		{
			name:   "plus sign in last",
			input:  "0-+100/200",
			wantOK: false,
		},
		{
			name:   "plus sign in length",
			input:  "0-100/+200",
			wantOK: false,
		},
		{
			name:   "decimal in first",
			input:  "0.5-100/200",
			wantOK: false,
		},
		{
			name:   "decimal in last",
			input:  "0-100.5/200",
			wantOK: false,
		},
		{
			name:   "decimal in length",
			input:  "0-100/200.5",
			wantOK: false,
		},
		{
			name:   "overflow int64 first",
			input:  "9223372036854775808-9223372036854775808/9223372036854775809",
			wantOK: false,
		},
		{
			name:   "overflow int64 last",
			input:  "0-9223372036854775808/9223372036854775809",
			wantOK: false,
		},
		{
			name:   "overflow int64 length",
			input:  "0-100/9223372036854775808",
			wantOK: false,
		},
		{
			name:   "unsatisfied range",
			input:  "*/100",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseBytesRangeResp(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("parseBytesRangeResp(%q) ok = %t, want %t", tt.input, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if got.First != tt.want.First {
				t.Errorf(
					"parseBytesRangeResp(%q) First = %d, want %d",
					tt.input, got.First, tt.want.First,
				)
			}
			if got.Last != tt.want.Last {
				t.Errorf(
					"parseBytesRangeResp(%q) Last = %d, want %d",
					tt.input, got.Last, tt.want.Last,
				)
			}
			if got.CompleteLength != tt.want.CompleteLength {
				t.Errorf(
					"parseBytesRangeResp(%q) = %d, want %d",
					tt.input, got.CompleteLength, tt.want.CompleteLength,
				)
			}
		})
	}
}
