package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/snapshot"
)

func TestNew_SetsFieldsAndCopiesPorts(t *testing.T) {
	before := time.Now().UTC()
	ports := []int{80, 443, 8080}
	s := snapshot.New(ports)

	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if s.CapturedAt.Before(before) {
		t.Errorf("CapturedAt %v is before test start %v", s.CapturedAt, before)
	}
	if len(s.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(s.Ports))
	}
	// Mutating the original slice must not affect the snapshot.
	ports[0] = 9999
	if s.Ports[0] == 9999 {
		t.Error("snapshot ports should be an independent copy")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New([]int{22, 80, 443})
	if err := orig.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil loaded snapshot")
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("ports length mismatch: got %d, want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestLoad_MissingFile_ReturnsNil(t *testing.T) {
	s, err := snapshot.Load("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o600)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}

func TestPortSet_ReturnsCorrectMap(t *testing.T) {
	s := snapshot.New([]int{22, 80, 443})
	set := s.PortSet()

	for _, p := range []int{22, 80, 443} {
		if _, ok := set[p]; !ok {
			t.Errorf("expected port %d in set", p)
		}
	}
	if _, ok := set[9999]; ok {
		t.Error("unexpected port 9999 in set")
	}
}
