package scanner

import (
	"testing"
)

func TestParseHexPort_ValidInput(t *testing.T) {
	tests := []struct {
		addr     string
		wantPort int
	}{
		{"00000000:0050", 80},
		{"00000000:01BB", 443},
		{"00000000:0016", 22},
		{"00000000:1F90", 8080},
	}

	for _, tc := range tests {
		t.Run(tc.addr, func(t *testing.T) {
			got, err := parseHexPort(tc.addr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantPort {
				t.Errorf("parseHexPort(%q) = %d; want %d", tc.addr, got, tc.wantPort)
			}
		})
	}
}

func TestParseHexPort_InvalidInput(t *testing.T) {
	invalidAddrs := []string{
		"",
		"nocolon",
		"addr:ZZZZ",
	}

	for _, addr := range invalidAddrs {
		t.Run(addr, func(t *testing.T) {
			_, err := parseHexPort(addr)
			if err == nil {
				t.Errorf("expected error for input %q, got nil", addr)
			}
		})
	}
}

func TestScanOpenPorts_ReturnsSlice(t *testing.T) {
	// This test exercises the real scanner on the host; it should not panic
	// and must return a non-error result on a Linux system.
	entries, err := ScanOpenPorts()
	if err != nil {
		t.Fatalf("ScanOpenPorts() returned error: %v", err)
	}
	// Entries may be empty in sandboxed environments; just verify the type.
	for _, e := range entries {
		if e.Port <= 0 || e.Port > 65535 {
			t.Errorf("invalid port number %d in entry %+v", e.Port, e)
		}
		if e.State != "LISTEN" {
			t.Errorf("expected state LISTEN, got %q", e.State)
		}
	}
}
