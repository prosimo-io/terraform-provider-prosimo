####### User2App Agent Onboarding on bulk public websites
resource "prosimo_app_onboarding_fqdn" "Agent_bulk_webs" {
  count            = length(var.web_list)
  app_name         = "bulk_${var.web_list[count.index]}_tf"
  idp_name         = local.idp_name

  app_urls {
    domain_type        = "custom"
    subdomain_included = true
    app_fqdn           = var.web_list[count.index]
    protocols {
      protocol  = "tcp"
      port_list = var.port_list
    }
    health_check_info {
      enabled = true
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = local.cloud_creds_name
      edge_regions {
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_dns      = true
        dns_custom {
          dns_server             = var.pub_dns_server
          is_healthcheck_enabled = local.is_healthcheck_enabled
        }
      }
    }

  }
  saml_rewrite{
    selected_auth_type = "other" #"oidc"
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app

}


##### User2App Agent Onboarding on a single
resource "prosimo_app_onboarding_ip" "Agent-single-IP" {
  app_name         = var.agent_app_names[2]
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name

  app_urls {
    app_fqdn = "${var.vm_ips[1]}/32"
    protocols {
      protocol  = "tcp"
      port_list = var.port_list
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = local.cloud_creds_name
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = local.conn_option
      }
    }
  }

  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

### User2App Agent Onboarding on bulk internal application servers via private DNS server IP
resource "prosimo_app_onboarding_fqdn" "Agent_bulk_apps_via_DNS_IP" {
  depends_on       = [prosimo_app_onboarding_ip.Agent_DNS_svr_tf]
  count            = length(var.app_list)
  app_name         = "Agent_IP_bulk_${var.app_list[count.index]}_tf"

  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name

  app_urls {
    domain_type         = "custom"
    subdomain_included  = false
    app_fqdn            = var.app_list[count.index]
    protocols {
      protocol          = "tcp"
      port_list         = var.port_list
    }
    health_check_info { enabled = false }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = local.cloud_creds_name
      edge_regions {
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_dns      = true
        dns_custom {
          dns_server             = [var.vm_ips[0]]
          is_healthcheck_enabled = local.is_healthcheck_enabled
        }
      }
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}
data "prosimo_app_onboarding" "bulk_internal_apps" {
  depends_on     = [prosimo_app_onboarding_fqdn.Agent_bulk_apps_via_DNS_IP]
  count            = length(var.app_list)
  input_app_name         = "Agent_IP_bulk_${var.app_list[count.index]}_tf"
}

output "app_onboard_output_bulk_internal_apps_0" {
  description = "applist"
#  value       = data.prosimo_app_onboarding.bulk_internal_apps
#  value       = [for output in data.prosimo_app_onboarding.bulk_internal_apps[0].onboarded_apps : [for appurl in output.app_urls : appurl.papp_fqdn]]
  value       = data.prosimo_app_onboarding.bulk_internal_apps[0].onboarded_apps
}


#### User2App Agent Onboarding on a single fqdn server via private DNS server IP solution
resource "prosimo_app_onboarding_fqdn" "Agent_single_fqdn_dns_IP" {
  depends_on       = [prosimo_app_onboarding_ip.Agent_DNS_svr_tf]
  app_name         = var.agent_app_names[4]
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name

  app_urls {
    domain_type        = "custom"
    subdomain_included = false
    app_fqdn           = "client.${var.domain_2}" #"${var.agent_app_names[4]}.${var.domain_2}"
    protocols {
      protocol         = "tcp"
      #              port = 22
      port_list        = var.port_list
    }
    health_check_info {
           enabled = true
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = local.cloud_creds_name
      edge_regions {
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_dns      = true
        dns_custom {
          dns_server             = [var.vm_ips[0]]
          is_healthcheck_enabled = local.is_healthcheck_enabled
        }
      }
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

#### User2App Agent Onboarding on private DNS server
resource "prosimo_app_onboarding_ip" "Agent_DNS_svr_tf" {
  app_name         = var.agent_app_names[5] #"agent_DNS_Server_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    app_fqdn       = "${var.vm_ips[0]}/32"
    protocols {
      protocol     = "dns"
    }

    cloud_config {
      connection_option             = local.connection_option
      cloud_creds_name              = local.cloud_creds_name
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = local.conn_option
      }
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

##### User2App Agent Onboarding on a single fqdn server via above DNS server app
resource "prosimo_app_onboarding_fqdn" "Agent_single_fqdn_dns_app" {
  depends_on       = [prosimo_app_onboarding_ip.Agent_DNS_svr_tf]
  app_name         = var.agent_app_names[6] #"Agent_single_fqdn_dns_app_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name

  app_urls {
    domain_type        = "custom"
    subdomain_included = false
    app_fqdn           = "dns-agent.${var.domain_2}"#"${var.agent_app_names[6]}.${var.domain_2}"
    protocols {
      protocol  = "tcp"
      #              port = 22
      port_list = var.port_list
    }
        health_check_info {
           enabled = true
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = local.cloud_creds_name
      edge_regions {
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_dns      = true
        dns_custom {
          dns_app                = var.agent_app_names[5]
          is_healthcheck_enabled = local.is_healthcheck_enabled
        }
      }
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

### User2App Agent Onboarding on bulk fqdn deployment via DNS server App
resource "prosimo_app_onboarding_fqdn" "Agent_DNS_bulk" {
  depends_on       = [prosimo_app_onboarding_ip.Agent_DNS_svr_tf]
  app_name         = var.agent_app_names[7] #"Agent_DNS_bulk_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    domain_type        = "custom"
    subdomain_included = true
    app_fqdn           = "${var.domain_2}"
    protocols {
      protocol  = "tcp"
      port_list = ["1-6000"] #var.port_list
    }
    health_check_info {
           enabled = true
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = local.cloud_creds_name
      edge_regions {
        backend_ip_address_discover = local.backend_ip_address_discover
        backend_ip_address_dns      = true
        dns_custom {
          dns_app                = var.agent_app_names[5]
          is_healthcheck_enabled = local.is_healthcheck_enabled
        }
      }
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}
#
######## User2App Agent Onboarding on an Azure ssh server via vWANHub connection---pass on 9/2/22
resource "prosimo_app_onboarding_ip" "agent_vwanhub_ssh_https" {
  app_name         = var.agent_app_names[8] #"vwan_ssh_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    app_fqdn = "${var.vm_ips[4]}/32"
    protocols {
      protocol  = "tcp"
      port_list = var.port_list
    }
    cloud_config {
      connection_option = local.connection_option
      cloud_creds_name  = var.az_app_account
      edge_regions {
        region_name                 = var.az_region
        region_type                 = local.region_type
        conn_option                 = "vwanHub"
        app_network_id              = var.vnet
        attach_point_id             = var.hub
      }
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agent
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}

####### User2App Agent Onboarding on an Azure subnets via transit vNET connection
#resource "prosimo_app_onboarding_ip" "agent_transit_vnet_subnets" {
#  app_name         = var.agent_app_names[9]
#  ip_pool_cidr     = local.ip_pool_cidr
#  idp_name         = local.idp_name
#  app_urls {
#    app_fqdn       = "${var.vm_ips[5]}/32" #var.vnet_ciders
#    protocols {
#      protocol     = "tcp"
#      port_list    = var.port_list
#    }
#    cloud_config {
#      connection_option = local.connection_option
#      cloud_creds_name  = var.az_app_account
#      edge_regions {
#        region_name                 = var.az_region
#        region_type                 = local.region_type
#        conn_option                 = "azureTransitVnet"
#        #        backend_ip_address_discover = true
##        app_network_id              = "${var.az_resouce}${var.az_spoke_1}, ${var.az_resouce}${var.az_spoke_2}"
#        attach_point_id             = "${var.az_resouce}${var.az_hub}"
#
#      }
#    }
#  }
#  optimization_option       = local.optimization_option
#  enable_multi_cloud_access = local.enable_multi_cloud_access
#  policy_name               = var.policy_name_agent
#  onboard_app               = local.onboard_app
#  decommission_app          = local.decommission_app
#}

### User2App Agent Onboarding on subnets via AWS TGW connection, the subnet has a DNS server in it.
#####NOTE: before connecting, make sure to goto Resource Access Manager to share TGW with the Infra account
resource "prosimo_app_onboarding_ip" "Agent_TGW" {
  app_name         = var.agent_app_names[10] #"Agent_subnets_TGW_DNS_IP_tf"
  ip_pool_cidr     = local.ip_pool_cidr
  idp_name         = local.idp_name
  app_urls {
    app_fqdn       = "${var.vm_ips[3]}/32"
    protocols {
      protocol     = "tcp"
      port_list    = var.port_list
    }
    cloud_config {
      connection_option             = local.connection_option
      cloud_creds_name              = local.cloud_creds_name
      edge_regions {
        region_name                 = local.region_name
        region_type                 = local.region_type
        conn_option                 = "transitGateway"
        app_network_id              = var.vpc_id_tgw
        attach_point_id             =  var.TGW_id
        tgw_app_routetable          = "MODIFY"
      }
    }
  }
  optimization_option       = local.optimization_option
  enable_multi_cloud_access = local.enable_multi_cloud_access
  policy_name               = var.policy_name_agentless
  onboard_app               = local.onboard_app
  decommission_app          = local.decommission_app
}


