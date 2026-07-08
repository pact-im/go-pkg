package solve

import (
	"go/types"
	"testing"

	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
	"go.pact.im/x/plumb/internal/packagestest"
)

// TestNoCmpTypeTiesInSortedSlices enforces the invariant the determinism of the
// solver rests on, at the point the doc comment on gotypes.CmpType names it:
// CmpType is a WEAK order (two distinct type parameters sharing name and index
// compare equal), and several solve sort sites (pendingDemands, shadowedDemand,
// input ordering, and sortTypesByProducer) sort types drawn from
// nondeterministically-ordered maps with a stable sort, so byte-identical output
// holds only because no two elements ever tie. That safety is emergent (it rests
// on lifted-parameter names being globally unique), not enforced at the sort
// site; a future change that let two same-name/same-index parameters co-occur in
// one of those slices would reintroduce nondeterminism that the shuffle test can
// miss (it fails only when the sampled shuffles happen to reorder the tied pair
// AND the two orders diverge in committed output).
//
// Every element sorted at those four sites is drawn from the input or value-output
// types of the plan’s instances (demand is the union of instance inputs; the
// injector inputs and value outputs are subsets), so if the whole population of
// instance types carries no CmpType-tie between non-Identical types, no sort site
// can. This checks that population over lift-heavy fixtures, catching a tie by its
// existence rather than by a divergence it happens to cause.
func TestNoCmpTypeTiesInSortedSlices(t *testing.T) {
	// Fixtures that exercise the type-parameter sorting paths: free/lifted
	// parameters carried onto the header, joint pinning, pinning revision across
	// rounds, and a constraint near-miss that leaves a generic injector input.
	fixtures := map[string]string{
		"free lifted param": `type Pool[T any] struct{}
//plumb:build
func NewPool[T any]() *Pool[T] { return &Pool[T]{} }
//plumb:build
func PoolSize[T any](p *Pool[T]) int { return 0 }`,

		"joint late-pin merge": `type Key[T any] struct{}
type Val[U any] struct{}
type Box[W any] struct{}
type Cog[X any] struct{}
type SA struct{}
type SB struct{}
//plumb:build
func Make[T, U any]() (Key[T], Val[U]) { return Key[T]{}, Val[U]{} }
//plumb:build
func UseKey(k Key[int]) *SA { return &SA{} }
//plumb:build
func Q2[W, X any](v Val[W]) (Box[W], Cog[X]) { return Box[W]{}, Cog[X]{} }
//plumb:build
func UsePair(b Box[string]) *SB { return &SB{} }`,

		// Two lift-arounds of one template parameter: Make lifts T at both the
		// string and byte instantiations, so the two Key[T] outputs carry two
		// distinct lifted parameters that share the base name "T" and stay
		// distinct only through suffixing (T, T2). Drop the suffixing and the two
		// outputs tie: the exact regression this test guards.
		"double lift-around of one param": `type Key[T any] struct{}
type Val[U any] struct{}
type S struct{}
//plumb:build
func Make[T, U any]() (Key[T], Val[U]) { return Key[T]{}, Val[U]{} }
//plumb:build
func KeyInt() Key[int] { return Key[int]{} }
//plumb:build
func KeyBool() Key[bool] { return Key[bool]{} }
//plumb:build
func NeedS(a Val[string], b Val[byte]) *S { return nil }`,

		"near-miss lift + free input param": `type Stringer interface{ String() string }
type Name string
func (n Name) String() string { return string(n) }
type Box[T any] struct{ V T }
//plumb:build
func Conv[T Stringer, U any](u U) Box[T] { return Box[T]{} }
//plumb:build
func UseName(b Box[Name]) int { return 0 }
//plumb:build
func UseInt(b Box[int]) string { return "" }`,
	}

	for name, src := range fixtures {
		t.Run(name, func(t *testing.T) {
			for _, pl := range solvePlans(t, src) {
				var pop []types.Type
				for _, in := range pl.Order {
					pop = append(pop, in.InputTypes()...)
					pop = append(pop, in.valueOuts()...)
				}
				if a, b, ok := cmpTypeTie(pop); ok {
					t.Errorf("CmpType-tie between non-Identical instance types %s and %s: a stable sort of these is load-order-dependent",
						gotypes.TypeName(a), gotypes.TypeName(b))
				}
			}
		})
	}
}

// cmpTypeTie returns the first pair of non-Identical types in ts that CmpType
// orders equal, and whether one exists. The slices are small, so the quadratic
// scan is fine.
func cmpTypeTie(ts []types.Type) (types.Type, types.Type, bool) {
	for i := range ts {
		for j := i + 1; j < len(ts); j++ {
			if gotypes.CmpType(ts[i], ts[j]) == 0 && !types.Identical(ts[i], ts[j]) {
				return ts[i], ts[j], true
			}
		}
	}
	return nil, nil, false
}

// solvePlans loads src as package p and resolves every set into a plan, treating
// a fresh separate package as the destination so no collision bookkeeping is
// needed. It fails the test on any load or solve error.
func solvePlans(t *testing.T, src string) []*Plan {
	t.Helper()
	loaded, err := packagestest.Load(map[string]string{
		"example.com/p/p.go": "package p\n" + src,
	})
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.TypeErrors) != 0 {
		t.Fatalf("fixture has type errors: %v", loaded.TypeErrors)
	}
	providers, derr := discover.Analyze(loaded.Packages, "example.com/gen", "")
	if derr != nil {
		t.Fatalf("discover: %v", derr)
	}
	bySet := map[string][]*discover.Provider{}
	order := []string{}
	for _, p := range providers {
		if _, seen := bySet[p.SetName]; !seen {
			order = append(order, p.SetName)
		}
		bySet[p.SetName] = append(bySet[p.SetName], p)
	}
	dest := &DestInfo{PkgName: "gen"} // separate package: not scanned, no names
	ctxt := types.NewContext()
	var plans []*Plan
	for _, name := range order {
		pl, perr := Set(name, bySet[name], "example.com/gen", "", dest, ctxt)
		if perr != nil {
			t.Fatalf("solve set %q: %v", name, perr)
		}
		plans = append(plans, pl)
	}
	return plans
}
