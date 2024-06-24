---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prosimo_network_prefix Data Source - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this data source to get information on summary route network prefixes.
---

# prosimo_network_prefix (Data Source)

Use this data source to get information on summary route network prefixes.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) Custom filters to scope specific results. Usage: filter = cloudregion==us-west-2

### Read-Only

- `id` (String) The ID of this resource.
- `network_prefix_count` (Number) Total Number of configured/deployed summary route network prefixes
- `network_prefix_list` (Set of Object) (see [below for nested schema](#nestedatt--network_prefix_list))

<a id="nestedatt--network_prefix_list"></a>
### Nested Schema for `network_prefix_list`

Read-Only:

- `cloudkeyid` (String)
- `cloudnetworkid` (String)
- `cloudnetworkname` (String)
- `cloudregion` (String)
- `createdtime` (String)
- `csp` (String)
- `enabled` (Boolean)
- `id` (String)
- `overwriteroute` (Boolean)
- `prefixroutetables` (Set of Object) (see [below for nested schema](#nestedobjatt--network_prefix_list--prefixroutetables))
- `status` (String)
- `teamid` (String)
- `updatedtime` (String)

<a id="nestedobjatt--network_prefix_list--prefixroutetables"></a>
### Nested Schema for `network_prefix_list.prefixroutetables`

Read-Only:

- `prefix` (String)
- `routetables` (Set of Object) (see [below for nested schema](#nestedobjatt--network_prefix_list--prefixroutetables--routetables))

<a id="nestedobjatt--network_prefix_list--prefixroutetables--routetables"></a>
### Nested Schema for `network_prefix_list.prefixroutetables.routetables`

Read-Only:

- `id` (String)
- `name` (String)

