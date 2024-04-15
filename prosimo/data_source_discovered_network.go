package prosimo

import (
	"context"
	"reflect"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"input_cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_account_name": {
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
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredDNList))
			if diags != nil {
				return diags
			}
			if flag {
				returnfilteredNetworks = append(returnfilteredNetworks, filteredDNList)
			}
		}
	} else {
		for _, filteredDNList := range discoveredNetworkList {
			returnfilteredNetworks = append(returnfilteredNetworks, filteredDNList)
		}
	}
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
