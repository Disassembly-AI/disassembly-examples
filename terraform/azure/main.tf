terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "this" {
  name     = "${var.prefix}-scan"
  location = var.location
}

# Container group that runs one scan (Never restart). Schedule it with a Logic App
# recurrence that starts this group, or re-create it on a pipeline schedule.
resource "azurerm_container_group" "scan" {
  name                = "${var.prefix}-scan"
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  os_type             = "Linux"
  restart_policy      = "Never"

  container {
    name     = "disassembly"
    image    = "ghcr.io/disassembly-ai/toolkit:latest"
    cpu      = "1"
    memory   = "2"
    commands = ["disassembly", "scan", var.target, "--ci"]
    secure_environment_variables = {
      DISASSEMBLY_API_KEY = var.api_key # pass via TF_VAR_api_key or a Key Vault data source
    }
  }
}
