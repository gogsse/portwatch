package snapshot_test

import (
	"path/filepath"
	"testing"

	"github.com/yourorg/portwatch/internal/snapshot"
)

// TestSnapshotCycle simulates a realistic watch-cycle: save, then load and
// compare against a new scan to detect added/removed ports.
func TestSnapshotCycle_DetectsChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	// First cycle — baseline ports.
	firstScan := []int{22, 80, 443}
	if err := snapshot.New(firstScan).Save(path); err != nil {
		t.Fatalf("first save: %v", err)
	}

	// Second cycle — port 80 closed, 8080 opened.
	secondScan := []int{22, 443, 8080}

	prev, err := snapshot.Load(path)
	if err != nil || prev == nil {
		t.Fatalf("load after first cycle: %v", err)
	}

	prevSet := prev.PortSet()
	currSet := snapshot.New(secondScan).PortSet()

	// Detect opened ports.
	var opened []int
	for p := range currSet {
		if _, existed := prevSet[p]; !existed {
			opened = append(opened, p)
		}
	}

	// Detect closed ports.
	var closed []int
	for p := range prevSet {
		if _, still := currSet[p]; !still {
			closed = append(closed, p)
		}
	}

	if len(opened) != 1 || opened[0] != 8080 {
		t.Errorf("opened: expected [8080], got %v", opened)
	}
	if len(closed) != 1 || closed[0] != 80 {
		t.Errorf("closed: expected [80], got %v", closed)
	}

	// Persist second scan for the next hypothetical cycle.
	if err := snapshot.New(secondScan).Save(path); err != nil {
		t.Fatalf("second save: %v", err)
	}
}
