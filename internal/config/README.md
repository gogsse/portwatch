# config

The `config` package provides runtime configuration loading and defaults for **portwatch**.

## Usage

### Default configuration

```go
cfg := config.DefaultConfig()
```

Defaults:

| Field              | Default        | Description                              |
|--------------------|----------------|------------------------------------------|
| `interval_seconds` | `10`           | Seconds between port scans               |
| `allowed_ports`    | `[22, 80, 443]`| Ports considered safe / expected         |
| `alert_on_close`   | `false`        | Alert when a previously open port closes |
| `log_file`         | `""`           | Log destination (empty = stdout)         |

### Load from file

```go
cfg, err := config.LoadFromFile("/etc/portwatch/config.json")
```

Example `config.json`:

```json
{
  "interval_seconds": 15,
  "allowed_ports": [22, 80, 443, 8080],
  "alert_on_close": true,
  "log_file": "/var/log/portwatch.log"
}
```

Unrecognised keys are silently ignored; missing keys retain their defaults.
