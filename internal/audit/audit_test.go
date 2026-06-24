package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/portwatch/internal/audit"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestRecord_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Record(8080, audit.EventOpened, "unexpected listener"); err != nil {
		t.Fatalf("Record returned error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if entry.Port != 8080 {
		t.Errorf("expected port 8080, got %d", entry.Port)
	}
	if entry.Kind != audit.EventOpened {
		t.Errorf("expected kind %q, got %q", audit.EventOpened, entry.Kind)
	}
	if entry.Note != "unexpected listener" {
		t.Errorf("unexpected note: %q", entry.Note)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_AppendsNewline(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Record(443, audit.EventClosed, "")

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("expected output to end with newline")
	}
}

func TestRecordOpened_SetsKind(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.RecordOpened(22, "ssh")

	var entry audit.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry)
	if entry.Kind != audit.EventOpened {
		t.Errorf("expected EventOpened, got %q", entry.Kind)
	}
}

func TestRecordClosed_SetsKind(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.RecordClosed(22, "")

	var entry audit.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry)
	if entry.Kind != audit.EventClosed {
		t.Errorf("expected EventClosed, got %q", entry.Kind)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Record(80, audit.EventOpened, "")
	_ = l.Record(443, audit.EventAllowed, "")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 log lines, got %d", len(lines))
	}
}
