// Package portmap maintains a human-readable mapping of well-known port
// numbers to service names, allowing alerts and reports to display
// "ssh (22)" instead of bare port numbers.
package portmap

import "fmt"

// Entry describes a known port.
type Entry struct {
	Port    int
	Service string
	Proto   string // "tcp" or "udp"
}

// Map provides lookups from port number to service metadata.
type Map struct {
	entries map[int]Entry
}

// builtIn is the default set of well-known ports embedded at compile time.
var builtIn = []Entry{
	{22, "ssh", "tcp"},
	{23, "telnet", "tcp"},
	{25, "smtp", "tcp"},
	{53, "dns", "udp"},
	{80, "http", "tcp"},
	{110, "pop3", "tcp"},
	{143, "imap", "tcp"},
	{443, "https", "tcp"},
	{445, "smb", "tcp"},
	{3306, "mysql", "tcp"},
	{3389, "rdp", "tcp"},
	{5432, "postgres", "tcp"},
	{6379, "redis", "tcp"},
	{8080, "http-alt", "tcp"},
	{8443, "https-alt", "tcp"},
	{27017, "mongodb", "tcp"},
}

// New returns a Map pre-populated with the built-in well-known entries.
func New() *Map {
	m := &Map{entries: make(map[int]Entry, len(builtIn))}
	for _, e := range builtIn {
		m.entries[e.Port] = e
	}
	return m
}

// Add registers a custom entry, overwriting any existing entry for that port.
func (m *Map) Add(e Entry) {
	m.entries[e.Port] = e
}

// Lookup returns the Entry for the given port and true when found.
func (m *Map) Lookup(port int) (Entry, bool) {
	e, ok := m.entries[port]
	return e, ok
}

// Label returns a display-friendly label such as "ssh (22)".
// If the port is unknown it returns "unknown (N)".
func (m *Map) Label(port int) string {
	if e, ok := m.entries[port]; ok {
		return fmt.Sprintf("%s (%d)", e.Service, port)
	}
	return fmt.Sprintf("unknown (%d)", port)
}

// Known reports whether the port has a registered entry.
func (m *Map) Known(port int) bool {
	_, ok := m.entries[port]
	return ok
}
