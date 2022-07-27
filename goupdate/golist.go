package main

import (
	"encoding/json"
	"time"
)

type goModule struct {
	Path      string         // module path
	Version   string         // module version
	Versions  []string       // available module versions (with -versions)
	Replace   *goModule      // replaced by this module
	Time      *time.Time     // time version was created
	Update    *goModule      // available update, if any (with -u)
	Main      bool           // is this the main module?
	Indirect  bool           // is this module only an indirect dependency of main module?
	Dir       string         // directory holding files for this module, if any
	GoMod     string         // path to go.mod file used when loading this module, if any
	GoVersion string         // go version used in module
	Retracted string         // retraction information, if any (with -retracted or -u)
	Error     *goModuleError // error loading module
}

type goModuleError struct {
	Err string // the error itself
}

func golist() ([]goModule, error) {
	buf, err := system("go", "list", "-m", "-u", "-json", "all")
	if err != nil {
		return nil, err
	}
	var out []goModule
	dec := json.NewDecoder(buf)
	for dec.More() {
		var m goModule
		if err := dec.Decode(&m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}
