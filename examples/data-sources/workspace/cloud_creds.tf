# Either of these filter options can be used.Max one filter condition per request. 
data "prosimo_cloud_creds" "cloud_creds" {
  # filter="nickname==prosimo-infra,cloudtype!=AWS"
  # filter="cloudtype=*aws"
  filter="nickname=@prosimo"

}

output "cloud_creds_output" {
  description = "cloud_creds"
  value       = data.prosimo_cloud_creds.cloud_creds
}