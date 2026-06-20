package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNotify_FormatsOutput(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	err := n.Notify(Event{
		Level:     LevelAlert,
		Port:      8080,
		Message:   "unexpected listener detected",
		Timestamp: fixedTime,
	})
	if err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", output)
	}
	if !strings.Contains(output, "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", output)
	}
	if !strings.Contains(output, "unexpected listener detected") {
		t.Errorf("expected message in output, got: %s", output)
	}
}

func TestNotifyOpened_WritesAlert(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	if err := n.NotifyOpened(443); err != nil {
		t.Fatalf("NotifyOpened returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ALERT") {
		t.Errorf("expected ALERT level, got: %s", output)
	}
	if !strings.Contains(output, "port=443") {
		t.Errorf("expected port=443, got: %s", output)
	}
}

func TestNotifyClosed_WritesInfo(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	if err := n.NotifyClosed(22); err != nil {
		t.Fatalf("NotifyClosed returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected INFO level, got: %s", output)
	}
	if !strings.Contains(output, "port=22") {
		t.Errorf("expected port=22, got: %s", output)
	}
}

func TestNewNotifier_DefaultsToStdout(t *testing.T) {
	n := NewNotifier(nil)
	if n.out == nil {
		t.Error("expected non-nil writer when nil is passed")
	}
}
