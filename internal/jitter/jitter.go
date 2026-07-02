// Package jitter adds randomised delay to polling intervals to avoid
// thundering-herd problems when multiple portwatch instances run in
// parallel on the same host.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Jitter wraps a base interval and returns a slightly randomised
// duration each time Next is called.  The spread is ±factor of the
// base duration, clamped so the result is always > 0.
type Jitter struct {
	mu      sync.Mutex
	base    time.Duration
	factor  float64 // 0 < factor <= 1
	rng     *rand.Rand
}

// New creates a Jitter with the given base interval and spread factor.
// factor=0.1 means the returned duration will be within ±10 % of base.
// Panics if base <= 0 or factor is outside (0, 1].
func New(base time.Duration, factor float64) *Jitter {
	if base <= 0 {
		panic("jitter: base must be positive")
	}
	if factor <= 0 || factor > 1 {
		panic("jitter: factor must be in (0, 1]")
	}
	return &Jitter{
		base:   base,
		factor: factor,
		//nolint:gosec // non-cryptographic randomness is fine for scheduling jitter
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Next returns the next jittered duration.  It is safe for concurrent use.
func (j *Jitter) Next() time.Duration {
	j.mu.Lock()
	spread := j.rng.Float64()*2 - 1 // [-1, 1)
	j.mu.Unlock()

	delta := time.Duration(float64(j.base) * j.factor * spread)
	d := j.base + delta
	if d <= 0 {
		d = time.Millisecond
	}
	return d
}

// Base returns the configured base interval.
func (j *Jitter) Base() time.Duration {
	return j.base
}

// Factor returns the configured spread factor.
func (j *Jitter) Factor() float64 {
	return j.factor
}
