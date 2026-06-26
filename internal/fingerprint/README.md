# fingerprint

The `fingerprint` package resolves a listening port number to the OS process
that owns the socket.

## How it works

1. Reads `/proc/net/tcp` to find the inode associated with a local port.
2. Walks `/proc/<pid>/fd/` for every running process, following symlinks to
   locate the file descriptor that points to `socket:[inode]`.
3. Reads `/proc/<pid>/comm` for a human-readable process name.

## Usage

```go
resolver := fingerprint.New()
info, err := resolver.Lookup(8080)
if err != nil {
    log.Println("could not fingerprint port:", err)
} else {
    fmt.Println(info) // port 8080 -> myapp (pid 3412)
}
```

## Types

### `Info`

| Field  | Type   | Description                        |
|--------|--------|------------------------------------|
| Port   | int    | The port number that was queried   |
| PID    | int    | PID of the owning process          |
| Name   | string | Process name from `/proc/pid/comm` |

### `Resolver`

Created via `New()`. Exposes a single method:

```go
func (r *Resolver) Lookup(port int) (Info, error)
```

## Limitations

- Linux-only (relies on `/proc/net/tcp`).
- Only inspects IPv4 listeners; IPv6 requires `/proc/net/tcp6`.
- Requires read access to `/proc/<pid>/fd` for all processes of interest;
  run as root or with `CAP_DAC_READ_SEARCH` for full coverage.
