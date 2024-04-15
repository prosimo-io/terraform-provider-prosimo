#Previous Setup
resource "prosimo_namespace" "test" {
    name = "test-ns"
    assign {
        source_networks = ["test"]
    }
    # export {
    #     source_network = "demo_network_new"
    #     namespaces = [ "namespace1" ]
    # }
}
# #New Setup
# resource "prosimo_namespace" "test-new" {
#     name = "test1"
#     assign {
#         source_networks = ["demo_network_new"]
#     }
#     export {
#         source_network = "demo_network_new"
#         namespaces = [ "namespace1" ]
#     }
# }