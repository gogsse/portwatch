// Package envelope wraps alert payloads with metadata before dispatch.
package envelope

import (
	"fmt"
	"time"
)

// Kind describes the nature of a port event.
type Kind string

const (
	KindOpened Kind = "opened"
	KindClosed Kind = "closed"
)

// Envelope is a decorated alert payload ready for external dispatch.
type Envelope struct {
	ID        string    `json:"id"`
	Kind      Kind      `json:"kind"`
	Port      int       `json:"port"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	Host      string    `json:"host"`
	Timestamp time.Time `json:"timestamp"`
}

// Builder constructs Envelope values with a fixed host and clock.
type Builder struct {
	host  string
	nowFn func() time.Time
}

// New returns a Builder stamping envelopes with the given host.
// Provide a nil nowFn to use time.Now.
func New(host string, nowFn func() time.Time) *Builder {
	if nowFn == nil {
		nowFn = time.Now
	}
	return &Builder{host: host, nowFn: nowFn}
}

// Wrap creates an Envelope for a port event.
func (b *Builder) Wrap(kind Kind, port int, severity, message string) Envelope {
	now := b.nowFn()
	return Envelope{
		ID:        fmt.Sprintf("%s-%d-%d", kind, port, now.UnixNano()),
		Kind:      kind,
		Port:      port,
		Severity:  severity,
		Message:   message,
		Host:      b.host,
		Timestamp: now,
	}
}

// WrapOpened is a convenience wrapper for port-opened events.
func (b *Builder) WrapOpened(port int, severity string) Envelope {
	return b.Wrap(KindOpened, port, severity, fmt.Sprintf("unexpected listener on port %d", port))
}

// WrapClosed is a convenience wrapper for port-closed events.
func (b *Builder) WrapClosed(port int) Envelope {
	return b.Wrap(KindClosed, port, "info", fmt.Sprintf("port %d is no longer listening", port))
}
