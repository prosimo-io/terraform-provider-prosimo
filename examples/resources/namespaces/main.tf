resource "prosimo_namespace" "test" {
    name = "test-ns"
    assign {
        source_networks = ["test"]
    }
}
