package gotypes

import (
	"fmt"
	"go/types"
)

// Subst returns t with every type parameter in m replaced by its mapped type.
// Only the nodes on the path to a substituted parameter are rebuilt; a subtree
// that mentions none of the mapped parameters is returned unchanged, so special
// types such as the predeclared comparable constraint and named cross-package
// interfaces keep their identity.
//
// ctxt deduplicates any generic types Subst must re-instantiate; pass the caller’s
// shared context so identical instantiations resolve to one instance. A nil ctxt
// is valid (a throwaway): correctness never depends on instance identity, because
// the result is rendered to source text and inspected structurally (go/types
// Identical, underlying-type assertions, reachability walks), never keyed by an
// instance pointer.
func Subst(ctxt *types.Context, t types.Type, m map[*types.TypeParam]types.Type) types.Type {
	return (&subster{ctxt: ctxt, m: m}).subst(t)
}

type subster struct {
	ctxt *types.Context
	m    map[*types.TypeParam]types.Type
}

func (s *subster) subst(t types.Type) types.Type {
	if t == nil || !s.mentions(t) {
		return t
	}
	// Past the guard the type mentions a substituted parameter, so it must be
	// rebuilt; collapse any alias to its target first (the alias name cannot
	// survive substitution of what it expands to).
	t = types.Unalias(t)
	switch u := t.(type) {
	case *types.TypeParam:
		if r, ok := s.m[u]; ok {
			return r
		}
		return u
	case *types.Pointer:
		return types.NewPointer(s.subst(u.Elem()))
	case *types.Slice:
		return types.NewSlice(s.subst(u.Elem()))
	case *types.Array:
		return types.NewArray(s.subst(u.Elem()), u.Len())
	case *types.Chan:
		return types.NewChan(u.Dir(), s.subst(u.Elem()))
	case *types.Map:
		return types.NewMap(s.subst(u.Key()), s.subst(u.Elem()))
	case *types.Union:
		terms := make([]*types.Term, u.Len())
		for i := range u.Len() {
			tm := u.Term(i)
			terms[i] = types.NewTerm(tm.Tilde(), s.subst(tm.Type()))
		}
		return types.NewUnion(terms)
	case *types.Interface:
		var methods []*types.Func
		for mth := range u.ExplicitMethods() {
			sig := s.subst(mth.Type()).(*types.Signature)
			methods = append(methods, types.NewFunc(mth.Pos(), mth.Pkg(), mth.Name(), sig))
		}
		var embeds []types.Type
		for etyp := range u.EmbeddedTypes() {
			embeds = append(embeds, s.subst(etyp))
		}
		ni := types.NewInterfaceType(methods, embeds)
		ni.Complete()
		return ni
	case *types.Signature:
		return types.NewSignatureType(nil, nil, nil,
			s.substTuple(u.Params()), s.substTuple(u.Results()), u.Variadic())
	case *types.Struct:
		n := u.NumFields()
		fields := make([]*types.Var, n)
		tags := make([]string, n)
		for i := range n {
			f := u.Field(i)
			fields[i] = types.NewField(f.Pos(), f.Pkg(), f.Name(), s.subst(f.Type()), f.Embedded())
			tags[i] = u.Tag(i)
		}
		return types.NewStruct(fields, tags)
	case *types.Named:
		args := u.TypeArgs()
		if args.Len() == 0 { // TypeList.Len is nil-safe
			return u
		}
		na := make([]types.Type, args.Len())
		for i := range args.Len() {
			na[i] = s.subst(args.At(i))
		}
		// validate=false: with a generic origin and the right number of type
		// arguments (both guaranteed here: u is an instantiated Named and na has one
		// entry per original arg), Instantiate cannot return an error. A non-nil error
		// would mean that invariant broke; returning the un-substituted type would
		// leak the template’s own parameters into a lifted constraint, so fail loud.
		inst, err := types.Instantiate(s.ctxt, u.Origin(), na, false)
		if err != nil {
			panic(fmt.Sprintf("plumb: Subst: re-instantiating %s failed: %v", u, err))
		}
		return inst
	default:
		// Unreachable: past the mentions guard t references a type parameter, so
		// it cannot be a kind (Basic, Tuple) that mentions none.
		panic(fmt.Sprintf("plumb: subst on unsupported type %T", t))
	}
}

func (s *subster) substTuple(tup *types.Tuple) *types.Tuple {
	vars := make([]*types.Var, tup.Len())
	for i := range tup.Len() {
		v := tup.At(i)
		vars[i] = types.NewVar(v.Pos(), v.Pkg(), v.Name(), s.subst(v.Type()))
	}
	return types.NewTuple(vars...)
}

// mentions reports whether t references any type parameter in the substitution
// map s.m. It shares typeContains’s exhaustive type walk, testing membership in
// s.m at each type-parameter leaf rather than an explicit params set.
func (s *subster) mentions(t types.Type) bool {
	return typeContains(t, &Set[types.Type]{}, func(x types.Type) bool {
		tp, ok := x.(*types.TypeParam)
		if !ok {
			return false
		}
		_, ok = s.m[tp]
		return ok
	})
}
