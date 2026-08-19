# Terraform examples

Infrastructure to run scheduled scans on each major cloud. All illustrative — run `terraform init`
with your own backend and credentials, review the plan, and wire the API key into the cloud's
secret store (never a variable in source control).

| Cloud | Pattern | Path |
|---|---|---|
| **AWS** | EventBridge → Fargate task, SARIF to S3 | [`aws/`](aws/) (uses [`modules/disassembly-scan`](modules/disassembly-scan/)) |
| **Azure** | Container Instances (schedule via Logic App) | [`azure/`](azure/) |
| **GCP** | Cloud Run job + Cloud Scheduler | [`gcp/`](gcp/) |

### Setting the API key (out-of-band, after apply)
```bash
# AWS
aws secretsmanager put-secret-value --secret-id disassembly/disassembly-api-key \
  --secret-string '{"DISASSEMBLY_API_KEY":"dsa_live_xxx"}'
# GCP
echo -n "dsa_live_xxx" | gcloud secrets versions add disassembly-api-key --data-file=-
# Azure: source var.api_key from a Key Vault data source or TF_VAR_api_key at apply time.
```

## Prefer a native resource?
There's also a **custom Terraform provider** with a first-class `disassembly_scan` resource (gate a deploy on a clean scan inside `terraform apply`): see [`../terraform-provider/`](../terraform-provider/).
