# suppression

The `suppression` package provides a thread-safe cooldown store that prevents
the daemon from emitting repeated alerts for the same port within a configurable
time window.

## Overview

When an unexpected port is first detected, the caller records it in the store.
On subsequent scan ticks, `IsSuppressed` returns `true` for that port until the
cooldown duration has elapsed, after which the port is eligible to alert again.

## Usage

```go
import "github.com/example/portwatch/internal/suppression"

// Create a store with a 5-minute cooldown.
store := suppression.New(5 * time.Minute)

for _, port := range unexpectedPorts {
    if store.IsSuppressed(port) {
        continue // already alerted recently
    }
    notifier.NotifyOpened(port)
    store.Record(port)
}

// Periodically flush expired entries to free memory.
store.Flush()
```

## API

| Method | Description |
|---|---|
| `New(cooldown)` | Create a new Store with the given window |
| `IsSuppressed(port)` | Returns true if port is within its cooldown |
| `Record(port)` | Mark port as alerted now |
| `Flush()` | Remove entries whose cooldown has expired |
| `Active()` | Return slice of currently suppressed ports |
