variable "region" {
  type    = string
  default = "us-east-1"
}
variable "prefix" {
  type    = string
  default = "disassembly"
}
variable "target" {
  type = string
}
variable "subnet_ids" {
  type = list(string)
}
variable "security_group_ids" {
  type = list(string)
}
