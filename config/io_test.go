package config

import (
	"errors"
	"io/fs"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

var _ fs.FS = openFunc(nil)

type openFunc func(name string) (fs.File, error)

func (f openFunc) Open(name string) (fs.File, error) {
	return f(name)
}

type brokenFile struct{ err error }

func (b brokenFile) Stat() (fs.FileInfo, error) { return nil, b.err }
func (b brokenFile) Read(_ []byte) (int, error) { return 0, b.err }
func (b brokenFile) Close() error               { return nil }

func TestLoadErrors(t *testing.T) {
	sentinelError := errors.New("sentinel test error")

	testCases := []struct {
		fsys openFunc
		conf *RootConfig
	}{
		{
			fsys: func(_ string) (fs.File, error) {
				return nil, fs.ErrNotExist
			},
			conf: &RootConfig{},
		},
		{
			fsys: func(_ string) (fs.File, error) {
				return nil, sentinelError
			},
		},
		{
			fsys: func(name string) (fs.File, error) {
				if name != defaultEnv {
					return nil, fs.ErrNotExist
				}
				return nil, sentinelError
			},
		},
		{
			fsys: func(name string) (fs.File, error) {
				if name != defaultEnv {
					return nil, fs.ErrNotExist
				}
				return brokenFile{sentinelError}, nil
			},
		},
	}

	for _, tc := range testCases {
		root := &RootConfig{}
		err := Load(tc.fsys, root)

		if tc.conf == nil {
			assert.ErrorIs(t, err, sentinelError)
		} else {
			assert.Check(t, is.DeepEqual(tc.conf, root))
		}
	}
}
