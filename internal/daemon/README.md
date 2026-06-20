# daemon

The `daemon` package wires together all portwatch subsystems and drives the
main monitoring loop.

## Overview

```
Daemon
├── scanner  – reads /proc/net/tcp* for open ports
├── watcher  – computes diff between scan cycles
├── notifier – writes human-readable alerts to stdout (or a custom writer)
└── reporter – buffers structured events and flushes them periodically
```

## Usage

```go
cfg := config.DefaultConfig()
cfg.Interval = 5 * time.Second

d := daemon.New(cfg)

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

if err := d.Run(ctx); err != nil && err != context.Canceled {
    log.Fatal(err)
}
```

## Behaviour

- An initial scan is performed immediately on startup.
- Subsequent scans run on every `cfg.Interval` tick.
- When the context is cancelled the reporter is flushed before the daemon
  returns.
- Scan errors are logged but do not stop the loop.
