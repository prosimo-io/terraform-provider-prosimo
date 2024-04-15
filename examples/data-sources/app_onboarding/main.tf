data "prosimo_app_onboarding" "test-app" {
# filter_cloud_type = ["AZURE", "GCP"]
# filter_cloud_region = ["us-east-2"]
# filter_protocaol = "http"
# filter_port = 90
# filter_internal_domain = "psonar-url-rewrite-wordpress.myeventarena.com"
# filter = "id==d7d4b477-b01c-45c4-b786-eca5a0f674ae&app_name==common-app"
# filter = "id==8209c912-977e-48fe-9007-04e48b49aca6"
# filter_app_fqdn = "sshtesting.access.myeveqingchen1657127431278.scnetworkers.info"
# filter = "app_access_type==agent&status==CONFIGURING"
# filter = "name==us-west-1"
# filter = "app_name!=us-east10"
# filter = "apponboardtype!=behind_fabric"
filter = "status!=CONFIGURED"
# filter = "status==CONFIGURED"
}

output "app_onboard_output" {
  description = "applist"
  value       = data.prosimo_app_onboarding.test-app
}
