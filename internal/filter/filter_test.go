package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

func TestIsAllowed_InSet(t *testing.T) {
	f := filter.New([]int{80, 443, 8080}, false)
	for _, p := range []int{80, 443, 8080} {
		if !f.IsAllowed(p) {
			t.Errorf("expected port %d to be allowed", p)
		}
	}
}

func TestIsAllowed_NotInSet(t *testing.T) {
	f := filter.New([]int{80, 443}, false)
	if f.IsAllowed(9090) {
		t.Error("expected port 9090 to be unexpected")
	}
}

func TestIsAllowed_PrivilegedExempt(t *testing.T) {
	f := filter.New([]int{}, true)
	if !f.IsAllowed(22) {
		t.Error("expected privileged port 22 to be exempt")
	}
}

func TestIsAllowed_PrivilegedNotExemptWhenDisabled(t *testing.T) {
	f := filter.New([]int{}, false)
	if f.IsAllowed(22) {
		t.Error("expected privileged port 22 to be flagged when exemption is off")
	}
}

func TestUnexpected_ReturnsOnlyUnknown(t *testing.T) {
	f := filter.New([]int{80, 443}, false)
	got := f.Unexpected([]int{80, 443, 9090, 3000})
	if len(got) != 2 {
		t.Fatalf("expected 2 unexpected ports, got %d: %v", len(got), got)
	}
}

func TestUnexpected_EmptyInput(t *testing.T) {
	f := filter.New([]int{80}, false)
	got := f.Unexpected([]int{})
	if got != nil && len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestNewFromStrings_SkipsInvalid(t *testing.T) {
	f := filter.NewFromStrings([]string{"80", "not-a-port", "443", "-1"}, false)
	if !f.IsAllowed(80) || !f.IsAllowed(443) {
		t.Error("valid ports should be allowed")
	}
	if f.IsAllowed(0) {
		t.Error("invalid port 0 should not be allowed")
	}
}

func TestAllowedPorts_Length(t *testing.T) {
	f := filter.New([]int{22, 80, 443}, false)
	ports := f.AllowedPorts()
	if len(ports) != 3 {
		t.Errorf("expected 3 allowed ports, got %d", len(ports))
	}
}
