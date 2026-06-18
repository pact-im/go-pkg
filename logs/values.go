package logs

import (
	"log/slog"
	"reflect"
	"strconv"
)

var _ slog.LogValuer = sliceGroup{}

// sliceGroup is a slice that implements [slog.LogValuer]. Each element is
// logged as an attribute within a group, with numeric keys "0", "1", and so on.
type sliceGroup []slog.Value

// LogValue returns a group value whose attributes are the elements of l,
// keyed by their index.
func (l sliceGroup) LogValue() slog.Value {
	attrs := make([]slog.Attr, len(l))
	for i, v := range l {
		attrs[i] = slog.Attr{
			Key:   strconv.Itoa(i),
			Value: v,
		}
	}
	return slog.GroupValue(attrs...)
}

// SliceGroupValue returns an [slog.Value] that logs values as a group with
// numeric keys "0", "1", and so on.
func SliceGroupValue(values ...slog.Value) slog.Value {
	return slog.AnyValue(sliceGroup(values))
}

// SliceGroup returns an [slog.Attr] that logs values as a group with numeric keys
// "0", "1", and so on.
func SliceGroup(key string, values ...slog.Value) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: SliceGroupValue(values...),
	}
}

// typedError wraps an error as a group with attributes "type" (the type
// string from reflection) and "value" (the error itself).
type typedError struct {
	Error error
}

// LogValue returns a group with "type" and "value" attributes
// for e.Error: the type string from reflection, and the error itself.
func (e typedError) LogValue() slog.Value {
	typ := reflect.TypeOf(e.Error)
	if typ == nil {
		return slog.Value{}
	}
	return slog.GroupValue(
		slog.String("type", typ.String()),
		slog.Any("value", e.Error),
	)
}

// Error returns an [slog.Attr] with key "error" for the given error. The
// value is a group with "type" (the error’s type string from reflection) and
// "value" (the error itself).
func Error(e error) slog.Attr {
	return NamedError("error", e)
}

// NamedError returns an [slog.Attr] with the given key for the error. The
// value is a group with "type" (the error’s type string from reflection) and
// "value" (the error itself).
func NamedError(key string, e error) slog.Attr {
	return slog.Any(key, typedError{e})
}
