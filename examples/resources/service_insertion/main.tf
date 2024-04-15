#AWS
resource "prosimo_service_insertion_AWS" "firewall" {
    name = "Terraform-test"
    service_name = "terraform-test"
    namespace = "ns-1" 
    source {
        networks{
            name = "terraform-si"
        }
    }
    target {
        networks{
            name = "0.0.0.0/0"
        }
    }
    ip_rules {
        source_addresses = ["any"]
        source_ports = ["any"]
        destination_addresses = ["any"]
        destination_ports = ["any"]
        protocols = ["ANY"]
    }
}

#AZURE
resource "prosimo_service_insertion_AZURE" "firewall" {
    name = "Terraform-test"
    service_name = "firewall_svc"
    namespace = "ns-1" 
    prosimo_managed_routing = true
    route_tables = ["/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Azure-Lab-Arjun/providers/Microsoft.Network/routeTables/fw-rtb"]
    source {
        networks{
            name = "src-eastus2"
        }

    }
    target {
        networks{
            name = "0.0.0.0/0"
        }

    }
    ip_rules {
        source_addresses = ["any"]
        source_ports = ["any"]
        destination_addresses = ["any"]
        destination_ports = ["any"]
        protocols = ["ANY"]
    }

}
