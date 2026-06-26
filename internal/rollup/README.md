# rollup

The `rollup` package groups repeated port events into summarised bursts,
reducing noise when the same port opens and closes repeatedly within a
short monitoring window.

## Concepts

| Term | Meaning |
|------|---------|
| **Event** | A rolled-up summary for a single `(kind, port)` pair |
| **Window** | The aggregation duration before a flush is expected |
| **Count** | Number of times the same event was recorded before flush |

## Usage

```go
rl := rollup.New(30 * time.Second)

// Call Record each time a port event is observed.
rl.Record("opened", 8080)
rl.Record("opened", 8080) // same port — increments count
rl.Record("closed", 9000)

// At the end of the interval, flush and process summaries.
for _, ev := range rl.Flush() {
    fmt.Printf("%s port %d (seen %d times)\n", ev.Kind, ev.Port, ev.Count)
}
```

## Integration

`Rollup` is intended to sit between the monitor watcher and the alert
notifier. Feed raw diff events into `Record`, then call `Flush` at the
end of each daemon tick and forward the resulting summaries to the
notifier or digest.

## Thread Safety

All methods are safe for concurrent use.
