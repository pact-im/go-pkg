package gen_test

import (
	"errors"
	"maps"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"

	"go.pact.im/x/plumb/internal/diag"
	"go.pact.im/x/plumb/internal/gen"
	"go.pact.im/x/plumb/internal/golden"
	"go.pact.im/x/plumb/internal/packagestest"
)

// A fixture is a txtar archive: source files at import-path-encoded paths, plus a
// "config" file giving the generation options. A success fixture’s .golden holds
// the expected generated source. A fixture whose config names a sentinel with
// "error-is" is instead an error fixture: generation must fail, and its .golden
// holds the diagnostic string. Both are regenerated with UPDATE=1.

type fixture struct {
	files     map[string]string
	opts      gen.Options
	errIs     string // sentinel name the error must wrap; set marks an error fixture
	noCompile bool   // skip the compile check: the fixture is invalid Go that plumb tolerates best-effort
}

// sentinels maps the name used in an "error-is" config line to the exported
// sentinel it must wrap, so an error fixture pins both the message (golden) and
// the programmatic classification.
var sentinels = map[string]error{
	"ErrNoDirectives":              diag.ErrNoDirectives,
	"ErrInvalidSetName":            diag.ErrInvalidSetName,
	"ErrMisplacedDirective":        diag.ErrMisplacedDirective,
	"ErrDuplicateDirective":        diag.ErrDuplicateDirective,
	"ErrInvalidConversion":         diag.ErrInvalidConversion,
	"ErrUntypedConstant":           diag.ErrUntypedConstant,
	"ErrBlankProvider":             diag.ErrBlankProvider,
	"ErrInitProvider":              diag.ErrInitProvider,
	"ErrDestShadowsPredeclared":    diag.ErrDestShadowsPredeclared,
	"ErrEmbeddedField":             diag.ErrEmbeddedField,
	"ErrEmbeddedInterface":         diag.ErrEmbeddedInterface,
	"ErrConstraintInterfaceMethod": diag.ErrConstraintInterfaceMethod,
	"ErrStructProvider":            diag.ErrStructProvider,
	"ErrAmbiguousProducer":         diag.ErrAmbiguousProducer,
	"ErrMultipleErrors":            diag.ErrMultipleErrors,
	"ErrDependencyCycle":           diag.ErrDependencyCycle,
	"ErrUnusedTemplate":            diag.ErrUnusedTemplate,
	"ErrBareTypeParamResult":       diag.ErrBareTypeParamResult,
	"ErrAmbiguousTemplates":        diag.ErrAmbiguousTemplates,
	"ErrNonTerminating":            diag.ErrNonTerminating,
	"ErrUnexportedProvider":        diag.ErrUnexportedProvider,
	"ErrUnreachableType":           diag.ErrUnreachableType,
	"ErrInvalidType":               diag.ErrInvalidType,
	"ErrReservedName":              diag.ErrReservedName,
	"ErrShadowsPredeclared":        diag.ErrShadowsPredeclared,
	"ErrSetNameCollision":          diag.ErrSetNameCollision,
}

func parseFixture(t *testing.T, data []byte) fixture {
	t.Helper()
	ar := txtar.Parse(data)
	fx := fixture{files: map[string]string{}}
	for _, f := range ar.Files {
		if f.Name == "config" {
			parseConfig(t, string(f.Data), &fx)
			continue
		}
		fx.files[f.Name] = string(f.Data)
	}
	return fx
}

func parseConfig(t *testing.T, s string, fx *fixture) {
	t.Helper()
	for line := range strings.SplitSeq(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		switch k {
		case "import-path":
			fx.opts.ImportPath = v
		case "package-name":
			fx.opts.PackageName = v
		case "output":
			fx.opts.OutputBase = v
		case "error-is":
			fx.errIs = v
		case "compile":
			b, err := strconv.ParseBool(v)
			if err != nil {
				t.Fatalf("config: invalid compile %q: %v", v, err)
			}
			fx.noCompile = !b
		}
	}
}

// determinismSeeds is how many distinct shuffled load orders each fixture is
// regenerated under to prove its output (or diagnostic) is independent of the
// order the loader presents packages and files in. Breadth across the whole
// corpus carries the determinism guarantee; the runtime re-randomizes map
// iteration on every run on top of the explicit shuffle.
const determinismSeeds = 8

// shuffleForSeed returns a deterministic shuffle for the given seed, so a
// determinism failure reproduces from the reported seed.
func shuffleForSeed(seed int) packagestest.ShuffleFunc {
	return rand.New(rand.NewPCG(uint64(seed)+1, 0)).Shuffle
}

func TestCorpus(t *testing.T) {
	paths, err := filepath.Glob(filepath.FromSlash("testdata/*.txtar"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("no corpus fixtures found")
	}
	for _, p := range paths {
		name := strings.TrimSuffix(filepath.Base(p), ".txtar")
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			fx := parseFixture(t, data)
			loaded, err := packagestest.Load(fx.files)
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			res, gerr := gen.Generate(fx.opts, loaded.Packages)
			goldenPath := strings.TrimSuffix(p, ".txtar") + ".golden"

			if fx.errIs != "" {
				checkErrorFixture(t, fx, goldenPath, gerr)
				return
			}

			if gerr != nil {
				t.Fatalf("generate: %v", gerr)
			}

			// Generated output must type-check in context, unless the fixture is
			// deliberately invalid Go that plumb tolerates best-effort.
			if !fx.noCompile {
				compileCheck(t, fx.files, fx.opts, res.Source)
			}

			// Golden comparison: the generated source and the -v report. The report
			// (discovered providers, instantiation count, inferred signature, import
			// aliases) is deterministic and pinned alongside the source, so a change
			// in resolution that the source alone renders opaquely shows up here too.
			golden.Check(t, goldenPath, res.Source)
			golden.Check(t, strings.TrimSuffix(p, ".txtar")+".report.golden", res.Report)

			// Determinism: the output is byte-identical no matter how the
			// loader orders packages and files (shuffled per seed), and
			// regardless of the runtime’s per-run map-iteration randomness.
			for seed := range determinismSeeds {
				l2, err := packagestest.LoadShuffled(fx.files, shuffleForSeed(seed))
				if err != nil {
					t.Fatalf("shuffled reload (seed %d): %v", seed, err)
				}
				r2, err := gen.Generate(fx.opts, l2.Packages)
				if err != nil {
					t.Fatalf("regenerate (seed %d): %v", seed, err)
				}
				if r2.Source != res.Source {
					t.Fatalf("non-deterministic source under shuffle (seed %d)", seed)
				}
				if r2.Report != res.Report {
					t.Fatalf("non-deterministic -v report under shuffle (seed %d):\n--- ref ---\n%s\n--- got ---\n%s", seed, res.Report, r2.Report)
				}
			}

			// Regenerating over the previous output must reproduce it byte-for-byte, so a
			// go generate + CI-diff workflow never reports a spurious change for unchanged
			// input. Feed the emitted file back into the destination package and require an
			// identical result with no false set-name collision against the output file’s
			// own re-scanned declarations. This holds in same-package mode and across the
			// separate-package transition, where the destination is unscanned on the first
			// run and scanned (now holding the prior output) on the second. A noCompile
			// fixture’s output may not reload cleanly, so it is skipped.
			if fx.opts.OutputBase != "" && !fx.noCompile {
				regen := map[string]string{}
				maps.Copy(regen, fx.files)
				regen[fx.opts.ImportPath+"/"+fx.opts.OutputBase] = res.Source
				l3, err := packagestest.Load(regen)
				if err != nil {
					t.Fatalf("regeneration reload: %v", err)
				}
				r3, err := gen.Generate(fx.opts, l3.Packages)
				if err != nil {
					t.Fatalf("regeneration generate: %v", err)
				}
				if r3.Source != res.Source {
					t.Fatalf("not idempotent: regenerating over the emitted file changed the output\n--- first ---\n%s\n--- second ---\n%s", res.Source, r3.Source)
				}
			}
		})
	}
}

// checkErrorFixture verifies an error fixture: generation must fail, the
// diagnostic must wrap the named sentinel, and its string is held in the golden.
// The message is checked for stability across re-runs, since source positions
// and error selection must be deterministic.
func checkErrorFixture(t *testing.T, fx fixture, goldenPath string, gerr error) {
	t.Helper()
	if gerr == nil {
		t.Fatalf("expected generation to fail, but it succeeded")
	}
	sentinel, ok := sentinels[fx.errIs]
	if !ok {
		t.Fatalf("config error-is: unknown sentinel %q", fx.errIs)
	}
	if !errors.Is(gerr, sentinel) {
		t.Errorf("error %q does not wrap sentinel %s", gerr.Error(), fx.errIs)
	}

	got := gerr.Error() + "\n"
	golden.Check(t, goldenPath, got)

	// Determinism: the diagnostic (its position and the provider it selects on
	// a tie) must be stable no matter how the loader orders packages and files.
	for seed := range determinismSeeds {
		l2, err := packagestest.LoadShuffled(fx.files, shuffleForSeed(seed))
		if err != nil {
			t.Fatal(err)
		}
		_, e2 := gen.Generate(fx.opts, l2.Packages)
		if e2 == nil {
			t.Fatalf("seed %d: expected generation to fail", seed)
		}
		if e2.Error()+"\n" != got {
			t.Fatalf("non-deterministic error under shuffle (seed %d):\n got: %s\nwant: %s", seed, e2.Error(), got)
		}
	}
}

// compileCheck overlays the generated source onto the fixture under the
// destination import path and asserts the augmented module type-checks. This is
// what turns “valid input ⇒ compilable output” into a tested guarantee.
func compileCheck(t *testing.T, files map[string]string, opts gen.Options, src string) {
	t.Helper()
	aug := map[string]string{}
	maps.Copy(aug, files)
	genFile := opts.ImportPath + "/plumb_generated_check.go"
	aug[genFile] = src
	loaded, err := packagestest.Load(aug)
	if err != nil {
		t.Fatalf("compile-check load failed: %v\n--- generated ---\n%s", err, src)
	}
	if len(loaded.TypeErrors) != 0 {
		var b strings.Builder
		for _, e := range loaded.TypeErrors {
			b.WriteString(e.Error())
			b.WriteByte('\n')
		}
		t.Fatalf("generated code does not type-check:\n%s--- generated ---\n%s", b.String(), src)
	}
}
