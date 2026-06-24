// Package digest produces a periodic summary of port activity across ticks.
package digest

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

// Entry holds a single port-event summary line.
type Entry struct {
	Port      int
	Event     string // "opened" | "closed"
	Count     int
	LastSeen  time.Time
}

// Digest accumulates port events and can flush a human-readable summary.
type Digest struct {
	out     io.Writer
	events  map[string]*Entry // key: "<port>:<event>"
}

// New returns a Digest that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Digest {
	if w == nil {
		w = os.Stdout
	}
	return &Digest{
		out:    w,
		events: make(map[string]*Entry),
	}
}

// Record accumulates an event for the given port.
func (d *Digest) Record(port int, event string, at time.Time) {
	key := fmt.Sprintf("%d:%s", port, event)
	if e, ok := d.events[key]; ok {
		e.Count++
		if at.After(e.LastSeen) {
			e.LastSeen = at
		}
		return
	}
	d.events[key] = &Entry{
		Port:     port,
		Event:    event,
		Count:    1,
		LastSeen: at,
	}
}

// Flush writes the accumulated digest to the writer and resets state.
// Nothing is written when there are no recorded events.
func (d *Digest) Flush(label string) {
	if len(d.events) == 0 {
		return
	}

	entries := make([]*Entry, 0, len(d.events))
	for _, e := range d.events {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Port != entries[j].Port {
			return entries[i].Port < entries[j].Port
		}
		return entries[i].Event < entries[j].Event
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== digest [%s] ===\n", label))
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("  port %-6d  %-7s  x%d  last=%s\n",
			e.Port, e.Event, e.Count, e.LastSeen.Format(time.RFC3339)))
	}
	fmt.Fprint(d.out, sb.String())

	d.events = make(map[string]*Entry)
}

// Len returns the number of distinct port+event pairs currently buffered.
func (d *Digest) Len() int { return len(d.events) }
