# Kubernetes

- **`deployment-sidecar.yaml`** — a per-pod sidecar that scans the app over `localhost`. Use for
  ephemeral/preview environments where you want a scan tied to the pod lifecycle.
- **`cronjob.yaml`** — a scheduled `CronJob` that scans a Service by DNS name. Preferred for
  long-lived environments (no idle sidecar burning resources).

Manage the API key with a real secrets tool (External Secrets Operator, Sealed Secrets, or your
cloud's CSI driver) rather than the inline `Secret` shown here.
