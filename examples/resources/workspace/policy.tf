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
    #  namespace = "default"
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
         
    #    apps {
    #         selecteditems {
    #             name = "agent-httpbin"
    #         }
    #     }
        # networks {
        #     selecteditems {
        #         name = "gcp-usw1-vpc1-tf"
        #     }
        # }
     }
 }

# output "policy_details" {
#   description = "policy details"
#   value       = data.prosimo_policy_access
# }