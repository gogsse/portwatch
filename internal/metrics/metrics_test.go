package metrics

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	m := New(nil)
	if m == nil {
		t.Fatal("expected non-nil Metrics")
	}
}

func TestIncScans_Increments(t *testing.T) {
	m := New(nil)
	m.IncScans()
	m.IncScans()
	if got := m.Snapshot().Scans; got != 2 {
		t.Errorf("Scans = %d; want 2", got)
	}
}

func TestIncOpened_Increments(t *testing.T) {
	m := New(nil)
	m.IncOpened(3)
	m.IncOpened(2)
	if got := m.Snapshot().Opened; got != 5 {
		t.Errorf("Opened = %d; want 5", got)
	}
}

func TestIncClosed_Increments(t *testing.T) {
	m := New(nil)
	m.IncClosed(4)
	if got := m.Snapshot().Closed; got != 4 {
		t.Errorf("Closed = %d; want 4", got)
	}
}

func TestIncAlerts_Increments(t *testing.T) {
	m := New(nil)
	m.IncAlerts()
	m.IncAlerts()
	m.IncAlerts()
	if got := m.Snapshot().Alerts; got != 3 {
		t.Errorf("Alerts = %d; want 3", got)
	}
}

func TestIncSuppressed_Increments(t *testing.T) {
	m := New(nil)
	m.IncSuppressed()
	if got := m.Snapshot().Suppressed; got != 1 {
		t.Errorf("Suppressed = %d; want 1", got)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	m := New(nil)
	m.IncScans()
	snap := m.Snapshot()
	m.IncScans()
	if snap.Scans != 1 {
		t.Errorf("snapshot was mutated: Scans = %d; want 1", snap.Scans)
	}
}

func TestPrint_ContainsAllFields(t *testing.T) {
	var buf bytes.Buffer
	m := New(&buf)
	m.IncScans()
	m.IncOpened(2)
	m.IncClosed(1)
	m.IncAlerts()
	m.IncSuppressed()
	m.Print()
	out := buf.String()
	for _, want := range []string{"scans=1", "opened=2", "closed=1", "alerts=1", "suppressed=1", "uptime="} {
		if !strings.Contains(out, want) {
			t.Errorf("Print() output missing %q; got: %s", want, out)
		}
	}
}

func TestPrint_ZeroCountersOnFreshInstance(t *testing.T) {
	var buf bytes.Buffer
	m := New(&buf)
	m.Print()
	out := buf.String()
	for _, want := range []string{"scans=0", "opened=0", "closed=0", "alerts=0", "suppressed=0"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print() output missing %q; got: %s", want, out)
		}
	}
}
