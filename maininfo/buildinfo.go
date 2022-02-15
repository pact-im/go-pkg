package maininfo

import (
	"runtime/debug"
	"sync"
)

var (
	modOnce    sync.Once
	modPath    string
	modVersion string
)

func readBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	mod := info.Main
	if v := mod.Replace; v != nil {
		mod = *v
	}
	modPath, modVersion = mod.Path, mod.Version
}

// Path returns main module’s path. It returns an empty string if an executable
// was built without module support.
//
// If the main module is replaced by another module, it returns the path of that
// module.
func Path() string {
	modOnce.Do(readBuildInfo)
	return modPath
}

// Version returns main module’s version. It returns an empty string if an
// executable was built without module support.
//
// If the main module is replaced by another module, it returns the version of
// that module.
func Version() string {
	modOnce.Do(readBuildInfo)
	return modVersion
}
