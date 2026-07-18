package gopackages_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"

	"go.pact.im/x/plumb/internal/gen"
	"go.pact.im/x/plumb/internal/golden"
	"go.pact.im/x/plumb/internal/gopackages"
)

// Each loader fixture is a txtar archive: a config file, a go.mod, and sources.
// The loader needs a real on-disk module (it drives go/packages), so the archive
// is materialized into a temp module per run, the path the packagestest-based
// gen corpus never exercises. A successful load is run through the core and its
// generated source compared to a sibling .golden (regenerate with
// PLUMB_GOLDEN_UPDATE=1); a
// fixture whose config sets "load-error: true" instead expects Load to fail.

type loaderConfig struct {
	opts       gen.Options
	loadError  bool
	typeErrors bool // expect a tolerated type error to be collected, not fatal
}

func parseLoaderConfig(t *testing.T, s string, cfg *loaderConfig) {
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
			cfg.opts.ImportPath = v
		case "package-name":
			cfg.opts.PackageName = v
		case "output":
			cfg.opts.OutputBase = v
		case "load-error":
			b, err := strconv.ParseBool(v)
			if err != nil {
				t.Fatalf("config: invalid load-error %q: %v", v, err)
			}
			cfg.loadError = b
		case "type-errors":
			b, err := strconv.ParseBool(v)
			if err != nil {
				t.Fatalf("config: invalid type-errors %q: %v", v, err)
			}
			cfg.typeErrors = b
		}
	}
}

func TestLoadCorpus(t *testing.T) {
	t.Setenv("GOWORK", "off")
	paths, err := filepath.Glob(filepath.FromSlash("testdata/*.txtar"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("no loader fixtures found")
	}
	for _, p := range paths {
		name := strings.TrimSuffix(filepath.Base(p), ".txtar")
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			ar := txtar.Parse(data)

			var cfg loaderConfig
			dir := t.TempDir()
			for _, f := range ar.Files {
				if f.Name == "config" {
					parseLoaderConfig(t, string(f.Data), &cfg)
					continue
				}
				// txtar names are slash-separated; convert for the host filesystem.
				dst := filepath.Join(dir, filepath.FromSlash(f.Name))
				if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(dst, f.Data, 0o644); err != nil {
					t.Fatal(err)
				}
			}

			res, lerr := gopackages.Load([]string{"."}, dir)
			if cfg.loadError {
				if lerr == nil {
					t.Fatal("expected Load to fail, got nil")
				}
				return
			}
			if lerr != nil {
				t.Fatalf("Load: %v", lerr)
			}
			if cfg.typeErrors && len(res.TypeErrors) == 0 {
				t.Fatal("expected a tolerated type error to be collected, got none")
			}

			out, gerr := gen.Generate(cfg.opts, res.Packages)
			if gerr != nil {
				t.Fatalf("Generate: %v", gerr)
			}
			golden.Check(t, strings.TrimSuffix(p, ".txtar")+".golden", out.Source)
		})
	}
}
