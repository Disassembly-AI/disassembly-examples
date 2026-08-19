# Quickstart

## Manual CLI (Docker — no install)
```bash
docker run --rm -e DISASSEMBLY_API_KEY \
  ghcr.io/disassembly-ai/toolkit:latest \
  scan https://staging.example.com

# CI mode: SARIF to stdout
docker run --rm -e DISASSEMBLY_API_KEY \
  ghcr.io/disassembly-ai/toolkit:latest \
  scan https://staging.example.com --ci > report.sarif
```

## Python library
See [`python/scan.py`](python/scan.py) — drive the engine programmatically and get typed findings.
```bash
pip install disassembly
DISASSEMBLY_API_KEY=... python python/scan.py https://staging.example.com
```
