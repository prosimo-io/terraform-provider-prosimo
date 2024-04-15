

# # Azure with VWAN hub
# resource "prosimo_network_onboarding" "testapp-azure" {

#     name = "demo_network_azure"
#     network_exportable_policy = false
#     namespace = "default"
#     public_cloud {
#         cloud_type = "public"
#         connection_option = "private"
#         cloud_creds_name = "prosimo-app"
#         region_name = "eastus2"
#         cloud_networks {
#           vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/app-azure-eastus2-1661236757258-rg/providers/Microsoft.Network/virtualNetworks/app-azure-eastus2-1661236757258-vnet"
#           hub_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"
#           connectivity_type = "vwan-hub"
#           connector_placement = "none"
#           subnets = ["192.168.128.0/25"]
#         }
#         cloud_networks {
#           vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Gitlab/providers/Microsoft.Network/virtualNetworks/Gitlab-vnet"
#           hub_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"
#           connectivity_type = "vwan-hub"
#           connector_placement = "Workload VNET"
#           subnets = ["10.3.5.0/24"]
#         }
#         connect_type = "connector"

#     }
#     onboard_app = false
#     decommission_app = false
# }


# # GCP Workload Connector Placement
# resource "prosimo_network_onboarding" "testapp-s3" {

#     name = "demo_network_gcp"
#     network_exportable_policy = false
#     namespace = "default"
#     public_cloud {
#         cloud_type = "public"
#         connection_option = "private"
#         cloud_creds_name = "prosimo-gcp-infra"
#         region_name = "us-west2"
#         cloud_networks {
#           vpc = "https://www.googleapis.com/compute/v1/projects/prosimo-test-infra/global/networks/default"
#           connector_placement = "Workload VPC"
#           connectivity_type = "vpc-peering"
#           subnets = ["10.168.0.0/20"]
#           connector_settings {
#             connector_subnets= ["10.168.0.0/20"]
#             bandwidth = "<1 Gbps"
#             instance_type = "e2-standard-2"
#           }
#         }
#         connect_type = "connector"

#     }
#     policies = ["ALLOW-ALL-NETWORKS"]
#     onboard_app = false
#     decommission_app = false
# }

# #Azure Workload Connector Placement
# resource "prosimo_network_onboarding" "testapp-azure-workload-vpc" {

#     name = "demo_network_new"
#     network_exportable_policy = false
#     namespace = "default"
#     public_cloud {
#         cloud_type = "public"
#         connection_option = "private"
#         cloud_creds_name = "prosimo-infra"
#         region_name = "westus"
#         cloud_networks {
#           vnet = "/subscriptions/77102da4-2e1f-4445-b74a-93e842dc8c3c/resourceGroups/DefaultResourceGroup-WUS/providers/Microsoft.Network/virtualNetworks/DefaultResourceGroupWUSvnet574"
#           connectivity_type = "vnet-peering"
#           connector_placement = "Workload VNET"
#           subnets = ["10.4.0.0/24"]
#           connector_settings {
#             connector_subnets= ["10.4.0.0/24"]
#           }
#         }
#         connect_type = "connector"

#     }
#     policies = ["ALLOW-ALL-NETWORKS"]
#     onboard_app = false
#     decommission_app = false
# }

# # #Azure Infra Connector Placement
# resource "prosimo_network_onboarding" "testapp-azure-infra-vpc" {

#     name = "demo_network_new"
#     network_exportable_policy = false
#     namespace = "default"
#     public_cloud {
#         cloud_type = "public"
#         connection_option = "private"
#         cloud_creds_name = "prosimo-infra"
#         region_name = "westus"
#         cloud_networks {
#           vnet = "/subscriptions/77102da4-2e1f-4445-b74a-93e842dc8c3c/resourceGroups/DefaultResourceGroup-WUS/providers/Microsoft.Network/virtualNetworks/DefaultResourceGroupWUSvnet574"
#           connectivity_type = "vnet-peering"
#           connector_placement = "Workload VNET"
#           subnets {
#             subnet = "10.4.0.0/24"
#             # virtual_subnet = "10.250.2.128/25"
#           }
#           connector_settings {
#             connector_subnets = ["10.4.0.0/24"]
#             bandwidth_range {
#                 min = 3
#                 max = 5
#             }
#           }
#         }
#         connect_type = "connector"

#     }
#     policies = ["ALLOW-ALL-NETWORKS"]
#     onboard_app = false
#     decommission_app = false
# }

resource "prosimo_network_onboarding" "testapp-azure-infra-vpc_vn" {
    name = "demo_network_new"
    network_exportable_policy = false
    namespace = "default"
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-infra"
        region_name = "westus"
        cloud_networks {
          vnet = "/subscriptions/77102da4-2e1f-4445-b74a-93e842dc8c3c/resourceGroups/psonar-app-westus-rg-common-app/providers/Microsoft.Network/virtualNetworks/psonar-app-westus-rg-vnet"
          connectivity_type = "vnet-peering"
          connector_placement = "Infra VNET"
          subnets {
            subnet = "10.3.10.0/24"
            # virtual_subnet = "10.68.0.0/24"
          }
          connector_settings {
            connector_subnets = ["10.3.10.0/24"]
            bandwidth_range {
                min = 3
                max = 5
            }
          }
        }
        connect_type = "connector"

    }
    policies = ["ALLOW-ALL-NETWORKS"]
    onboard_app = false
    decommission_app = false
}

# #AWS with transit gateway and workload vpc
# resource "prosimo_network_onboarding" "testapp-AWS-WorkLoad-vpc" {

#     name = "demo_network_aws"
#     namespace = "default"
#     network_exportable_policy = false
#     public_cloud {
#         cloud_type = "public"
#         connection_option = "private"
#         cloud_creds_name = "prosimo-aws-iam"
#         region_name = "us-east-2"
#         cloud_networks {
#           vpc = "vpc-a8892dc3"
#           hub_id = "tgw-04d69a6cd846cd26b"
#           connector_placement = "Workload VPC"
#           connectivity_type = "transit-gateway"
#           service_insertion_endpoint_subnets = "auto"
#           subnets = ["172.31.0.0/20"]
#           connector_settings {
#             bandwidth = "1-5 Gbps"
#             instance_type = "c5a.large"
#             connector_subnets = ["172.31.0.0/20"]
#           }
#         }

#         connect_type = "connector"

#     }
#     policies = ["ALLOW-ALL-NETWORKS"]
#     onboard_app = false
#     decommission_app = false
# }

# #AWS with transit gateway and infra vpc
# resource "prosimo_network_onboarding" "testapp-AWS-Infra-vpc" {

#     name = "demo_network_aws"
#     namespace = "default"
#     network_exportable_policy = true
#     public_cloud {
#         cloud_type = "public"
#         connection_option = "private"
#         cloud_creds_name = "prosimo-aws-app-iam"
#         region_name = "eu-west-2"
#         cloud_networks {
#           vpc = "vpc-b2197fda"
#           # hub_id = "tgw-04d69a6cd846cd26b"
#           connector_placement = "Workload VPC"
#           connectivity_type = "vpc-peering"
#           subnets = ["172.31.0.0/20"]
#           service_insertion_endpoint_subnets = "auto"
#           connector_settings {
#             bandwidth = "<1 Gbps"
#             instance_type = "t3.medium"
#             connector_subnets = ["172.31.0.0/20"]
#           }
#         }

#         connect_type = "connector"

#     }
#     policies = ["ALLOW-ALL-NETWORKS"]
#     onboard_app = true
#     decommission_app = false
# }

# #AWS with transit gateway and none
# resource "prosimo_network_onboarding" "testapp-AWS-none" {

#     name = "demo_network_aws"
#     namespace = "default"
#     network_exportable_policy = false
#     public_cloud {
#         cloud_type = "public"
#         connection_option = "private"
#         cloud_creds_name = "prosimo-aws-app-iam"
#         region_name = "eu-west-2"
#         cloud_networks {
#           vpc = "vpc-05fcc0826ec4cf188"
#           hub_id = "tgw-0186760ef53675f35"
#           connector_placement = "Infra VPC"
#           connectivity_type = "transit-gateway"
#           subnets {
#             subnet = "10.250.2.128/25"
#             virtual_subnet = "10.250.2.128/25"
#           }
#           connector_settings {
#             bandwidth_range {
#                 min = 3
#                 max = 5
#             }
#           }
#         }

#         connect_type = "connector"

#     }
#     policies = ["ALLOW-ALL-NETWORKS"]
#     onboard_app = true
#     decommission_app = false
# }

# #AWS with transit gateway and none
resource "prosimo_network_onboarding" "testapp-AWS-public" {

    name = "demo_network_aws"
    namespace = "default"
    network_exportable_policy = false
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-aws-app-iam"
        region_name = "us-east-1"
        cloud_networks {
          vpc = "vpc-09654a32f8e0aebdd"
          # hub_id = ""
          connector_placement = "Workload VPC"
          connectivity_type = "vpc-peering"
          service_insertion_endpoint_subnets = "auto"
          subnets {
            subnet = "10.40.0.0/24"
          }
          connector_settings {
            connector_subnets = ["10.40.0.0/24"]
            bandwidth_range {
                min = 1
                max = 1
            }
          }
        }

        connect_type = "connector"

    }
    policies = ["ALLOW-ALL-NETWORKS"]
    internet_egress_controls = ["multimatch"]
    onboard_app = true
    decommission_app = false
    force_offboard = false
}


# #PrivateDC Network Onboarding
# resource "prosimo_network_onboarding" "privateDC" {
#   network_exportable_policy = false
#   name = "private-network-test"
#   private_cloud {
#     cloud_creds_name  = "PrivateDC"
#      subnets           = ["10.0.0.2/32"]
#   }
#   policies                = ["ALLOW-ALL-NETWORKS"]
#   onboard_app             = false
#   decommission_app        = false
# }

