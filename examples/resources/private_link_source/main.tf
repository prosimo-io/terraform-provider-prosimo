resource "prosimo_private_link_source" "tf1" {
    name = "terraform-test"
    cloud_creds_name = "prosimo-aws-app-iam"
    cloud_region = "us-east-2"
    cloud_sources {
        cloud_network {
            name = "cloud-census-us-east-2-vpc"
        }
        subnets {
            cidr = "192.13.0.0/24"
        }
    }
}