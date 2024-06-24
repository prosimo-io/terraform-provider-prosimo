---
page_title: "prosimo_internet_egress Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  The Internet Egress in Prosimo helps control internet access rules for applications. Use this resource to create/modify egress rules.
---

# prosimo_internet_egress (Resource)

The Internet Egress in Prosimo helps control internet access rules for applications. Use this resource to create/modify egress rules.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
resource "prosimo_internet_egress" "test-internet-egress-policy" {
    name = "psonar-internet-egress-policy"
    action = "allow"
    matches {
        match_entries {
            property = "FQDN"
            operation = "Is"
            type = "fqdn"
            values {
                inputitems {
                    id = "youtube.com"
                }
            }
        }
        match_entries {
            property = "FQDN"
            operation = "Is NOT"
            type = "fqdn"
            values {
                inputitems {
                    id = "prosimo.io"
                }
            }
        }
        match_entries {
            property = "Domain"
            operation = "Is"
            type = "fqdn"
            values {
                inputitems {
                    id = "prosimo.io"
                }
            }
        }
    }
    namespaces {
        namespace_entries {
            name = "default"
        }
    }
    networks {
        network_entries {
            name = "test"
        }
    }
    network_groups {
        network_group_entries {
            name = "test-network-group"
        }
    }
}

# output "policy_details" {
#   description = "policy details"
#   value       = data.prosimo_policy_access
# }
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) Policy action, e.g: allow, deny
- `matches` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--matches))
- `name` (String) Name of Internet Egress Policy

### Optional

- `namespaces` (Block Set) Policy Namespace where the policy can be in the action (see [below for nested schema](#nestedblock--namespaces))
- `network_groups` (Block Set) Network group details to attach to the policy (see [below for nested schema](#nestedblock--network_groups))
- `networks` (Block Set) Network details to attach to the policy (see [below for nested schema](#nestedblock--networks))

### Read-Only

- `createdtime` (String)
- `id` (String) The ID of this resource.
- `updatedtime` (String)

<a id="nestedblock--matches"></a>
### Nested Schema for `matches`

Optional:

- `match_entries` (Block Set) (see [below for nested schema](#nestedblock--matches--match_entries))

<a id="nestedblock--matches--match_entries"></a>
### Nested Schema for `matches.match_entries`

Optional:

- `operation` (String) Operation of the selected property, available options are Is, Is NOT
- `property` (String) Select property of selected type
- `type` (String) Select policy match condition type i.e. - fqdn
- `values` (Block Set) (see [below for nested schema](#nestedblock--matches--match_entries--values))

<a id="nestedblock--matches--match_entries--values"></a>
### Nested Schema for `matches.match_entries.values`

Optional:

- `inputitems` (Block Set) (see [below for nested schema](#nestedblock--matches--match_entries--values--inputitems))

<a id="nestedblock--matches--match_entries--values--inputitems"></a>
### Nested Schema for `matches.match_entries.values.inputitems`

Optional:

- `id` (String) Input domain/fqdn value





<a id="nestedblock--namespaces"></a>
### Nested Schema for `namespaces`

Optional:

- `namespace_entries` (Block Set) (see [below for nested schema](#nestedblock--namespaces--namespace_entries))

<a id="nestedblock--namespaces--namespace_entries"></a>
### Nested Schema for `namespaces.namespace_entries`

Optional:

- `name` (String) Name of the network



<a id="nestedblock--network_groups"></a>
### Nested Schema for `network_groups`

Optional:

- `network_group_entries` (Block Set) (see [below for nested schema](#nestedblock--network_groups--network_group_entries))

<a id="nestedblock--network_groups--network_group_entries"></a>
### Nested Schema for `network_groups.network_group_entries`

Optional:

- `name` (String) Name of the network-group



<a id="nestedblock--networks"></a>
### Nested Schema for `networks`

Optional:

- `network_entries` (Block Set) (see [below for nested schema](#nestedblock--networks--network_entries))

<a id="nestedblock--networks--network_entries"></a>
### Nested Schema for `networks.network_entries`

Optional:

- `name` (String) Name of the network
