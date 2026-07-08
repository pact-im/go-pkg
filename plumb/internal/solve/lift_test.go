package solve

import (
	"go/token"
	"go/types"
	"testing"

	"go.pact.im/x/plumb/internal/discover"
)

// newParam makes a standalone type parameter with an empty-interface (any)
// constraint, enough for the lift bookkeeping under test.
func newParam(name string) *types.TypeParam {
	tp := types.NewTypeParam(types.NewTypeName(token.NoPos, nil, name, nil), nil)
	tp.SetConstraint(types.NewInterfaceType(nil, nil).Complete())
	return tp
}

func TestCompatibleBind(t *testing.T) {
	tT, tU := newParam("T"), newParam("U")
	intT, strT := types.Typ[types.Int], types.Typ[types.String]
	tests := []struct {
		name string
		a, b map[*types.TypeParam]types.Type
		want bool
	}{
		{"disjoint keys", map[*types.TypeParam]types.Type{tT: intT}, map[*types.TypeParam]types.Type{tU: strT}, true},
		{"agree on shared", map[*types.TypeParam]types.Type{tT: intT}, map[*types.TypeParam]types.Type{tT: intT, tU: strT}, true},
		{"conflict on shared", map[*types.TypeParam]types.Type{tT: intT}, map[*types.TypeParam]types.Type{tT: strT}, false},
		{"empty b", map[*types.TypeParam]types.Type{tT: intT}, map[*types.TypeParam]types.Type{}, true},
		{"empty a", map[*types.TypeParam]types.Type{}, map[*types.TypeParam]types.Type{tT: intT}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := compatibleBind(tc.a, tc.b); got != tc.want {
				t.Errorf("compatibleBind = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestRollbackLiftsReleasesNames pins lift.go’s invariant: rolling back a
// speculative lift truncates liftedMeta AND releases the names it reserved in
// liftedNames, so a later lift can reuse a name a rolled-back attempt had taken.
// The two slices must stay in sync or a near-miss would leak a phantom parameter.
func TestRollbackLiftsReleasesNames(t *testing.T) {
	s := &solver{
		liftedNames: map[string]bool{},
		avoid:       map[string]bool{},
	}
	prov := &discover.Provider{Pos: token.Position{Filename: "x.go", Line: 1}}

	first := s.liftOne(newParam("T"), prov)
	if first.Obj().Name() != "T" {
		t.Fatalf("first lift = %q, want T", first.Obj().Name())
	}

	// A second lift of another T while the first still holds "T" must disambiguate.
	mark := len(s.liftedMeta)
	second := s.liftOne(newParam("T"), prov)
	if second.Obj().Name() != "T2" {
		t.Fatalf("second lift = %q, want T2 (T is taken)", second.Obj().Name())
	}

	// Roll back the second lift: it must shrink liftedMeta and release "T2", while
	// leaving the retained "T".
	s.rollbackLifts(mark)
	if len(s.liftedMeta) != mark {
		t.Fatalf("liftedMeta len = %d after rollback, want %d", len(s.liftedMeta), mark)
	}
	if s.liftedNames["T2"] {
		t.Error("rollback did not release the reserved name T2")
	}
	if !s.liftedNames["T"] {
		t.Error("rollback wrongly released the retained name T")
	}

	// A fresh lift now reuses the released T2 (the next free name after T).
	third := s.liftOne(newParam("T"), prov)
	if third.Obj().Name() != "T2" {
		t.Errorf("re-lift = %q, want T2 (reused after rollback)", third.Obj().Name())
	}
}
