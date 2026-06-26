package correlation

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestAdd_AccumulatesWithinWindow(t *testing.T) {
	c := New(5 * time.Second)
	e1 := Event{Port: 8080, Kind: Opened, ObservedAt: t0}
	e2 := Event{Port: 8081, Kind: Opened, ObservedAt: t0.Add(2 * time.Second)}

	if burst := c.Add(e1); burst != nil {
		t.Fatal("expected nil burst during window")
	}
	if burst := c.Add(e2); burst != nil {
		t.Fatal("expected nil burst during window")
	}

	burst := c.Flush()
	if burst == nil {
		t.Fatal("expected burst from Flush")
	}
	if len(burst.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(burst.Events))
	}
}

func TestAdd_TriggersNewBurstAfterWindow(t *testing.T) {
	c := New(3 * time.Second)
	c.Add(Event{Port: 80, Kind: Opened, ObservedAt: t0})

	late := Event{Port: 443, Kind: Opened, ObservedAt: t0.Add(10 * time.Second)}
	burst := c.Add(late)
	if burst == nil {
		t.Fatal("expected expired burst to be returned")
	}
	if len(burst.Events) != 1 || burst.Events[0].Port != 80 {
		t.Fatalf("unexpected burst contents: %+v", burst.Events)
	}
}

func TestFlush_ReturnsNilWhenEmpty(t *testing.T) {
	c := New(time.Second)
	if b := c.Flush(); b != nil {
		t.Fatal("expected nil for empty correlator")
	}
}

func TestFlush_StartsAtFirstEvent(t *testing.T) {
	c := New(10 * time.Second)
	c.Add(Event{Port: 22, Kind: Opened, ObservedAt: t0})
	c.Add(Event{Port: 23, Kind: Opened, ObservedAt: t0.Add(time.Second)})

	b := c.Flush()
	if !b.StartedAt.Equal(t0) {
		t.Fatalf("expected StartedAt %v, got %v", t0, b.StartedAt)
	}
}
