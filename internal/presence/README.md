# presence

The `presence` package tracks how long individual ports have been continuously
observed as open. It is used by the daemon to distinguish freshly-opened ports
from long-standing listeners, enabling duration-aware alerting.

## Overview

| Symbol | Description |
|--------|-------------|
| `Tracker` | Maintains per-port first-seen / last-seen timestamps |
| `Entry` | Snapshot of a single port's presence data |
| `New()` | Constructs a zero-value `Tracker` |

## Usage

```go
tr := presence.New()

// Call once per scan tick for every open port.
for _, p := range openPorts {
    tr.Observe(p)
}

// Remove ports that are no longer open.
for _, closed := range closedPorts {
    tr.Evict(closed)
}

// Check stability before raising a persistent-listener alert.
if tr.Stable(port, 10*time.Minute) {
    // port has been open for at least 10 minutes
}
```

## Thread Safety

All methods are safe for concurrent use.
