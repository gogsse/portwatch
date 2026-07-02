package backoff_test

import (
	"testing"
	"time"

	"github.com/your-org/portwatch/internal/backoff"
)

func TestNew_PanicsOnZeroBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero base")
		}
	}()
	backoff.New(0, time.Second, 2, false)
}

func TestNew_PanicsWhenMaxLessThanBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when max < base")
		}
	}()
	backoff.New(time.Second, time.Millisecond, 2, false)
}

func TestNew_PanicsOnFactorBelowOne(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for factor < 1")
		}
	}()
	backoff.New(time.Millisecond, time.Second, 0.5, false)
}

func TestNext_IncreasesDuration(t *testing.T) {
	b := backoff.New(10*time.Millisecond, time.Second, 2, false)

	d1 := b.Next()
	d2 := b.Next()
	d3 := b.Next()

	if d2 <= d1 {
		t.Errorf("expected d2 (%v) > d1 (%v)", d2, d1)
	}
	if d3 <= d2 {
		t.Errorf("expected d3 (%v) > d2 (%v)", d3, d2)
	}
}

func TestNext_CapsAtMax(t *testing.T) {
	max := 50 * time.Millisecond
	b := backoff.New(10*time.Millisecond, max, 10, false)

	for i := 0; i < 10; i++ {
		if d := b.Next(); d > max {
			t.Fatalf("duration %v exceeds max %v on attempt %d", d, max, i)
		}
	}
}

func TestReset_ZeroesAttempts(t *testing.T) {
	b := backoff.New(10*time.Millisecond, time.Second, 2, false)
	b.Next()
	b.Next()

	if b.Attempts() != 2 {
		t.Fatalf("expected 2 attempts, got %d", b.Attempts())
	}

	b.Reset()
	if b.Attempts() != 0 {
		t.Fatalf("expected 0 attempts after reset, got %d", b.Attempts())
	}
}

func TestReset_RestoresBaseDelay(t *testing.T) {
	base := 10 * time.Millisecond
	b := backoff.New(base, time.Second, 2, false)

	for i := 0; i < 5; i++ {
		b.Next()
	}
	b.Reset()

	d := b.Next()
	if d != base {
		t.Errorf("expected base delay %v after reset, got %v", base, d)
	}
}

func TestNext_JitterWithinBounds(t *testing.T) {
	base := 100 * time.Millisecond
	b := backoff.New(base, time.Second, 2, true)

	for i := 0; i < 50; i++ {
		d := b.Next()
		if d < 0 || d > time.Second {
			t.Errorf("jittered duration %v out of bounds", d)
		}
	}
}

func TestAttempts_TracksNextCalls(t *testing.T) {
	b := backoff.New(time.Millisecond, time.Second, 2, false)
	for i := 1; i <= 5; i++ {
		b.Next()
		if got := b.Attempts(); got != i {
			t.Errorf("attempt %d: Attempts() = %d", i, got)
		}
	}
}
