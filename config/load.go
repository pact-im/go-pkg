package config

import (
	"io/fs"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const (
	defaultJSON = "config.json"
	defaultHCL  = "config.hcl"
	defaultEnv  = ".env"
)

var defaultPaths = []string{defaultJSON, defaultHCL}

// EnvAppender is an interface for configs that can accept environment variables.
type EnvAppender interface {
	ApplyEnv(env []string)
}

// Load loads the config into the provided schema pointer.
// If the schema implements EnvAppender, it appends environment variables from defaultEnv.
func Load(fsys fs.FS, schemaPtr any) error {
	for _, path := range defaultPaths {
		src, err := fs.ReadFile(fsys, path)
		if err != nil {
			if shouldSkipError(err) {
				continue
			}
			return err
		}
		if err := hclsimple.Decode(path, src, newEvalContext(fsys), schemaPtr); err != nil {
			return err
		}
	}

	f, err := fsys.Open(defaultEnv)
	if err != nil {
		if !shouldSkipError(err) {
			return err
		}
		// defaultEnv file missing is not an error, just return.
		return nil
	}
	defer func() { _ = f.Close() }()

	env, err := parseEnv(f)
	if err != nil {
		return err
	}

	if appender, ok := schemaPtr.(EnvAppender); ok {
		appender.ApplyEnv(env)
	}

	return nil
}
