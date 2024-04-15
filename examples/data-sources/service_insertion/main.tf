# Either of these filter options can be used.Max one filter condition per request. 
data "prosimo_service_insertion" "test-service-insertion" {
# filter = "cloudregion==us-west-2"
filter = "status==DEPLOYED"
# filter = "id==a25cf952-acdc-4e7d-8fd2-5a8e6f406bc0"
}

output "service_insertion_output" {
  description = "service_insertion_list"
  value       = data.prosimo_service_insertion.test-service-insertion.service_insertions
}