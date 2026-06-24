# internal/triage

The `triage` package classifies open-port events by **severity** so downstream
components (alerts, reports) can prioritise their output.

## Severity levels

| Level      | Meaning                                                  |
|------------|----------------------------------------------------------|
| `INFO`     | Port is unknown but not considered dangerous             |
| `WARNING`  | Port is in the caller-supplied watch list                |
| `CRITICAL` | Port matches a built-in list of high-risk listeners      |

## Built-in critical ports

The classifier ships with a curated list that includes common remote-access and
malware-associated ports (22, 23, 445, 3389, 4444, 5900, …). See `triage.go`
for the full set.

## Usage

```go
c := triage.New([]int{8080, 9090}) // optional warning-level ports

sev := c.Classify(4444)  // → triage.SeverityCritical
fmt.Println(sev)         // → "CRITICAL"
```

## Integration with the daemon

The `daemon` package calls `Classify` for every newly-opened port detected by
the watcher and passes the resulting `Severity` to the notifier and reporter so
alerts include an urgency label.
