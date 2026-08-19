# CI/CD build runners (per environment)

One runner image (`Dockerfile`) that carries the toolkit + your build tooling, wired into each
environment's runner so scans run identically everywhere.

| Environment | How to use the runner |
|---|---|
| **GitHub Actions** | `container: ghcr.io/disassembly-ai/runner:latest` on the job, or register it as a self-hosted runner. See [`../.github/workflows/scan.yml`](../.github/workflows/scan.yml). |
| **AWS CodeBuild** | Set the project **image** to `ghcr.io/disassembly-ai/runner:latest` (or push to ECR). See [`../ci/aws-codebuild/buildspec.yml`](../ci/aws-codebuild/buildspec.yml). |
| **Azure Pipelines** | `container: ghcr.io/disassembly-ai/runner:latest` on the job, or a self-hosted agent from this image. See [`../ci/azure-pipelines/azure-pipelines.yml`](../ci/azure-pipelines/azure-pipelines.yml). |
| **GCP Cloud Build** | Use it as the step `name:`. See [`../ci/gcp-cloudbuild/cloudbuild.yaml`](../ci/gcp-cloudbuild/cloudbuild.yaml). |

Build & publish:
```bash
docker build -t ghcr.io/disassembly-ai/runner:latest runners/
docker push ghcr.io/disassembly-ai/runner:latest
```
