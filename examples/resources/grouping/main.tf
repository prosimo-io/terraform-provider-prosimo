resource "prosimo_grouping" "newgrp" {
    name = "demo3"
    type = "APP"
    details {
        apps {
            name = "common-app"
        }
    }
}

resource "prosimo_grouping" "newgrp1" {
    name = "demo3"
    type = "USER"
    details {
        apps {
            name = "common-app"
        }
    }
}

resource "prosimo_grouping" "newgrp9" {
    name = "demo7"
    sub_type = "Device OS"
    type = "DEVICE"
    details {
        names {
            name = "Linux"
        }
    }
}

resource "prosimo_grouping" "newgrp5" {
    name = "demo7"
    type = "TIME"
    details {
        time {
            from = "2022-02-03 00:00:00"
            to = "2022-11-02 00:00:00"
            timezone = "Africa/Accra"
        }
    }
}

resource "prosimo_grouping" "newgrp6" {
    name = "demo7"
    type = "IP_RANGE"
    details {
        ranges = ["10.10.10.0/22"]
    }
}








