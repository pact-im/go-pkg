package solve

import (
	"go/types"
	"testing"
)

// TestErrRestartFormatSafe pins that the restart sentinel is safe to format.
// errRestart rides the *diag.Error return channel but is a control-flow signal
// compared by identity; a stray fmt/log of an in-flight error must print its
// message rather than nil-panic (which a zero-value *diag.Error would).
func TestErrRestartFormatSafe(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("errRestart.Error() panicked: %v", r)
		}
	}()
	if errRestart.Error() == "" {
		t.Error("errRestart formats as the empty string; want a message")
	}
}

// TestPinningRevisionCommitsUnion pins revision.go's contract directly on the
// resolved plan: when a joint template's pins arrive across separate fixpoint
// rounds, the solver commits their union and replays, so the template is
// instantiated ONCE at the merged pinning that serves every consumer, never
// split into half-pinned calls and never leaking a phantom lifted parameter. The
// corpus (joint_late_pin_merges, joint_late_pin_staircase) checks the generated
// source; this checks the mechanism's output (the committed union on the
// instance) without running emit, and covers both a single revision and the
// commitment-upgrade staircase whose regression would run the restart cap.
func TestPinningRevisionCommitsUnion(t *testing.T) {
	intT, strT := types.Typ[types.Int], types.Typ[types.String]
	mapT := types.NewMap(intT, strT)

	tests := []struct {
		name      string
		prov      string
		wantTargs []types.Type
		src       string
	}{
		{
			// One revision: Make[T, U] produces Key[T] and Val[U]. UseKey pins T=int
			// immediately; Val[string] is demanded only a round later (Q2 → UsePair),
			// so U=string is the late pin merged into the existing Make instance.
			name:      "single late pin merges",
			prov:      "Make",
			wantTargs: []types.Type{intT, strT},
			src: `type Key[T any] struct{}
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
		},
		{
			// Two sequential revisions upgrading one commitment in place: [int] grows
			// to [int, string], then to [int, string, map[int]string], the third slot
			// (~map[T]U) viable only once the first two pins are in force. Appending
			// instead of upgrading would re-derive the extension forever and hit the
			// restart cap, so a converged [int, string, map[int]string] here also pins
			// that guard.
			name:      "staircase upgrades commitment in place",
			prov:      "NewTriple",
			wantTargs: []types.Type{intT, strT, mapT},
			src: `type R1[T any] struct{}
type R2[U any] struct{}
type R3[V any] struct{}
type B1[W any] struct{}
type C1[X any] struct{}
type B2[W any] struct{}
type C2[X any] struct{}
type SA struct{}
type SB struct{}
type SC struct{}
//plumb:build
func NewTriple[T comparable, U any, V interface{ ~map[T]U }]() (R1[T], R2[U], R3[V]) {
	return R1[T]{}, R2[U]{}, R3[V]{}
}
//plumb:build
func UseR1(r R1[int]) *SA { return &SA{} }
//plumb:build
func QU[W, X any](r R2[W]) (B1[W], C1[X]) { return B1[W]{}, C1[X]{} }
//plumb:build
func UseB1(b B1[string]) *SB { return &SB{} }
//plumb:build
func QV[W, X any](r R3[W]) (B2[W], C2[X]) { return B2[W]{}, C2[X]{} }
//plumb:build
func UseB2(b B2[map[int]string]) *SC { return &SC{} }`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var insts []*Instance
			for _, pl := range solvePlans(t, tc.src) {
				for _, in := range pl.Order {
					if in.Prov.Name == tc.prov {
						insts = append(insts, in)
					}
				}
			}
			if len(insts) != 1 {
				t.Fatalf("%s instantiated %d times; want 1 (the committed union, not a split)", tc.prov, len(insts))
			}
			if got := insts[0].Targs; !identicalTypes(got, tc.wantTargs) {
				t.Fatalf("%s pinned at %v; want %v (the committed union)", tc.prov, got, tc.wantTargs)
			}
		})
	}
}

func identicalTypes(a, b []types.Type) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !types.Identical(a[i], b[i]) {
			return false
		}
	}
	return true
}
