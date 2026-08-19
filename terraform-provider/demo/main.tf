# Full local demo: build the provider, run the mock API, then `tofu apply`.
# See ../README.md → "Local demo (mock API)".
terraform {
  required_providers {
    disassembly = { source = "Disassembly-AI/disassembly" }
  }
}

provider "disassembly" {
  endpoint = "http://localhost:8080" # the mock server
  # api_key comes from DISASSEMBLY_API_KEY (any non-empty value for the mock)
}

resource "disassembly_scan" "demo" {
  target           = "https://staging.example.com"
  effort           = "high"
  fail_on_severity = "none" # set to "high" to watch the apply FAIL on the IDOR
}

output "summary" {
  value = {
    total  = disassembly_scan.demo.findings_count
    high   = disassembly_scan.demo.high_count
    medium = disassembly_scan.demo.medium_count
    report = disassembly_scan.demo.report_url
  }
}
