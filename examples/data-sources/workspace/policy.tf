# Either of these filter options can be used.Max one filter condition per request. 
# data "prosimo_policy_access" "policy_list" {
#     # filter = "name!=ALLOW-ALL-USERS"
#     # filter = "name=@Bypass"
#     # filter = "name==ALLOW-ALL-NETWORKS&id==cc157573-2db6-43f2-86b3-f6bdd6ff8bf5-network,name==DENY-ALL-NETWORKS"
#     # filter = "name=*allow-all-users"
# }

data "prosimo_policy_access" "policy_list" {
    # filter = "name!=ALLOW-ALL-USERS"
    # filter = "name=@ProAccess"
    # filter = "name==ALLOW-ALL-NETWORKS&id==cc157573-2db6-43f2-86b3-f6bdd6ff8bf5-network,name==DENY-ALL-NETWORKS"
    filter = "name=*allow-all-users"
}

output "policy_details" {
  description = "policy details"
  value       = data.prosimo_policy_access.policy_list
}


data "prosimo_policy_transit" "policy_list_transit" {
#     # filter = "name==ALLOW-ALL-NETWORKS"
#     # filter = "name!=ALLOW-ALL-NETWORKS"
#     # filter = "name==ALLOW-ALL-NETWORKS,id==cc157573-2db6-43f2-86b3-f6bdd6ff8bf5-network"
#     # filter = "name==ALLOW-ALL-NETWORKS&id==cc157573-2db6-43f2-86b3-f6bdd6ff8bf5-network"
#     # filter = "name==ALLOW-ALL-NETWORKS,id==cc157573-2db6-43f2-86b3-f6bdd6ff8bf5-network,name==DENY-ALL-NETWORKS"
      filter = "name=*allow-ALL-NETWORKS"
#     # filter = "name=@new"
}

output "policy_details_transit" {
  description = "policy details_transit"
  value       = data.prosimo_policy_transit.policy_list_transit
}

