# backoff

Provides **exponential back-off** with optional jitter for portwatch subsystems
that need to retry transient failures (webhook delivery, scanner I/O errors,
etc.).

## Usage

```go
import "github.com/your-org/portwatch/internal/backoff"

// base=100 ms, max=30 s, factor=2, jitter=true
b := backoff.New(100*time.Millisecond, 30*time.Second, 2.0, true)

for {
    if err := doWork(); err == nil {
        b.Reset()
        break
    }
    time.Sleep(b.Next())
}
```

## Parameters

| Parameter | Description |
|-----------|-------------|
| `base`    | Initial delay for the first retry. |
| `max`     | Upper bound on any single delay. |
| `factor`  | Multiplicative growth rate (must be ≥ 1). |
| `jitter`  | When true, randomises each delay in `[d/2, d]` to avoid thundering-herd. |

## Methods

- **`Next() time.Duration`** — Returns the next delay and increments the
  internal attempt counter.
- **`Reset()`** — Resets the attempt counter so the next call to `Next`
  returns the base delay again.
- **`Attempts() int`** — Returns the number of `Next` calls since the last
  `Reset`.

## Panics

The constructor panics on invalid arguments rather than returning an error,
because these are always programming mistakes caught at start-up.
