resource "prosimo_private_link_mapping" "tfmap" {
    source = "terraform-test"
    target = "azure-common-agent-app"
    hosted_zones {
        vpc_name = "nw-azure-eastus-2-rg/azure-vnet-centralus"
        domain_name = "common-agent-app-azure-1689008453602.myeventarena.com"
        private_hosted_zone = "myeventarena.com."
    }
    #     hosted_zones {
    #     vpc_name = "cloud-census-us-east-2-vpc"
    #     domain_name = "speedtest-server-us-west-2-1681064493897.myeventarena.com"
    #     private_hosted_zone = "myeventarena.com."
    # }
    # hosted_zones {
    #     vpc_name = "punit-vpc"
    #     domain_name = "app-aws-us-west-2-1681064493897.myeventarena.com"
    #     private_hosted_zone = "myeventarena.com."
    # }
    # hosted_zones {
    #     vpc_name = "punit-vpc"
    #     domain_name = "speedtest-server-us-west-2-1681064493897.myeventarena.com"
    #     private_hosted_zone = "myeventarena.com."
    # }
}