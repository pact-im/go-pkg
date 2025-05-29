package config

import (
	"embed"
	"io/fs"
	"testing"
	"testing/fstest"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

type RootConfig struct {
	Nodes []NodeConfig `hcl:"node,block"`
	Env   []string     `hcl:"env,optional"`
}

type NodeConfig struct {
	ID               string         `hcl:"id,label"`
	Debug            bool           `hcl:"debug,optional"`
	Listen           []ListenConfig `hcl:"listen,block"`
	ReachableAddress string         `hcl:"reachable_address,optional"`
	Env              []string       `hcl:"env,optional"`
}

type ListenConfig struct {
	Type string `hcl:"type,label"`
	Addr string `hcl:"addr"`
}

func (r *RootConfig) ApplyEnv(env []string) {
	r.Env = append(env, r.Env...)
}

const (
	emptyJSON = "{}"
	emptyHCL  = ""
)

//go:embed _assets
var assets embed.FS

func TestExample(t *testing.T) {
	fsys, err := fs.Sub(assets, "_assets")
	assert.NilError(t, err)

	root := &RootConfig{}
	err = Load(fsys, root)
	assert.NilError(t, err)

	listen := []ListenConfig{
		{
			Type: "tcp6",
			Addr: "[::1]:80",
		},
		{
			Type: "unix",
			Addr: "/run/service/unix.sock",
		},
	}
	assert.Check(t, is.DeepEqual(&RootConfig{
		Nodes: []NodeConfig{{
			ID:               "example",
			Debug:            true,
			Listen:           listen,
			ReachableAddress: "example.com:666",
		}},
		Env: []string{
			"OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317",
		},
	}, root))
}

func TestLoadFS(t *testing.T) {
	testCases := []struct {
		File string
		Data string
		Conf *RootConfig
	}{
		// Invalid JSON data.
		{
			File: defaultJSON,
			Data: ``,
		},
		{
			File: defaultJSON,
			Data: `!`,
		},
		// Invalid schema.
		{
			File: defaultJSON,
			Data: `{
				"foobar": {
				},
			}`,
		},
		// Empty JSON is OK.
		{
			File: defaultJSON,
			Data: emptyJSON,
			Conf: &RootConfig{},
		},
		// Empty HCL is OK.
		{
			File: defaultHCL,
			Data: emptyHCL,
			Conf: &RootConfig{},
		},
		// HCL is OK.
		{
			File: defaultHCL,
			Data: `
				node "single" {
					env = [
						"FOO=BAR",
					]
				}

				env = [
					"BAR=FOO",
				]
			`,
			Conf: &RootConfig{
				Nodes: []NodeConfig{{
					ID:  "single",
					Env: []string{"FOO=BAR"},
				}},
				Env: []string{"BAR=FOO"},
			},
		},
		// The same HCL as JSON is OK.
		{
			File: defaultJSON,
			Data: `{
				"node": {
					"single": {
						"env": [ "FOO=BAR" ]
					}
				},
				"env": [ "BAR=FOO" ]
			}`,
			Conf: &RootConfig{
				Nodes: []NodeConfig{{
					ID:  "single",
					Env: []string{"FOO=BAR"},
				}},
				Env: []string{"BAR=FOO"},
			},
		},
		// Does not load JSON as HCL.
		{
			File: defaultHCL,
			Data: `{
				"node": {
					"single": {
					}
				}
			}`,
		},
		// Does not load HCL as JSON.
		{
			File: defaultJSON,
			Data: `
				node "single" {
				}
			`,
		},
	}
	for _, tc := range testCases {
		root := &RootConfig{}
		err := Load(fstest.MapFS{
			tc.File: &fstest.MapFile{
				Data: []byte(tc.Data),
			},
		}, root)

		if tc.Conf != nil {
			assert.Check(t, err)
			assert.Check(t, is.DeepEqual(tc.Conf, root))
		} else {
			assert.Check(t, is.ErrorContains(err, ""))
		}
	}
}

func TestLoadMerge(t *testing.T) {
	testCases := []struct {
		JSON string
		HCL  string
		Conf *RootConfig
	}{
		// Both empties are OK.
		{
			JSON: emptyJSON,
			HCL:  emptyHCL,
			Conf: &RootConfig{},
		},
		// Either invalid is an error.
		{
			JSON: "invalid",
			HCL:  emptyHCL,
		},
		{
			JSON: emptyJSON,
			HCL:  "invalid",
		},
		// HCL and JSON variable overrides, overwrites and inheritance.
		{
			JSON: `{
				"env": [ "FOO" ]
			}`,
			HCL: `
				env = [
					"BAR",
				]
			`,
			Conf: &RootConfig{
				Env: []string{"BAR"},
			},
		},
		{
			JSON: `{
				"node": {
					"foo": {
					}
				}
			}`,
			HCL: `
				node "bar" {
				}
			`,
			Conf: &RootConfig{
				Nodes: []NodeConfig{{
					ID: "bar",
				}},
			},
		},
		{
			JSON: `{
				"node": {
					"bar": {
						"debug": true
					}
				}
			}`,
			HCL: `
				node "bar" {
				}
			`,
			Conf: &RootConfig{
				Nodes: []NodeConfig{{
					ID:    "bar",
					Debug: true,
				}},
			},
		},
	}
	for _, tc := range testCases {
		root := &RootConfig{}
		err := Load(fstest.MapFS{
			defaultJSON: &fstest.MapFile{
				Data: []byte(tc.JSON),
			},
			defaultHCL: &fstest.MapFile{
				Data: []byte(tc.HCL),
			},
		}, root)
		if tc.Conf != nil {
			assert.Check(t, err)
			assert.Check(t, is.DeepEqual(tc.Conf, root))
		} else {
			assert.Check(t, is.ErrorContains(err, ""))
		}
	}
}

func TestLoadEnvNil(t *testing.T) {
	root := &RootConfig{}
	err := Load(fstest.MapFS{
		defaultHCL: &fstest.MapFile{
			Data: []byte(`
				node "" {
				}
			`),
		},
	}, root)
	assert.NilError(t, err)
	assert.DeepEqual(t, &RootConfig{Nodes: []NodeConfig{{}}}, root)
	assert.Check(t, is.Nil(root.Env))
	assert.Check(t, is.Nil(root.Nodes[0].Env))
}

func TestLoadNotFile(t *testing.T) {
	testCases := []fstest.MapFS{
		{
			// Empty FS.
		},
		{
			defaultJSON: &fstest.MapFile{
				Mode: fs.ModeDir,
			},
		},
		{
			defaultHCL: &fstest.MapFile{
				Mode: fs.ModeDir,
			},
		},
	}
	for _, tc := range testCases {
		root := &RootConfig{}
		err := Load(tc, root)
		assert.Check(t, err)
	}
}

func TestEvalFileDoesNotExist(t *testing.T) {
	root := &RootConfig{}
	err := Load(fstest.MapFS{
		defaultHCL: &fstest.MapFile{
			Data: []byte(`
				node "" {
					reachable_address = file("missing.txt")
				}
			`),
		},
	}, root)
	assert.Assert(t, is.ErrorContains(err, ""))
	assert.Assert(t, is.Contains(err.Error(), fs.ErrNotExist.Error()))
}

func TestEvalFileInvalidUTF8(t *testing.T) {
	root := &RootConfig{}
	err := Load(fstest.MapFS{
		defaultHCL: &fstest.MapFile{
			Data: []byte(`
				node "" {
					reachable_address = file("invalid-utf8.txt")
				}
			`),
		},
		"invalid-utf8.txt": &fstest.MapFile{
			Data: []byte{0xA0, 0xA1},
		},
	}, root)
	assert.Assert(t, is.ErrorContains(err, ""))
	assert.Assert(t, is.Contains(err.Error(), "contents of invalid-utf8.txt are not valid UTF-8"))
}

func TestLoadEnv(t *testing.T) {
	root := &RootConfig{}
	err := Load(fstest.MapFS{
		defaultHCL: &fstest.MapFile{
			Data: []byte(`env=["B"]`),
		},
		defaultEnv: &fstest.MapFile{
			Data: []byte("A"),
		},
	}, root)
	assert.NilError(t, err)
	assert.Check(t, is.DeepEqual(&RootConfig{Env: []string{"A", "B"}}, root))
}
