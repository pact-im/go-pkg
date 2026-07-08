// Package gen is the orchestrator of plumb’s pure core: it runs the discover →
// solve → emit pipeline over loaded, type-checked packages to produce generated
// injector source, performing no I/O, holding no global state, and consulting no
// clock or randomness. Its output is a deterministic function of its inputs.
// The phases live in sibling packages (discover, solve, emit) over shared leaves
// (diag, gotypes); this package only wires them together.
package gen

import (
	"go/ast"
	"go/token"
	"go/types"
	"maps"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/emit"
	"go.pact.im/x/plumb/internal/gotypes"
	"go.pact.im/x/plumb/internal/solve"
)

// Options configures a single generation run. The cli builds it from flags; gen
// distributes the fields each phase needs.
type Options struct {
	// ImportPath is the destination package’s import path. It decides what is
	// referenced unqualified (the destination) and what is imported and
	// qualified (everything else). Required and never empty by the time the
	// core runs.
	ImportPath string
	// PackageName is the identifier emitted in the generated `package <name>`
	// clause. Required and never empty by the time the core runs.
	PackageName string
	// OutputBase is the base name of the file being written (e.g.
	// "plumb_gen.go"), used to identify the file being overwritten for the
	// same-package function-name collision check. Empty when writing to stdout,
	// in which case that check is skipped.
	OutputBase string
}

// Generate is the pure entry point: it analyzes the loaded packages, resolves
// every set, and returns the generated source and report. It performs no I/O.
// The returned error, when non-nil, is always a *diag.Error with a source
// position where one applies.
func Generate(opts Options, pkgs []*discover.Package) (*emit.Result, error) {
	// PackageName and ImportPath are required; the cli always resolves them, so an
	// empty value is an internal-caller invariant violation, not a user error.
	if opts.PackageName == "" {
		panic("plumb: gen.Options.PackageName is required")
	}
	if opts.ImportPath == "" {
		panic("plumb: gen.Options.ImportPath is required")
	}
	providers, err := discover.Analyze(pkgs, opts.ImportPath, opts.OutputBase)
	if err != nil {
		return nil, err
	}
	if len(providers) == 0 {
		return nil, diag.Errorf(token.Position{}, diag.ErrNoDirectives, "nothing to generate")
	}

	bySet := map[string][]*discover.Provider{}
	for _, p := range providers {
		bySet[p.SetName] = append(bySet[p.SetName], p)
	}

	dest, err := buildDestInfo(opts, pkgs)
	if err != nil {
		return nil, err
	}
	ctxt := types.NewContext()

	var plans []*solve.Plan
	for _, name := range slices.Sorted(maps.Keys(bySet)) {
		pl, perr := solve.Set(name, bySet[name], opts.ImportPath, opts.OutputBase, dest, ctxt)
		if perr != nil {
			return nil, perr
		}
		plans = append(plans, pl)
	}

	return emit.File(opts.ImportPath, opts.PackageName, pkgs, plans, dest), nil
}

// buildDestInfo gathers what the emitter needs to know about the destination
// package: whether it was scanned, its name, and the file each top-level
// identifier is declared in (for the collision safeguards). It also rejects a
// destination declaration that shadows a predeclared identifier, a
// whole-destination property, so it is checked here once rather than per set.
func buildDestInfo(opts Options, pkgs []*discover.Package) (*solve.DestInfo, *diag.Error) {
	di := &solve.DestInfo{PkgName: opts.PackageName, Names: map[string]string{}, Imports: map[string]string{}}
	for _, pkg := range pkgs {
		if pkg.PkgPath != opts.ImportPath || pkg.Types == nil {
			continue
		}
		di.Scanned = true
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			base := filepath.Base(diag.PosIn(pkg.Fset, obj.Pos()).Filename)
			if base == opts.OutputBase {
				// The file plumb is about to overwrite declares only generated
				// injectors, which no injector body references. Recording their names
				// as reserved would make a local or lifted parameter derived from a
				// set name (e.g. a "server" local in set "server") free on the first
				// generation but reserved once the output file exists, breaking
				// idempotent regeneration. Skip them here: gen owns output-file
				// exclusion (the collision check trusts these maps and does not
				// re-filter), and import aliases re-add the set names independently
				// (assignAliases).
				continue
			}
			// The generated file spells predeclared identifiers unqualified (error and
			// nil in fallible wiring, basic-type names in signatures) and no
			// qualification can restore a shadowed builtin, so a destination that
			// declares a universe name at top level is rejected, regardless of whether
			// any one set’s signature happens to spell it. scope.Names() is sorted, so
			// the reported name is deterministic; anchoring at the declaration points at
			// the offending source directly.
			if gotypes.IsUniverseName(name) {
				return nil, diag.Errorf(diag.PosIn(pkg.Fset, obj.Pos()), diag.ErrDestShadowsPredeclared, "destination package %q declares %s (%s), shadowing the predeclared identifier; rename the declaration", opts.ImportPath, name, base)
			}
			di.Names[name] = base
		}
		collectImportQualifiers(pkg, opts.OutputBase, di.Imports)
	}
	return di, nil
}

// collectImportQualifiers records, for each import qualifier used by a hand-written
// file in the destination package, the base name of a file that uses it (a hint for
// the collision diagnostic). The file plumb will overwrite (outputBase) is skipped
// entirely: plumb rewrites its imports, so a qualifier used only there can never
// collide with a generated set name; only one a hand-written sibling still uses
// can. This makes gen the sole owner of output-file exclusion, for import
// qualifiers as for top-level names (see buildDestInfo). Among the eligible files
// the first base in sorted order wins, so the result never depends on the loader’s
// file-iteration order.
func collectImportQualifiers(pkg *discover.Package, outputBase string, out map[string]string) {
	pathName := map[string]string{}
	for _, ip := range pkg.Types.Imports() {
		pathName[ip.Path()] = ip.Name()
	}
	files := slices.SortedFunc(slices.Values(pkg.Syntax), func(a, b *ast.File) int {
		return strings.Compare(discover.FileBase(pkg, a), discover.FileBase(pkg, b))
	})
	for _, file := range files {
		base := discover.FileBase(pkg, file)
		if base == outputBase {
			continue
		}
		for _, imp := range file.Imports {
			q := ""
			if imp.Name != nil {
				if imp.Name.Name == "_" || imp.Name.Name == "." {
					continue // blank/dot imports introduce no qualifier
				}
				q = imp.Name.Name
			} else {
				p, _ := strconv.Unquote(imp.Path.Value)
				q = pathName[p]
			}
			if q == "" {
				continue
			}
			if _, seen := out[q]; !seen {
				out[q] = base
			}
		}
	}
}
