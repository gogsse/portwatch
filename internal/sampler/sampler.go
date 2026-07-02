// Package sampler provides probabilistic sampling for port events,
// reducing noise during high-frequency scan periods.
package sampler

import (
	"math/rand"
	"sync"
	"time"
)

// Sampler decides whether a given event should be emitted based on a
// configured sample rate. A rate of 1.0 always emits; 0.0 never emits.
type Sampler struct {
	mu      sync.Mutex
	rate    float64
	counts  map[int]uint64
	rng     *rand.Rand
}

// New returns a Sampler with the given sample rate clamped to [0.0, 1.0].
// Panics if rate is negative.
func New(rate float64) *Sampler {
	if rate < 0 {
		panic("sampler: rate must be >= 0")
	}
	if rate > 1.0 {
		rate = 1.0
	}
	//nolint:gosec // non-cryptographic use is intentional
	return &Sampler{
		rate:   rate,
		counts: make(map[int]uint64),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Allow reports whether the event for port should be emitted.
// It records every call so callers can inspect total volume via Count.
func (s *Sampler) Allow(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counts[port]++

	if s.rate >= 1.0 {
		return true
	}
	if s.rate <= 0.0 {
		return false
	}
	return s.rng.Float64() < s.rate
}

// Count returns the total number of times Allow has been called for port.
func (s *Sampler) Count(port int) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counts[port]
}

// Reset clears all per-port counters.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[int]uint64)
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
