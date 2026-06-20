package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.IntervalSeconds != 10 {
		t.Errorf("expected interval 10, got %d", cfg.IntervalSeconds)
	}
	if len(cfg.AllowedPorts) == 0 {
		t.Error("expected non-empty default allowed ports")
	}
	if cfg.AlertOnClose {
		t.Error("expected AlertOnClose to default to false")
	}
}

func TestInterval(t *testing.T) {
	cfg := &Config{IntervalSeconds: 5}
	if cfg.Interval() != 5*time.Second {
		t.Errorf("expected 5s, got %v", cfg.Interval())
	}
}

func TestAllowedSet(t *testing.T) {
	cfg := &Config{AllowedPorts: []int{22, 80, 443}}
	set := cfg.AllowedSet()
	for _, p := range []int{22, 80, 443} {
		if _, ok := set[p]; !ok {
			t.Errorf("expected port %d in allowed set", p)
		}
	}
	if _, ok := set[8080]; ok {
		t.Error("unexpected port 8080 in allowed set")
	}
}

func TestLoadFromFile_Valid(t *testing.T) {
	content := `{"interval_seconds":30,"allowed_ports":[22,8080],"alert_on_close":true,"log_file":"/tmp/portwatch.log"}`
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	cfg, err := LoadFromFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.IntervalSeconds != 30 {
		t.Errorf("expected 30, got %d", cfg.IntervalSeconds)
	}
	if !cfg.AlertOnClose {
		t.Error("expected AlertOnClose true")
	}
	if len(cfg.AllowedPorts) != 2 {
		t.Errorf("expected 2 allowed ports, got %d", len(cfg.AllowedPorts))
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
