#AWS
resource "prosimo_shared_services" "firewall" {
    name = "firewall_svc"
    region {
        cloud_region = "us-west-2"
        gateway_lb = "com.amazonaws.vpce.us-west-2.vpce-svc-0fdda8395ea088814"
        cloud_creds_name = "prosimo-aws-app-iam"
    }
    onboard = true
    decommission = false
}

#AZURE
resource "prosimo_shared_services" "firewall-Azure" {
    name = "firewall_svc"
    region {
        cloud_region = "eastus2"
        gateway_lb = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/Azure-Lab-Arjun/providers/Microsoft.Network/loadBalancers/vmseries-public-lb"
        cloud_creds_name = "prosimo-app"
        resource_group = "Azure-Lab-Arjun"
    }
    onboard = true
    decommission = false
}