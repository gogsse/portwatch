package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortState holds a snapshot of open ports at a given time.
type PortState struct {
	Ports     []int
	Timestamp time.Time
}

// Watcher monitors open ports at a regular interval and reports changes.
type Watcher struct {
	Interval  time.Duration
	AlertFunc func(opened, closed []int)
	stopCh    chan struct{}
}

// NewWatcher creates a Watcher with the given polling interval and alert callback.
func NewWatcher(interval time.Duration, alertFunc func(opened, closed []int)) *Watcher {
	return &Watcher{
		Interval:  interval,
		AlertFunc: alertFunc,
		stopCh:    make(chan struct{}),
	}
}

// Start begins the monitoring loop. It blocks until Stop is called.
func (w *Watcher) Start() {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	prev, err := scanner.ScanOpenPorts()
	if err != nil {
		log.Printf("portwatch: initial scan failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			curr, err := scanner.ScanOpenPorts()
			if err != nil {
				log.Printf("portwatch: scan error: %v", err)
				continue
			}
			opened, closed := diff(prev, curr)
			if (len(opened) > 0 || len(closed) > 0) && w.AlertFunc != nil {
				w.AlertFunc(opened, closed)
			}
			prev = curr
		case <-w.stopCh:
			return
		}
	}
}

// Stop signals the monitoring loop to exit.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

// diff returns ports that were opened and ports that were closed
// between the previous and current snapshots.
func diff(prev, curr []int) (opened, closed []int) {
	prevSet := toSet(prev)
	currSet := toSet(curr)

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
