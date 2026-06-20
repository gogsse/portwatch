package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Baseline holds a snapshot of known-good open ports.
type Baseline struct {
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New creates a new Baseline from the given port list.
func New(ports []int) *Baseline {
	now := time.Now().UTC()
	return &Baseline{
		Ports:     ports,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Save writes the baseline to the given file path as JSON.
func (b *Baseline) Save(path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a baseline from the given file path.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Update replaces the port list and bumps UpdatedAt.
func (b *Baseline) Update(ports []int) {
	b.Ports = ports
	b.UpdatedAt = time.Now().UTC()
}

// ToSet returns a map[int]struct{} of all baselined ports for O(1) lookup.
func (b *Baseline) ToSet() map[int]struct{} {
	set := make(map[int]struct{}, len(b.Ports))
	for _, p := range b.Ports {
		set[p] = struct{}{}
	}
	return set
}
