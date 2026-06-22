# filter

The `filter` package provides port filtering logic used to distinguish expected
listeners from unexpected ones.

## Types

### Filter

Holds an allowed set of ports and optional privileged-port exemption.

```go
f := filter.New(map[int]struct{}{80: {}, 443: {}}, true)
```

### NewFromStrings

Parses a slice of string port numbers into a Filter.

```go
f, err := filter.NewFromStrings([]string{"22", "80", "443"}, false)
```

## Methods

| Method | Description |
|---|---|
| `IsAllowed(port int) bool` | Returns true if port is in the allowed set or is a privileged port (when exemption enabled) |
| `Unexpected(ports []int) []int` | Returns ports not covered by the filter |

## Privileged Exemption

When `ExemptPrivileged` is true, ports below 1024 are always considered allowed
regardless of the allowed set. This reduces noise from well-known system services.
