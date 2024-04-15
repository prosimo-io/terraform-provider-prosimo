data "prosimo_discovered_networks" "discovered" {
  filter="name==asia-northeast3"
}

output "discovered_network_details" {
  description = "Discovered Network Details"
  value       = data.prosimo_discovered_networks.discovered 
  
}