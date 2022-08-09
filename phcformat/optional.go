package phcformat

// OptionalString represents a string that may be unset. It is used instead of
// pointers to avoid allocations.
type OptionalString struct {
	// IsSet is true if Value contains a string value.
	IsSet bool
	// Value is the value of string. It is ignored if IsSet is false.
	Value string
}

// String returns an OptionalString instance with the given value.
func String(s string) OptionalString {
	return OptionalString{true, s}
}
