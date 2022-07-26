//go:build linux
// +build linux

package zapjournal

import (
	"fmt"
	"net"
	"os"

	"go.uber.org/zap/zapcore"
	"golang.org/x/sys/unix"
)

// memfdName is a name used for memfd file descriptor.
const memfdName = "zapjournal"

// checkEnabled returns true if journald socket exists on the current system.
func checkEnabled() (bool, error) {
	_, err := os.Stat(socketPath)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func newCoreWithConfig(conn UnixConn, c Config) zapcore.Core {
	enc := getVarsEncoder(c.Prefix)
	return &journalCore{
		LevelEnabler: c.Level,
		path:         c.Path,
		conn:         conn,
		enc:          enc,
	}
}

type journalCore struct {
	zapcore.LevelEnabler
	path string
	conn UnixConn
	enc  *varsEncoder
}

func (c *journalCore) With(fields []zapcore.Field) zapcore.Core {
	e := cloneVarsEncoder(c.enc)
	addFields(e, fields)
	return &journalCore{
		LevelEnabler: c.LevelEnabler,
		path:         c.path,
		conn:         c.conn,
		enc:          e,
	}
}

func (c *journalCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	enc := c.enc.encodeEntry(ent, fields)
	defer putVarsEncoder(enc)

	socketAddr := &net.UnixAddr{
		Name: c.path,
		Net:  "unixgram",
	}

	_, _, err := c.conn.WriteMsgUnix(enc.buf, nil, socketAddr)
	if err == nil {
		return nil
	}
	if !isSocketSpaceError(err) {
		return err
	}

	fd, err := unix.MemfdCreate(memfdName, unix.MFD_CLOEXEC|unix.MFD_ALLOW_SEALING)
	if err != nil {
		return fmt.Errorf("create memfd: %w", err)
	}

	f := os.NewFile(uintptr(fd), memfdName)
	defer func() { _ = f.Close() }()

	if _, err := f.Write(enc.buf); err != nil {
		return fmt.Errorf("write to memfd: %w", err)
	}

	seals := unix.F_SEAL_SHRINK | unix.F_SEAL_GROW | unix.F_SEAL_WRITE | unix.F_SEAL_SEAL
	if _, err := unix.FcntlInt(uintptr(fd), unix.F_ADD_SEALS, seals); err != nil {
		return fmt.Errorf("seal memfd: %w", err)
	}

	oob := unix.UnixRights(fd)
	_, _, err = c.conn.WriteMsgUnix(nil, oob, socketAddr)
	return err
}

func (c *journalCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *journalCore) Sync() error {
	return nil
}

// isSocketSpaceError checks whether the error is signaling an "overlarge
// message" condition.
func isSocketSpaceError(err error) bool {
	opErr, ok := err.(*net.OpError)
	if !ok || opErr == nil {
		return false
	}

	sysErr, ok := opErr.Err.(*os.SyscallError)
	if !ok || sysErr == nil {
		return false
	}

	return sysErr.Err == unix.EMSGSIZE || sysErr.Err == unix.ENOBUFS
}
