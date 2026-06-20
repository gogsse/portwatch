package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewReporter_DefaultsToStdout(t *testing.T) {
	r := NewReporter(nil)
	if r.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestRecord_AddsEvent(t *testing.T) {
	r := NewReporter(&bytes.Buffer{})
	r.Record(8080, "opened", false)
	if len(r.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(r.events))
	}
	e := r.events[0]
	if e.Port != 8080 || e.Action != "opened" || e.Allowed {
		t.Errorf("unexpected event values: %+v", e)
	}
}

func TestFlush_WritesAndClearsBuffer(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Record(443, "opened", true)
	r.Record(9999, "closed", false)

	if err := r.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "port=443") {
		t.Errorf("expected port=443 in output, got: %s", output)
	}
	if !strings.Contains(output, "port=9999") {
		t.Errorf("expected port=9999 in output, got: %s", output)
	}
	if !strings.Contains(output, "ALLOWED") {
		t.Errorf("expected ALLOWED status in output")
	}
	if !strings.Contains(output, "UNEXPECTED") {
		t.Errorf("expected UNEXPECTED status in output")
	}
	if len(r.events) != 0 {
		t.Errorf("expected events buffer to be cleared after Flush")
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	r := NewReporter(&bytes.Buffer{})
	r.Record(22, "opened", true)

	evs := r.Events()
	evs[0].Port = 9999

	if r.events[0].Port == 9999 {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestFlush_EmptyBuffer_NoWrite(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	if err := r.Flush(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty flush, got: %s", buf.String())
	}
}
