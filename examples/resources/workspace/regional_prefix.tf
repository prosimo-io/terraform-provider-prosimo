resource "prosimo_regional_prefix" "new" {
    cidr = ["10.10.0.0/24"]
    all_regions = true
}

resource "prosimo_regional_prefix" "new_test" {
    cidr = ["10.10.0.0/24"]
    all_regions = false
    selected_regions {
      cloud_type = "AWS"
      cloud_region = [ "us-east-2" ]
    }
}