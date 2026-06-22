// Package suppression provides a cooldown mechanism to prevent repeated
// alerts for the same port within a configurable time window.
package suppression

import (
	"sync"
	"time"
)

// Record tracks the last alert time for a port.
type Record struct {
	Port      int
	LastAlert time.Time
}

// Store holds suppression state for observed ports.
type Store struct {
	mu      sync.Mutex
	cooldown time.Duration
	entries  map[int]time.Time
}

// New creates a Store with the given cooldown duration.
func New(cooldown time.Duration) *Store {
	return &Store{
		cooldown: cooldown,
		entries:  make(map[int]time.Time),
	}
}

// IsSuppressed returns true if an alert for port was already issued
// within the cooldown window.
func (s *Store) IsSuppressed(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	last, ok := s.entries[port]
	if !ok {
		return false
	}
	return time.Since(last) < s.cooldown
}

// Record marks port as alerted at the current time.
func (s *Store) Record(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = time.Now()
}

// Flush removes suppression records for ports that have aged out.
func (s *Store) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for port, last := range s.entries {
		if time.Since(last) >= s.cooldown {
			delete(s.entries, port)
		}
	}
}

// Active returns a snapshot of all currently suppressed ports.
func (s *Store) Active() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]int, 0, len(s.entries))
	for port, last := range s.entries {
		if time.Since(last) < s.cooldown {
			out = append(out, port)
		}
	}
	return out
}
