---
page_title: "prosimo_private_link_source Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this resource to create/modify Private Link Sources.
---

# prosimo_private_link_source (Resource)

Use this resource to create/modify Private Link Sources.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
resource "prosimo_private_link_source" "tf1" {
    name = "terraform-test"
    cloud_creds_name = "prosimo-aws-app-iam"
    cloud_region = "us-east-2"
    cloud_sources {
        cloud_network {
            name = "cloud-census-us-east-2-vpc"
        }
        subnets {
            cidr = "192.13.0.0/24"
        }
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_creds_name` (String) cloud account under which application is hosted
- `cloud_region` (String) EX: us-west-2, eu-east-1
- `name` (String) Name of Private Link Source

### Optional

- `cloud_sources` (Block List) (see [below for nested schema](#nestedblock--cloud_sources))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) Resource ID

<a id="nestedblock--cloud_sources"></a>
### Nested Schema for `cloud_sources`

Required:

- `subnets` (Block List, Min: 1) (see [below for nested schema](#nestedblock--cloud_sources--subnets))

Optional:

- `cloud_network` (Block Set) (see [below for nested schema](#nestedblock--cloud_sources--cloud_network))

Read-Only:

- `id` (String) The ID of this resource.

<a id="nestedblock--cloud_sources--subnets"></a>
### Nested Schema for `cloud_sources.subnets`

Required:

- `cidr` (String) Subnet Details


<a id="nestedblock--cloud_sources--cloud_network"></a>
### Nested Schema for `cloud_sources.cloud_network`

Required:

- `name` (String) Name of source VPC.



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
