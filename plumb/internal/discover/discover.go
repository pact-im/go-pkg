// Package discover is plumb’s first phase: it scans loaded packages for
// //plumb:<name> directives and turns each tagged declaration into a Provider,
// the surface-independent description the solve phase wires together. It reads
// syntax and type information but makes no wiring decisions.
package discover

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strings"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/gotypes"
)

// directiveNames extracts the set names from a comment group attached to a
// declaration. It returns the recognized plumb set names in order, plus any
// located error for a recognized-but-invalid name (e.g. //plumb:123).
//
// Recognition uses go/ast.ParseDirective, the same rule gofmt applies: a comment
// gofmt would demote by inserting a space after "//" (and block comments) is not
// a directive and is silently ignored. A //plumb: directive whose name is not a
// Go identifier, or that carries arguments, is rejected.
func directiveNames(doc *ast.CommentGroup, fset *token.FileSet) ([]string, *diag.Error) {
	return appendDirectiveNames(nil, doc, fset)
}

// appendDirectiveNames appends the set names from doc to prior. prior holds the
// sets the declaration already joins (an enclosing group’s directives, or an
// earlier comment in the same doc), so a directive naming one of them (whether
// the repeat is within doc or across the group/spec boundary) is a duplicate and
// is rejected identically. The returned slice reuses prior’s backing array; the
// caller must not retain prior separately.
func appendDirectiveNames(prior []string, doc *ast.CommentGroup, fset *token.FileSet) ([]string, *diag.Error) {
	if doc == nil {
		return prior, nil
	}
	names := prior
	for _, c := range doc.List {
		d, ok := ast.ParseDirective(c.Slash, c.Text)
		if !ok || d.Tool != "plumb" {
			continue // not a //plumb: directive (gofmt-demoted, block comment, or another tool)
		}
		// Check the name before the trailing-argument case: //plumb:123 extra has
		// both faults, and the invalid name is the more fundamental one to report.
		// ParseDirective already requires the name to start with [a-z0-9], so it
		// never yields the blank identifier; token.IsIdentifier is the whole check.
		if !token.IsIdentifier(d.Name) {
			return nil, diag.Errorf(diag.PosIn(fset, c.Pos()), diag.ErrInvalidSetName,
				"//%s; the set name must be a valid Go identifier", strings.TrimPrefix(c.Text, "//"))
		}
		if d.Args != "" {
			return nil, diag.Errorf(diag.PosIn(fset, c.Pos()), diag.ErrInvalidSetName,
				"//%s; a plumb directive names only a set, with no trailing arguments", strings.TrimPrefix(c.Text, "//"))
		}
		if slices.Contains(names, d.Name) {
			// A declaration joins a set at most once. Rather than silently deduping
			// (or producing a confusing "X collides with X" ambiguity), reject the
			// repeated directive with a clear, dedicated diagnostic.
			return nil, diag.Errorf(diag.PosIn(fset, c.Pos()), diag.ErrDuplicateDirective,
				"//%s; the declaration already joins set %q: remove the duplicate directive", strings.TrimPrefix(c.Text, "//"), d.Name)
		}
		names = append(names, d.Name)
	}
	return names, nil
}

// allBlank reports whether every declared name is the blank identifier, so the
// spec declares nothing that could be wired.
func allBlank(idents []*ast.Ident) bool {
	for _, id := range idents {
		if id.Name != "_" {
			return false
		}
	}
	return len(idents) > 0
}

// Analyze scans every package for directives and builds the provider list. It
// returns a located error for any malformed provider.
//
// Packages are scanned in import-path order and files in name order, and scanFile
// returns the source-earliest fault within a file, so when the input has more than
// one fault the diagnostic returned is a fixed choice: the source-earliest fault
// of the first file (in that order) to have one. This choice does not depend on the order
// the loader presented packages and files (which it does not promise); every phase
// anchors ordering to source position. The sorts run over copies, leaving the
// loaded packages untouched. Provider order is irrelevant downstream (solve
// re-sorts each set by position), so only diagnostic selection is affected.
//
// The file being overwritten (the prior generated output, identified by the
// base name of -output within the destination package) is skipped, so plumb
// never feeds its own previous output back as input. (plumb does not emit
// directives, so this only matters for a hand-edited or stray directive in that
// file.)
func Analyze(pkgs []*Package, destPath, outputBase string) ([]*Provider, *diag.Error) {
	var providers []*Provider
	for _, pkg := range slices.SortedFunc(slices.Values(pkgs), func(a, b *Package) int {
		return cmp.Compare(a.PkgPath, b.PkgPath)
	}) {
		for _, file := range slices.SortedFunc(slices.Values(pkg.Syntax), func(a, b *ast.File) int {
			return cmp.Compare(FileBase(pkg, a), FileBase(pkg, b))
		}) {
			if outputBase != "" && pkg.PkgPath == destPath && FileBase(pkg, file) == outputBase {
				continue // the file plumb is about to overwrite
			}
			ps, err := scanFile(pkg, file)
			if err != nil {
				return nil, err
			}
			providers = append(providers, ps...)
		}
	}
	return providers, nil
}

// FileBase returns the base name of the file containing the given syntax tree.
func FileBase(pkg *Package, file *ast.File) string {
	tf := pkg.Fset.File(file.Pos())
	if tf == nil {
		return ""
	}
	return filepath.Base(tf.Name())
}

// scanFile walks one file’s top-level declarations (locals are never providers)
// and the type members nested in them, returning the source-earliest fault it
// finds. A stray directive is caught only by the whole-file sweep, which must run
// after the declaration scan, so the two can surface out of source order: a stray
// in an early function body precedes a duplicate directive on a later declaration,
// yet the declaration scan produces its fault first. Comparing the two by position
// makes the reported diagnostic the one a reader would fix first. Within each phase
// the earliest already comes first (declarations are visited in source order, so
// the first declaration fault is the earliest of them, and the sweep walks comments
// in position order), so only the cross-phase pair needs comparing.
func scanFile(pkg *Package, file *ast.File) ([]*Provider, *diag.Error) {
	var out []*Provider
	var declFault *diag.Error
	for _, decl := range file.Decls {
		var ps []*Provider
		var err *diag.Error
		switch d := decl.(type) {
		case *ast.FuncDecl:
			ps, err = scanFunc(pkg, d)
		case *ast.GenDecl:
			ps, err = scanGenDecl(pkg, d)
		default:
			continue
		}
		if err != nil {
			declFault = err
			break
		}
		out = append(out, ps...)
	}
	if fault := diag.Earlier(declFault, reportStrayDirectives(pkg, file)); fault != nil {
		return nil, fault
	}
	return out, nil
}

// reportStrayDirectives flags any //plumb: directive that sits in a position the
// targeted scan above never reads, so it would otherwise be silently ignored: a
// directive in a function body, on an import, or floating free between
// declarations. The scan only consults the doc comments of declarations and
// members it recognizes, so a directive anywhere else is a mistake worth a
// located error rather than a no-op.
//
// The sweep walks file.Comments (every comment group in the file) because
// comments inside bodies and free-floating ones live only there and are never
// attached to a node’s Doc; an AST node walk would miss exactly the positions we
// want to catch. collectProviderDocs marks the groups the scan does read, and
// anything else carrying a plumb directive is reported.
func reportStrayDirectives(pkg *Package, file *ast.File) *diag.Error {
	read := collectProviderDocs(file)
	for _, grp := range file.Comments {
		if read[grp] {
			continue
		}
		for _, c := range grp.List {
			d, ok := ast.ParseDirective(c.Slash, c.Text)
			if !ok || d.Tool != "plumb" {
				continue
			}
			return diag.Errorf(diag.PosIn(pkg.Fset, c.Pos()), diag.ErrMisplacedDirective,
				"//%s; a directive attaches to a package-level function, method, variable, constant, conversion, struct field, struct type, or interface method",
				strings.TrimPrefix(c.Text, "//"))
		}
	}
	return nil
}

// collectProviderDocs returns the set of comment groups the targeted scan reads
// for directives: the doc of each top-level func, of each non-import GenDecl and
// its value/type specs, and of each struct field and interface method. A plumb
// directive in any other group is in an unsupported position.
func collectProviderDocs(file *ast.File) map[*ast.CommentGroup]bool {
	read := make(map[*ast.CommentGroup]bool)
	mark := func(cg *ast.CommentGroup) {
		if cg != nil {
			read[cg] = true
		}
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			mark(d.Doc)
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				continue // import specs are never providers; directives on them are stray
			}
			mark(d.Doc)
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					mark(s.Doc)
				case *ast.TypeSpec:
					mark(s.Doc)
					switch t := s.Type.(type) {
					case *ast.StructType:
						for _, f := range t.Fields.List {
							mark(f.Doc)
						}
					case *ast.InterfaceType:
						for _, m := range t.Methods.List {
							mark(m.Doc)
						}
					}
				}
			}
		}
	}
	return read
}

// specDirectiveNames returns the sets a spec joins: the union of the directives on
// the enclosing GenDecl (which apply to every spec in the group) and the spec’s
// own directives. A group directive that precedes the keyword attaches to
// GenDecl.Doc for both a single declaration and a parenthesized group; a directive
// inside the group attaches to the spec’s own doc. Unioning them (rather than
// letting the spec’s doc shadow the group’s) means a spec that joins an extra set
// does not silently drop out of the group’s set. A set named at both levels is a
// duplicate (the group already joins the spec to it), so it is rejected like any
// other repeated directive rather than silently deduped.
func specDirectiveNames(gd *ast.GenDecl, specDoc *ast.CommentGroup, fset *token.FileSet) ([]string, *diag.Error) {
	names, err := directiveNames(gd.Doc, fset)
	if err != nil {
		return nil, err
	}
	return appendDirectiveNames(names, specDoc, fset)
}

func scanFunc(pkg *Package, fd *ast.FuncDecl) ([]*Provider, *diag.Error) {
	names, err := directiveNames(fd.Doc, pkg.Fset)
	if err != nil || len(names) == 0 {
		return nil, err
	}
	obj, _ := pkg.Info.Defs[fd.Name].(*types.Func)
	if obj == nil {
		// A declared name with no object did not type-check: a redeclaration
		// loses its Defs entry. The loader already reports the real cause; a
		// recognized directive must still fail loudly, never silently vanish.
		return nil, diag.Errorf(diag.PosIn(pkg.Fset, fd.Name.Pos()), diag.ErrInvalidType, "declaration of %s did not type-check; cannot wire it", fd.Name.Name)
	}
	sig := obj.Type().(*types.Signature)
	pos := diag.PosIn(pkg.Fset, fd.Name.Pos())
	// go/types records an object for a blank name too, so a directive on func _()
	// or a blank method reaches here. A blank identifier cannot be referenced (and
	// a blank method is absent from its type’s method set, so it would later panic),
	// so reject it rather than emit an unusable "_" call.
	if obj.Name() == "_" {
		what := "function"
		if sig.Recv() != nil {
			what = "method"
		}
		return nil, diag.Errorf(pos, diag.ErrBlankProvider, "give the %s a name; a blank identifier declares nothing to wire", what)
	}
	// init is the same trap with a different spelling: Go forbids referring to it
	// from anywhere, so an emitted init() call can never compile, and a package
	// may declare several. A method named init is an ordinary selector and stays
	// allowed. (A set named init is unrelated and supported.)
	if obj.Name() == "init" && sig.Recv() == nil {
		return nil, diag.Errorf(pos, diag.ErrInitProvider, "init cannot be referenced, so it cannot provide; rename the function")
	}
	var out []*Provider
	for _, name := range names {
		p := &Provider{
			SetName: name,
			Pos:     pos,
			Name:    obj.Name(),
			Pkg:     pkg.Types,
			Fn:      obj,
		}
		if recv := sig.Recv(); recv != nil {
			// A method. Resolve the receiver’s named type and its type params. A
			// concrete and an interface receiver share one kind; solve resolves the
			// receiver input type from the owner at instantiation time.
			named := receiverNamed(recv.Type())
			if named == nil {
				// The receiver type did not type-check: a typo’d or undefined
				// receiver in tolerated-invalid input. The loader already reports the
				// real cause, so surface a located error rather than crashing.
				return nil, diag.Errorf(pos, diag.ErrInvalidType, "method %s has an unresolvable receiver type", obj.Name())
			}
			p.Kind = KindMethod
			p.Owner = named
			p.Tparams = named.TypeParams()
			// Diagnostics and the report name a method as Owner.Method, matching
			// the interface-method path: two same-named methods on different
			// receivers stay distinguishable without leaning on positions.
			p.Name = named.Obj().Name() + "." + obj.Name()
		} else {
			p.Kind = KindFunc
			p.Tparams = sig.TypeParams()
		}
		out = append(out, p)
	}
	return out, nil
}

// receiverNamed unwraps a method receiver type to its origin named type, or nil
// if the receiver is not a defined (named) type.
func receiverNamed(t types.Type) *types.Named {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	n, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return nil
	}
	return n.Origin()
}

func scanGenDecl(pkg *Package, gd *ast.GenDecl) ([]*Provider, *diag.Error) {
	// An empty group such as var (), const (), or type () has no spec a directive could
	// attach to, and collectProviderDocs marks the group doc as read, so the
	// stray-directive sweep never sees it. Judge the directive here: an invalid
	// name surfaces as such, and a well-formed one is rejected as misplaced
	// rather than silently ignored. Import groups stay with the sweep, which
	// already flags them.
	if len(gd.Specs) == 0 && gd.Tok != token.IMPORT {
		names, err := directiveNames(gd.Doc, pkg.Fset)
		if err != nil {
			return nil, err
		}
		if len(names) > 0 {
			for _, c := range gd.Doc.List {
				if d, ok := ast.ParseDirective(c.Slash, c.Text); ok && d.Tool == "plumb" {
					return nil, diag.Errorf(diag.PosIn(pkg.Fset, c.Pos()), diag.ErrMisplacedDirective,
						"//%s; the declaration group is empty and declares nothing to wire: add a declaration or remove the directive",
						strings.TrimPrefix(c.Text, "//"))
				}
			}
		}
	}
	var out []*Provider
	for _, spec := range gd.Specs {
		switch s := spec.(type) {
		case *ast.ValueSpec:
			ps, err := scanValueSpec(pkg, gd, s)
			if err != nil {
				return nil, err
			}
			out = append(out, ps...)
		case *ast.TypeSpec:
			ps, err := scanTypeSpec(pkg, gd, s)
			if err != nil {
				return nil, err
			}
			out = append(out, ps...)
		}
	}
	return out, nil
}

func scanValueSpec(pkg *Package, gd *ast.GenDecl, s *ast.ValueSpec) ([]*Provider, *diag.Error) {
	names, err := specDirectiveNames(gd, s.Doc, pkg.Fset)
	if err != nil || len(names) == 0 {
		return nil, err
	}
	switch gd.Tok {
	case token.VAR:
		return scanVar(pkg, s, names)
	case token.CONST:
		return scanConst(pkg, s, names)
	}
	return nil, nil
}

func scanVar(pkg *Package, s *ast.ValueSpec, names []string) ([]*Provider, *diag.Error) {
	if len(s.Names) == 0 {
		return nil, nil
	}
	// Conversion provider: a single blank variable with an explicit target type.
	if len(s.Names) == 1 && s.Names[0].Name == "_" {
		return scanConversion(pkg, s, s.Names[0], names)
	}
	// Past the conversion case, an all-blank spec (var _, _ = ...) declares nothing
	// to wire; the directive is inert, so reject it rather than silently dropping it.
	if allBlank(s.Names) {
		return nil, diag.Errorf(diag.PosIn(pkg.Fset, s.Names[0].Pos()), diag.ErrBlankProvider, "give the variable a name; a blank identifier declares nothing to wire")
	}

	// The directive applies to the spec; every named variable it declares becomes
	// a value provider. (Two of the same type then collide as an ambiguity, which
	// is correct.)
	var out []*Provider
	for _, ident := range s.Names {
		if ident.Name == "_" {
			continue
		}
		obj, _ := pkg.Info.Defs[ident].(*types.Var)
		if obj == nil {
			// Redeclared names lose their Defs entry; fail loudly, as scanFunc does.
			return nil, diag.Errorf(diag.PosIn(pkg.Fset, ident.Pos()), diag.ErrInvalidType, "declaration of %s did not type-check; cannot wire it", ident.Name)
		}
		pos := diag.PosIn(pkg.Fset, ident.Pos())
		for _, name := range names {
			out = append(out, &Provider{
				SetName: name,
				Kind:    KindSymbol,
				Pos:     pos,
				Name:    obj.Name(),
				Pkg:     pkg.Types,
				Sym:     obj,
			})
		}
	}
	return out, nil
}

func scanConversion(pkg *Package, s *ast.ValueSpec, ident *ast.Ident, names []string) ([]*Provider, *diag.Error) {
	pos := diag.PosIn(pkg.Fset, ident.Pos())
	if s.Type == nil {
		return nil, diag.Errorf(pos, diag.ErrInvalidConversion, "must declare an explicit target type")
	}
	targetType := pkg.Info.TypeOf(s.Type)
	if targetType == nil {
		// Unreachable even for tolerated-invalid input: go/types records a
		// (possibly Invalid) type for every type expression, never nothing.
		panic(fmt.Sprintf("plumb: conversion target type at %s has no recorded type", pos))
	}
	if len(s.Values) == 0 {
		return nil, diag.Errorf(pos, diag.ErrInvalidConversion, "needs a source like (*T)(nil); write the source type explicitly")
	}
	// A conversion provider goes from the source type to the declared target type:
	// it consumes typeof(Expr) and produces T, emitting T(src). Any conversion the
	// blank-var assignment type-checks is allowed (Go makes assignability itself a
	// conversion rule), so the target may be an interface the source implements, a
	// directional channel, a pointer with an identical base, and so on. The only
	// rejected form is an untyped nil, which names no source type at all.
	srcType := pkg.Info.TypeOf(s.Values[0])
	if srcType == nil {
		// An undefined identifier in value position gets no recorded type at all
		// (unlike a type position, which records Invalid). The loader already
		// reports the real cause; the directive must still fail loudly.
		return nil, diag.Errorf(pos, diag.ErrInvalidType, "conversion source did not type-check; cannot wire it")
	}
	if gotypes.IsUntypedNil(srcType) {
		return nil, diag.Errorf(pos, diag.ErrInvalidConversion, "bare nil source; write (*Concrete)(nil) or name the source type")
	}
	// An identity conversion (source and target the same type) consumes and
	// produces the same type, which would otherwise surface as a self-referential
	// dependency cycle; reject it with the real reason instead.
	if types.Identical(srcType, targetType) {
		// A constant source that is not itself written as a conversion (var _
		// MyInt = 5) already carries the target’s type (go/types converts the
		// literal implicitly), so "both MyInt" would baffle someone who wrote 5.
		// Name the real problem instead.
		if tv := pkg.Info.Types[s.Values[0]]; tv.Value != nil && !isConversionExpr(pkg, s.Values[0]) {
			return nil, diag.Errorf(pos, diag.ErrInvalidConversion, "the source must name a type, e.g. (*T)(nil); a constant takes the target's own type")
		}
		return nil, diag.Errorf(pos, diag.ErrInvalidConversion, "identity conversion: source and target are both %s", gotypes.TypeName(targetType))
	}
	// On valid input the blank-var assignment type-checks, so the source is
	// assignable (hence convertible) to the target. When it does not (tolerated
	// type-error input), a non-convertible pair would emit an uncompilable T(src);
	// report it here instead. Skip the check if either side did not type-check; that
	// invalid type is caught downstream as ErrInvalidType.
	if !gotypes.ContainsInvalid(srcType) && !gotypes.ContainsInvalid(targetType) && !types.ConvertibleTo(srcType, targetType) {
		return nil, diag.Errorf(pos, diag.ErrInvalidConversion, "cannot convert %s to %s", gotypes.TypeName(srcType), gotypes.TypeName(targetType))
	}
	var out []*Provider
	for _, name := range names {
		out = append(out, &Provider{
			SetName:     name,
			Kind:        KindConvert,
			Pos:         pos,
			Name:        gotypes.TypeName(targetType),
			Pkg:         pkg.Types,
			ConvertTo:   targetType,
			ConvertFrom: srcType,
		})
	}
	return out, nil
}

func scanConst(pkg *Package, s *ast.ValueSpec, names []string) ([]*Provider, *diag.Error) {
	if allBlank(s.Names) {
		return nil, diag.Errorf(diag.PosIn(pkg.Fset, s.Names[0].Pos()), diag.ErrBlankProvider, "give the constant a name; a blank identifier declares nothing to wire")
	}
	var out []*Provider
	for _, ident := range s.Names {
		if ident.Name == "_" {
			continue
		}
		obj, _ := pkg.Info.Defs[ident].(*types.Const)
		if obj == nil {
			// Redeclared names lose their Defs entry; fail loudly, as scanFunc does.
			return nil, diag.Errorf(diag.PosIn(pkg.Fset, ident.Pos()), diag.ErrInvalidType, "declaration of %s did not type-check; cannot wire it", ident.Name)
		}
		pos := diag.PosIn(pkg.Fset, ident.Pos())
		if gotypes.IsUntyped(obj.Type()) {
			return nil, diag.Errorf(pos, diag.ErrUntypedConstant, "constant %s must have an explicit type", obj.Name())
		}
		for _, name := range names {
			out = append(out, &Provider{
				SetName: name,
				Kind:    KindSymbol,
				Pos:     pos,
				Name:    obj.Name(),
				Pkg:     pkg.Types,
				Sym:     obj,
			})
		}
	}
	return out, nil
}

func scanTypeSpec(pkg *Package, gd *ast.GenDecl, s *ast.TypeSpec) ([]*Provider, *diag.Error) {
	var out []*Provider

	// 1. A directive on the type declaration itself: a struct provider.
	names, err := specDirectiveNames(gd, s.Doc, pkg.Fset)
	if err != nil {
		return nil, err
	}
	tn, _ := pkg.Info.Defs[s.Name].(*types.TypeName)
	blank := s.Name.Name == "_"
	if len(names) > 0 {
		if blank {
			return nil, diag.Errorf(diag.PosIn(pkg.Fset, s.Name.Pos()), diag.ErrBlankProvider, "give the type a name; a blank identifier declares nothing to wire")
		}
		if tn == nil {
			// A declared name with no object did not type-check: a redeclaration
			// loses its Defs entry. Fail loudly rather than drop the directive.
			return nil, diag.Errorf(diag.PosIn(pkg.Fset, s.Name.Pos()), diag.ErrInvalidType, "declaration of %s did not type-check; cannot wire it", s.Name.Name)
		}
		ps, err := structProviders(pkg, s, tn, names)
		if err != nil {
			return nil, err
		}
		out = append(out, ps...)
	}

	// 2. Directives on the type’s members: struct fields or interface methods. The
	// owner is the declared type: a *types.Named for a defined type, or a
	// *types.Alias when the directive is on a member of an alias to an anonymous
	// composite (type C = struct{...} / interface{...}). This member scan only fires
	// when s.Type is a struct or interface literal, so an alias to a *named*
	// composite (which has an ast.Ident type) never reaches here. The gotypes.*
	// helpers read the origin, type parameters, and declaring object off either
	// form uniformly, so a field or method binds to a value of the alias just as it
	// does to one of a defined type.
	if tn != nil {
		owner := tn.Type()
		var members []*Provider
		switch st := s.Type.(type) {
		case *ast.StructType:
			members, err = scanStructFields(pkg, st, owner)
		case *ast.InterfaceType:
			members, err = scanInterfaceMethods(pkg, st, owner)
		}
		if err != nil {
			return nil, err
		}
		// A blank-named type cannot be referenced, so a member provider bound to it
		// would emit an unusable receiver; reject a member directive on one.
		if len(members) > 0 && blank {
			return nil, diag.Errorf(diag.PosIn(pkg.Fset, s.Name.Pos()), diag.ErrBlankProvider, "give the type a name; its members cannot be reached through a blank identifier")
		}
		out = append(out, members...)
	} else if c := memberDirective(s.Type); c != nil {
		// The owner’s declaration did not type-check (a redeclaration loses its
		// Defs entry), so the member scan cannot run. A member directive
		// must still fail loudly, never silently vanish with its owner.
		return nil, diag.Errorf(diag.PosIn(pkg.Fset, c.Pos()), diag.ErrInvalidType, "declaration of %s did not type-check; cannot wire its members", s.Name.Name)
	}
	return out, nil
}

// memberDirective returns the first //plumb: directive comment attached to a
// member of a struct or interface literal, or nil.
func memberDirective(typ ast.Expr) *ast.Comment {
	var fields []*ast.Field
	switch t := typ.(type) {
	case *ast.StructType:
		fields = t.Fields.List
	case *ast.InterfaceType:
		fields = t.Methods.List
	default:
		return nil
	}
	for _, f := range fields {
		if f.Doc == nil {
			continue
		}
		for _, c := range f.Doc.List {
			if d, ok := ast.ParseDirective(c.Slash, c.Text); ok && d.Tool == "plumb" {
				return c
			}
		}
	}
	return nil
}

func structProviders(pkg *Package, s *ast.TypeSpec, tn *types.TypeName, names []string) ([]*Provider, *diag.Error) {
	pos := diag.PosIn(pkg.Fset, s.Name.Pos())
	// The declared type may be a defined type or an alias; either way it must be a
	// struct at its core. An alias to an anonymous struct (type S = struct{...}) has
	// no *types.Named of its own but is still a struct, so gate on the underlying
	// type rather than requiring a defined type.
	under := types.Unalias(tn.Type()).Underlying()
	if _, ok := under.(*types.Struct); !ok {
		return nil, diag.Errorf(pos, diag.ErrStructProvider, "%s is %s; struct-type providers are only supported on struct types",
			tn.Name(), gotypes.KindOfType(under))
	}
	// Keep the declared type as written (a *Named, or a *Alias when the directive
	// is on an alias) so the generated composite literal names it rather than the
	// underlying type; the type parameters and origin come from that declared type.
	declared := tn.Type()
	var out []*Provider
	for _, name := range names {
		out = append(out, &Provider{
			SetName:  name,
			Kind:     KindStruct,
			Pos:      pos,
			Name:     tn.Name(),
			Pkg:      pkg.Types,
			Declared: declared,
			Tparams:  gotypes.TypeParamsOf(declared),
		})
	}
	return out, nil
}

// scanStructFields scans the fields of a struct declaration for member
// directives. owner is the declared type (a *types.Named or, for an alias to an
// anonymous struct, a *types.Alias), which the field provider binds to as its
// receiver.
func scanStructFields(pkg *Package, st *ast.StructType, owner types.Type) ([]*Provider, *diag.Error) {
	var out []*Provider
	for _, field := range st.Fields.List {
		names, err := directiveNames(field.Doc, pkg.Fset)
		if err != nil {
			return nil, err
		}
		if len(names) == 0 {
			continue
		}
		pos := diag.PosIn(pkg.Fset, field.Pos())
		if len(field.Names) == 0 {
			return nil, diag.Errorf(pos, diag.ErrEmbeddedField, "give the field a name to use it")
		}
		// An all-blank directed field group (//plumb + _ int) declares nothing to
		// wire; reject it rather than silently drop the directive, as scanVar does.
		if allBlank(field.Names) {
			return nil, diag.Errorf(pos, diag.ErrBlankProvider, "give the field a name; a blank identifier declares nothing to wire")
		}
		// Every name in a multi-name field group (A, B int) is its own provider,
		// mirroring multi-name var/const handling; same-type names then collide as
		// an ambiguity rather than being silently narrowed to the first. A blank
		// among them is skipped, the same way scanVar skips one.
		for _, ident := range field.Names {
			if ident.Name == "_" {
				continue
			}
			obj, _ := pkg.Info.Defs[ident].(*types.Var)
			if obj == nil {
				// Redeclared fields lose their Defs entry; fail loudly, as scanFunc
				// does, rather than silently drop the directive on the losing field.
				return nil, diag.Errorf(diag.PosIn(pkg.Fset, ident.Pos()), diag.ErrInvalidType, "declaration of %s did not type-check; cannot wire it", ident.Name)
			}
			for _, name := range names {
				out = append(out, &Provider{
					SetName: name,
					Kind:    KindField,
					Pos:     pos,
					Name:    gotypes.TypeNameOf(owner).Name() + "." + obj.Name(),
					Pkg:     pkg.Types,
					Sym:     obj,
					Owner:   gotypes.GenericOrigin(owner),
					Tparams: gotypes.TypeParamsOf(owner),
				})
			}
		}
	}
	return out, nil
}

// scanInterfaceMethods scans the methods of an interface declaration for member
// directives. owner is the declared type (a *types.Named or, for an alias to an
// anonymous interface, a *types.Alias), which the method provider binds to as its
// receiver.
func scanInterfaceMethods(pkg *Package, it *ast.InterfaceType, owner types.Type) ([]*Provider, *diag.Error) {
	var out []*Provider
	for _, field := range it.Methods.List {
		names, err := directiveNames(field.Doc, pkg.Fset)
		if err != nil {
			return nil, err
		}
		if len(names) == 0 {
			continue
		}
		pos := diag.PosIn(pkg.Fset, field.Pos())
		if len(field.Names) == 0 {
			return nil, diag.Errorf(pos, diag.ErrEmbeddedInterface, "put the directive on a method instead")
		}
		// A non-basic (constraint) interface cannot be the type of a value.
		if iface, ok := owner.Underlying().(*types.Interface); ok && !iface.IsMethodSet() {
			return nil, diag.Errorf(pos, diag.ErrConstraintInterfaceMethod, "%s: a constraint type cannot be the type of a value, so it cannot be a receiver", gotypes.TypeNameOf(owner).Name())
		}
		ident := field.Names[0]
		obj, _ := pkg.Info.Defs[ident].(*types.Func)
		if obj == nil {
			// No reachable trigger while the interface type-checks, but mirror the
			// sibling scanners rather than silently drop a recognized directive.
			return nil, diag.Errorf(diag.PosIn(pkg.Fset, ident.Pos()), diag.ErrInvalidType, "declaration of %s did not type-check; cannot wire it", ident.Name)
		}
		for _, name := range names {
			out = append(out, &Provider{
				SetName: name,
				Kind:    KindMethod,
				Pos:     pos,
				Name:    gotypes.TypeNameOf(owner).Name() + "." + obj.Name(),
				Pkg:     pkg.Types,
				Fn:      obj,
				Owner:   gotypes.GenericOrigin(owner),
				Tparams: gotypes.TypeParamsOf(owner),
			})
		}
	}
	return out, nil
}

// isConversionExpr reports whether e is written as a conversion (a call whose
// operand names a type), so MyInt(5) keeps the identity-conversion wording
// while a bare constant gets the constant-specific one.
func isConversionExpr(pkg *Package, e ast.Expr) bool {
	call, ok := ast.Unparen(e).(*ast.CallExpr)
	if !ok {
		return false
	}
	return pkg.Info.Types[ast.Unparen(call.Fun)].IsType()
}
