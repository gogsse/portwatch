// Package correlation groups related port events into correlated bursts,
// helping distinguish a service restart (many ports changing at once) from
// a genuinely unexpected new listener.
package correlation

import (
	"sync"
	"time"
)

// EventKind distinguishes port-opened from port-closed events.
type EventKind string

const (
	Opened EventKind = "opened"
	Closed EventKind = "closed"
)

// Event is a single port change captured during one scan tick.
type Event struct {
	Port      int
	Kind      EventKind
	ObservedAt time.Time
}

// Burst is a group of events that arrived within the correlation window.
type Burst struct {
	Events    []Event
	StartedAt time.Time
}

// Correlator accumulates events and groups them into bursts.
type Correlator struct {
	mu      sync.Mutex
	window  time.Duration
	pending []Event
	first   time.Time
}

// New returns a Correlator that groups events arriving within window.
func New(window time.Duration) *Correlator {
	return &Correlator{window: window}
}

// Add records a new event. Returns a completed Burst when the window expires,
// otherwise returns nil.
func (c *Correlator) Add(e Event) *Burst {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := e.ObservedAt
	if len(c.pending) == 0 {
		c.first = now
	}

	if now.Sub(c.first) > c.window && len(c.pending) > 0 {
		burst := c.flush()
		c.pending = []Event{e}
		c.first = now
		return burst
	}

	c.pending = append(c.pending, e)
	return nil
}

// Flush drains any remaining pending events into a Burst regardless of window.
// Returns nil if there are no pending events.
func (c *Correlator) Flush() *Burst {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.pending) == 0 {
		return nil
	}
	return c.flush()
}

func (c *Correlator) flush() *Burst {
	b := &Burst{
		Events:    make([]Event, len(c.pending)),
		StartedAt: c.first,
	}
	copy(b.Events, c.pending)
	c.pending = c.pending[:0]
	return b
}
