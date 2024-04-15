data "prosimo_ip_address" "ip_details" {
  filter="cloudtype==AZURE"
  }

  output "ip_details" {
  description = "ip"
  value       = data.prosimo_ip_address.ip_details
}