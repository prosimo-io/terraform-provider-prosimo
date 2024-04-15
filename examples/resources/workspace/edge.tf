# resource "prosimo_cloud_creds" "awskey" {
#   cloud_type        = "AWS"
#   nickname = "aws-test-key"
# }

# resource "prosimo_ip_addresses" "aws_ip_address" {
#   cidr        = "172.16.0.0/16"
#   cloud_type = prosimo_cloud_creds.awskey.cloud_type
# }

resource "prosimo_edge" "edge" {
  cloud_name        = "prosimo-aws-iam"
  cloud_region = "us-east-2"
  ip_range = "192.168.16.0/23"
  node_size_settings {
    bandwidth_range {
        min = 7
        max = 9
    }
    }
  deploy_edge = true
  decommission_edge = false
  wait_for_rollout = true
}

resource "prosimo_edge" "edge_azure" {
  cloud_name        = "prosimo-infra"
  cloud_region = "westus3"
  ip_range = "192.169.0.0/23"
  node_size_settings {
    bandwidth_range {
        min = 1
        max = 2
    }
    }
  deploy_edge = true
  decommission_edge = false
}

# resource "prosimo_edge" "edge_gcp" {
#   cloud_name        = "prosimo-infra"
#   cloud_region = "westus"
#   ip_range = "100.80.0.0/23"
# #   node_size_settings {
# #     bandwidth = "1-5 Gbps"
# #     instance_type = "c5a.large"
# #     }
#   deploy_edge = true
#   decommission_edge = false
# }