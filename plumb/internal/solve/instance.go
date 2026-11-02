package solve

import (
	"fmt"
	"go/token"
	"go/types"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
)

// ResultKind classifies one entry of a provider’s result tuple.
type ResultKind int

// The result kinds for classifying entries of a provider’s result tuple.
const (
	ResultKindValue   ResultKind = iota // an ordinary value output
	ResultKindCleanup                   // a bare func() / func() error teardown hook
	ResultKindError                     // the predeclared error channel
)

// ResultSlot is one classified entry of a provider’s result tuple.
type ResultSlot struct {
	Kind     ResultKind
	Typ      types.Type // ResultKindValue: the value type; otherwise nil
	Failable bool       // ResultKindCleanup: true if func() error
}

// InputSlot is one consumed type, in call order.
type InputSlot struct {
	typ types.Type

	// Name is the source parameter/field name, a hint for naming an injector
	// input. For a struct-type provider it is the exported field name, which emit
	// renders into the composite literal.
	Name string

	// flexRecv marks the implicit receiver of a struct-field provider, whose
	// value-vs-pointer form is chosen during solving from the rest of the set.
	flexRecv bool
	flexBase types.Type // the struct type (value form) for a flexible receiver: a *Named, or a *Alias when the field is on an alias to an anonymous struct
}

// Instance is a concrete (fully instantiated) provider ready for wiring: pure
// data describing what to call with what. emit turns it into Go source. For a
// non-generic provider there is exactly one; for a template there is one per
// pinning.
type Instance struct {
	Prov    *discover.Provider // the source provider this instantiates
	Rargs   []types.Type       // receiver type arguments, aligned with Prov.Rparams
	Targs   []types.Type       // explicit type arguments, aligned with Prov.Tparams
	Inputs  []InputSlot        // consumed types, in call order
	Results []ResultSlot       // classified results, in result order

	pos token.Position
}

// typeArgs yields receiver arguments followed by declaration arguments.
func (in *Instance) typeArgs(yield func(types.Type) bool) {
	for _, args := range [][]types.Type{in.Rargs, in.Targs} {
		for _, arg := range args {
			if !yield(arg) {
				return
			}
		}
	}
}

// valueOuts returns the value-output types of the instance, in result order.
func (in *Instance) valueOuts() []types.Type {
	var out []types.Type
	for _, r := range in.Results {
		if r.Kind == ResultKindValue {
			out = append(out, r.Typ)
		}
	}
	return out
}

// InputTypes returns the instance’s input types, in call order (for the report).
func (in *Instance) InputTypes() []types.Type {
	var out []types.Type
	for _, i := range in.Inputs {
		out = append(out, i.typ)
	}
	return out
}

// instantiate produces the concrete instance of p for the given type arguments
// (receiver arguments followed by declaration arguments, empty for a concrete
// provider): it classifies the provider’s inputs and results, with no rendering:
// emit decides how each kind is written. It returns
// a located error (*diag.Error) for user-visible problems (e.g. multiple error
// results), and a non-nil miss error (the type-checker’s own *types.ArgumentError)
// when the instantiation is infeasible because arguments violate a type-parameter
// constraint, which the solver treats as a near-miss non-match. Both nil means a
// matched instance. The per-kind helpers panic on a post-instantiation structural
// failure (a method/field that vanished), which is unreachable for valid input.
func instantiate(p *discover.Provider, ctxt *types.Context, args []types.Type) (in *Instance, err *diag.Error, miss error) {
	rargs, targs := args[:p.Rparams.Len()], args[p.Rparams.Len():]
	in = &Instance{Prov: p, Rargs: rargs, Targs: targs, pos: p.Pos}
	switch p.Kind {
	case discover.KindFunc, discover.KindMethod:
		var sig *types.Signature
		var e error
		if p.Kind == discover.KindMethod {
			var receiver types.Type
			receiver, sig, e = methodSignature(p, ctxt, rargs, targs)
			in.Inputs = append(in.Inputs, InputSlot{typ: receiver})
		} else {
			sig, e = instSignature(ctxt, p.Fn.Type().(*types.Signature), targs)
		}
		if e != nil {
			return nil, nil, e
		}
		classifyParams(in, sig)
		if err := classifyResults(p, in, sig); err != nil {
			return nil, err, nil
		}
	case discover.KindSymbol:
		in.Results = []ResultSlot{{Kind: ResultKindValue, Typ: p.Sym.Type()}}
	case discover.KindConvert:
		in.Inputs = []InputSlot{{typ: p.ConvertFrom}}
		in.Results = []ResultSlot{{Kind: ResultKindValue, Typ: p.ConvertTo}}
	case discover.KindField:
		fieldType, recvBase, e := fieldTypes(p, ctxt, rargs)
		if e != nil {
			return nil, nil, e
		}
		in.Inputs = []InputSlot{{flexRecv: true, flexBase: recvBase, typ: recvBase}}
		in.Results = []ResultSlot{{Kind: ResultKindValue, Typ: fieldType}}
	case discover.KindStruct:
		declared, fields, e := structFields(p, ctxt, targs)
		if e != nil {
			return nil, nil, e
		}
		for _, f := range fields {
			in.Inputs = append(in.Inputs, InputSlot{typ: f.Type(), Name: f.Name()})
		}
		in.Results = []ResultSlot{{Kind: ResultKindValue, Typ: declared}}
	default:
		panic(fmt.Sprintf("plumb: unhandled provider kind %d", p.Kind))
	}
	return in, nil, nil
}

// instSignature instantiates a generic signature, or returns it unchanged when
// there are no type args. A non-nil error is the constraint near-miss reason.
func instSignature(ctxt *types.Context, sig *types.Signature, targs []types.Type) (*types.Signature, error) {
	if len(targs) == 0 {
		return sig, nil
	}
	t, e := types.Instantiate(ctxt, sig, targs, true)
	if e != nil {
		return nil, e
	}
	return t.(*types.Signature), nil
}

// methodSignature instantiates the receiver, then the method on that receiver.
// Go substitutes receiver arguments into the method’s constraints before checking
// its own arguments. Interface methods use the declared owner as their input,
// preserving aliases to anonymous interfaces.
func methodSignature(p *discover.Provider, ctxt *types.Context, rargs, targs []types.Type) (types.Type, *types.Signature, error) {
	owner := p.Owner
	if len(rargs) > 0 {
		var err error
		owner, err = types.Instantiate(ctxt, gotypes.GenericOrigin(owner), rargs, true)
		if err != nil {
			return nil, nil, err
		}
	}
	if iface, ok := owner.Underlying().(*types.Interface); ok {
		for method := range iface.Methods() {
			if method.Name() == p.Fn.Name() {
				return owner, method.Type().(*types.Signature), nil
			}
		}
	} else {
		for method := range owner.(*types.Named).Methods() {
			if method.Name() == p.Fn.Name() {
				sig, err := instSignature(ctxt, method.Type().(*types.Signature), targs)
				if err != nil {
					return nil, nil, err
				}
				return sig.Recv().Type(), sig, nil
			}
		}
	}
	panic(fmt.Sprintf("plumb: method %s of provider %s vanished after instantiation", p.Fn.Name(), p.Name))
}

// fieldTypes returns the field’s type and the (possibly instantiated) receiver
// type for a field provider: a *types.Named, or a *types.Alias when the field is
// on an alias to an anonymous struct. A non-nil error is the constraint near-miss
// reason.
func fieldTypes(p *discover.Provider, ctxt *types.Context, rargs []types.Type) (types.Type, types.Type, error) {
	recv := p.Owner
	if len(rargs) > 0 {
		t, e := types.Instantiate(ctxt, gotypes.GenericOrigin(p.Owner), rargs, true)
		if e != nil {
			return nil, nil, e
		}
		recv = t
	}
	st, ok := recv.Underlying().(*types.Struct)
	if !ok {
		panic(fmt.Sprintf("plumb: receiver of field provider %s is not a struct after instantiation", p.Name))
	}
	for f := range st.Fields() {
		if f.Name() == p.Sym.Name() {
			return f.Type(), recv, nil
		}
	}
	panic(fmt.Sprintf("plumb: field %s of provider %s vanished after instantiation", p.Sym.Name(), p.Name))
}

// structFields returns the declared type to render (the directive’s named or
// alias type, instantiated at targs) and its exported struct fields. The alias is
// kept as written so the composite literal names it (and an exported alias to an
// unexported struct stays reachable across packages); the fields come from its
// flattened underlying struct. A non-nil error is the constraint near-miss reason.
func structFields(p *discover.Provider, ctxt *types.Context, targs []types.Type) (declared types.Type, fields []*types.Var, _ error) {
	declared = p.Declared
	if len(targs) > 0 {
		t, e := types.Instantiate(ctxt, gotypes.GenericOrigin(p.Declared), targs, true)
		if e != nil {
			return nil, nil, e
		}
		declared = t
	}
	st, ok := declared.Underlying().(*types.Struct)
	if !ok {
		panic(fmt.Sprintf("plumb: struct-type provider %s is not a struct after instantiation", p.Name))
	}
	for f := range st.Fields() {
		if f.Exported() {
			fields = append(fields, f)
		}
	}
	return declared, fields, nil
}

// classifyParams appends the signature parameters as inputs. A variadic parameter
// needs nothing special here: the tuple already types it as its slice ([]T), and
// the call-site spread (x...) is emit’s job.
func classifyParams(in *Instance, sig *types.Signature) {
	for v := range sig.Params().Variables() {
		in.Inputs = append(in.Inputs, InputSlot{typ: v.Type(), Name: v.Name()})
	}
}

// classifyResults classifies the signature results into value/cleanup/error
// slots, enforcing the at-most-one-error rule.
func classifyResults(p *discover.Provider, in *Instance, sig *types.Signature) *diag.Error {
	res := sig.Results()
	errCount := 0
	for v := range res.Variables() {
		t := v.Type()
		switch {
		case gotypes.IsErrorType(t):
			errCount++
			if errCount > 1 {
				return diag.Errorf(p.Pos, diag.ErrMultipleErrors, "provider %s returns more than one error result; at most one is supported", p.Name)
			}
			in.Results = append(in.Results, ResultSlot{Kind: ResultKindError})
		case gotypes.IsBareCleanup(t):
			in.Results = append(in.Results, ResultSlot{Kind: ResultKindCleanup, Failable: false})
		case gotypes.IsFailableCleanup(t):
			in.Results = append(in.Results, ResultSlot{Kind: ResultKindCleanup, Failable: true})
		default:
			in.Results = append(in.Results, ResultSlot{Kind: ResultKindValue, Typ: t})
		}
	}
	return nil
}
