package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

// TestFilter_WithRealisticPortSet simulates a typical portwatch cycle:
// the scanner returns a set of open ports, and the filter identifies
// which ones are not in the configured allow-list.
func TestFilter_WithRealisticPortSet(t *testing.T) {
	allowed := []int{22, 80, 443, 5432}
	open := []int{22, 80, 443, 5432, 6379, 27017} // redis + mongo unexpected

	f := filter.New(allowed, false)
	unexpected := f.Unexpected(open)

	if len(unexpected) != 2 {
		t.Fatalf("expected 2 unexpected ports, got %d: %v", len(unexpected), unexpected)
	}

	expectedPorts := map[int]bool{6379: true, 27017: true}
	for _, p := range unexpected {
		if !expectedPorts[p] {
			t.Errorf("port %d should not be in unexpected list", p)
		}
	}
}

// TestFilter_PrivilegedExemptionReducesNoise ensures that enabling
// ignorePrivileged prevents low-numbered ports from generating alerts.
func TestFilter_PrivilegedExemptionReducesNoise(t *testing.T) {
	// No ports explicitly allowed, but privileged ports are exempt.
	f := filter.New([]int{}, true)
	open := []int{22, 80, 443, 8080, 9200}

	unexpected := f.Unexpected(open)

	// Only 8080 and 9200 should be flagged.
	if len(unexpected) != 2 {
		t.Fatalf("expected 2 unexpected ports, got %d: %v", len(unexpected), unexpected)
	}
	for _, p := range unexpected {
		if p < 1024 {
			t.Errorf("privileged port %d should be exempt", p)
		}
	}
}
