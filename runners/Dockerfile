# Shared build-runner image: the toolkit + common build tooling. Use it as a
# self-hosted runner / build image across GitHub Actions, AWS CodeBuild,
# Azure Pipelines, and GCP Cloud Build so scans run in the same environment
# everywhere.
FROM python:3.12-slim
RUN apt-get update && apt-get install -y --no-install-recommends git curl ca-certificates && rm -rf /var/lib/apt/lists/*
RUN pip install --no-cache-dir disassembly
# (add your language toolchains here: node, go, jdk, …)
ENTRYPOINT ["disassembly"]
