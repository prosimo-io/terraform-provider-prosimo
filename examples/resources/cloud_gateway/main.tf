resource "prosimo_cloud_gateway" "test" {
    region = "us-west-2"
    attach_point = "tgw-05136861d2429935d"
    attachment = "tgw-attach-0c44c52943c2693c6"
    route_table = "tgw-rtb-0c68aa88db7de77b2"
}
