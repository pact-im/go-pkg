package process

import (
	"errors"
)

var (
	// ErrManagerNotRunning is an error that is returned when a process
	// manager temporarily unavailable (i.e. it is not running).
	ErrManagerNotRunning = errors.New("process: manager is not running")

	// ErrProcessNotFound is an error that is returned if the requested
	// process was not found.
	ErrProcessNotFound = errors.New("process: process not found")

	// ErrProcessNotRunning is an error that is returned when a process is
	// not running (i.e. exists but is in the starting state).
	ErrProcessNotRunning = errors.New("process: process is not running")

	// ErrProcessExists is an error that is returned if the process already
	// exists for the given ID.
	ErrProcessExists = errors.New("process: process already exists")

	// ErrProcessInvalidState is an error that is returned if the process
	// is not in the valid state for the operation.
	ErrProcessInvalidState = errors.New("process: invalid process state")
)
