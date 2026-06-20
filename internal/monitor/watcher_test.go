package monitor

import (
	"testing"
)

func TestDiff_NewPortsOpened(t *testing.T) {
	prev := []int{80, 443}
	curr := []int{80, 443, 8080}

	opened, closed := diff(prev, curr)

	if len(opened) != 1 || opened[0] != 8080 {
		t.Errorf("expected opened=[8080], got %v", opened)
	}
	if len(closed) != 0 {
		t.Errorf("expected no closed ports, got %v", closed)
	}
}

func TestDiff_PortsClosed(t *testing.T) {
	prev := []int{80, 443, 8080}
	curr := []int{80, 443}

	opened, closed := diff(prev, curr)

	if len(opened) != 0 {
		t.Errorf("expected no opened ports, got %v", opened)
	}
	if len(closed) != 1 || closed[0] != 8080 {
		t.Errorf("expected closed=[8080], got %v", closed)
	}
}

func TestDiff_NoChange(t *testing.T) {
	prev := []int{22, 80, 443}
	curr := []int{22, 80, 443}

	opened, closed := diff(prev, curr)

	if len(opened) != 0 {
		t.Errorf("expected no opened ports, got %v", opened)
	}
	if len(closed) != 0 {
		t.Errorf("expected no closed ports, got %v", closed)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	prev := []int{}
	curr := []int{80, 443}

	opened, closed := diff(prev, curr)

	if len(opened) != 2 {
		t.Errorf("expected 2 opened ports, got %v", opened)
	}
	if len(closed) != 0 {
		t.Errorf("expected no closed ports, got %v", closed)
	}
}

func TestToSet(t *testing.T) {
	ports := []int{22, 80, 443}
	s := toSet(ports)

	for _, p := range ports {
		if !s[p] {
			t.Errorf("expected port %d in set", p)
		}
	}
	if s[8080] {
		t.Error("did not expect port 8080 in set")
	}
}

func TestNewWatcher_DefaultFields(t *testing.T) {
	alertCalled := false
	w := NewWatcher(0, func(opened, closed []int) {
		alertCalled = true
	})

	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
	if w.AlertFunc == nil {
		t.Error("expected AlertFunc to be set")
	}
	_ = alertCalled
}
