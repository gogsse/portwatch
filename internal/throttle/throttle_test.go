package throttle_test

import (
	"testing"
	"time"

	"github.com/yourusername/portwatch/internal/throttle"
)

// newTestThrottle returns a Throttle with a controllable clock.
func newTestThrottle(max int, window time.Duration, now func() time.Time) *throttle.Throttle {
	t := throttle.New(max, window)
	t.SetClock(now) // exposed via a test-only setter added below
	return t
}

func TestAllow_UnderLimit(t *testing.T) {
	th := throttle.New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !th.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestAllow_AtLimitBlocks(t *testing.T) {
	th := throttle.New(2, time.Minute)
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("expected Allow()=false when limit reached")
	}
}

func TestAllow_WindowExpiryRestoresCapacity(t *testing.T) {
	var current time.Time
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	current = base

	th := throttle.New(2, 10*time.Second)
	th.SetClock(func() time.Time { return current })

	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("should be blocked at capacity")
	}

	// Advance past the window so earlier tokens expire.
	current = base.Add(11 * time.Second)
	if !th.Allow() {
		t.Fatal("expected Allow()=true after window expiry")
	}
}

func TestRemaining_ReflectsCurrentUsage(t *testing.T) {
	th := throttle.New(5, time.Minute)
	if got := th.Remaining(); got != 5 {
		t.Fatalf("want 5 remaining, got %d", got)
	}
	th.Allow()
	th.Allow()
	if got := th.Remaining(); got != 3 {
		t.Fatalf("want 3 remaining, got %d", got)
	}
}

func TestReset_RestoresFullCapacity(t *testing.T) {
	th := throttle.New(2, time.Minute)
	th.Allow()
	th.Allow()
	th.Reset()
	if got := th.Remaining(); got != 2 {
		t.Fatalf("want 2 after reset, got %d", got)
	}
	if !th.Allow() {
		t.Fatal("expected Allow()=true after reset")
	}
}
