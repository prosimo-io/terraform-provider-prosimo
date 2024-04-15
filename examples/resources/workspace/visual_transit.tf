resource "prosimo_visual_transit" "test" {
    transit_input {
    cloud_type   = "AWS"  
    cloud_region = "us-east-2"
    transit_deployment {
        tgws {
            name = "manual-tgw"
            action = "MOD"
            connection {
                type = "EDGE"
                action = "ADD"
            }
        }
        vpcs {
            name = "pktest-1"
            action = "ADD"
        }
    }
  }
deploy_transit_setup = true
}
resource "prosimo_visual_transit" "test1"{
 transit_input {
    cloud_type   = "AZURE"  
    cloud_region = "westus"
    transit_deployment {
        vhubs {
            name = "test-vhub"
            action = "ADD"
            account = "prosimo-app"
            vwan = "/subscriptions/77102da4-2e1f-4445-b74a-93e842dc8c3c/resourceGroups/josh-infra-app/providers/Microsoft.Network/virtualWans/josh-infra-hub"
            address_space = "10.0.0.0/22"
            connection {
                type = "EDGE"
                action = "ADD"
            }

        }
        vnets {
            name = "staging-agent-apps/stagingagentappsvnet146"
            action = "ADD"
        }
    }
 }
 deploy_transit_setup = true
       
}