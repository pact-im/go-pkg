package gotypes_test

import (
	"go/types"
	"slices"
	"testing"

	"go.pact.im/x/plumb/internal/gotypes"
	"go.pact.im/x/plumb/internal/packagestest"
)

// fixture is the shared source the predicate tests look types up from. It is
// loaded once via the in-memory loader; each test pulls the named types and
// function signatures it needs from package p's scope.
const fixture = `
type T struct{}
type Box[A any] struct{ V A }
type Pair[A, B any] struct{ X A; Y B }

// shapes exposes a spread of type forms as parameter types to read back.
func shapes(
	i int,
	pt *T,
	sl []int,
	bi Box[int],
	bb Box[Box[int]],
	mp map[int]*T,
	pp **T,
) {}

// genBox is generic so its signature carries a type parameter A and the pattern
// type Box[A]; concBox fixes A to int.
func genBox[A any](x Box[A]) {}
func concBox(x Box[int])     {}

// genPair / concPair give a two-parameter pattern and its concretion.
func genPair[A, B any](x *Pair[A, B]) {}
func concPair(x *Pair[int, string])   {}

// repeat pins the same parameter twice, for the conflicting-bind case.
func repeat[A any](x Pair[A, A]) {}
`

func load(t *testing.T) *types.Package {
	t.Helper()
	loaded, err := packagestest.Load(map[string]string{
		"example.com/p/p.go": "package p\n" + fixture,
	})
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.TypeErrors) != 0 {
		t.Fatalf("fixture has type errors: %v", loaded.TypeErrors)
	}
	for _, p := range loaded.Packages {
		if p.PkgPath == "example.com/p" {
			return p.Types
		}
	}
	t.Fatal("package p not loaded")
	return nil
}

func named(t *testing.T, pkg *types.Package, name string) types.Type {
	t.Helper()
	obj := pkg.Scope().Lookup(name)
	if obj == nil {
		t.Fatalf("type %q not found", name)
	}
	return obj.Type()
}

func sig(t *testing.T, pkg *types.Package, fn string) *types.Signature {
	t.Helper()
	obj := pkg.Scope().Lookup(fn)
	f, ok := obj.(*types.Func)
	if !ok {
		t.Fatalf("%q is not a func", fn)
	}
	return f.Type().(*types.Signature)
}

// param returns the i-th parameter type of function fn.
func param(t *testing.T, pkg *types.Package, fn string, i int) types.Type {
	t.Helper()
	return sig(t, pkg, fn).Params().At(i).Type()
}

func TestDualType(t *testing.T) {
	pkg := load(t)
	T := named(t, pkg, "T")
	ptrT := types.NewPointer(T)

	t.Run("value to pointer", func(t *testing.T) {
		got, ok := gotypes.DualType(T)
		if !ok || !types.Identical(got, ptrT) {
			t.Fatalf("DualType(T) = %v, %v; want *T, true", got, ok)
		}
	})
	t.Run("pointer to value", func(t *testing.T) {
		got, ok := gotypes.DualType(ptrT)
		if !ok || !types.Identical(got, T) {
			t.Fatalf("DualType(*T) = %v, %v; want T, true", got, ok)
		}
	})
	t.Run("basic type bridges", func(t *testing.T) {
		got, ok := gotypes.DualType(types.Typ[types.Int])
		if !ok || !types.Identical(got, types.NewPointer(types.Typ[types.Int])) {
			t.Fatalf("DualType(int) = %v, %v; want *int, true", got, ok)
		}
		got, ok = gotypes.DualType(types.NewPointer(types.Typ[types.Int]))
		if !ok || !types.Identical(got, types.Typ[types.Int]) {
			t.Fatalf("DualType(*int) = %v, %v; want int, true", got, ok)
		}
	})
	t.Run("composite has no dual", func(t *testing.T) {
		// An inline composite ([]int) is not a name, so it has no value/pointer dual.
		if _, ok := gotypes.DualType(types.NewSlice(types.Typ[types.Int])); ok {
			t.Fatal("DualType([]int) reported a dual; want none")
		}
	})
	t.Run("double pointer has no dual", func(t *testing.T) {
		// The bridge is one level: **T is neither bridgeable nor a pointer to a
		// bridgeable type (its element *T is itself a pointer).
		if _, ok := gotypes.DualType(types.NewPointer(ptrT)); ok {
			t.Fatal("DualType(**T) reported a dual; want none (bridge is one level)")
		}
	})
}

func TestTypeDepth(t *testing.T) {
	pkg := load(t)
	cases := []struct {
		param int
		want  int
		desc  string
	}{
		{0, 1, "int"},
		{1, 2, "*T"},
		{2, 2, "[]int"},
		{3, 2, "Box[int]"},
		{4, 3, "Box[Box[int]]"},
		{5, 3, "map[int]*T"}, // 1 + max(depth int=1, depth *T=2)
		{6, 3, "**T"},
	}

	for _, c := range cases {
		got := gotypes.TypeDepth(param(t, pkg, "shapes", c.param))
		if got != c.want {
			t.Errorf("TypeDepth(%s) = %d; want %d", c.desc, got, c.want)
		}
	}
}

func TestCmpType(t *testing.T) {
	pkg := load(t)
	ts := []types.Type{
		types.Typ[types.Int],
		types.Typ[types.String],
		named(t, pkg, "T"),
		types.NewPointer(named(t, pkg, "T")),
		param(t, pkg, "concBox", 0),  // Box[int]
		param(t, pkg, "concPair", 0), // *Pair[int, string]
		named(t, pkg, "Box"),
	}

	t.Run("reflexive zero", func(t *testing.T) {
		for _, a := range ts {
			if c := gotypes.CmpType(a, a); c != 0 {
				t.Errorf("CmpType(a, a) = %d; want 0 for %v", c, a)
			}
		}
	})
	t.Run("antisymmetric", func(t *testing.T) {
		for _, a := range ts {
			for _, b := range ts {
				if c, r := gotypes.CmpType(a, b), gotypes.CmpType(b, a); c != -r {
					t.Errorf("CmpType not antisymmetric for %v,%v: %d vs %d", a, b, c, r)
				}
			}
		}
	})
	t.Run("no ties among these representative types", func(t *testing.T) {
		// CmpType is deliberately a weak order (see cmp.go): distinct type
		// parameters sharing a name and index compare equal, so a zero result does
		// not in general imply identity. These fixtures happen to contain no such
		// tying pair, so a zero here does mean identical, but that is a property of
		// this sample, not a totality guarantee. The tie-freedom the solver's
		// determinism actually rests on is pinned by
		// solve.TestNoCmpTypeTiesInSortedSlices.
		for i, a := range ts {
			for j, b := range ts {
				if i == j {
					continue
				}
				if gotypes.CmpType(a, b) == 0 && !types.Identical(a, b) {
					t.Errorf("CmpType(%v,%v)=0 but not identical", a, b)
				}
			}
		}
	})
	t.Run("sort stable under input order", func(t *testing.T) {
		a := slices.SortedFunc(slices.Values(ts), gotypes.CmpType)
		rev := slices.Clone(ts)
		slices.Reverse(rev)
		b := slices.SortedFunc(slices.Values(rev), gotypes.CmpType)
		if !slices.EqualFunc(a, b, types.Identical) {
			t.Error("sorted order depends on input order")
		}
	})
}

func TestUnify(t *testing.T) {
	pkg := load(t)
	// genBox's signature gives the pattern Box[A] and the parameter A.
	genBox := sig(t, pkg, "genBox")
	A := genBox.TypeParams().At(0)
	patBox := genBox.Params().At(0).Type() // Box[A]
	boxInt := param(t, pkg, "concBox", 0)  // Box[int]

	t.Run("binds single parameter", func(t *testing.T) {
		bind := map[*types.TypeParam]types.Type{}
		if !gotypes.Unify(patBox, boxInt, map[*types.TypeParam]bool{A: true}, bind) {
			t.Fatal("Box[A] did not unify with Box[int]")
		}
		if got := bind[A]; got == nil || !types.Identical(got, types.Typ[types.Int]) {
			t.Fatalf("bound A = %v; want int", got)
		}
	})
	t.Run("no match against different shape", func(t *testing.T) {
		bind := map[*types.TypeParam]types.Type{}
		sl := param(t, pkg, "shapes", 2) // []int
		if gotypes.Unify(patBox, sl, map[*types.TypeParam]bool{A: true}, bind) {
			t.Fatal("Box[A] wrongly unified with []int")
		}
	})

	genPair := sig(t, pkg, "genPair")
	pA, pB := genPair.TypeParams().At(0), genPair.TypeParams().At(1)
	patPair := genPair.Params().At(0).Type() // *Pair[A, B]
	pairConc := param(t, pkg, "concPair", 0) // *Pair[int, string]

	t.Run("binds multiple parameters", func(t *testing.T) {
		bind := map[*types.TypeParam]types.Type{}
		if !gotypes.Unify(patPair, pairConc, map[*types.TypeParam]bool{pA: true, pB: true}, bind) {
			t.Fatal("*Pair[A,B] did not unify with *Pair[int,string]")
		}
		if !types.Identical(bind[pA], types.Typ[types.Int]) || !types.Identical(bind[pB], types.Typ[types.String]) {
			t.Fatalf("bound A=%v B=%v; want int, string", bind[pA], bind[pB])
		}
	})

	repeat := sig(t, pkg, "repeat")
	rA := repeat.TypeParams().At(0)
	patRepeat := repeat.Params().At(0).Type() // Pair[A, A]

	t.Run("rejects conflicting binding", func(t *testing.T) {
		bind := map[*types.TypeParam]types.Type{}
		// Pair[A, A] cannot unify with Pair[int, string]: A would be both.
		if gotypes.Unify(patRepeat, pairConcValue(t, pkg), map[*types.TypeParam]bool{rA: true}, bind) {
			t.Fatal("Pair[A,A] wrongly unified with Pair[int,string]")
		}
	})
}

// pairConcValue returns the value type Pair[int, string] (concPair's parameter is
// the pointer form, so dereference it).
func pairConcValue(t *testing.T, pkg *types.Package) types.Type {
	t.Helper()
	ptr := param(t, pkg, "concPair", 0).(*types.Pointer)
	return ptr.Elem()
}

func TestWalkNamed(t *testing.T) {
	pkg := load(t)
	// Box[T] nests the named type T as Box's sole type argument, so T is reachable
	// only by descending below Box, the signal that distinguishes “descended” from
	// “pruned”.
	boxT, err := types.Instantiate(types.NewContext(), named(t, pkg, "Box"), []types.Type{named(t, pkg, "T")}, false)
	if err != nil {
		t.Fatalf("instantiate Box[T]: %v", err)
	}

	walk := func(prune string) []string {
		var got []string
		gotypes.WalkNamed(boxT, func(tn *types.TypeName) bool {
			got = append(got, tn.Name())
			return tn.Name() != prune
		})
		return got
	}

	t.Run("descends when visit returns true", func(t *testing.T) {
		if got := walk(""); !slices.Equal(got, []string{"Box", "T"}) {
			t.Fatalf("walk visited %v; want [Box T]", got)
		}
	})
	t.Run("prunes below a node when visit returns false", func(t *testing.T) {
		// Returning false at Box must stop descent into its type argument T.
		if got := walk("Box"); !slices.Equal(got, []string{"Box"}) {
			t.Fatalf("walk visited %v; want [Box] (T pruned)", got)
		}
	})
}

// identical returns two distinct *types.Pointer values that denote the same
// type: a plain Go map keys them apart, a type-identity map collapses them. This
// is the property Map and Set exist to provide, so every keying assertion below
// uses the second value to probe an entry stored under the first.
func identical(t *testing.T) (k1, k2 types.Type) {
	t.Helper()
	T := named(t, load(t), "T")
	k1, k2 = types.NewPointer(T), types.NewPointer(T)
	if k1 == k2 {
		t.Fatal("want distinct pointer values for the identity test")
	}
	return k1, k2
}

func TestMap(t *testing.T) {
	k1, k2 := identical(t)

	var m gotypes.Map[types.Type, int]
	if got, ok := m.At(k1); ok {
		t.Errorf("At on empty map = %d, %v; want 0, false", got, ok)
	}
	m.Set(k1, 1)
	if got, ok := m.At(k1); !ok || got != 1 {
		t.Errorf("At(k1) = %d, %v; want 1, true", got, ok)
	}
	// k2 is a different pointer denoting the same type: it hits k1's entry.
	if got, ok := m.At(k2); !ok || got != 1 {
		t.Errorf("At(k2) = %d, %v; want 1, true (identity keying)", got, ok)
	}
	m.Set(k2, 2) // overwrites the same entry, does not add a second
	if got, _ := m.At(k1); got != 2 {
		t.Errorf("At(k1) after Set(k2, 2) = %d; want 2", got)
	}
	if m.Len() != 1 {
		t.Errorf("Len = %d; want 1 (k1 and k2 are the same key)", m.Len())
	}
}

func TestMapAllBreak(t *testing.T) {
	pkg := load(t)
	var m gotypes.Map[types.Type, int]
	m.Set(named(t, pkg, "T"), 0)
	m.Set(named(t, pkg, "Box"), 1)
	m.Set(named(t, pkg, "Pair"), 2)

	// All emulates early termination over typeutil.Map.Iterate, which cannot stop
	// early itself; breaking after the first entry must visit exactly one.
	n := 0
	for range m.All {
		n++
		break
	}
	if n != 1 {
		t.Errorf("All visited %d entries after break; want 1", n)
	}
	if got := m.Len(); got != 3 {
		t.Errorf("Len = %d; want 3 (break must not drop entries)", got)
	}
}

func TestSet(t *testing.T) {
	k1, k2 := identical(t)

	var s gotypes.Set[types.Type]
	if s.Contains(k1) {
		t.Error("empty set Contains(k1) = true; want false")
	}
	if !s.Add(k1) {
		t.Error("Add(k1) on empty set = false; want true (newly added)")
	}
	if s.Add(k1) {
		t.Error("Add(k1) repeat = true; want false (already present)")
	}
	// k2 denotes the same type, so it is already present.
	if s.Add(k2) {
		t.Error("Add(k2) = true; want false (identity keying)")
	}
	if !s.Contains(k2) {
		t.Error("Contains(k2) = false; want true (identity keying)")
	}
}

func TestSetElements(t *testing.T) {
	pkg := load(t)
	var s gotypes.Set[types.Type]
	for _, n := range []string{"T", "Box", "Pair"} {
		s.Add(named(t, pkg, n))
	}

	// Break after the first: Elements must stop, visiting exactly one.
	n := 0
	for range s.Elements {
		n++
		break
	}
	if n != 1 {
		t.Errorf("Elements visited %d after break; want 1", n)
	}
	// A full pass sees every member.
	full := 0
	for range s.Elements {
		full++
	}
	if full != 3 {
		t.Errorf("Elements visited %d; want 3", full)
	}
}

// TestSetTupleKeys exercises the Set[*types.Tuple] instantiation the solver keys
// its per-provider done-sets by, where ListKey packs a type list into the tuple
// key: two lists that are element-wise identical collapse, a reordered one does
// not.
func TestSetTupleKeys(t *testing.T) {
	pkg := load(t)
	T := named(t, pkg, "T")
	intT := types.Typ[types.Int]
	k1 := gotypes.ListKey([]types.Type{T, intT})
	k2 := gotypes.ListKey([]types.Type{T, intT}) // distinct tuple, identical contents
	k3 := gotypes.ListKey([]types.Type{intT, T}) // reordered: a different key

	var s gotypes.Set[*types.Tuple]
	if !s.Add(k1) {
		t.Error("Add(k1) = false; want true")
	}
	if s.Add(k2) {
		t.Error("Add(k2) = true; want false (element-wise tuple identity)")
	}
	if !s.Add(k3) {
		t.Error("Add(k3 reordered) = false; want true (distinct tuple)")
	}
}

func TestSubst(t *testing.T) {
	pkg := load(t)
	ctxt := types.NewContext()
	genBox := sig(t, pkg, "genBox")
	A := genBox.TypeParams().At(0)
	patBox := genBox.Params().At(0).Type() // Box[A]
	boxInt := param(t, pkg, "concBox", 0)  // Box[int]

	t.Run("substitutes parameter", func(t *testing.T) {
		got := gotypes.Subst(ctxt, patBox, map[*types.TypeParam]types.Type{A: types.Typ[types.Int]})
		if !types.Identical(got, boxInt) {
			t.Fatalf("Subst(Box[A], A->int) = %v; want Box[int]", got)
		}
	})
	t.Run("leaves unrelated type unchanged", func(t *testing.T) {
		T := named(t, pkg, "T")
		got := gotypes.Subst(ctxt, T, map[*types.TypeParam]types.Type{A: types.Typ[types.Int]})
		if !types.Identical(got, T) {
			t.Fatalf("Subst(T, A->int) = %v; want T unchanged", got)
		}
	})
	t.Run("substitutes under a pointer", func(t *testing.T) {
		// *A with A→*T yields **T.
		ptrA := types.NewPointer(A)
		ptrT := types.NewPointer(named(t, pkg, "T"))
		got := gotypes.Subst(ctxt, ptrA, map[*types.TypeParam]types.Type{A: ptrT})
		if !types.Identical(got, types.NewPointer(ptrT)) {
			t.Fatalf("Subst(*A, A->*T) = %v; want **T", got)
		}
	})
}

// TestNamedAliasHelpersPanicOnOtherShapes pins the invariant enforced by
// TypeNameOf, GenericOrigin, and TypeParamsOf: they accept only a defined type
// or an alias (the two forms a provider's declared/receiver type can take) and
// panic on anything else, rather than returning a fallback that would mask a
// broken invariant: for TypeNameOf, a typed-nil that defeats a caller's
// obj != nil guard.
func TestNamedAliasHelpersPanicOnOtherShapes(t *testing.T) {
	// A bare basic type is neither a *types.Named nor a *types.Alias.
	notNamed := types.Typ[types.Int]
	cases := []struct {
		name string
		call func()
	}{
		{"TypeNameOf", func() { gotypes.TypeNameOf(notNamed) }},
		{"GenericOrigin", func() { gotypes.GenericOrigin(notNamed) }},
		{"TypeParamsOf", func() { gotypes.TypeParamsOf(notNamed) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("%s(*types.Basic) did not panic", c.name)
				}
			}()
			c.call()
		})
	}
}

// TestMentionsParamsCoverage pins that the typeContains walk shared by
// MentionsParams and subster.mentions descends into every composite go/types
// kind: a type parameter buried under any one of them must be found, and an
// otherwise-identical shape free of the parameter must not be. The negative twin
// proves each positive came from descending, not from the kind alone, so
// dropping a case in the shared walk flips a row and fails here, guarding against
// a silent coverage narrowing.
func TestMentionsParamsCoverage(t *testing.T) {
	pkg := load(t)
	A := sig(t, pkg, "genBox").TypeParams().At(0) // a real *types.TypeParam
	params := map[*types.TypeParam]bool{A: true}
	tInt := types.Typ[types.Int]

	// method builds interface{ M(x) }, so a positive requires descending the
	// method's signature.
	method := func(x types.Type) types.Type {
		mSig := types.NewSignatureType(nil, nil, nil, types.NewTuple(types.NewVar(0, pkg, "", x)), nil, false)
		iface := types.NewInterfaceType([]*types.Func{types.NewFunc(0, pkg, "M", mSig)}, nil)
		iface.Complete()
		return iface
	}
	// union builds interface{ []x | []int }, exercising the interface-embedded and
	// union-term descents together.
	union := func(x types.Type) types.Type {
		u := types.NewUnion([]*types.Term{
			types.NewTerm(false, types.NewSlice(x)),
			types.NewTerm(false, types.NewSlice(tInt)),
		})
		iface := types.NewInterfaceType(nil, []types.Type{u})
		iface.Complete()
		return iface
	}
	boxOf := func(x types.Type) types.Type {
		inst, err := types.Instantiate(types.NewContext(), named(t, pkg, "Box"), []types.Type{x}, false)
		if err != nil {
			t.Fatalf("instantiate Box[%v]: %v", x, err)
		}
		return inst
	}

	kinds := []struct {
		name string
		of   func(types.Type) types.Type
	}{
		{"pointer", func(x types.Type) types.Type { return types.NewPointer(x) }},
		{"slice", func(x types.Type) types.Type { return types.NewSlice(x) }},
		{"array", func(x types.Type) types.Type { return types.NewArray(x, 2) }},
		{"chan", func(x types.Type) types.Type { return types.NewChan(types.SendRecv, x) }},
		{"map key", func(x types.Type) types.Type { return types.NewMap(x, tInt) }},
		{"map elem", func(x types.Type) types.Type { return types.NewMap(tInt, x) }},
		{"named type-arg", boxOf},
		{"signature param and result", func(x types.Type) types.Type {
			return types.NewSignatureType(nil, nil, nil, types.NewTuple(types.NewVar(0, pkg, "", x)), types.NewTuple(types.NewVar(0, pkg, "", x)), false)
		}},
		{"struct field", func(x types.Type) types.Type {
			return types.NewStruct([]*types.Var{types.NewField(0, pkg, "F", x, false)}, nil)
		}},
		{"interface method", method},
		{"interface union term", union},
	}

	for _, k := range kinds {
		t.Run(k.name, func(t *testing.T) {
			if !gotypes.MentionsParams(k.of(A), params) {
				t.Errorf("MentionsParams did not descend into %s to find the parameter", k.name)
			}
			if gotypes.MentionsParams(k.of(tInt), params) {
				t.Errorf("MentionsParams(%s of int) = true; want false (descended but found no parameter)", k.name)
			}
		})
	}
}
