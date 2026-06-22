package history_test

import (
	"testing"

	"github.com/yourusername/portwatch/internal/history"
)

// TestHistoryCycle simulates a realistic sequence of port events across
// multiple daemon ticks and verifies the full round-trip.
func TestHistoryCycle_MultipleTicksRoundTrip(t *testing.T) {
	p := tmpPath(t)

	h, err := history.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	events := []struct {
		port int
		kind history.EventKind
	}{
		{8080, history.EventOpened},
		{9090, history.EventOpened},
		{8080, history.EventClosed},
		{3000, history.EventOpened},
		{9090, history.EventClosed},
	}

	for _, ev := range events {
		if err := h.Add(ev.port, ev.kind); err != nil {
			t.Fatalf("Add(%d, %s): %v", ev.port, ev.kind, err)
		}
	}

	// Reload and verify full history preserved.
	h2, err := history.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}

	if got := len(h2.Entries()); got != len(events) {
		t.Fatalf("expected %d entries, got %d", len(events), got)
	}

	for i, e := range h2.Entries() {
		if e.Port != events[i].port || e.Kind != events[i].kind {
			t.Errorf("entry[%d]: got (%d,%s), want (%d,%s)",
				i, e.Port, e.Kind, events[i].port, events[i].kind)
		}
		if e.Timestamp.IsZero() {
			t.Errorf("entry[%d]: timestamp is zero", i)
		}
	}

	// Recent(3) should return the last three events.
	recent := h2.Recent(3)
	if len(recent) != 3 {
		t.Fatalf("Recent(3): got %d", len(recent))
	}
	if recent[2].Port != 9090 || recent[2].Kind != history.EventClosed {
		t.Errorf("last recent entry unexpected: %+v", recent[2])
	}
}
