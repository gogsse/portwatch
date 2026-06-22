# internal/history

Persistent append-only log of port-change events observed by the daemon.

## Overview

`History` records every port **opened** or **closed** event to a JSON file so
that operators can audit activity after the fact without relying on live output.

## Types

### `EventKind`

String enum — either `"opened"` or `"closed"`.

### `Entry`

```go
type Entry struct {
    Timestamp time.Time `json:"timestamp"`
    Port      int       `json:"port"`
    Kind      EventKind `json:"kind"`
}
```

### `History`

| Method | Description |
|---|---|
| `New(path)` | Load existing history from *path* (creates empty if missing). |
| `Add(port, kind)` | Append an event and persist to disk immediately. |
| `Entries()` | Return an independent copy of all events. |
| `Recent(n)` | Return the *n* most-recent events. |

## Usage

```go
h, err := history.New("/var/lib/portwatch/history.json")
if err != nil { log.Fatal(err) }

// called by the daemon each tick
for _, p := range newPorts {
    _ = h.Add(p, history.EventOpened)
}
```

## Persistence

Entries are stored as a JSON array and written atomically via `os.WriteFile`.
The file is created automatically on the first `Add` call.
