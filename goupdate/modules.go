package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

// updateGoModVersion updates go directive in go.mod file.
func updateGoModVersion(moduleDir, goVersion string) (string, error) {
	old, err := loadGoModVersion(moduleDir)
	if err != nil {
		return "", err
	}

	c := exec.Command("go", "mod", "edit", "-go", goVersion)
	c.Stderr = os.Stderr
	c.Dir = moduleDir
	return old, c.Run()
}

// loadGoModVersion returns the value of the go directive in go.mod file.
func loadGoModVersion(moduleDir string) (string, error) {
	var buf bytes.Buffer
	c := exec.Command("go", "mod", "edit", "-json")
	c.Stdout = &buf
	c.Stderr = os.Stderr
	c.Dir = moduleDir
	if err := c.Run(); err != nil {
		return "", err
	}

	var m struct{ Go string }
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		return "", err
	}
	return m.Go, nil
}
