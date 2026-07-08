package discover

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Package is one loaded, type-checked Go package as the core consumes it. The
// loaders (the real go/packages loader and the in-memory test loader) produce
// values of this shape; the core never loads anything itself. It is discover’s
// input contract (the first phase defines what a loaded package must carry),
// and the later phases and the cli read the same shape.
type Package struct {
	// PkgPath is the package’s import path (e.g. "example.com/app"). It is the
	// identity used to decide qualification against the destination.
	PkgPath string
	// Name is the package clause name (e.g. "app").
	Name string
	// Fset is the file set the package’s positions live in. All scanned
	// packages share one Fset so positions are globally comparable.
	Fset *token.FileSet
	// Syntax is the parsed source of the package (with comments), used to find
	// directives and the declarations they attach to.
	Syntax []*ast.File
	// Types is the type-checked package, used to resolve provider and signature
	// types.
	Types *types.Package
	// Info carries the type information gathered during checking (Defs, Uses,
	// Types, Selections).
	Info *types.Info
}

// ProviderKind enumerates the surface forms a directive can tag. They all reduce
// to the same abstraction (inputs, results, a render rule), which is what lets
// the wiring code in the solve package ignore the kind almost entirely.
type ProviderKind int

// The provider kinds, one per surface form a directive can tag.
const (
	KindFunc    ProviderKind = iota // package-level function
	KindMethod                      // method on a concrete or interface receiver
	KindSymbol                      // reference to a package-level var or typed const
	KindField                       // struct field
	KindStruct                      // struct type that builds itself via a composite literal
	KindConvert                     // conversion provider: var _ T = expr
)

func (k ProviderKind) String() string {
	switch k {
	case KindFunc:
		return "func"
	case KindMethod:
		return "method"
	case KindSymbol:
		return "symbol"
	case KindField:
		return "field"
	case KindStruct:
		return "struct"
	case KindConvert:
		return "convert"
	}
	return "?"
}

// Provider is one discovered directive occurrence. A declaration tagged with N
// directives yields N providers (one per set name), so each provider belongs to
// exactly one set.
type Provider struct {
	SetName string
	Kind    ProviderKind
	Pos     token.Position // identifying position, for ordering and diagnostics
	Name    string         // human-readable name for the -v report and diagnostics
	Pkg     *types.Package // the declaring package

	// Origin objects, by kind. Only the fields relevant to Kind are set.
	Fn          *types.Func  // KindFunc, KindMethod
	Sym         types.Object // KindSymbol (a var or typed const), KindField (the field)
	Owner       types.Type   // KindMethod (receiver), KindField (struct): a *Named, or a *Alias when the member is on an alias to an anonymous composite
	Declared    types.Type   // KindStruct: the directive’s declared type (a *Named or *Alias), rendered as written
	ConvertTo   types.Type   // KindConvert: the target type produced
	ConvertFrom types.Type   // KindConvert: the source type consumed

	// Tparams are the provider’s type parameters (from the function or the owner
	// type). nil/empty means the provider is concrete.
	Tparams *types.TypeParamList
}

// Generic reports whether the provider is a template.
func (p *Provider) Generic() bool {
	return p.Tparams != nil && p.Tparams.Len() > 0
}
