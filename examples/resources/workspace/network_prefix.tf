resource "prosimo_network_prefix" "new" {
    cloud_account = "prosimo-aws-app-iam"
    cloud_region = "us-east-1"
    cloud_network = "vpc-544cde2e"
    prefix_route_tables {
        ip_prefix = "10.10.0.0/24"
        route_tables {
            route_table = "rtb-a9c5dcd6"
        }
    }

    prefix_route_tables {
        ip_prefix = "10.12.0.0/24"
        route_tables {
            route_table = "rtb-a9c5dcd6"
        }
    }
    enable = true
}
