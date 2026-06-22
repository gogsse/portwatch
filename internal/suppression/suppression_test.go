package suppression_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/suppression"
)

func TestIsSuppressed_FalseWhenNotRecorded(t *testing.T) {
	s := suppression.New(5 * time.Second)
	if s.IsSuppressed(8080) {
		t.Error("expected port 8080 not to be suppressed before any record")
	}
}

func TestRecord_SuppressesPort(t *testing.T) {
	s := suppression.New(5 * time.Second)
	s.Record(8080)
	if !s.IsSuppressed(8080) {
		t.Error("expected port 8080 to be suppressed after Record")
	}
}

func TestIsSuppressed_FalseAfterCooldown(t *testing.T) {
	s := suppression.New(10 * time.Millisecond)
	s.Record(9090)
	time.Sleep(20 * time.Millisecond)
	if s.IsSuppressed(9090) {
		t.Error("expected port 9090 to no longer be suppressed after cooldown")
	}
}

func TestFlush_RemovesExpiredEntries(t *testing.T) {
	s := suppression.New(10 * time.Millisecond)
	s.Record(1234)
	s.Record(5678)
	time.Sleep(20 * time.Millisecond)
	s.Flush()

	if got := s.Active(); len(got) != 0 {
		t.Errorf("expected empty active list after flush, got %v", got)
	}
}

func TestActive_ReturnsCurrentlySuppressed(t *testing.T) {
	s := suppression.New(5 * time.Second)
	s.Record(80)
	s.Record(443)

	active := s.Active()
	if len(active) != 2 {
		t.Errorf("expected 2 active suppressed ports, got %d", len(active))
	}
}

func TestRecord_IndependentPorts(t *testing.T) {
	s := suppression.New(5 * time.Second)
	s.Record(22)

	if s.IsSuppressed(80) {
		t.Error("port 80 should not be suppressed when only port 22 was recorded")
	}
}
