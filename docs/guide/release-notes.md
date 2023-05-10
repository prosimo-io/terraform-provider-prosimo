## 3.3.4(April 27 2023)
### Enhancements:
- `Terraform destroy` now would both offboard and delete an app/network if it's in onboarded state.
- Data Source filter enhancements.


## 3.3.3(April 14 2023)

### Features:
- Terraform support for Firewall Service Insertion.
    New Resources:
        `prosimo_shared_services`
        `prosimo_service_insertion`
- Terraform support for Private Link Mapping.
    New Resources:
        `prosimo_private_link_source`
        `prosimo_private_link_mapping`
### Notes:
- Refer documentation for examples of each.
## 3.2.3(April 04 2023)

### Features:
- Terraform support for onboarding apps hosted in PrivateDC.
- Terraform support for onboarding Networks in PrivateDC.
### Notes:
# resource_app_onboarding_*:
- Introduced new fields `app_hosted_type` and  `dc_app_ip` .
    - `app_hosted_type` (String) Wheather app is hosted in Public cloud like AWS/AZURE/GCP or private DC. Available options PRIVATE/PUBLIC.
    - `dc_app_ip` (String) Applicable only if  app_hosted_type is PRIVATE, IP of the app hosted in PRIVATE DC.
- With PrivateDC support, field `edge_regions` is optional now.
- Field `idp_name` is Optional now instead of required.
# resource_network_onboarding:
- Introduced new field `private_cloud`:
```hcl
    "private_cloud": {
        Type:     schema.TypeSet,
        Optional: true,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "cloud_creds_name": {
                    Type:        schema.TypeString,
                    Required:    true,
                    Description: "cloud application account name.",
                },
                "subnets": {
                    Type:        schema.TypeList,
                    Optional:    true,
                    Elem:        &schema.Schema{Type: schema.TypeString},
                    Description: "subnet cider list",
                },
            },
        },
    },
```
- Filed `cloud_type` is Optional now with default value being `public`.
Refer docs for examples of each.



## 3.1.3(March 29 2023)
### Notes:
- Supported Terraform version: v1.x

### Features:
- API endpoint update for policy read api. Old API: "/api/policy" with GET call, New API: "/api/policy/search" with POST call"
- In resource prosimo_network_onboarding connector setting options are only applicable for aws cloud now. For azure the default size is incorporated through terraform.
```hcl
    connector_settings {
        bandwidth = "small"
        bandwidth_name = "<1 Gbps"
        instance_type = "Standard_A2_v2"
    }
```
## 3.1.2(March 03 2023)
### Notes:
- Supported Terraform version: v1.x

### Features:
- Implemented support for datasource filters.Ref documentation for more details.

