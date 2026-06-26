// Package dedupe provides event deduplication to avoid re-alerting on
// ports that have already been reported within a configurable window.
package dedupe

import (
	"sync"
	"time"
)

// entry tracks when a port event was last seen.
type entry struct {
	seenAt  time.Time
	eventKey string
}

// Deduplicator suppresses duplicate port events within a time window.
type Deduplicator struct {
	mu      sync.Mutex
	window  time.Duration
	seen    map[string]entry
	nowFunc func() time.Time
}

// New returns a Deduplicator with the given deduplication window.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		window:  window,
		seen:    make(map[string]entry),
		nowFunc: time.Now,
	}
}

// IsDuplicate reports whether the given (port, kind) pair has been seen
// within the deduplication window. If not, it records the event and
// returns false.
func (d *Deduplicator) IsDuplicate(port int, kind string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := eventKey(port, kind)
	now := d.nowFunc()

	if e, ok := d.seen[key]; ok {
		if now.Sub(e.seenAt) < d.window {
			return true
		}
	}

	d.seen[key] = entry{seenAt: now, eventKey: key}
	return false
}

// Flush removes all entries whose window has expired.
func (d *Deduplicator) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	for k, e := range d.seen {
		if now.Sub(e.seenAt) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Len returns the number of currently tracked events.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

func eventKey(port int, kind string) string {
	return kind + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
