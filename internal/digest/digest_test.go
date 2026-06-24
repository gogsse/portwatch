package digest

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestRecord_AddsEntry(t *testing.T) {
	d := New(nil)
	d.Record(8080, "opened", epoch)
	if d.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", d.Len())
	}
}

func TestRecord_AccumulatesCount(t *testing.T) {
	d := New(nil)
	d.Record(8080, "opened", epoch)
	d.Record(8080, "opened", epoch.Add(time.Minute))

	key := "8080:opened"
	e := d.events[key]
	if e == nil {
		t.Fatal("entry missing")
	}
	if e.Count != 2 {
		t.Errorf("expected count 2, got %d", e.Count)
	}
	if !e.LastSeen.Equal(epoch.Add(time.Minute)) {
		t.Errorf("LastSeen not updated")
	}
}

func TestFlush_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	d := New(&buf)
	d.Record(443, "opened", epoch)
	d.Record(22, "closed", epoch)
	d.Flush("test-run")

	out := buf.String()
	if !strings.Contains(out, "digest [test-run]") {
		t.Errorf("missing digest header, got: %s", out)
	}
	if !strings.Contains(out, "443") {
		t.Errorf("missing port 443")
	}
	if !strings.Contains(out, "22") {
		t.Errorf("missing port 22")
	}
}

func TestFlush_ResetsBuffer(t *testing.T) {
	var buf bytes.Buffer
	d := New(&buf)
	d.Record(9090, "opened", epoch)
	d.Flush("first")
	buf.Reset()
	d.Flush("second")

	if buf.Len() != 0 {
		t.Errorf("expected empty output after reset flush, got: %s", buf.String())
	}
}

func TestFlush_EmptyBuffer_NoWrite(t *testing.T) {
	var buf bytes.Buffer
	d := New(&buf)
	d.Flush("empty")
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty digest")
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	d := New(nil)
	if d.out == nil {
		t.Error("expected non-nil writer")
	}
}
