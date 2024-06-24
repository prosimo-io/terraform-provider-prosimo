---
page_title: "prosimo_app_onboarding_web Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this resource to onboard web apps.
---

# prosimo_app_onboarding_web (Resource)

Use this resource to onboard web apps.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
# terraform {
#   required_providers {
#     prosimo = {
#       version = "1.0.0"
#       source  = "prosimo.io/prosimo/prosimo"
#     }
#   }
# }
# provider "prosimo" {
#     token = "K4K_JXv1WjAjchzAXAciKEqE-JhYwIVvt90YtVmfvkI="
#     insecure = true
#     base_url = "https://myevesachinsp1715634065479.dashboard.psonar.us/"
# }
resource "prosimo_app_onboarding_web" "azure-private-wordpress" {
    # depends_on = [prosimo_app_onboarding_web.azure-vhub-wordpress]
    app_name = "azure-privatelink-wordpress"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "wordpress-azure-1715634348750.myeventarena.com"
        domain_type = "custom"
        app_fqdn = "wordpress-azure-1715634348750.myeventarena.com"
        subdomain_included = false
        protocols {
            protocol = "https"
            port = 443
        }
        health_check_info {
          enabled = true
        }
        cloud_config {
            connection_option = "private"
            cloud_creds_name = "prosimo-app"
            edge_regions {
                region_name = "westus"
                region_type = "active"
                conn_option = "azurePrivateLink"
                backend_ip_address_discover = false
                backend_ip_address_manual = ["10.250.31.28"]
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
    enable_multi_cloud_access = true
    policy_name = ["ALLOW-ALL-USERS"]
    onboard_app = true
    decommission_app = false
    wait_for_rollout = false
}
resource "prosimo_app_onboarding_web" "azure-vhub-wordpress" {
    app_name = "azure-vhub-wordpress"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "wordpress-azure-1715635122046.myeventarena.com"
        domain_type = "custom"
        app_fqdn = "wordpress-azure-1715635122046.myeventarena.com"
        subdomain_included = false
        protocols {
            protocol = "https"
            port = 443
        }
        health_check_info {
          enabled = true
        }
        cloud_config {
            connection_option = "private"
            cloud_creds_name = "prosimo-app"
            edge_regions {
                region_name = "westus"
                region_type = "active"
                conn_option = "vwanHub"
                backend_ip_address_discover = false
                backend_ip_address_manual = ["13.87.133.10"]
                app_network_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/wordpress-azure-1715635122046-rg/providers/Microsoft.Network/virtualNetworks/wordpress-azure-1715635122046-vnet"
                attach_point_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/wordpress-azure-1715635122046-rg/providers/Microsoft.Network/virtualHubs/wordpress-azure-1715635122046-vhub"
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
    enable_multi_cloud_access = true
    policy_name = ["ALLOW-ALL-USERS"]
    onboard_app = true
    decommission_app = false
    wait_for_rollout = false
}
# resource "prosimo_app_onboarding_cloudsvc" "psonar-s3-app" {
#     app_name = "psonar-s3-app"
#     idp_name         = "azure_ad"
#     cloud_svc        = "amazon-s3"
#     app_urls {
#         internal_domain = "s3-aws-1715634065779.myeventarena.com"
#         domain_type = "custom"
#         app_fqdn = "s3-aws-1715634065779.myeventarena.com"
#         subdomain_included = false
#         health_check_info {
#           enabled = false
#         }
#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-aws-app-iam"
#             edge_regions {
#                 region_type = "active"
#                 region_name = "us-east-1"
#                 buckets = ["psonar-s3-wordpress"]
#             }
#         }
#         dns_service {
#             type = "manual"
#         }
#         ssl_cert {
#           generate_cert = true
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true
#     policy_name = ["ALLOW-ALL-USERS"]
#     onboard_app = true
#     decommission_app = false
# }
# resource "prosimo_app_onboarding_web" "aws-vpc_peering-wordpress" {
#     app_name = "aws-vpc_peering-wordpress"
#     idp_name = "azure_ad"
#     app_urls {
#         internal_domain = "wordpress-aws-1715648697124.myeventarena.com"
#         domain_type = "custom"
#         app_fqdn = "wordpress-aws-1715648697124.myeventarena.com"
#         subdomain_included = false
#         protocols {
#             protocol = "https"
#             port = 443
#         }
#         health_check_info {
#           enabled = true
#         }
#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-aws-app-iam"
#             edge_regions {
#                 region_name = "us-west-1"
#                 region_type = "active"
#                 conn_option = "private"
#                 backend_ip_address_discover = false
#                 backend_ip_address_manual = ["54.151.116.212"]
#             }
#         }
#         dns_service {
#             type = "manual"
#         }
#         ssl_cert {
#             generate_cert = true
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true
#     policy_name = ["ALLOW-ALL-USERS"]
#     onboard_app = true
#     decommission_app = false
# }
# resource "prosimo_app_onboarding_web" "aws-privatelink-wordpress" {
#     app_name = "aws-privatelink-wordpress"
#     idp_name = "azure_ad"
#     app_urls {
#         internal_domain = "wordpress-aws-1715640653389.myeventarena.com"
#         domain_type = "custom"
#         app_fqdn = "wordpress-aws-1715640653389.myeventarena.com"
#         subdomain_included = false
#         protocols {
#             protocol = "https"
#             port = 443
#         }
#         health_check_info {
#           enabled = true
#         }
#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-aws-app-iam"
#             edge_regions {
#                 region_name = "us-west-1"
#                 region_type = "active"
#                 conn_option = "awsPrivateLink"
#                 backend_ip_address_discover = true
#             }
#         }
#         dns_service {
#             type = "manual"
#         }
#         ssl_cert {
#             generate_cert = true
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true
#     policy_name = ["ALLOW-ALL-USERS"]
#     onboard_app = true
#     decommission_app = false
# }
# resource "prosimo_app_onboarding_web" "aws-vpn-gateway-wordpress" {
#     app_name = "aws-vpn-gateway-wordpress"
#     idp_name = "azure_ad"
#     app_urls {
#         internal_domain = "wordpress-aws-1715648954579.myeventarena.com"
#         domain_type = "custom"
#         app_fqdn = "wordpress-aws-1715648954579.myeventarena.com"
#         subdomain_included = false
#         protocols {
#             protocol = "https"
#             port = 443
#         }
#         health_check_info {
#           enabled = true
#         }
#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-aws-app-iam"
#             edge_regions {
#                 region_name = "us-west-1"
#                 region_type = "active"
#                 conn_option = "awsVpnGateway"
#                 backend_ip_address_discover = false
#                 backend_ip_address_manual = ["10.250.40.9"]
#                 app_network_id = "vpc-0204ddfc1848ea929"
#             }
#         }
#         dns_service {
#             type = "manual"
#         }
#         ssl_cert {
#             generate_cert = true
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true
#     policy_name = ["ALLOW-ALL-USERS"]
#     onboard_app = true
#     decommission_app = false
# }
# resource "prosimo_app_onboarding_web" "aws-tgw-wordpress" {
#     app_name = "aws-tgw-wordpress"
#     idp_name = "azure_ad"
#     app_urls {
#         internal_domain = "wordpress-aws-1715649214334.myeventarena.com"
#         domain_type = "custom"
#         app_fqdn = "wordpress-aws-1715649214334.myeventarena.com"
#         subdomain_included = false
#         protocols {
#             protocol = "https"
#             port = 443
#         }
#         health_check_info {
#           enabled = true
#         }
#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-aws-app-iam"
#             edge_regions {
#                 region_name = "us-west-1"
#                 region_type = "active"
#                 conn_option = "transitGateway"
#                 backend_ip_address_discover = false
#                 backend_ip_address_manual = ["10.250.41.60"]
#                 app_network_id = "vpc-0ebc3bb3cff0dbb09"
#                 attach_point_id = "tgw-0a1871b2a38c38c8f"
#             }
#         }
#         dns_service {
#             type = "manual"
#         }
#         ssl_cert {
#             generate_cert = true
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true
#     policy_name = ["ALLOW-ALL-USERS"]
#     onboard_app = true
#     decommission_app = false
# }
# resource "prosimo_app_onboarding_web" "common-app" {
#     app_name = "common-app"
#     idp_name = "azure_ad"
#     app_urls {
#         internal_domain = "app-gcp-us-west1-1715634065485.myeventarena.com"
#         domain_type = "custom"
#         app_fqdn = "app-gcp-us-west1-1715634065485.myeventarena.com"
#         subdomain_included = false
#         protocols {
#             protocol = "https"
#             port = 443
#         }
#         health_check_info {
#           enabled = false
#         }
#         cloud_config {
#             connection_option = "public"
#             cloud_creds_name = "prosimo-gcp-infra"
#             edge_regions {
#                 region_name = "us-west1"
#                 region_type = "active"
#                 conn_option = "public"
#                 backend_ip_address_discover = false
#                 backend_ip_address_manual = ["23.99.84.98"]
#             }
#         }
#         dns_service {
#             type = "manual"
#         }
#         ssl_cert {
#             generate_cert = true
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true
#     policy_name = ["ALLOW-ALL-USERS"]
#     onboard_app = true
#     decommission_app = false
# }
# resource "prosimo_app_onboarding_jumpbox" "gcp-peering-jumpbox-server" {
#     app_name = "gcp-peering-jumpbox-server"
#     idp_name         = "azure_ad"
#     app_urls {
#         internal_domain = "jumpbox-server-gcp-1715637227320.myeventarena.com"
#         domain_type = "custom"
#         app_fqdn = "jumpbox-server-gcp-1715637227320.myeventarena.com"
#         subdomain_included = false
#         health_check_info {
#           enabled = false
#         }
#         cloud_config {
#             connection_option = "private"
#             cloud_creds_name = "prosimo-gcp-app"
#             edge_regions {
#                 region_type = "active"
#                 region_name = "us-west1"
#                 conn_option = "private"
#                 backend_ip_address_discover = false
#                 backend_ip_address_manual = ["35.247.55.92"]
#             }
#         }
#         dns_service {
#             type = "manual"
#         }
#         ssl_cert {
#           generate_cert = true
#         }
#     }
#     optimization_option = "PerformanceEnhanced"
#     enable_multi_cloud_access = true
#     policy_name = ["ALLOW-ALL-USERS"]
#     onboard_app = true
#     decommission_app = false
# }
resource "prosimo_app_onboarding_web" "aws-privatelink-wordpress" {
    app_name = "aws-privatelink-wordpress"
    idp_name = "azure_ad"
    app_urls {
        internal_domain = "wordpress-aws-1715640653389.myeventarena.com"
        domain_type = "custom"
        app_fqdn = "wordpress-aws-1715640653389.myeventarena.com"
        subdomain_included = false
        protocols {
            protocol = "https"
            port = 443
        }
        health_check_info {
          enabled = true
        }
        cloud_config {
            connection_option = "public"
            cloud_creds_name = "prosimo-aws-iam"
            edge_regions {
                region_name = "us-east-1"
                region_type = "active"
                conn_option = "public"
                backend_ip_address_discover = false
                backend_ip_address_manual = ["23.99.84.98"]
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
    enable_multi_cloud_access = true
    policy_name = ["ALLOW-ALL-USERS"]
    onboard_app = false
    decommission_app = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_name` (String) name for the application
- `app_urls` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--app_urls))
- `decommission_app` (Boolean) Set this to true if you would like app to be offboarded from fabric
- `onboard_app` (Boolean) Set this to true if you would like app to be onboarded to fabric
- `optimization_option` (String) Optimization option for app: e.g: CostSaving, PerformanceEnhanced, FastLane
- `policy_name` (List of String) Select policy name.e.g: ALLOW-ALL-USERS, DENY-ALL-USERS or CUSTOMIZE.Conditional access policies and Web Application Firewall policies for the application.

### Optional

- `client_cert` (String) Client Cert details
- `customize_policy` (Block Set, Max: 1) Choose any custom policy created from the policy library or create one. (see [below for nested schema](#nestedblock--customize_policy))
- `enable_multi_cloud_access` (Boolean) Setting this to true would leverage multi clouds to optimize the app performance
- `force_offboard` (Boolean) Force app offboarding incase of normal offboarding failure.
- `idp_name` (String) IDP provider name.
- `saml_rewrite` (Block Set, Max: 1) App authentication option while selecting prosimo domain (see [below for nested schema](#nestedblock--saml_rewrite))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `wait_for_rollout` (Boolean) Wait for the rollout of the task to complete. Defaults to true.

### Read-Only

- `app_access_type` (String) e.g: Agent or Agentless
- `app_type` (String) e.g: type of app onboarded, e.g: citrix, web, fqdn, jumpbox
- `id` (String) The ID of this resource.

<a id="nestedblock--app_urls"></a>
### Nested Schema for `app_urls`

Required:

- `app_fqdn` (String) Fqdn of the app that user would access after onboarding
- `cloud_config` (Block Set, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--app_urls--cloud_config))
- `domain_type` (String) Type of Domain: e.g custom or prosimo
- `health_check_info` (Block Set, Min: 1, Max: 1) Application health check config from edge (see [below for nested schema](#nestedblock--app_urls--health_check_info))
- `internal_domain` (String) server domain name or IP
- `protocols` (Block Set, Min: 1, Max: 1) Protocol that prosimo edge uses to connect to App (see [below for nested schema](#nestedblock--app_urls--protocols))

Optional:

- `cache_rule` (String) Cache Rules for your App Domains
- `dns_custom` (Block Set, Max: 1) Custom DNS setup (see [below for nested schema](#nestedblock--app_urls--dns_custom))
- `dns_service` (Block Set, Max: 1) In order to enable users to access an application using the external domain via the Prosimo fabric, you need to set up a new canonical name (CNAME) record redirect in your origin domain name system (DNS) record. (see [below for nested schema](#nestedblock--app_urls--dns_service))
- `ssl_cert` (Block Set, Max: 1) set up secure communication between the user and the application via the fabric, there are 3 options: Upload a Certificate, Generate a new certificate or Use an existing certificate (see [below for nested schema](#nestedblock--app_urls--ssl_cert))
- `subdomain_included` (Boolean) Set True to onboard subdomains of the application else False
- `waf_policy_name` (String) WAF Policies for your App Domains, applicable when the Edge to App Protocol is either HTTP or HTTPS.

Read-Only:

- `id` (String) The ID of this resource.

<a id="nestedblock--app_urls--cloud_config"></a>
### Nested Schema for `app_urls.cloud_config`

Required:

- `cloud_creds_name` (String) cloud account under which application is hosted
- `connection_option` (String) Public, if the app domain has a public IP address / DNS A record on the internet currently, and the Prosimo Edge should connect to the application using a public connection.Private, if the application only has a private IP address, and Edge should connect to it over a private connection.

Optional:

- `app_hosted_type` (String) Wheather app is hosted in Public cloud like AWS/AZURE/GCP or private DC. Available options PRIVATE/PUBLIC
- `dc_app_ip` (String) Applicable only if  app_hosted_type is PRIVATE, IP of the app hosted in PRIVATE DC
- `edge_regions` (Block List) (see [below for nested schema](#nestedblock--app_urls--cloud_config--edge_regions))
- `has_private_connection_options` (Boolean)
- `is_show_connection_options` (Boolean)

<a id="nestedblock--app_urls--cloud_config--edge_regions"></a>
### Nested Schema for `app_urls.cloud_config.edge_regions`

Required:

- `backend_ip_address_discover` (Boolean) if Set to true, auto discoverers available endpoints

Optional:

- `app_network_id` (String) App network id details
- `attach_point_id` (String) Attach Point id details
- `backend_ip_address_dns` (Boolean)
- `backend_ip_address_manual` (List of String) Pass endpoints manually.
- `conn_option` (String) Connection option for private connection: e.g: peering/transitGateway/awsPrivateLink/azurePrivateLink/azureTransitVnet/vwanHub
- `region_name` (String) Name of the region where app is available
- `region_type` (String) Type of region: e.g:active, backup etc
- `tgw_app_routetable` (String)



<a id="nestedblock--app_urls--health_check_info"></a>
### Nested Schema for `app_urls.health_check_info`

Optional:

- `enabled` (Boolean)
- `endpoint` (String) HealthCheck Endpoints


<a id="nestedblock--app_urls--protocols"></a>
### Nested Schema for `app_urls.protocols`

Required:

- `port` (Number) target port number
- `protocol` (String) Protocol type, e.g: “http”, “https”, “ssh”, “vnc”, or “rdp

Optional:

- `paths` (List of String) Customized websocket paths
- `web_socket_enabled` (Boolean) Set to true if tou would like prosimo edges to communicate with app via websocket

Read-Only:

- `is_valid_protocol_port` (Boolean)


<a id="nestedblock--app_urls--dns_custom"></a>
### Nested Schema for `app_urls.dns_custom`

Optional:

- `dns_app` (String) DNS App name
- `dns_server` (List of String) DNS Server List
- `is_healthcheck_enabled` (Boolean) Health check to ensure application domains being resolved by dns servers


<a id="nestedblock--app_urls--dns_service"></a>
### Nested Schema for `app_urls.dns_service`

Required:

- `type` (String) Type of DNS service: e.g: manual, route 53, prosimo

Optional:

- `aws_route53_cloud_creds_name` (String) Cloud creds for route 53


<a id="nestedblock--app_urls--ssl_cert"></a>
### Nested Schema for `app_urls.ssl_cert`

Optional:

- `existing_cert` (String) Select from already existing certificates(In Certificate TAB)
- `generate_cert` (Boolean) Set this to true if you want prosimo to generate new certificates
- `upload_cert` (Block Set) Upload the certificate if the certificates are already available for application (see [below for nested schema](#nestedblock--app_urls--ssl_cert--upload_cert))

<a id="nestedblock--app_urls--ssl_cert--upload_cert"></a>
### Nested Schema for `app_urls.ssl_cert.upload_cert`

Optional:

- `cert_path` (String) Path to certificate
- `private_key_path` (String) Path to private key




<a id="nestedblock--customize_policy"></a>
### Nested Schema for `customize_policy`

Optional:

- `name` (String)


<a id="nestedblock--saml_rewrite"></a>
### Nested Schema for `saml_rewrite`

Optional:

- `metadata` (String) Required while selecting SAML based authentication
- `metadata_url` (String) Required while selecting SAML based authentication
- `selected_auth_type` (String) Type of authentication: e.g. SAML, OIDC, Others


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)

