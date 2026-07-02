// Package backoff provides an exponential back-off strategy with optional
// jitter for use when retrying transient failures (e.g. webhook delivery,
// scanner I/O errors).
package backoff

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// Backoff computes successive wait durations using exponential back-off.
type Backoff struct {
	mu       sync.Mutex
	base     time.Duration
	max      time.Duration
	factor   float64
	jitter   bool
	attempts int
}

// New returns a Backoff with the given base delay, maximum delay, growth
// factor and jitter flag. It panics if base <= 0, max < base, or factor < 1.
func New(base, max time.Duration, factor float64, jitter bool) *Backoff {
	if base <= 0 {
		panic("backoff: base must be > 0")
	}
	if max < base {
		panic("backoff: max must be >= base")
	}
	if factor < 1 {
		panic("backoff: factor must be >= 1")
	}
	return &Backoff{base: base, max: max, factor: factor, jitter: jitter}
}

// Next returns the next back-off duration and increments the internal attempt
// counter. The returned value is capped at max.
func (b *Backoff) Next() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	d := float64(b.base) * math.Pow(b.factor, float64(b.attempts))
	if d > float64(b.max) {
		d = float64(b.max)
	}
	if b.jitter {
		// Uniform jitter: [d/2, d]
		d = d/2 + rand.Float64()*(d/2) //nolint:gosec
	}
	b.attempts++
	return time.Duration(d)
}

// Reset brings the back-off back to its initial state.
func (b *Backoff) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.attempts = 0
}

// Attempts returns the number of Next calls since the last Reset.
func (b *Backoff) Attempts() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.attempts
}
