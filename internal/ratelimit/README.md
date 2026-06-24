# ratelimit

Per-port alert rate limiting for portwatch.

## Purpose

Prevents notification floods when a port oscillates between open and closed
states within a short time window. Each port is individually tracked; once an
alert has been emitted the port enters a cooldown period during which further
alerts are suppressed.

## Usage

```go
import "github.com/yourorg/portwatch/internal/ratelimit"

// Create a limiter with a 30-second cooldown per port.
limiter := ratelimit.New(30 * time.Second)

// In your alert path:
if limiter.Allow(port) {
    notifier.NotifyOpened(port)
}
```

## API

| Function | Description |
|---|---|
| `New(cooldown)` | Returns a new Limiter with the given cooldown duration |
| `Allow(port)` | Returns true if an alert may be emitted; updates timestamp |
| `Reset(port)` | Clears the record for a port so the next alert is immediate |
| `Flush()` | Removes expired entries to prevent unbounded map growth |
| `Active()` | Returns ports currently within their cooldown window |

## Integration

Call `Flush()` periodically (e.g. once per daemon tick) to keep memory usage
stable during long-running sessions. `Active()` is useful for diagnostics and
status reporting.
