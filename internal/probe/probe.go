// Package probe attempts lightweight TCP dial checks against open ports
// to confirm they are genuinely accepting connections rather than merely
// appearing in /proc/net/tcp.
package probe

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single port probe.
type Result struct {
	Port      int
	Reachable bool
	Latency   time.Duration
	Err       error
}

// String returns a human-readable summary of the result.
func (r Result) String() string {
	if r.Reachable {
		return fmt.Sprintf("port %d reachable (latency %s)", r.Port, r.Latency.Round(time.Millisecond))
	}
	return fmt.Sprintf("port %d unreachable: %v", r.Port, r.Err)
}

// Prober dials ports to verify they are actively accepting connections.
type Prober struct {
	host    string
	timeout time.Duration
}

// New returns a Prober that connects to host with the given dial timeout.
// If timeout is zero, a 2-second default is used.
func New(host string, timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return &Prober{host: host, timeout: timeout}
}

// Probe dials the given port and returns a Result.
func (p *Prober) Probe(port int) Result {
	addr := fmt.Sprintf("%s:%d", p.host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Port: port, Reachable: false, Latency: latency, Err: err}
	}
	_ = conn.Close()
	return Result{Port: port, Reachable: true, Latency: latency}
}

// ProbeAll probes each port in the slice and returns all results.
func (p *Prober) ProbeAll(ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Probe(port))
	}
	return results
}

// Reachable returns only those results where the port accepted a connection.
func Reachable(results []Result) []Result {
	out := results[:0:0]
	for _, r := range results {
		if r.Reachable {
			out = append(out, r)
		}
	}
	return out
}
