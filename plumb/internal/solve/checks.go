package solve

import (
	"fmt"
	"go/types"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
)

// checkReachability checks the names written in generated signatures and
// expressions. They must be accessible without changing the type’s identities.
func (s *solver) checkReachability(pl *Plan) *diag.Error {
	for _, in := range s.instances {
		p := in.Prov
		// A type plumb must render that does not type-check is fatal (distinct
		// from a tolerated stale-output type error, which never reaches here).
		for ta := range in.typeArgs {
			if gotypes.ContainsInvalid(ta) {
				return diag.Errorf(p.Pos, diag.ErrInvalidType, "provider %s is instantiated at %s", p.Name, gotypes.TypeName(ta))
			}
		}
		if p.Kind == discover.KindConvert && gotypes.ContainsInvalid(p.ConvertTo) {
			return diag.Errorf(p.Pos, diag.ErrInvalidType, "conversion %s", p.Name)
		}
		// A provider referenced across the boundary must itself be exported. A
		// same-package symbol is named unqualified, so this gate is correct; but
		// the type checks below are NOT gated, because a type a provider renders
		// (a type argument, a conversion target) may come from a third package
		// even when the provider itself is in the destination.
		if p.Pkg != nil && p.Pkg.Path() != s.destPath {
			if obj := crossBoundarySymbol(p); obj != nil && !obj.Exported() {
				return diag.Errorf(p.Pos, diag.ErrUnexportedProvider, "provider %s referenced from package %q; export it or generate into its own package", p.Name, s.destPath)
			}
		}
		// Check the same arguments rendering names. Receiver arguments carried
		// only by local values remain exempt; signature types are checked below.
		for _, ta := range in.Targs {
			if obj := findUnreachable(ta, s.destPath); obj != nil {
				return diag.Errorf(p.Pos, diag.ErrUnreachableType, "type argument %s of provider %s, from package %q (%s is unexported)", gotypes.TypeName(ta), p.Name, obj.Pkg().Path(), obj.Name())
			}
		}
		// The conversion’s target type is written in the body.
		if p.Kind == discover.KindConvert {
			if obj := findUnreachable(p.ConvertTo, s.destPath); obj != nil {
				return diag.Errorf(p.Pos, diag.ErrUnreachableType, "conversion target type %s from package %q (%s is unexported)", gotypes.TypeName(p.ConvertTo), obj.Pkg().Path(), obj.Name())
			}
		}
	}
	// 2. Signature types must type-check and be reachable.
	for _, in := range pl.Inputs {
		t := in.Type
		if gotypes.ContainsInvalid(t) {
			return diag.Errorf(pl.pos, diag.ErrInvalidType, "injector input type %s", gotypes.TypeName(t))
		}
		if obj := findUnreachable(t, s.destPath); obj != nil {
			return diag.Errorf(pl.pos, diag.ErrUnreachableType, "injector input type %s, from package %q (%s is unexported)", gotypes.TypeName(t), obj.Pkg().Path(), obj.Name())
		}
	}
	for _, t := range pl.Outputs {
		if gotypes.ContainsInvalid(t) {
			return diag.Errorf(pl.pos, diag.ErrInvalidType, "injector output type %s", gotypes.TypeName(t))
		}
		if obj := findUnreachable(t, s.destPath); obj != nil {
			return diag.Errorf(pl.pos, diag.ErrUnreachableType, "injector output type %s, from package %q (%s is unexported)", gotypes.TypeName(t), obj.Pkg().Path(), obj.Name())
		}
	}
	// 3. Lifted free parameters’ constraints appear in the header, except one that
	// emit collapses to the bare "any", which names no package (the same predicate
	// gates emit’s rendering and import collection, so the two cannot drift).
	for _, tp := range pl.Lifted {
		if gotypes.ConstraintCollapsesToAny(tp.Constraint()) {
			continue
		}
		// A constraint that did not type-check (a typo’d or undefined name in
		// tolerated-invalid input) would render as "invalid type" in the header and
		// fail to format; catch it here, as inputs and outputs are checked above.
		if gotypes.ContainsInvalid(tp.Constraint()) {
			return diag.Errorf(pl.pos, diag.ErrInvalidType, "lifted type parameter %s has a constraint that does not type-check", tp.Obj().Name())
		}
		if obj := findUnreachable(tp.Constraint(), s.destPath); obj != nil {
			return diag.Errorf(pl.pos, diag.ErrUnreachableType, "lifted type parameter %s has a constraint from package %q (%s is unexported)", tp.Obj().Name(), obj.Pkg().Path(), obj.Name())
		}
	}
	return nil
}

// crossBoundarySymbol returns the named symbol a provider renders across the
// boundary, or nil if the provider renders no qualified symbol of its own.
func crossBoundarySymbol(p *discover.Provider) types.Object {
	switch p.Kind {
	case discover.KindFunc, discover.KindMethod:
		return p.Fn
	case discover.KindSymbol, discover.KindField:
		return p.Sym
	case discover.KindStruct:
		return gotypes.TypeNameOf(p.Declared)
	case discover.KindConvert:
		// The conversion’s only foreign reference is the target type, checked by
		// the type walk; it renders no qualified symbol of its own.
		return nil
	}
	panic(fmt.Sprintf("plumb: unhandled provider kind %d", p.Kind))
}

// findUnreachable returns the first foreign, unexported object written in t’s
// source representation, or nil. Repeating a member name from another package
// would change its identity.
func findUnreachable(t types.Type, destPath string) types.Object {
	var found types.Object
	gotypes.WalkObjects(t, func(obj types.Object) bool {
		if found != nil {
			return false
		}
		if pkg := obj.Pkg(); pkg != nil && pkg.Path() != destPath && !obj.Exported() {
			found = obj
			return false
		}
		return true
	})
	return found
}

// checkReservedAndCollision enforces the reserved-name and set-name collision
// rules.
func (s *solver) checkReservedAndCollision(pl *Plan) *diag.Error {
	nonEmpty := len(pl.Inputs) > 0 || len(pl.Outputs) > 0 || pl.AnyCleanup || pl.Fallible || len(pl.Lifted) > 0

	if s.name == "init" && nonEmpty {
		return diag.Errorf(pl.pos, diag.ErrReservedName, "set %q would generate a func init; init must take no parameters and return nothing", s.name)
	}
	if s.name == "main" && s.dest.PkgName == "main" && nonEmpty {
		return diag.Errorf(pl.pos, diag.ErrReservedName, "set %q would generate a func main in package main", s.name)
	}

	// A generated function is a package-block declaration, so naming a set after a
	// predeclared identifier (int, new, error, ...) shadows that identifier for the
	// entire destination package, breaking the generated body and any sibling file
	// that uses the builtin. In separate-package mode there are no siblings today,
	// but the name remains a latent trap, and the function name cannot be aliased.
	// So reject it outright in both modes, regardless of whether the generated
	// signature happens to mention the identifier.
	if gotypes.IsUniverseName(s.name) {
		return diag.Errorf(pl.pos, diag.ErrShadowsPredeclared, "set %q would generate a top-level func %q shadowing the builtin for the whole destination package; rename the set", s.name, s.name)
	}

	// The reverse shadowing direction (a destination declaration capturing a
	// predeclared name the generated body spells unqualified) is a property of
	// the destination, not of any one set, so gen checks it once before solving
	// (see buildDestInfo).

	// init is never entered into the package scope, so it never collides.
	if s.name == "init" || !s.dest.Scanned {
		return nil
	}
	// Collision with an existing package-level declaration. gen already drops the
	// overwritten file’s own declarations from dest.Names, so a hit here is always a
	// genuine collision; but for standard output no file is overwritten and every
	// declaration is recorded, so the check cannot tell a real collision from the
	// set’s own future output and is skipped.
	if s.outputBase != "" {
		if base, ok := s.dest.Names[s.name]; ok {
			return diag.Errorf(pl.pos, diag.ErrSetNameCollision, "set %q in package %q (%s); rename the set", s.name, s.destPath, base)
		}
	}
	// Collision with an import qualifier (file block) in a hand-written file: Go
	// forbids one identifier in both the file and package block. gen drops the
	// overwritten file’s own qualifiers from dest.Imports (plumb rewrites them), so
	// a hit is a hand-written sibling’s qualifier, a genuine collision, checked
	// even for standard output.
	if base, ok := s.dest.Imports[s.name]; ok {
		return diag.Errorf(pl.pos, diag.ErrSetNameCollision, "set %q matches an import qualifier in package %q (%s); rename the set", s.name, s.destPath, base)
	}
	return nil
}
