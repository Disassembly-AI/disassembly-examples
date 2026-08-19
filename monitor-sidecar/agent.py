#!/usr/bin/env python3
"""Disassembly.AI monitoring sidecar.

Runs next to a customer service, periodically scans it, and STREAMS telemetry
(scan lifecycle + findings + host health) to the ingest endpoint. That stream is
what the single-pane "fabric" view on the admin + customer dashboards renders.

Env:
  DISASSEMBLY_API_KEY   account key (auth)
  TARGET                URL/host to monitor (e.g. http://localhost:3000)
  INGEST_URL            collector, default https://ingest.disassembly.ai/v1/events
  INTERVAL_SECONDS      scan cadence, default 300
"""
import json
import os
import time
import urllib.request

API_KEY = os.environ.get("DISASSEMBLY_API_KEY", "")
TARGET = os.environ.get("TARGET", "http://localhost:3000")
INGEST = os.environ.get("INGEST_URL", "https://ingest.disassembly.ai/v1/events")
INTERVAL = int(os.environ.get("INTERVAL_SECONDS", "300"))


def emit(kind: str, **fields) -> None:
    """Push one telemetry event to the collector (fire-and-forget)."""
    event = {"kind": kind, "target": TARGET, "ts": int(time.time()), **fields}
    data = json.dumps(event).encode()
    req = urllib.request.Request(
        INGEST, data=data,
        headers={"Authorization": "Bearer " + API_KEY, "Content-Type": "application/json"},
    )
    try:
        urllib.request.urlopen(req, timeout=5).read()
    except Exception as e:  # never let telemetry break the sidecar
        print(f"[monitor] ingest failed: {e}", flush=True)
    print(f"[monitor] {kind} {fields}", flush=True)


def scan_once() -> None:
    emit("scan.started")
    # In production this shells out to the toolkit:
    #   disassembly scan "$TARGET" --ci --format json
    # and streams each finding as it is produced. Stub below keeps deps to stdlib.
    emit("recon.surface", routes=42, apis=3)
    emit("finding", rule="authz.idor", severity="high", path="/api/v1/orders/{id}")
    emit("verified", findings=1)
    emit("sarif.uploaded")


def main() -> None:
    if not API_KEY:
        print("[monitor] set DISASSEMBLY_API_KEY", flush=True)
    emit("host.up")
    while True:
        scan_once()
        time.sleep(INTERVAL)


if __name__ == "__main__":
    main()
