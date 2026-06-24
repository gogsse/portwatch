# digest

The `digest` package provides a lightweight event-aggregation buffer that accumulates port-change events across multiple watcher ticks and periodically flushes a human-readable summary.

## Overview

Rather than emitting a separate alert for every individual tick that observes the same port opening or closing, the `Digest` type collects events in memory and writes a consolidated report when `Flush` is called.

## Usage

```go
w := digest.New(os.Stdout) // pass nil to default to stdout

// called inside the daemon tick handler:
for _, p := range newPorts {
    w.Record(p, "opened", time.Now())
}
for _, p := range closedPorts {
    w.Record(p, "closed", time.Now())
}

// called on a slower cadence (e.g. every N ticks or on SIGUSR1):
w.Flush(time.Now().Format(time.RFC3339))
```

## API

| Symbol | Description |
|--------|-------------|
| `New(w io.Writer) *Digest` | Create a new Digest; pass `nil` to use `os.Stdout` |
| `Record(port int, event string, at time.Time)` | Accumulate an event; duplicate port+event pairs increment the counter |
| `Flush(label string)` | Write the summary and reset the buffer |
| `Len() int` | Number of distinct port+event pairs buffered |
