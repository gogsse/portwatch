package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/portwatch/internal/history"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNew_EmptyWhenFileMissing(t *testing.T) {
	h, err := history.New(tmpPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Entries()) != 0 {
		t.Errorf("expected 0 entries, got %d", len(h.Entries()))
	}
}

func TestAdd_PersistsEntry(t *testing.T) {
	p := tmpPath(t)
	h, _ := history.New(p)

	if err := h.Add(8080, history.EventOpened); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Reload from disk.
	h2, err := history.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(h2.Entries()) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(h2.Entries()))
	}
	e := h2.Entries()[0]
	if e.Port != 8080 || e.Kind != history.EventOpened {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	h, _ := history.New(tmpPath(t))
	_ = h.Add(443, history.EventClosed)

	a := h.Entries()
	a[0].Port = 9999
	b := h.Entries()
	if b[0].Port == 9999 {
		t.Error("Entries should return an independent copy")
	}
}

func TestRecent_LimitsResults(t *testing.T) {
	h, _ := history.New(tmpPath(t))
	for _, port := range []int{80, 443, 8080, 9090} {
		_ = h.Add(port, history.EventOpened)
	}

	got := h.Recent(2)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	if got[0].Port != 8080 || got[1].Port != 9090 {
		t.Errorf("unexpected ports in Recent: %v", got)
	}
}

func TestRecent_AllWhenNLarger(t *testing.T) {
	h, _ := history.New(tmpPath(t))
	_ = h.Add(22, history.EventOpened)

	if len(h.Recent(10)) != 1 {
		t.Error("Recent(10) should return all 1 entry")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	_ = os.WriteFile(p, []byte("not-json{"), 0o644)
	_, err := history.New(p)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
