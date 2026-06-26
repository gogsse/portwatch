# correlation

Groups consecutive port-change events into **bursts** and classifies them to
reduce false-positive alerts.

## Concepts

| Term | Description |
|------|-------------|
| **Burst** | A collection of events observed within the correlation window. |
| **Window** | Maximum gap between the first event and any subsequent event that still belongs to the same burst. |
| **BurstClass** | Inferred intent: `noise`, `restart`, `intrusion`, or `sweep`. |

## Usage

```go
corr := correlation.New(3 * time.Second)

for _, e := range events {
    if burst := corr.Add(e); burst != nil {
        cls := correlation.Classify(burst, 2)
        fmt.Println("burst class:", cls)
    }
}

// drain remaining events at end of tick
if burst := corr.Flush(); burst != nil {
    cls := correlation.Classify(burst, 2)
    fmt.Println("final burst class:", cls)
}
```

## Classification rules

- **noise** — fewer events than the configured threshold.
- **restart** — both opened and closed events present (service bounce).
- **intrusion** — opened events outnumber closed events.
- **sweep** — closed events only (mass service shutdown).
