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
| `Allowed() []int` | Returns a sorted slice of explicitly allowed ports |

## Privileged Exemption

When `ExemptPrivileged` is true, ports below 1024 are always considered allowed
regardless of the allowed set. This reduces noise from well-known system services.

## Error Handling

`NewFromStrings` returns an error if any string in the slice cannot be parsed as
an integer, or if a port number falls outside the valid range (1–65535).

```go
f, err := filter.NewFromStrings([]string{"80", "99999"}, false)
if err != nil {
    // err: port 99999 out of valid range [1, 65535]
}
```
