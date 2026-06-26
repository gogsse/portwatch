// Package heartbeat emits periodic liveness signals so external supervisors
// (systemd watchdog, uptime monitors, etc.) can detect a stalled daemon.
package heartbeat

import (
	"context"
	"io"
	"os"
	"sync/atomic"
	"time"
)

// Beat is written to the output writer on every tick.
const Beat = "."

// Heartbeat emits a periodic signal and exposes a missed-beat counter.
type Heartbeat struct {
	interval time.Duration
	out      io.Writer
	missed   atomic.Int64
	ticks    atomic.Int64
}

// New returns a Heartbeat that writes to w every interval.
// If w is nil, os.Stderr is used.
func New(interval time.Duration, w io.Writer) *Heartbeat {
	if w == nil {
		w = os.Stderr
	}
	return &Heartbeat{
		interval: interval,
		out:      w,
	}
}

// Run blocks until ctx is cancelled, emitting a beat on every tick.
// If a write takes longer than the interval the missed-beat counter is
// incremented so callers can surface the issue through metrics or alerts.
func (h *Heartbeat) Run(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			done := make(chan struct{}, 1)
			go func() {
				_, _ = io.WriteString(h.out, Beat)
				done <- struct{}{}
			}()

			select {
			case <-done:
				h.ticks.Add(1)
			case <-time.After(h.interval):
				h.missed.Add(1)
			}
		}
	}
}

// Ticks returns the total number of successful beats emitted.
func (h *Heartbeat) Ticks() int64 { return h.ticks.Load() }

// Missed returns the number of beats that were not acknowledged within the
// configured interval.
func (h *Heartbeat) Missed() int64 { return h.missed.Load() }
