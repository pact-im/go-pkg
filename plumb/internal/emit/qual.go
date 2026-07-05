package emit

import (
	"fmt"
	"go/types"
	"slices"
	"strings"
)

// qualifier renders go/types types and symbols as Go source, deciding what is
// unqualified (the destination package) and what is imported and qualified.
//
// It runs in two modes. In recording mode (built with newRecording) it captures
// every package a render would reference, so the import block and aliases can be
// computed; the strings it returns in that mode are throwaway. In final mode
// (built with newFinal) it returns the chosen qualifier for each package.
type qualifier struct {
	destPath string
	record   map[string]*types.Package // path → package (recording mode)
	alias    map[string]string         // path → chosen qualifier (final mode)
}

// newRecording returns a qualifier that records every package referenced when
// rendering, for computing the import block.
func newRecording(destPath string) *qualifier {
	return &qualifier{destPath: destPath, record: map[string]*types.Package{}}
}

// newFinal returns a qualifier that renders with the chosen per-package alias.
func newFinal(destPath string, alias map[string]string) *qualifier {
	return &qualifier{destPath: destPath, alias: alias}
}

// qualify is the go/types.Qualifier used with types.TypeString.
func (q *qualifier) qualify(p *types.Package) string {
	if p == nil || p.Path() == q.destPath {
		return ""
	}
	if q.record != nil {
		q.record[p.Path()] = p
	}
	if q.alias != nil {
		// Every foreign package the render reaches must have been recorded by the
		// earlier import pass; a miss means that pass and this one drifted, which
		// would silently emit an unqualified, unimported name. Fail loud instead.
		a, ok := q.alias[p.Path()]
		if !ok {
			panic(fmt.Sprintf("plumb: package %q referenced by the render but not recorded for import", p.Path()))
		}
		return a
	}
	return p.Name()
}

// typeString renders t as Go source under this qualifier.
func (q *qualifier) typeString(t types.Type) string {
	return types.TypeString(t, q.qualify)
}

// recordedPackages returns the recorded packages sorted by import path.
func (q *qualifier) recordedPackages() []*types.Package {
	pkgs := make([]*types.Package, 0, len(q.record))
	for _, p := range q.record {
		pkgs = append(pkgs, p)
	}
	slices.SortFunc(pkgs, func(a, b *types.Package) int { return strings.Compare(a.Path(), b.Path()) })
	return pkgs
}

// aliasNames returns the chosen import qualifiers (final mode), so a caller can
// reserve them as taken identifiers when allocating local names.
func (q *qualifier) aliasNames() []string {
	out := make([]string, 0, len(q.alias))
	for _, a := range q.alias {
		out = append(out, a)
	}
	return out
}

// objQual renders a package-level object reference (function, variable,
// constant), qualified if it lives outside the destination.
func (q *qualifier) objQual(obj types.Object) string {
	pre := q.qualify(obj.Pkg())
	if pre == "" {
		return obj.Name()
	}
	return pre + "." + obj.Name()
}
