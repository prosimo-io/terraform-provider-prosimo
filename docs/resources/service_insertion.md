---
page_title: "prosimo_service_insertion Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this resource to create/modify service insertion policy.
---

# prosimo_service_insertion (Resource)

Use this resource to create/modify service insertion policy.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
# #AWS
# resource "prosimo_service_insertion_AWS" "firewall" {
#     name = "Terraform-test"
#     service_name = "terraform-test"
#     namespace = "ns-1" 
#     source {
#         networks{
#             name = "terraform-si"
#         }
#     }
#     target {
#         networks{
#             name = "0.0.0.0/0"
#         }
#     }
#     ip_rules {
#         source_addresses = ["any"]
#         source_ports = ["any"]
#         destination_addresses = ["any"]
#         destination_ports = ["any"]
#         protocols = ["ANY"]
#     }
# }

# #AZURE
# resource "prosimo_service_insertion_AZURE" "firewall" {
#     name = "Terraform-test"
#     service_name = "firewall_svc"
#     namespace = "ns-1" 
#     prosimo_managed_routing = true
#     route_tables = ["/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Azure-Lab-Arjun/providers/Microsoft.Network/routeTables/fw-rtb"]
#     source {
#         networks{
#             name = "src-eastus2"
#         }

#     }
#     target {
#         networks{
#             name = "0.0.0.0/0"
#         }

#     }
#     ip_rules {
#         source_addresses = ["any"]
#         source_ports = ["any"]
#         destination_addresses = ["any"]
#         destination_ports = ["any"]
#         protocols = ["ANY"]
#     }

# }

#AWS
resource "prosimo_service_insertion" "firewall" {
    name = "Terraform-test1"
    service_name = "aws-fw"
    namespace = "default" 
    source {
        networks{
            name = "net-us-west-2-infra"
        }
    }
    target {
        networks{
            name = "Internet"
        }
    }
    ip_rules {
        source_addresses = ["any"]
        source_ports = ["10101"]
        destination_addresses = ["20.0.0.0/24"]
        destination_ports = ["30000"]
        protocols = ["TCP"]
    }
    ip_rules {
        source_addresses = ["any"]
        source_ports = ["10102"]
        destination_addresses = ["20.2.0.0/24"]
        destination_ports = ["40000"]
        protocols = ["TCP"]
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of Service Insertion
- `service_name` (String) Name of the Shared Service

### Optional

- `ip_rules` (Block List) Insertion Policy Rules (see [below for nested schema](#nestedblock--ip_rules))
- `namespace` (String) Policy Namespace, Defaults to default
- `prosimo_managed_routing` (Boolean) TRUE if you would like Prosimo to update Firewal VNET Roue Table
- `route_tables` (List of String) List of Route Table ID
- `source` (Block List) (see [below for nested schema](#nestedblock--source))
- `target` (Block List) (see [below for nested schema](#nestedblock--target))
- `type` (String) Service Insertion Type

### Read-Only

- `id` (String) Resource ID
- `status` (String) Service Insertion Deployment Status

<a id="nestedblock--ip_rules"></a>
### Nested Schema for `ip_rules`

Optional:

- `destination_addresses` (List of String) Target Ip Address
- `destination_ports` (List of String) Destination Port
- `protocols` (List of String) Protocols
- `source_addresses` (List of String) Source Ip Address
- `source_ports` (List of String) Source Port


<a id="nestedblock--source"></a>
### Nested Schema for `source`

Optional:

- `networks` (Block Set) (see [below for nested schema](#nestedblock--source--networks))

<a id="nestedblock--source--networks"></a>
### Nested Schema for `source.networks`

Optional:

- `name` (String) Source Network Name



<a id="nestedblock--target"></a>
### Nested Schema for `target`

Optional:

- `apps` (Block Set) (see [below for nested schema](#nestedblock--target--apps))
- `networks` (Block Set) (see [below for nested schema](#nestedblock--target--networks))

<a id="nestedblock--target--apps"></a>
### Nested Schema for `target.apps`

Optional:

- `name` (String) Target App Name


<a id="nestedblock--target--networks"></a>
### Nested Schema for `target.networks`

Optional:

- `name` (String) Target Network Name

