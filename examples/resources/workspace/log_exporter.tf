resource "prosimo_log_exporter" "splunk" {
    name = "demo3"
    ip = "10.10.10.10"
    tcp_port = 8085
    tls_enabled = false
    description = "splunk"
    auth_token = "var.auth_token"
}