---
page_title: "prosimo_edr_integration Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this resource to create/modify end point security integrations.
---

# prosimo_edr_integration (Resource)

Use this resource to create/modify end point security integrations.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
resource "prosimo_edr_profile" "crowdstrike" {
    name = "demo3"
    vendor = "CrowdStrike"
    auth {
      client_id =var.clientID
      client_secret = var.clientSecret
      base_url = "https://demo.crowdstrike.com"
      customer_id = "customer_id"
      mssp = false
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `crowdstrike` (Block List, Min: 1) (see [below for nested schema](#nestedblock--crowdstrike))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--crowdstrike"></a>
### Nested Schema for `crowdstrike`

Required:

- `criteria` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--crowdstrike--criteria))
- `name` (String) Endpoint Security Integration name
- `vendor` (String) Select EDR Vendor, for now only CrowdStrike is supported.

Read-Only:

- `id` (String) The ID of this resource.

<a id="nestedblock--crowdstrike--criteria"></a>
### Nested Schema for `crowdstrike.criteria`

Required:

- `sensor_active` (String) Activate sensor, e.g: enabled, disabled
- `status` (String) Status, e.g: enabled, disabled
- `zta_score` (Block Set, Min: 1) Zero Trust Access Score (see [below for nested schema](#nestedblock--crowdstrike--criteria--zta_score))

<a id="nestedblock--crowdstrike--criteria--zta_score"></a>
### Nested Schema for `crowdstrike.criteria.zta_score`

Required:

- `from` (Number)
- `to` (Number)
