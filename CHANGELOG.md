## 4.5.3(Sep 27 2024)
### Feature:
- Terraform support for gcp shared service.
## 4.4.3(Septemebr 10 2024)
### Bugfix:
- PRO-22649: Terraform - prosimo_network_onboarding resource - Azure network onboarding allows random subnets to be onboarded as a part of the network
## 4.4.2(Septemebr 10 2024)
### Bugfix:
- PRO-21161: [Terraform] Multiple endpoints attached for a single domain
## 4.4.1(Jun 28 2024)
### Feature:
- Terraform support for network Rules with Internet Access
## 4.3.1(May 16 2024)
### Feature:
- Terraform support for app onboarding 2.0
## 4.2.1(April 18   2024)
### :
- Feature: Terraform support for Managed firewall and Internet egress.
    ```hcl   
    resource "prosimo_firewall_manager" "fm" {
        integration_type = "panorama"
        ip_address = "52.68.81.235"
        api_key = "********"
        license_settings {
        license_mode = "BYOL"
        firewall_family = "VM-100"
        instance_family= "computeOptimized"
        }
    }

    resource "prosimo_managed_firewall" "mf" {
        name = "tf"
        firewall_type = "vmseries"
        cloud_creds_name = "prosimo-aws-app-iam"
        cloud_region = "us-west-2"
        cidr = "10.220.0.0/16"
        instance_size = "c4.xlarge"
        version = "11.0.2"
        auth_key = "0TrustNext"
        auth_code = "0TrustNext"
        scaling_settings {
        desired = 1
        min = 1
        max = 2
        }
        assignments {
        template_name = "test-stac-1"
        device_group = "aws-fw-1"
        }
        access_details {
        username = "***"
        password = "****"
        select_option_for_ssh = "new key pair"
        }
        onboard = false
        decommission = true

    }

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
    ```
    ```hcl

## 4.1.1(March 18   2024)
### :
- BugFix https://prosimoio.atlassian.net/browse/PRO-19277.
  resource `prosimo_namespaces` has been devided in to two resources,`prosimo_namespaces`(Includes namespace creation and assignment) and `prosimo_namespaces_export`(Includes Export and Withdraws). 
```hcl   
resource "prosimo_namespace" "test" {
    name = "test-ns"
    assign {
        source_networks = ["test"]
    }
}
```
```hcl    
resource "prosimo_namespace_export" "test" {
    name = "test-ns"
    export {
        source_network = "test"
        namespaces = [ "default" ]
    }
}
```
- BugFix https://prosimoio.atlassian.net/browse/PRO-19360.
- BugFix https://prosimoio.atlassian.net/browse/PRO-18733
- BugFix https://prosimoio.atlassian.net/browse/PRO-18819 Introduced new tgw `id` filed in transit deployment.
```hcl    
    transit_deployment {
        tgws {
            id= "tgw-ob9531fe31541c"
            action = "MOD"
            connection {
                type = "EDGE"
                action = "ADD"
            }
        }
```
## 4.1.0(January 10   2024)
### Features:
- Field `connector_settings` in resource `prosimo_network_onboarding` and field `node_size_settings` in resource `prosimo_edge` for cloud AZURE have been modified.Here are the changes.
```hcl   // For prosimo_edge
 node_size_settings {
    bandwidth_range {
        min = 7
        max = 9
    }
    }
```
```hcl     // For prosimo_network_onboarding
    connector_settings {
    bandwidth_range {
        min = 3
        max = 5
    }
    }
```
## 3.10.9(December 24   2023)
### Features:
- Field `subnets` in resource `prosimo_network_onboarding` have been changed, the new value is a list of block with inclusion of `virtual_subnet` optional field.
```hcl
    subnets {
    subnet = "10.4.0.0/24"
    virtual_subnet = "10.4.0.0/24"
    }
```
- Field `connector_settings` in resource `prosimo_network_onboarding` and field `node_size_settings` in resource `prosimo_edge` for cloud AWS have been modified.Here are the changes.
```hcl   // For prosimo_edge
 node_size_settings {
    bandwidth_range {
        min = 7
        max = 9
    }
    }
```
```hcl     // For prosimo_network_onboarding
    connector_settings {
    bandwidth_range {
        min = 3
        max = 5
    }
    }
```

## 3.9.9(November 14   2023)
### BugFix:
- Azure Network Onboarding fails with error "Invalid Connectivity Options
## 3.9.8(October 30   2023)
### Features:
- Terraform support for regional prefix. Created resource named `prosimo_regional_prefix`. 
- With latest ui changes, resource `"prosimo_namespace"` does not have withdraw field anymore, the withdrawal would be taken care with the help of export api.
- Terraform support for force decom of apps,networks and edges. With the new feature, optional field `force_offboard` has been    introduced  which enables force decom option. The value defaults to "true" and can be overwritten. A of now force decom workflow would only be triggered if normal offboarding has failed.
Refer docs for more info and sample HQLs.
## 3.8.8(October 06   2023)
### Enhancements:
With latest changes in Edge Api resource `prosimo_ip_address` has been Obsolated. 
### BugFix:
- for resource `prosimo_network_onboarding`, `connector_settings` field dependancy has been removed when connector placement type is "NONE".
## 3.8.7(Septemeber 30   2023)
### Features:
- Terraform support for Latest API changes for edge bringup.
    Introduced a new field names `ip_range`, `deploy_edge`, `decommission_edge` similar to app onboarding.
## 3.7.7(Septemeber 06 2023)
### Features:
- Terraform support for Visual Transit.
Resource name `prosimo_visual_transit`. Refer documentation for more details.
- Resource `prosimo_network_onboarding` enhancements: As part this change connector setting fields have been updated as per latest UI. 
    - filed `bandwidth_name` has been removed, now used need to share the required bandwidth and instance type fields.
    - Doc improvements around connector setting options. 
## 3.6.7(August 23 2023)
### Features:
- Terraform support for IP-based Service Core: As part of this feature resource `prosimo_app_onboarding_ip` has been renamed to  `prosimo_app_onboarding_dns`. In addition there are couple fields being added.
`service_ip_type`: Select if the target needs to be assigned a specific IP address or it could be auto-generated. Even if manually assigned, the address needs to be from the service core IP pool. Default method is to auto generate an IP address from the service core pool.
`service_ip`: ` Service Ip Address.
  
## 3.5.7(July 24 2023)
### Features:
- Terraform support for AZURE cloud in `prosimo_shared_service` and `prosimo_service_insertion` resources.Ref docs for examples.
### BugFix:
- Allow Connector Setting inputs for AZURE and GCP Cloud in resource `prosimo_network_onboarding`.
## 3.4.6(June 30 2023)
### BugFix:
-  Better handling of api error response for resource  `prosimo_network_onboarding`.The terraform execution flow would now stop once it receives error from api.
        
## 3.4.5(June 29 2023)
### BugFix:
-  Allow multiple protocol configuration using terraform for following resources.(Please refer docs for related examples)
         `prosimo_app_onboarding_fqdn`
         `prosimo_app_onboarding_ip`
## 3.4.4(June 24 2023)
### Features:
- Terraform support for network namespaces and cloud gateway .
    New Resources:
        `prosimo_namespace`, 
         `prosimo_cloud_gateway`
- New field `namespace` included in following resources.
        `prosimo_network_onboarding`, 
         `prosimo_policy`
- Additionally following fields have been included in resource `prosimo_network_onboarding`: 
        `connector_subnets` (List of String) connector subnet cider list
         and `service_insertion_endpoint_subnets` (String) Service Insertion Endpoint, applicable when connector is placed in Workload VPC
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

