# Either of these filter options can be used.Max one filter condition per request. 
data "prosimo_private_links" "test-pvt-link" {
# filter_cloud_type = [""]
# filter_cloud_region = ["us-east-2"]
# filter = "pamcname==a-o0pt6l.network.prosimo-eng.prosimoedge.us,id==8524d150-9a13-4a2a-97ee-9ae6f1203315"
# filter = "name==private-DC,status==DEPLOYED"
filter = "region!=us-east-1"
# filter = "id==9d5863fe-5c6c-4765-a54d-3f827ffbbfd7,connectionoption==private"
}

output "private_link_output" {
  description = "pvtlinklist"
  value       = data.prosimo_private_links.test-pvt-link.private_links
}