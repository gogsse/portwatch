# filter

The `filter` package determines whether an open port is **expected** or **unexpected** based on a user-supplied allow-list and optional privileged-port exemptions.

## Types

### `Filter`

Holds the allow-list and configuration flags.

```go
f := filter.New([]int{80, 443, 8080}, false)
```

Or build from string slices (e.g. loaded from config):

```go
f := filter.NewFromStrings(cfg.AllowedPorts, cfg.IgnorePrivileged)
```

## Key Methods

| Method | Description |
|---|---|
| `IsAllowed(port int) bool` | Returns true if the port is on the allow-list or exempt |
| `Unexpected(ports []int) []int` | Filters a slice down to only unexpected ports |
| `AllowedPorts() []int` | Returns the current allow-list snapshot |

## Privileged Port Exemption

When `ignorePrivileged` is `true`, any port below **1024** is considered allowed regardless of the explicit list. This avoids noise on systems that expose many well-known service ports.
