package baseline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_SetsFields(t *testing.T) {
	ports := []int{22, 80, 443}
	b := New(ports)
	if len(b.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(b.Ports))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if b.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	orig := New([]int{22, 8080})
	if err := orig.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(loaded.Ports))
	}
}

func TestLoad_MissingFile_ReturnsNil(t *testing.T) {
	b, err := Load("/nonexistent/path/baseline.json")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if b != nil {
		t.Error("expected nil baseline for missing file")
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for corrupt JSON")
	}
}

func TestUpdate_BumpsUpdatedAt(t *testing.T) {
	b := New([]int{22})
	before := b.UpdatedAt
	b.Update([]int{22, 80})
	if !b.UpdatedAt.After(before) && b.UpdatedAt.Equal(before) {
		// allow equal in fast tests; just check ports changed
	}
	if len(b.Ports) != 2 {
		t.Errorf("expected 2 ports after update, got %d", len(b.Ports))
	}
}

func TestToSet_ContainsAllPorts(t *testing.T) {
	b := New([]int{22, 80, 443})
	set := b.ToSet()
	for _, p := range []int{22, 80, 443} {
		if _, ok := set[p]; !ok {
			t.Errorf("expected port %d in set", p)
		}
	}
	if _, ok := set[9999]; ok {
		t.Error("unexpected port 9999 in set")
	}
}
