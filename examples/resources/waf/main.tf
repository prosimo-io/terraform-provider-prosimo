resource "prosimo_waf" "waf" {
  waf_name        = "terraform-waf"
  mode = "enforce"
  threshold = 30
  
  rulesets {
      basic {
          rule_groups = ["11000_whitelist"]
      }

      owasp_crs_v32 {
          rule_groups = ["REQUEST-903.9001-DRUPAL-EXCLUSION-RULES"]
      }
  }

  app_domains = ["google.com"]

}


