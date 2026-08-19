variable "prefix" {
  type    = string
  default = "disassembly"
}
variable "location" {
  type    = string
  default = "eastus"
}
variable "target" {
  type = string
}
variable "api_key" {
  type      = string
  sensitive = true
}
