package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

type workspace struct {
	// Dir, if non-empty, is a directory of the Go workspace.
	Dir string
	// Paths is a list of modules in the current workspace. If Dir is empty,
	// it contains a single directory for the main module.
	Paths []string
}

// Root returns workspace root directory.
func (w *workspace) Root() string {
	if w.Dir == "" {
		return w.Paths[0]
	}
	return w.Dir
}

func loadWorkspace() (*workspace, error) {
	env, err := goenv("GOMOD", "GOWORK")
	if err != nil {
		return nil, err
	}

	gowork, ok := env["GOWORK"]
	if !ok || gowork == "" || gowork == "off" {
		gomod, ok := env["GOMOD"]
		if !ok || gomod == "" || gomod == os.DevNull {
			return nil, errors.New("no go.work or go.mod files found")
		}
		return &workspace{
			Paths: []string{
				filepath.Dir(gomod),
			},
		}, nil
	}

	workspaceDir := filepath.Dir(gowork)

	var buf bytes.Buffer
	c := exec.Command("go", "work", "edit", "-json")
	c.Stdout = &buf
	c.Stderr = os.Stderr
	c.Dir = workspaceDir
	if err := c.Run(); err != nil {
		return nil, err
	}
	var work struct {
		Use []struct {
			DiskPath string
		}
	}
	if err := json.Unmarshal(buf.Bytes(), &work); err != nil {
		return nil, err
	}

	paths := make([]string, len(work.Use))
	for i, m := range work.Use {
		paths[i] = filepath.Join(workspaceDir, filepath.FromSlash(m.DiskPath))
	}
	return &workspace{
		Dir:   workspaceDir,
		Paths: paths,
	}, nil
}

// syncWorkspace runs go mod tidy and go work sync for the given workspace.
func syncWorkspace(w *workspace) error {
	for _, moduleDir := range w.Paths {
		c := exec.Command("go", "mod", "tidy")
		c.Stderr = os.Stderr
		c.Dir = moduleDir
		if err := c.Run(); err != nil {
			return err
		}
	}
	return goworksync(w)
}

// goworksync removes current go.work.sum file and runs go work sync.
func goworksync(w *workspace) error {
	if w.Dir == "" {
		return nil
	}
	if err := os.Remove(filepath.Join(w.Dir, "go.work.sum")); err != nil && !os.IsNotExist(err) {
		return err
	}
	c := exec.Command("go", "work", "sync")
	c.Stderr = os.Stderr
	c.Dir = w.Dir
	return c.Run()
}

// goenv returns the given Go environment variables.
func goenv(vars ...string) (map[string]string, error) {
	var buf bytes.Buffer
	c := exec.Command("go", append([]string{"env", "-json", "--"}, vars...)...)
	c.Stdout = &buf
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return nil, err
	}
	out := make(map[string]string)
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		return nil, err
	}
	return out, nil
}
