// Package audit provides structured event logging for portwatch,
// recording port lifecycle events with timestamps for forensic review.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventOpened    EventKind = "opened"
	EventClosed    EventKind = "closed"
	EventSuppressed EventKind = "suppressed"
	EventAllowed   EventKind = "allowed"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Kind      EventKind `json:"kind"`
	Note      string    `json:"note,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New returns a Logger that writes to out.
// If out is nil, os.Stdout is used.
func New(out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{out: out}
}

// Record writes a single audit entry to the underlying writer.
func (l *Logger) Record(port int, kind EventKind, note string) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Port:      port,
		Kind:      kind,
		Note:      note,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err = fmt.Fprintf(l.out, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// RecordOpened is a convenience wrapper for EventOpened.
func (l *Logger) RecordOpened(port int, note string) error {
	return l.Record(port, EventOpened, note)
}

// RecordClosed is a convenience wrapper for EventClosed.
func (l *Logger) RecordClosed(port int, note string) error {
	return l.Record(port, EventClosed, note)
}
