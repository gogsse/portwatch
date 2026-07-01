package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/webhook"
)

func TestNew_ReturnsNonNil(t *testing.T) {
	w := webhook.New("http://example.com/hook", 5*time.Second)
	if w == nil {
		t.Fatal("expected non-nil Webhook")
	}
}

func TestSend_PostsJSON(t *testing.T) {
	var received webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: got %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("content-type: got %q, want application/json", ct)
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	wh := webhook.New(srv.URL, 5*time.Second)
	p := webhook.Payload{Port: 8080, Kind: "opened", Severity: "warning", Timestamp: time.Now().UTC()}
	if err := wh.Send(p); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if received.Port != 8080 {
		t.Errorf("port: got %d, want 8080", received.Port)
	}
	if received.Kind != "opened" {
		t.Errorf("kind: got %q, want \"opened\"", received.Kind)
	}
}

func TestSend_ErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	wh := webhook.New(srv.URL, 5*time.Second)
	if err := wh.Send(webhook.Payload{Port: 443, Kind: "opened", Severity: "info"}); err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestSend_ErrorOnNetworkFailure(t *testing.T) {
	wh := webhook.New("http://127.0.0.1:1", 100*time.Millisecond)
	if err := wh.Send(webhook.Payload{Port: 22, Kind: "opened", Severity: "danger"}); err == nil {
		t.Fatal("expected network error, got nil")
	}
}

func TestSendOpened_SetsKindAndSeverity(t *testing.T) {
	var received webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	wh := webhook.New(srv.URL, 5*time.Second)
	if err := wh.SendOpened(9090, "warning"); err != nil {
		t.Fatalf("SendOpened: %v", err)
	}
	if received.Kind != "opened" {
		t.Errorf("kind: got %q, want \"opened\"", received.Kind)
	}
	if received.Severity != "warning" {
		t.Errorf("severity: got %q, want \"warning\"", received.Severity)
	}
}

func TestSendClosed_SetsKindAndInfoSeverity(t *testing.T) {
	var received webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	wh := webhook.New(srv.URL, 5*time.Second)
	if err := wh.SendClosed(3306); err != nil {
		t.Fatalf("SendClosed: %v", err)
	}
	if received.Kind != "closed" {
		t.Errorf("kind: got %q, want \"closed\"", received.Kind)
	}
	if received.Severity != "info" {
		t.Errorf("severity: got %q, want \"info\"", received.Severity)
	}
}
