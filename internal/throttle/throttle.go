// Package throttle provides a token-bucket style throttle that limits how
// many alerts can be emitted across all ports within a rolling time window.
package throttle

import (
	"sync"
	"time"
)

// Throttle enforces a global cap on the number of events allowed within a
// sliding window. Once the cap is reached, Allow returns false until enough
// time has passed that older tokens have expired.
type Throttle struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	timestamps []time.Time
	now      func() time.Time // injectable for testing
}

// New creates a Throttle that permits at most max events per window duration.
func New(max int, window time.Duration) *Throttle {
	return &Throttle{
		max:    max,
		window: window,
		now:    time.Now,
	}
}

// Allow reports whether a new event may proceed. If the number of events
// recorded within the current window is below max, the event is recorded and
// true is returned. Otherwise false is returned and nothing is recorded.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	t.evict(now)

	if len(t.timestamps) >= t.max {
		return false
	}

	t.timestamps = append(t.timestamps, now)
	return true
}

// Remaining returns the number of additional events that may be admitted
// within the current window.
func (t *Throttle) Remaining() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.evict(t.now())
	rem := t.max - len(t.timestamps)
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears all recorded timestamps, immediately restoring full capacity.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.timestamps = t.timestamps[:0]
}

// evict removes timestamps that have fallen outside the current window.
// Caller must hold t.mu.
func (t *Throttle) evict(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.timestamps) && t.timestamps[i].Before(cutoff) {
		i++
	}
	t.timestamps = t.timestamps[i:]
}
