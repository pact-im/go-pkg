package gotypes

import (
	"fmt"
	"go/types"
)

// WalkNamed visits every named type reachable from t, invoking visit for each type
// name. It does not descend into a type parameter's constraint; a caller that needs
// the constraint walks it separately. When visit returns false for a name, the walk does
// not descend into that name's type arguments (it prunes below the node); it
// still visits the rest of the tree, so this controls descent, not whole-walk
// termination. A caller that wants to stop entirely self-guards inside visit.
func WalkNamed(t types.Type, visit func(*types.TypeName) bool) {
	walkNamed(t, &Set[types.Type]{}, visit)
}

// walkNamed is WalkNamed's recursion, threading the cycle-guard seen set that the
// exported entry point seeds fresh.
func walkNamed(t types.Type, seen *Set[types.Type], visit func(*types.TypeName) bool) {
	// Rendering names an alias, not its expansion, so the alias's own name is
	// what reachability must judge, and it is judged before the seen check:
	// the seen set keys on types.Identical, under which an alias equals its
	// target, so whichever spelling were walked first would swallow the other
	// (an exported alias hiding its unexported target, or vice versa). Aliases
	// bypass the dedup entirely; they cannot cycle (the type checker rejects
	// alias cycles), so the walk still covers only the finite spelling tree.
	// Only the rendered type arguments are walked, never Unalias(u).
	if u, ok := t.(*types.Alias); ok {
		if !visit(u.Obj()) {
			return
		}
		if args := u.TypeArgs(); args != nil {
			for t := range args.Types() {
				walkNamed(t, seen, visit)
			}
		}
		return
	}
	if !seen.Add(t) {
		return
	}
	switch u := t.(type) {
	case *types.Named:
		if !visit(u.Obj()) {
			return
		}
		if args := u.TypeArgs(); args != nil {
			for t := range args.Types() {
				walkNamed(t, seen, visit)
			}
		}
	case *types.Pointer:
		walkNamed(u.Elem(), seen, visit)
	case *types.Slice:
		walkNamed(u.Elem(), seen, visit)
	case *types.Array:
		walkNamed(u.Elem(), seen, visit)
	case *types.Chan:
		walkNamed(u.Elem(), seen, visit)
	case *types.Map:
		walkNamed(u.Key(), seen, visit)
		walkNamed(u.Elem(), seen, visit)
	case *types.Signature:
		walkTuple(u.Params(), seen, visit)
		walkTuple(u.Results(), seen, visit)
	case *types.Struct:
		for field := range u.Fields() {
			walkNamed(field.Type(), seen, visit)
		}
	case *types.Interface:
		for etyp := range u.EmbeddedTypes() {
			walkNamed(etyp, seen, visit)
		}
		for method := range u.ExplicitMethods() {
			walkNamed(method.Type(), seen, visit)
		}
	case *types.Union:
		for term := range u.Terms() {
			walkNamed(term.Type(), seen, visit)
		}
	case *types.Basic:
		// A predeclared or basic type names no type; nothing to visit or descend.
	case *types.TypeParam:
		// constraints are handled separately for lifted params
	default:
		panic(fmt.Sprintf("plumb: WalkNamed: unhandled type kind %T", t))
	}
}

func walkTuple(tup *types.Tuple, seen *Set[types.Type], visit func(*types.TypeName) bool) {
	for v := range tup.Variables() {
		walkNamed(v.Type(), seen, visit)
	}
}

// ContainsInvalid reports whether t mentions the invalid type anywhere, meaning
// it failed to type-check and plumb must not render it.
func ContainsInvalid(t types.Type) bool {
	return typeContains(t, &Set[types.Type]{}, func(x types.Type) bool {
		b, ok := x.(*types.Basic)
		return ok && b.Kind() == types.Invalid
	})
}

// ContainsTypeParam reports whether t mentions any type parameter anywhere:
// the test for a pinning that cannot survive a resolution restart, since lifted
// parameters are per-run objects.
func ContainsTypeParam(t types.Type) bool {
	return typeContains(t, &Set[types.Type]{}, func(x types.Type) bool {
		_, ok := x.(*types.TypeParam)
		return ok
	})
}

// typeContains walks t (guarding against cycles) and reports whether pred holds
// for any node. The switch is exhaustive over every types.Type kind: leaf kinds
// terminate explicitly and any unhandled kind panics, never silently reports
// false. A silent false here is precisely how an unrendered invalid type would
// slip past ContainsInvalid and reach the emitter.
func typeContains(t types.Type, seen *Set[types.Type], pred func(types.Type) bool) bool {
	t = types.Unalias(t)
	if !seen.Add(t) {
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

func tupleContains(tup *types.Tuple, seen *Set[types.Type], pred func(types.Type) bool) bool {
	for v := range tup.Variables() {
		if typeContains(v.Type(), seen, pred) {
			return true
		}
	}
	return false
}

// universeNames is the set of predeclared identifiers, backing IsUniverseName.
var universeNames = func() map[string]bool {
	m := map[string]bool{}
	for _, n := range types.Universe.Names() {
		m[n] = true
	}
	return m
}()

// IsUniverseName reports whether name is a predeclared identifier.
func IsUniverseName(name string) bool { return universeNames[name] }
