package prosimo

import (
	"context"
	"fmt"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourceNetworkDiscovered() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on discovered networks.",
		ReadContext: dataSourceNetworkDiscoveredRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"discovered_networks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudtype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accountname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"regions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_creds_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"vpcs": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cidr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"network": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"subnet_count": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkDiscoveredRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	var returnfilteredNetworks []*client.DiscoveredNetworks

	discoveredNetworkList, err := prosimoClient.GetDiscoveredNetworks(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	filter := d.Get("filter").(string)
	if filter != "" {
		for _, filteredDNList := range discoveredNetworkList {
			filteredMap := map[string]interface{}{}
			err := mapstructure.Decode(filteredDNList, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, filteredMap)
			if diags != nil {
				return diags
			}
			if flag {
				returnfilteredNetworks = append(returnfilteredNetworks, filteredDNList)
			}
		}
		if len(returnfilteredNetworks) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})

			return diags
		}
	} else {
		for _, filteredDNList := range discoveredNetworkList {
			returnfilteredNetworks = append(returnfilteredNetworks, filteredDNList)
		}
	}

	// if len(cloud_type) > 0 && len(account_name) > 0 {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Error,
	// 		Summary:  "Invalid Input, either of input_cloud_type/input_account_name is expected",
	// 		Detail:   fmt.Sprintln("Invalid Input, please enter either of the mentioned inputs: input_cloud_type/input_account_name"),
	// 	})

	// 	return diags
	// }

	// if len(cloud_type) > 0 {
	// 	flag := false
	// 	for _, discoveredNetworkDetails := range discoveredNetworkList {
	// 		if discoveredNetworkDetails.CloudType == cloud_type {
	// 			returnfilteredNetworks = append(returnfilteredNetworks, discoveredNetworkDetails)
	// 			// fmt.Println("Discovered Networks:", returnfilteredNetworks)
	// 			flag = true
	// 		}
	// 	}
	// 	if !flag {
	// 		diags = append(diags, diag.Diagnostic{
	// 			Severity: diag.Error,
	// 			Summary:  "Cloud Type does not exists",
	// 			Detail:   fmt.Sprintln("Given Cloud Type does not exists"),
	// 		})

	// 		return diags
	// 	}
	// } else if len(account_name) > 0 {
	// 	flag := false
	// 	for _, discodiscoveredNetworkDetails := range discoveredNetworkList {
	// 		if discodiscoveredNetworkDetails.AccountName == account_name {
	// 			returnfilteredNetworks = append(returnfilteredNetworks, discodiscoveredNetworkDetails)
	// 			// fmt.Println("Discovered Networks:", returnfilteredNetworks)
	// 			flag = true
	// 		}
	// 	}
	// 	if !flag {
	// 		diags = append(diags, diag.Diagnostic{
	// 			Severity: diag.Error,
	// 			Summary:  "Account Name does not exists",
	// 			Detail:   fmt.Sprintln("Given Account Name does not exists"),
	// 		})

	// 	}
	// } else {
	// 	for _, discodiscoveredNetworkDetails := range discoveredNetworkList {
	// 		returnfilteredNetworks = append(returnfilteredNetworks, discodiscoveredNetworkDetails)
	// 	}
	// }

	d.SetId(time.Now().Format(time.RFC850))
	discoveredNetworkItems := flattenDNItemsData(returnfilteredNetworks)
	d.Set("discovered_networks", discoveredNetworkItems)
	return diags

}

func flattenDNItemsData(DiscoveredNetworkItems []*client.DiscoveredNetworks) []interface{} {
	if DiscoveredNetworkItems != nil {
		ois := make([]interface{}, len(DiscoveredNetworkItems), len(DiscoveredNetworkItems))

		for i, DiscoveredNetworkItem := range DiscoveredNetworkItems {
			oi := make(map[interface{}]interface{})
			oi["id"] = DiscoveredNetworkItem.ID
			oi["name"] = DiscoveredNetworkItem.Name
			oi["cloudtype"] = DiscoveredNetworkItem.CloudType
			oi["accountname"] = DiscoveredNetworkItem.AccountName
			regionItems := flattenRegionItemsData(DiscoveredNetworkItem.Regions)
			oi["regions"] = regionItems
			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}

func flattenRegionItemsData(RegionItems []client.Regions) []interface{} {
	if RegionItems != nil {
		ois := make([]interface{}, len(RegionItems), len(RegionItems))

		for i, RegionItem := range RegionItems {
			oi := make(map[interface{}]interface{})
			oi["id"] = RegionItem.ID
			oi["cloud_creds_id"] = RegionItem.CloudCredsID
			oi["name"] = RegionItem.Name
			oi["vpc_count"] = RegionItem.VpcCount
			vpcsList := flattenVpcsItemsData(RegionItem.Vpcs)
			oi["vpcs"] = vpcsList
			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}

func flattenVpcsItemsData(VpcsItems []client.VPCs) []interface{} {
	if VpcsItems != nil {
		ois := make([]interface{}, len(VpcsItems), len(VpcsItems))

		for i, VpcsItem := range VpcsItems {
			oi := make(map[interface{}]interface{})
			oi["id"] = VpcsItem.ID
			oi["region_id"] = VpcsItem.RegionID
			oi["cidr"] = VpcsItem.Cidr
			oi["network"] = VpcsItem.Network
			oi["name"] = VpcsItem.Name
			oi["subnet_count"] = VpcsItem.SubnetCount
			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}
