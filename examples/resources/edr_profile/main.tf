
resource "prosimo_edr_profile" "crowdstrike" {
    name = "demo3"
    vendor = "CrowdStrike"
    auth {
      client_id =var.clientID
      client_secret = var.clientSecret
      base_url = "https://demo.crowdstrike.com"
      customer_id = "customer_id"
      mssp = false
    }
}








