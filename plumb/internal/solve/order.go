package solve

import (
	"fmt"
	"go/types"
	"slices"
	"strings"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/gotypes"
)

// topoOrder returns the instances in dependency order: every producer precedes
// the consumers of its outputs. Ties are broken by source position so the order
// is deterministic and independent of how packages and files were read. A cycle
// is reported as a located path.
//
// The selection loop below is O(V²·arity): each emitted node rescans all
// instances and re-checks their dependency sets. That is quadratic, but V is
// bounded by maxInstantiations and tiny for real sets (tens of instances; a
// deepening generic template trips maxTypeDepth long before it manufactures
// thousands of instances), so ordering stays sub-millisecond in practice. Kahn’s
// O(V+E) would be premature. The rescan also gives the deterministic
// position-ordered pick and feeds the cycle report for free.
func (s *solver) topoOrder(pl *Plan) ([]*Instance, *diag.Error) {
	insts := slices.SortedStableFunc(slices.Values(s.instances), func(a, b *Instance) int { return diag.CmpPos(a.pos, b.pos) })

	// deps[in] = the producer instances in must follow.
	deps := map[*Instance]map[*Instance]bool{}
	for _, in := range insts {
		deps[in] = map[*Instance]bool{}
		for _, ref := range pl.Args[in] {
			if ref.isParam {
				continue
			}
			// A self-edge (a provider consuming a type it produces, directly or
			// via the value/pointer bridge) is a single-node cycle, so it is kept.
			if prod, ok := s.supply.At(ref.SrcType); ok {
				deps[in][prod] = true
			}
		}
	}

	emitted := map[*Instance]bool{}
	var order []*Instance
	for len(order) < len(insts) {
		progress := false
		for _, in := range insts { // already position-sorted: deterministic pick
			if emitted[in] {
				continue
			}
			ready := true
			for d := range deps[in] {
				if !emitted[d] {
					ready = false
					break
				}
			}
			if ready {
				emitted[in] = true
				order = append(order, in)
				progress = true
				break
			}
		}
		if !progress {
			return nil, s.cycleError(pl, insts, emitted, deps)
		}
	}
	return order, nil
}

// cycleError finds a dependency cycle among the not-yet-emitted instances and
// reports it as a path, naming the type each provider needs from the next.
func (s *solver) cycleError(pl *Plan, insts []*Instance, emitted map[*Instance]bool, deps map[*Instance]map[*Instance]bool) *diag.Error {
	var remaining []*Instance
	for _, in := range insts {
		if !emitted[in] {
			remaining = append(remaining, in)
		}
	}
	const (
		white = iota // unvisited
		gray         // on the current DFS stack
		black        // fully explored
	)
	color := map[*Instance]int{}
	var stack []*Instance
	var cycle []*Instance

	var dfs func(in *Instance) bool
	dfs = func(in *Instance) bool {
		color[in] = gray
		stack = append(stack, in)
		// deterministic neighbor order
		var nbrs []*Instance
		for d := range deps[in] {
			if !emitted[d] {
				nbrs = append(nbrs, d)
			}
		}
		slices.SortStableFunc(nbrs, func(a, b *Instance) int { return diag.CmpPos(a.pos, b.pos) })
		for _, d := range nbrs {
			switch color[d] {
			case gray:
				// found a cycle: slice the stack from d onward
				for i, v := range slices.Backward(stack) {
					if v == d {
						cycle = append([]*Instance{}, stack[i:]...)
						return true
					}
				}
				return true
			case white:
				if dfs(d) {
					return true
				}
			}
		}
		color[in] = black
		stack = stack[:len(stack)-1]
		return false
	}

	for _, in := range remaining {
		if color[in] == white {
			if dfs(in) {
				break
			}
		}
	}

	if len(cycle) == 0 {
		// Should be unreachable: no progress implies a cycle exists.
		panic("plumb: ordering stalled without a detectable cycle")
	}
	// Render the cycle as "A needs T1 → B needs T2 → A", naming on each edge the
	// type the provider consumes from the next one round the loop.
	var parts []string
	for i, in := range cycle {
		next := cycle[(i+1)%len(cycle)]
		if t, bridged := s.edgeType(pl, in, next); t != nil {
			label := gotypes.TypeName(t)
			if bridged {
				label += " (bridged)"
			}
			parts = append(parts, fmt.Sprintf("%s needs %s", in.Prov.Name, label))
		} else {
			parts = append(parts, in.Prov.Name)
		}
	}
	parts = append(parts, cycle[0].Prov.Name)
	return diag.Errorf(cycle[0].pos, diag.ErrDependencyCycle, "set %q: %s", s.name, strings.Join(parts, " → "))
}

// edgeType returns the type consumer declares as the input that producer supplies
// (the dependency that makes consumer follow producer) and whether that input is
// satisfied through the value/pointer bridge. It reports the consumer’s demanded
// type (its InputSlot), not the producer’s SrcType, so the diagnostic names the
// type as written in the source; on a bridged edge the two are duals. The type is
// nil if no such edge exists (defensive; cycle edges always have one).
func (s *solver) edgeType(pl *Plan, consumer, producer *Instance) (types.Type, bool) {
	for i, ref := range pl.Args[consumer] {
		if ref.isParam {
			continue
		}
		if prod, ok := s.supply.At(ref.SrcType); ok && prod == producer {
			return consumer.Inputs[i].typ, ref.Coerce != CoerceNone
		}
	}
	return nil, false
}
