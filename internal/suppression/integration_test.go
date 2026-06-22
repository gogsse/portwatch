package suppression_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/suppression"
)

// TestSuppressionCycle_SimulatesAlertLoop verifies the typical daemon loop
// behaviour: first encounter alerts, subsequent encounters within the cooldown
// window are suppressed, and ports re-alert after the window expires.
func TestSuppressionCycle_SimulatesAlertLoop(t *testing.T) {
	cooldown := 30 * time.Millisecond
	s := suppression.New(cooldown)

	ports := []int{8080, 8443, 9000}

	// First tick — none suppressed, record all.
	for _, p := range ports {
		if s.IsSuppressed(p) {
			t.Errorf("tick 1: port %d should not be suppressed", p)
		}
		s.Record(p)
	}

	// Second tick — all should be suppressed within cooldown.
	for _, p := range ports {
		if !s.IsSuppressed(p) {
			t.Errorf("tick 2: port %d should be suppressed", p)
		}
	}

	// Wait for cooldown to expire.
	time.Sleep(cooldown + 10*time.Millisecond)

	// Third tick — suppression should have lifted.
	for _, p := range ports {
		if s.IsSuppressed(p) {
			t.Errorf("tick 3: port %d should no longer be suppressed after cooldown", p)
		}
	}

	// Flush should leave active list empty.
	s.Flush()
	if n := len(s.Active()); n != 0 {
		t.Errorf("expected 0 active entries after flush, got %d", n)
	}
}
