terraform {
  required_providers {
    prosimo = {
      version = "1.0.0"
      source  = "prosimo.com/prosimo/prosimo"
    }
  }
}

provider "prosimo" {
  token    = var.token
  insecure = var.insecure
  base_url = var.url
}

resource "prosimo_edr_config" "crowdstrike" {
    name = "demo3"
    vendor = "CrowdStrike"
    auth {
      client_id = "08dc7dc0e2174b478174c6040169b4a1"
      client_secret = var.clientSecret
      base_url = "https://demo.crowdstrike.com"
      customer_id = "customer_id"
      mssp = false
    }
}








