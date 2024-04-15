
# Citrix VDI access (domain_type = "custom", subdomain_included = false)
resource "prosimo_app_onboarding_citrixvdi" "AgentlessAppOnboarding" {

    app_name = "common-app-new"
    idp_name = "azure_ad"
    citrix_ip = ["10.1.1.1"]
    app_urls {
        internal_domain = "google.com"
        domain_type = "custom"
        app_fqdn = "google.com"
        subdomain_included = true

        protocols {
            protocol = "http"
            port = 80
        }

        health_check_info {
          enabled = false
        }

        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-infra"
            edge_regions {
                backend_ip_address_discover = false
            }
        }
        dns_custom {                                   
          dns_app = "agent-DNS-Server-tf"
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

    onboard_app =false
    decommission_app = false
}




# Citrix VDI access  (domain_type = "Custom", subdomain_included = false)

resource "prosimo_app_onboarding_citrixvdi" "test" {

	app_name = "common-app"
	idp_name = "azure_ad"
	citrix_ip = ["10.1.1.1"]
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
			enabled = false
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-infra"
			edge_regions {
				region_type = "active"
				region_name = "westus2"
				conn_option = "public"
				backend_ip_address_discover = true
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

	onboard_app =false
	decommission_app = false
}


## Citrix VDI access (app hosted in privateDC.)
resource "prosimo_app_onboarding_citrixvdi" "private-dc" {

    app_name = "common-app-private"
    app_urls {
        domain_type = "custom"
        app_fqdn = "alex-app-101.abc.com"
        subdomain_included = false
        internal_domain = "alex-app-101.abc.com"

		protocols {
			protocol = "http"
			port = 80
		}

        health_check_info {
          enabled = false
        }

        cloud_config {
            app_hosted_type = "PRIVATE"
            connection_option = "public"
            cloud_creds_name = "PrivateDC"
            dc_app_ip = "10.1.1.5"
        }
    
		dns_service {
			type = "manual"
		}

		ssl_cert {
			generate_cert = true
		}
    }
    optimization_option = "PerformanceEnhanced"

    policy_name = ["DENY-ALL-USERS"]

    onboard_app = true
    decommission_app = false
}