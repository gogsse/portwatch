package correlation

import (
	"testing"
	"time"
)

func makeBurst(kinds ...EventKind) *Burst {
	events := make([]Event, len(kinds))
	for i, k := range kinds {
		events[i] = Event{Port: 8000 + i, Kind: k, ObservedAt: time.Now()}
	}
	return &Burst{Events: events, StartedAt: time.Now()}
}

func TestClassify_NilBurst(t *testing.T) {
	if got := Classify(nil, 1); got != ClassNoise {
		t.Fatalf("expected ClassNoise, got %v", got)
	}
}

func TestClassify_BelowThreshold(t *testing.T) {
	b := makeBurst(Opened)
	if got := Classify(b, 3); got != ClassNoise {
		t.Fatalf("expected ClassNoise for small burst, got %v", got)
	}
}

func TestClassify_IntrusionPattern(t *testing.T) {
	b := makeBurst(Opened, Opened, Opened)
	if got := Classify(b, 2); got != ClassIntrusion {
		t.Fatalf("expected ClassIntrusion, got %v", got)
	}
}

func TestClassify_RestartPattern(t *testing.T) {
	b := makeBurst(Opened, Closed, Opened, Closed)
	if got := Classify(b, 2); got != ClassRestart {
		t.Fatalf("expected ClassRestart, got %v", got)
	}
}

func TestClassify_SweepPattern(t *testing.T) {
	b := makeBurst(Closed, Closed, Closed)
	if got := Classify(b, 2); got != ClassSweep {
		t.Fatalf("expected ClassSweep, got %v", got)
	}
}

func TestBurstClass_String(t *testing.T) {
	cases := map[BurstClass]string{
		ClassNoise:     "noise",
		ClassRestart:   "restart",
		ClassIntrusion: "intrusion",
		ClassSweep:     "sweep",
	}
	for cls, want := range cases {
		if got := cls.String(); got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	}
}
