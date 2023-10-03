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

