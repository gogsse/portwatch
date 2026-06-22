// Package filter provides port filtering logic based on allowed lists and ignore rules.
package filter

import "strconv"

// Filter decides which ports are considered unexpected given a set of allowed ports.
type Filter struct {
	allowed map[int]struct{}
	ignorePrivileged bool
}

// New creates a Filter from a slice of allowed port numbers.
// If ignorePrivileged is true, ports below 1024 are never flagged.
func New(allowedPorts []int, ignorePrivileged bool) *Filter {
	allowed := make(map[int]struct{}, len(allowedPorts))
	for _, p := range allowedPorts {
		allowed[p] = struct{}{}
	}
	return &Filter{allowed: allowed, ignorePrivileged: ignorePrivileged}
}

// NewFromStrings creates a Filter from string port representations.
// Invalid entries are silently skipped.
func NewFromStrings(ports []string, ignorePrivileged bool) *Filter {
	nums := make([]int, 0, len(ports))
	for _, s := range ports {
		n, err := strconv.Atoi(s)
		if err == nil && n > 0 && n <= 65535 {
			nums = append(nums, n)
		}
	}
	return New(nums, ignorePrivileged)
}

// IsAllowed returns true when port p is in the allowed set or is exempted
// by the ignorePrivileged rule.
func (f *Filter) IsAllowed(p int) bool {
	if f.ignorePrivileged && p < 1024 {
		return true
	}
	_, ok := f.allowed[p]
	return ok
}

// Unexpected returns only those ports from the provided slice that are not allowed.
func (f *Filter) Unexpected(ports []int) []int {
	var out []int
	for _, p := range ports {
		if !f.IsAllowed(p) {
			out = append(out, p)
		}
	}
	return out
}

// AllowedPorts returns a sorted snapshot of the allowed port set.
func (f *Filter) AllowedPorts() []int {
	out := make([]int, 0, len(f.allowed))
	for p := range f.allowed {
		out = append(out, p)
	}
	return out
}
