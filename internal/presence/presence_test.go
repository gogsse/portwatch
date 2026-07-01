package presence

import (
	"testing"
	"time"
)

func newFixedTracker(base time.Time) *Tracker {
	t := New()
	t.now = func() time.Time { return base }
	return t
}

func TestObserve_RecordsFirstSeen(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := newFixedTracker(base)

	tr.Observe(8080)
	e, ok := tr.Lookup(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !e.FirstSeen.Equal(base) {
		t.Errorf("FirstSeen = %v, want %v", e.FirstSeen, base)
	}
}

func TestObserve_AdvancesLastSeen(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := newFixedTracker(base)
	tr.Observe(9000)

	later := base.Add(5 * time.Minute)
	tr.now = func() time.Time { return later }
	tr.Observe(9000)

	e, _ := tr.Lookup(9000)
	if !e.FirstSeen.Equal(base) {
		t.Errorf("FirstSeen changed unexpectedly: got %v", e.FirstSeen)
	}
	if !e.LastSeen.Equal(later) {
		t.Errorf("LastSeen = %v, want %v", e.LastSeen, later)
	}
}

func TestEvict_RemovesEntry(t *testing.T) {
	base := time.Now()
	tr := newFixedTracker(base)
	tr.Observe(443)
	tr.Evict(443)

	if _, ok := tr.Lookup(443); ok {
		t.Error("expected entry to be removed after Evict")
	}
}

func TestStable_TrueWhenOldEnough(t *testing.T) {
	base := time.Now()
	tr := newFixedTracker(base)
	tr.Observe(22)

	tr.now = func() time.Time { return base.Add(10 * time.Minute) }
	tr.Observe(22)

	if !tr.Stable(22, 5*time.Minute) {
		t.Error("expected port 22 to be stable after 10 minutes")
	}
}

func TestStable_FalseWhenTooNew(t *testing.T) {
	base := time.Now()
	tr := newFixedTracker(base)
	tr.Observe(22)

	if tr.Stable(22, 5*time.Minute) {
		t.Error("expected port 22 to be unstable immediately after first observation")
	}
}

func TestActive_ReturnsAllEntries(t *testing.T) {
	tr := New()
	tr.Observe(80)
	tr.Observe(443)
	tr.Observe(8080)

	active := tr.Active()
	if len(active) != 3 {
		t.Errorf("Active() returned %d entries, want 3", len(active))
	}
}

func TestEntry_Duration(t *testing.T) {
	base := time.Now()
	e := Entry{
		Port:      80,
		FirstSeen: base,
		LastSeen:  base.Add(3 * time.Minute),
	}
	if e.Duration() != 3*time.Minute {
		t.Errorf("Duration() = %v, want 3m", e.Duration())
	}
}
