// main is the entry point for the portwatch daemon.
// It wires together configuration, baseline, filtering, scanning,
// alerting, and reporting into a running daemon process.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/baseline"
	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/daemon"
	"github.com/yourorg/portwatch/internal/filter"
	"github.com/yourorg/portwatch/internal/monitor"
	"github.com/yourorg/portwatch/internal/report"
	"github.com/yourorg/portwatch/internal/scanner"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("portwatch: %v", err)
	}
}

func run() error {
	// --- Flags ---
	configPath := flag.String("config", "", "path to YAML config file (optional)")
	baselinePath := flag.String("baseline", "portwatch.baseline.json", "path to baseline state file")
	initBaseline := flag.Bool("init", false, "initialise baseline from current open ports and exit")
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *version {
		fmt.Println("portwatch dev")
		return nil
	}

	// --- Config ---
	cfg := config.DefaultConfig()
	if *configPath != "" {
		loaded, err := config.LoadFromFile(*configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		cfg = loaded
	}

	// --- Baseline ---
	bl, err := baseline.Load(*baselinePath)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}

	if *initBaseline {
		return initialiseBaseline(*baselinePath)
	}

	if bl == nil {
		log.Println("portwatch: no baseline found — run with -init to create one, or alerts may be noisy")
		bl = baseline.New(*baselinePath, nil)
	}

	// --- Sub-systems ---
	allowedPorts := cfg.AllowedSet()
	f := filter.New(allowedPorts, cfg.ExemptPrivileged)

	notifier := alert.NewNotifier(os.Stdout)
	reporter := report.NewReporter(os.Stdout)
	watcher := monitor.NewWatcher(bl.Ports)

	scanFn := func() ([]int, error) {
		return scanner.ScanOpenPorts()
	}

	d := daemon.New(daemon.Options{
		Interval: cfg.Interval,
		Scan:     scanFn,
		Watcher:  watcher,
		Filter:   f,
		Notifier: notifier,
		Reporter: reporter,
		Baseline: bl,
	})

	// --- Signal handling ---
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("portwatch: started (interval=%s, baseline=%s)", cfg.Interval, *baselinePath)
	d.Run(ctx)
	log.Println("portwatch: stopped")
	return nil
}

// initialiseBaseline performs a single scan, persists the result as the new
// baseline, and exits. Useful for first-run setup.
func initialiseBaseline(path string) error {
	ports, err := scanner.ScanOpenPorts()
	if err != nil {
		return fmt.Errorf("scan ports: %w", err)
	}
	bl := baseline.New(path, ports)
	if err := bl.Save(); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	log.Printf("portwatch: baseline initialised with %d ports → %s", len(ports), path)
	return nil
}
