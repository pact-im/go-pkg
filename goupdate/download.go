package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

type downloadedModule struct {
	Path     string // module path
	Version  string // module version
	Error    string // error loading module
	Info     string // absolute path to cached .info file
	GoMod    string // absolute path to cached .mod file
	Zip      string // absolute path to cached .zip file
	Dir      string // absolute path to cached source root directory
	Sum      string // checksum for path, version (as in go.sum)
	GoModSum string // checksum for go.mod (as in go.sum)
}

func queryUpgrades(workspaceDir string, modulePaths []string) ([]downloadedModule, error) {
	upgrades := make([]string, len(modulePaths))
	for i, modulePath := range modulePaths {
		upgrades[i] = modulePath + "@upgrade"
	}

	var buf bytes.Buffer
	c := exec.Command("go", append([]string{"mod", "download", "-json"}, upgrades...)...)
	c.Stderr = os.Stderr
	c.Stdout = &buf
	c.Dir = workspaceDir
	if err := c.Run(); err != nil {
		return nil, err
	}

	var out []downloadedModule
	for dec := json.NewDecoder(&buf); ; {
		var m downloadedModule
		err := dec.Decode(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}
