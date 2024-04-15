terraform {
  required_providers {
    prosimo = {
      version = "1.0.0"
      source  = "prosimo.com/prosimo/prosimo"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

###qing-nwob
provider "prosimo" {
  token    = var.token
  insecure = var.insecure
  base_url = var.url
}
variable "ip_pool_cidr" {default = "192.168.8.0/22"}
variable "prosimo_domain" { default = ".access.myeveqingchen1662950588667.scnetworkers.info"}
variable "az_app_account" {default = "prosimo-app"}
variable "cloud_creds_name" {default = "prosimo-aws-app-iam"}
