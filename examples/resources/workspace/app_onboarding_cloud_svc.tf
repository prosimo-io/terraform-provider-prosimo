# CloudSVC(subdomain_included = false)
resource "prosimo_app_onboarding_cloudsvc" "testapp-s3" {

    app_name = "common-app-s3"
    cloud_svc = "amazon-s3"
    # ip_pool_cidr = "10.253.0.0/16"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "qingscnetworkers.info"
        app_fqdn = "qingscnetworkers.info"
        subdomain_included = false
        cloud_config {
            # connection_option = "private"
            cloud_creds_name = "prosimo-aws-app-iam"
            edge_regions {
                region_name = "us-east-2"
                region_type = "active"
                buckets = ["arun-vpc-s3"]
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
    force_offboard = true
}

# CloudSVC(subdomain_included = true)
resource "prosimo_app_onboarding_cloudsvc" "testapp-s4" {

    app_name = "common-app-s3"
    # cloud_svc = "amazon-s3"
    # ip_pool_cidr = "192.168.0.0/22"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "qingscnetworkers.info"
        app_fqdn = "qingscnetworkers.info"
        subdomain_included = true
        cloud_config {
            # connection_option = "private"
            cloud_creds_name = "prosimo-aws-app-iam"
            edge_regions {
            }
        }
        dns_custom {                                   
          # dns_app = "agent-DNS-Server-tf"
          dns_server = ["10.100.0.5"]
          is_healthcheck_enabled = true
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

# Onboarding for app hosted in privateDC.
resource "prosimo_app_onboarding_cloudsvc" "private-dc" {

    app_name = "common-app-private"
    app_urls {
        app_fqdn = "alex-app-101.abc.com"
        subdomain_included = false
        internal_domain = "alex-app-101.abc.com"


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
}


