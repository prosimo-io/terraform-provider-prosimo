
# Agent WEB access(subdomain_included = true, domain_type = "custom")
resource "prosimo_app_onboarding_fqdn" "AgentlessAppOnboarding" {

    app_name = "common-app-agent-fqdn"
    # idp_name = "azure_ad"
    app_urls {
        domain_type = "custom"
        app_fqdn = "alex-app-101.abc.com"
        subdomain_included = true

        protocols {
            protocol = "tcp"
            port_list = ["80", "90"]
        }
        protocols {
            protocol = "udp"
            port_list = ["80", "90"]
        }

        health_check_info {
          enabled = false
        }

        cloud_config {
            connection_option = "private"
            cloud_creds_name = "prosimo-gcp-infra"
            edge_regions {
                backend_ip_address_discover = false
                backend_ip_address_dns = false
                dns_custom {                                   
                    dns_server = ["8.8.8.8"]
                    is_healthcheck_enabled = true
                  }
            }
        }
    }
    saml_rewrite{
      selected_auth_type = "oidc"
    }
    optimization_option = "PerformanceEnhanced"

    policy_name =[ "ALLOW-ALL-USERS"]

    onboard_app = false
    decommission_app = false
     force_offboard = false
}



# Agent WEB access(subdomain_included = false, domain_type = "custom", backend_ip_address_dns = true)
# resource "prosimo_app_onboarding_fqdn" "AgentlessAppOnboarding" {

#     app_name = "common-app-agent-fqdn"
#     idp_name = "azure_ad"
#     app_urls {
#         domain_type = "custom"
#         app_fqdn = "alex-app-101.abc.com"
#         subdomain_included = false

#         protocols {
#             protocol = "tcp"
#             port_list = ["80", "90"]
#         }

#         health_check_info {
#           enabled = false
#         }

#         cloud_config {
#             connection_option = "public"
#             cloud_creds_name = "prosimo-gcp-infra"
#             edge_regions {
#                 region_type = "active"
#                 region_name = "us-west2"
#                 conn_option = "public"
#                 backend_ip_address_discover = false
#                 backend_ip_address_manual = ["23.99.84.98"]
#                 # dns_custom {                                   
#                 #     dns_app = "agent-DNS-Server-tf"
#                 #     is_healthcheck_enabled = true
#                 #   }
#             }
#         }
#     }
#     saml_rewrite{
#       selected_auth_type = "oidc"
#     }
#     optimization_option = "PerformanceEnhanced"

#     policy_name = ["ALLOW-ALL-USERS"]

#     onboard_app = false
#     decommission_app = false
# }


# Onboarding for app hosted in privateDC.
resource "prosimo_app_onboarding_fqdn" "private-dc" {

    app_name = "common-app-private"
    app_urls {
        domain_type = "custom"
        app_fqdn = "alex-app-101.abc.com"
        subdomain_included = false

        protocols {
            protocol = "tcp"
            port_list = ["80", "90"]
        }

        health_check_info {
          enabled = false
        }

        cloud_config {
            app_hosted_type = "PRIVATE"
            connection_option = "public"
            cloud_creds_name = "PrivateDC"
            dc_app_ip = "10.1.1.2"
        }
    }
    saml_rewrite{
      selected_auth_type = "oidc"
    }
    optimization_option = "PerformanceEnhanced"

    policy_name = ["DENY-ALL-USERS"]

    onboard_app = false
    decommission_app = false
}