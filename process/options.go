package process

import (
	"go.pact.im/x/clock"
)

// Options is a set of options for Manager constructor.
type Options struct {
	// Clock is the clock to use. Defaults to system clock.
	Clock *clock.Clock
}

// setDefaults sets default values for unspecified options.
func (o *Options) setDefaults() {
	if o.Clock == nil {
		o.Clock = clock.System()
	}
}
