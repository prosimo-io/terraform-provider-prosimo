# Either of these filter options can be used.Max one filter condition per request. 
data "prosimo_app_onboarding" "test-app" { 
# filter_cloud_type = ["AZURE", "GCP"]
# filter_cloud_region = ["us-east-2"]
# filter_protocol = "http"
# filter_port = 90
# filter_internal_domain = "psonar-url-rewrite-wordpress.myeventarena.com"
# filter_papp_fqdn = "a-dok1a1.agentaccess.prosimo-eng.prosimoedge.us"
# filter_app_fqdn = "sshtesting.access.myeveqingchen1657127431278.scnetworkers.info"
# filter = "app_access_type==agent&status==CONFIGURING"
# filter = "app_access_type==agentless"
# filter = "status==CONFIGURING"

}

output "app_onboard_output" {
  description = "applist"
  value       = data.prosimo_app_onboarding.test-app
}
