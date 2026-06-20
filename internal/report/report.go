package report

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a port state change event for reporting.
type Event struct {
	Timestamp time.Time
	Port      int
	Action    string // "opened" or "closed"
	Allowed   bool
}

// Reporter writes structured port event reports to an output sink.
type Reporter struct {
	out    io.Writer
	events []Event
}

// NewReporter creates a Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w}
}

// Record appends an event to the reporter's history.
func (r *Reporter) Record(port int, action string, allowed bool) {
	r.events = append(r.events, Event{
		Timestamp: time.Now(),
		Port:      port,
		Action:    action,
		Allowed:   allowed,
	})
}

// Flush writes all recorded events to the output and clears the buffer.
func (r *Reporter) Flush() error {
	for _, e := range r.events {
		if err := r.writeEvent(e); err != nil {
			return err
		}
	}
	r.events = r.events[:0]
	return nil
}

// Events returns a copy of all buffered events.
func (r *Reporter) Events() []Event {
	copy := make([]Event, len(r.events))
	for i, e := range r.events {
		copy[i] = e
	}
	return copy
}

func (r *Reporter) writeEvent(e Event) error {
	status := "ALLOWED"
	if !e.Allowed {
		status = "UNEXPECTED"
	}
	_, err := fmt.Fprintf(
		r.out,
		"[%s] port=%d action=%s status=%s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Port,
		e.Action,
		status,
	)
	return err
}
