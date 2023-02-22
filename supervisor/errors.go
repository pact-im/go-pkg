package supervisor

import (
	"errors"
)

var (
	// ErrNotRunning is an error that is returned when a process
	// supervisor temporarily unavailable (i.e. it is not running).
	ErrNotRunning = errors.New("supervisor: not running")

	// ErrProcessNotFound is an error that is returned if the requested
	// process was not found.
	ErrProcessNotFound = errors.New("supervisor: process not found")

	// ErrProcessNotRunning is an error that is returned when a process is
	// not running (i.e. exists but is in the starting state).
	ErrProcessNotRunning = errors.New("supervisor: process is not running")

	// ErrProcessExists is an error that is returned if the process already
	// exists for the given ID.
	ErrProcessExists = errors.New("supervisor: process already exists")
)
