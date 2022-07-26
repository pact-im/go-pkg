package zapjournal

import (
	"go.uber.org/zap/zapcore"
)

// defaultPrefix is the default variable name prefix. It is also used if the
// variable name starts with with underscore or is empty in appendVarName.
const defaultPrefix = "X"

// socketPath is the default socket path.
const socketPath = "/run/systemd/journal/socket"

// defaultConfig is the default configuration that we use if Config is nil.
var defaultConfig = Config{
	Level:  zapcore.DebugLevel,
	Prefix: defaultPrefix,
	Path:   socketPath,
}

// Config contains the options for the zapcore.Core implementation.
type Config struct {
	// Level decides whether a given logging level is enabled. Defaults to
	// all levels enabled (i.e. debug level).
	Level zapcore.LevelEnabler
	// Prefix is the variable name prefix for journald fields. Defaults to
	// "X".
	Prefix string
	// Path specifies the journal socket path to send logs to. Defaults to
	// "/run/systemd/journal/socket".
	Path string
}

// Option modifies the given configuration for the zapcore.Core implementation.
type Option func(*Config)

// WithLevel sets the Level configuration option.
func WithLevel(level zapcore.LevelEnabler) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithPrefix sets the Prefix configuration option.
func WithPrefix(prefix string) Option {
	return func(c *Config) {
		c.Prefix = prefix
	}
}

// WithPath sets the Path configuration option.
func WithPath(path string) Option {
	return func(c *Config) {
		c.Path = path
	}
}
