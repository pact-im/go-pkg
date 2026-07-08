// Package cli is the command-line boundary for plumb: it parses flags, drives
// the loader, resolves the destination, runs the pure core, and writes output.
// All diagnostics flow through the streams passed to Run so the behavior is
// testable.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"go/token"
	"io"
	"path/filepath"

	"golang.org/x/mod/module"

	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gen"
	"go.pact.im/x/plumb/internal/gopackages"
	"go.pact.im/x/plumb/internal/writefile"
)

// Sentinel errors for destination resolution, so its failures are classified in
// tests with errors.Is rather than by matching message text.
var (
	errNoPackages           = errors.New("no packages scanned")
	errAmbiguousDestination = errors.New("ambiguous destination")
	errSyntheticImportPath  = errors.New("no real import path")
	errInvalidImportPath    = errors.New("invalid import path")
	errPackageNameRequired  = errors.New("package name required")
	errInvalidPackageName   = errors.New("invalid package name")
	errPackageNameMismatch  = errors.New("package name does not match destination")
)

// syntheticImportPath is the placeholder import path the go/packages loader
// assigns to a package it synthesizes from file arguments (e.g. `plumb a.go`)
// rather than resolving from an import-path pattern. It is syntactically a valid
// path, so it passes module.CheckImportPath; plumb must reject it separately
// because it is not a real, importable destination to qualify against.
const syntheticImportPath = "command-line-arguments"

const usage = `plumb is a demand-based compile-time dependency injection generator.

Usage:
  plumb [-package-name=<name>] [-import-path=<path>] [-output=<file>] [-v] [packages...]

plumb scans the matched packages for //plumb:<name> directives and writes one
generated function per set. With no package pattern it scans the current
directory.
`

// Run executes plumb with the given arguments and streams, returning the process
// exit code. It never calls os.Exit itself. Output is split across three writers:
// stdout carries the generated source (when -output is unset), stderr carries
// diagnostics and the tolerated-error notes, and report carries the -v
// discovery-and-inference report. That report is the deterministic,
// plumb-controlled stream, kept separate from the variable diagnostics so it can
// be consumed on its own.
// The command wires report to the same destination as stderr.
func Run(args []string, stdout, stderr, report io.Writer) int {
	fs := flag.NewFlagSet("plumb", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { _, _ = fmt.Fprint(stderr, usage) }

	packageName := fs.String("package-name", "", "package name for the generated `package` clause")
	importPath := fs.String("import-path", "", "import `path` of the destination package")
	output := fs.String("output", "", "output `file` (default: standard output)")
	verbose := fs.Bool("v", false, "print a discovery-and-inference report to standard error")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0 // explicitly requested help is a success, not a usage error
		}
		return 2
	}

	res, err := gopackages.Load(fs.Args(), "")
	if err != nil {
		fail(stderr, err.Error())
		return 1
	}

	destPath, pkgName, err := resolveDestination(res.Packages, *importPath, *packageName)
	if err != nil {
		fail(stderr, err.Error())
		return 1
	}

	opts := gen.Options{
		ImportPath:  destPath,
		PackageName: pkgName,
		OutputBase:  baseName(*output),
	}

	result, gerr := gen.Generate(opts, res.Packages)
	if gerr != nil {
		fail(stderr, gerr.Error())
		// Surface any type errors tolerated at load time, so a genuine provider
		// type problem (not just a stale generated file) is not swallowed when
		// generation then fails for a related reason.
		noteToleratedErrors(stderr, res.TypeErrors)
		return 1
	}

	if *verbose {
		// On success the load errors are tolerated, but under -v surface them so a
		// not-fully-type-correct input (a malformed body, a type used only
		// internally) is visible rather than silently swallowed.
		noteToleratedErrors(stderr, res.TypeErrors)
		_, _ = fmt.Fprint(report, result.Report)
	}

	if err := writeOutput(*output, result.Source, stdout); err != nil {
		fail(stderr, err.Error())
		return 1
	}
	return 0
}

// resolveDestination determines the destination import path and emitted package
// name from the flags and scanned packages, inferring each when its flag is
// omitted and rejecting the ambiguous or contradictory combinations.
func resolveDestination(pkgs []*discover.Package, importPathFlag, packageNameFlag string) (importPath, packageName string, err error) {
	importPath = importPathFlag
	if importPath == "" {
		switch len(pkgs) {
		case 1:
			importPath = pkgs[0].PkgPath
		case 0:
			return "", "", fmt.Errorf("%w: cannot infer -import-path", errNoPackages)
		default:
			return "", "", fmt.Errorf("%w: more than one package scanned; -import-path is required", errAmbiguousDestination)
		}
	}
	if importPath == syntheticImportPath {
		// The synthetic path can arrive two ways: inferred from a file-argument load
		// (the user gave no -import-path), or passed verbatim. Say which, so the
		// advice fits: "pass -import-path" is wrong when they already did.
		if importPathFlag == "" {
			return "", "", fmt.Errorf("%w: could not infer one from the arguments; pass -import-path", errSyntheticImportPath)
		}
		return "", "", fmt.Errorf("%w: %q is not a real, importable package", errSyntheticImportPath, syntheticImportPath)
	}
	if err := module.CheckImportPath(importPath); err != nil {
		return "", "", fmt.Errorf("%w: %q: %v", errInvalidImportPath, importPath, err)
	}

	// Find the destination among scanned packages (same-package generation).
	var scanned *discover.Package
	for _, p := range pkgs {
		if p.PkgPath == importPath {
			scanned = p
			break
		}
	}

	packageName = packageNameFlag
	if packageName == "" {
		if scanned == nil {
			return "", "", fmt.Errorf("%w: -package-name is required when generating into a package that is not scanned (%q)", errPackageNameRequired, importPath)
		}
		packageName = scanned.Name
	}
	if !token.IsIdentifier(packageName) || packageName == "_" {
		return "", "", fmt.Errorf("%w: %q must be a valid Go identifier", errInvalidPackageName, packageName)
	}
	// Same-package mode: the destination’s real name is known, and the generated
	// file must share it. Honoring a conflicting flag would emit an uncompilable
	// package clause, so reject the mismatch rather than emit broken code. (Checked
	// after validity so a malformed value is still reported as invalid.)
	if packageNameFlag != "" && scanned != nil && packageNameFlag != scanned.Name {
		return "", "", fmt.Errorf("%w: -package-name %q contradicts the scanned destination package %q; omit it in same-package mode", errPackageNameMismatch, packageNameFlag, scanned.Name)
	}
	return importPath, packageName, nil
}

func writeOutput(path, src string, stdout io.Writer) error {
	if path == "" {
		_, err := io.WriteString(stdout, src)
		return err
	}
	// Write atomically so an interrupted run cannot leave a truncated file where
	// the previous good output was. plumb does not create intermediate directories;
	// writefile.Write preserves that (its temp-file create fails on a missing parent).
	return writefile.Write(path, src)
}

func baseName(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

func fail(stderr io.Writer, msg string) {
	_, _ = fmt.Fprintf(stderr, "plumb: %s\n", msg)
}

// noteHeader is the fixed line noteToleratedErrors prints above the tolerated
// go/types errors. It is deterministic (unlike the error bodies below it), so
// tests key on it to confirm the note surfaced.
const noteHeader = "plumb: note: tolerated load errors:"

// noteToleratedErrors prints the load-time type errors plumb tolerated, under a
// single stderr note header with one indented line per error. They are surfaced
// on failure (a tolerated error may be the real cause) and under -v on success
// (so a not-fully-type-correct load is visible).
//
// The errors print in the order the loader encountered them, which is not stable
// across package load orders. That is deliberate: unlike the generated source
// and the -v report (which plumb guarantees byte-identical across read orders),
// these notes are a diagnostic aid, and showing the errors in the order they
// actually occurred is more useful than imposing an artificial, stable sort.
func noteToleratedErrors(stderr io.Writer, typeErrors []error) {
	if len(typeErrors) == 0 {
		return
	}
	_, _ = fmt.Fprintln(stderr, noteHeader)
	for _, te := range typeErrors {
		_, _ = fmt.Fprintf(stderr, "  %s\n", te)
	}
}
