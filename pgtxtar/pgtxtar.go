// Package pgtxtar loads migrations for github.com/go-pg/migrations/v8 package
// in ad-hoc txtar format from a fs.FS filesystem.
package pgtxtar

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"strings"

	"golang.org/x/tools/txtar"

	"github.com/go-pg/migrations/v8"
	"go.uber.org/multierr"
)

// LoadFS loads migrations from SQL script archives in filesystem.
func LoadFS(fsys fs.FS) ([]*migrations.Migration, error) {
	var parseErrors error

	var ms []*migrations.Migration
	err := fs.WalkDir(fsys, ".", func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // fatal
		}
		if d.IsDir() {
			return nil
		}

		fileBytes, err := fs.ReadFile(fsys, fpath)
		if err != nil {
			return err // fatal
		}

		name := strings.TrimSuffix(path.Base(fpath), path.Ext(fpath))

		parts := strings.SplitN(name, "_", 2)
		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			parseErrors = multierr.Append(parseErrors, fmt.Errorf("parse version: %q", err))
			return nil
		}

		ar := txtar.Parse(fileBytes)
		if err := checkArchive(ar); err != nil {
			parseErrors = multierr.Append(parseErrors, fmt.Errorf("migration %q: %v", name, err))
			return nil
		}

		m := &migrations.Migration{
			Version: version,
			UpTx:    true,
			DownTx:  true,
		}
		for _, file := range ar.Files {
			switch file.Name {
			case "up":
				m.Up = migrationExec(file.Data)
			case "down":
				m.Down = migrationExec(file.Data)
			case "disable_tx":
				m.UpTx = false
				m.DownTx = false
			}
		}
		ms = append(ms, m)
		return nil
	})
	return ms, multierr.Append(parseErrors, err)
}

func checkArchive(ar *txtar.Archive) error {
	var err error
	seen := map[string]int{}
	for _, file := range ar.Files {
		name := file.Name

		switch name {
		case "up", "down":
			// ok
		case "disable_tx":
			data := bytes.TrimSpace(file.Data)
			if len(data) != 0 {
				err = multierr.Append(err, fmt.Errorf(
					"section %q must be empty", name,
				))
			}
		default:
			err = multierr.Append(err, fmt.Errorf(
				"unknown section %q", name,
			))
		}

		if seen[name] == 1 {
			err = multierr.Append(err, fmt.Errorf(
				"duplicate section %q", name,
			))
		}
		seen[name]++
	}

	required := []string{"up"}
	for _, name := range required {
		if seen[name] != 0 {
			continue
		}
		err = multierr.Append(err, fmt.Errorf(
			"missing required section %q", name,
		))
	}
	return err
}

func migrationExec(b []byte) func(db migrations.DB) error {
	return func(db migrations.DB) error {
		_, err := db.Exec(string(b))
		return err
	}
}
