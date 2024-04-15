variable "onboard_app" { default = true }
variable "decommission_app" { default = false }
#variable "onboard_app" { default = false }
#variable "decommission_app" { default = true }

variable "policy_name_agentless" { default = "ALLOW-ALL"} #"allow-users" } #
variable "policy_name_agent" { default = "ALLOW-ALL"} #"agent-allow-users" } #



variable "port_list" { default = ["22-443"] }
variable "domain_1" { default = "psonar.us" }
variable "domain_2" { default = "psonar.local" }
###"10.100.0.22"--DNS server, "10.100.0.142"--access VM, "10.101.0.141"--for PrivateLink in another VPC,
###"10.10.0.36"-- AWS TGW peering VM, "172.17.2.4"-- Azure access VM for vwanhub
variable "vm_ips" { default = ["10.100.0.22", "10.100.0.142", "10.101.0.141", "10.10.0.36", "172.17.2.4","172.19.1.4"] }
variable "pub_dns_server" { default = ["8.8.8.8", "1.1.1.1"] }
variable "vpc_cider" { default = "10.100.0.0/24" }
variable "vnet_ciders" {default = "172.19.1.0/24,172.19.2.0/24"}

###AWS TGW info
variable "TGW_id" { default = "tgw-04d5146c784e4d6d6" }
variable "vpc_id_tgw" { default = "vpc-0ede62043db530015" }
#variable "TGW_id" { default = "tgw-02c98e2afd03758d7" }
#variable "vpc_id_tgw" { default = "vpc-067648a8b45cd2369" }
#variable "agent_vpc_id_tgw" { default = "vpc-0aa9baf4e2c6802d3"}

#NLB for PrivatrLink info
variable "NLB_name" { default = "qing-NLB-e074b26d5dbd1c34.elb.us-east-2.amazonaws.com" }
variable "vpc_id_nlb" { default = "vpc-03bd678f593aa8866" }

###vwanhub info
variable "vnet" {
  default = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vm-2-rg/providers/Microsoft.Network/virtualNetworks/qing-vm-2-rg-vnet"
}
variable "hub" {
  default = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"
}
variable "az_region" { default = "eastus2"}

##### Azure transit vNET info
variable "az_resouce" {default = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing_transit_rg/providers/Microsoft.Network/virtualNetworks/"}
variable "az_hub" { default = "qing-transit-hub-vnet"}
variable "az_spoke_1" { default = "transit_spoke1"}
variable "az_spoke_2" { default = "transit_spoke2"}


variable "web_list" {
  default = [ "fastly.com","azure-api.us"]
#    "fastly.com", "digitalocean.com", "sendgrid.com", "sendgrid.net", "oraclecloud.com", "cloudfront.net",
#    "keycdn.com", "kxcdn.com", "yahoo.com", "azure-api.us"
#  ]
}
#aws.amazon.com
variable "app_list" {
  default = [ "app1.psonar.local", "app2.psonar.local"]
#    "app1.psonar.local", "app2.psonar.local", "app3.psonar.local", "app4.psonar.local", "app5.psonar.local",
#    "app6.psonar.local", "app7.psonar.local", "app8.psonar.local", "app9.psonar.local", "app10.psonar.local"
#  ]
}

variable "agentless_app_names" {
  default = [
    "agentless", "agentless-ssh-tf", "agentless-multi-VMs-tf", "url-rewrite-ssh-tf", "url-rewrite-https-tf",
    "jumpbox-peering-tf", "aws-s3-tf", "jumpbox-TGW-tf", "ssh-PrivateLink-tf", "vwanhub-ssh-tf", "vwanhub-jumpbox-tf"
  ]
}

variable "agent_app_names" {
  default = [
    "agent", "Agent-bulk-webs-tf", "Agent-single-IP-tf", "Agent-bulk-apps-via-DNS-IP-tf", "Agent-single-fqdn-dns-IP-tf",
    "Agent-DNS-svr-tf", "Agent-single-fqdn-dns-app-tf", "Agent-DNS-bulk-tf", "agent-vwanhub-ssh-https-tf",
    "agent-transit-vnet-subnets-tf", "Agent-TGW-tf"
  ]
}

locals {
  app_access_type_agentless = "agentless"
  app_access_type_agent = "agent"
  ip_pool_cidr = var.ip_pool_cidr
  idp_name = "azure_ad"
  onboard_type = "behind_fabric"
  interaction_type = "userToApp"
  address_type = "fqdn"
  domain_type = "prosimo"
  connection_option = "private"
  cloud_creds_name = var.cloud_creds_name
  region_name = "us-west-1" #"us-east-2"
  region_type = "active"
  conn_option = "peering"
  backend_ip_address_discover = false
  is_healthcheck_enabled = true
  optimization_option = "PerformanceEnhanced"
  enable_multi_cloud_access = true
  onboard_app = var.onboard_app
  decommission_app = var.decommission_app
}
