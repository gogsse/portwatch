package pid_test

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"portwatch/internal/pid"
)

func TestAcquire_WritesCurrentPID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "portwatch.pid")
	f := pid.New(path)

	if err := f.Acquire(); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	defer f.Release()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	got, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		t.Fatalf("parse pid: %v", err)
	}
	if got != os.Getpid() {
		t.Errorf("pid = %d, want %d", got, os.Getpid())
	}
}

func TestAcquire_ReturnsErrAlreadyRunning_WhenLive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "portwatch.pid")
	f := pid.New(path)

	if err := f.Acquire(); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	defer f.Release()

	if err := f.Acquire(); !errors.Is(err, pid.ErrAlreadyRunning) {
		t.Errorf("second Acquire = %v, want ErrAlreadyRunning", err)
	}
}

func TestAcquire_OverwritesStalePIDFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "portwatch.pid")
	// Write a PID that almost certainly does not correspond to a live process.
	const deadPID = 999999
	if err := os.WriteFile(path, []byte(strconv.Itoa(deadPID)+"\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	f := pid.New(path)
	if err := f.Acquire(); err != nil {
		t.Errorf("Acquire on stale file: %v", err)
	}
	defer f.Release()
}

func TestRelease_RemovesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "portwatch.pid")
	f := pid.New(path)

	if err := f.Acquire(); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if err := f.Release(); err != nil {
		t.Fatalf("Release: %v", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error("file still exists after Release")
	}
}

func TestRelease_NoopWhenFileMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "portwatch.pid")
	f := pid.New(path)
	if err := f.Release(); err != nil {
		t.Errorf("Release on absent file: %v", err)
	}
}

func TestPath_ReturnsConfiguredValue(t *testing.T) {
	const want = "/var/run/portwatch.pid"
	if got := pid.New(want).Path(); got != want {
		t.Errorf("Path() = %q, want %q", got, want)
	}
}
