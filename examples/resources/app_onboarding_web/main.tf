# #Agentless Web access (domain_type = "custom", subdomain_included = false)
resource "prosimo_app_onboarding_web" "agentless_multi_VMs" {

    app_name = "agentless-multi-VMs-tf"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "10.100.0.142"
        domain_type = "custom"
        app_fqdn = "alex-app-101.abc.com"
        subdomain_included = false

        protocols {
            protocol = "ssh"
            port = 22
        }

        health_check_info {
          enabled = true
        }

        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-aws-iam"
            edge_regions {
                region_type = "active"
                region_name = "us-west-1"
                conn_option = "public"
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

    onboard_app =false
    decommission_app = false
}

#Agentless Web access (domain_type = "Prosimo", subdomain_included = false)
resource "prosimo_app_onboarding_web" "AgentlessAppOnboarding-prosimo_domain" {

    app_name = "common-app-issue-new1"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "ssh-server-us-west2-1657650573897.myeventarena.com"
        domain_type = "prosimo"
        app_fqdn = "demo.access.myevekapildev1660666246049.scnetworkers.info"
        subdomain_included = false

        protocols {
            protocol = "ssh"
            port = 80
        }

        health_check_info {
          enabled = false
        }

        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-gcp-infra"
            edge_regions {
                region_type = "active"
                region_name = "us-west2"
                conn_option = "public"
                backend_ip_address_discover = false
                backend_ip_address_manual = ["23.99.84.98"]
            }
        }

        dns_service {
            type = "manual"
        }

        ssl_cert {
            upload_cert {
                cert_path = "path/to/certificate"
                private_key_path = "path/to/key"
            }
        }
    }

    optimization_option = "PerformanceEnhanced"

    policy_name = ["ALLOW-ALL-USERS"]

    onboard_app =false
    decommission_app = false
}

# Agentless Web access (domain_type = "Prosimo", subdomain_included = true)
resource "prosimo_app_onboarding_web" "AgentlessAppOnboarding-prosimo_domain_bulk" {

    app_name = "common-app-issue-new2"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "ssh-server-us-west2-1657650573897.myeventarena.com"
        domain_type = "prosimo"
        app_fqdn = "demo.access.myevekapildev1660666246049.scnetworkers.info"
        subdomain_included = true

        protocols {
            protocol = "ssh"
            port = 80
        }

        health_check_info {
          enabled = false
        }

        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-gcp-infra"
            edge_regions {
                backend_ip_address_discover = false
            }
        }

        dns_custom {
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

    onboard_app =false
    decommission_app = false
}


# Agentless Web access (domain_type = "Custom", subdomain_included = true)
resource "prosimo_app_onboarding_web" "AgentlessAppOnboarding-custom_domain_bulk" {

    app_name = "common-app-issue-new3"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "ssh-server-us-west2-1657650573897.myeventarena.com"
        domain_type = "custom"
        app_fqdn = "ssh-server-us-west2-1657650573897.myeventarena.com"
        subdomain_included = true

        protocols {
            protocol = "ssh"
            port = 80
        }

        health_check_info {
          enabled = false
        }

        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-gcp-infra"
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
            upload_cert {
                cert_path = "path/to/cert"
                private_key_path = "path/to/key"
            }
        }
    }

    optimization_option = "PerformanceEnhanced"

    policy_name = ["ALLOW-ALL-USERS"]

    onboard_app =false
    decommission_app = false
}

# Agentless Web access(app hosted in privateDC.)
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
}


