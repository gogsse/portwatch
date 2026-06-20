# baseline

The `baseline` package manages a **persistent snapshot** of known-good open ports for portwatch.

## Purpose

When portwatch first runs (or when the operator explicitly captures a baseline), the current set of open ports is saved to a JSON file. On subsequent runs the live scan is compared against this baseline to distinguish expected listeners from unexpected ones.

## File format

```json
{
  "ports": [22, 80, 443],
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

## Usage

```go
// Capture a new baseline from the current scan.
b := baseline.New(ports)
_ = b.Save("/var/lib/portwatch/baseline.json")

// Load an existing baseline at startup.
b, err := baseline.Load("/var/lib/portwatch/baseline.json")
if b == nil {
    // No baseline yet — first run.
}

// Check if a port is baselined (O(1)).
set := b.ToSet()
if _, known := set[port]; !known {
    // Unexpected listener!
}
```

## Default path

The default baseline path is configured via `config.Config.BaselinePath` (see `internal/config`).
