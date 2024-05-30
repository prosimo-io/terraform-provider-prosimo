package prosimo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourceNetworkPrefix() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on summary route network prefixes.",
		ReadContext: dataSourceNetworkPrefixRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = cloudregion==us-west-2",
			},
			"network_prefix_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of configured/deployed summary route network prefixes",
			},
			"network_prefix_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudkeyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"csp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudregion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudnetworkid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudnetworkname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"overwriteroute": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"teamid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"createdtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updatedtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefixroutetables": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"routetables": {
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

func dataSourceNetworkPrefixRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnRouteEntryNetworkList []client.Route_entry_network

	RouteEntryNetworkList, err := prosimoClient.GetNetworkRouteEntry(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)

	if filter != "" {
		for _, routeEntryNetwork := range RouteEntryNetworkList {
			var filteredMap *client.Route_entry_network

			err := mapstructure.Decode(routeEntryNetwork, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))
			if diags != nil {
				return diags
			}
			if flag {
				returnRouteEntryNetworkList = append(returnRouteEntryNetworkList, *routeEntryNetwork)
			}
		}
		if len(returnRouteEntryNetworkList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, routeEntryNetwork := range RouteEntryNetworkList {
			returnRouteEntryNetworkList = append(returnRouteEntryNetworkList, *routeEntryNetwork)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	routeEntryNetworkItems := flattenRouteEntryNetworkItemsData(returnRouteEntryNetworkList)
	d.Set("network_prefix_list", routeEntryNetworkItems)
	d.Set("network_prefix_count", len(returnRouteEntryNetworkList))
	return diags
}

func flattenRouteEntryNetworkItemsData(RouteEntryNetworkItems []client.Route_entry_network) []interface{} {
	if RouteEntryNetworkItems != nil {
		ois := make([]interface{}, len(RouteEntryNetworkItems), len(RouteEntryNetworkItems))

		for i, RouteEntryNetworkItem := range RouteEntryNetworkItems {
			oi := make(map[string]interface{})
			oi["id"] = RouteEntryNetworkItem.ID
			oi["cloudkeyid"] = RouteEntryNetworkItem.CloudKeyID
			oi["csp"] = RouteEntryNetworkItem.CSP
			oi["cloudregion"] = RouteEntryNetworkItem.CloudRegion
			oi["cloudnetworkid"] = RouteEntryNetworkItem.CloudNetworkID
			oi["cloudnetworkname"] = RouteEntryNetworkItem.CloudNetworkName
			oi["enabled"] = RouteEntryNetworkItem.Enabled
			oi["overwriteroute"] = RouteEntryNetworkItem.OverwriteRoute
			oi["status"] = RouteEntryNetworkItem.Status
			oi["teamid"] = RouteEntryNetworkItem.TeamID
			oi["prefixroutetables"] = flattenPrefixRouteIDItemsData(*RouteEntryNetworkItem.PrefixesRT)
			oi["createdtime"] = RouteEntryNetworkItem.CreatedTime
			oi["updatedtime"] = RouteEntryNetworkItem.UpdatedTime

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenPrefixRouteIDItemsData(PrefixRouteIDItems []client.Prefix_RouteID) []interface{} {
	if PrefixRouteIDItems != nil {
		ois := make([]interface{}, len(PrefixRouteIDItems), len(PrefixRouteIDItems))

		for i, PrefixRouteIDItem := range PrefixRouteIDItems {
			oi := make(map[string]interface{})
			oi["prefix"] = PrefixRouteIDItem.Prefix
			oi["routetables"] = flattenPrefixRouteTableItemsData(PrefixRouteIDItem.RouteTableIDS)

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenPrefixRouteTableItemsData(PrefixRouteTableItems []client.Route_Tables) []interface{} {
	if PrefixRouteTableItems != nil {
		ois := make([]interface{}, len(PrefixRouteTableItems), len(PrefixRouteTableItems))

		for i, PrefixRouteTableItem := range PrefixRouteTableItems {
			oi := make(map[string]interface{})
			oi["id"] = PrefixRouteTableItem.ID
			oi["name"] = PrefixRouteTableItem.Name

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
