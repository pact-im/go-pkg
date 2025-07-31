// Package names provides an abstraction for generating short and human-readable
// pseudo-random names for objects.
package names

import (
	"context"
)

// NamerBuilder constructs Namer instances.
type NamerBuilder interface {
	Build() Namer
}

// Namer generates short and human-readable names for objects.
type Namer interface {
	Name(ctx context.Context) (string, error)
}

// noopNamerBuilder is a NamerBuilder for a Namer implementation that returns
// empty strings.
type noopNamerBuilder struct{}

// NewNoopNamerBuilder returns a NamerBuilder for Namer that returns empty
// strings.
func NewNoopNamerBuilder() NamerBuilder {
	return (*noopNamerBuilder)(nil)
}

// Build implements the NamerBuilder interface.
func (*noopNamerBuilder) Build() Namer {
	return NewNoopNamer()
}

// noopNamer is a Namer that returns empty strings.
type noopNamer struct{}

// NewNoopNamer returns a Namer that returns empty strings.
func NewNoopNamer() Namer {
	return (*noopNamer)(nil)
}

// Name implements the Namer interface.
func (*noopNamer) Name(_ context.Context) (string, error) {
	return "", nil
}
