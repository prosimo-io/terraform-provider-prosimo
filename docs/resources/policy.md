---
page_title: "prosimo_policy Resource - terraform-provider-prosimo"
subcategory: ""
description: |-
  The policy engine in Prosimo helps control access rules between users, applications, and networking. Use this resource to create/modify policies.
---

# prosimo_policy (Resource)

The policy engine in Prosimo helps control access rules between users, applications, and networking. Use this resource to create/modify policies.

This resource is usually used along with `terraform-provider-prosimo`.



## Example Usage

```terraform
# resource "prosimo_policy" "testapp1" {
#     name = "demo80"
#     app_access_type = "access"
#     device_posture_configured = true
#     details {
#         actions = "deny"
#         lock_users = false
#         alert = false
#         apps {
#             selecteditems {
#                 name = "common-app"
#                 //id   = "1cfac9ee-9358-43ff-ba93-1a14cd094a64"
#             }
#         }
#          matches {
#             match_entries {
#                 property = "User"
#                 operation = "Does NOT contain"
#                 type     = "users"
#                 values {
#                     inputitems {
#                         name = "siba1"                       
#                     }
#                     inputitems {
#                         name = "siba#"                      
#                     }
#                 }
#             }
#             match_entries {
#                 property = "User"
#                 operation = "Is"
#                 type     = "users"
#                 values {
#                     inputitems {
#                         name = "siba"                      
#                     }
#                     inputitems {
#                         name = "siba@"                       
#                     }
#                 }
#             }

#             match_entries {
#                 property = "Device OS"
#                 operation = "Is NOT"
#                 type     = "devices"
#                 values {
#                     selecteditems {
#                         name = "Linux"                      
#                     }
#                     selecteditems {
#                         name = "Windows"                      
#                     }
#                     selecteditems {
#                         name = "iOS"                      
#                     }
#                 }
#             }
#             match_entries {
#                 property = "Device OS"
#                 operation = "Is NOT"
#                 type     = "devices"
#                 values {
#                     selecteditems {
#                         name = "Windows"                      
#                     }
#                     selecteditems {
#                         name = "Linux"                      
#                     }
#                     selecteditems {
#                         name = "iOS"                      
#                     }
#                 }
#             }
#             match_entries {
#                 property = "Device Category"
#                 operation = "Is NOT"
#                 type     = "devices"
#                 values {
#                     selecteditems {
#                         name = "Desktop"                      
#                     }
#                     selecteditems {
#                         name = "Mobile"                      
#                     }
#                 }
#             }

#             match_entries {
#                 property = "Browser Version"
#                 operation = "Is at least"
#                 type     = "devices"
#                 values {
#                     selecteditems {
#                         name = "10"                      
#                     }
#                 }
#              }
#             match_entries {
#                 property = "Country"
#                 operation = "Is"
#                 type     = "location"
#                 values {
#                     selecteditems {
#                         # name = "10.10.10.10"
#                         #id  = "af"
#                         country_name = "Pakistan"    
#                         # state_name = "Tripura"
#                         # city_name = "Abhanga"                    
#                     }
#                 }
#              }

#             match_entries {
#                 property = "City"
#                 operation = "Is"
#                 type     = "location"
#                 values {
#                     selecteditems {
#                         # name = "10.10.10.10"
#                         #id  = "af"
#                         country_name = "India"    
#                         state_name = "Odisha"
#                         city_name = "Cuttack"                    
#                     }
#                 }
#              }

#             match_entries {
#                 property = "IP Prefix/Address"
#                 operation = "In"
#                 type     = "location"
#                 values {
#                     inputitems {
#                         name = "10.10.10.10/16"
#                         #id  = "af"
#                         # country_name = "India"    
#                         # state_name = "Bihar"                    
#                     }
#                 }
#              }

#             match_entries {
#                 property = "Risk Level"
#                 operation = "Is NOT"
#                 type     = "device-posture"
#                 values {
#                     selecteditems {
#                         name = "High"
#                         #id  = "af"                        
#                     }
#                 }
#              }

#             match_entries {
#                 property = "Risk Level"
#                 operation = "Is"
#                 type     = "device-posture"
#                 values {
#                     selecteditems {
#                         name = "Low"
#                         #id  = "af"                        
#                     }
#                 }
#              }
             
#             match_entries {
#                 property = "Time"
#                 operation = "Between"
#                 type     = "time"
#                 values {
#                     inputitems {
#                         name = "10"                    
#                     }
#                 }
#              }
#             match_entries {
#                 property = "URL"
#                 operation = "Does NOT contain"
#                 type     = "url"
#                 values {
#                     inputitems {
#                         name = "123"                      
#                     }
#                 }
#             }
#             match_entries {
#                 property = "FQDN"
#                 operation = "Is"
#                 type     = "fqdn"
#                 values {
#                     inputitems {
#                         name = "1234"                    
#                     }
#                 }
#             }
#             match_entries {
#                 property = "HTTP Method"
#                 operation = "Is"
#                 type     = "advanced"
#                 values {
#                     selecteditems {
#                         name = "GET"                     
#                     }
#                     selecteditems {
#                         name = "POST"                     
#                     }
#                     selecteditems {
#                         name = "HEAD"                     
#                     }
#                     selecteditems {
#                         name = "DELETE"                     
#                     }
#                     selecteditems {
#                         name = "CONNECT"                     
#                     }
#                 }
#              }
#             match_entries {
#                 property = "abc"
#                 operation = "Ends with"
#                 type     = "idp"
#                 values {
#                     inputitems {
#                         name = "adc"                       
#                     }
#                 }
            
#         }
#       }
#     }
# }
#  resource "prosimo_policy" "policy-for-common-app" {
#      name = "psonar-test-policy"
#      app_access_type = "access"
#      details {
#         actions = "allow"
#         lock_users = false
#         alert = true
#         mfa = true
#         # lock_users = false
#           matches {
#              match_entries {
#                  property = "User"
#                  operation = "Does NOT contain"
#                  type     = "users"
#                  values {
#                      inputitems {
#                          name = "siba1"
#                      }
#                      inputitems {
#                          name = "siba#"
#                      }
#                  }
#              }
#          }
#        apps {
#             selecteditems {
#                 name = "app"
#             }
#         }
#         networks {
#             selecteditems {
#                 name = "azure-subnets-tf"
#             }
#             selecteditems {
#                 name = "azure-subnets-tf"
#             }
#         }
#      }
#  }


  resource "prosimo_policy" "policy-for-common-app-new" {
     name = "psonar-test-policy-transit"
     app_access_type = "transit"
     namespace = "default"
     details {
        actions = "allow"
        # lock_users = false
        # alert = true
        # mfa = true
        # lock_users = false
          matches {
            match_entries {
                property = "URL"
                operation = "Does NOT contain"
                type     = "url"
                values {
                    inputitems {
                        name = "1235678"                      
                    }
                }
            }
            match_entries {
                property = "Network"
                operation = "Is"
                type     = "networks"
                values {
                    selecteditems {
                        name = "aws-uswest1-spoke1-infra-tf"                  
                    }
                }
            }
            match_entries {
                type     = "networkacl"
                values {
                    inputitems {
                       ip_details {
                        source_ip = ["any"]
                        target_ip = ["any"]
                        protocol = ["tcp"]
                        source_port = ["any"]
                        target_port = ["any"]
                       }               
                    }
                }
            }
            match_entries {
                property = "Time"
                operation = "Between"
                type     = "time"
                values {
                    inputitems {
                        name = "10"                    
                    }
                }
             }
          
          match_entries {
                property = "FQDN"
                operation = "Is"
                type     = "fqdn"
                values {
                    inputitems {
                        name = "1234"                    
                    }
                }
            }
             match_entries {
                property = "HTTP Method"
                operation = "Is"
                type     = "advanced"
                values {
                    selecteditems {
                        name = "GET"                     
                    }
                    selecteditems {
                        name = "POST"                     
                    }
                    selecteditems {
                        name = "HEAD"                     
                    }
                    selecteditems {
                        name = "DELETE"                     
                    }
                    selecteditems {
                        name = "CONNECT"                     
                    }
                }
             }
          }
        internet_traffic_enabled = true
    #    apps {
    #         selecteditems {
    #             name = "agent-httpbin"
    #         }
    #     }
        networks {
            selecteditems {
                name = "gcp-usw1-vpc1-tf"
            }
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

- `details` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--details))
- `name` (String) name of policy

### Optional

- `app_access_type` (String) app access type, e.g: access, transit
- `device_posture_configured` (Boolean) only applicable for access app access type, set it to true to enable device posture
- `namespace` (String) Policy Namespace, only applicable for transit app_access_type
- `teamid` (String)
- `types` (String) type of policy, e.g: default, managed

### Read-Only

- `createdtime` (String)
- `id` (String) The ID of this resource.
- `updatedtime` (String)

<a id="nestedblock--details"></a>
### Nested Schema for `details`

Required:

- `actions` (String) policy action, e.g: allow, deny
- `matches` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--details--matches))

Optional:

- `alert` (Boolean) set this to true to trigger the alert as per policy config
- `apps` (Block Set) App details to attach to the policy (see [below for nested schema](#nestedblock--details--apps))
- `bypass` (Boolean) set this to true to bypass policy
- `internet_traffic_enabled` (Boolean) set it to true to enable internet access
- `lock_users` (Boolean) set this to true to lock the user defined in policy
- `mfa` (Boolean) set this to true to trigger
- `networks` (Block Set) Network details to attach to the policy (see [below for nested schema](#nestedblock--details--networks))
- `skipwaf` (Boolean) set this to true to skip waf

<a id="nestedblock--details--matches"></a>
### Nested Schema for `details.matches`

Optional:

- `match_entries` (Block Set) (see [below for nested schema](#nestedblock--details--matches--match_entries))

<a id="nestedblock--details--matches--match_entries"></a>
### Nested Schema for `details.matches.match_entries`

Required:

- `type` (String) Select policy match condition type, for access policy options are users, location, idp, devices, time, url, device-posture, fqdn and advanced. For transit type options are time, url, networkacl, fqdn, networks and advanced

Optional:

- `operation` (String) Operation of the selected property, available options are Id, Is NOT, Contains, Does NOT contain, Starts with, Ends with, In, NOT in, Is at least, Between
- `property` (String) Select property of selected type
- `values` (Block Set) (see [below for nested schema](#nestedblock--details--matches--match_entries--values))

<a id="nestedblock--details--matches--match_entries--values"></a>
### Nested Schema for `details.matches.match_entries.values`

Optional:

- `inputitems` (Block Set) (see [below for nested schema](#nestedblock--details--matches--match_entries--values--inputitems))
- `selectedgroups` (Block Set) (see [below for nested schema](#nestedblock--details--matches--match_entries--values--selectedgroups))
- `selecteditems` (Block Set) (see [below for nested schema](#nestedblock--details--matches--match_entries--values--selecteditems))

<a id="nestedblock--details--matches--match_entries--values--inputitems"></a>
### Nested Schema for `details.matches.match_entries.values.inputitems`

Optional:

- `ip_details` (Block Set) Only applicable for type networkacl (see [below for nested schema](#nestedblock--details--matches--match_entries--values--inputitems--ip_details))
- `name` (String) Input value name

<a id="nestedblock--details--matches--match_entries--values--inputitems--ip_details"></a>
### Nested Schema for `details.matches.match_entries.values.inputitems.ip_details`

Optional:

- `protocol` (List of String) List of protocols
- `source_ip` (List of String) Source IP list
- `source_port` (List of String) Source port list
- `target_ip` (List of String) Target IP list
- `target_port` (List of String) Target port list



<a id="nestedblock--details--matches--match_entries--values--selectedgroups"></a>
### Nested Schema for `details.matches.match_entries.values.selectedgroups`

Optional:

- `name` (String) Input value name


<a id="nestedblock--details--matches--match_entries--values--selecteditems"></a>
### Nested Schema for `details.matches.match_entries.values.selecteditems`

Optional:

- `city_name` (String) City name, only applicable for type location
- `country_name` (String) Country name, only applicable for type location
- `name` (String) Selected value name
- `state_name` (String) State name, only applicable for type location





<a id="nestedblock--details--apps"></a>
### Nested Schema for `details.apps`

Optional:

- `selecteditems` (Block Set) (see [below for nested schema](#nestedblock--details--apps--selecteditems))

<a id="nestedblock--details--apps--selecteditems"></a>
### Nested Schema for `details.apps.selecteditems`

Optional:

- `name` (String) Name of the app



<a id="nestedblock--details--networks"></a>
### Nested Schema for `details.networks`

Optional:

- `selecteditems` (Block Set) (see [below for nested schema](#nestedblock--details--networks--selecteditems))

<a id="nestedblock--details--networks--selecteditems"></a>
### Nested Schema for `details.networks.selecteditems`

Optional:

- `name` (String) Name of the network

