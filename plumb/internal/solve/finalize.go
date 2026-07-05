package solve

import (
	"cmp"
	"go/token"
	"go/types"
	"slices"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/gotypes"
)

// injectorInput is a deduplicated injector parameter: its type, the source-name
// hint and position of its earliest consumer (for ordering and naming).
type injectorInput struct {
	typ  types.Type
	pos  token.Position
	name string
}

// finalize resolves every input, classifies the signature, runs reachability
// and reserved-name checks, orders the instances, and builds the plan.
func (s *solver) finalize() (*Plan, *diag.Error) {
	pl := &Plan{
		Name:          s.name,
		pos:           s.setPos(),
		Args:          map[*Instance][]ArgRef{},
		ProviderCount: len(s.provs),
	}

	var consumed gotypes.Set[types.Type]                  // value types some instance consumes
	var inputInfo gotypes.Map[types.Type, *injectorInput] // injector input type → its dedup record

	for _, in := range s.instances {
		refs := make([]ArgRef, len(in.Inputs))
		for i, inp := range in.Inputs {
			ref, err := s.resolveInput(in, inp)
			if err != nil {
				return nil, err
			}
			refs[i] = ref
			if ref.isParam {
				if cur, ok := inputInfo.At(ref.SrcType); !ok {
					inputInfo.Set(ref.SrcType, &injectorInput{typ: ref.SrcType, pos: in.pos, name: inp.Name})
				} else if diag.CmpPos(in.pos, cur.pos) < 0 {
					cur.pos = in.pos
					cur.name = inp.Name
				}
			} else {
				consumed.Add(ref.SrcType)
			}
		}
		pl.Args[in] = refs
	}

	// Value outputs: produced but not consumed.
	var outs []types.Type
	for _, in := range s.instances {
		for _, vo := range in.valueOuts() {
			if !consumed.Contains(vo) {
				outs = append(outs, vo)
			}
		}
	}
	sortTypesByProducer(outs, &s.supply)
	pl.Outputs = outs

	// Inputs: deduped, ordered by earliest consumer then type.
	inputs := make([]*injectorInput, 0, inputInfo.Len())
	for _, v := range inputInfo.All {
		inputs = append(inputs, v)
	}
	slices.SortFunc(inputs, func(a, b *injectorInput) int {
		return cmp.Or(diag.CmpPos(a.pos, b.pos), gotypes.CmpType(a.typ, b.typ))
	})
	for _, in := range inputs {
		pl.Inputs = append(pl.Inputs, Input{Type: in.typ, Name: in.name})
	}

	// Cleanup / error aggregation.
	for _, in := range s.instances {
		for _, r := range in.Results {
			switch r.Kind {
			case ResultKindCleanup:
				pl.AnyCleanup = true
				if r.Failable {
					pl.CleanupFailable = true
				}
			case ResultKindError:
				pl.Fallible = true
			}
		}
	}

	// Lifted parameters, in deterministic header order.
	pl.Lifted = s.orderedLifted()

	if err := s.checkReachability(pl); err != nil {
		return nil, err
	}
	if err := s.checkReservedAndCollision(pl); err != nil {
		return nil, err
	}

	order, err := s.topoOrder(pl)
	if err != nil {
		return nil, err
	}
	pl.Order = order
	return pl, nil
}

// resolveInput maps one provider input to its source: an exact producer, a
// value/pointer-bridged producer, or an injector parameter.
func (s *solver) resolveInput(in *Instance, inp InputSlot) (ArgRef, *diag.Error) {
	d := inp.typ
	if err := s.checkInputType(in, d); err != nil {
		return ArgRef{}, err
	}
	if _, ok := s.supply.At(d); ok {
		return ArgRef{SrcType: d, Coerce: CoerceNone}, nil
	}
	if dt, ok := gotypes.DualType(d); ok {
		if _, ok := s.supply.At(dt); ok {
			return ArgRef{SrcType: dt, Coerce: bridgeDir(d)}, nil
		}
	}
	return ArgRef{isParam: true, SrcType: d}, nil
}

func sortTypesByProducer(ts []types.Type, supply *gotypes.Map[types.Type, *Instance]) {
	producer := func(t types.Type) *Instance {
		in, _ := supply.At(t)
		return in
	}
	slices.SortStableFunc(ts, func(a, b types.Type) int {
		pa := producer(a)
		pb := producer(b)
		if pa != nil && pb != nil {
			if c := diag.CmpPos(pa.pos, pb.pos); c != 0 {
				return c
			}
		}
		return gotypes.CmpType(a, b)
	})
}
