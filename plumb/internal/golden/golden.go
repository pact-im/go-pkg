// Package golden compares test output against stored golden files, rewriting
// them instead when the PLUMB_GOLDEN_UPDATE environment variable is set.
package golden

import (
	"os"
	"testing"

	"go.pact.im/x/plumb/internal/writefile"
)

// updateEnv, when set to a non-empty value, switches the helpers from comparing
// against golden files to rewriting them.
const updateEnv = "PLUMB_GOLDEN_UPDATE"

// updating reports whether golden files should be rewritten this run.
func updating() bool {
	return os.Getenv(updateEnv) != ""
}

// Check compares got against the golden file at path and fails the test on a
// mismatch. When updating() is true it rewrites the file instead (creating it if
// absent); otherwise a missing file is fatal and names the env var to set.
func Check(t *testing.T, path, got string) {
	t.Helper()
	if updating() {
		if err := writefile.Write(path, got); err != nil {
			t.Fatalf("writing golden %s: %v", path, err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading golden %s (run with %s=1 to create it): %v", path, updateEnv, err)
	}
	if got != string(want) {
		t.Errorf("output differs from golden %s:\n--- got ---\n%s\n--- want ---\n%s", path, got, want)
	}
}
