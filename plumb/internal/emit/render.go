package emit

import (
	"fmt"
	"go/types"
	"strings"

	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/solve"
)

// renderInstance renders the value-producing expression for one instance: the
// provider call or reference, given a qualifier and the already-coerced argument
// expressions in input order. It is the single place provider kinds turn into Go
// source: solve decides the wiring, emit renders it.
func renderInstance(in *solve.Instance, q *qualifier, args []string) string {
	p := in.Prov
	switch p.Kind {
	case discover.KindFunc:
		var b strings.Builder
		b.WriteString(q.objQual(p.Fn))
		writeTypeArgs(&b, q, in.Targs)
		b.WriteByte('(')
		writeArgs(&b, args, p.Fn.Type().(*types.Signature).Variadic())
		b.WriteByte(')')
		return b.String()
	case discover.KindMethod:
		// args[0] is the receiver expression. Variadic-ness is stable under
		// instantiation, so the origin method signature is authoritative.
		var b strings.Builder
		b.WriteString(asReceiver(args[0]))
		b.WriteByte('.')
		b.WriteString(p.Fn.Name())
		b.WriteByte('(')
		writeArgs(&b, args[1:], p.Fn.Type().(*types.Signature).Variadic())
		b.WriteByte(')')
		return b.String()
	case discover.KindSymbol:
		return q.objQual(p.Sym)
	case discover.KindConvert:
		target := q.typeString(p.ConvertTo)
		if convNeedsParens(p.ConvertTo) {
			target = "(" + target + ")"
		}
		return fmt.Sprintf("%s(%s)", target, args[0])
	case discover.KindField:
		return asReceiver(args[0]) + "." + p.Sym.Name()
	case discover.KindStruct:
		// The value result is the declared type, rendered as written; the inputs
		// carry the exported field names in composite-literal order.
		var b strings.Builder
		b.WriteString(q.typeString(in.Results[0].Typ))
		b.WriteByte('{')
		for i, a := range args {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(in.Inputs[i].Name)
			b.WriteString(": ")
			b.WriteString(a)
		}
		b.WriteByte('}')
		return b.String()
	}
	panic(fmt.Sprintf("plumb: unhandled provider kind %d", p.Kind))
}

// asReceiver parenthesizes a selector receiver when it is a value/pointer bridge.
// The bridge renders *T → T as "*x" and T → *T as "&x"; used bare as a
// receiver, "*x.M()" / "&x.M()" parse as "*(x.M())" / "&(x.M())", so they must
// become "(*x).M()" / "(&x).M()". A plain local name is already a primary
// expression and binds the selector correctly.
func asReceiver(expr string) string {
	if strings.HasPrefix(expr, "*") || strings.HasPrefix(expr, "&") {
		return "(" + expr + ")"
	}
	return expr
}

// convNeedsParens reports whether a conversion to t must parenthesize the type so
// that T(x) is not mis-parsed. T(x) goes wrong two ways: a leading operator (a
// pointer *T reads as a dereference, a receive-only <-chan T as a receive), and a
// rendered form ending in a result-less func (func()(x) reparses as func() (x),
// its result list swallowing the argument). The latter is reachable through slice,
// array, map-value, channel, and pointer tails. gofmt rescues neither: it leaves
// *T(x) as a dereference and breaks func()-tailed conversions rather than fixing
// them. (It does add parens around a top-level func *with* a result, e.g.
// func() int(x) → (func() int)(x), so that one case needs no help here.) Every
// other target (named/alias/interface/struct/basic types, and composites not
// ending in a bare func) renders unambiguously.
func convNeedsParens(t types.Type) bool {
	switch u := t.(type) {
	case *types.Pointer:
		return true
	case *types.Chan:
		if u.Dir() == types.RecvOnly {
			return true
		}
	}
	return endsInResultlessFunc(t)
}

// endsInResultlessFunc reports whether t’s rendered form ends in a result-less
// func type, reached directly or through a slice/array/map-value/channel/pointer
// element or a func’s sole result. (A two-or-more result list renders in
// parentheses, which closes the type, so it does not end in a bare func.) A named
// or alias type renders as its name, so it is opaque and stops the recursion.
func endsInResultlessFunc(t types.Type) bool {
	switch u := t.(type) {
	case *types.Signature:
		switch u.Results().Len() {
		case 0:
			return true
		case 1:
			return endsInResultlessFunc(u.Results().At(0).Type())
		default:
			return false
		}
	case *types.Slice:
		return endsInResultlessFunc(u.Elem())
	case *types.Array:
		return endsInResultlessFunc(u.Elem())
	case *types.Map:
		return endsInResultlessFunc(u.Elem())
	case *types.Chan:
		return endsInResultlessFunc(u.Elem())
	case *types.Pointer:
		return endsInResultlessFunc(u.Elem())
	}
	return false
}

func writeArgs(b *strings.Builder, args []string, variadic bool) {
	for i, a := range args {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(a)
		if variadic && i == len(args)-1 {
			b.WriteString("...")
		}
	}
}

func writeTypeArgs(b *strings.Builder, q *qualifier, targs []types.Type) {
	if len(targs) == 0 {
		return
	}
	b.WriteByte('[')
	for i, t := range targs {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(q.typeString(t))
	}
	b.WriteByte(']')
}
