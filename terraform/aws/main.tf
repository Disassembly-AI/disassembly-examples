terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.60"
    }
  }
}

provider "aws" {
  region = var.region
}

resource "aws_s3_bucket" "reports" {
  bucket = "${var.prefix}-scan-reports"
}

resource "aws_secretsmanager_secret" "api_key" {
  name = "${var.prefix}/disassembly-api-key"
}
# Set the value out-of-band:
#   aws secretsmanager put-secret-value --secret-id <name> \
#     --secret-string '{"DISASSEMBLY_API_KEY":"dsa_live_xxx"}'

module "scan" {
  source             = "../modules/disassembly-scan"
  name               = "${var.prefix}-scan"
  target             = var.target
  subnet_ids         = var.subnet_ids
  security_group_ids = var.security_group_ids
  api_key_secret_arn = aws_secretsmanager_secret.api_key.arn
  reports_bucket     = aws_s3_bucket.reports.bucket
}

output "scan_schedule" {
  value = module.scan.schedule
}
