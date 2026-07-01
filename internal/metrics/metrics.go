// Package metrics tracks runtime counters for the portwatch daemon.
package metrics

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Counters holds a point-in-time snapshot of accumulated metric values.
type Counters struct {
	Scans      int64
	Opened     int64
	Closed     int64
	Alerts     int64
	Suppressed int64
}

// Metrics is a goroutine-safe counter set for daemon telemetry.
type Metrics struct {
	mu        sync.Mutex
	counters  Counters
	startedAt time.Time
	out       io.Writer
}

// New returns a Metrics instance that writes summaries to out.
// If out is nil, os.Stderr is used.
func New(out io.Writer) *Metrics {
	if out == nil {
		out = os.Stderr
	}
	return &Metrics{
		startedAt: time.Now(),
		out:       out,
	}
}

// IncScans records one completed port scan cycle.
func (m *Metrics) IncScans() {
	m.mu.Lock()
	m.counters.Scans++
	m.mu.Unlock()
}

// IncOpened adds n to the count of newly-opened ports observed.
func (m *Metrics) IncOpened(n int) {
	m.mu.Lock()
	m.counters.Opened += int64(n)
	m.mu.Unlock()
}

// IncClosed adds n to the count of ports observed closing.
func (m *Metrics) IncClosed(n int) {
	m.mu.Lock()
	m.counters.Closed += int64(n)
	m.mu.Unlock()
}

// IncAlerts records one alert that was dispatched.
func (m *Metrics) IncAlerts() {
	m.mu.Lock()
	m.counters.Alerts++
	m.mu.Unlock()
}

// IncSuppressed records one event that was silenced by the suppression layer.
func (m *Metrics) IncSuppressed() {
	m.mu.Lock()
	m.counters.Suppressed++
	m.mu.Unlock()
}

// Snapshot returns a copy of the current counters without resetting them.
func (m *Metrics) Snapshot() Counters {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.counters
}

// Print writes a single-line summary to the configured writer.
func (m *Metrics) Print() {
	c := m.Snapshot()
	uptime := time.Since(m.startedAt).Truncate(time.Second)
	fmt.Fprintf(m.out,
		"metrics uptime=%s scans=%d opened=%d closed=%d alerts=%d suppressed=%d\n",
		uptime, c.Scans, c.Opened, c.Closed, c.Alerts, c.Suppressed)
}
