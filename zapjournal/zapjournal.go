// Package zapjournal provides a zapcore.Core implementation that sends logs to
// systemd-journald socket.
package zapjournal

import (
	"net"

	"go.uber.org/zap/zapcore"
)

// UnixConn is the subset of the net.UnixConn interface used by the zap.Core
// implementation that sends logs to journald.
type UnixConn interface {
	WriteMsgUnix(b, oob []byte, addr *net.UnixAddr) (n, oobn int, err error)
}

// Bind is a convenience function that returns a new unixgram client connection
// for communicating with journald.
func Bind() (*net.UnixConn, error) {
	return net.ListenUnixgram("unixgram", &net.UnixAddr{Net: "unixgram"})
}

// NewCore returns a new core that uses default configuration with the given
// options applied and sends logs to journald socket. On non-Linux platforms
// it always returns no-op implementation.
func NewCore(conn UnixConn, opts ...Option) zapcore.Core {
	c := defaultConfig
	for _, o := range opts {
		o(&c)
	}
	return newCoreWithConfig(conn, c)
}

// NewCoreWithConfig returns a new core that uses the given configuration and
// sends logs to journald socket. If c is nil, default values are used. On
// non-Linux platforms it always returns no-op implementation.
func NewCoreWithConfig(conn UnixConn, c *Config) zapcore.Core {
	cc := defaultConfig
	if c != nil {
		cc = *c
	}
	return newCoreWithConfig(conn, cc)
}

// Available checks whether the journald socket exists at the default path on
// the current system.
func Available() (bool, error) {
	return checkEnabled()
}
