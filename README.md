# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected listeners.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a baseline of allowed ports:

```bash
portwatch --allow 22,80,443 --interval 30s
```

Run once and print all currently open ports:

```bash
portwatch --scan
```

Alert via stdout when an unexpected listener is detected:

```
[ALERT] 2024/01/15 14:32:01 Unexpected listener detected: 0.0.0.0:4444 (PID 9821)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--allow` | `""` | Comma-separated list of allowed ports |
| `--interval` | `60s` | How often to scan for open ports |
| `--scan` | `false` | Run a single scan and exit |
| `--verbose` | `false` | Enable verbose logging |

## License

MIT © [yourusername](https://github.com/yourusername)