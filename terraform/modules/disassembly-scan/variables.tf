variable "name" {
  type    = string
  default = "disassembly-scan"
}
variable "target" {
  type        = string
  description = "URL/host to scan (authorized only)"
}
variable "schedule" {
  type    = string
  default = "cron(0 6 ? * MON *)"
}
variable "image" {
  type    = string
  default = "ghcr.io/disassembly-ai/toolkit:latest"
}
variable "subnet_ids" {
  type = list(string)
}
variable "security_group_ids" {
  type = list(string)
}
variable "api_key_secret_arn" {
  type        = string
  description = "Secrets Manager ARN holding DISASSEMBLY_API_KEY"
}
variable "reports_bucket" {
  type        = string
  description = "S3 bucket for SARIF reports"
}
