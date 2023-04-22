---
page_title: "prosimo_network_onboarding Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this resource to onboard networks.
---

# prosimo_network_onboarding (Resource)

Use this resource to onboard networks.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
# Azure with VWAN hub
resource "prosimo_network_onboarding" "testapp-azure" {

    name = "demo_network_azure"
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-app"
        region_name = "eastus2"
        cloud_networks {
          vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/app-azure-eastus2-1661236757258-rg/providers/Microsoft.Network/virtualNetworks/app-azure-eastus2-1661236757258-vnet"
          hub_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"
          connectivity_type = "vwan-hub"
          connector_placement = "none"
          subnets = ["192.168.128.0/25"]
        }
        cloud_networks {
          vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Gitlab/providers/Microsoft.Network/virtualNetworks/Gitlab-vnet"
          hub_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"
          connectivity_type = "vwan-hub"
          connector_placement = "Workload VPC"
          subnets = ["10.3.5.0/24"]
        }
        connect_type = "connector"

    }
    onboard_app = false
    decommission_app = false
}


# # # Azure with VNET peering
resource "prosimo_network_onboarding" "testapp-s3" {

    name = "demo_network_new"
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-app"
        region_name = "eastus"
        cloud_networks {
          vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Gitlab/providers/Microsoft.Network/virtualNetworks/Gitlabvnet696"
          # hub_id = "tgw-04db5eac6fe3de45e"
           connector_placement = "Workload VPC"
          connectivity_type = "vnet-peering"
          subnets = ["10.3.7.0/24"]
        }
        connect_type = "connector"

    }
    policies = ["ProAccess-policy-tf"]
    onboard_app = false
    decommission_app = false
}

#AWS with transit gateway
resource "prosimo_network_onboarding" "testapp-s4" {

    name = "demo_network_new"
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-aws-app-iam"
        region_name = "us-west-2"
        cloud_networks {
          vpc = "vpc-033019a1cab5c5086"
          hub_id = "tgw-02b93ffe5733e94cd"
          connector_placement = "Infra VPC"
          connectivity_type = "transit-gateway"
          subnets = ["10.11.0.0/20"]
          connector_settings {
            bandwidth = "small"
            bandwidth_name = "<1 Gbps"
            instance_type = "t3.medium"
          }
        }

        connect_type = "connector"

    }
    policies = ["DENY-ALL-NETWORKS"]
    onboard_app = true
    decommission_app = false
}
#AWS with transit gateway and connector placement as none.
resource "prosimo_network_onboarding" "aws_u2n_euspoke3" {
  name = "aws-u2n-euspoke3-tf"
  public_cloud {
    cloud_type        = "public"
    connection_option = "private"
    cloud_creds_name  = "prosimo-aws-app-iam"
    region_name       = "eu-west-1"
    cloud_networks {
      hub_id              = "tgw-06d2d8db5d344a1ed"
      vpc                 = "vpc-02669f01859cd3545"
      connectivity_type   = "transit-gateway"
      connector_placement = "none"
      subnets             = ["10.24.3.0/28","10.24.3.32/28"]
    }
  }
  policies         = ["ALLOW-ALL-NETWORKS"]
  onboard_app      = true
  decommission_app = false
}

#AWS with transit gateway and connector placement in Infra VPC.
resource "prosimo_network_onboarding" "sin-subnet-1" {
  name = "sin-subnet-tf"
  public_cloud {
    cloud_type        = "public"
    connection_option = "private"
    cloud_creds_name  = "prosimo-aws-app-iam"
    region_name       = "ap-southeast-1"
    cloud_networks {
      vpc               = "vpc-01556b89470488af8" 
      connectivity_type = "transit-gateway"
      connector_placement = "Infra VPC"
      subnets           = ["192.168.250.0/26"]
      connector_settings {
        bandwidth = "small"
        bandwidth_name = "<1 Gbps"
        instance_type = "t3.medium"
      }
    }

 connect_type = "connector"
  }
  policies                = ["DENY-ALL-NETWORKS"]
  onboard_app             = false
  decommission_app        = false
}


#PrivateDC Network Onboarding
resource "prosimo_network_onboarding" "privateDC" {
  name = "private-network-test"
  private_cloud {
    cloud_creds_name  = "PrivateDC"
     subnets           = ["10.0.0.2/32"]
  }
  policies                = ["ALLOW-ALL-NETWORKS"]
  onboard_app             = false
  decommission_app        = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `decommission_app` (Boolean) Set this to true if you would like the network  to be offboarded from fabric
- `name` (String) name for the application
- `onboard_app` (Boolean) Set this to true if you would like the network to be onboarded to fabric

### Optional

- `policies` (List of String) Select policy name.e.g: ALLOW-ALL-NETWORKS, DENY-ALL-NETWORKS or Custom Policies
- `private_cloud` (Block Set) (see [below for nested schema](#nestedblock--private_cloud))
- `public_cloud` (Block Set) (see [below for nested schema](#nestedblock--public_cloud))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `wait_for_rollout` (Boolean) Wait for the rollout of the task to complete. Defaults to true.

### Read-Only

- `deployed` (Boolean)
- `id` (String) The ID of this resource.
- `pam_cname` (String)
- `status` (String)

<a id="nestedblock--private_cloud"></a>
### Nested Schema for `private_cloud`

Required:

- `cloud_creds_name` (String) cloud application account name.

Optional:

- `subnets` (List of String) subnet cider list


<a id="nestedblock--public_cloud"></a>
### Nested Schema for `public_cloud`

Required:

- `cloud_creds_name` (String) cloud application account name.
- `cloud_networks` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--public_cloud--cloud_networks))
- `connection_option` (String) public or private cloud
- `region_name` (String) Name of cloud region

Optional:

- `cloud_type` (String) public or private cloud
- `connect_type` (String) connector

Read-Only:

- `id` (String) The ID of this resource.

<a id="nestedblock--public_cloud--cloud_networks"></a>
### Nested Schema for `public_cloud.cloud_networks`

Required:

- `connector_placement` (String) Infra VPC, Workload VPC or none.

Optional:

- `connectivity_type` (String) transit-gateway, vpc-peering
- `connector_settings` (Block Set) (see [below for nested schema](#nestedblock--public_cloud--cloud_networks--connector_settings))
- `hub_id` (String) (Required if transit-gateway is selected) tgw-id
- `subnets` (List of String) subnet cider list
- `vnet` (String) VNET ID
- `vpc` (String) VPC ID

<a id="nestedblock--public_cloud--cloud_networks--connector_settings"></a>
### Nested Schema for `public_cloud.cloud_networks.connector_settings`

Required:

- `bandwidth` (String) EX: small, medium, large
- `bandwidth_name` (String) EX: <1 Gbps, >1 Gbps
- `instance_type` (String) EX: t3.medium, t3.large




<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)

