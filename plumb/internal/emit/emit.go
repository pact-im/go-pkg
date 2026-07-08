// Package emit is plumb’s final phase: it turns the resolved plans into the
// gofmt-canonical generated file and the -v report. It chooses import aliases,
// allocates local names, renders each injector via the plan’s instances, and
// aggregates cleanup and error handling. It produces text; it makes no wiring
// decisions and reports no diagnostics (a malformed render is an invariant
// violation and panics).
package emit

import (
	"fmt"
	"go/format"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"

	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/gotypes"
	"go.pact.im/x/plumb/internal/solve"
)

// Result is the outcome of a successful generation. It is emit’s output, handed
// back through gen to the caller.
type Result struct {
	// Source is the gofmt-canonical generated Go source.
	Source string
	// Report is the human-readable discovery-and-inference report (the body of
	// the -v output). It is informational and not contractual.
	Report string
}

// File turns the resolved plans into the gofmt-canonical generated file and the
// -v report.
func File(importPath, packageName string, pkgs []*discover.Package, plans []*solve.Plan, dest *solve.DestInfo) *Result {
	needErrors := slices.ContainsFunc(plans, planNeedsErrors)
	allLifted := liftedNamesOf(plans)

	// Pass 1: collect the packages every function references by rendering each
	// plan under a recording qualifier and discarding the text (see
	// recordPlanPackages).
	rec := newRecording(importPath)
	for _, pl := range plans {
		recordPlanPackages(pl, rec, dest, allLifted)
	}
	pkgList := rec.recordedPackages()

	// Assign import aliases and (if needed) the errors qualifier. Imports live in
	// file scope and must avoid the generated function (set) names, which are
	// package-scope declarations in the same file: an import qualifier equal to a
	// set name would be a redeclaration.
	setNames := map[string]bool{}
	for _, pl := range plans {
		setNames[pl.Name] = true
	}
	aliasByPath, errorsAlias, importLines := assignAliases(pkgList, needErrors, importPath, allLifted, setNames, dest)

	// Pass 2: render each function.
	q := newFinal(importPath, aliasByPath)
	var funcs []string
	for _, pl := range plans {
		funcs = append(funcs, renderPlan(pl, q, errorsAlias, dest, allLifted))
	}

	src := assembleFile(packageName, importLines, funcs)
	formatted, err := format.Source([]byte(src))
	if err != nil {
		// A formatting failure means plumb produced syntactically invalid Go,
		// an internal invariant violation, never reachable from valid input.
		panic(fmt.Sprintf("plumb: generated source did not format: %v\n---\n%s", err, src))
	}

	report := buildReport(importPath, pkgs, plans, importLines, aliasByPath)
	return &Result{Source: string(formatted), Report: report}
}

// recordPlanPackages records, into the recording qualifier q, every package the
// generated function will reference. It does so by rendering the plan under q and
// discarding the text: the render is the single authoritative enumeration of
// every qualified reference (signature, lifted constraints, and each instance’s
// producing expression), so recording through it means the import block cannot
// drift from the emitted body: there is no second enumeration to keep in step.
// The errors qualifier is irrelevant while recording (whether errors is imported
// is decided separately by planNeedsErrors, and it is written as a literal alias,
// never through q), so it is left empty.
func recordPlanPackages(pl *solve.Plan, q *qualifier, dest *solve.DestInfo, lifted map[string]bool) {
	_ = renderPlan(pl, q, "", dest, lifted)
}

func liftedNamesOf(plans []*solve.Plan) map[string]bool {
	names := map[string]bool{}
	for _, pl := range plans {
		for _, tp := range pl.Lifted {
			names[tp.Obj().Name()] = true
		}
	}
	return names
}

// assignAliases assigns a deterministic qualifier to each imported package,
// aliasing on collision with destination identifiers, lifted parameter names,
// or another import.
func assignAliases(pkgs []*types.Package, needErrors bool, destPath string, lifted, setNames map[string]bool, dest *solve.DestInfo) (map[string]string, string, []string) {
	taken := reservedIdents(dest, lifted)
	for n := range setNames {
		taken[n] = true
	}

	type imp struct {
		path  string
		name  string
		alias string
	}
	var imps []imp
	seen := map[string]bool{}
	add := func(path, name string) {
		if path == destPath || seen[path] {
			return
		}
		seen[path] = true
		imps = append(imps, imp{path: path, name: name})
	}
	for _, p := range pkgs {
		add(p.Path(), p.Name())
	}
	if needErrors {
		add("errors", "errors")
	}
	slices.SortFunc(imps, func(a, b imp) int { return strings.Compare(a.path, b.path) })

	aliasByPath := map[string]string{}
	for i := range imps {
		base := imps[i].name
		if base == "" {
			base = "pkg"
		}
		name := base
		for n := 2; taken[name] || token.IsKeyword(name); n++ {
			name = base + strconv.Itoa(n)
		}
		taken[name] = true
		imps[i].alias = name
		aliasByPath[imps[i].path] = name
	}

	var lines []string
	for _, m := range imps {
		if m.alias == m.name {
			lines = append(lines, fmt.Sprintf("%q", m.path))
		} else {
			lines = append(lines, fmt.Sprintf("%s %q", m.alias, m.path))
		}
	}
	errorsAlias := ""
	if needErrors {
		errorsAlias = aliasByPath["errors"]
	}
	return aliasByPath, errorsAlias, lines
}

// renderPlan renders one set’s injector function.
func renderPlan(pl *solve.Plan, q *qualifier, errorsAlias string, dest *solve.DestInfo, lifted map[string]bool) string {
	alloc := newAllocator(dest, q, lifted)
	// localOf maps a value type (by identity) to the local/param/result name
	// holding it. local() reads it; the name is always present by the time it is
	// read, so a missing entry is an internal-invariant violation.
	var localOf gotypes.Map[types.Type, string]
	setLocal := func(t types.Type, name string) { localOf.Set(t, name) }
	local := func(t types.Type) string {
		n, ok := localOf.At(t)
		if !ok {
			panic(fmt.Sprintf("plumb: emit: no local bound for %s", gotypes.TypeName(t)))
		}
		return n
	}

	// Named results: outputs, then cleanup, then error.
	var resultDecls []string
	outNames := make([]string, len(pl.Outputs))
	for i, o := range pl.Outputs {
		n := alloc.alloc(baseName(o))
		outNames[i] = n
		resultDecls = append(resultDecls, n+" "+q.typeString(o))
	}
	cleanupName := ""
	if pl.AnyCleanup {
		cleanupName = alloc.alloc("cleanup")
		ct := "func()"
		if pl.CleanupFailable {
			ct = "func() error"
		}
		resultDecls = append(resultDecls, cleanupName+" "+ct)
	}
	errName := ""
	if pl.Fallible {
		errName = alloc.alloc("err")
		resultDecls = append(resultDecls, errName+" error")
	}

	// Parameters. Prefer the source parameter/field name as a readable base,
	// falling back to a type-derived name.
	var paramDecls []string
	for _, in := range pl.Inputs {
		base := baseName(in.Type)
		if h := in.Name; h != "" && h != "_" {
			base = lowerCamel(h)
		}
		n := alloc.alloc(base)
		setLocal(in.Type, n)
		paramDecls = append(paramDecls, n+" "+q.typeString(in.Type))
	}

	// Body.
	var body []string
	var acquired []cleanupRef
	errLocal := ""
	errDeclared := false

	for _, in := range pl.Order {
		args := make([]string, len(in.Inputs))
		for i, ref := range pl.Args[in] {
			args[i] = coerceExpr(local(ref.SrcType), ref.Coerce)
		}
		rhs := renderInstance(in, q, args)

		var lhs []string
		newVar := false
		hasErr := false
		var newCleanups []cleanupRef
		for _, r := range in.Results {
			switch r.Kind {
			case solve.ResultKindValue:
				vn := alloc.alloc(baseName(r.Typ))
				setLocal(r.Typ, vn)
				lhs = append(lhs, vn)
				newVar = true
			case solve.ResultKindCleanup:
				cn := alloc.alloc(cleanupBaseName(in.Prov.Fn))
				lhs = append(lhs, cn)
				newVar = true
				newCleanups = append(newCleanups, cleanupRef{name: cn, failable: r.Failable})
			case solve.ResultKindError:
				if errLocal == "" {
					errLocal = alloc.alloc("e")
				}
				lhs = append(lhs, errLocal)
				hasErr = true
			default:
				panic(fmt.Sprintf("plumb: unhandled result kind %d", r.Kind))
			}
		}

		switch {
		case len(in.Results) == 0:
			body = append(body, rhs)
		default:
			op := ":="
			if !newVar && errDeclared {
				op = "="
			}
			body = append(body, strings.Join(lhs, ", ")+" "+op+" "+rhs)
		}

		if hasErr {
			errDeclared = true
			stmts, errExpr := renderUnwind(acquired, errLocal, pl, errorsAlias, alloc)
			block := []string{"if " + errLocal + " != nil {"}
			block = append(block, stmts...)
			block = append(block, errName+" = "+errExpr, "return", "}")
			body = append(body, block...)
		}
		acquired = append(acquired, newCleanups...)
	}

	// Success path: assign outputs, build the aggregated cleanup, return.
	for i, o := range pl.Outputs {
		body = append(body, outNames[i]+" = "+local(o))
	}
	if pl.AnyCleanup {
		body = append(body, cleanupName+" = "+renderAggregate(acquired, pl, errorsAlias, alloc))
	}
	// A naked return is required to surface the named results; a result-less,
	// side-effect-only injector (an init/main-style set) has none, so the trailing
	// return would be dead syntax.
	if len(resultDecls) > 0 {
		body = append(body, "return")
	}

	var b strings.Builder
	b.WriteString("func ")
	b.WriteString(pl.Name)
	b.WriteString(renderLiftedHeader(pl.Lifted, q.typeString))
	b.WriteByte('(')
	b.WriteString(strings.Join(paramDecls, ", "))
	b.WriteByte(')')
	if len(resultDecls) > 0 {
		b.WriteString(" (")
		b.WriteString(strings.Join(resultDecls, ", "))
		b.WriteByte(')')
	}
	b.WriteString(" {\n")
	for _, line := range body {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteString("}")
	return b.String()
}

type cleanupRef struct {
	name     string
	failable bool
}

// renderUnwind renders the teardown of already-acquired cleanups on a mid-build
// error: newest first, with the build error joined in when cleanups are failable.
func renderUnwind(acquired []cleanupRef, errLocal string, pl *solve.Plan, errorsAlias string, alloc *allocator) (stmts []string, errExpr string) {
	rev := reverseCleanups(acquired)
	// Without any failable cleanup to unwind, the error flows through directly:
	// non-failable cleanups run as plain statements and err = e.
	if !pl.CleanupFailable || !anyFailable(rev) {
		for _, c := range rev {
			stmts = append(stmts, c.name+"()")
		}
		return stmts, errLocal
	}
	cs, errLocals := cleanupStmts(rev, alloc)
	return cs, joinErrs(errorsAlias, append([]string{errLocal}, errLocals...))
}

func anyFailable(cs []cleanupRef) bool {
	for _, c := range cs {
		if c.failable {
			return true
		}
	}
	return false
}

// renderAggregate renders the single aggregated cleanup returned on success.
func renderAggregate(acquired []cleanupRef, pl *solve.Plan, errorsAlias string, alloc *allocator) string {
	rev := reverseCleanups(acquired)
	if !pl.CleanupFailable {
		var b strings.Builder
		b.WriteString("func() {\n")
		for _, c := range rev {
			b.WriteString(c.name)
			b.WriteString("()\n")
		}
		b.WriteString("}")
		return b.String()
	}
	cs, errLocals := cleanupStmts(rev, alloc)
	var b strings.Builder
	b.WriteString("func() error {\n")
	for _, s := range cs {
		b.WriteString(s)
		b.WriteByte('\n')
	}
	b.WriteString("return ")
	b.WriteString(joinErrs(errorsAlias, errLocals))
	b.WriteString("\n}")
	return b.String()
}

// cleanupStmts renders each cleanup call in run order: a failable cleanup
// captures its error in a "cleanupErr" local, a non-failable one is called bare.
// It returns the statements and the captured error-local names, in run order.
// The error locals live in this branch’s scope alone, so they are allocated from
// a fresh per-branch allocator and reset across branches.
func cleanupStmts(rev []cleanupRef, alloc *allocator) (stmts, errLocals []string) {
	ba := alloc.branch()
	for _, c := range rev {
		if c.failable {
			errl := ba.alloc("cleanupErr")
			stmts = append(stmts, errl+" := "+c.name+"()")
			errLocals = append(errLocals, errl)
		} else {
			stmts = append(stmts, c.name+"()")
		}
	}
	return stmts, errLocals
}

// joinErrs returns the lone error expression when there is one, or an
// errors.Join over all of them when there are several.
func joinErrs(errorsAlias string, errs []string) string {
	if len(errs) == 1 {
		return errs[0]
	}
	return errorsAlias + ".Join(" + strings.Join(errs, ", ") + ")"
}

// planNeedsErrors reports whether the plan’s generated code emits errors.Join,
// and so must import the errors package. It mirrors the two join sites exactly,
// since claiming the import without an emitted Join (or vice versa) breaks the
// generated build:
//
//   - renderAggregate joins on the success path when two or more cleanups are
//     failable (a single failable cleanup is returned directly);
//   - renderUnwind joins the build error with the failable cleanups acquired so
//     far when a fallible provider fails after at least one was acquired.
func planNeedsErrors(pl *solve.Plan) bool {
	if !pl.CleanupFailable {
		return false
	}
	failable := 0
	for _, in := range pl.Order {
		for _, r := range in.Results {
			if r.Kind == solve.ResultKindCleanup && r.Failable {
				failable++
			}
		}
	}
	if failable >= 2 {
		return true
	}
	// A single failable cleanup still forces a Join if a later provider can fail
	// after it was acquired, because the unwind joins the build error with it.
	acquired := 0
	for _, in := range pl.Order {
		if acquired >= 1 && instanceFallible(in) {
			return true
		}
		for _, r := range in.Results {
			if r.Kind == solve.ResultKindCleanup && r.Failable {
				acquired++
			}
		}
	}
	return false
}

func instanceFallible(in *solve.Instance) bool {
	for _, r := range in.Results {
		if r.Kind == solve.ResultKindError {
			return true
		}
	}
	return false
}

func reverseCleanups(in []cleanupRef) []cleanupRef {
	out := slices.Clone(in)
	slices.Reverse(out)
	return out
}

// coerceExpr renders the value/pointer bridge for an argument: &x passes a value
// local where a pointer is wanted (every consumer takes the address of the same
// local, so they share one instance), *x derefs a pointer where a value is
// wanted. As a method receiver these prefixed forms are parenthesized by
// asReceiver.
func coerceExpr(local string, c solve.Coerce) string {
	switch c {
	case solve.CoerceNone:
		return local
	case solve.CoerceToPtr:
		return "&" + local
	case solve.CoerceToVal:
		return "*" + local
	default:
		panic(fmt.Sprintf("plumb: unhandled coercion %d", c))
	}
}

// renderLiftedHeader renders a generic injector’s type-parameter list
// [T constraint, ...] from its lifted free parameters, collapsing an
// empty-interface constraint to the bare "any"; qual renders a non-collapsing
// constraint. An empty list renders nothing. Both the generated header (via
// q.typeString) and the -v report (via its fallback qualifier) route through
// here so the two spellings cannot drift.
func renderLiftedHeader(lifted []*types.TypeParam, qual func(types.Type) string) string {
	if len(lifted) == 0 {
		return ""
	}
	var parts []string
	for _, tp := range lifted {
		c := "any"
		if !gotypes.ConstraintCollapsesToAny(tp.Constraint()) {
			c = qual(tp.Constraint())
		}
		parts = append(parts, tp.Obj().Name()+" "+c)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func assembleFile(pkgName string, importLines, funcs []string) string {
	var b strings.Builder
	b.WriteString("// Code generated by plumb. DO NOT EDIT.\n\n")
	b.WriteString("package ")
	b.WriteString(pkgName)
	b.WriteString("\n\n")
	if len(importLines) > 0 {
		b.WriteString("import (\n")
		for _, l := range importLines {
			b.WriteString(l)
			b.WriteByte('\n')
		}
		b.WriteString(")\n\n")
	}
	b.WriteString(strings.Join(funcs, "\n\n"))
	b.WriteString("\n")
	return b.String()
}
