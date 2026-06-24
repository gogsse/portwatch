// Package ratelimit provides per-port alert rate limiting to prevent
// notification floods when a port repeatedly appears and disappears.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks how recently an alert was emitted for a given port
// and suppresses further alerts until a cooldown period has elapsed.
type Limiter struct {
	mu       sync.Mutex
	last     map[int]time.Time
	cooldown time.Duration
	clock    func() time.Time
}

// New returns a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		last:     make(map[int]time.Time),
		cooldown: cooldown,
		clock:    time.Now,
	}
}

// Allow reports whether an alert for the given port should be emitted.
// If allowed, the port's last-seen timestamp is updated.
func (l *Limiter) Allow(port int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	if t, ok := l.last[port]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[port] = now
	return true
}

// Reset clears the rate-limit record for the given port, allowing the
// next alert to be emitted immediately regardless of cooldown.
func (l *Reset) Reset(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, port)
}

// Flush removes entries whose cooldown has fully elapsed, keeping the
// internal map from growing unboundedly over long daemon runs.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	for port, t := range l.last {
		if now.Sub(t) >= l.cooldown {
			delete(l.last, port)
		}
	}
}

// Active returns the ports that are currently within their cooldown window.
func (l *Limiter) Active() []int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	var ports []int
	for port, t := range l.last {
		if now.Sub(t) < l.cooldown {
			ports = append(ports, port)
		}
	}
	return ports
}
