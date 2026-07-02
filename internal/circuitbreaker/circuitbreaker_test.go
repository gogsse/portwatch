package circuitbreaker

import (
	"testing"
	"time"
)

func newTestBreaker(threshold int, cooldown time.Duration) (*Breaker, *time.Time) {
	b := New(threshold, cooldown)
	current := time.Now()
	b.now = func() time.Time { return current }
	return b, &current
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b, _ := newTestBreaker(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestOpenAfterThresholdFailures(t *testing.T) {
	b, _ := newTestBreaker(3, time.Second)
	for i := 0; i < 3; i++ {
		_ = b.Allow()
		b.RecordFailure()
	}
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen, got %v", b.State())
	}
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestRejectsWhileOpen(t *testing.T) {
	b, _ := newTestBreaker(1, time.Minute)
	_ = b.Allow()
	b.RecordFailure()
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestTransitionsToHalfOpenAfterCooldown(t *testing.T) {
	b, ts := newTestBreaker(1, time.Second)
	_ = b.Allow()
	b.RecordFailure()

	*ts = ts.Add(2 * time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestRecoveryFromHalfOpen(t *testing.T) {
	b, ts := newTestBreaker(1, time.Second)
	_ = b.Allow()
	b.RecordFailure()
	*ts = ts.Add(2 * time.Second)

	_ = b.Allow()       // half-open probe
	b.RecordSuccess()   // single success closes it

	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed after recovery, got %v", b.State())
	}
}

func TestHalfOpen_FailureReopens(t *testing.T) {
	b, ts := newTestBreaker(1, time.Second)
	_ = b.Allow()
	b.RecordFailure()
	*ts = ts.Add(2 * time.Second)

	_ = b.Allow()
	b.RecordFailure()

	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after half-open failure, got %v", b.State())
	}
}

func TestSuccessResetsClosed(t *testing.T) {
	b, _ := newTestBreaker(3, time.Second)
	_ = b.Allow()
	b.RecordFailure()
	_ = b.Allow()
	b.RecordFailure()
	_ = b.Allow()
	b.RecordSuccess() // resets counter
	_ = b.Allow()
	b.RecordFailure()
	_ = b.Allow()
	b.RecordFailure()
	// only 2 failures after reset — should still be closed
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", b.State())
	}
}
