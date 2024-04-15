resource "prosimo_network_onboarding" "aws-subnets" {
  name = "aws-subnets-tf"
  public_cloud {
    cloud_type        = "public"
    connection_option = local.connection_option
    cloud_creds_name  = local.cloud_creds_name
    region_name       = local.region_name

    cloud_networks {
      vpc               = "vpc-04ad26e924bc5e138" #"vpc-0aa9baf4e2c6802d3"
      # hub_id = "tgw-04db5eac6fe3de45e"
      connectivity_type = "vpc-peering"
      subnets           = ["10.100.0.16/28", "10.100.0.128/28"]
    }
    cloud_networks {
      vpc               = var.vpc_id_tgw #"vpc-067648a8b45cd2369"
      # hub_id = "tgw-04db5eac6fe3de45e"
      connectivity_type = "vpc-peering"
      subnets           = ["10.10.0.0/28", "10.10.0.32/28"]
    }
    connect_type = "connector"

  }
  policies         = ["DEFAULT-MCN-POLICY", "DEMO2"]
  onboard_app      = local.onboard_app
  decommission_app = local.decommission_app
}

#resource "prosimo_network_onboarding" "testapp-azure" {
#
#    name = "azure_network"
#    public_cloud {
#        cloud_type = "public"
#        connection_option = "private"
#        cloud_creds_name = "prosimo-app"
#        region_name = "eastus2"
#        cloud_networks {
#          vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/app-azure-eastus2-1661236757258-rg/providers/Microsoft.Network/virtualNetworks/app-azure-eastus2-1661236757258-vnet"
#          # hub_id = "tgw-04db5eac6fe3de45e"
#          connectivity_type = "vnet-peering"
#          subnets = ["192.168.128.0/25"]
#        }
#        cloud_networks {
#          vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Gitlab/providers/Microsoft.Network/virtualNetworks/Gitlab-vnet"
#          # hub_id = "tgw-04db5eac6fe3de45e"
#          connectivity_type = "vnet-peering"
#          subnets = ["10.3.5.0/24"]
#        }
#        connect_type = "connector"
#
#    }
#    policies = ["DEFAULT-MCN-POLICY"]
#    onboard_app = true
#    decommission_app = false
#}