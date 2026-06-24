package triage_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/triage"
)

// TestTriage_RealisticUnexpectedPorts simulates the daemon discovering a mix
// of unexpected ports and verifies that each is classified correctly.
func TestTriage_RealisticUnexpectedPorts(t *testing.T) {
	c := triage.New([]int{8080, 9090})

	cases := []struct {
		port int
		want triage.Severity
	}{
		{4444, triage.SeverityCritical}, // reverse-shell favourite
		{3389, triage.SeverityCritical}, // RDP
		{8080, triage.SeverityWarning},  // custom warning
		{9090, triage.SeverityWarning},  // custom warning
		{51820, triage.SeverityInfo},    // WireGuard — unusual but not in lists
		{12345, triage.SeverityInfo},    // completely unknown
	}

	for _, tc := range cases {
		got := c.Classify(tc.port)
		if got != tc.want {
			t.Errorf("port %d: got %s, want %s", tc.port, got, tc.want)
		}
	}
}

// TestTriage_EmptyWarningList ensures the classifier still works with no
// custom warning ports provided.
func TestTriage_EmptyWarningList(t *testing.T) {
	c := triage.New(nil)
	if got := c.Classify(9999); got != triage.SeverityInfo {
		t.Errorf("expected INFO for unlisted port, got %s", got)
	}
}
