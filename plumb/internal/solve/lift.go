package solve

import (
	"cmp"
	"fmt"
	"go/token"
	"go/types"
	"slices"
	"strconv"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
)

// liftAll lifts every type parameter of an anchor template, returning the fresh
// parameters as type arguments.
func (s *solver) liftAll(p *discover.Provider) []types.Type {
	targs := make([]types.Type, p.Tparams.Len())
	fresh := make([]int, p.Tparams.Len())
	for i := range p.Tparams.Len() {
		targs[i] = s.liftOne(p.Tparams.At(i), p)
		fresh[i] = i
	}
	s.setLiftedConstraints(p, targs, fresh)
	return targs
}

// setLiftedConstraints rewrites the constraint of each freshly lifted parameter
// (the indices in fresh) so that references to the template's other parameters
// point at their resolved targets: concrete types for pinned parameters, the
// co-lifted fresh parameters for the rest. The original template parameters are
// not in scope in the generated header, so an inter-parameter constraint such as
// U interface{ ~[]T } must travel as ~[]<pin or lifted T>, not the verbatim T.
// Leaving it verbatim renders a dangling identifier (uncompilable) and makes
// types.Instantiate's validation reject the lift.
func (s *solver) setLiftedConstraints(p *discover.Provider, targs []types.Type, fresh []int) {
	m := make(map[*types.TypeParam]types.Type, p.Tparams.Len())
	for i := range p.Tparams.Len() {
		m[p.Tparams.At(i)] = targs[i]
	}
	for _, i := range fresh {
		lifted := targs[i].(*types.TypeParam)
		if c := p.Tparams.At(i).Constraint(); c != nil {
			lifted.SetConstraint(gotypes.Subst(s.ctxt, c, m))
		}
	}
}

// liftOne creates a fresh injector type parameter standing in for an unpinned
// template parameter. Its name is chosen to avoid the destination's identifiers
// and the predeclared names the body emits, so it never shadows anything the
// generated body references unqualified.
func (s *solver) liftOne(orig *types.TypeParam, p *discover.Provider) *types.TypeParam {
	name := s.freshLiftedName(orig.Obj().Name())
	// Globally-unique live lifted names are the whole basis of the solver's
	// determinism: they are the only source of a CmpType weak-order tie (two
	// lifted parameters compare equal iff their names match; see gotypes.CmpType),
	// and a stable sort over the map-ordered demand/input/output slices would
	// otherwise be load-order-dependent. freshLiftedName guarantees uniqueness, so
	// a duplicate here means that guarantee (or rollbackLifts' name release)
	// regressed; fail loud at the source rather than emit nondeterministic code.
	if s.liftedNames[name] {
		panic(fmt.Sprintf("plumb: lifted name %q reused; determinism relies on unique live lifted names", name))
	}
	tn := types.NewTypeName(token.NoPos, nil, name, nil)
	tp := types.NewTypeParam(tn, nil)
	tp.SetConstraint(orig.Constraint())
	s.liftedNames[name] = true
	s.liftedMeta = append(s.liftedMeta, liftedParam{tp: tp, pos: p.Pos, idx: orig.Index()})
	return tp
}

// rollbackLifts undoes the speculative lifts appended since mark: it both
// truncates liftedMeta and releases the names those lifts reserved in liftedNames.
// Releasing the names is what lets a later successful lift reuse a name a failed
// attempt had taken (e.g. T rather than T2); liftedMeta and liftedNames must be
// rolled back together or the two go out of sync.
func (s *solver) rollbackLifts(mark int) {
	for _, m := range s.liftedMeta[mark:] {
		delete(s.liftedNames, m.tp.Obj().Name())
	}
	s.liftedMeta = s.liftedMeta[:mark]
}

// freshLiftedName returns a valid identifier derived from base that collides
// with nothing the body references unqualified.
func (s *solver) freshLiftedName(base string) string {
	if base == "" || base == "_" {
		base = "T"
	}
	name := base
	for i := 2; s.avoid[name] || s.liftedNames[name]; i++ {
		name = base + strconv.Itoa(i)
	}
	return name
}

// orderedLifted returns the lifted parameters in deterministic header order:
// by the source position of their origin template, then parameter index.
func (s *solver) orderedLifted() []*types.TypeParam {
	meta := slices.SortedStableFunc(slices.Values(s.liftedMeta), func(a, b liftedParam) int {
		return cmp.Or(diag.CmpPos(a.pos, b.pos), cmp.Compare(a.idx, b.idx))
	})
	out := make([]*types.TypeParam, len(meta))
	for i, m := range meta {
		out[i] = m.tp
	}
	return out
}
