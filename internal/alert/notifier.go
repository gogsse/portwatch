package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change event to be reported.
type Event struct {
	Level     Level
	Port      int
	Message   string
	Timestamp time.Time
}

// Notifier writes alert events to an output destination.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// Pass nil to default to os.Stdout.
func NewNotifier(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// Notify formats and writes an Event to the output destination.
func (n *Notifier) Notify(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s port=%d msg=%q\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port,
		e.Message,
	)
	return err
}

// NotifyOpened emits an ALERT event for a newly opened port.
func (n *Notifier) NotifyOpened(port int) error {
	return n.Notify(Event{
		Level:   LevelAlert,
		Port:    port,
		Message: "unexpected listener detected",
	})
}

// NotifyClosed emits an INFO event for a port that is no longer open.
func (n *Notifier) NotifyClosed(port int) error {
	return n.Notify(Event{
		Level:   LevelInfo,
		Port:    port,
		Message: "port closed",
	})
}
