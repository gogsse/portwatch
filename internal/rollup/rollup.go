// Package rollup groups repeated port events into summarised bursts,
// reducing noise when the same ports flap repeatedly within a window.
package rollup

import (
	"sync"
	"time"
)

// Event represents a rolled-up summary for a single port.
type Event struct {
	Port      int
	Kind      string // "opened" | "closed"
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Rollup accumulates port events and emits summaries after a flush window.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]*Event // key: "kind:port"
	now     func() time.Time
}

// New creates a Rollup with the given aggregation window.
func New(window time.Duration) *Rollup {
	return &Rollup{
		window:  window,
		entries: make(map[string]*Event),
		now:     time.Now,
	}
}

func key(kind string, port int) string {
	return kind + ":" + itoa(port)
}

// Record adds a port event to the current accumulation window.
func (r *Rollup) Record(kind string, port int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	k := key(kind, port)
	now := r.now()
	if e, ok := r.entries[k]; ok {
		e.Count++
		e.LastSeen = now
	} else {
		r.entries[k] = &Event{
			Port:      port,
			Kind:      kind,
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
		}
	}
}

// Flush returns all accumulated events and resets internal state.
// Events whose last occurrence is older than the window are also included.
func (r *Rollup) Flush() []Event {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Event, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, *e)
	}
	r.entries = make(map[string]*Event)
	return out
}

// Len returns the number of distinct port/kind pairs currently buffered.
func (r *Rollup) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
