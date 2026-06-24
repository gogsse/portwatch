// Package triage classifies port events by severity based on well-known
// dangerous ports and configurable rules.
package triage

import "sort"

// Severity represents the urgency level of a port event.
type Severity int

const (
	SeverityInfo     Severity = iota // Expected or low-risk port
	SeverityWarning                  // Uncommon but not immediately dangerous
	SeverityCritical                 // Known dangerous or unexpected privileged port
)

// String returns a human-readable label for the severity.
func (s Severity) String() string {
	switch s {
	case SeverityCritical:
		return "CRITICAL"
	case SeverityWarning:
		return "WARNING"
	default:
		return "INFO"
	}
}

// defaultDangerousPorts is a curated list of ports commonly associated with
// malware, remote-access tools, or high-value attack targets.
var defaultDangerousPorts = []int{
	22, 23, 135, 139, 445, 1080, 3389, 4444, 5900, 6666, 6667, 8080, 9001,
}

// Classifier assigns a Severity to a given port number.
type Classifier struct {
	dangerous map[int]struct{}
	warning   map[int]struct{}
}

// New returns a Classifier seeded with the built-in dangerous port list.
// Additional warning-level ports can be supplied via warningPorts.
func New(warningPorts []int) *Classifier {
	dan := make(map[int]struct{}, len(defaultDangerousPorts))
	for _, p := range defaultDangerousPorts {
		dan[p] = struct{}{}
	}
	warn := make(map[int]struct{}, len(warningPorts))
	for _, p := range warningPorts {
		warn[p] = struct{}{}
	}
	return &Classifier{dangerous: dan, warning: warn}
}

// Classify returns the Severity for the given port.
func (c *Classifier) Classify(port int) Severity {
	if _, ok := c.dangerous[port]; ok {
		return SeverityCritical
	}
	if _, ok := c.warning[port]; ok {
		return SeverityWarning
	}
	return SeverityInfo
}

// CriticalPorts returns a sorted slice of all ports currently marked critical.
func (c *Classifier) CriticalPorts() []int {
	out := make([]int, 0, len(c.dangerous))
	for p := range c.dangerous {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}
