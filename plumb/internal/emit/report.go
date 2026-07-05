package emit

import (
	"fmt"
	"go/types"
	"slices"
	"strings"

	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/solve"
)

func buildReport(importPath string, pkgs []*discover.Package, plans []*solve.Plan, importLines []string, aliasByPath map[string]string) string {
	var b strings.Builder
	var paths []string
	for _, p := range pkgs {
		paths = append(paths, p.PkgPath)
	}
	slices.Sort(paths)
	fmt.Fprintf(&b, "plumb: scanned %d package(s): %s\n", len(pkgs), strings.Join(paths, ", "))
	fmt.Fprintf(&b, "plumb: generating into package %q\n", importPath)
	// Qualify foreign packages by the same import alias the generated source uses,
	// leaving the destination's own types unqualified, so the report mirrors the
	// source even when two imported packages share a name (config, config2). A
	// path missing from the alias map falls back to the package name rather than
	// panicking as the source qualifier does: the report is a diagnostic aid, not
	// contractual output.
	pathQual := func(t types.Type) string {
		return types.TypeString(t, func(p *types.Package) string {
			if p == nil || p.Path() == importPath {
				return ""
			}
			if a, ok := aliasByPath[p.Path()]; ok {
				return a
			}
			return p.Name()
		})
	}
	for _, pl := range plans {
		fmt.Fprintf(&b, "\nset %q (%s, %s):\n", pl.Name, plural(pl.ProviderCount, "provider"), plural(len(pl.Order), "instantiation"))
		for _, in := range pl.Order {
			fmt.Fprintf(&b, "  %s %s in(%s) → out(%s)\n",
				kindLabel(in.Prov), in.Prov.Name, joinTypes(in.InputTypes(), pathQual), joinResults(in, pathQual))
		}
		fmt.Fprintf(&b, "  inject %s%s in(%s) → out(%s)\n", pl.Name, renderLiftedHeader(pl.Lifted, pathQual),
			joinTypes(pl.InputTypes(), pathQual), injectorOut(pl, pathQual))
	}
	if len(importLines) > 0 {
		b.WriteString("  imports:\n")
		for _, l := range importLines {
			fmt.Fprintf(&b, "    %s\n", l)
		}
	}
	return b.String()
}

// plural renders a count with its noun, pluralizing for any count but one:
// "1 provider", "2 providers", "0 providers".
func plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

// kindLabel is the provider-kind label for the report. KindSymbol recovers
// var-vs-const from the referenced object's dynamic type, since the two share one
// kind everywhere else.
func kindLabel(p *discover.Provider) string {
	if p.Kind == discover.KindSymbol {
		if _, ok := p.Sym.(*types.Const); ok {
			return "const"
		}
		return "var"
	}
	return p.Kind.String()
}

func joinTypes(ts []types.Type, q func(types.Type) string) string {
	var parts []string
	for _, t := range ts {
		parts = append(parts, q(t))
	}
	return strings.Join(parts, ", ")
}

func joinResults(in *solve.Instance, q func(types.Type) string) string {
	var parts []string
	for _, r := range in.Results {
		switch r.Kind {
		case solve.ResultKindValue:
			parts = append(parts, "val "+q(r.Typ))
		case solve.ResultKindCleanup:
			if r.Failable {
				parts = append(parts, "cleanup func() error")
			} else {
				parts = append(parts, "cleanup func()")
			}
		case solve.ResultKindError:
			parts = append(parts, "err error")
		default:
			panic(fmt.Sprintf("plumb: unhandled result kind %d", r.Kind))
		}
	}
	return strings.Join(parts, ", ")
}

func injectorOut(pl *solve.Plan, q func(types.Type) string) string {
	var parts []string
	for _, o := range pl.Outputs {
		parts = append(parts, "val "+q(o))
	}
	if pl.AnyCleanup {
		if pl.CleanupFailable {
			parts = append(parts, "cleanup func() error")
		} else {
			parts = append(parts, "cleanup func()")
		}
	}
	if pl.Fallible {
		parts = append(parts, "err error")
	}
	return strings.Join(parts, ", ")
}
