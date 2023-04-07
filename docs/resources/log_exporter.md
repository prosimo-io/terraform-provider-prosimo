---
page_title: "prosimo_log_exporter Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this resource to create/modify log exporter settings.
---

# prosimo_log_exporter (Resource)

Use this resource to create/modify log exporter settings.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
resource "prosimo_log_exporter" "splunk" {
    name = "demo3"
    ip = "10.10.10.10"
    tcp_port = 8085
    tls_enabled = false
    description = "splunk"
    auth_token = var.auth_token
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auth_token` (String, Sensitive) Authentication token from receiver endpoint
- `description` (String) Description about log receiver
- `ip` (String) IP address of log receiver endpoint
- `name` (String) Name of log receiver endpoint
- `tcp_port` (Number) port of log receiver endpoint
- `tls_enabled` (Boolean) Defaults to false, set it true to enable tls verification

### Read-Only

- `created_time` (String)
- `id` (String) The ID of this resource.
- `status` (String)
- `team_id` (String)
- `updated_time` (String)
