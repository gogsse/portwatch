# internal/snapshot

The `snapshot` package captures a timestamped point-in-time view of the
currently open ports on the host.

## Types

### `Snapshot`

```go
type Snapshot struct {
    CapturedAt time.Time `json:"captured_at"`
    Ports      []int     `json:"ports"`
}
```

## Functions

| Function | Description |
|---|---|
| `New(ports []int) *Snapshot` | Creates a new snapshot, copying the port slice. |
| `(s *Snapshot) Save(path string) error` | Persists the snapshot as JSON. |
| `Load(path string) (*Snapshot, error)` | Loads a snapshot from disk; returns `nil, nil` if the file is absent. |
| `(s *Snapshot) PortSet() map[int]struct{}` | Returns ports as a set for O(1) membership checks. |

## Usage

```go
snap := snapshot.New(currentPorts)
if err := snap.Save("/var/lib/portwatch/last.json"); err != nil {
    log.Println("could not save snapshot:", err)
}

prev, _ := snapshot.Load("/var/lib/portwatch/last.json")
if prev != nil {
    set := prev.PortSet()
    // compare with current scan …
}
```
