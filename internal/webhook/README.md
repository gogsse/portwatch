# webhook

Delivers port-event alerts to an external HTTP endpoint via JSON `POST`.

## Overview

The `webhook` package complements the in-process `alert` notifier with outbound
HTTP delivery, enabling integrations with Slack incoming webhooks, PagerDuty,
custom SIEM collectors, or any HTTP receiver.

## Payload schema

```json
{
  "port":      8080,
  "kind":      "opened",
  "severity":  "warning",
  "timestamp": "2024-01-15T12:34:56Z"
}
```

| Field       | Values                          |
|-------------|----------------------------------|
| `kind`      | `opened` \| `closed`             |
| `severity`  | `info` \| `warning` \| `danger`  |

## Usage

```go
wh := webhook.New("https://hooks.example.com/portwatch", 5*time.Second)

// Send a pre-built payload.
wh.Send(webhook.Payload{Port: 4444, Kind: "opened", Severity: "danger"})

// Convenience helpers.
wh.SendOpened(8080, "warning")
wh.SendClosed(8080)
```

## Error handling

`Send` returns an error for network failures and any non-2xx HTTP status code.
Callers are responsible for retry logic; the package intentionally stays
stateless to compose cleanly with `ratelimit` or `throttle`.
