// Package monikernames provides a names.Namer implementation that uses the
// github.com/technosophos/moniker package.
package monikernames

import (
	"context"
	"sync"

	"github.com/technosophos/moniker"

	"go.pact.im/x/names"
)

var (
	_ names.NamerBuilder = (*Builder)(nil)
	_ names.Namer        = (*Namer)(nil)
)

// Builder constructs Namer instances for moniker package that are safe for
// concurrent use.
type Builder struct{}

// New returns a new NamerBuilder that constructs Namer instances for the
// moniker package.
func New() *Builder {
	return (*Builder)(nil)
}

// Build implements the names.NamerBuilder interface.
func (n *Builder) Build() names.Namer {
	return NewNamer(moniker.New())
}

// Namer is an adapter for moniker’s Namer type that is safe for concurrent
// use.
type Namer struct {
	mu sync.Mutex
	n  moniker.Namer
}

// NewNamer returns a new names.Namer adapter for the given Namer.
func NewNamer(n moniker.Namer) *Namer {
	return &Namer{
		n: n,
	}
}

// Name implements the names.Namer interface.
func (n *Namer) Name(_ context.Context) (string, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.n.Name(), nil
}
