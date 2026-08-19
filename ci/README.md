# CI pipeline examples

Every file runs `disassembly scan <target> --ci` and publishes the SARIF to that platform's
code-scanning / artifact surface. In all of them:

- **`DISASSEMBLY_API_KEY`** comes from the platform's secret store (masked variable, Secrets
  Manager, Key Vault, credential, or context) — never committed.
- **target** is a variable so you can scan staging on PRs and production on a schedule.
- a **high-severity finding exits non-zero**, so you can gate merges/deploys on it.

| Platform | File |
|---|---|
| GitHub Actions | [`../.github/workflows/scan.yml`](../.github/workflows/scan.yml) |
| GitLab CI | [`gitlab/.gitlab-ci.yml`](gitlab/.gitlab-ci.yml) |
| AWS CodeBuild | [`aws-codebuild/buildspec.yml`](aws-codebuild/buildspec.yml) |
| Azure Pipelines | [`azure-pipelines/azure-pipelines.yml`](azure-pipelines/azure-pipelines.yml) |
| Google Cloud Build | [`gcp-cloudbuild/cloudbuild.yaml`](gcp-cloudbuild/cloudbuild.yaml) |
| CircleCI | [`circleci/config.yml`](circleci/config.yml) |
| Jenkins | [`jenkins/Jenkinsfile`](jenkins/Jenkinsfile) |
