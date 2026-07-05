// Package solve is plumb's second phase: it turns the providers of one set into a
// fully resolved Plan, instantiating generics by demand, inferring the injector
// signature, ordering by dependency, and enforcing reachability and naming rules.
// It is a deterministic function of its providers; it reads no files and consults
// no clock.
//
// # Pipeline
//
// Set resolves one set by running run on a solver:
//
//   - classify: split providers into concrete / anchor / result-generic and
//     validate each result shape.
//   - seed: register the concrete providers and instantiate the anchors (lifting
//     their free parameters); registering an instance seeds demand from its inputs.
//   - solveDemands: the demand fixpoint. Each pass first works the worklist (a
//     demand a single result pins completely), then, if that made no progress, the
//     joint step (demands clustered across a template's results, lifting the
//     parameters no demand pins), until a whole pass adds nothing.
//   - checkUnusedTemplates / resolveFlexReceivers / finalize: reject a template no
//     consumer used, choose value-vs-pointer form for flexible struct receivers,
//     then resolve inputs, order instances by dependency, and build the Plan.
//
// # Pinning-revision restart
//
// A template whose pins arrive across different fixpoint rounds must be served by
// one instantiation at the union of those pins, not split into half-pinned calls.
// When a late pin would extend an existing instance's pinning, reviseOrSplit
// commits the union and the solve restarts: it signals with the errRestart
// sentinel, which rides the *diag.Error return channel up through
// instantiateTemplate → viableSingle/tryWorklist/tryJoint → solveDemands → run,
// compared by identity at each layer, until Set replays run on a fresh solver
// with the accumulated commitments in force. Resolution is therefore a pure
// function of (providers, commitments). A replay is the ordinary fixpoint, not a
// mutation of the aborted attempt, and maxRestarts bounds the loop. See
// revision.go.
//
// # Solver memo state
//
// Every memo table on solver is rebuilt fresh each attempt EXCEPT commitments and
// refusedUnions, which persist across restarts (a committed union must stay in
// force; a union already found non-viable must not re-trigger). The per-field roles
// are documented on the solver struct.
package solve

import (
	"fmt"
	"go/token"
	"go/types"
	"maps"
	"slices"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
)

// DestInfo describes the destination package as far as solve's collision checks
// and emit's qualification decisions need it. gen builds it; solve and emit
// (which imports solve) consume it.
type DestInfo struct {
	Scanned bool              // the destination is one of the scanned packages
	PkgName string            // the destination's package name
	Names   map[string]string // top-level identifier → base filename declaring it
	// Imports maps each import qualifier used in the destination's existing files
	// to the base name of a file using it. A generated function name (package
	// block) cannot coexist with an import qualifier (file block) of the same
	// name in any one file, so this feeds the set-name collision check.
	Imports map[string]string
}

// ReservedNames returns the base set of identifiers a name generated for this
// destination must avoid: every predeclared universe name and every destination
// top-level identifier (so, in same-package mode, a generated name never shadows
// a builtin the body emits or a provider the body calls unqualified). The
// receiver may be nil. Callers that reserve further categories (lifted-parameter
// names, import aliases, set names) add them to the returned map.
func (d *DestInfo) ReservedNames() map[string]bool {
	r := map[string]bool{}
	for _, n := range types.Universe.Names() {
		r[n] = true
	}
	if d != nil {
		for n := range d.Names {
			r[n] = true
		}
	}
	return r
}

// Bounds that protect against adversarial input. They are far above any real
// program and exist only so a non-terminating generic solve fails loudly rather
// than looping.
const (
	maxInstantiations = 5000
	maxTypeDepth      = 40
)

// Coerce describes a value/pointer bridge applied to an argument.
type Coerce int

// Coercion directions for the value/pointer bridge.
const (
	CoerceNone  Coerce = iota
	CoerceToPtr        // value T → *T via &x
	CoerceToVal        // *T → T via *x
)

// ArgRef records how one provider input is satisfied.
type ArgRef struct {
	isParam bool       // satisfied by an injector parameter
	SrcType types.Type // the type whose local holds the value
	Coerce  Coerce
}

// Input is one resolved injector parameter: its type and the source-name hint
// (the earliest consumer's parameter or field name, or "") used to name the
// generated parameter.
type Input struct {
	Type types.Type
	Name string
}

// Plan is the fully resolved wiring for one set, ready for emission.
type Plan struct {
	Name string         // the set name; the generated injector func is named this
	pos  token.Position // the set's source position, for diagnostics

	Order   []*Instance            // instances to emit, in dependency (topological) order
	Args    map[*Instance][]ArgRef // each instance's inputs, one ArgRef per input in order
	Inputs  []Input                // injector parameters, in signature order
	Outputs []types.Type           // value outputs, in signature order

	AnyCleanup      bool // some provider yields a cleanup, so the injector returns an aggregated cleanup func
	CleanupFailable bool // some cleanup can itself fail, so the aggregated cleanup returns an error
	Fallible        bool // some provider returns an error, so the injector returns a trailing error

	Lifted        []*types.TypeParam // free type parameters carried onto the generic injector header
	ProviderCount int                // source providers in the set, for the -v report; not len(Order), the instantiation count
}

// InputTypes returns the injector parameter types in signature order.
func (pl *Plan) InputTypes() []types.Type {
	ts := make([]types.Type, len(pl.Inputs))
	for i, in := range pl.Inputs {
		ts[i] = in.Type
	}
	return ts
}

// solver holds the mutable state of resolving one set.
type solver struct {
	name       string
	provs      []*discover.Provider
	destPath   string
	outputBase string
	dest       *DestInfo
	ctxt       *types.Context

	concrete []*Instance // concrete providers, instantiated during classification
	anchors  []*discover.Provider
	resGen   []*discover.Provider

	instances []*Instance
	supply    gotypes.Map[types.Type, *Instance] // type → producing instance (exact), by type identity

	demand       gotypes.Set[types.Type]                           // demanded types, by type identity
	done         map[*discover.Provider]*gotypes.Set[*types.Tuple] // provider → set of instantiated type-arg lists (keyed by ListKey)
	jointDone    map[*discover.Provider]*gotypes.Set[*types.Tuple] // provider → set of joint-cluster pinnings (keyed by jointBindKey)
	nearMissDone map[*discover.Provider]*gotypes.Set[*types.Tuple] // provider → set of bindings that failed to instantiate (constraint near-miss)
	skel         map[*discover.Provider]*Instance                  // provider → memoized own-params skeleton
	resGenUsed   map[*discover.Provider]bool
	nearMiss     map[*discover.Provider]error // template → the type-checker error a structural match failed to instantiate with

	liftedMeta  []liftedParam
	liftedNames map[string]bool
	avoid       map[string]bool // identifiers lifted params must not equal
	instCount   int

	// commitments and refusedUnions persist across pinning-revision restarts;
	// every other field is rebuilt fresh each attempt. commitments hold, per
	// template and in commitment order, the union pinnings a replay merges into
	// compatible bindings at instantiation time; refusedUnions memoizes unions
	// found non-viable or lift-entangled, so a refused merge never re-triggers.
	commitments   map[*discover.Provider][]map[*types.TypeParam]types.Type
	refusedUnions map[*discover.Provider]*gotypes.Set[*types.Tuple]

	// pinnings records each template instance's pinned slots as instantiated
	// (after commitment merging), so the revision trigger can tell an
	// instance's own lifted slot from a pin that merely mentions another
	// instance's lifted parameter.
	pinnings map[*Instance]map[*types.TypeParam]types.Type
}

// liftedParam records a lifted free type parameter and where it came from, for
// deterministic header ordering.
type liftedParam struct {
	tp  *types.TypeParam
	pos token.Position
	idx int // index within its origin provider's type parameters
}

// Set resolves one set into a plan. destPath is the destination import path
// (what is referenced unqualified) and outputBase is the base name of the file
// being overwritten, used by the set-name collision check (empty for stdout).
func Set(name string, provs []*discover.Provider, destPath, outputBase string, dest *DestInfo, ctxt *types.Context) (*Plan, *diag.Error) {
	// Sort a copy: a resolve entry point must not reorder its caller's slice.
	provs = slices.SortedStableFunc(slices.Values(provs), func(a, b *discover.Provider) int {
		return diag.CmpPos(a.Pos, b.Pos)
	})
	// Pinning revision replays the whole solve with the commitments in force,
	// on a fresh solver. Resolution stays a pure function of (providers,
	// commitments), and a replay is the ordinary fixpoint, not a mutation of
	// the aborted one. Only the commitments and refusals survive a restart.
	commitments := map[*discover.Provider][]map[*types.TypeParam]types.Type{}
	refused := map[*discover.Provider]*gotypes.Set[*types.Tuple]{}
	for range maxRestarts {
		s := &solver{
			name:          name,
			provs:         provs,
			destPath:      destPath,
			outputBase:    outputBase,
			dest:          dest,
			ctxt:          ctxt,
			done:          map[*discover.Provider]*gotypes.Set[*types.Tuple]{},
			jointDone:     map[*discover.Provider]*gotypes.Set[*types.Tuple]{},
			nearMissDone:  map[*discover.Provider]*gotypes.Set[*types.Tuple]{},
			skel:          map[*discover.Provider]*Instance{},
			resGenUsed:    map[*discover.Provider]bool{},
			nearMiss:      map[*discover.Provider]error{},
			liftedNames:   map[string]bool{},
			avoid:         dest.ReservedNames(),
			commitments:   commitments,
			refusedUnions: refused,
			pinnings:      map[*Instance]map[*types.TypeParam]types.Type{},
		}
		pl, err := s.run()
		if err == errRestart {
			continue
		}
		return pl, err
	}
	return nil, diag.Errorf(provs[0].Pos, diag.ErrNonTerminating, "set %q: pinning revision did not converge after %d restarts", name, maxRestarts)
}

func (s *solver) run() (*Plan, *diag.Error) {
	if err := s.classify(); err != nil {
		return nil, err
	}
	if err := s.seed(); err != nil {
		return nil, err
	}
	if err := s.solveDemands(); err != nil {
		return nil, err
	}
	if err := s.checkUnusedTemplates(); err != nil {
		return nil, err
	}
	s.resolveFlexReceivers()
	return s.finalize()
}

// classify splits providers into concrete / anchor / result-generic, validating
// each provider's result shape in source-position order (so the first invalid
// provider by position is the one reported, regardless of class) and rejecting
// bare type-parameter results.
func (s *solver) classify() *diag.Error {
	for _, p := range s.provs {
		if !p.Generic() {
			in, err, miss := instantiate(p, s.ctxt, nil)
			if err != nil {
				return err
			}
			if miss != nil {
				// Unreachable for valid input: a concrete provider has no type
				// arguments to validate.
				panic(fmt.Sprintf("plumb: provider %s could not be instantiated: %v", p.Name, miss))
			}
			s.concrete = append(s.concrete, in)
			continue
		}
		skel, err, miss := instantiate(p, s.ctxt, ownParams(p))
		if err != nil {
			return err
		}
		if miss != nil {
			// Unreachable for valid input: the skeleton instantiation uses the
			// provider's own type parameters, which always satisfy their constraints.
			panic(fmt.Sprintf("plumb: provider %s could not be analyzed: %v", p.Name, miss))
		}
		params := paramSet(p)
		resultGeneric := false
		for _, vo := range skel.valueOuts() {
			if gotypes.IsBareTypeParam(vo, params) {
				return diag.Errorf(p.Pos, diag.ErrBareTypeParamResult, "provider %s would match every demand", p.Name)
			}
			if gotypes.IsPointerToBareTypeParam(vo, params) {
				return diag.Errorf(p.Pos, diag.ErrBareTypeParamResult, "provider %s produces *T for a bare type parameter, matching every pointer demand and, through the value/pointer bridge, every value demand", p.Name)
			}
			if gotypes.MentionsParams(vo, params) {
				resultGeneric = true
			}
		}
		if resultGeneric {
			s.resGen = append(s.resGen, p)
		} else {
			s.anchors = append(s.anchors, p)
		}
	}
	return nil
}

// seed registers every concrete instance and instantiates every anchor template,
// registering their outputs and queuing their inputs.
func (s *solver) seed() *diag.Error {
	for _, in := range s.concrete {
		if err := s.addInstance(in); err != nil {
			return err
		}
	}
	for _, p := range s.anchors {
		targs := s.liftAll(p)
		in, err, miss := instantiate(p, s.ctxt, targs)
		if err != nil {
			return err
		}
		if miss != nil {
			// Unreachable for valid input: an anchor's free parameters are lifted
			// to fresh parameters that carry the original constraints.
			panic(fmt.Sprintf("plumb: anchor template %s could not be instantiated: %v", p.Name, miss))
		}
		if err := s.addInstance(in); err != nil {
			return err
		}
	}
	return nil
}

// solveDemands runs the demand-driven instantiation fixpoint. Each successful
// instantiation restarts the scan, and pendingDemands re-derives the whole demand
// set each pass, so this is O(V²) in the instantiation count, bounded by
// maxInstantiations and sub-millisecond for real sets (tens of instances), as with
// topoOrder; a rewrite to incremental tracking would be premature.
func (s *solver) solveDemands() *diag.Error {
	for {
		progress := false
		for _, d := range s.pendingDemands() {
			made, err := s.tryWorklist(d)
			if err != nil {
				return err
			}
			if made {
				progress = true
				break
			}
		}
		if progress {
			continue
		}
		made, err := s.tryJoint()
		if err != nil {
			return err
		}
		if made {
			continue
		}
		break
	}
	return nil
}

// pendingDemands returns the still-unsatisfied demanded types, sorted.
func (s *solver) pendingDemands() []types.Type {
	var out []types.Type
	for t := range s.demand.Elements {
		if !s.satisfied(t) {
			out = append(out, t)
		}
	}
	slices.SortStableFunc(out, gotypes.CmpType)
	return out
}

// satisfied reports whether demand t is met by an exact or value/pointer-dual
// supply.
func (s *solver) satisfied(t types.Type) bool {
	if _, ok := s.supply.At(t); ok {
		return true
	}
	if d, ok := gotypes.DualType(t); ok {
		if _, ok := s.supply.At(d); ok {
			return true
		}
	}
	return false
}

// tryWorklist attempts to satisfy demand d by instantiating a simple or
// monadic result-generic template fully pinned by a single demand.
func (s *solver) tryWorklist(d types.Type) (bool, *diag.Error) {
	candidates, err := s.viableSingle(d)
	if err != nil {
		return false, err
	}
	if len(candidates) == 0 {
		// Try the value/pointer dual: instantiate to the dual and coerce. The
		// retry runs whenever d's exact form yields no viable candidate: a
		// structural match whose constraint rejects the pin cannot build d, so it
		// must not block a dual producer that can.
		if dt, ok := gotypes.DualType(d); ok {
			candidates, err = s.viableSingle(dt)
			if err != nil {
				return false, err
			}
		}
	}
	if len(candidates) > 1 {
		return false, ambiguousTemplates(candidates)
	}
	if len(candidates) == 1 {
		c := candidates[0]
		return s.instantiateTemplate(c.prov, c.bind)
	}
	return false, nil // d stays open, for the joint step or the signature
}

// viableSingle returns, per template, the first value result (in result order)
// that unifies with d in a covering, viable binding, so a template contributes
// at most one candidate and ambiguity is decided among distinct templates. A
// covering result whose binding fails the template's constraint is recorded as a
// near-miss before the scan moves on: if the template ends the solve unused, the
// actionable reason is its constraint rejecting the pin, not the rival that
// served the demand.
func (s *solver) viableSingle(d types.Type) ([]matchCand, *diag.Error) {
	var out []matchCand
	for _, p := range s.resGen {
		params := paramSet(p)
		skel := s.skeleton(p)
		for _, vo := range skel.valueOuts() {
			b := map[*types.TypeParam]types.Type{}
			if !gotypes.Unify(vo, d, params, b) || !s.coversPinnable(p, b) {
				continue
			}
			if s.bindInstantiates(p, b) {
				out = append(out, matchCand{prov: p, bind: b})
				break
			}
			if _, err := s.instantiateTemplate(p, b); err != nil { // records the near-miss
				return nil, err
			}
		}
	}
	return out, nil
}

// skeleton returns p instantiated at its own type parameters, computed once per
// provider and reused across the fixpoint. The skeleton is a pure function of the
// provider, and classify already proved it instantiates (panicking otherwise), so
// no error or near-miss can occur here; callers only read it (valueOuts).
func (s *solver) skeleton(p *discover.Provider) *Instance {
	if in, ok := s.skel[p]; ok {
		return in
	}
	in, _, _ := instantiate(p, s.ctxt, ownParams(p))
	s.skel[p] = in
	return in
}

// compatibleBind reports whether two partial bindings agree on every type
// parameter they share, so they describe one instantiation rather than two.
func compatibleBind(a, b map[*types.TypeParam]types.Type) bool {
	for k, v := range b {
		if prev, ok := a[k]; ok && !types.Identical(prev, v) {
			return false
		}
	}
	return true
}

// bindInstantiates reports whether p cleanly instantiates at the given partial
// binding, lifting the unpinned parameters to test, then rolling those lifts back
// so the probe leaves no trace. It is used to decide whether merging two demands'
// pins into one cluster would poison it (a near-miss on one pin failing the whole
// instantiation). Both a constraint near-miss and a hard instantiation error count
// as “does not cleanly instantiate”; the real instantiation in tryJoint surfaces
// either for the pin's own cluster.
func (s *solver) bindInstantiates(p *discover.Provider, bind map[*types.TypeParam]types.Type) bool {
	mark := len(s.liftedMeta)
	targs := make([]types.Type, p.Tparams.Len())
	var fresh []int
	for i := range p.Tparams.Len() {
		tp := p.Tparams.At(i)
		if v, ok := bind[tp]; ok {
			targs[i] = v
		} else {
			targs[i] = s.liftOne(tp, p)
			fresh = append(fresh, i)
		}
	}
	s.setLiftedConstraints(p, targs, fresh)
	_, err, miss := instantiate(p, s.ctxt, targs)
	s.rollbackLifts(mark)
	return err == nil && miss == nil
}

// instantiateTemplate instantiates p at the given partial binding, lifting any
// unpinned parameters, and registers the resulting instance. A binding that a
// prior pass found to violate the template's constraint (a near-miss) is skipped:
// it would lift, fail, and roll back again on every fixpoint pass, leaking a
// phantom lifted parameter and inflating the instantiation count each time.
//
// Lifts mutate shared solver state (liftedMeta records the generated header's
// type-parameter list), so a binding that does not yield a retained instance
// must undo the lifts it speculatively appended; otherwise the failed attempt's
// parameters survive into the signature.
func (s *solver) instantiateTemplate(p *discover.Provider, bind map[*types.TypeParam]types.Type) (bool, *diag.Error) {
	bind = s.applyCommitment(p, bind)
	if err := s.reviseOrSplit(p, bind); err != nil {
		return false, err
	}
	if s.nearMissDoneAt(p, bind) {
		return false, nil
	}
	mark := len(s.liftedMeta)
	targs := make([]types.Type, p.Tparams.Len())
	var fresh []int
	for i := range p.Tparams.Len() {
		tp := p.Tparams.At(i)
		if v, ok := bind[tp]; ok {
			targs[i] = v
		} else {
			targs[i] = s.liftOne(tp, p)
			fresh = append(fresh, i)
		}
	}
	s.setLiftedConstraints(p, targs, fresh)
	if s.instanceDone(p, targs) {
		s.rollbackLifts(mark) // already instantiated; undo speculative lifts
		return false, nil
	}
	in, err, miss := instantiate(p, s.ctxt, targs)
	if err != nil {
		return false, err
	}
	if miss != nil {
		s.rollbackLifts(mark) // constraint near-miss; undo speculative lifts
		s.markNearMissDone(p, bind)
		// miss is the type-checker's constraint violation; keep the first one so an
		// otherwise-unused template reports the real reason rather than a misleading
		// "no consumer pins its result type".
		if s.nearMiss[p] == nil {
			s.nearMiss[p] = miss
		}
		return false, nil
	}
	for _, vo := range in.valueOuts() {
		if gotypes.TypeDepth(vo) > maxTypeDepth {
			return false, diag.Errorf(p.Pos, diag.ErrNonTerminating, "set %q: produced type %s exceeds depth %d", s.name, gotypes.TypeName(vo), maxTypeDepth)
		}
	}
	// Count only retained instantiations, so the bound reflects distinct instances
	// rather than (now-memoized) attempts.
	if s.instCount++; s.instCount > maxInstantiations {
		return false, diag.Errorf(p.Pos, diag.ErrNonTerminating, "set %q exceeded %d instantiations; a provider appears to manufacture unboundedly larger types", s.name, maxInstantiations)
	}
	s.markInstanceDone(p, targs)
	s.resGenUsed[p] = true
	s.pinnings[in] = maps.Clone(bind)
	if err := s.addInstance(in); err != nil {
		return false, err
	}
	return true, nil
}

// nearMissDoneAt reports whether binding (p, bind) was already found to violate
// p's constraint. The key matches jointClusterDone's: pinned slots by type,
// unpinned slots by the jointUnpinned sentinel, so it is stable across the
// freshly lifted parameters an unpinned slot takes on each pass.
func (s *solver) nearMissDoneAt(p *discover.Provider, bind map[*types.TypeParam]types.Type) bool {
	m := s.nearMissDone[p]
	return m != nil && m.Contains(jointBindKey(p, bind))
}

func (s *solver) markNearMissDone(p *discover.Provider, bind map[*types.TypeParam]types.Type) {
	m := s.nearMissDone[p]
	if m == nil {
		m = new(gotypes.Set[*types.Tuple])
		s.nearMissDone[p] = m
	}
	m.Add(jointBindKey(p, bind))
}

// instanceDone reports whether p has already been instantiated at this exact
// list of type arguments. Concrete arguments make the key stable across passes;
// freshly lifted parameters are distinct objects, so a lifted instantiation is
// never deduplicated here (the worklist and joint-cluster passes gate those).
func (s *solver) instanceDone(p *discover.Provider, targs []types.Type) bool {
	m := s.done[p]
	return m != nil && m.Contains(gotypes.ListKey(targs))
}

func (s *solver) markInstanceDone(p *discover.Provider, targs []types.Type) {
	m := s.done[p]
	if m == nil {
		m = new(gotypes.Set[*types.Tuple])
		s.done[p] = m
	}
	m.Add(gotypes.ListKey(targs))
}

// addInstance registers an instance's outputs and queues its inputs.
func (s *solver) addInstance(in *Instance) *diag.Error {
	for _, vo := range in.valueOuts() {
		if prev, ok := s.supply.At(vo); ok {
			if prev == in {
				// One provider returns two values of the same type; name it once
				// rather than reporting it as colliding “with itself”.
				return diag.Errorf(in.pos, diag.ErrAmbiguousProducer, "provider %s produces multiple values of type %s", in.Prov.Name, gotypes.TypeName(vo))
			}
			if prev.Prov == in.Prov {
				// Two instantiations of one template: a result whose type does not
				// depend on the type parameter is produced identically by each, so it
				// collides across instances. Reporting it through ambiguousProducers
				// would name the same provider and position twice ("Make and Make"),
				// which reads as colliding with itself; explain the real cause instead.
				return diag.Errorf(in.pos, diag.ErrAmbiguousProducer, "provider %s produces %s at more than one instantiation; this result does not depend on the type parameter: give it a type-parameter-dependent type or split the provider", in.Prov.Name, gotypes.TypeName(vo))
			}
			return ambiguousProducers(prev, in, vo)
		}
		s.supply.Set(vo, in)
	}
	for _, inp := range in.Inputs {
		t := inp.typ
		if inp.flexRecv {
			t = inp.flexBase // demand the value form; the dual bridge covers *S
		}
		s.demand.Add(t)
	}
	s.instances = append(s.instances, in)
	return nil
}

// checkUnusedTemplates rejects result-generic templates that nothing pinned.
func (s *solver) checkUnusedTemplates() *diag.Error {
	for _, p := range s.resGen {
		if s.resGenUsed[p] {
			continue
		}
		if reason := s.nearMiss[p]; reason != nil {
			return diag.Errorf(p.Pos, diag.ErrUnusedTemplate, "provider %s in set %q is never instantiated: %w", p.Name, s.name, reason)
		}
		if d, in, dual := s.shadowedDemand(p); in != nil {
			// A direct producer supplies d itself; a bridge producer supplies d's
			// dual and only the value/pointer bridge yields d, so name the dual, or
			// the message would claim the producer makes a type it does not.
			by := "produced by " + in.Prov.Name
			if dual != nil {
				by = fmt.Sprintf("covered via the value/pointer bridge by %s (which produces %s)", in.Prov.Name, gotypes.TypeName(dual))
			}
			return diag.Errorf(p.Pos, diag.ErrUnusedTemplate, "provider %s in set %q: %s is already %s, so no demand is left for it: remove the template or the other provider", p.Name, s.name, gotypes.TypeName(d), by)
		}
		return diag.Errorf(p.Pos, diag.ErrUnusedTemplate, "provider %s in set %q: no consumer pins its result type: add a consumer or remove it", p.Name, s.name)
	}
	return nil
}

// shadowedDemand explains the common reason a template goes unused: a consumer
// does pin its result, but another provider already supplies that type, so the
// demand was met before the template could run. It returns the shadowed demand,
// its producer, and the dual type that producer supplies when the value/pointer
// bridge (not an exact producer) is what preempted the demand, and nil when the
// producer supplies the demand directly. in is nil when no such producer exists
// (the result is genuinely unconsumed). Demands are scanned in CmpType order so
// the named producer is deterministic.
func (s *solver) shadowedDemand(p *discover.Provider) (types.Type, *Instance, types.Type) {
	params := paramSet(p)
	skel := s.skeleton(p)
	demands := slices.SortedStableFunc(s.demand.Elements, gotypes.CmpType)
	for _, d := range demands {
		if !s.templateServes(p, skel, params, d) {
			continue
		}
		// The template would have served d; find the producer that already does:
		// directly, or through the value/pointer bridge over a producer of d's dual
		// (a *Node[int] demand covered by a Node[int] producer, or vice versa).
		if in, ok := s.supply.At(d); ok && in.Prov != p {
			return d, in, nil
		}
		if dt, ok := gotypes.DualType(d); ok {
			if in, ok := s.supply.At(dt); ok && in.Prov != p {
				return d, in, dt
			}
		}
	}
	return nil, nil, nil
}

// templateServes reports whether template p could produce demand d, directly or
// through the value/pointer bridge by producing d's dual, mirroring tryWorklist's
// dual retry. It is the same unify+coversPinnable test viableSingle uses, so a hit
// means the template really would have been a candidate for d.
func (s *solver) templateServes(p *discover.Provider, skel *Instance, params map[*types.TypeParam]bool, d types.Type) bool {
	cands := []types.Type{d}
	if dt, ok := gotypes.DualType(d); ok {
		cands = append(cands, dt)
	}
	for _, cand := range cands {
		for _, vo := range skel.valueOuts() {
			b := map[*types.TypeParam]types.Type{}
			if gotypes.Unify(vo, cand, params, b) && s.coversPinnable(p, b) {
				return true
			}
		}
	}
	return false
}

// resolveFlexReceivers fixes each struct-field provider's receiver to value or
// pointer form, preferring the pointer when the set uses it elsewhere.
func (s *solver) resolveFlexReceivers() {
	// Collect the exact types the set uses for non-flexible inputs and all
	// outputs.
	var used gotypes.Set[types.Type]
	for _, in := range s.instances {
		for _, vo := range in.valueOuts() {
			used.Add(vo)
		}
		for _, inp := range in.Inputs {
			if !inp.flexRecv {
				used.Add(inp.typ)
			}
		}
	}
	for _, in := range s.instances {
		for i := range in.Inputs {
			if !in.Inputs[i].flexRecv {
				continue
			}
			ptr := types.NewPointer(in.Inputs[i].flexBase)
			if used.Contains(ptr) {
				in.Inputs[i].typ = ptr
			} else {
				in.Inputs[i].typ = in.Inputs[i].flexBase
			}
		}
	}
}

func (s *solver) checkInputType(in *Instance, t types.Type) *diag.Error {
	if gotypes.ContainsInvalid(t) {
		return diag.Errorf(in.pos, diag.ErrInvalidType, "provider %s references %s", in.Prov.Name, gotypes.TypeName(t))
	}
	return nil
}

// --- helpers -----------------------------------------------------------------

type matchCand struct {
	prov *discover.Provider
	bind map[*types.TypeParam]types.Type
}

func ambiguousTemplates(cs []matchCand) *diag.Error {
	slices.SortFunc(cs, func(a, b matchCand) int { return diag.CmpPos(a.prov.Pos, b.prov.Pos) })
	return diag.Errorf(cs[0].prov.Pos, diag.ErrAmbiguousTemplates, "templates %s and %s both match one demand",
		cs[0].prov.Name, cs[1].prov.Name)
}

// checkTemplateAmbiguity reports an error if a still-unsatisfied demand can be
// produced (cleanly, under each one's constraint) by two or more distinct
// templates. Such a demand has no single producer, and picking one by source
// order would be arbitrary. When no template's result matches the demand's exact
// form, the same test runs against its dual, and the rejection names the
// produced dual form. A template that merely lifts around a demand another
// provider already supplies is not in tension here: a supplied demand is no longer
// pending, so it never reaches this check. This is the joint-path analog of the
// worklist's single-result ambiguity test.
func (s *solver) checkTemplateAmbiguity() *diag.Error {
	for _, d := range s.pendingDemands() {
		producers, matched := s.viableProducers(d)
		if !matched {
			if dt, ok := gotypes.DualType(d); ok {
				producers, _ = s.viableProducers(dt)
				d = dt
			}
		}
		if len(producers) >= 2 {
			return ambiguousTemplatesFor(producers, d)
		}
	}
	return nil
}

// viableProducers returns the distinct templates with some value result that
// unifies with d in a viable binding, and whether any result of any template
// unified at all, viable or not: the caller falls through to the dual only
// when d's exact form matched nothing.
func (s *solver) viableProducers(d types.Type) (ps []*discover.Provider, matched bool) {
	for _, p := range s.resGen {
		params := paramSet(p)
		added := false
		for _, vo := range s.skeleton(p).valueOuts() {
			b := map[*types.TypeParam]types.Type{}
			if !gotypes.Unify(vo, d, params, b) {
				continue
			}
			matched = true
			if !added && s.bindInstantiates(p, b) {
				ps = append(ps, p)
				added = true
			}
		}
	}
	return ps, matched
}

func ambiguousTemplatesFor(ps []*discover.Provider, d types.Type) *diag.Error {
	slices.SortFunc(ps, func(a, b *discover.Provider) int { return diag.CmpPos(a.Pos, b.Pos) })
	return diag.Errorf(ps[0].Pos, diag.ErrAmbiguousTemplates, "templates %s and %s both produce %s",
		ps[0].Name, ps[1].Name, gotypes.TypeName(d))
}

func ambiguousProducers(a, b *Instance, t types.Type) *diag.Error {
	first, second := a, b
	if diag.CmpPos(b.pos, a.pos) < 0 {
		first, second = b, a
	}
	return diag.Errorf(first.pos, diag.ErrAmbiguousProducer, "type %s is produced by both %s (%s) and %s (%s); plumb never picks a winner",
		gotypes.TypeName(t), first.Prov.Name, first.pos, second.Prov.Name, second.pos)
}

// ownParams returns the provider's own type parameters as type arguments (used
// to build the generic skeleton).
func ownParams(p *discover.Provider) []types.Type {
	if !p.Generic() {
		return nil
	}
	out := make([]types.Type, p.Tparams.Len())
	for i := range p.Tparams.Len() {
		out[i] = p.Tparams.At(i)
	}
	return out
}

func paramSet(p *discover.Provider) map[*types.TypeParam]bool {
	m := map[*types.TypeParam]bool{}
	if p.Tparams == nil {
		return m
	}
	for tparam := range p.Tparams.TypeParams() {
		m[tparam] = true
	}
	return m
}

// coversPinnable reports whether bind pins every parameter that appears in a
// value result.
func (s *solver) coversPinnable(p *discover.Provider, bind map[*types.TypeParam]types.Type) bool {
	outs := s.skeleton(p).valueOuts()
	for tp := range paramSet(p) {
		if _, pinned := bind[tp]; pinned {
			continue
		}
		for _, vo := range outs {
			if gotypes.MentionsParam(vo, tp) {
				return false // tp appears in a value result but is unpinned
			}
		}
	}
	return true
}

// bridgeDir returns the coercion that turns a value into the want form, where the
// available value and want are a value/pointer dual; the direction is fixed by want.
// want is unaliased first so an alias to a pointer (type P = *T) is recognized as
// the pointer form, matching pointerElem/DualType, which decide the dual the same way.
func bridgeDir(want types.Type) Coerce {
	if _, ok := types.Unalias(want).(*types.Pointer); ok {
		return CoerceToPtr // want *N, have N
	}
	return CoerceToVal // want N, have *N
}

func (s *solver) setPos() token.Position {
	return slices.MinFunc(s.provs, func(a, b *discover.Provider) int {
		return diag.CmpPos(a.Pos, b.Pos)
	}).Pos
}
