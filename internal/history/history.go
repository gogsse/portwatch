package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventKind describes whether a port was opened or closed.
type EventKind string

const (
	EventOpened EventKind = "opened"
	EventClosed EventKind = "closed"
)

// Entry records a single port-change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Kind      EventKind `json:"kind"`
}

// History holds an ordered log of port-change events backed by a JSON file.
type History struct {
	path    string
	entries []Entry
}

// New creates a History that persists to path, loading any existing entries.
func New(path string) (*History, error) {
	h := &History{path: path}
	if err := h.load(); err != nil {
		return nil, err
	}
	return h, nil
}

// Add appends a new event and immediately persists it.
func (h *History) Add(port int, kind EventKind) error {
	h.entries = append(h.entries, Entry{
		Timestamp: time.Now().UTC(),
		Port:      port,
		Kind:      kind,
	})
	return h.save()
}

// Entries returns a copy of all recorded events.
func (h *History) Entries() []Entry {
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Recent returns up to n most-recent entries.
func (h *History) Recent(n int) []Entry {
	if n >= len(h.entries) {
		return h.Entries()
	}
	out := make([]Entry, n)
	copy(out, h.entries[len(h.entries)-n:])
	return out
}

func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("history: read %s: %w", h.path, err)
	}
	if err := json.Unmarshal(data, &h.entries); err != nil {
		return fmt.Errorf("history: parse %s: %w", h.path, err)
	}
	return nil
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	if err := os.WriteFile(h.path, data, 0o644); err != nil {
		return fmt.Errorf("history: write %s: %w", h.path, err)
	}
	return nil
}
