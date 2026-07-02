// Package circuitbreaker implements a simple circuit breaker that prevents
// repeated alerting when a downstream sink (webhook, notifier, etc.) is
// continuously failing.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the current circuit state.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing; requests are rejected
	StateHalfOpen              // probing whether the downstream has recovered
)

// ErrOpen is returned when the circuit is open and the call is rejected.
var ErrOpen = errors.New("circuitbreaker: circuit is open")

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	cooldown     time.Duration
	openedAt     time.Time
	successes    int
	probeNeeded  int
	now          func() time.Time
}

// New returns a Breaker that opens after threshold consecutive failures and
// attempts recovery after cooldown has elapsed.
func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{
		threshold:   threshold,
		cooldown:    cooldown,
		probeNeeded: 1,
		now:         time.Now,
	}
}

// Allow reports whether the call should proceed. It transitions state as
// needed and returns ErrOpen when the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if b.now().Sub(b.openedAt) >= b.cooldown {
			b.state = StateHalfOpen
			b.successes = 0
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess notifies the breaker that the last call succeeded.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateHalfOpen {
		b.successes++
		if b.successes >= b.probeNeeded {
			b.state = StateClosed
			b.failures = 0
		}
		return
	}
	b.failures = 0
}

// RecordFailure notifies the breaker that the last call failed.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateHalfOpen {
		b.state = StateOpen
		b.openedAt = b.now()
		return
	}
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// State returns the current circuit state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
