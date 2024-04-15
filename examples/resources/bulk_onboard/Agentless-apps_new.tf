
resource "prosimo_app_onboarding_web" "agentless_multi_VMs" {
  app_name     = var.agentless_app_names[2]
  ip_pool_cidr = local.ip_pool_cidr
  idp_name     = local.idp_name
  app_urls {
    internal_domain    = var.vm_ips[0] #"app-gcp-us-west2-1661340185614.myeventarena.com"
    domain_type        = local.domain_type #"custom"
    app_fqdn           = "multi-VM-https${var.prosimo_domain}"
    #"app-gcp-us-west2-1661340185614.myeventarena.com"
    subdomain_included = false #true
    protocols {
      protocol = "https"
      port     = 443
    }
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option = local.connection_option #"public"
      cloud_creds_name  = local.cloud_creds_name #"prosimo-gcp-infra"
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = local.conn_option #"public"
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_manual = [var.vm_ips[0]]
      }
    }
    dns_service {
      type = "manual"
    }
    ssl_cert { generate_cert = false }
  }
  app_urls {
    internal_domain    = var.vm_ips[1]
    domain_type        = local.domain_type #"custom"
    app_fqdn           = "${var.agentless_app_names[2]}${var.prosimo_domain}"
    subdomain_included = false #true
    protocols {
      protocol = "ssh"
      port     = 22
    }
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option = local.connection_option #"public"
      cloud_creds_name  = local.cloud_creds_name #"prosimo-gcp-infra"
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = local.conn_option
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_manual = [var.vm_ips[1]]
      }
    }
    dns_service {
      type = "manual"
    }
    ssl_cert {
      generate_cert = true
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agentless
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}


### User2App Agentless Onboarding on URL rewriting for ssh
resource "prosimo_app_onboarding_web" "url_rewrite_ssh" {
  app_name         = var.agentless_app_names[3] #"url_rewrite_ssh_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    internal_domain    = "url-rewrite.${var.domain_1}"
    domain_type        = local.domain_type
    app_fqdn           = "${var.agentless_app_names[3]}${var.prosimo_domain}"
    subdomain_included = false
    protocols {
      protocol = "ssh"
      port     = 22
    }
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option             = local.connection_option
      cloud_creds_name              = local.cloud_creds_name
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = local.conn_option
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_manual = [var.vm_ips[1]]
      }
    }
    ssl_cert {
      generate_cert = false
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agentless
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

## User2App Agentless Onboarding on URL rewriting for https
resource "prosimo_app_onboarding_web" "url_rewrite_https" {
  app_name         = var.agentless_app_names[4] #"url_rewrite_https_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    internal_domain    = "url-rewrite.${var.domain_1}"
    domain_type        = local.domain_type
    app_fqdn           = "${var.agentless_app_names[4]}${var.prosimo_domain}"
    subdomain_included = false
    protocols {
      protocol = "https"
      port     = 443
    }
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option             = local.connection_option
      cloud_creds_name              = local.cloud_creds_name
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = local.conn_option
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_manual = [var.vm_ips[1]]
      }
    }
    ssl_cert {
      generate_cert = false
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agentless
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}
#
#### User2App Agentless Onboarding on AWS jumpbox
resource "prosimo_app_onboarding_jumpbox" "app_jumpbox" {
  app_name         = var.agentless_app_names[5] #"jumpbox_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    internal_domain    = "${var.agentless_app_names[5]}.${var.domain_1}"
    domain_type        = "custom"
    app_fqdn           = "${var.agentless_app_names[5]}.${var.domain_1}"
    subdomain_included = false
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option             = local.connection_option
      cloud_creds_name              = local.cloud_creds_name
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = local.conn_option
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_manual   = [var.vm_ips[1]]
      }
    }
    dns_service {
      type = "manual"
    }
    ssl_cert {
      generate_cert = true
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agentless
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

####get S3 bucket list
data "prosimo_s3bucket" "s3_bucket" {
  input_nickname    = local.cloud_creds_name
  input_region      = local.region_name
}
# User2App Agentless Onboarding on websites in AWS S3
resource "prosimo_app_onboarding_cloudsvc" "app_s3" {
  app_name         = var.agentless_app_names[6] #"aws_s3_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  cloud_svc       ="amazon-s3"
  app_urls {
    internal_domain    = "${var.agentless_app_names[6]}.${var.domain_1}"
    domain_type        = "custom"
    app_fqdn           = "${var.agentless_app_names[6]}.${var.domain_1}"
    subdomain_included = false
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option             = local.connection_option
      cloud_creds_name              = local.cloud_creds_name
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
#        backend_ip_address_discover = local.backend_ip_address_discover
        buckets                     = data.prosimo_s3bucket.s3_bucket.data
      }
    }
    dns_service {
      type = "manual"
    }
    ssl_cert {
      generate_cert = true
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agentless
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}
#
##### User2App Agentless Onboarding on a jumpbox server via AWS TGW connection
#########NOTE: before connecting, make sure to goto Resource Access Manager to share TGW with the Infra account
#resource "prosimo_app_onboarding_jumpbox" "app_jumpbox_TGW" {
#  app_name         = var.agentless_app_names[7] #"tgw_jumpbox_tf"
#  ip_pool_cidr     = local.ip_pool_cidr
#  idp_name         = local.idp_name
#  app_urls {
#    internal_domain    = "${var.agentless_app_names[7]}.${var.domain_1}"
#    domain_type        = "custom"
#    app_fqdn           = "${var.agentless_app_names[7]}.${var.domain_1}"
#    subdomain_included = false
#    health_check_info {
#      enabled = true
#    }
#    cloud_config {
#      connection_option             = local.connection_option
#      cloud_creds_name              = local.cloud_creds_name
#      edge_regions {
#        region_name                 = local.region_name
#        region_type                 = local.region_type
#        backend_ip_address_discover = local.backend_ip_address_discover
#        conn_option                 = "transitGateway"
#        backend_ip_address_manual   = [var.vm_ips[3]]
#        app_network_id              = var.vpc_id_tgw
#        attach_point_id             = var.TGW_id
#        tgw_app_routetable          = "MODIFY"
#      }
#    }
#    dns_service {
#      type = "manual"
#    }
#    ssl_cert {
#      generate_cert = true
#    }
#  }
#  optimization_option       = local.optimization_option
#  enable_multi_cloud_access = local.enable_multi_cloud_access
#  policy_name               = var.policy_name_agentless
#  onboard_app               = local.onboard_app
#  decommission_app          = local.decommission_app
#}

### User2App Agentless Onboarding on an AWS ssh server via PrivateLink
resource "prosimo_app_onboarding_web" "ssh_PrivateLink" {
  app_name         = var.agentless_app_names[8] #"ssh_privatrlink_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    internal_domain    = var.vm_ips[2]
    domain_type        = local.domain_type
    app_fqdn           = "${var.agentless_app_names[8]}${var.prosimo_domain}"
    subdomain_included = false
    protocols {
      protocol = "ssh"
      port     = 22
    }
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = local.cloud_creds_name
      edge_regions {
        region_name                 = "us-east-2" #local.region_name
        region_type                 = local.region_type
        backend_ip_address_discover = local.backend_ip_address_discover
        conn_option                 = "awsPrivateLink"
        backend_ip_address_manual   = [var.NLB_name]
        app_network_id              = var.vpc_id_nlb
      }
    }
    ssl_cert {
      generate_cert = false
    }
  }
  optimization_option = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name         = var.policy_name_agentless
  onboard_app         = local.onboard_app
  decommission_app    = local.decommission_app
}


##### User2App Agentless Onboarding on an Azure ssh server via vWANHub connection
resource "prosimo_app_onboarding_web" "vwanhub_ssh" {
  app_name         = var.agentless_app_names[9] #"vwan_ssh_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    internal_domain    = var.vm_ips[4]  #172.17.2.4
    domain_type        = local.domain_type
    app_fqdn           = "${var.agentless_app_names[9]}${var.prosimo_domain}"
    subdomain_included = false
    protocols {
      protocol = "ssh"
      port     = 22
    }
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = var.az_app_account
      edge_regions {
        region_name                 = var.az_region
        region_type                 = local.region_type
        conn_option                 = "vwanHub"
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_manual   = [var.vm_ips[4]]
        app_network_id              = var.vnet
        attach_point_id             = var.hub
      }
    }
    ssl_cert {
      generate_cert = true
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

## User2App Agentless Onboarding on Azure jumpbox via vWANHub connection
resource "prosimo_app_onboarding_jumpbox" "vwanhub_jumpbox" {
  app_name         = var.agentless_app_names[10] #"vwanhub_jumpbox_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    internal_domain    = "${var.agentless_app_names[10]}.${var.domain_1}"
    domain_type        = "custom"
    app_fqdn           = "${var.agentless_app_names[10]}.${var.domain_1}"
    subdomain_included = false
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = var.az_app_account
      edge_regions {
        region_name                 = "eastus2"
        region_type                 = local.region_type
        conn_option                 = "vwanHub"
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_manual   = [var.vm_ips[4]]
        app_network_id              = var.vnet
        attach_point_id             = var.hub
      }
    }
    dns_service {
      type = "manual"
    }
    ssl_cert {
      generate_cert = true
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

data "prosimo_app_onboarding" "R53_vwanhub" {
  depends_on     = [prosimo_app_onboarding_jumpbox.vwanhub_jumpbox]
  input_app_name = var.agentless_app_names[10]
}

output "app_onboard_output_vwanhub" {
  description = "applist"
  value       = [for output in data.prosimo_app_onboarding.R53_vwanhub.onboarded_apps : [for appurl in output.app_urls : appurl.papp_fqdn]]
}



