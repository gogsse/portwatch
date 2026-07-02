// Package cooldown provides a per-port cooldown tracker that prevents
// re-alerting on the same port until a configurable quiet period has elapsed.
package cooldown

import (
	"sync"
	"time"
)

// entry holds the timestamp of the last alert for a port.
type entry struct {
	lastAlert time.Time
}

// Tracker tracks the last-alert time for each port and reports whether
// a new alert should be suppressed because the cooldown has not expired.
type Tracker struct {
	mu       sync.Mutex
	entries  map[int]entry
	cooldown time.Duration
	now      func() time.Time
}

// New returns a Tracker with the given cooldown duration.
// Panics if cooldown is zero or negative.
func New(cooldown time.Duration) *Tracker {
	if cooldown <= 0 {
		panic("cooldown: duration must be positive")
	}
	return &Tracker{
		entries:  make(map[int]entry),
		cooldown: cooldown,
		now:      time.Now,
	}
}

// Allow returns true and records the alert time when the port is not in
// cooldown. Returns false without updating state when the cooldown period
// has not yet elapsed since the last alert.
func (t *Tracker) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if e, ok := t.entries[port]; ok {
		if now.Sub(e.lastAlert) < t.cooldown {
			return false
		}
	}
	t.entries[port] = entry{lastAlert: now}
	return true
}

// Reset removes the cooldown record for a port so the next alert is
// always allowed through, regardless of when the previous one occurred.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Active returns the set of ports currently within their cooldown window.
func (t *Tracker) Active() []int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	var ports []int
	for port, e := range t.entries {
		if now.Sub(e.lastAlert) < t.cooldown {
			ports = append(ports, port)
		}
	}
	return ports
}

// Flush removes all entries whose cooldown has expired.
func (t *Tracker) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	for port, e := range t.entries {
		if now.Sub(e.lastAlert) >= t.cooldown {
			delete(t.entries, port)
		}
	}
}
