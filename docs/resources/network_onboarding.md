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
    network_exportable_policy = false
    namespace = "default"
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
          subnets {
            subnet = "10.4.0.0/24"
          }
        }
        cloud_networks {
          vnet = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Gitlab/providers/Microsoft.Network/virtualNetworks/Gitlab-vnet"
          hub_id = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"
          connectivity_type = "vwan-hub"
          connector_placement = "Workload VPC"
          subnets {
            subnet = "10.3.5.0/24"
            virtual_subnet = "10.250.2.128/25"
          }
          connector_settings {
            connector_subnets= ["10.99.0.0/24"]
            bandwidth_range {
                min = 3
                max = 5
            }
          }
          
        }
        connect_type = "connector"

    }
    onboard_app = false
    decommission_app = false
}


# GCP Workload Connector Placement
resource "prosimo_network_onboarding" "testapp-s3" {

    name = "demo_network_gcp"
    network_exportable_policy = false
    namespace = "default"
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-gcp-app"
        region_name = "us-west2"
        cloud_networks {
          vpc = "https://www.googleapis.com/compute/v1/projects/prosimo-test-infra/global/networks/default"
          connector_placement = "Workload VPC"
          connectivity_type = "vpc-peering"
          subnets {
            subnet = "10.4.0.0/24"
          }
          connector_settings {
            connector_subnets= ["10.168.0.0/20"]
            bandwidth = "<1 Gbps"
            instance_type = "e2-standard-2"
          }
        }
        connect_type = "connector"

    }
    policies = ["ALLOW-ALL-NETWORKS"]
    onboard_app = false
    decommission_app = false
}

#Azure Workload Connector Placement
resource "prosimo_network_onboarding" "testapp-azure-workload-vpc" {

    name = "demo_network_new"
    network_exportable_policy = false
    namespace = "default"
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-infra"
        region_name = "westus"
        cloud_networks {
          vnet = "/subscriptions/77102da4-2e1f-4445-b74a-93e842dc8c3c/resourceGroups/DefaultResourceGroup-WUS/providers/Microsoft.Network/virtualNetworks/DefaultResourceGroupWUSvnet574"
          connectivity_type = "vnet-peering"
          connector_placement = "Workload VPC"
          subnets {
            subnet = "10.4.0.0/24"
            virtual_subnet = "10.250.2.128/25"
          }
          connector_settings {
            connector_subnets = ["10.4.0.0/24"]
            bandwidth_range {
                min = 3
                max = 5
            }
          }
        }
        connect_type = "connector"

    }
    policies = ["ALLOW-ALL-NETWORKS"]
    onboard_app = false
    decommission_app = false
}

#Azure Infra Connector Placement
resource "prosimo_network_onboarding" "testapp-azure-infra-vpc" {

    name = "demo_network_new"
    network_exportable_policy = false
    namespace = "default"
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-infra"
        region_name = "westus"
        cloud_networks {
          vnet = "/subscriptions/77102da4-2e1f-4445-b74a-93e842dc8c3c/resourceGroups/DefaultResourceGroup-WUS/providers/Microsoft.Network/virtualNetworks/DefaultResourceGroupWUSvnet574"
          connectivity_type = "vnet-peering"
          connector_placement = "Infra VPC"
          subnets {
            subnet = "10.4.0.0/24"
            virtual_subnet = "10.168.0.0/20"
          }
          connector_settings {
            bandwidth_range {
                min = 3
                max = 5
            }
          }
          
        }
        connect_type = "connector"

    }
    policies = ["ALLOW-ALL-NETWORKS"]
    onboard_app = false
    decommission_app = false
}

#AWS with transit gateway and workload vpc
resource "prosimo_network_onboarding" "testapp-AWS-WorkLoad-vpc" {

    name = "demo_network_aws"
    namespace = "default"
    network_exportable_policy = false
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-aws-iam"
        region_name = "us-east-2"
        cloud_networks {
          vpc = "vpc-a8892dc3"
          hub_id = "tgw-04d69a6cd846cd26b"
          connector_placement = "Workload VPC"
          connectivity_type = "transit-gateway"
          service_insertion_endpoint_subnets = "auto"
          subnets {
            subnet = "10.250.2.128/25"
            virtual_subnet = "10.168.0.0/20"
          }
          connector_settings {
            connector_subnets = ["10.4.0.0/24"]
            bandwidth_range {
                min = 3
                max = 5
            }
          }
        }

        connect_type = "connector"

    }
    policies = ["ALLOW-ALL-NETWORKS"]
    onboard_app = false
    decommission_app = false
}

#AWS with transit gateway and infra vpc
resource "prosimo_network_onboarding" "testapp-AWS-Infra-vpc" {

    name = "demo_network_aws"
    namespace = "default"
    network_exportable_policy = false
    public_cloud {
        cloud_type = "public"
        connection_option = "private"
        cloud_creds_name = "prosimo-aws-iam"
        region_name = "us-east-2"
        cloud_networks {
          vpc = "vpc-a8892dc3"
          hub_id = "tgw-04d69a6cd846cd26b"
          connector_placement = "Infra VPC"
          connectivity_type = "transit-gateway"
          subnets {
            subnet = "10.250.2.128/25"
            virtual_subnet = "10.168.0.0/20"
          }
          subnets {
            subnet = "10.250.3.128/25"
            virtual_subnet = "10.168.1.0/20"
          }
          connector_settings {
            bandwidth_range {
                min = 3
                max = 5
            }
          }
        }

        connect_type = "connector"

    }
    policies = ["ALLOW-ALL-NETWORKS"]
    onboard_app = false
    decommission_app = false
}


#PrivateDC Network Onboarding
resource "prosimo_network_onboarding" "privateDC" {
  network_exportable_policy = false
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
- `network_exportable_policy` (Boolean) Mark Network Exportable in Policy
- `onboard_app` (Boolean) Set this to true if you would like the network to be onboarded to fabric

### Optional

- `force_offboard` (Boolean) Force app offboarding incase of normal offboarding failure.
- `namespace` (String) Assigned Namespace
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

- `connectivity_type` (String) transit-gateway, vpc-peering & public(Only applicable if connector placement is in WorkLoad VPC)
- `connector_settings` (Block Set) (see [below for nested schema](#nestedblock--public_cloud--cloud_networks--connector_settings))
- `hub_id` (String) (Required if transit-gateway is selected) tgw-id
- `service_insertion_endpoint_subnets` (String) Service Insertion Endpoint, applicable when connector is placed in Workload VPC
- `subnets` (Block List) subnet cider list (see [below for nested schema](#nestedblock--public_cloud--cloud_networks--subnets))
- `vnet` (String) VNET ID
- `vpc` (String) VPC ID

<a id="nestedblock--public_cloud--cloud_networks--connector_settings"></a>
### Nested Schema for `public_cloud.cloud_networks.connector_settings`

Optional:

- `bandwidth` (String) Available Options: <1 Gbps, 1-5 Gbps, 5-10 Gbps, >10 Gbps
- `bandwidth_range` (Block Set) Applicable for AWS (see [below for nested schema](#nestedblock--public_cloud--cloud_networks--connector_settings--bandwidth_range))
- `connector_subnets` (List of String) connector subnet cider list, Applicable when connector placement is in workload VPC/VNET
- `instance_type` (String) Available Options wrt cloud and bandwidth :Cloud_Provider: AWS:Bandwidth:  <1 Gbps, Available Options: t3.medium/t3a.medium/c5.largeBandwidth:  1-5 Gbps, Available Options: c5a.large/c5.xlarge/c5a.xlarge/c5n.xlargeBandwidth: 5-10 Gbps, Available Options: c5a.8xlarge/c5.9xlargeBandwidth: >10 Gbps, Available Options: c5n.9xlarge/c5a.16xlarge/c5.18xlarge/c5n.18xlargeCloud_Provider: AZURE:For AZURE Default Connector settings are used,hence user does not have to specify is explicitlyProvided values: Bandwidth: <1 Gbps, Instance Type: Standard_A2_v2Cloud_Provider: GCP:Bandwidth:  <1 Gbps, Available Options: e2-standard-2Bandwidth:  1-5 Gbps, Available Options: e2-standard-4Bandwidth: 5-10 Gbps, Available Options: e2-standard-8/e2-standard-16Bandwidth: >10 Gbps, Available Options: c2-standard-16

<a id="nestedblock--public_cloud--cloud_networks--connector_settings--bandwidth_range"></a>
### Nested Schema for `public_cloud.cloud_networks.connector_settings.bandwidth_range`

Required:

- `max` (Number) Minimum Bandwidth Range
- `min` (Number) Minimum Bandwidth Range



<a id="nestedblock--public_cloud--cloud_networks--subnets"></a>
### Nested Schema for `public_cloud.cloud_networks.subnets`

Optional:

- `subnet` (String) Ip Range
- `virtual_subnet` (String) Virtual Subnet




<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)

