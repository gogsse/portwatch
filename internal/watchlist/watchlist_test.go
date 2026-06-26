package watchlist_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/watchlist"
)

func TestNew_Empty(t *testing.T) {
	wl := watchlist.New()
	if len(wl.Ports()) != 0 {
		t.Fatal("expected empty watchlist")
	}
}

func TestNewFromSlice_Populated(t *testing.T) {
	wl := watchlist.NewFromSlice([]int{22, 80, 443})
	ports := wl.Ports()
	if len(ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(ports))
	}
	// Ports() must be sorted
	if ports[0] != 22 || ports[1] != 80 || ports[2] != 443 {
		t.Fatalf("unexpected order: %v", ports)
	}
}

func TestContains_TrueAndFalse(t *testing.T) {
	wl := watchlist.NewFromSlice([]int{8080})
	if !wl.Contains(8080) {
		t.Error("expected 8080 to be in watchlist")
	}
	if wl.Contains(9090) {
		t.Error("expected 9090 not to be in watchlist")
	}
}

func TestAdd_Remove(t *testing.T) {
	wl := watchlist.New()
	wl.Add(3306)
	if !wl.Contains(3306) {
		t.Fatal("expected 3306 after Add")
	}
	wl.Remove(3306)
	if wl.Contains(3306) {
		t.Fatal("expected 3306 absent after Remove")
	}
}

func TestFilter_ReturnsIntersection(t *testing.T) {
	wl := watchlist.NewFromSlice([]int{22, 443, 8443})
	got := wl.Filter([]int{80, 443, 8443, 9000})
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d: %v", len(got), got)
	}
}

func TestLoadFromFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watchlist.json")
	data, _ := json.Marshal([]int{22, 3389})
	_ = os.WriteFile(path, data, 0o644)

	wl, err := watchlist.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !wl.Contains(22) || !wl.Contains(3389) {
		t.Error("expected loaded ports to be present")
	}
}

func TestLoadFromFile_Missing_ReturnsEmpty(t *testing.T) {
	wl, err := watchlist.LoadFromFile("/nonexistent/path/watchlist.json")
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if len(wl.Ports()) != 0 {
		t.Fatal("expected empty watchlist for missing file")
	}
}

func TestLoadFromFile_Corrupt_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := watchlist.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
