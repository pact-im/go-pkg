package main

import (
	"encoding/json"
	"errors"
	"path/filepath"
)

func gowork() (string, []string, error) {
	env, err := goenv("GOMOD", "GOWORK")
	if err != nil {
		return "", nil, err
	}

	gowork, ok := env["GOWORK"]
	if !ok || gowork == "" {
		gomod, ok := env["GOMOD"]
		if !ok || gomod == "" {
			return "", nil, errors.New("no go.work or go.mod found")
		}
		return "", []string{filepath.Dir(gomod)}, nil
	}

	var work struct {
		Use []struct {
			DiskPath string
		}
	}
	buf, err := system("go", "work", "edit", "-json")
	if err != nil {
		return "", nil, err
	}
	if err := json.Unmarshal(buf.Bytes(), &work); err != nil {
		return "", nil, err
	}

	dir := filepath.Dir(gowork)
	paths := make([]string, len(work.Use))
	for i, m := range work.Use {
		paths[i] = filepath.Join(dir, filepath.FromSlash(m.DiskPath))
	}
	return dir, paths, nil
}
