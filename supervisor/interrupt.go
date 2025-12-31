package supervisor

import (
	"context"
	"sync"
)

// supervisorInterrupter manages the interruption state between
// [flaky.Executor], [process.Runner] and [process.Callback].
type supervisorInterrupter struct {
	done   chan struct{}
	cancel context.CancelFunc

	stopMu  sync.Mutex
	stopped bool
	running bool
}

// Interrupt signals the supervisor to stop execution.
func (s *supervisorInterrupter) Interrupt() {
	s.stopMu.Lock()
	defer s.stopMu.Unlock()

	if s.stopped {
		return
	}
	s.stopped = true

	if s.running {
		close(s.done)
	} else {
		s.cancel()
	}
}

// shouldStopBeforeRun checks if an interrupt was requested before runner
// execution starts. It sets the running flag to true and returns whether
// execution should be stopped. It panics if the executor wasn’t properly
// interrupted after a stop was requested.
func (s *supervisorInterrupter) shouldStopBeforeRunner() bool {
	s.stopMu.Lock()
	defer s.stopMu.Unlock()
	if s.stopped && s.running {
		panic("supervisor: executor was not interrupted")
	}
	s.running = true
	return s.stopped
}

// afterRunner handles cleanup after the runner completes execution. It cancels the
// callback context if an interrupt was requested during execution.
func (s *supervisorInterrupter) afterRunner() {
	s.stopMu.Lock()
	defer s.stopMu.Unlock()

	// Stop executor after runner returns. Note that we do not unset running
	// flag to handle faulty executors (see panic in shouldStopBeforeRun).
	if s.stopped {
		s.cancel()
		return
	}

	s.running = false
}
