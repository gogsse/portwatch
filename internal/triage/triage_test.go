package triage

import (
	"testing"
)

func TestClassify_KnownDangerousPort(t *testing.T) {
	c := New(nil)
	if got := c.Classify(4444); got != SeverityCritical {
		t.Errorf("expected CRITICAL for port 4444, got %s", got)
	}
}

func TestClassify_WarningPort(t *testing.T) {
	c := New([]int{8888, 9999})
	if got := c.Classify(8888); got != SeverityWarning {
		t.Errorf("expected WARNING for port 8888, got %s", got)
	}
}

func TestClassify_InfoPort(t *testing.T) {
	c := New(nil)
	if got := c.Classify(12345); got != SeverityInfo {
		t.Errorf("expected INFO for port 12345, got %s", got)
	}
}

func TestClassify_DangerousTakesPrecedenceOverWarning(t *testing.T) {
	// Port 22 is in the dangerous list; adding it as warning should not downgrade it.
	c := New([]int{22})
	if got := c.Classify(22); got != SeverityCritical {
		t.Errorf("expected CRITICAL, got %s", got)
	}
}

func TestSeverityString(t *testing.T) {
	cases := []struct {
		s    Severity
		want string
	}{
		{SeverityInfo, "INFO"},
		{SeverityWarning, "WARNING"},
		{SeverityCritical, "CRITICAL"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestCriticalPorts_IsSorted(t *testing.T) {
	c := New(nil)
	ports := c.CriticalPorts()
	for i := 1; i < len(ports); i++ {
		if ports[i] < ports[i-1] {
			t.Errorf("CriticalPorts not sorted at index %d: %v", i, ports)
		}
	}
}

func TestCriticalPorts_ContainsKnownPorts(t *testing.T) {
	c := New(nil)
	set := make(map[int]struct{})
	for _, p := range c.CriticalPorts() {
		set[p] = struct{}{}
	}
	for _, expected := range []int{22, 3389, 4444} {
		if _, ok := set[expected]; !ok {
			t.Errorf("expected port %d in CriticalPorts", expected)
		}
	}
}
