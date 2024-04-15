#AWS
resource "prosimo_edge" "edge" {
  cloud_name        = "prosimo-aws-iam"
  cloud_region = "us-west-1"
  ip_range = "100.190.0.0/23"
  node_size_settings {
    bandwidth_range {
        min = 7
        max = 9
    }
    }
  deploy_edge = true
  decommission_edge = false
}

#AZURE
resource "prosimo_edge" "edge_azure" {
  cloud_name        = "prosimo-infra"
  cloud_region = "westus"
  ip_range = "100.80.0.0/23"
    node_size_settings {
    bandwidth_range {
        min = 7
        max = 9
    }
    }
  deploy_edge = true
  decommission_edge = false
}

#GCP
resource "prosimo_edge" "edge_gcp" {
  cloud_name        = "prosimo-gcp-infra"
  cloud_region = "westus"
  ip_range = "100.80.0.0/23"
  deploy_edge = true
  decommission_edge = false
}