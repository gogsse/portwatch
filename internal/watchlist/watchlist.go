// Package watchlist maintains a set of ports that should always be monitored
// with elevated attention, regardless of the baseline or filter configuration.
package watchlist

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

// Watchlist holds ports that are explicitly tracked for mandatory alerting.
type Watchlist struct {
	mu    sync.RWMutex
	ports map[int]struct{}
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{
		ports: make(map[int]struct{}),
	}
}

// NewFromSlice constructs a Watchlist pre-populated with the given ports.
func NewFromSlice(ports []int) *Watchlist {
	wl := New()
	for _, p := range ports {
		wl.ports[p] = struct{}{}
	}
	return wl
}

// LoadFromFile reads a JSON array of port numbers from path and returns a
// Watchlist. Returns an empty Watchlist if the file does not exist.
func LoadFromFile(path string) (*Watchlist, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(), nil
	}
	if err != nil {
		return nil, err
	}
	var ports []int
	if err := json.Unmarshal(data, &ports); err != nil {
		return nil, err
	}
	return NewFromSlice(ports), nil
}

// Contains reports whether port is on the watchlist.
func (w *Watchlist) Contains(port int) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, ok := w.ports[port]
	return ok
}

// Add inserts port into the watchlist.
func (w *Watchlist) Add(port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ports[port] = struct{}{}
}

// Remove deletes port from the watchlist.
func (w *Watchlist) Remove(port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.ports, port)
}

// Ports returns a sorted slice of all watched port numbers.
func (w *Watchlist) Ports() []int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]int, 0, len(w.ports))
	for p := range w.ports {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}

// Filter returns only those ports from candidates that appear in the watchlist.
func (w *Watchlist) Filter(candidates []int) []int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var matched []int
	for _, p := range candidates {
		if _, ok := w.ports[p]; ok {
			matched = append(matched, p)
		}
	}
	return matched
}
