package envelope_test

import (
	"strings"
	"testing"
	"time"

	"github.com/example/portwatch/internal/envelope"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 42, time.UTC)

func fixedNow() func() time.Time {
	return func() time.Time { return fixedTime }
}

func TestWrap_FieldsArePopulated(t *testing.T) {
	b := envelope.New("testhost", fixedNow())
	e := b.Wrap(envelope.KindOpened, 8080, "warning", "test message")

	if e.Kind != envelope.KindOpened {
		t.Errorf("expected kind opened, got %s", e.Kind)
	}
	if e.Port != 8080 {
		t.Errorf("expected port 8080, got %d", e.Port)
	}
	if e.Severity != "warning" {
		t.Errorf("expected severity warning, got %s", e.Severity)
	}
	if e.Host != "testhost" {
		t.Errorf("expected host testhost, got %s", e.Host)
	}
	if e.Message != "test message" {
		t.Errorf("unexpected message: %s", e.Message)
	}
	if !e.Timestamp.Equal(fixedTime) {
		t.Errorf("unexpected timestamp: %v", e.Timestamp)
	}
}

func TestWrap_IDContainsKindAndPort(t *testing.T) {
	b := envelope.New("h", fixedNow())
	e := b.Wrap(envelope.KindClosed, 443, "info", "msg")

	if !strings.HasPrefix(e.ID, "closed-443-") {
		t.Errorf("unexpected ID format: %s", e.ID)
	}
}

func TestWrapOpened_DefaultMessage(t *testing.T) {
	b := envelope.New("h", fixedNow())
	e := b.WrapOpened(9000, "danger")

	if e.Kind != envelope.KindOpened {
		t.Errorf("expected opened, got %s", e.Kind)
	}
	if !strings.Contains(e.Message, "9000") {
		t.Errorf("message should reference port 9000: %s", e.Message)
	}
}

func TestWrapClosed_DefaultMessage(t *testing.T) {
	b := envelope.New("h", fixedNow())
	e := b.WrapClosed(22)

	if e.Kind != envelope.KindClosed {
		t.Errorf("expected closed, got %s", e.Kind)
	}
	if e.Severity != "info" {
		t.Errorf("expected severity info, got %s", e.Severity)
	}
	if !strings.Contains(e.Message, "22") {
		t.Errorf("message should reference port 22: %s", e.Message)
	}
}

func TestNew_NilNowFnDefaultsToTimeNow(t *testing.T) {
	b := envelope.New("h", nil)
	e := b.WrapOpened(80, "info")

	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
