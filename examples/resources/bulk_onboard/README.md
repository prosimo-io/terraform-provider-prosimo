<!-- BEGIN_AUTOMATED_TF_DOCS_BLOCK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 4.0 |
| <a name="requirement_prosimo"></a> [prosimo](#requirement\_prosimo) | 1.0.0 |
## Usage
Basic usage of this module is as follows:
```hcl
module "example" {

	 # Optional variables
	 NLB_name  = "qing-NLB-e074b26d5dbd1c34.elb.us-east-2.amazonaws.com"
	 TGW_id  = "tgw-04d5146c784e4d6d6"
	 agent_app_names  = [
  "agent",
  "Agent-bulk-webs-tf",
  "Agent-single-IP-tf",
  "Agent-bulk-apps-via-DNS-IP-tf",
  "Agent-single-fqdn-dns-IP-tf",
  "Agent-DNS-svr-tf",
  "Agent-single-fqdn-dns-app-tf",
  "Agent-DNS-bulk-tf",
  "agent-vwanhub-ssh-https-tf",
  "agent-transit-vnet-subnets-tf",
  "Agent-TGW-tf"
]
	 agentless_app_names  = [
  "agentless",
  "agentless-ssh-tf",
  "agentless-multi-VMs-tf",
  "url-rewrite-ssh-tf",
  "url-rewrite-https-tf",
  "jumpbox-peering-tf",
  "aws-s3-tf",
  "jumpbox-TGW-tf",
  "ssh-PrivateLink-tf",
  "vwanhub-ssh-tf",
  "vwanhub-jumpbox-tf"
]
	 app_list  = [
  "app1.psonar.local",
  "app2.psonar.local"
]
	 az_app_account  = "prosimo-app"
	 az_hub  = "qing-transit-hub-vnet"
	 az_region  = "eastus2"
	 az_resouce  = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing_transit_rg/providers/Microsoft.Network/virtualNetworks/"
	 az_spoke_1  = "transit_spoke1"
	 az_spoke_2  = "transit_spoke2"
	 cloud_creds_name  = "prosimo-aws-app-iam"
	 decommission_app  = false
	 domain_1  = "psonar.us"
	 domain_2  = "psonar.local"
	 hub  = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"
	 ip_pool_cidr  = "192.168.8.0/22"
	 onboard_app  = true
	 policy_name_agent  = "ALLOW-ALL"
	 policy_name_agentless  = "ALLOW-ALL"
	 port_list  = [
  "22-443"
]
	 prosimo_domain  = ".access.myeveqingchen1662950588667.scnetworkers.info"
	 pub_dns_server  = [
  "8.8.8.8",
  "1.1.1.1"
]
	 vm_ips  = [
  "10.100.0.22",
  "10.100.0.142",
  "10.101.0.141",
  "10.10.0.36",
  "172.17.2.4",
  "172.19.1.4"
]
	 vnet  = "/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vm-2-rg/providers/Microsoft.Network/virtualNetworks/qing-vm-2-rg-vnet"
	 vnet_ciders  = "172.19.1.0/24,172.19.2.0/24"
	 vpc_cider  = "10.100.0.0/24"
	 vpc_id_nlb  = "vpc-03bd678f593aa8866"
	 vpc_id_tgw  = "vpc-0ede62043db530015"
	 web_list  = [
  "fastly.com",
  "azure-api.us"
]
}
```
## Resources

| Name | Type |
|------|------|
| prosimo_app_onboarding_cloudsvc.app_s3 | resource |
| prosimo_app_onboarding_fqdn.Agent_DNS_bulk | resource |
| prosimo_app_onboarding_fqdn.Agent_bulk_apps_via_DNS_IP | resource |
| prosimo_app_onboarding_fqdn.Agent_bulk_webs | resource |
| prosimo_app_onboarding_fqdn.Agent_single_fqdn_dns_IP | resource |
| prosimo_app_onboarding_fqdn.Agent_single_fqdn_dns_app | resource |
| prosimo_app_onboarding_ip.Agent-single-IP | resource |
| prosimo_app_onboarding_ip.Agent_DNS_svr_tf | resource |
| prosimo_app_onboarding_ip.Agent_TGW | resource |
| prosimo_app_onboarding_ip.agent_vwanhub_ssh_https | resource |
| prosimo_app_onboarding_jumpbox.app_jumpbox | resource |
| prosimo_app_onboarding_jumpbox.vwanhub_jumpbox | resource |
| prosimo_app_onboarding_web.agentless_multi_VMs | resource |
| prosimo_app_onboarding_web.ssh_PrivateLink | resource |
| prosimo_app_onboarding_web.url_rewrite_https | resource |
| prosimo_app_onboarding_web.url_rewrite_ssh | resource |
| prosimo_app_onboarding_web.vwanhub_ssh | resource |
| prosimo_network_onboarding.aws-subnets | resource |
| prosimo_app_onboarding.R53_vwanhub | data source |
| prosimo_app_onboarding.bulk_internal_apps | data source |
| prosimo_s3bucket.s3_bucket | data source |
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_NLB_name"></a> [NLB\_name](#input\_NLB\_name) | NLB for PrivatrLink info | `string` | `"qing-NLB-e074b26d5dbd1c34.elb.us-east-2.amazonaws.com"` | no |
| <a name="input_TGW_id"></a> [TGW\_id](#input\_TGW\_id) | ##AWS TGW info | `string` | `"tgw-04d5146c784e4d6d6"` | no |
| <a name="input_agent_app_names"></a> [agent\_app\_names](#input\_agent\_app\_names) | n/a | `list` | <pre>[<br>  "agent",<br>  "Agent-bulk-webs-tf",<br>  "Agent-single-IP-tf",<br>  "Agent-bulk-apps-via-DNS-IP-tf",<br>  "Agent-single-fqdn-dns-IP-tf",<br>  "Agent-DNS-svr-tf",<br>  "Agent-single-fqdn-dns-app-tf",<br>  "Agent-DNS-bulk-tf",<br>  "agent-vwanhub-ssh-https-tf",<br>  "agent-transit-vnet-subnets-tf",<br>  "Agent-TGW-tf"<br>]</pre> | no |
| <a name="input_agentless_app_names"></a> [agentless\_app\_names](#input\_agentless\_app\_names) | n/a | `list` | <pre>[<br>  "agentless",<br>  "agentless-ssh-tf",<br>  "agentless-multi-VMs-tf",<br>  "url-rewrite-ssh-tf",<br>  "url-rewrite-https-tf",<br>  "jumpbox-peering-tf",<br>  "aws-s3-tf",<br>  "jumpbox-TGW-tf",<br>  "ssh-PrivateLink-tf",<br>  "vwanhub-ssh-tf",<br>  "vwanhub-jumpbox-tf"<br>]</pre> | no |
| <a name="input_app_list"></a> [app\_list](#input\_app\_list) | aws.amazon.com | `list` | <pre>[<br>  "app1.psonar.local",<br>  "app2.psonar.local"<br>]</pre> | no |
| <a name="input_az_app_account"></a> [az\_app\_account](#input\_az\_app\_account) | n/a | `string` | `"prosimo-app"` | no |
| <a name="input_az_hub"></a> [az\_hub](#input\_az\_hub) | n/a | `string` | `"qing-transit-hub-vnet"` | no |
| <a name="input_az_region"></a> [az\_region](#input\_az\_region) | n/a | `string` | `"eastus2"` | no |
| <a name="input_az_resouce"></a> [az\_resouce](#input\_az\_resouce) | #### Azure transit vNET info | `string` | `"/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing_transit_rg/providers/Microsoft.Network/virtualNetworks/"` | no |
| <a name="input_az_spoke_1"></a> [az\_spoke\_1](#input\_az\_spoke\_1) | n/a | `string` | `"transit_spoke1"` | no |
| <a name="input_az_spoke_2"></a> [az\_spoke\_2](#input\_az\_spoke\_2) | n/a | `string` | `"transit_spoke2"` | no |
| <a name="input_cloud_creds_name"></a> [cloud\_creds\_name](#input\_cloud\_creds\_name) | n/a | `string` | `"prosimo-aws-app-iam"` | no |
| <a name="input_decommission_app"></a> [decommission\_app](#input\_decommission\_app) | n/a | `bool` | `false` | no |
| <a name="input_domain_1"></a> [domain\_1](#input\_domain\_1) | n/a | `string` | `"psonar.us"` | no |
| <a name="input_domain_2"></a> [domain\_2](#input\_domain\_2) | n/a | `string` | `"psonar.local"` | no |
| <a name="input_hub"></a> [hub](#input\_hub) | n/a | `string` | `"/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vwan-rg/providers/Microsoft.Network/virtualHubs/qing-hub-useast-2"` | no |
| <a name="input_ip_pool_cidr"></a> [ip\_pool\_cidr](#input\_ip\_pool\_cidr) | n/a | `string` | `"192.168.8.0/22"` | no |
| <a name="input_onboard_app"></a> [onboard\_app](#input\_onboard\_app) | n/a | `bool` | `true` | no |
| <a name="input_policy_name_agent"></a> [policy\_name\_agent](#input\_policy\_name\_agent) | n/a | `string` | `"ALLOW-ALL"` | no |
| <a name="input_policy_name_agentless"></a> [policy\_name\_agentless](#input\_policy\_name\_agentless) | n/a | `string` | `"ALLOW-ALL"` | no |
| <a name="input_port_list"></a> [port\_list](#input\_port\_list) | n/a | `list` | <pre>[<br>  "22-443"<br>]</pre> | no |
| <a name="input_prosimo_domain"></a> [prosimo\_domain](#input\_prosimo\_domain) | n/a | `string` | `".access.myeveqingchen1662950588667.scnetworkers.info"` | no |
| <a name="input_pub_dns_server"></a> [pub\_dns\_server](#input\_pub\_dns\_server) | n/a | `list` | <pre>[<br>  "8.8.8.8",<br>  "1.1.1.1"<br>]</pre> | no |
| <a name="input_vm_ips"></a> [vm\_ips](#input\_vm\_ips) | ##"10.100.0.22"--DNS server, "10.100.0.142"--access VM, "10.101.0.141"--for PrivateLink in another VPC, ##"10.10.0.36"-- AWS TGW peering VM, "172.17.2.4"-- Azure access VM for vwanhub | `list` | <pre>[<br>  "10.100.0.22",<br>  "10.100.0.142",<br>  "10.101.0.141",<br>  "10.10.0.36",<br>  "172.17.2.4",<br>  "172.19.1.4"<br>]</pre> | no |
| <a name="input_vnet"></a> [vnet](#input\_vnet) | ##vwanhub info | `string` | `"/subscriptions/2de14016-6ebc-426e-848e-62a10837ce40/resourceGroups/qing-vm-2-rg/providers/Microsoft.Network/virtualNetworks/qing-vm-2-rg-vnet"` | no |
| <a name="input_vnet_ciders"></a> [vnet\_ciders](#input\_vnet\_ciders) | n/a | `string` | `"172.19.1.0/24,172.19.2.0/24"` | no |
| <a name="input_vpc_cider"></a> [vpc\_cider](#input\_vpc\_cider) | n/a | `string` | `"10.100.0.0/24"` | no |
| <a name="input_vpc_id_nlb"></a> [vpc\_id\_nlb](#input\_vpc\_id\_nlb) | n/a | `string` | `"vpc-03bd678f593aa8866"` | no |
| <a name="input_vpc_id_tgw"></a> [vpc\_id\_tgw](#input\_vpc\_id\_tgw) | n/a | `string` | `"vpc-0ede62043db530015"` | no |
| <a name="input_web_list"></a> [web\_list](#input\_web\_list) | n/a | `list` | <pre>[<br>  "fastly.com",<br>  "azure-api.us"<br>]</pre> | no |
## Outputs

| Name | Description |
|------|-------------|
| <a name="output_app_onboard_output_bulk_internal_apps_0"></a> [app\_onboard\_output\_bulk\_internal\_apps\_0](#output\_app\_onboard\_output\_bulk\_internal\_apps\_0) | applist |
| <a name="output_app_onboard_output_vwanhub"></a> [app\_onboard\_output\_vwanhub](#output\_app\_onboard\_output\_vwanhub) | applist |
<!-- END_AUTOMATED_TF_DOCS_BLOCK -->