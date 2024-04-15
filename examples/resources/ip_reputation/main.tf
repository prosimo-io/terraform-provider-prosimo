
resource "prosimo_ip_reputation" "ip_reputation" {
  enabled        = true
  allowlist = ["1.1.1.3/16", "1.1.1.2/16", "1.1.1.4/16"]
}

