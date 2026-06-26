# watchlist

The `watchlist` package provides a thread-safe set of port numbers that
should receive **mandatory alerting** regardless of the baseline or filter
configuration.

## Purpose

Some ports are so high-risk (e.g. `22`, `3389`, `3306`) that you always want
an alert when they appear, even if they were previously seen or are within an
allowed range defined elsewhere. The `Watchlist` type sits alongside the
`filter` and `triage` packages to give operators explicit, persistent control
over those ports.

## Usage

```go
// Load from a JSON file (array of ints).
wl, err := watchlist.LoadFromFile("/etc/portwatch/watchlist.json")
if err != nil {
    log.Fatal(err)
}

// Or build programmatically.
wl = watchlist.NewFromSlice([]int{22, 3306, 3389})

// Check a single port.
if wl.Contains(port) {
    // always alert
}

// Get watched ports from a candidate slice.
hits := wl.Filter(openPorts)
```

## File format

The watchlist file is a plain JSON array of integers:

```json
[22, 3306, 3389, 5900]
```

If the file does not exist `LoadFromFile` returns an empty `Watchlist`
without error, so the daemon starts cleanly with no watchlist configured.
