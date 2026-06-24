package ratelimit

import (
	"testing"
	"time"
)

func newTestLimiter(cooldown time.Duration) (*Limiter, *time.Time) {
	now := time.Now()
	l := New(cooldown)
	l.clock = func() time.Time { return now }
	return l, &now
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	l, _ := newTestLimiter(5 * time.Second)
	if !l.Allow(8080) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	l, _ := newTestLimiter(5 * time.Second)
	l.Allow(8080)
	if l.Allow(8080) {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_CallAfterCooldownPermitted(t *testing.T) {
	now := time.Now()
	l := New(5 * time.Second)
	l.clock = func() time.Time { return now }
	l.Allow(8080)

	now = now.Add(6 * time.Second)
	if !l.Allow(8080) {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	l, _ := newTestLimiter(5 * time.Second)
	l.Allow(8080)
	if !l.Allow(9090) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestFlush_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	l := New(5 * time.Second)
	l.clock = func() time.Time { return now }
	l.Allow(8080)
	l.Allow(9090)

	now = now.Add(6 * time.Second)
	l.Flush()

	if !l.Allow(8080) {
		t.Fatal("expected port to be allowed after flush + cooldown")
	}
}

func TestActive_ReturnsPortsInCooldown(t *testing.T) {
	now := time.Now()
	l := New(5 * time.Second)
	l.clock = func() time.Time { return now }
	l.Allow(8080)
	l.Allow(9090)

	active := l.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active ports, got %d", len(active))
	}
}

func TestActive_ExcludesExpiredPorts(t *testing.T) {
	now := time.Now()
	l := New(5 * time.Second)
	l.clock = func() time.Time { return now }
	l.Allow(8080)

	now = now.Add(6 * time.Second)
	active := l.Active()
	if len(active) != 0 {
		t.Fatalf("expected 0 active ports after cooldown, got %d", len(active))
	}
}
