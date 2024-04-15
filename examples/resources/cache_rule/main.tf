resource "prosimo_cache_rules" "cacherule" {
    name = "demo80"
    default = true
    editable = true
    share_static_content = false
    cache_control_ignored = true
    bypass_cache = true

    path_patterns {
        path = "*"
        bypass_uri = false
        is_default = false
        is_new_path = false
        status = "existing"
         settings {
            type = "Dynamic"
            user_id_ignored = false
            cache_control_ignored = false
            cookie_ignored = false
            query_parameter_ignored = false 
            ttl {
                enabled = true
                time = 24
                time_unit = "Hours"
            }


        }

    }
    path_patterns {
        path = "/abcdefg"
        bypass_uri = false
        is_default = false
        is_new_path = false
        status = "existing"

        settings {
            type = "Dynamic"
            user_id_ignored = false
            cache_control_ignored = false
            cookie_ignored = false
            query_parameter_ignored = false 
            ttl {
                enabled = true
                time = 24
                time_unit = "Hours"
            }


        }


    }
    bypass_info {
		resp_hdrs {
			x_jenkins_session = [""]
		}
    }
    app_domains {
        domain = "speedtest-server-eastus2-1625128238856.myeventarena.com"
    }
}