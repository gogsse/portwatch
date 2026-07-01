// Package webhook delivers alert payloads to an HTTP endpoint via POST.
package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload is the JSON body sent to the webhook endpoint on each event.
type Payload struct {
	Port      int       `json:"port"`
	Kind      string    `json:"kind"`      // "opened" | "closed"
	Severity  string    `json:"severity"`  // "info" | "warning" | "danger"
	Timestamp time.Time `json:"timestamp"`
}

// Webhook posts alert payloads to a configured HTTP endpoint.
type Webhook struct {
	url    string
	client *http.Client
}

// New returns a Webhook that will POST to url with the given per-request timeout.
func New(url string, timeout time.Duration) *Webhook {
	return &Webhook{
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

// Send serialises p as JSON and POSTs it to the configured webhook URL.
// Any non-2xx status code is returned as an error.
func (w *Webhook) Send(p Payload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: server returned %d", resp.StatusCode)
	}
	return nil
}

// SendOpened fires an "opened" event for the given port with the provided severity.
func (w *Webhook) SendOpened(port int, severity string) error {
	return w.Send(Payload{
		Port:      port,
		Kind:      "opened",
		Severity:  severity,
		Timestamp: time.Now().UTC(),
	})
}

// SendClosed fires a "closed" event for the given port (always info severity).
func (w *Webhook) SendClosed(port int) error {
	return w.Send(Payload{
		Port:      port,
		Kind:      "closed",
		Severity:  "info",
		Timestamp: time.Now().UTC(),
	})
}
