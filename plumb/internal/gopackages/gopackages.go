// Package gopackages is the real package-loading boundary: it drives the Go
// toolchain's package loader to turn patterns into the typed packages the pure
// core consumes. It is deliberately thin and kept out of the core so the core
// stays a pure, deterministic function of already-loaded inputs.
package gopackages

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"

	"go.pact.im/x/plumb/internal/discover"
)

// Result is the outcome of loading.
type Result struct {
	Packages []*discover.Package
	// TypeErrors are the tolerated type-checking errors (e.g. a stale generated
	// file whose old wiring no longer compiles). They are collected rather than
	// fatal so regeneration still works; a caller may surface them as warnings.
	// Parse errors and other structural failures are returned as the error result
	// instead. This mirrors packagestest.Loaded.TypeErrors.
	TypeErrors []error
}

// loadMode is the smallest set of fields the core needs. NeedTypes|NeedTypesInfo
// resolve provider and signature types (cross-package types come from export
// data), which is all the core consumes, the same resolution model the
// in-memory test loader uses. We deliberately omit NeedImports|NeedDeps: they
// would type-check the entire transitive dependency graph in process (an order of
// magnitude slower on dependency-heavy packages) only to keep type errors and
// broken-dependency errors in distinct Error.Kinds. The classifier below does not
// rely on that distinction, so the dependency type-check is pure cost.
const loadMode = packages.NeedName |
	packages.NeedSyntax |
	packages.NeedTypes |
	packages.NeedTypesInfo

// Load resolves and type-checks the packages matching patterns, rooted at dir.
// A structural failure is fatal: a broken module or missing package (reported by
// the go command), or a parse error that could hide a //plumb: directive.
// Everything else, notably a type error, is tolerated and collected in
// Result.TypeErrors: regenerating over a stale generated file whose old wiring no
// longer type-checks must still work, and an unresolved import surfaces precisely
// downstream as an invalid-type diagnostic if it reaches the generated output.
func Load(patterns []string, dir string) (*Result, error) {
	if len(patterns) == 0 {
		patterns = []string{"."}
	}
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode:      loadMode,
		Dir:       dir,
		Fset:      fset,
		ParseFile: parseFile,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}

	res := &Result{}
	for _, p := range pkgs {
		for _, e := range p.Errors {
			// A parse error could hide a //plumb: directive, so the package cannot
			// be scanned safely: fatal. Every other error (a type error, including,
			// without NeedDeps, an unresolved-import "could not import") is tolerated
			// and collected; a stale generated file about to be overwritten is the
			// motivating case. An unresolved import that actually reaches the output
			// is caught by the core's invalid-type check.
			//
			// A per-package structural error (packages.ListError) is also tolerated
			// here, not treated as fatal: whole-load breakage (a broken module, no
			// package matching the patterns) surfaces either through the returned
			// err above or the empty-res.Packages check below, which is the “load as
			// a whole” failure plumb treats as fatal. A ListError attached to one of
			// several otherwise-usable packages should not abort generation.
			if e.Kind == packages.ParseError {
				return nil, fmt.Errorf("loading package %q: %w", p.PkgPath, e)
			}
			res.TypeErrors = append(res.TypeErrors, e)
		}

		// Expected to be populated for every root under NeedTypes|NeedTypesInfo:
		// go/packages assigns Types and TypesInfo unconditionally, and the sole
		// nil-TypesInfo path (sources missing) carries a ParseError the fatal branch
		// above catches first. But this rests on an external tool's documented
		// behavior, not a plumb invariant, so an unexpected nil here is a fatal load
		// failure, not a broken internal invariant: skipping silently would hide
		// every directive in the package, so fail on the returned-error side (the
		// same side as "no packages matched" below) rather than crash on user input.
		if p.Types == nil || p.TypesInfo == nil {
			return nil, fmt.Errorf("package %q loaded without type information", p.PkgPath)
		}

		res.Packages = append(res.Packages, &discover.Package{
			PkgPath: p.PkgPath,
			Name:    p.Name,
			Fset:    fset,
			Syntax:  p.Syntax,
			Types:   p.Types,
			Info:    p.TypesInfo,
		})
	}
	if len(res.Packages) == 0 {
		return nil, fmt.Errorf("no packages matched %v", patterns)
	}
	return res, nil
}

// parseFile parses a source file for go/packages. SkipObjectResolution drops the
// unused ast.Object graph; ParseComments keeps the //plumb: directives.
func parseFile(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	const mode = parser.ParseComments | parser.SkipObjectResolution
	return parser.ParseFile(fset, filename, src, mode)
}
