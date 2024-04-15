
# Agent IP Address
resource "prosimo_app_onboarding_dns" "AgentlessAppOnboarding" {

    app_name = "agent-tgw"
    # idp_name = "azure_ad"
    app_urls {
        app_fqdn = "10.1.1.1"
        service_ip_type = "auto"
        service_ip = "100.127.10.1"
        protocols {
            protocol = "dns"
            port_list = ["53"]
        }

        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-aws-iam"
            edge_regions {
                region_type = "active"
                region_name = "us-west-1"
                conn_option = "public"
                # app_network_id = "vpc-067648a8b45cd2369"
                # attach_point_id = "tgw-02c98e2afd03758d7"
                # tgw_app_routetable = "MODIFY"
                
            }
        }
    }
    optimization_option = "PerformanceEnhanced"

    policy_name = ["ALLOW-ALL-USERS"]

    onboard_app = false
    decommission_app = false
     force_offboard = false
}

# resource "prosimo_app_onboarding_dns" "AgentlessAppOnboarding_new" {

#     app_name = "agent-vwanhub-ssh-https-tf"
#     # ip_pool_cidr = "192.168.0.0/22"
#     idp_name = "azure_ad"
#     app_urls {
#         app_fqdn = "172.17.1.4/32"
#         protocols {
#             protocol = "tcp"
#             port_list = ["443", "22"]
#         }

#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-app"
#             edge_regions {
#                 region_type = "active"
#                 region_name = "eastus"
#                 conn_option = "vwanHub"
#                 app_network_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vm-1-rg/providers/Microsoft.Network/virtualNetworks/qing-vm-1-rg-vnet"
#                 attach_point_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-1"
#                 # tgw_app_routetable = "MODIFY"
                
#             }
#         }
#     }
#     optimization_option = "PerformanceEnhanced"

#     policy_name = [ "ALLOW-ALL"]

#     onboard_app = true
#     decommission_app = false
# }

# resource "prosimo_app_onboarding_dns" "AgentlessAppOnboarding_TransitVnet" {

#     app_name = "agent-transit-vnet-subnets-tf"
#     # ip_pool_cidr = "192.168.8.0/22"
#     idp_name = "azure_ad"
#     app_urls {
#         app_fqdn = "172.19.1.4"
#         protocols {
#             protocol = "tcp"
#             port_list = ["443", "22"]
#         }

#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-app"
#             edge_regions {
#                 region_type = "active"
#                 region_name = "eastus2"
#                 conn_option = "azureTransitVnet"
#                 app_network_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing_vtransit_eastus2_rg/providers/Microsoft.Network/virtualNetworks/qing_vtransit__eastus2_spoke1"
#                 attach_point_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing_vtransit_eastus2_rg/providers/Microsoft.Network/virtualNetworks/qing_vtransit_eastus2_hub"
#                 # tgw_app_routetable = "MODIFY"
                
#             }
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true

#     policy_name =[ "ALLOW-ALL"]

#     onboard_app = false
#     decommission_app = false
# }

# # Onboarding for app hosted in privateDC.
# resource "prosimo_app_onboarding_dns" "private-dc" {

#     app_name = "common-app-private"
#     # idp_name = "azure_ad"
#     app_urls {
#         # domain_type = "custom"
#         app_fqdn = "10.0.0.1"
#         # subdomain_included = false

#         protocols {
#             protocol = "tcp"
#             port_list = ["80", "90"]
#         }

#         # health_check_info {
#         #   enabled = false
#         # }

#         cloud_config {
#             app_hosted_type = "PRIVATE"
#             connection_option = "public"
#             cloud_creds_name = "PrivateDC"
#             dc_app_ip = "10.1.1.1"
#         }
#     }
#     # saml_rewrite{
#     #   selected_auth_type = "oidc"
#     # }
#     optimization_option = "PerformanceEnhanced"

#     policy_name = ["ALLOW-ALL-USERS"]

#     onboard_app = false
#     decommission_app = false
# }