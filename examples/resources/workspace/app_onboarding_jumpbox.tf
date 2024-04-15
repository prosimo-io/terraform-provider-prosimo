

# Jumpbox
resource "prosimo_app_onboarding_jumpbox" "testapp-jumpbox" {

    app_name = "common-app-jumpbox"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "tf-jumpbox.psonar.us"
        app_fqdn = "tf-jumpbox.psonar.us"
        cloud_config {
            connection_option = "private"
            cloud_creds_name = "prosimo-aws-app-iam"
            edge_regions {
                region_name = "us-east-2"
                region_type = "active"
                conn_option = "peering"
                backend_ip_address_discover = false
                backend_ip_address_manual = ["10.100.0.142"]
            }
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
}


# Onboarding for app hosted in privateDC.
resource "prosimo_app_onboarding_jumpbox" "private-dc" {

    app_name = "common-app-private"
    app_urls {
        internal_domain = "alex-app-101.abc.com"
        app_fqdn = "alex-app-101.abc.com"


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
