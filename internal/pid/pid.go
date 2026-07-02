package pid

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// ErrAlreadyRunning is returned when a PID file is present and the recorded
// process is still alive.
var ErrAlreadyRunning = errors.New("another instance is already running")

// File manages a PID file for the portwatch daemon process.
type File struct {
	path string
}

// New returns a File manager for the given path.
// Nothing is written to disk until Acquire is called.
func New(path string) *File {
	return &File{path: path}
}

// Acquire writes the current process PID to the file.
// If the file already exists and belongs to a live process, ErrAlreadyRunning
// is returned. Stale files (dead PIDs) are silently replaced.
func (f *File) Acquire() error {
	if existing, err := f.readPID(); err == nil {
		if processExists(existing) {
			return fmt.Errorf("%w: pid %d", ErrAlreadyRunning, existing)
		}
		_ = os.Remove(f.path) // stale — safe to overwrite
	}
	return os.WriteFile(f.path, []byte(strconv.Itoa(os.Getpid())+"\n"), 0o644)
}

// Release removes the PID file. It is a no-op when the file is absent.
func (f *File) Release() error {
	err := os.Remove(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// Path returns the configured file path.
func (f *File) Path() string { return f.path }

// readPID reads and parses the PID stored in the file.
func (f *File) readPID() (int, error) {
	b, err := os.ReadFile(f.path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(b)))
}

// processExists reports whether the process with the given PID is alive.
// Signal 0 performs validity checking without delivering a signal.
// EPERM means the process exists but is owned by another user.
func processExists(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}
