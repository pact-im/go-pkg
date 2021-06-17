package buildinfo

import (
	"runtime/debug"
)

// Version is the version string of this build.
var Version = func() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Version
}()
