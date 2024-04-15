data "prosimo_waf_policy" "waf_policy_list" {
    input_name = "Default WAF Policy"
    # input_mode = "Qing-agent-allow"
}

output "waf_policy_details" {
  description = "WAF policy details"
  value       = data.prosimo_waf_policy.waf_policy_list
}