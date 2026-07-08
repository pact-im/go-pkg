// Package packagestest is an in-memory Go package loader for tests. It type-checks
// source provided as a map of path → content, resolving imports against the
// in-memory module and a small synthetic subset of the standard library (see
// fakeStd), never GOROOT, a subprocess, or the disk. That keeps the fast unit
// tests (table-driven, golden, compile-check, fuzz, determinism) hermetic: they
// pass even when the test binary is built with -trimpath, which strips the
// baked-in GOROOT that the default disk-backed importer would otherwise need.
//
// File paths encode the package: a file at "example.com/app/x.go" belongs to
// the package whose import path is "example.com/app". The package name is taken
// from the package clause.
package packagestest

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"slices"
	"strings"

	"go.pact.im/x/plumb/internal/discover"
)

// fakeStd is the standard-library subset the corpus and plumb’s generated output
// import, given as synthetic source. Resolving stdlib here rather than through
// go/importer keeps loads self-contained (no GOROOT, no disk), so the tests
// type-check under -trimpath. The declared surface is exactly what the fixtures
// and generated code reference: a wider import fails to resolve like any missing
// package, which points at the fixture that needs a new entry added here.
var fakeStd = map[string]string{
	"database/sql": "package sql\n\ntype DB struct{}\n",
	"errors":       "package errors\n\nfunc Join(errs ...error) error { return nil }\n",
	"fmt":          "package fmt\n\nfunc Sprint(a ...any) string { return \"\" }\n",
	"io":           "package io\n\ntype Reader interface {\n\tRead(p []byte) (n int, err error)\n}\n",
	"log":          "package log\n\ntype Logger struct{}\n",
	"net/url":      "package url\n\ntype URL struct{}\n",
	// os.File carries a Read method so provider_kinds’ assertion that *os.File
	// satisfies io.Reader type-checks as it does against the real standard library.
	"os": "package os\n\n" +
		"type File struct{}\n\n" +
		"func (f *File) Read(p []byte) (n int, err error) { return 0, nil }\n\n" +
		"var Stdin *File\n",
	"strings": "package strings\n\ntype Replacer struct{}\n",
}

// ErrImportCycle is returned by Load when type-checking hits an import cycle. It
// is a sentinel so callers classify the failure with errors.Is rather than
// matching message text. (The type-checker also reports the cycle as a per-import
// type error, but that loses the error chain.)
var ErrImportCycle = errors.New("import cycle")

// Loaded is the result of loading an in-memory module.
type Loaded struct {
	Packages []*discover.Package
	// TypeErrors are non-fatal type-checking errors (e.g. a stale generated
	// file). Parse errors are returned as the error result instead.
	TypeErrors []error
}

// ShuffleFunc reorders a sequence of n elements by swapping pairs. It has the
// signature of math/rand/v2’s (*Rand).Shuffle, so a *Rand’s Shuffle method can
// be passed directly.
type ShuffleFunc func(n int, swap func(i, j int))

// Load parses and type-checks the given files, grouped into packages by their
// directory. The companion LoadShuffled additionally reorders the files, the
// load order, and the returned package slice to exercise determinism.
func Load(files map[string]string) (*Loaded, error) {
	return LoadShuffled(files, nil)
}

// LoadShuffled is Load with an optional reordering of files within packages and
// of the packages themselves, used by the determinism test. When shuffle is nil
// the deterministic sorted order is used as-is.
func LoadShuffled(files map[string]string, shuffle ShuffleFunc) (*Loaded, error) {
	fset := token.NewFileSet()

	// Group file paths by package import path (their directory).
	byPkg := map[string][]string{}
	for p := range files {
		dir := path.Dir(p)
		byPkg[dir] = append(byPkg[dir], p)
	}

	pkgPaths := make([]string, 0, len(byPkg))
	for p := range byPkg {
		pkgPaths = append(pkgPaths, p)
	}

	slices.Sort(pkgPaths)
	applyShuffle(pkgPaths, shuffle)

	for _, ip := range byPkg {
		slices.Sort(ip)
		applyShuffle(ip, shuffle)
	}

	// Parse every file.
	parsed := map[string][]*ast.File{} // import path → files
	names := map[string]string{}
	for _, ip := range pkgPaths {
		for _, fp := range byPkg[ip] {
			const mode = parser.ParseComments | parser.SkipObjectResolution
			f, err := parser.ParseFile(fset, fp, files[fp], mode)
			if err != nil {
				return nil, fmt.Errorf("parse %s: %w", fp, err)
			}
			parsed[ip] = append(parsed[ip], f)
			names[ip] = f.Name.Name
		}
	}

	// Seed the synthetic standard library. These packages resolve on demand
	// through recChecker.Import, exactly like an intra-module import, but are not
	// part of the loaded module: they are excluded from pkgPaths, so they are
	// neither type-checked eagerly nor returned in l.Packages. Seeded after the
	// module files so module source positions are unaffected by this map’s order.
	for ip, src := range fakeStd {
		f, err := parser.ParseFile(fset, ip+"/std.go", src, parser.SkipObjectResolution)
		if err != nil {
			panic(fmt.Sprintf("packagestest: fakeStd[%q] does not parse: %v", ip, err))
		}
		parsed[ip] = []*ast.File{f}
		names[ip] = f.Name.Name
	}

	l := &Loaded{}
	checker := &recChecker{
		fset:    fset,
		parsed:  parsed,
		names:   names,
		cache:   map[string]*types.Package{},
		infos:   map[string]*types.Info{},
		typeErr: &l.TypeErrors,
	}

	for _, ip := range pkgPaths {
		tpkg, info, err := checker.check(ip)
		if err != nil {
			return nil, err
		}
		l.Packages = append(l.Packages, &discover.Package{
			PkgPath: ip,
			Name:    names[ip],
			Fset:    fset,
			Syntax:  parsed[ip],
			Types:   tpkg,
			Info:    info,
		})
	}

	// An import cycle is a structural load failure (the type-checker also records
	// it as a per-import type error, but that string loses the chain). Surface it
	// as the sentinel-wrapped error so callers classify with errors.Is.
	if checker.cycleErr != nil {
		return nil, checker.cycleErr
	}

	slices.SortFunc(l.Packages, func(a, b *discover.Package) int {
		return strings.Compare(a.PkgPath, b.PkgPath)
	})
	applyShuffle(l.Packages, shuffle)

	return l, nil
}

// recChecker type-checks the in-memory packages, resolving imports (both
// intra-module packages and the synthetic standard library) recursively from
// its parsed set. An import outside that set is an unknown package.
type recChecker struct {
	fset     *token.FileSet
	parsed   map[string][]*ast.File
	names    map[string]string
	cache    map[string]*types.Package
	infos    map[string]*types.Info
	inProg   map[string]bool
	typeErr  *[]error
	cycleErr error // first import cycle seen, wrapping ErrImportCycle
}

func (c *recChecker) Import(path string) (*types.Package, error) {
	if _, ok := c.parsed[path]; ok {
		tpkg, _, err := c.check(path)
		return tpkg, err
	}
	return nil, fmt.Errorf("cannot find package %q", path)
}

func (c *recChecker) appendError(err error) {
	*c.typeErr = append(*c.typeErr, err)
}

func (c *recChecker) check(ip string) (*types.Package, *types.Info, error) {
	if p, ok := c.cache[ip]; ok {
		return p, c.infos[ip], nil
	}
	if c.inProg == nil {
		c.inProg = map[string]bool{}
	}
	if c.inProg[ip] {
		if c.cycleErr == nil {
			c.cycleErr = fmt.Errorf("import cycle through %s: %w", ip, ErrImportCycle)
		}
		return nil, nil, c.cycleErr
	}
	c.inProg[ip] = true
	defer delete(c.inProg, ip)

	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
		Instances:  map[*ast.Ident]types.Instance{},
		Implicits:  map[ast.Node]types.Object{},
		Scopes:     map[ast.Node]*types.Scope{},
	}
	conf := types.Config{
		Importer: c,
		Error:    c.appendError,
	}
	tpkg, _ := conf.Check(ip, c.fset, c.parsed[ip], info)
	c.cache[ip] = tpkg
	c.infos[ip] = info
	return tpkg, info, nil
}

// applyShuffle reorders s in place using the given shuffle function. A nil
// shuffle leaves s untouched.
func applyShuffle[S ~[]E, E any](s S, shuffle ShuffleFunc) {
	if shuffle != nil {
		shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	}
}
