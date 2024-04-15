
# Agent IP Address (AWS)

resource "prosimo_app_onboarding_dns" "AgentlessAppOnboarding" {

    app_name = "dns-pub"
    app_urls {
        app_fqdn = "10.1.1.1"
        service_ip_type = "manual"
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
            }
        }
    }
    optimization_option = "PerformanceEnhanced"

    policy_name = ["ALLOW-ALL-USERS"]

    onboard_app = false
    decommission_app = false
}
resource "prosimo_app_onboarding_ip" "AgentlessAppOnboarding" {

    app_name = "agent-tgw"
    idp_name = "azure_ad"
    app_urls {
        app_fqdn = "10.10.0.36/32"
        service_ip_type = "auto"
        protocols {
            protocol = "dns"
            port_list = ["53"]
        }

        cloud_config {
            connection_option = "private"
            cloud_creds_name = "prosimo-aws-app-iam"
            edge_regions {
                region_type = "active"
                region_name = "us-east-2"
                conn_option = "transitGateway"
                app_network_id = "vpc-067648a8b45cd2369"
                attach_point_id = "tgw-02c98e2afd03758d7"
                tgw_app_routetable = "MODIFY"
                
            }
        }
    }
    optimization_option = "PerformanceEnhanced"

    policy_name = "ALLOW-ALL"

    onboard_app = false
    decommission_app = false
}


# Agent IP Address (AZURE with AzureTransitVnet)
resource "prosimo_app_onboarding_ip" "AgentlessAppOnboarding_TransitVnet" {

    app_name = "agent-transit-vnet-subnets-tf"
    # ip_pool_cidr = "192.168.8.0/22"
    idp_name = "azure_ad"
    app_urls {
        app_fqdn = "172.19.1.4"
        service_ip_type = "auto"
        protocols {
            protocol = "dns"
            port_list = ["53"]
        }

        cloud_config {
            connection_option = "private"
            cloud_creds_name = "prosimo-app"
            edge_regions {
                region_type = "active"
                region_name = "eastus2"
                conn_option = "azureTransitVnet"
                app_network_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing_vtransit_eastus2_rg/providers/Microsoft.Network/virtualNetworks/qing_vtransit__eastus2_spoke1"
                attach_point_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing_vtransit_eastus2_rg/providers/Microsoft.Network/virtualNetworks/qing_vtransit_eastus2_hub"
                # tgw_app_routetable = "MODIFY"
                
            }
        }
    }
    optimization_option = "PerformanceEnhanced"
    enable_multi_cloud_access = true

    policy_name = "ALLOW-ALL"

    onboard_app = false
    decommission_app = false
}

# Agent IP Address(app hosted in privateDC.)
resource "prosimo_app_onboarding_ip" "private-dc" {

    app_name = "common-app-private"
    app_urls {
        app_fqdn = "10.0.0.1"
        service_ip_type = "auto"
        protocols {
            protocol = "dns"
            port_list = ["53"]
        }

        cloud_config {
            app_hosted_type = "PRIVATE"
            connection_option = "public"
            cloud_creds_name = "PrivateDC"
            dc_app_ip = "10.1.1.1"
        }
    }
    optimization_option = "PerformanceEnhanced"

    policy_name = ["ALLOW-ALL-USERS"]

    onboard_app = false
    decommission_app = false
}

