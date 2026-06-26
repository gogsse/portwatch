package dedupe

import (
	"testing"
	"time"
)

func newTestDeduplicator(window time.Duration) (*Deduplicator, *time.Time) {
	now := time.Now()
	d := New(window)
	d.nowFunc = func() time.Time { return now }
	return d, &now
}

func TestIsDuplicate_FirstCallNotDuplicate(t *testing.T) {
	d, _ := newTestDeduplicator(5 * time.Minute)
	if d.IsDuplicate(8080, "opened") {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinWindowIsDuplicate(t *testing.T) {
	d, _ := newTestDeduplicator(5 * time.Minute)
	d.IsDuplicate(8080, "opened")
	if !d.IsDuplicate(8080, "opened") {
		t.Fatal("expected second call within window to be a duplicate")
	}
}

func TestIsDuplicate_AfterWindowExpiry_NotDuplicate(t *testing.T) {
	now := time.Now()
	d := New(1 * time.Minute)
	d.nowFunc = func() time.Time { return now }

	d.IsDuplicate(9090, "opened")

	// Advance time beyond the window.
	now = now.Add(2 * time.Minute)
	d.nowFunc = func() time.Time { return now }

	if d.IsDuplicate(9090, "opened") {
		t.Fatal("expected event after window expiry to not be a duplicate")
	}
}

func TestIsDuplicate_DifferentKinds_Independent(t *testing.T) {
	d, _ := newTestDeduplicator(5 * time.Minute)
	d.IsDuplicate(443, "opened")
	if d.IsDuplicate(443, "closed") {
		t.Fatal("opened and closed events should be tracked independently")
	}
}

func TestIsDuplicate_DifferentPorts_Independent(t *testing.T) {
	d, _ := newTestDeduplicator(5 * time.Minute)
	d.IsDuplicate(80, "opened")
	if d.IsDuplicate(443, "opened") {
		t.Fatal("different ports should be tracked independently")
	}
}

func TestFlush_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	d := New(1 * time.Minute)
	d.nowFunc = func() time.Time { return now }

	d.IsDuplicate(22, "opened")
	d.IsDuplicate(80, "opened")

	if d.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", d.Len())
	}

	now = now.Add(2 * time.Minute)
	d.nowFunc = func() time.Time { return now }
	d.Flush()

	if d.Len() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", d.Len())
	}
}

func TestFlush_RetainsActiveEntries(t *testing.T) {
	now := time.Now()
	d := New(5 * time.Minute)
	d.nowFunc = func() time.Time { return now }

	d.IsDuplicate(3000, "opened")
	d.Flush()

	if d.Len() != 1 {
		t.Fatalf("expected 1 active entry after flush, got %d", d.Len())
	}
}
