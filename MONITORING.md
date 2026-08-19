# Realtime monitoring — single-pane architecture

**Goal:** one live "fabric" view of every customer's scans and infrastructure —
shown to the customer (their own) and to admins (the whole fleet).

```
 monitoring sidecars           ingest + fan-out            single pane (SSE/WS)
 (per customer service)   ─▶   collector API          ─▶   • customer portal /monitor
   scan.started                (auth by instance key)      • admin dashboard /live
   finding / verified          normalizes + persists
   host.up / host.down         → event log + last-state
```

## Components
1. **Sidecar** (`monitor-sidecar/`) — streams telemetry events tagged with the account's instance key.
2. **Ingest API** — authenticates by instance key, writes to an event log (e.g. Kinesis/Kafka →
   Timestream/Dynamo for last-state), and fans out over **SSE or WebSocket**.
3. **Single pane** — the admin **Live** and portal **Monitor** pages subscribe to the fan-out.
   Today they run a client-side simulated stream; swap the simulator for
   `new EventSource("/api/stream?scope=…")` and the render layer is unchanged.

## Scoping & auth
- Customer stream: filtered to that account's assets (JWT → account id).
- Admin stream: fleet-wide, gated on the `Admins` group.
- The instance key authenticates sidecars; the Cognito session authenticates viewers.

## Status
Front-end single pane: **built and deployed** (simulated stream). Ingest + fan-out backend:
**to build** — tracked on our internal roadmap.
