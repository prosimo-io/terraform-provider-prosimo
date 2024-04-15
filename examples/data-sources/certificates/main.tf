data "prosimo_certificates" "cert_list"{
  filter = "ca==Pebble Intermediate CA 71d5c3"
}

output "cert_details" {
  description = "certifiates"
  value       = data.prosimo_certificates.cert_list
}
