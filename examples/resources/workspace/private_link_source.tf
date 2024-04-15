# resource "prosimo_private_link_source" "tf1" {
#     name = "terraform-test"
#     cloud_creds_name = "prosimo-app"
#     cloud_region = "centralus"
#     cloud_sources {
#         cloud_network {
#             name = "nw-azure-eastus-2-rg/azure-vnet-centralus"
#         }
#         subnets {
#             cidr = "10.15.0.0/24"
#         }
#     }
#     # cloud_sources {
#     #     cloud_network {
#     #         name = "punit-vpc"
#     #     }
#     #     subnets {
#     #         cidr = "10.30.1.0/24"
#     #     }
#     # }
# }