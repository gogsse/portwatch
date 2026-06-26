package rollup

import (
	"testing"
	"time"
)

func TestRecord_AddsNewEntry(t *testing.T) {
	r := New(time.Minute)
	r.Record("opened", 8080)

	if r.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Len())
	}
}

func TestRecord_AccumulatesCount(t *testing.T) {
	r := New(time.Minute)
	r.Record("opened", 8080)
	r.Record("opened", 8080)
	r.Record("opened", 8080)

	events := r.Flush()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Count != 3 {
		t.Errorf("expected count 3, got %d", events[0].Count)
	}
}

func TestRecord_DifferentKindsTreatedSeparately(t *testing.T) {
	r := New(time.Minute)
	r.Record("opened", 9000)
	r.Record("closed", 9000)

	if r.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", r.Len())
	}
}

func TestRecord_DifferentPortsTreatedSeparately(t *testing.T) {
	r := New(time.Minute)
	r.Record("opened", 80)
	r.Record("opened", 443)

	if r.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", r.Len())
	}
}

func TestFlush_ResetsBuffer(t *testing.T) {
	r := New(time.Minute)
	r.Record("opened", 8080)
	r.Flush()

	if r.Len() != 0 {
		t.Errorf("expected empty buffer after flush, got %d", r.Len())
	}
}

func TestFlush_EmptyBuffer_ReturnsEmpty(t *testing.T) {
	r := New(time.Minute)
	events := r.Flush()
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestRecord_SetsTimestamps(t *testing.T) {
	before := time.Now()
	r := New(time.Minute)
	r.Record("opened", 3000)
	after := time.Now()

	events := r.Flush()
	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	e := events[0]
	if e.FirstSeen.Before(before) || e.FirstSeen.After(after) {
		t.Errorf("FirstSeen %v out of expected range", e.FirstSeen)
	}
	if e.LastSeen.Before(before) || e.LastSeen.After(after) {
		t.Errorf("LastSeen %v out of expected range", e.LastSeen)
	}
}

func TestLen_ReflectsCurrentState(t *testing.T) {
	r := New(time.Minute)
	if r.Len() != 0 {
		t.Errorf("expected 0, got %d", r.Len())
	}
	r.Record("opened", 1234)
	if r.Len() != 1 {
		t.Errorf("expected 1, got %d", r.Len())
	}
}
