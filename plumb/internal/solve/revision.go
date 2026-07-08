package solve

import (
	"go/types"
	"maps"
	"slices"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
)

// This file holds the pinning-revision machinery layered on top of the base
// demand fixpoint: the restart signal and bound, and the solver methods that
// decide when to commit a merged pinning and replay. Its persist-across-restart
// state (the commitments, refusedUnions, and pinnings fields) lives on the
// solver in solve.go; the replay loop that consumes errRestart is in Set.

// maxRestarts bounds pinning-revision restarts. Each restart commits a pinning
// strictly larger than an existing instance’s (or is preempted by a one-shot
// refusal), so the bound is far above any real program.
const maxRestarts = 100

// errRestart signals, up through the fixpoint to Set, that a new pinning was
// committed and resolution must replay from seed with it in force. It rides the
// *diag.Error return channel but is a control-flow signal, not a diagnostic:
// every layer compares it by identity (err == errRestart) and Set consumes it
// before it can escape. It is a diag.Sentinel rather than a zero-value *Error so
// that a stray format of an in-flight error prints the message instead of
// nil-panicking.
var errRestart = diag.Sentinel("plumb: internal: pinning-revision restart")

// applyCommitment merges bind with the first committed pinning of p it is
// compatible with, so a replay instantiates the union where the aborted run
// would have split. A binding incompatible with every commitment (a distinct
// instantiation of the same template) passes through unchanged.
func (s *solver) applyCommitment(p *discover.Provider, bind map[*types.TypeParam]types.Type) map[*types.TypeParam]types.Type {
	for _, c := range s.commitments[p] {
		if compatibleBind(c, bind) {
			merged := maps.Clone(c)
			maps.Copy(merged, bind)
			return merged
		}
	}
	return bind
}

// reviseOrSplit is the pinning-revision trigger: about to instantiate p at
// bind, it looks for an existing instance of p whose pinning is compatible
// with bind and would gain pins from it. When their union is viable and
// mentions no lifted parameter, the union becomes a committed pinning and the
// solve restarts to replay with it in force; otherwise the union is refused
// once (recorded so it never re-triggers) and the split stands, the same
// result the resolver reaches when no revision applies. Pinnings that conflict
// on a shared slot are ordinary multi-instantiation and never trigger.
func (s *solver) reviseOrSplit(p *discover.Provider, bind map[*types.TypeParam]types.Type) *diag.Error {
	for _, in := range s.instances {
		if in.Prov != p {
			continue
		}
		pin := s.pinnings[in]
		if pin == nil || !compatibleBind(pin, bind) {
			continue
		}
		union := maps.Clone(pin)
		maps.Copy(union, bind)
		if len(union) == len(pin) {
			continue // no new pins; instance dedup covers a subsumed binding
		}
		if s.unionRefused(p, union) {
			continue
		}
		entangled := slices.ContainsFunc(slices.Collect(maps.Values(union)), gotypes.ContainsTypeParam)
		if entangled || !s.bindInstantiates(p, union) {
			// A lifted-parameter pin cannot survive the restart, and a
			// constraint-violating union cannot instantiate; either way the
			// fallback is the split, never a rejection.
			s.markUnionRefused(p, union)
			continue
		}
		// A union extending an existing commitment upgrades it in place:
		// applyCommitment must find the largest form, or every replay merges
		// with the smaller one and re-derives the extension forever, a
		// staircase that runs a convergent program into the restart cap. The
		// upgrade keeps a template’s commitments pairwise incompatible, so the
		// first compatible commitment is the only one. A union the replay
		// re-derives without extending anything is already in force: commit
		// nothing, do not restart.
		for i, c := range s.commitments[p] {
			if !compatibleBind(c, union) {
				continue
			}
			upgraded := maps.Clone(c)
			maps.Copy(upgraded, union)
			if len(upgraded) == len(c) {
				return nil // union ⊆ c: already committed
			}
			s.commitments[p][i] = upgraded
			return errRestart
		}
		s.commitments[p] = append(s.commitments[p], union)
		return errRestart
	}
	return nil
}

func (s *solver) unionRefused(p *discover.Provider, union map[*types.TypeParam]types.Type) bool {
	m := s.refusedUnions[p]
	return m != nil && m.Contains(jointBindKey(p, union))
}

func (s *solver) markUnionRefused(p *discover.Provider, union map[*types.TypeParam]types.Type) {
	m := s.refusedUnions[p]
	if m == nil {
		m = new(gotypes.Set[*types.Tuple])
		s.refusedUnions[p] = m
	}
	m.Add(jointBindKey(p, union))
}
