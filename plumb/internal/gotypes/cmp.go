package gotypes

import (
	"cmp"
	"fmt"
	"go/types"
	"strings"
)

// CmpType is a deterministic, environment-independent order over types, used to
// sort type slices into a stable order, never for identity (that goes through
// typeutil.Map) and never for output (that goes through a qualifier). Named types
// are ordered by package path + name rather than pointer, and every level is
// unaliased before comparison, so the order never depends on pointer addresses,
// load order, or whether a type was first seen through an alias. It is the
// ordering authority for types, so gotypes never has to serialize a type to a
// string except when rendering output or diagnostics.
//
// It is a WEAK order, not a structural total order: types that are not
// types.Identical can still compare equal. Named types are keyed on package
// path + name + type args (not their full definition), and type parameters on
// name + index (see the TypeParam case). Determinism therefore rests entirely on
// the absence of ties, not on the sort: every sorted slice holds types drawn from
// identity-keyed sets of provider-signature types that already differ in those
// keys, so no two elements ever compare equal, and the determinism tests guard
// that empirically: solve's TestNoCmpTypeTiesInSortedSlices asserts no tie
// reaches those sort sites. A structural string tiebreaker cannot substitute for
// this invariant: two distinct type parameters that share a name render
// identically, so no deterministic key separates them (their only distinguisher
// is object identity, which is exactly what this order avoids). Stability cannot
// be the safeguard either: several of these slices are collected from
// nondeterministically-ordered maps, where a stable sort would only preserve a
// nondeterministic input order. Callers still use a stable sort, but for
// consistency with plumb's other ordering sites rather than for determinism here.
func CmpType(a, b types.Type) int {
	a = types.Unalias(a)
	b = types.Unalias(b)
	if ra, rb := typeRank(a), typeRank(b); ra != rb {
		return cmp.Compare(ra, rb)
	}
	switch a := a.(type) {
	case *types.Basic:
		return cmp.Compare(a.Kind(), b.(*types.Basic).Kind())
	case *types.Pointer:
		return CmpType(a.Elem(), b.(*types.Pointer).Elem())
	case *types.Slice:
		return CmpType(a.Elem(), b.(*types.Slice).Elem())
	case *types.Array:
		b := b.(*types.Array)
		if c := cmp.Compare(a.Len(), b.Len()); c != 0 {
			return c
		}
		return CmpType(a.Elem(), b.Elem())
	case *types.Chan:
		b := b.(*types.Chan)
		if c := cmp.Compare(a.Dir(), b.Dir()); c != 0 {
			return c
		}
		return CmpType(a.Elem(), b.Elem())
	case *types.Map:
		b := b.(*types.Map)
		if c := CmpType(a.Key(), b.Key()); c != 0 {
			return c
		}
		return CmpType(a.Elem(), b.Elem())
	case *types.Signature:
		b := b.(*types.Signature)
		if c := cmpBool(a.Variadic(), b.Variadic()); c != 0 {
			return c
		}
		if c := cmpTuple(a.Params(), b.Params()); c != 0 {
			return c
		}
		return cmpTuple(a.Results(), b.Results())
	case *types.Tuple:
		return cmpTuple(a, b.(*types.Tuple))
	case *types.Struct:
		b := b.(*types.Struct)
		if c := cmp.Compare(a.NumFields(), b.NumFields()); c != 0 {
			return c
		}
		for i := range a.NumFields() {
			fa, fb := a.Field(i), b.Field(i)
			// Id() not Name() so an unexported field carries its package, matching
			// the identity semantics unify.go compares fields by.
			if c := strings.Compare(fa.Id(), fb.Id()); c != 0 {
				return c
			}
			if c := cmpBool(fa.Embedded(), fb.Embedded()); c != 0 {
				return c
			}
			if c := strings.Compare(a.Tag(i), b.Tag(i)); c != 0 {
				return c
			}
			if c := CmpType(fa.Type(), fb.Type()); c != 0 {
				return c
			}
		}
		return 0
	case *types.Interface:
		b := b.(*types.Interface)
		if c := cmp.Compare(a.NumMethods(), b.NumMethods()); c != 0 {
			return c
		}
		for i := range a.NumMethods() {
			ma, mb := a.Method(i), b.Method(i)
			if c := strings.Compare(ma.Id(), mb.Id()); c != 0 {
				return c
			}
			if c := CmpType(ma.Type(), mb.Type()); c != 0 {
				return c
			}
		}
		if c := cmp.Compare(a.NumEmbeddeds(), b.NumEmbeddeds()); c != 0 {
			return c
		}
		for i := range a.NumEmbeddeds() {
			if c := CmpType(a.EmbeddedType(i), b.EmbeddedType(i)); c != 0 {
				return c
			}
		}
		return 0
	case *types.Named:
		b := b.(*types.Named)
		if c := cmpTypeName(a.Obj(), b.Obj()); c != 0 {
			return c
		}
		return cmpTypeList(a.TypeArgs(), b.TypeArgs())
	case *types.TypeParam:
		b := b.(*types.TypeParam)
		// Distinct parameters with the same name and index compare equal; a stable
		// sort leaves them in input order.
		if c := strings.Compare(a.Obj().Name(), b.Obj().Name()); c != 0 {
			return c
		}
		return cmp.Compare(a.Index(), b.Index())
	case *types.Union:
		b := b.(*types.Union)
		if c := cmp.Compare(a.Len(), b.Len()); c != 0 {
			return c
		}
		for i := range a.Len() {
			ta, tb := a.Term(i), b.Term(i)
			if c := cmpBool(ta.Tilde(), tb.Tilde()); c != 0 {
				return c
			}
			if c := CmpType(ta.Type(), tb.Type()); c != 0 {
				return c
			}
		}
		return 0
	}
	// Unreachable: typeRank already rejected any kind absent from this switch.
	panic(fmt.Sprintf("plumb: CmpType on unsupported type %T", a))
}

// typeRank orders the type kinds, the primary key of CmpType.
func typeRank(t types.Type) int {
	switch t.(type) {
	case *types.Basic:
		return 0
	case *types.Pointer:
		return 1
	case *types.Slice:
		return 2
	case *types.Array:
		return 3
	case *types.Chan:
		return 4
	case *types.Map:
		return 5
	case *types.Signature:
		return 6
	case *types.Tuple:
		return 7
	case *types.Struct:
		return 8
	case *types.Interface:
		return 9
	case *types.Named:
		return 10
	case *types.TypeParam:
		return 11
	case *types.Union:
		return 12
	}
	panic(fmt.Sprintf("plumb: typeRank on unsupported type %T", t))
}

func cmpBool(a, b bool) int {
	switch {
	case a == b:
		return 0
	case !a:
		return -1
	default:
		return 1
	}
}

func cmpTuple(a, b *types.Tuple) int {
	if c := cmp.Compare(a.Len(), b.Len()); c != 0 {
		return c
	}
	for i := range a.Len() {
		if c := CmpType(a.At(i).Type(), b.At(i).Type()); c != 0 {
			return c
		}
	}
	return 0
}

// cmpTypeName orders named types' declaring objects by package path then name:
// environment-independent, unlike pointer identity.
func cmpTypeName(a, b *types.TypeName) int {
	pa, pb := "", ""
	if a.Pkg() != nil {
		pa = a.Pkg().Path()
	}
	if b.Pkg() != nil {
		pb = b.Pkg().Path()
	}
	if c := strings.Compare(pa, pb); c != 0 {
		return c
	}
	return strings.Compare(a.Name(), b.Name())
}

func cmpTypeList(a, b *types.TypeList) int {
	// TypeList.Len is safe on a nil receiver (returns 0).
	la, lb := a.Len(), b.Len()
	if c := cmp.Compare(la, lb); c != 0 {
		return c
	}
	for i := range la {
		if c := CmpType(a.At(i), b.At(i)); c != 0 {
			return c
		}
	}
	return 0
}
