variable "project" {
  type = string
}
variable "region" {
  type    = string
  default = "us-central1"
}
variable "prefix" {
  type    = string
  default = "disassembly"
}
variable "target" {
  type = string
}
variable "scheduler_sa_email" {
  type        = string
  description = "SA with run.jobs.run permission"
}
