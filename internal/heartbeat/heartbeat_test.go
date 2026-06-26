package heartbeat_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rjocoleman/portwatch/internal/heartbeat"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	hb := heartbeat.New(time.Second, nil)
	if hb == nil {
		t.Fatal("expected non-nil Heartbeat")
	}
}

func TestRun_EmitsBeats(t *testing.T) {
	var buf bytes.Buffer
	hb := heartbeat.New(20*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()

	hb.Run(ctx)

	got := buf.String()
	if !strings.Contains(got, heartbeat.Beat) {
		t.Fatalf("expected at least one beat, got %q", got)
	}
}

func TestRun_TicksIncrement(t *testing.T) {
	var buf bytes.Buffer
	hb := heartbeat.New(15*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	hb.Run(ctx)

	if hb.Ticks() < 2 {
		t.Fatalf("expected at least 2 ticks, got %d", hb.Ticks())
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	var buf bytes.Buffer
	hb := heartbeat.New(50*time.Millisecond, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	start := time.Now()
	hb.Run(ctx)
	elapsed := time.Since(start)

	if elapsed > 200*time.Millisecond {
		t.Fatalf("Run did not stop promptly after cancel: %v", elapsed)
	}
}

func TestMissed_ZeroOnNormalOperation(t *testing.T) {
	var buf bytes.Buffer
	hb := heartbeat.New(20*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()

	hb.Run(ctx)

	if hb.Missed() != 0 {
		t.Fatalf("expected 0 missed beats on fast writer, got %d", hb.Missed())
	}
}
