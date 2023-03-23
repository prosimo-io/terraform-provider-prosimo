---
page_title: "DataSource Filters"
subcategory: ""
description: |-
  Use DataSource Filters to scope out specific results.
---

# Filter results of datasource.

Filter results of the list type datasources.

## Filter
**Syntax**: `[key][operator][pattern]`

The following filter operators are supported:

* `==` - Pattern must be identical to the value, case-sensitive.
* `=*` - Pattern must be identical to the value, case-insensitive.
* `!=` - Pattern does not match the value, case-insensitive.
* `=@` - Pattern found within value, case-insensitive.




### Custom Filter Examples:

```HCL
# To display Onboarded apps where app access type is agentless, use:
data "prosimo_app_onboarding" "test-app" {
    filter = "app_access_type==agentless"
}

output "app_onboard_output" {
  description = "applist"
  value       = data.prosimo_app_onboarding.test-app
}

# To display all policies except ALLOW-ALL-NETWORKS, use:
data "prosimo_policy_transit" "policy_list_transit" {
    filter = "name!=ALLOW-ALL-NETWORKS"
}

output "policy_details_transit" {
  description = "policy details_transit"
  value       = data.prosimo_policy_transit.policy_list_transit
}

# To match policy name  ALLOW-ALL-NETWORKS case insensitive, use:
data "prosimo_policy_transit" "policy_list_transit" {
    filter = "name=*allow-ALL-NETWORKS"
}

output "policy_details_transit" {
  description = "policy details_transit"
  value       = data.prosimo_policy_transit.policy_list_transit
}

# To match all policy names  which contains "ALL" keyword, use:
data "prosimo_policy_transit" "policy_list_transit" {
    filter = "name=@ALL"
}

output "policy_details_transit" {
  description = "policy details_transit"
  value       = data.prosimo_policy_transit.policy_list_transit
}

```

## Combination

To create a complex query, filters can be combined as follows:

* `Logical OR` - Separate filters using commas `,`.
* `Logical AND` - Separate filters using ampersands `&`.
* `Combining AND and OR` - Separate filters using commas `,` and ampersands `&`.

### Examples:
```HCL
# Logical OR Case:
data "prosimo_policy_transit" "policy_list_transit" {
    filter = "name==ALLOW-ALL-NETWORKS,id==cc157573-2db6-43f2-86b3-f6bdd6ff8bf5-network"
}

output "policy_details_transit" {
  description = "policy details_transit"
  value       = data.prosimo_policy_transit.policy_list_transit
}

# Logical AND Case:
data "prosimo_policy_transit" "policy_list_transit" {
    filter = "name==ALLOW-ALL-NETWORKS&id==cc157573-2db6-43f2-86b3-f6bdd6ff8bf5-network"
}

output "policy_details_transit" {
  description = "policy details_transit"
  value       = data.prosimo_policy_transit.policy_list_transit
}
```

### DataSource specific filters:

These can be used to filter field specific values.The fields generally start with `filter` keyword. e.g>
*   `filter_internal_domain`    = "psonar-url-rewrite-wordpress.myeventarena.com"
*   `filter_app_fqdn`           = "sshtesting.access.myeveqingchen1657127431278.scnetworkers.info"

### Examples:
```HCL
To filter Onboarded app with Internal Domain "psonar-url-rewrite-wordpress.myeventarena.com" use:
data "prosimo_app_onboarding" "test-app" {
    filter_internal_domain = "psonar-url-rewrite-wordpress.myeventarena.com"
}

output "app_onboard_output" {
  description = "applist"
  value       = data.prosimo_app_onboarding.test-app
}

To filter Onboarded app with pApp Fqdn  "a-dok1a1.agentaccess.prosimo-eng.prosimoedge.us" use:
data "prosimo_app_onboarding" "test-app" {
    filter_papp_fqdn = "a-dok1a1.agentaccess.prosimo-eng.prosimoedge.us"
}

output "app_onboard_output" {
  description = "applist"
  value       = data.prosimo_app_onboarding.test-app
}

```