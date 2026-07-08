package solve

import (
	"go/token"
	"go/types"
	"maps"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
)

// tryJoint instantiates result-generic templates from the demands that are still
// unserved. The pending demands are clustered into one binding per distinct
// pinning, and the template is instantiated once per cluster, lifting the
// parameters no demand pins, so one template may serve several consumers at
// different types (the multi-result analog of pinning a single-result template at
// several types).
//
// This is the general path; the worklist is the fast case ahead of it. The
// worklist first handles every demand that a single result pins completely.
// Whatever it cannot fully pin from one demand (a demand that matches only a
// partial result while another result mentions more parameters, or pinnings that
// arrive jointly across disjoint results) falls to clustering here, where the
// unpinned parameters are lifted. Dispatch is therefore driven by what the
// demands actually pin, not by whether some result structurally covers every
// parameter: a template whose all-covering result is never the one demanded is
// still served, by lifting the parameters its demanded (partial) result leaves free.
//
// Two clusters never produce the same output type, because every distinct pin for
// a parameter lands in exactly one cluster, and clusters seed only from still-
// unsatisfied demands, so a type an existing instance already supplies never
// drives a second instantiation; a genuine duplicate across providers still
// surfaces downstream as an ambiguous producer.
func (s *solver) tryJoint() (bool, *diag.Error) {
	if err := s.checkTemplateAmbiguity(); err != nil {
		return false, err
	}
	for _, p := range s.resGen {
		made := false
		for _, bind := range s.jointClusters(p) {
			if s.jointClusterDone(p, bind) {
				continue
			}
			s.markJointClusterDone(p, bind)
			ok, err := s.instantiateTemplate(p, bind)
			if err != nil {
				return false, err
			}
			if ok {
				made = true
			}
		}
		if made {
			return true, nil
		}
	}
	return false, nil
}

// jointClusters groups the still-unsatisfied demands that pin p’s results into
// consistent bindings: one per distinct pinning. Each demand contributes the
// binding it forces on p’s result parameters; demands whose bindings agree merge
// into one cluster, and conflicting ones form separate clusters, each a distinct
// instantiation. Clusters are built in deterministic demand order, and a
// parameter no demand pins is left for instantiateTemplate to lift. Seeding from
// the unsatisfied demands (not every demand) keeps a type an existing instance
// already supplies from driving a redundant (and possibly conflicting)
// instantiation as later demands arrive across fixpoint passes.
func (s *solver) jointClusters(p *discover.Provider) []map[*types.TypeParam]types.Type {
	params := paramSet(p)
	skel := s.skeleton(p)
	outs := skel.valueOuts()

	var clusters []map[*types.TypeParam]types.Type
	for _, d := range s.pendingDemands() {
		cands, matched := clusterCands(outs, d, params)
		if !matched {
			// None of the results unifies with d’s exact form: try the dual, so a
			// pointer demand can pin a template producing the value form: the
			// instantiation supplies the dual and d rides the bridge.
			if dt, ok := gotypes.DualType(d); ok {
				cands, _ = clusterCands(outs, dt, params)
			}
		}
		if len(cands) == 0 {
			continue
		}

		// Prefer merging one candidate into an existing cluster, filling a slot that
		// cluster leaves open. This is joint pinning: different demands pinning
		// different slots of one instance. Validate that the merged binding actually
		// instantiates, so a near-miss pin falls into its own cluster to lift around
		// rather than poisoning this one (compatibleBind alone would merge a disjoint
		// pin whose near-miss then drops the whole cluster, stranding the valid pin).
		placed := false
		for _, c := range clusters {
			for _, b := range cands {
				if !compatibleBind(c, b) {
					continue
				}
				trial := maps.Clone(c)
				maps.Copy(trial, b)
				if !s.bindInstantiates(p, trial) {
					continue
				}
				maps.Copy(c, b)
				placed = true
				break
			}
			if placed {
				break
			}
		}
		if placed {
			continue
		}

		// New cluster: seed it from the first candidate that instantiates cleanly,
		// falling back to the first slot so a demand whose every pin is a near-miss
		// still forms a cluster that instantiateTemplate records as such, leaving the
		// demand to become an injector input.
		seed := cands[0]
		for _, b := range cands {
			if s.bindInstantiates(p, b) {
				seed = b
				break
			}
		}
		clusters = append(clusters, seed)
	}
	return clusters
}

// clusterCands collects one candidate binding per value result the demand
// unifies with, in result order, counting only bindings that pin at least one
// parameter. A demand pins at most one result slot: merging its per-slot
// bindings into one would, when two slots share an outer generic type (A[T] and
// A[U] both matched by A[int]), pin both parameters from a single demand and
// make the template produce that type twice, a spurious self-collision. matched
// reports whether any result unified at all: a parameter-free result yields no
// binding but still means the exact form matched, so the caller must not fall
// through to the dual.
func clusterCands(outs []types.Type, d types.Type, params map[*types.TypeParam]bool) (cands []map[*types.TypeParam]types.Type, matched bool) {
	for _, vo := range outs {
		b := map[*types.TypeParam]types.Type{}
		if gotypes.Unify(vo, d, params, b) {
			matched = true
			if len(b) > 0 {
				cands = append(cands, b)
			}
		}
	}
	return cands, matched
}

// jointClusterDone reports whether p has already been instantiated for a cluster
// with these pinned bindings. A cluster is identified by its pinned bindings
// alone: parameters it leaves to be freshly lifted would otherwise vary the
// instance identity every pass, so the unpinned slots are filled with a
// sentinel that is identical only to itself.
func (s *solver) jointClusterDone(p *discover.Provider, bind map[*types.TypeParam]types.Type) bool {
	m := s.jointDone[p]
	return m != nil && m.Contains(jointBindKey(p, bind))
}

func (s *solver) markJointClusterDone(p *discover.Provider, bind map[*types.TypeParam]types.Type) {
	m := s.jointDone[p]
	if m == nil {
		m = new(gotypes.Set[*types.Tuple])
		s.jointDone[p] = m
	}
	m.Add(jointBindKey(p, bind))
}

// jointUnpinned is a unique sentinel type standing for a cluster parameter no
// demand pinned. Being a fresh named type, it is identical only to itself, so it
// never coincides with a real pinned type.
var jointUnpinned types.Type = types.NewNamed(
	types.NewTypeName(token.NoPos, nil, "unpinned", nil),
	types.NewStruct(nil, nil), nil)

func jointBindKey(p *discover.Provider, bind map[*types.TypeParam]types.Type) *types.Tuple {
	ts := make([]types.Type, p.Tparams.Len())
	for i := range ts {
		if v, ok := bind[p.Tparams.At(i)]; ok {
			ts[i] = v
		} else {
			ts[i] = jointUnpinned
		}
	}
	return gotypes.ListKey(ts)
}
