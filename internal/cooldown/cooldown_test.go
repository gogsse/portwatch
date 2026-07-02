package cooldown

import (
	"testing"
	"time"
)

func newTestTracker(d time.Duration) *Tracker {
	t := New(d)
	return t
}

func TestNew_PanicsOnZeroDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero duration")
		}
	}()
	New(0)
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	tr := newTestTracker(time.Second)
	if !tr.Allow(8080) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	tr := newTestTracker(time.Second)
	tr.Allow(8080)
	if tr.Allow(8080) {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_CallAfterCooldownPermitted(t *testing.T) {
	fixed := time.Now()
	tr := newTestTracker(time.Second)
	tr.now = func() time.Time { return fixed }
	tr.Allow(8080)

	tr.now = func() time.Time { return fixed.Add(2 * time.Second) }
	if !tr.Allow(8080) {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	tr := newTestTracker(time.Second)
	tr.Allow(80)
	if !tr.Allow(443) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_AllowsImmediateRealert(t *testing.T) {
	tr := newTestTracker(time.Minute)
	tr.Allow(9000)
	tr.Reset(9000)
	if !tr.Allow(9000) {
		t.Fatal("expected Allow after Reset to be permitted")
	}
}

func TestActive_ReturnsPortsInCooldown(t *testing.T) {
	fixed := time.Now()
	tr := newTestTracker(time.Minute)
	tr.now = func() time.Time { return fixed }
	tr.Allow(1234)
	tr.Allow(5678)

	active := tr.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active, got %d", len(active))
	}
}

func TestFlush_RemovesExpiredEntries(t *testing.T) {
	fixed := time.Now()
	tr := newTestTracker(time.Second)
	tr.now = func() time.Time { return fixed }
	tr.Allow(3000)

	tr.now = func() time.Time { return fixed.Add(2 * time.Second) }
	tr.Flush()

	if len(tr.entries) != 0 {
		t.Fatalf("expected empty entries after flush, got %d", len(tr.entries))
	}
}

func TestFlush_PreservesActiveEntries(t *testing.T) {
	fixed := time.Now()
	tr := newTestTracker(time.Minute)
	tr.now = func() time.Time { return fixed }
	tr.Allow(4000)

	tr.now = func() time.Time { return fixed.Add(5 * time.Second) }
	tr.Flush()

	if len(tr.entries) != 1 {
		t.Fatalf("expected 1 entry preserved after flush, got %d", len(tr.entries))
	}
}
