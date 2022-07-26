//go:build !linux
// +build !linux

package zapjournal

import (
	"go.uber.org/zap/zapcore"
)

func checkEnabled() (bool, error)                     { return false, nil }
func newCoreWithConfig(UnixConn, Config) zapcore.Core { return zapcore.NewNopCore() }
