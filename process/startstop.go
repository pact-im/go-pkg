package process

import (
	"context"
)

// Start starts the process with the given key from the table. It returns
// ErrManagerNotRunning if the manager is not running and ErrProcessExists
// error if the process already exists. Otherwise it returns an error from
// process initialization.
func (m *Manager[K, P]) Start(ctx context.Context, pk K) (P, error) {
	p, err := m.startProcessForKey(ctx, pk)
	if err != nil {
		var zero P
		return zero, err
	}
	return p.proc, nil
}

// Stop stops the process with the given key. It returns ErrProcessNotFound if
// the process does not exist, and an error from running the process otherwise.
func (m *Manager[K, P]) Stop(ctx context.Context, pk K) error {
	return m.stopProcess(ctx, pk)
}

// Get returns a running process or either a ErrProcessNotFound error if the
// process does not exist or ErrProcessNotRunning is the process exists but is
// not running.
func (m *Manager[K, P]) Get(ctx context.Context, pk K) (P, error) {
	p, ok := m.processes.Load(pk)
	if !ok || p == nil {
		var zero P
		return zero, ErrProcessNotFound
	}
	if p.State() != StateRunning {
		return p.proc, ErrProcessNotRunning
	}
	return p.proc, nil
}
