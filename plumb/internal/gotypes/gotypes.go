// Package gotypes is plumb’s go/types toolkit: the type-graph walks, structural
// matching, substitution, and predicates the core needs to reason about loaded
// types. Everything here is a pure function of its go/types inputs (no plumb
// concepts, no diagnostics), so it can sit below every phase and be tested on its
// own. It is to go/types what the gopackages package is to go/packages.
package gotypes

import (
	"fmt"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/types/typeutil"
)

// Map and Set are keyed by type identity (types.Identical), not the pointer or ==
// identity a plain Go map uses, because go/types does not canonicalize types: two
// distinct *types.Named values can denote the same type. Both wrap typeutil.Map,
// giving typed keys and values so callers avoid the any-casts a bare typeutil.Map
// forces. Once Go ships a types.Hash (go.dev/issues/69559) both can become plain
// generic maps; see typeutil.Map’s own TODO.

// Map is a type-keyed map with typed values. The zero value is an empty map ready
// to use.
type Map[K types.Type, V any] struct{ m typeutil.Map }

// At returns the value stored for key and whether it was present.
func (m *Map[K, V]) At(key K) (V, bool) {
	if v := m.m.At(key); v != nil {
		return v.(V), true
	}
	var zero V
	return zero, false
}

// Set stores value for key, replacing any previous value.
func (m *Map[K, V]) Set(key K, value V) { m.m.Set(key, value) }

// Len returns the number of entries.
func (m *Map[K, V]) Len() int { return m.m.Len() }

// All iterates the entries in unspecified order; it satisfies iter.Seq2[K, V].
func (m *Map[K, V]) All(yield func(K, V) bool) {
	stop := false
	m.m.Iterate(func(k types.Type, v any) {
		if !stop && !yield(k.(K), v.(V)) {
			stop = true
		}
	})
}

// Set is a type-keyed set. The zero value is an empty set ready to use; it is the
// visited-set for cyclic-type-safe graph walks and the presence-set the solver
// keeps of demanded, consumed, and already-instantiated types.
type Set[T types.Type] struct{ m typeutil.Map }

// Add records t and reports whether it was newly added: false if t was already
// present. A cyclic-type walk guards recursion with `if !seen.Add(t) { return }`,
// which stops on a revisit.
func (s *Set[T]) Add(t T) bool {
	if s.m.At(t) != nil {
		return false
	}
	s.m.Set(t, true)
	return true
}

// Contains reports whether t is in the set.
func (s *Set[T]) Contains(t T) bool {
	return s.m.At(t) != nil
}

// Elements iterates the members in unspecified order; it satisfies iter.Seq[T].
// Callers that need a deterministic order sort what they collect.
func (s *Set[T]) Elements(yield func(T) bool) {
	stop := false
	s.m.Iterate(func(k types.Type, _ any) {
		if !stop && !yield(k.(T)) {
			stop = true
		}
	})
}

// ListKey packs a type list into a tuple usable as a typeutil.Map key:
// types.Identical compares tuples element-wise (ignoring the synthetic names), so
// two identical lists collapse to one entry. It backs the (provider, type-args)
// instantiation keys, where the provider supplies the outer map dimension.
func ListKey(ts []types.Type) *types.Tuple {
	vars := make([]*types.Var, len(ts))
	for i, t := range ts {
		vars[i] = types.NewVar(token.NoPos, nil, "", t)
	}
	return types.NewTuple(vars...)
}

// pointerElem returns the element of a pointer to a bridgeable type and true, when
// t is *E for some bridgeable E (a named type, a type parameter, or a predeclared
// basic; see bridgeable). Used for the value/pointer bridge.
func pointerElem(t types.Type) (types.Type, bool) {
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return nil, false
	}
	if bridgeable(ptr.Elem()) {
		return ptr.Elem(), true
	}
	return nil, false
}

// bridgeable reports whether t is a type whose value/pointer dual is meaningful
// for the bridge: a named type, a type parameter, or a predeclared basic type
// (so int⇄*int bridges just like MyInt⇄*MyInt). A composite type written inline
// (a slice, map, struct, func, etc.) has no dual; only a pointer to a bridgeable
// type does.
func bridgeable(t types.Type) bool {
	switch types.Unalias(t).(type) {
	case *types.Named, *types.TypeParam, *types.Basic:
		return true
	}
	return false
}

// DualType returns the value/pointer dual of a bridgeable type. The argument
// may be either form: for a pointer to a bridgeable type it returns the element
// (*E → E), for a bridgeable non-pointer the pointer type (T → *T), and for
// everything else (nil, false).
func DualType(t types.Type) (types.Type, bool) {
	if elem, ok := pointerElem(t); ok {
		return elem, true
	}
	if bridgeable(t) {
		return types.NewPointer(t), true
	}
	return nil, false
}

// --- type predicates ---------------------------------------------------------

var errorType = types.Universe.Lookup("error").Type()

// comparableType is the predeclared comparable constraint, which must not be
// collapsed to "any".
var comparableType = types.Universe.Lookup("comparable").Type()

// ConstraintCollapsesToAny reports whether a type-parameter constraint renders as
// the bare keyword "any": an interface with no methods and no embeddeds that is
// not the predeclared comparable. emit renders such a constraint as "any" and
// names neither it nor its package, so solve’s reachability check must agree.
// This is the single predicate both consult, so they cannot drift (a collapsed
// unexported constraint must not be rejected, and a rendered one must be imported).
func ConstraintCollapsesToAny(c types.Type) bool {
	iface, ok := c.Underlying().(*types.Interface)
	return ok && iface.NumMethods() == 0 && iface.NumEmbeddeds() == 0 &&
		!types.Identical(c, comparableType)
}

// IsErrorType reports whether t is exactly the predeclared error interface.
func IsErrorType(t types.Type) bool {
	return types.Identical(t, errorType)
}

// IsBareCleanup reports whether t is the bare type func() (no params, no
// results). A *named* function type is not bare; a type *alias* to func() is,
// since an alias and its target are the same type.
func IsBareCleanup(t types.Type) bool {
	// No variadic check: a variadic signature always has at least one parameter,
	// so the zero-parameter test already excludes it.
	sig, ok := types.Unalias(t).(*types.Signature)
	return ok && sig.Params().Len() == 0 && sig.Results().Len() == 0
}

// IsFailableCleanup reports whether t is the bare type func() error.
func IsFailableCleanup(t types.Type) bool {
	sig, ok := types.Unalias(t).(*types.Signature)
	if !ok || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}
	return IsErrorType(sig.Results().At(0).Type())
}

// IsUntyped reports whether t is an untyped basic type (e.g. an untyped
// constant’s default-less type).
func IsUntyped(t types.Type) bool {
	b, ok := t.(*types.Basic)
	return ok && b.Info()&types.IsUntyped != 0
}

// IsUntypedNil reports whether t is the predeclared untyped nil.
func IsUntypedNil(t types.Type) bool {
	b, ok := t.(*types.Basic)
	return ok && b.Kind() == types.UntypedNil
}

// TypeName renders t for diagnostics, qualifying foreign packages by name.
func TypeName(t types.Type) string {
	return types.TypeString(t, func(p *types.Package) string {
		if p == nil {
			return ""
		}
		return p.Name()
	})
}

// KindOfType names the broad shape of t for a "not a struct" diagnostic:
// interfaces are called out by name, and everything else is "a non-struct type".
func KindOfType(t types.Type) string {
	if _, ok := t.(*types.Interface); ok {
		return "an interface"
	}
	return "a non-struct type"
}

// GenericOrigin returns the generic origin of t. t is a provider’s declared or
// receiver type, which is always a *types.Named or a *types.Alias; any other
// shape is a broken invariant and panics rather than passing t through, which
// would mask the bug downstream.
func GenericOrigin(t types.Type) types.Type {
	switch u := t.(type) {
	case *types.Named:
		return u.Origin()
	case *types.Alias:
		return u.Origin()
	}
	panic(fmt.Sprintf("plumb: GenericOrigin on %T; want a defined type or alias", t))
}

// TypeParamsOf returns the type parameters of t. As with GenericOrigin, t is a
// provider’s declared or receiver type and is always a *types.Named or a
// *types.Alias; any other shape panics.
func TypeParamsOf(t types.Type) *types.TypeParamList {
	switch u := t.(type) {
	case *types.Named:
		return u.TypeParams()
	case *types.Alias:
		return u.TypeParams()
	}
	panic(fmt.Sprintf("plumb: TypeParamsOf on %T; want a defined type or alias", t))
}

// TypeNameOf returns the declaring object of t. t is a provider’s declared or
// receiver type, always a *types.Named or a *types.Alias; any other shape
// panics. Returning a nil *types.TypeName instead would wrap a nil pointer in a
// non-nil types.Object interface (a typed-nil that slips past a caller’s
// obj != nil guard and nil-panics on first use), so the invariant is enforced
// here, where every caller then gets a non-nil result.
func TypeNameOf(t types.Type) *types.TypeName {
	switch u := t.(type) {
	case *types.Named:
		return u.Obj()
	case *types.Alias:
		return u.Obj()
	}
	panic(fmt.Sprintf("plumb: TypeNameOf on %T; want a defined type or alias", t))
}
