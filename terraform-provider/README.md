# terraform-provider-disassembly

A **custom Terraform provider** that exposes the Disassembly.AI pentest engine as a first-class
resource:

```hcl
resource "disassembly_scan" "staging" {
  target           = "https://staging.example.com"
  effort           = "high"
  fail_on_severity = "high"   # apply FAILS if a high-severity finding exists
}
```

Because the scan is a resource, other resources can `depends_on` it — so you can **gate a deploy on a
clean scan** entirely within `terraform apply`, no separate CI step required.

## Resource: `disassembly_scan`

| Attribute | | Description |
|---|---|---|
| `target` | required, forces replace | URL/host to scan (**authorized only**) |
| `effort` | optional | `low`\|`medium`\|`high`\|`xhigh`\|`max` (default `high`) |
| `fail_on_severity` | optional | `none`\|`medium`\|`high` (default `high`) — apply fails at/above this |
| `id` · `findings_count` · `high_count` · `medium_count` · `sarif` · `report_url` | computed | scan outputs |

Provider config: `api_key` (or `DISASSEMBLY_API_KEY`), `endpoint` (default `https://api.disassembly.ai`).

## Build & run locally

> Requires Go 1.22+. This example is written against `terraform-plugin-framework` v1.x. Run
> `go mod tidy` once to resolve indirect deps.

```bash
make tidy      # resolve deps
make build     # -> ./terraform-provider-disassembly
```

Then wire up a **dev override** so Terraform uses your local binary (no registry needed): copy
[`terraformrc.example`](terraformrc.example) into `~/.terraformrc`, set the absolute path, and:

```bash
cd examples
DISASSEMBLY_API_KEY=dsa_live_xxx terraform plan   # dev_overrides skips `terraform init`
```

## Local demo (mock API) — run a full `apply` with no backend

The engine's REST API isn't live yet, so a mock is included that returns canned findings — enough
to run a real `terraform apply` end to end.

```bash
# 1. build the provider and the mock
make build
go build -o /tmp/disa-mock ./mockserver

# 2. run the mock API (leave it running)
/tmp/disa-mock            # listens on :8080

# 3. point Terraform at the local provider (dev override) and apply the demo
cp terraformrc.example ~/.terraformrc      # edit the path to this dir
cd demo
export DISASSEMBLY_API_KEY=anything        # mock accepts any non-empty key
tofu apply -auto-approve                   # or: terraform apply
```

You'll see the scan resource created with `high_count = 1`, `medium_count = 1`, and a report URL.
Flip `fail_on_severity` to `"high"` in `demo/main.tf` and re-apply to watch the **apply fail** on the
high-severity finding — that's the deploy gate.

## Verified
Compile + load tested (this repo):
- Go 1.26 — `go vet` clean, `gofmt` clean, `go build` produces a valid plugin binary.
- OpenTofu 1.12 (Terraform-compatible) — `tofu validate` = **valid**, and `tofu plan`/`apply` load the
  provider and run a full scan against the bundled mock — including the deploy gate failing the
  apply on a high-severity finding (via the `dev_overrides` in `terraformrc.example`).

## Status
The provider talks to the Disassembly.AI REST API (`POST /v1/scans`, `GET /v1/scans/{id}`). That
API is part of the product build — until it's live, use this as the integration contract and for
`terraform plan`. Publishing to the Terraform Registry (`Disassembly-AI/disassembly`) is a later step.
