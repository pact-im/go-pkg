package gotypes

import (
	"fmt"
	"go/types"
	"slices"
)

// Subst returns t with every type parameter in m replaced by its mapped type.
// Only the nodes on the path to a substituted parameter are rebuilt; a subtree
// that mentions none of the mapped parameters is returned unchanged, so special
// types such as the predeclared comparable constraint and named cross-package
// interfaces keep their identity. Aliases keep their names, with substitution
// applied to their type arguments.
//
// ctxt deduplicates any generic types Subst must re-instantiate; pass the caller’s
// shared context so identical instantiations resolve to one instance. A nil ctxt
// is valid (a throwaway): correctness never depends on instance identity, because
// the result is rendered to source text and inspected structurally (go/types
// Identical, underlying-type assertions, reachability walks), never keyed by an
// instance pointer.
func Subst(ctxt *types.Context, t types.Type, m map[*types.TypeParam]types.Type) types.Type {
	if len(m) == 0 {
		return t
	}
	return (&subster{ctxt: ctxt, m: m}).subst(t)
}

type subster struct {
	ctxt *types.Context
	m    map[*types.TypeParam]types.Type
}

func (s *subster) subst(t types.Type) types.Type {
	switch u := t.(type) {
	case nil, *types.Basic:
		return t
	case *types.TypeParam:
		r, ok := s.m[u]
		if !ok {
			return u
		}
		return r
	case *types.Pointer:
		elem := s.subst(u.Elem())
		if elem == u.Elem() {
			return u
		}
		return types.NewPointer(elem)
	case *types.Slice:
		elem := s.subst(u.Elem())
		if elem == u.Elem() {
			return u
		}
		return types.NewSlice(elem)
	case *types.Array:
		elem := s.subst(u.Elem())
		if elem == u.Elem() {
			return u
		}
		return types.NewArray(elem, u.Len())
	case *types.Chan:
		elem := s.subst(u.Elem())
		if elem == u.Elem() {
			return u
		}
		return types.NewChan(u.Dir(), elem)
	case *types.Map:
		key, elem := s.subst(u.Key()), s.subst(u.Elem())
		if key == u.Key() && elem == u.Elem() {
			return u
		}
		return types.NewMap(key, elem)
	case *types.Union:
		return s.substUnion(u)
	case *types.Interface:
		return s.substInterface(u)
	case *types.Signature:
		params, results := s.substTuple(u.Params()), s.substTuple(u.Results())
		if params == u.Params() && results == u.Results() {
			return u
		}
		return types.NewSignatureType(nil, nil, nil,
			params, results, u.Variadic())
	case *types.Struct:
		return s.substStruct(u)
	case *types.Named:
		return s.substInstance(u, u.Origin(), u.TypeArgs())
	case *types.Alias:
		return s.substInstance(u, u.Origin(), u.TypeArgs())
	default:
		// Unreachable: tuples are handled by substTuple, and every
		// other type kind is covered above.
		panic(fmt.Sprintf("plumb: subst on unsupported type %T", t))
	}
}

func (s *subster) substUnion(u *types.Union) *types.Union {
	changed := false

	terms := slices.Collect(u.Terms())
	for i, tm := range terms {
		typ := s.subst(tm.Type())
		if typ == tm.Type() {
			continue
		}
		changed = true
		terms[i] = types.NewTerm(tm.Tilde(), typ)
	}

	if !changed {
		return u
	}

	return types.NewUnion(terms)
}

func (s *subster) substInterface(u *types.Interface) *types.Interface {
	changed := false

	methods := slices.Collect(u.ExplicitMethods())
	for i, mth := range methods {
		sig := s.subst(mth.Type()).(*types.Signature)
		if sig == mth.Type() {
			continue
		}
		changed = true
		methods[i] = types.NewFunc(mth.Pos(), mth.Pkg(), mth.Name(), sig)
	}

	embeds := slices.Collect(u.EmbeddedTypes())
	for i, etyp := range embeds {
		typ := s.subst(etyp)
		if typ == etyp {
			continue
		}
		changed = true
		embeds[i] = typ
	}

	if !changed {
		return u
	}

	ni := types.NewInterfaceType(methods, embeds)
	if u.IsImplicit() {
		ni.MarkImplicit()
	}
	ni.Complete()
	return ni
}

func (s *subster) substStruct(u *types.Struct) *types.Struct {
	changed := false

	fields := slices.Collect(u.Fields())
	for i, f := range fields {
		typ := s.subst(f.Type())
		if typ == f.Type() {
			continue
		}
		changed = true
		fields[i] = types.NewField(f.Pos(), f.Pkg(), f.Name(), typ, f.Embedded())
	}

	if !changed {
		return u
	}

	tags := make([]string, len(fields))
	for i := range fields {
		tags[i] = u.Tag(i)
	}
	return types.NewStruct(fields, tags)
}

func (s *subster) substTuple(tup *types.Tuple) *types.Tuple {
	changed := false

	vars := slices.Collect(tup.Variables())
	for i, v := range vars {
		typ := s.subst(v.Type())
		if typ == v.Type() {
			continue
		}
		changed = true
		vars[i] = types.NewVar(v.Pos(), v.Pkg(), v.Name(), typ)
	}

	if !changed {
		return tup
	}

	return types.NewTuple(vars...)
}

func (s *subster) substInstance(t, origin types.Type, args *types.TypeList) types.Type {
	changed := false

	targs := slices.Collect(args.Types())
	for i, arg := range targs {
		typ := s.subst(arg)
		if typ == arg {
			continue
		}
		changed = true
		targs[i] = typ
	}

	if !changed {
		return t
	}

	// The origin and argument count come from an existing instance. With
	// validation disabled, instantiation cannot fail unless that invariant breaks.
	inst, err := types.Instantiate(s.ctxt, origin, targs, false)
	if err != nil {
		panic(fmt.Sprintf("plumb: Subst: re-instantiating %s failed: %v", t, err))
	}
	return inst
}
