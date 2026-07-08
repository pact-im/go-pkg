// Package writefile writes a file atomically, so an interrupted or failed write
// cannot leave a truncated file where a previous good one was.
package writefile

import (
	"os"
	"path/filepath"
)

// Write atomically writes data to name with mode 0o644. It writes to a temp file
// in name’s own directory and renames it over name: an interrupted write cannot
// leave a partial file in place of the previous good one, and sharing the
// directory keeps the rename a same-filesystem swap. The parent directory must
// already exist (os.CreateTemp fails when it does not, so Write never creates
// it), and a failed write removes the temp file rather than leaving it behind.
func Write(name, data string) error {
	tmp, err := os.CreateTemp(filepath.Dir(name), "."+filepath.Base(name)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }() // no-op after a successful rename
	if _, err := tmp.WriteString(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// CreateTemp creates the file 0o600; set the conventional 0o644.
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpName, name)
}
