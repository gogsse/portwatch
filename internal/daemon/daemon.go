package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/scanner"
)

// Daemon orchestrates the port monitoring loop.
type Daemon struct {
	cfg      *config.Config
	watcher  *monitor.Watcher
	notifier *alert.Notifier
	reporter *report.Reporter
}

// New creates a new Daemon with the provided configuration.
func New(cfg *config.Config) *Daemon {
	return &Daemon{
		cfg:      cfg,
		watcher:  monitor.NewWatcher(cfg),
		notifier: alert.NewNotifier(nil),
		reporter: report.NewReporter(nil),
	}
}

// Run starts the monitoring loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch daemon starting (interval=%s)\n", d.cfg.Interval)

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	// Perform an initial scan immediately.
	if err := d.tick(); err != nil {
		log.Printf("initial scan error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch daemon stopping")
			d.reporter.Flush()
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("scan error: %v", err)
			}
		}
	}
}

// tick performs a single scan cycle.
func (d *Daemon) tick() error {
	ports, err := scanner.ScanOpenPorts()
	if err != nil {
		return err
	}

	opened, closed := d.watcher.Diff(ports)

	for _, p := range opened {
		d.notifier.NotifyOpened(p)
		d.reporter.Record("opened", p)
	}
	for _, p := range closed {
		d.notifier.NotifyClosed(p)
		d.reporter.Record("closed", p)
	}

	return nil
}
