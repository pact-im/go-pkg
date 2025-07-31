// Package dockernames provides a names.Namer implementation that uses the
// Docker’s namesgenerator package.
package dockernames

import (
	"context"
	"strings"

	"github.com/go-x-pkg/namesgenerator"

	"go.pact.im/x/names"
)

var _ interface {
	names.NamerBuilder
	names.Namer
} = (*Namer)(nil)

// Namer provides names using the Docker’s namesgenerator package. It implements
// both names.Namer and names.NamerBuilder interfaces.
type Namer struct{}

// New returns a new Namer instance that uses Docker’s namesgenerator package.
func New() *Namer {
	return (*Namer)(nil)
}

// Build implements the names.NamerBuilder interface.
func (n *Namer) Build() names.Namer {
	return n
}

// Name implements the names.Namer interface.
func (n *Namer) Name(_ context.Context) (string, error) {
	name := namesgenerator.GetRandomName(0)
	return strings.Replace(name, "-", " ", 1), nil
}
