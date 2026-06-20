package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// IntervalSeconds is how often to scan open ports.
	IntervalSeconds int `json:"interval_seconds"`

	// AllowedPorts is the list of ports considered expected/safe.
	AllowedPorts []int `json:"allowed_ports"`

	// AlertOnClose controls whether closing a previously open port triggers an alert.
	AlertOnClose bool `json:"alert_on_close"`

	// LogFile is the optional path to write alerts to. Empty means stdout.
	LogFile string `json:"log_file"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		IntervalSeconds: 10,
		AllowedPorts:    []int{22, 80, 443},
		AlertOnClose:    false,
		LogFile:         "",
	}
}

// Interval returns the scan interval as a time.Duration.
func (c *Config) Interval() time.Duration {
	return time.Duration(c.IntervalSeconds) * time.Second
}

// AllowedSet returns a set (map) of allowed ports for O(1) lookup.
func (c *Config) AllowedSet() map[int]struct{} {
	set := make(map[int]struct{}, len(c.AllowedPorts))
	for _, p := range c.AllowedPorts {
		set[p] = struct{}{}
	}
	return set
}

// LoadFromFile reads and parses a JSON config file at the given path.
func LoadFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
