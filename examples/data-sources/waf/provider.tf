terraform {
  required_providers {
    prosimo = {
      version = "1.0.0"
      source  = "prosimo.io/prosimo/prosimo"
    }
  }
}


provider "prosimo" {
  token    = var.token
  insecure = var.insecure
  base_url = var.url
}
