---
page_title: "prosimo_cache_rules Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  Use this resource to create/modify cache rules.
---

# prosimo_cache_rules (Resource)

Use this resource to create/modify cache rules.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
resource "prosimo_cache_rules" "cacherule" {
    name = "demo80"
    default = true
    editable = true
    share_static_content = false
    cache_control_ignored = true
    bypass_cache = true

    path_patterns {
        path = "*"
        bypass_uri = false
        is_default = false
        is_new_path = false
        status = "existing"
         settings {
            type = "Dynamic"
            user_id_ignored = false
            cache_control_ignored = false
            cookie_ignored = false
            query_parameter_ignored = false 
            ttl {
                enabled = true
                time = 24
                time_unit = "Hours"
            }


        }

    }
    path_patterns {
        path = "/abcdefg"
        bypass_uri = false
        is_default = false
        is_new_path = false
        status = "existing"

        settings {
            type = "Dynamic"
            user_id_ignored = false
            cache_control_ignored = false
            cookie_ignored = false
            query_parameter_ignored = false 
            ttl {
                enabled = true
                time = 24
                time_unit = "Hours"
            }


        }


    }
    bypass_info {
		resp_hdrs {
			x_jenkins_session = [""]
		}
    }
    app_domains {
        domain = "speedtest-server-eastus2-1625128238856.myeventarena.com"
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of cache rule
- `path_patterns` (Block List, Min: 1) Path pattern list. (see [below for nested schema](#nestedblock--path_patterns))

### Optional

- `app_domains` (Block List) (see [below for nested schema](#nestedblock--app_domains))
- `bypass_cache` (Boolean) Defaults to false, set it to true if you want to bypass cache.
- `bypass_info` (Block Set) (see [below for nested schema](#nestedblock--bypass_info))
- `cache_control_ignored` (Boolean) Defaults to false, set it to true if you want to skip cache control.
- `default` (Boolean)
- `editable` (Boolean)
- `is_new` (Boolean)
- `share_static_content` (Boolean) Defaults to false, set it to true if you want to share static content.
- `teamid` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `last_updated` (String)

<a id="nestedblock--path_patterns"></a>
### Nested Schema for `path_patterns`

Required:

- `bypass_uri` (Boolean)
- `is_default` (Boolean)
- `is_new_path` (Boolean)
- `path` (String) Path to store cache
- `settings` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--path_patterns--settings))
- `status` (String)

<a id="nestedblock--path_patterns--settings"></a>
### Nested Schema for `path_patterns.settings`

Required:

- `cache_control_ignored` (Boolean)
- `cookie_ignored` (Boolean)
- `query_parameter_ignored` (Boolean)
- `ttl` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--path_patterns--settings--ttl))
- `type` (String)
- `user_id_ignored` (Boolean)

<a id="nestedblock--path_patterns--settings--ttl"></a>
### Nested Schema for `path_patterns.settings.ttl`

Required:

- `enabled` (Boolean)
- `time` (Number)
- `time_unit` (String)




<a id="nestedblock--app_domains"></a>
### Nested Schema for `app_domains`

Optional:

- `domain` (String)

Read-Only:

- `id` (String) The ID of this resource.


<a id="nestedblock--bypass_info"></a>
### Nested Schema for `bypass_info`

Optional:

- `resp_hdrs` (Block Set) (see [below for nested schema](#nestedblock--bypass_info--resp_hdrs))

<a id="nestedblock--bypass_info--resp_hdrs"></a>
### Nested Schema for `bypass_info.resp_hdrs`

Optional:

- `content_type` (List of String)
- `x_jenkins_session` (List of String)

