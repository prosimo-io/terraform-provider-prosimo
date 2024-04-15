# Either of these filter options can be used.Max one filter condition per request. 
data "prosimo_network_onboarding" "test-app" {
# filter_cloud_type = [""]
# filter_cloud_region = ["us-east-2"]
# filter = "pamcname==a-o0pt6l.network.prosimo-eng.prosimoedge.us,id==8524d150-9a13-4a2a-97ee-9ae6f1203315"
# filter = "name==private-DC,status==DEPLOYED"
filter = "status==DEPLOYED"
# filter = "id==9d5863fe-5c6c-4765-a54d-3f827ffbbfd7,connectionoption==private"
}

output "network_onboard_output" {
  description = "applist"
  value       = data.prosimo_network_onboarding.test-app.onboarded_networks
}