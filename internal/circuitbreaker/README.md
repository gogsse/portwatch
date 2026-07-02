# circuitbreaker

Provides a thread-safe circuit breaker used to protect downstream sinks
(webhook endpoints, notifiers, audit writers) from repeated failure storms.

## States

| State | Meaning |
|-------|---------|
| `Closed` | Normal operation — calls are allowed through. |
| `Open` | Downstream is unhealthy — calls are rejected with `ErrOpen`. |
| `Half-Open` | Cooldown elapsed — one probe call is allowed to test recovery. |

## Usage

```go
br := circuitbreaker.New(5, 30*time.Second)

if err := br.Allow(); err != nil {
    // circuit is open; skip the expensive downstream call
    return err
}

if err := webhook.Send(payload); err != nil {
    br.RecordFailure()
    return err
}
br.RecordSuccess()
```

## Parameters

- **threshold** — number of consecutive failures before the circuit opens.
- **cooldown** — duration to wait in the open state before probing recovery.

## Thread safety

All methods are safe for concurrent use.
