package gotypes

import "go/types"

// Unify matches a pattern type (which may mention the template’s type
// parameters) against a concrete demand type, recording the type-parameter
// bindings it discovers. It reports whether the match succeeds. Only the type
// parameters in params are treated as variables; every other type node must be
// structurally identical.
func Unify(pattern, concrete types.Type, params map[*types.TypeParam]bool, bind map[*types.TypeParam]types.Type) bool {
	pattern = types.Unalias(pattern)
	concrete = types.Unalias(concrete)

	if tp, ok := pattern.(*types.TypeParam); ok && params[tp] {
		if prev, seen := bind[tp]; seen {
			return types.Identical(prev, concrete)
		}
		bind[tp] = concrete
		return true
	}

	// If the pattern mentions no template parameter, identity is enough.
	if !MentionsParams(pattern, params) {
		return types.Identical(pattern, concrete)
	}

	switch p := pattern.(type) {
	case *types.Pointer:
		c, ok := concrete.(*types.Pointer)
		return ok && Unify(p.Elem(), c.Elem(), params, bind)
	case *types.Slice:
		c, ok := concrete.(*types.Slice)
		return ok && Unify(p.Elem(), c.Elem(), params, bind)
	case *types.Array:
		c, ok := concrete.(*types.Array)
		return ok && p.Len() == c.Len() && Unify(p.Elem(), c.Elem(), params, bind)
	case *types.Chan:
		c, ok := concrete.(*types.Chan)
		return ok && p.Dir() == c.Dir() && Unify(p.Elem(), c.Elem(), params, bind)
	case *types.Map:
		c, ok := concrete.(*types.Map)
		return ok && Unify(p.Key(), c.Key(), params, bind) && Unify(p.Elem(), c.Elem(), params, bind)
	case *types.Named:
		c, ok := concrete.(*types.Named)
		if !ok || !types.Identical(p.Origin(), c.Origin()) {
			return false
		}
		pa, ca := p.TypeArgs(), c.TypeArgs()
		if pa.Len() != ca.Len() { // TypeList.Len is nil-safe
			return false
		}
		for i := range pa.Len() {
			if !Unify(pa.At(i), ca.At(i), params, bind) {
				return false
			}
		}
		return true
	case *types.Signature:
		c, ok := concrete.(*types.Signature)
		return ok && p.Variadic() == c.Variadic() &&
			unifyTuple(p.Params(), c.Params(), params, bind) &&
			unifyTuple(p.Results(), c.Results(), params, bind)
	case *types.Struct:
		c, ok := concrete.(*types.Struct)
		if !ok || p.NumFields() != c.NumFields() {
			return false
		}
		for i := range p.NumFields() {
			pf, cf := p.Field(i), c.Field(i)
			if pf.Id() != cf.Id() || pf.Embedded() != cf.Embedded() || p.Tag(i) != c.Tag(i) {
				return false
			}
			if !Unify(pf.Type(), cf.Type(), params, bind) {
				return false
			}
		}
		return true
	case *types.Interface:
		// Value-type interfaces are basic (method sets only); compare the
		// complete, name-sorted method sets and unify each signature.
		c, ok := concrete.(*types.Interface)
		if !ok || p.NumMethods() != c.NumMethods() {
			return false
		}
		for i := range p.NumMethods() {
			pm, cm := p.Method(i), c.Method(i)
			if pm.Id() != cm.Id() || !Unify(pm.Type(), cm.Type(), params, bind) {
				return false
			}
		}
		return true
	default:
		// Any remaining form (e.g. a union, which only appears in a constraint,
		// never a value type) is matched by identity, which fails safely.
		return types.Identical(pattern, concrete)
	}
}

// unifyTuple unifies two tuples element-wise; they must have equal arity.
func unifyTuple(p, c *types.Tuple, params map[*types.TypeParam]bool, bind map[*types.TypeParam]types.Type) bool {
	if p.Len() != c.Len() {
		return false
	}
	for i := range p.Len() {
		if !Unify(p.At(i).Type(), c.At(i).Type(), params, bind) {
			return false
		}
	}
	return true
}

// MentionsParams reports whether t mentions any of the given type parameters.
// It shares typeContains’s exhaustive type walk (which panics on an unhandled
// go/types kind rather than silently reporting false), testing membership in
// params at each type-parameter leaf.
func MentionsParams(t types.Type, params map[*types.TypeParam]bool) bool {
	return typeContains(t, &Set[types.Type]{}, func(x types.Type) bool {
		tp, ok := x.(*types.TypeParam)
		return ok && params[tp]
	})
}

// MentionsParam reports whether t mentions the single type parameter tp. It is
// the singular form of MentionsParams over the same typeContains walk, for
// callers testing one parameter at a time without building a one-element set.
func MentionsParam(t types.Type, tp *types.TypeParam) bool {
	return typeContains(t, &Set[types.Type]{}, func(x types.Type) bool {
		p, ok := x.(*types.TypeParam)
		return ok && p == tp
	})
}

// IsBareTypeParam reports whether t is exactly one of the given type parameters.
func IsBareTypeParam(t types.Type, params map[*types.TypeParam]bool) bool {
	tp, ok := types.Unalias(t).(*types.TypeParam)
	return ok && params[tp]
}

// IsPointerToBareTypeParam reports whether t is *T for a type parameter T in
// params. Such a result matches every pointer demand, and through the
// value/pointer bridge every value demand too, so, like a bare T, it cannot be
// a demand-driven producer.
func IsPointerToBareTypeParam(t types.Type, params map[*types.TypeParam]bool) bool {
	p, ok := types.Unalias(t).(*types.Pointer)
	return ok && IsBareTypeParam(p.Elem(), params)
}

// TypeDepth returns a structural nesting depth used to bound non-terminating
// generic instantiation.
func TypeDepth(t types.Type) int {
	return typeDepthRec(t, &Set[types.Type]{})
}

func typeDepthRec(t types.Type, seen *Set[types.Type]) int {
	t = types.Unalias(t)
	if !seen.Add(t) {
		return 0
	}
	switch u := t.(type) {
	case *types.Pointer:
		return 1 + typeDepthRec(u.Elem(), seen)
	case *types.Slice:
		return 1 + typeDepthRec(u.Elem(), seen)
	case *types.Array:
		return 1 + typeDepthRec(u.Elem(), seen)
	case *types.Chan:
		return 1 + typeDepthRec(u.Elem(), seen)
	case *types.Map:
		return 1 + max(typeDepthRec(u.Key(), seen), typeDepthRec(u.Elem(), seen))
	case *types.Named:
		d := 0
		if args := u.TypeArgs(); args != nil {
			for t := range args.Types() {
				d = max(d, typeDepthRec(t, seen))
			}
		}
		return 1 + d
	case *types.Signature:
		// A template can grow through a function type (Box[func(T)] → Box[T]), so
		// the bound must see through signatures, not treat them as flat leaves.
		d := 0
		for v := range u.Params().Variables() {
			d = max(d, typeDepthRec(v.Type(), seen))
		}
		for v := range u.Results().Variables() {
			d = max(d, typeDepthRec(v.Type(), seen))
		}
		return 1 + d
	case *types.Struct:
		d := 0
		for f := range u.Fields() {
			d = max(d, typeDepthRec(f.Type(), seen))
		}
		return 1 + d
	case *types.Interface:
		d := 0
		for m := range u.ExplicitMethods() {
			d = max(d, typeDepthRec(m.Type(), seen))
		}
		for e := range u.EmbeddedTypes() {
			d = max(d, typeDepthRec(e, seen))
		}
		return 1 + d
	}
	return 1
}
