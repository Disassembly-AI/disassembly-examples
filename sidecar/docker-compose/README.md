# Docker Compose sidecar

`app` starts, becomes healthy, then the `disassembly` sidecar scans it over the compose network
(`http://app:3000`) and writes `reports/report.sarif`. Great for local pre-push checks and
ephemeral CI environments.
