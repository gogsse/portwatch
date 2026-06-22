# history

The `history` package provides a persistent, append-only log of port change events detected by portwatch.

## Overview

Each time the daemon detects a port opening or closing, a timestamped entry is recorded in a JSON-lines file on disk. The history can be queried for recent events or exported for audit purposes.

## Types

### `Entry`

Represents a single recorded event:

```go
type Entry struct {
    Timestamp time.Time `json:"timestamp"`
    Port      int       `json:"port"`
    Event     string    `json:"event"` // "opened" or "closed"
    Proto     string    `json:"proto"` // e.g. "tcp"
}
```

### `History`

Manages loading, appending, and querying entries.

## Usage

```go
h, err := history.New("/var/lib/portwatch/history.jsonl")
if err != nil {
    log.Fatal(err)
}

h.Add(history.Entry{
    Timestamp: time.Now(),
    Port:      8080,
    Event:     "opened",
    Proto:     "tcp",
})

// Retrieve the 20 most recent entries
recent := h.Recent(20)
for _, e := range recent {
    fmt.Printf("%s port %d %s\n", e.Timestamp.Format(time.RFC3339), e.Port, e.Event)
}
```

## Persistence

Entries are appended to the file in JSON-lines format (one JSON object per line). The file is created automatically if it does not exist. On startup, existing entries are loaded into memory.

## Notes

- `Recent(n)` returns at most `n` entries ordered newest-first.
- `Entries()` returns a full copy of all in-memory entries.
- The history file grows indefinitely; operators should arrange log rotation if needed.
