// Package supervisor provides a supervisor implementation for starting,
// stopping, and monitoring its child processes.
package supervisor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.pact.im/x/clock"
	"go.pact.im/x/process"
)

// Supervisor is responsible for starting, stopping, and monitoring its child
// processes.
type Supervisor[K comparable, P process.Runner] struct {
	table Table[K, P]
	clock *clock.Clock

	// processes is a map of managed processes. It is used to track process
	// state and allows returning an ErrProcessExists error to guarantee
	// that at most one processes is active per key.
	processes typedMap[K, *managedProcess[P]]

	// runLock ensures that at most one Run method is executing at a time.
	runLock chanLock

	// startMu guards startProcess and startProcessForKey calls when
	// Supervisor is not running. It also allows waiting for the ongoing
	// calls to complete on shutdown.
	startMu sync.RWMutex
	start   bool

	// parent is the parent context for all processes. It is set to the
	// context passed to Run method and is guarded by startMu.
	parent context.Context

	// wg is the wait group for running processes and watchdogs. It is
	// indirectly guarded by startMu and start.
	wg sync.WaitGroup
}

// managedProcess contains a process.Process and associated [process.Runner]
// managed by Supervisor.
type managedProcess[P process.Runner] struct {
	*process.Process

	// runner is the underlying process entrypoint with parametrized type P.
	runner P

	// stopped is used by Supervisor to remove the process instance from
	// internal map at most once.
	stopped atomic.Bool
}

// NewSupervisor returns a new Supervisor instance. The given table is used to
// lookup managed processes and periodically restart failed units (or processes
// that were added externally). Managed processes are uniquely identifiable by
// key.
//
// A managed process may remove itself from the Supervisor by deleting the
// associated entry from the table before terminating. Likewise, to stop a
// process, it must be removed from the table prior to Stop call. That is,
// processes must be aware of being managed and the removal is tighly coupled
// with the table.
//
// As a rule of thumb, to keep the underlying table consistent, processes should
// not be re-added to table after being removed from the table. It is possible
// to implement re-adding on top of the Supervisor but that requires handling
// possible orderings of table removal, addition, re-addition and process
// startup, shutdown and self-removal (or a subset of these operations depending
// on the use cases).
func NewSupervisor[K comparable, P process.Runner](t Table[K, P], o Options) *Supervisor[K, P] {
	o.setDefaults()

	return &Supervisor[K, P]{
		clock:   o.Clock,
		table:   t,
		runLock: newChanLock(),
	}
}

// Run starts the supervisor and executes callback on successful initialization.
func (m *Supervisor[K, P]) Run(ctx context.Context, callback process.Callback) error {
	if err := m.runLock.Acquire(ctx); err != nil {
		return err
	}
	defer m.runLock.Release()

	m.parent = ctx
	defer func() { m.parent = nil }()

	// Allow startProcess and startProcessForKey calls.
	m.startMu.Lock()
	m.start = true
	m.startMu.Unlock()

	// Restore last state from the storage.
	m.restartInitial(ctx)

	// Run restartLoop to keep the current state up-to-date with changes in
	// the storage.
	stopRestartLoop := m.spawnRestartLoop(ctx)

	// Invoke callback.
	err := callback(ctx)

	// Wait for restart loop since it uses wg to spawn background tasks.
	stopRestartLoop()

	// Block until all ongoing startProcess and startProcessForKey calls are
	// complete and forbid subsequent calls.
	m.startMu.Lock()
	m.start = false
	m.startMu.Unlock()

	// Stop all processes and wait for shutdown completion. At this point
	// we are guaranteed that new processes would not be started.
	m.stopAll(ctx)
	m.wg.Wait()
	return err
}

// startProcessForKey starts the process for the given key. An error is returned
// if Supervisor’s Run method is not currently running.
func (m *Supervisor[K, P]) startProcessForKey(ctx context.Context, pk K) (*managedProcess[P], error) {
	m.startMu.RLock()
	defer m.startMu.RUnlock()
	if !m.start {
		return nil, ErrNotRunning
	}
	if _, exists := m.processes.LoadOrStore(pk, nil); exists {
		return nil, ErrProcessExists
	}

	r, err := m.table.Get(ctx, pk)
	if err != nil {
		m.processes.Delete(pk)
		return nil, fmt.Errorf("get process from table: %w", err)
	}
	return m.startProcessUnlocked(ctx, pk, r)
}

// startProcess starts the process for the given key. Unlike startProcessForKey,
// it uses the given r [process.Runner] instance instead of getting it from the
// table.
func (m *Supervisor[K, P]) startProcess(ctx context.Context, pk K, r P) (*managedProcess[P], error) {
	m.startMu.RLock()
	defer m.startMu.RUnlock()
	if !m.start {
		return nil, ErrNotRunning
	}
	if _, exists := m.processes.LoadOrStore(pk, nil); exists {
		return nil, ErrProcessExists
	}
	return m.startProcessUnlocked(ctx, pk, r)
}

// startProcessUnlocked starts the given process assuming that the start lock
// was acquired and an entry in the processes map exists. It removes this entry
// on error.
func (m *Supervisor[K, P]) startProcessUnlocked(ctx context.Context, pk K, r P) (*managedProcess[P], error) {
	p := &managedProcess[P]{
		Process: process.NewProcess(m.parent, r),
		runner:  r,
	}
	m.processes.Store(pk, p)
	if err := p.Start(ctx); err != nil {
		m.processes.Delete(pk)
		return nil, fmt.Errorf("start process: %w", err)
	}
	m.wg.Add(1)
	go m.watchdog(pk, p)
	return p, nil
}

// watchdog waits for process shutdown and removes it from processes map on such
// event.
func (m *Supervisor[K, P]) watchdog(pk K, p *managedProcess[P]) {
	defer m.wg.Done()

	<-p.Done()

	// Remove process from the processes map unless we have been stopped
	// externally.
	if p.stopped.Swap(true) {
		return
	}
	m.processes.Delete(pk)

	_ = p.Err() // TODO: log error
}

// stopProcess stops the process with the given key. If the processes does not
// exist, it returns ErrProcessNotFound.
func (m *Supervisor[K, P]) stopProcess(ctx context.Context, pk K) error {
	p, ok := m.processes.Load(pk)
	if !ok || p == nil {
		return ErrProcessNotFound
	}

	// Something else has already stopped the given process. Do nothing.
	if p.stopped.Swap(true) {
		return ErrProcessNotFound
	}
	m.processes.Delete(pk)

	return p.Stop(ctx)
}

// stopAll stops all the processes in the underlying map. It does not wait for
// processes to complete the termination.
func (m *Supervisor[K, P]) stopAll(ctx context.Context) {
	m.processes.Range(func(pk K, _ *managedProcess[P]) bool {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			_ = m.stopProcess(ctx, pk)
		}()
		return true
	})
}
