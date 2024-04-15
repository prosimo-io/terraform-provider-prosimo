# #Agentless Web access (domain_type = "custom", subdomain_included = false)
resource "prosimo_app_onboarding_web" "agentless_multi_VMs" {

    app_name = "agentless-multi-VMs-tf"
    # idp_name = "azure_ad"
    app_urls {
        internal_domain = "google.com"
        domain_type = "custom"
        app_fqdn = "google.com"
        subdomain_included = false

        protocols {
            protocol = "http"
            port = 80
        }

        health_check_info {
          enabled = true
        }

        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-aws-iam"
            edge_regions {
                region_type = "active"
                region_name = "us-east-2"
                conn_option = "public"
                backend_ip_address_discover = false
                backend_ip_address_manual = ["173.194.202.138"]
            }
        }
        dns_service {
            type = "manual"
        }

        ssl_cert {
          generate_cert = true
        }
        cache_rule = "Jenkins"
    }
    optimization_option = "CostSaving"

    policy_name = ["ALLOW-ALL-USERS"]

    onboard_app = false
    decommission_app = false
}

# Onboarding for app hosted in privateDC.
resource "prosimo_app_onboarding_web" "private-dc" {

    app_name = "common-app-private"
    app_urls {
        subdomain_included = false
        domain_type = "custom"
        internal_domain = "alex-app-101.abc.com"
        app_fqdn = "alex-app-101.abc.com"

        protocols {
            protocol = "ssh"
            port = 22
        }

        health_check_info {
          enabled = true
        }


        cloud_config {
            app_hosted_type = "PRIVATE"
            connection_option = "public"
            cloud_creds_name = "PrivateDC"
            dc_app_ip = "10.1.1.1"
        }
                dns_service {
            type = "manual"
         }
        ssl_cert {
            generate_cert = true
        }
    }

    optimization_option = "PerformanceEnhanced"

    policy_name = ["ALLOW-ALL-USERS"]

    onboard_app = false
    decommission_app = false
     force_offboard = false
}