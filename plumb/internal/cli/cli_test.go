package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"go.pact.im/x/plumb/internal/discover"
	"go.pact.im/x/plumb/internal/golden"
)

func pkgs(specs ...[2]string) []*discover.Package {
	var out []*discover.Package
	for _, s := range specs {
		out = append(out, &discover.Package{PkgPath: s[0], Name: s[1]})
	}
	return out
}

func TestResolveDestination(t *testing.T) {
	tests := []struct {
		name        string
		pkgs        []*discover.Package
		importFlag  string
		nameFlag    string
		wantPath    string
		wantName    string
		wantErr     error
		wantErrText string
	}{
		{
			name:     "single_package_infers_both",
			pkgs:     pkgs([2]string{"example.com/app", "app"}),
			wantPath: "example.com/app",
			wantName: "app",
		},
		{
			name:    "multiple_packages_requires_import_path",
			pkgs:    pkgs([2]string{"example.com/a", "a"}, [2]string{"example.com/b", "b"}),
			wantErr: errAmbiguousDestination,
		},
		{
			name:    "no_packages_cannot_infer",
			pkgs:    nil,
			wantErr: errNoPackages,
		},
		{
			name:       "invalid_import_path",
			pkgs:       pkgs([2]string{"example.com/app", "app"}),
			importFlag: "bad path",
			wantErr:    errInvalidImportPath,
		},
		{
			name:       "import_path_picks_same_package",
			pkgs:       pkgs([2]string{"example.com/a", "a"}, [2]string{"example.com/b", "b"}),
			importFlag: "example.com/b",
			wantPath:   "example.com/b",
			wantName:   "b",
		},
		{
			name:       "separate_package_requires_package_name",
			pkgs:       pkgs([2]string{"example.com/app", "app"}),
			importFlag: "example.com/wire",
			wantErr:    errPackageNameRequired,
		},
		{
			name:       "separate_package_with_name",
			pkgs:       pkgs([2]string{"example.com/app", "app"}),
			importFlag: "example.com/wire",
			nameFlag:   "wire",
			wantPath:   "example.com/wire",
			wantName:   "wire",
		},
		{
			name:        "command_line_arguments_inferred_rejected",
			pkgs:        pkgs([2]string{"command-line-arguments", "app"}),
			wantErr:     errSyntheticImportPath,
			wantErrText: "no real import path: could not infer one from the arguments; pass -import-path",
		},
		{
			name:        "command_line_arguments_explicit_rejected",
			pkgs:        pkgs([2]string{"command-line-arguments", "app"}),
			importFlag:  "command-line-arguments",
			wantErr:     errSyntheticImportPath,
			wantErrText: `no real import path: "command-line-arguments" is not a real, importable package`,
		},
		{
			name:     "invalid_package_name",
			pkgs:     pkgs([2]string{"example.com/app", "app"}),
			nameFlag: "123",
			wantErr:  errInvalidPackageName,
		},
		{
			name:     "blank_package_name_rejected",
			pkgs:     pkgs([2]string{"example.com/app", "app"}),
			nameFlag: "_",
			wantErr:  errInvalidPackageName,
		},
		{
			name:     "mismatched_package_name_rejected",
			pkgs:     pkgs([2]string{"example.com/app", "app"}),
			nameFlag: "wrongname",
			wantErr:  errPackageNameMismatch,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path, name, err := resolveDestination(tc.pkgs, tc.importFlag, tc.nameFlag)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("want error wrapping %v, got %v", tc.wantErr, err)
				}
				if tc.wantErrText != "" && err.Error() != tc.wantErrText {
					t.Fatalf("error text\n got: %q\nwant: %q", err.Error(), tc.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if path != tc.wantPath || name != tc.wantName {
				t.Fatalf("got (%q, %q), want (%q, %q)", path, name, tc.wantPath, tc.wantName)
			}
		})
	}
}

// TestRunHelpAndBadFlag pins the exit codes of the two flag-parse outcomes:
// explicitly requested help is a success, a genuine parse error is a usage
// error.
func TestRunHelpAndBadFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"-h"}, &stdout, &stderr, &stderr); code != 0 {
		t.Errorf("Run(-h) exit %d, want 0", code)
	}
	if stderr.Len() == 0 {
		t.Errorf("Run(-h) printed no usage text")
	}
	stderr.Reset()
	if code := Run([]string{"-no-such-flag"}, &stdout, &stderr, &stderr); code != 2 {
		t.Errorf("Run(-no-such-flag) exit %d, want 2", code)
	}
}

// TestRunEndToEnd exercises the real package loader and the CLI writing a file.
func TestRunEndToEnd(t *testing.T) {
	// Capture the package dir before t.Chdir so the golden path still resolves.
	pkgDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := newTestModule(t)
	appDir := filepath.Join(dir, "app")
	if err := os.Mkdir(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(appDir, "app.go"), `package app

type Config struct{ Addr string }
type Server struct{}

//plumb:build
func NewConfig() *Config { return &Config{Addr: ":8080"} }

//plumb:build
func NewServer(c *Config) (*Server, error) { return &Server{}, nil }
`)

	t.Chdir(appDir)
	var stdout, stderr bytes.Buffer
	code := Run([]string{"-output=plumb_gen.go", "."}, &stdout, &stderr, &stderr)
	if code != 0 {
		t.Fatalf("Run exit %d, stderr:\n%s", code, stderr.String())
	}
	out, err := os.ReadFile(filepath.Join(appDir, "plumb_gen.go"))
	if err != nil {
		t.Fatal(err)
	}
	golden.Check(t, filepath.Join(pkgDir, "testdata", "endtoend.golden"), string(out))
}

// TestRunIgnoresTestFileDirectives locks in a documented limitation: the real
// loader never scans _test.go files, so a //plumb: directive there produces no
// provider and no diagnostic. The package’s
// regular file consumes *Config and its _test.go file provides it; because the
// provider is invisible, *Config becomes an injector input instead of being
// wired from NewConfig. The golden would differ (a NewConfig() call) if the
// test file were scanned. Corpus fixtures cannot express this: packagestest
// loads every file by name, so only the real loader exercises the exclusion.
func TestRunIgnoresTestFileDirectives(t *testing.T) {
	pkgDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := newTestModule(t)
	appDir := filepath.Join(dir, "app")
	if err := os.Mkdir(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(appDir, "app.go"), `package app

type Config struct{ Addr string }
type Server struct{}

//plumb:build
func NewServer(c *Config) *Server { return &Server{} }
`)
	// The directive here must be ignored: if it were scanned, NewConfig would
	// supply *Config and the injector would take no parameter.
	mustWrite(t, filepath.Join(appDir, "app_test.go"), `package app

//plumb:build
func NewConfig() *Config { return &Config{Addr: ":8080"} }
`)

	t.Chdir(appDir)
	var stdout, stderr bytes.Buffer
	code := Run([]string{"-output=plumb_gen.go", "."}, &stdout, &stderr, &stderr)
	if code != 0 {
		t.Fatalf("Run exit %d, stderr:\n%s", code, stderr.String())
	}
	out, err := os.ReadFile(filepath.Join(appDir, "plumb_gen.go"))
	if err != nil {
		t.Fatal(err)
	}
	golden.Check(t, filepath.Join(pkgDir, "testdata", "ignore_test_directives.golden"), string(out))
}

func TestRunVerboseSurfacesToleratedErrors(t *testing.T) {
	// Capture the package dir before t.Chdir so the golden path still resolves.
	pkgDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := newTestModule(t)
	appDir := filepath.Join(dir, "app")
	if err := os.Mkdir(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// A clean provider plus an unrelated body with a type error the loader
	// tolerates and that never reaches the injector signature, so generation
	// succeeds but the load was not fully type-correct.
	mustWrite(t, filepath.Join(appDir, "app.go"), `package app

type Server struct{}

//plumb:build
func NewServer() *Server { return &Server{} }

func unrelated() { var _ int = "not an int" }
`)
	t.Chdir(appDir)

	// Default run: succeeds and writes nothing to stderr (the tolerated error stays
	// quiet without -v).
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"."}, &stdout, &stderr, &stderr); code != 0 {
		t.Fatalf("Run exit %d, stderr:\n%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("default run wrote to stderr:\n%s", stderr.String())
	}

	// -v run: the discovery report goes to its own writer, deterministic and
	// plumb-controlled, so pin it whole as a golden. The tolerated-error note goes
	// to stderr; its go/types bodies carry temp-dir paths and go-version-dependent
	// wording, so assert only that the note surfaced, not its variable contents.
	stdout.Reset()
	stderr.Reset()
	var report bytes.Buffer
	if code := Run([]string{"-v", "."}, &stdout, &stderr, &report); code != 0 {
		t.Fatalf("Run -v exit %d, stderr:\n%s", code, stderr.String())
	}
	golden.Check(t, filepath.Join(pkgDir, "testdata", "verbose_report.golden"), report.String())
	if !slices.Contains(strings.Split(stderr.String(), "\n"), noteHeader) {
		t.Errorf("-v did not surface the tolerated-error note; stderr:\n%s", stderr.String())
	}
}

// TestRunWritesToStdout covers writeOutput’s stdout branch: with no -output, the
// generated source is written to stdout and nothing to stderr.
func TestRunWritesToStdout(t *testing.T) {
	// Capture the package dir before t.Chdir so the golden path still resolves.
	pkgDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := newTestModule(t)
	appDir := filepath.Join(dir, "app")
	if err := os.Mkdir(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(appDir, "app.go"), `package app

type Server struct{}

//plumb:build
func NewServer() *Server { return &Server{} }
`)
	t.Chdir(appDir)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"."}, &stdout, &stderr, &stderr)
	if code != 0 {
		t.Fatalf("Run exit %d, want 0; stderr:\n%s", code, stderr.String())
	}
	// The generated source is fully deterministic (no temp paths), so pin it whole.
	golden.Check(t, filepath.Join(pkgDir, "testdata", "stdout_source.golden"), stdout.String())
	if stderr.Len() != 0 {
		t.Errorf("stderr not empty on a quiet success:\n%s", stderr.String())
	}
}

// TestRunWriteFailureReturnsError covers cli.go’s writeOutput error branch: a
// -output whose parent directory does not exist makes the atomic write fail
// (writefile never creates intermediate directories), so Run reports the failure
// on stderr and exits 1 without leaving any file behind.
func TestRunWriteFailureReturnsError(t *testing.T) {
	dir := newTestModule(t)
	appDir := filepath.Join(dir, "app")
	if err := os.Mkdir(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(appDir, "app.go"), `package app

type Server struct{}

//plumb:build
func NewServer() *Server { return &Server{} }
`)
	t.Chdir(appDir)

	// The parent directory "missing" does not exist; plumb never creates it, so the
	// atomic write fails when it tries to create its temp file there.
	var stdout, stderr bytes.Buffer
	code := Run([]string{"-output=missing/plumb_gen.go", "."}, &stdout, &stderr, &stderr)
	if code != 1 {
		t.Fatalf("Run exit %d, want 1; stderr:\n%s", code, stderr.String())
	}
	// fail() writes the diagnostic as "plumb: <err>\n"; the err body carries a
	// random temp-file name and an OS-dependent errno string, so pin only plumb’s
	// own diagnostic marker, not the variable tail. This branch writes nothing else
	// to stderr, so the prefix is exact.
	if !strings.HasPrefix(stderr.String(), "plumb: ") {
		t.Errorf("write failure produced no plumb diagnostic; stderr:\n%s", stderr.String())
	}
	// The write is atomic and creates no directories: nothing is left behind.
	if _, err := os.Stat(filepath.Join(appDir, "missing")); !os.IsNotExist(err) {
		t.Errorf("plumb created the missing parent directory; stat err = %v", err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func newTestModule(t *testing.T) string {
	t.Helper()
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "go.mod"), "module example.com/e2e\n\ngo 1.26.4\n")
	return dir
}
