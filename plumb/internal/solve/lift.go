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

// liftTypeParams uses supplied bindings and lifts the remaining parameters.
// Constraints are rewritten after all bindings are known.
func (s *solver) liftTypeParams(p *discover.Provider, bind map[*types.TypeParam]types.Type) []types.Type {
	count := p.Rparams.Len() + p.Tparams.Len()
	targs := make([]types.Type, 0, count)
	resolved := make(map[*types.TypeParam]types.Type, count)
	var fresh []*types.TypeParam
	for param := range p.TypeParams {
		arg, bound := bind[param]
		if !bound {
			arg = s.liftOne(param, p)
			fresh = append(fresh, param)
		}
		resolved[param] = arg
		targs = append(targs, arg)
	}
	for _, param := range fresh {
		lifted := resolved[param].(*types.TypeParam)
		if c := param.Constraint(); c != nil {
			lifted.SetConstraint(gotypes.Subst(s.ctxt, c, resolved))
		}
	}
	return targs
}

// liftOne creates a fresh injector type parameter standing in for an unpinned
// template parameter. Its name is chosen to avoid the destination’s identifiers
// and the predeclared names the body emits, so it never shadows anything the
// generated body references unqualified.
func (s *solver) liftOne(orig *types.TypeParam, p *discover.Provider) *types.TypeParam {
	name := s.freshLiftedName(orig.Obj().Name())
	// Globally-unique live lifted names are the whole basis of the solver’s
	// determinism: they are the only source of a CmpType weak-order tie (two
	// lifted parameters compare equal iff their names match; see gotypes.CmpType),
	// and a stable sort over the map-ordered demand/input/output slices would
	// otherwise be load-order-dependent. freshLiftedName guarantees uniqueness, so
	// a duplicate here means that guarantee (or rollbackLifts’ name release)
	// regressed; fail loud at the source rather than emit nondeterministic code.
	if s.liftedNames[name] {
		panic(fmt.Sprintf("plumb: lifted name %q reused; determinism relies on unique live lifted names", name))
	}
	tn := types.NewTypeName(token.NoPos, nil, name, nil)
	tp := types.NewTypeParam(tn, nil)
	tp.SetConstraint(orig.Constraint())
	s.liftedNames[name] = true
	s.liftedMeta = append(s.liftedMeta, liftedParam{tp: tp, origin: orig, pos: p.Pos})
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
// by provider position, then parameter declaration position.
func (s *solver) orderedLifted() []*types.TypeParam {
	meta := slices.SortedStableFunc(slices.Values(s.liftedMeta), func(a, b liftedParam) int {
		return cmp.Or(
			diag.CmpPos(a.pos, b.pos),
			cmp.Compare(a.origin.Obj().Pos(), b.origin.Obj().Pos()),
		)
	})
	out := make([]*types.TypeParam, len(meta))
	for i, m := range meta {
		out[i] = m.tp
	}
	return out
}
