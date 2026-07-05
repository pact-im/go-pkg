package writefile_test

import (
	"os"
	"path/filepath"
	"testing"

	"go.pact.im/x/plumb/internal/writefile"
)

func TestWrite(t *testing.T) {
	t.Run("creates a new file with the given contents", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "out.txt")
		if err := writefile.Write(path, "hello"); err != nil {
			t.Fatalf("Write: %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if string(got) != "hello" {
			t.Errorf("contents = %q, want %q", got, "hello")
		}
	})

	t.Run("overwrites an existing file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "out.txt")
		if err := os.WriteFile(path, []byte("old contents, longer"), 0o600); err != nil {
			t.Fatalf("seed: %v", err)
		}
		if err := writefile.Write(path, "new"); err != nil {
			t.Fatalf("Write: %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if string(got) != "new" {
			t.Errorf("contents = %q, want %q", got, "new")
		}
	})

	t.Run("leaves no temp file behind on success", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "out.txt")
		if err := writefile.Write(path, "x"); err != nil {
			t.Fatalf("Write: %v", err)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ReadDir: %v", err)
		}
		if len(entries) != 1 || entries[0].Name() != "out.txt" {
			names := make([]string, len(entries))
			for i, e := range entries {
				names[i] = e.Name()
			}
			t.Errorf("dir entries = %v, want [out.txt] only (a stray temp file was left)", names)
		}
	})

	t.Run("returns an error and writes nothing when the parent is missing", func(t *testing.T) {
		dir := t.TempDir()
		missing := filepath.Join(dir, "nope")
		path := filepath.Join(missing, "out.txt")
		if err := writefile.Write(path, "x"); err == nil {
			t.Fatal("Write into a missing directory succeeded; want an error")
		}
		if _, err := os.Stat(missing); !os.IsNotExist(err) {
			t.Errorf("missing parent %q now exists; Write must not create intermediate directories", missing)
		}
	})

	t.Run("removes the temp file when the write fails", func(t *testing.T) {
		// Force a failure after the temp file is created and written, so the deferred
		// Remove is what must clean up: make the target an existing directory, which
		// the final rename cannot replace. No .tmp residue may remain.
		dir := t.TempDir()
		target := filepath.Join(dir, "out")
		if err := os.Mkdir(target, 0o755); err != nil {
			t.Fatalf("mkdir target: %v", err)
		}
		if err := writefile.Write(target, "x"); err == nil {
			t.Fatal("Write over an existing directory succeeded; want an error")
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ReadDir: %v", err)
		}
		if len(entries) != 1 || entries[0].Name() != "out" {
			names := make([]string, len(entries))
			for i, e := range entries {
				names[i] = e.Name()
			}
			t.Errorf("dir entries = %v, want [out] only (a temp file was left after the failed write)", names)
		}
	})
}
