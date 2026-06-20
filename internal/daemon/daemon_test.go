package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestNew_ReturnsNonNil(t *testing.T) {
	cfg := config.DefaultConfig()
	d := New(cfg)
	if d == nil {
		t.Fatal("expected non-nil Daemon")
	}
	if d.cfg == nil {
		t.Error("expected cfg to be set")
	}
	if d.watcher == nil {
		t.Error("expected watcher to be set")
	}
	if d.notifier == nil {
		t.Error("expected notifier to be set")
	}
	if d.reporter == nil {
		t.Error("expected reporter to be set")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Interval = 500 * time.Millisecond

	d := New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestRun_TicksWithinInterval(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Interval = 100 * time.Millisecond

	d := New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	start := time.Now()
	d.Run(ctx) //nolint:errcheck
	elapsed := time.Since(start)

	// Should have run for roughly the timeout duration.
	if elapsed < 300*time.Millisecond {
		t.Errorf("daemon exited too early: %s", elapsed)
	}
}
