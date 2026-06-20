# report

The `report` package provides structured, buffered event reporting for port state changes detected by portwatch.

## Overview

A `Reporter` collects port events (opened/closed) and writes them to any `io.Writer` (default: stdout).

Events are buffered via `Record` and flushed together with `Flush`, making it easy to batch-write reports at the end of each scan cycle.

## Usage

```go
import "portwatch/internal/report"

r := report.NewReporter(os.Stdout)

// Record events during a scan cycle
r.Record(8080, "opened", false) // unexpected listener
r.Record(443,  "opened", true)  // allowed listener

// Flush all buffered events at end of cycle
if err := r.Flush(); err != nil {
    log.Printf("report flush error: %v", err)
}
```

## Output Format

```
[2024-01-15T10:30:00Z] port=8080 action=opened status=UNEXPECTED
[2024-01-15T10:30:00Z] port=443  action=opened status=ALLOWED
```

## Types

- `Event` — timestamped record of a port action with allowed/unexpected status
- `Reporter` — buffers and flushes events to a configurable writer
