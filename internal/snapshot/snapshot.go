// Package snapshot provides point-in-time captures of open ports
// that can be compared across watch cycles.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a timestamped set of open ports.
type Snapshot struct {
	CapturedAt time.Time `json:"captured_at"`
	Ports      []int     `json:"ports"`
}

// New creates a Snapshot from the given port list, stamped with the current time.
func New(ports []int) *Snapshot {
	copied := make([]int, len(ports))
	copy(copied, ports)
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Ports:      copied,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
// Returns nil, nil when the file does not exist.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// PortSet returns the snapshot's port list as a map for O(1) lookup.
func (s *Snapshot) PortSet() map[int]struct{} {
	set := make(map[int]struct{}, len(s.Ports))
	for _, p := range s.Ports {
		set[p] = struct{}{}
	}
	return set
}
