package monitor

import (
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// Watcher periodically scans open ports and emits alerts on changes.
type Watcher struct {
	interval time.Duration
	notifier *alert.Notifier
	prev     []int
	stop     chan struct{}
}

// NewWatcher creates a Watcher that scans at the given interval.
func NewWatcher(interval time.Duration, notifier *alert.Notifier) *Watcher {
	return &Watcher{
		interval: interval,
		notifier: notifier,
		stop:     make(chan struct{}),
	}
}

// Start begins the watch loop. It blocks until Stop is called.
func (w *Watcher) Start() error {
	ports, err := scanner.ScanOpenPorts()
	if err != nil {
		return err
	}
	w.prev = ports

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.tick(); err != nil {
				return err
			}
		case <-w.stop:
			return nil
		}
	}
}

// Stop signals the watch loop to exit.
// It is safe to call Stop only once; calling it more than once will panic.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) tick() error {
	current, err := scanner.ScanOpenPorts()
	if err != nil {
		return err
	}

	opened, closed := diff(w.prev, current)

	for _, p := range opened {
		if err := w.notifier.NotifyOpened(p); err != nil {
			return err
		}
	}
	for _, p := range closed {
		if err := w.notifier.NotifyClosed(p); err != nil {
			return err
		}
	}

	w.prev = current
	return nil
}

// diff computes the ports that were opened (present in current but not prev)
// and the ports that were closed (present in prev but not current).
func diff(prev, current []int) (opened, closed []int) {
	prevSet := toSet(prev)
	currSet := toSet(current)

	for p := range currSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}
	return
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
