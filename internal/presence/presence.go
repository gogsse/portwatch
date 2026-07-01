// Package presence tracks how long a port has been continuously open,
// enabling duration-aware alerting and stability classification.
package presence

import (
	"sync"
	"time"
)

// Entry records when a port was first seen and when it was last observed.
type Entry struct {
	Port      int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Duration returns how long the port has been continuously observed.
func (e Entry) Duration() time.Duration {
	return e.LastSeen.Sub(e.FirstSeen)
}

// Tracker maintains first-seen timestamps for open ports.
type Tracker struct {
	mu      sync.Mutex
	entries map[int]Entry
	now     func() time.Time
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[int]Entry),
		now:     time.Now,
	}
}

// Observe marks a port as seen at the current time.
// If the port is new, its FirstSeen is recorded; otherwise only LastSeen advances.
func (t *Tracker) Observe(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if e, ok := t.entries[port]; ok {
		e.LastSeen = now
		t.entries[port] = e
	} else {
		t.entries[port] = Entry{Port: port, FirstSeen: now, LastSeen: now}
	}
}

// Evict removes a port that is no longer open.
func (t *Tracker) Evict(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Lookup returns the Entry for a port and whether it was found.
func (t *Tracker) Lookup(port int) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	return e, ok
}

// Stable reports whether a port has been continuously open for at least minAge.
func (t *Tracker) Stable(port int, minAge time.Duration) bool {
	e, ok := t.Lookup(port)
	if !ok {
		return false
	}
	return e.Duration() >= minAge
}

// Active returns a snapshot of all currently tracked entries.
func (t *Tracker) Active() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}
