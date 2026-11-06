package gotypes

import (
	"fmt"
	"go/types"

	"golang.org/x/exp/collection/hashed"
)

// WalkObjects visits the objects named in t’s source representation: type names,
// struct fields, and explicit interface methods. It follows type arguments but
// not the definitions of named types or aliases. Type parameter constraints are
// walked separately. Returning false skips an object’s children; other branches
// are still visited.
func WalkObjects(t types.Type, visit func(types.Object) bool) {
	walkObjects(t, map[types.Type]bool{}, visit)
}

// walkObjects distinguishes type nodes even when types.Identical considers them
// equal, since their source representations can name different objects.
func walkObjects(t types.Type, seen map[types.Type]bool, visit func(types.Object) bool) {
	if seen[t] {
		return
	}
	seen[t] = true
	switch u := t.(type) {
	case *types.Named:
		if !visit(u.Obj()) {
			return
		}
		for arg := range u.TypeArgs().Types() {
			walkObjects(arg, seen, visit)
		}
	case *types.Alias:
		if !visit(u.Obj()) {
			return
		}
		for arg := range u.TypeArgs().Types() {
			walkObjects(arg, seen, visit)
		}
	case *types.Pointer:
		walkObjects(u.Elem(), seen, visit)
	case *types.Slice:
		walkObjects(u.Elem(), seen, visit)
	case *types.Array:
		walkObjects(u.Elem(), seen, visit)
	case *types.Chan:
		walkObjects(u.Elem(), seen, visit)
	case *types.Map:
		walkObjects(u.Key(), seen, visit)
		walkObjects(u.Elem(), seen, visit)
	case *types.Signature:
		walkTuple(u.Params(), seen, visit)
		walkTuple(u.Results(), seen, visit)
	case *types.Struct:
		for field := range u.Fields() {
			if !visit(field) {
				continue
			}
			walkObjects(field.Type(), seen, visit)
		}
	case *types.Interface:
		for etyp := range u.EmbeddedTypes() {
			walkObjects(etyp, seen, visit)
		}
		for method := range u.ExplicitMethods() {
			if !visit(method) {
				continue
			}
			walkObjects(method.Type(), seen, visit)
		}
	case *types.Union:
		for term := range u.Terms() {
			walkObjects(term.Type(), seen, visit)
		}
	case *types.Basic, *types.TypeParam:
		// Basic types name no objects; type parameter constraints are handled separately.
	default:
		panic(fmt.Sprintf("plumb: WalkObjects: unhandled type kind %T", t))
	}
}

func walkTuple(tup *types.Tuple, seen map[types.Type]bool, visit func(types.Object) bool) {
	for v := range tup.Variables() {
		walkObjects(v.Type(), seen, visit)
	}
}

// ContainsInvalid reports whether t mentions the invalid type anywhere, meaning
// it failed to type-check and plumb must not render it.
func ContainsInvalid(t types.Type) bool {
	return typeContains(t, hashed.NewSet(types.Hasher{}), func(x types.Type) bool {
		b, ok := x.(*types.Basic)
		return ok && b.Kind() == types.Invalid
	})
}

// ContainsTypeParam reports whether t mentions any type parameter anywhere:
// the test for a pinning that cannot survive a resolution restart, since lifted
// parameters are per-run objects.
func ContainsTypeParam(t types.Type) bool {
	return typeContains(t, hashed.NewSet(types.Hasher{}), func(x types.Type) bool {
		_, ok := x.(*types.TypeParam)
		return ok
	})
}

// typeContains walks t (guarding against cycles) and reports whether pred holds
// for any node. The switch is exhaustive over every types.Type kind: leaf kinds
// terminate explicitly and any unhandled kind panics, never silently reports
// false. A silent false here is precisely how an unrendered invalid type would
// slip past ContainsInvalid and reach the emitter.
func typeContains(t types.Type, seen *hashed.Set[types.Type], pred func(types.Type) bool) bool {
	t = types.Unalias(t)
	if !seen.Insert(t) {
		return false
	}
	if pred(t) {
		return true
	}
	switch u := t.(type) {
	case *types.Basic, *types.TypeParam:
		// Leaf types: nothing to descend into. A type parameter renders as its
		// name; its constraint is validated separately (see checks.go).
		return false
	case *types.Pointer:
		return typeContains(u.Elem(), seen, pred)
	case *types.Slice:
		return typeContains(u.Elem(), seen, pred)
	case *types.Array:
		return typeContains(u.Elem(), seen, pred)
	case *types.Chan:
		return typeContains(u.Elem(), seen, pred)
	case *types.Map:
		return typeContains(u.Key(), seen, pred) || typeContains(u.Elem(), seen, pred)
	case *types.Named:
		if args := u.TypeArgs(); args != nil {
			for t := range args.Types() {
				if typeContains(t, seen, pred) {
					return true
				}
			}
		}
		return false
	case *types.Signature:
		return tupleContains(u.Params(), seen, pred) || tupleContains(u.Results(), seen, pred)
	case *types.Struct:
		for field := range u.Fields() {
			if typeContains(field.Type(), seen, pred) {
				return true
			}
		}
		return false
	case *types.Interface:
		for etyp := range u.EmbeddedTypes() {
			if typeContains(etyp, seen, pred) {
				return true
			}
		}
		for method := range u.ExplicitMethods() {
			if typeContains(method.Type(), seen, pred) {
				return true
			}
		}
		return false
	case *types.Union:
		for term := range u.Terms() {
			if typeContains(term.Type(), seen, pred) {
				return true
			}
		}
		return false
	default:
		panic(fmt.Sprintf("plumb: typeContains: unhandled type kind %T", t))
	}
}

func tupleContains(tup *types.Tuple, seen *hashed.Set[types.Type], pred func(types.Type) bool) bool {
	for v := range tup.Variables() {
		if typeContains(v.Type(), seen, pred) {
			return true
		}
	}
	return false
}

// IsUniverseName reports whether name is a predeclared identifier.
func IsUniverseName(name string) bool {
	return types.Universe.Lookup(name) != nil
}
