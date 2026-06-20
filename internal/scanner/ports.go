package scanner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single open port found on the system.
type PortEntry struct {
	Protocol string
	LocalAddress string
	Port     int
	PID      int
	State    string
}

// ScanOpenPorts reads /proc/net/tcp and /proc/net/tcp6 to list open ports.
func ScanOpenPorts() ([]PortEntry, error) {
	var entries []PortEntry

	for _, proto := range []string{"tcp", "tcp6"} {
		path := fmt.Sprintf("/proc/net/%s", proto)
		file, err := os.Open(path)
		if err != nil {
			// Not all systems expose both files; skip gracefully.
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		// Skip header line.
		scanner.Scan()

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}

			// State "0A" == LISTEN
			state := fields[3]
			if state != "0A" {
				continue
			}

			port, err := parseHexPort(fields[1])
			if err != nil {
				continue
			}

			entries = append(entries, PortEntry{
				Protocol:     proto,
				LocalAddress: fields[1],
				Port:         port,
				State:        "LISTEN",
			})
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
	}

	return entries, nil
}

// parseHexPort extracts the port from a hex address string like "0F02:1F90".
func parseHexPort(addr string) (int, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid address format: %s", addr)
	}
	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return 0, fmt.Errorf("parsing port hex %q: %w", parts[1], err)
	}
	return int(port), nil
}
