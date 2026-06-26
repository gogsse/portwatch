// Package fingerprint identifies processes listening on open ports by
// correlating port numbers with entries from /proc/net/tcp.
package fingerprint

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Info holds process information for a listening port.
type Info struct {
	Port int
	PID  int
	Name string
}

// String returns a human-readable representation of the fingerprint.
func (i Info) String() string {
	if i.Name == "" {
		return fmt.Sprintf("port %d (pid %d)", i.Port, i.PID)
	}
	return fmt.Sprintf("port %d -> %s (pid %d)", i.Port, i.Name, i.PID)
}

// Resolver looks up process information for a given port.
type Resolver struct {
	procRoot string
}

// New returns a Resolver using the real /proc filesystem.
func New() *Resolver {
	return &Resolver{procRoot: "/proc"}
}

// newWithRoot returns a Resolver with a custom proc root (for testing).
func newWithRoot(root string) *Resolver {
	return &Resolver{procRoot: root}
}

// Lookup returns fingerprint Info for the given port, or an error if the
// process cannot be identified.
func (r *Resolver) Lookup(port int) (Info, error) {
	inode, err := r.inodeForPort(port)
	if err != nil {
		return Info{Port: port}, err
	}
	pid, name, err := r.pidForInode(inode)
	if err != nil {
		return Info{Port: port}, err
	}
	return Info{Port: port, PID: pid, Name: name}, nil
}

// inodeForPort scans /proc/net/tcp (IPv4) for the inode bound to port.
func (r *Resolver) inodeForPort(port int) (string, error) {
	data, err := os.ReadFile(filepath.Join(r.procRoot, "net", "tcp"))
	if err != nil {
		return "", fmt.Errorf("fingerprint: read tcp table: %w", err)
	}
	hex := fmt.Sprintf("%04X", port)
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		// local address is field[1]: addr:port in hex
		parts := strings.SplitN(fields[1], ":", 2)
		if len(parts) == 2 && strings.EqualFold(parts[1], hex) {
			return fields[9], nil
		}
	}
	return "", fmt.Errorf("fingerprint: no inode found for port %d", port)
}

// pidForInode walks /proc/<pid>/fd to find a socket matching inode.
func (r *Resolver) pidForInode(inode string) (int, string, error) {
	entries, err := os.ReadDir(r.procRoot)
	if err != nil {
		return 0, "", fmt.Errorf("fingerprint: read proc: %w", err)
	}
	target := fmt.Sprintf("socket:[%s]", inode)
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		fdDir := filepath.Join(r.procRoot, e.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}
		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == target {
				name := r.commForPID(pid)
				return pid, name, nil
			}
		}
	}
	return 0, "", fmt.Errorf("fingerprint: no process found for inode %s", inode)
}

// commForPID reads the comm name for a PID.
func (r *Resolver) commForPID(pid int) string {
	data, err := os.ReadFile(filepath.Join(r.procRoot, strconv.Itoa(pid), "comm"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
