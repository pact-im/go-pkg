package logs

import (
	"errors"
	"log/slog"
	"testing"
)

func TestSliceGroup(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		input []slog.Value
		want  slog.Attr
	}{
		{
			name:  "empty",
			key:   "items",
			input: nil,
			want:  slog.GroupAttrs("items"),
		},
		{
			name:  "single",
			key:   "data",
			input: []slog.Value{slog.IntValue(42)},
			want:  slog.GroupAttrs("data", slog.Int("0", 42)),
		},
		{
			name: "multiple",
			key:  "data",
			input: []slog.Value{
				slog.IntValue(1),
				slog.StringValue("two"),
			},
			want: slog.GroupAttrs("data",
				slog.Int("0", 1),
				slog.String("1", "two"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceGroup(tt.key, tt.input...)
			if !slogAttrEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

type testError struct {
	msg string
}

func (e testError) Error() string { return e.msg }

func TestError(t *testing.T) {
	sentinelError := errors.New("boom")
	tests := []struct {
		name  string
		input error
		want  slog.Attr
	}{
		{
			name:  "sentinel error",
			input: sentinelError,
			want: slog.GroupAttrs("error",
				slog.String("type", "*errors.errorString"),
				slog.Any("value", sentinelError),
			),
		},
		{
			name:  "typed error",
			input: testError{msg: "fail"},
			want: slog.GroupAttrs("error",
				slog.String("type", "logs.testError"),
				slog.Any("value", testError{msg: "fail"}),
			),
		},
		{
			name:  "nil error",
			input: nil,
			want:  slog.Attr{Key: "error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Error(tt.input)
			if !slogAttrEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
