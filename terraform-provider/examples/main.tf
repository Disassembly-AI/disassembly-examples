terraform {
  required_providers {
    disassembly = {
      source = "Disassembly-AI/disassembly"
    }
  }
}

provider "disassembly" {
  # api_key = var.api_key   # or export DISASSEMBLY_API_KEY
  # endpoint = "https://api.disassembly.ai"
}

# A scan modeled as a first-class resource. The apply FAILS if a high-severity
# finding exists — so you can make infrastructure depend on a clean scan.
resource "disassembly_scan" "staging" {
  target           = "https://staging.example.com"
  effort           = "high"
  fail_on_severity = "high"
}

# Gate a deploy on the scan: this only creates after the scan passes.
resource "null_resource" "deploy" {
  depends_on = [disassembly_scan.staging]
  provisioner "local-exec" {
    command = "echo scan clean — proceeding with deploy"
  }
}

output "scan" {
  value = {
    id      = disassembly_scan.staging.id
    total   = disassembly_scan.staging.findings_count
    high    = disassembly_scan.staging.high_count
    medium  = disassembly_scan.staging.medium_count
    report  = disassembly_scan.staging.report_url
  }
}
