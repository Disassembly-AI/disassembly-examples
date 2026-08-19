# Monitoring sidecar → single-pane view

Runs beside your service, scans it on a cadence, and **streams telemetry** to
Disassembly.AI. That stream powers the live "fabric" single pane on both the
customer portal (**Monitor**) and the admin dashboard (**Live** / fleet).

```
[ your app ] ── localhost ──> [ disassembly-monitor sidecar ]
                                        │  scan.started · finding · verified · host.up
                                        ▼
                             INGEST (SSE/WebSocket fan-out)
                                        │
                        ┌───────────────┴───────────────┐
                 customer portal /monitor        admin /live (fleet)
```

Deploy: [`docker-compose.yml`](docker-compose.yml) or [`k8s/monitor-sidecar.yaml`](k8s/monitor-sidecar.yaml).
See [`../MONITORING.md`](../MONITORING.md) for the end-to-end architecture.
