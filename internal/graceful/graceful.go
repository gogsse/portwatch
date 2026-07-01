// Package graceful provides utilities for coordinating clean shutdown of
// long-running goroutines when a context is cancelled or a signal is received.
package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Shutdown coordinates a graceful shutdown sequence.
// It listens for OS signals (SIGINT, SIGTERM) and cancels the provided context.
type Shutdown struct {
	cancel  context.CancelFunc
	timeout time.Duration
	wg      sync.WaitGroup
	mu      sync.Mutex
	done    chan struct{}
}

// New creates a new Shutdown coordinator. The provided cancel function is
// called when a termination signal is received. timeout specifies the maximum
// duration to wait for registered goroutines to finish before giving up.
func New(cancel context.CancelFunc, timeout time.Duration) *Shutdown {
	return &Shutdown{
		cancel:  cancel,
		timeout: timeout,
		done:    make(chan struct{}),
	}
}

// Register increments the internal WaitGroup, signalling that a goroutine
// is running and must complete before shutdown is considered clean.
func (s *Shutdown) Register() {
	s.wg.Add(1)
}

// Done decrements the internal WaitGroup. Each call to Register must be
// paired with exactly one call to Done (typically via defer).
func (s *Shutdown) Done() {
	s.wg.Done()
}

// Wait blocks until all registered goroutines have called Done or the
// configured timeout elapses, whichever comes first. It returns true if
// all goroutines finished cleanly, or false if the timeout was exceeded.
func (s *Shutdown) Wait() bool {
	finished := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		return true
	case <-time.After(s.timeout):
		return false
	}
}

// ListenAndServe starts a background goroutine that watches for SIGINT or
// SIGTERM. When a signal arrives the cancel function is invoked, giving all
// context-aware components the opportunity to stop cleanly.
//
// The returned stop function deregisters the signal handler and should be
// called when the process exits to free OS resources.
func (s *Shutdown) ListenAndServe() (stop func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigCh:
			s.cancel()
		case <-s.done:
		}
	}()

	return func() {
		signal.Stop(sigCh)
		s.mu.Lock()
		defer s.mu.Unlock()
		select {
		case <-s.done:
		default:
			close(s.done)
		}
	}
}
