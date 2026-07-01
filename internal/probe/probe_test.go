package probe_test

import (
	"net"
	"testing"
	"time"

	"github.com/rnemeth90/portwatch/internal/probe"
)

// startListener opens a random TCP port on localhost and returns it along with
// a closer function.
func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startListener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { _ = ln.Close() }
}

func TestNew_DefaultTimeout(t *testing.T) {
	p := probe.New("", 0)
	if p == nil {
		t.Fatal("expected non-nil Prober")
	}
}

func TestProbe_ReachablePort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := probe.New("127.0.0.1", 500*time.Millisecond)
	r := p.Probe(port)

	if !r.Reachable {
		t.Fatalf("expected port %d to be reachable, got err: %v", port, r.Err)
	}
	if r.Port != port {
		t.Errorf("expected Port=%d, got %d", port, r.Port)
	}
	if r.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbe_UnreachablePort(t *testing.T) {
	// Port 1 is almost certainly not open in test environments.
	p := probe.New("127.0.0.1", 200*time.Millisecond)
	r := p.Probe(1)

	if r.Reachable {
		t.Skip("port 1 unexpectedly open; skipping unreachable test")
	}
	if r.Err == nil {
		t.Error("expected a non-nil error for unreachable port")
	}
}

func TestProbeAll_ReturnsAllResults(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := probe.New("127.0.0.1", 500*time.Millisecond)
	results := p.ProbeAll([]int{port, 1})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestReachable_FiltersResults(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := probe.New("127.0.0.1", 500*time.Millisecond)
	all := p.ProbeAll([]int{port, 1})
	got := probe.Reachable(all)

	for _, r := range got {
		if !r.Reachable {
			t.Errorf("Reachable returned an unreachable result: port %d", r.Port)
		}
	}
}

func TestResult_String_Reachable(t *testing.T) {
	r := probe.Result{Port: 8080, Reachable: true, Latency: 3 * time.Millisecond}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string from Result.String()")
	}
}

func TestResult_String_Unreachable(t *testing.T) {
	r := probe.Result{Port: 22, Reachable: false, Err: net.ErrClosed}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string from Result.String()")
	}
}
