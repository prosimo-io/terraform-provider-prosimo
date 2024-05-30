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

func dataSourceRegionalPrefix() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on summary route regional prefixes.",
		ReadContext: dataSourceRegionalPrefixRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = status==CONFIGURED",
			},
			"regional_prefix_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of configured/deployed summary route regional prefixes",
			},
			"regional_prefix_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefixes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"regions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"all": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"selected": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"csp": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"names": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"type": {
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

						// "credentials": {
						// 	Type:     schema.TypeSet,
						// 	Computed: true,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id": {
						// 				Type:     schema.TypeString,
						// 				Computed: true,
						// 			},
						// 			"cloud": {
						// 				Type:     schema.TypeString,
						// 				Computed: true,
						// 			},
						// 			"type": {
						// 				Type:     schema.TypeString,
						// 				Computed: true,
						// 			},
						// 			"credentials": {
						// 				Type:     schema.TypeString,
						// 				Computed: true,
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},
		},
	}
}

func dataSourceRegionalPrefixRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnRouteEntryList []client.Route_entry_region

	RouteEntryList, err := prosimoClient.GetRouteEntry(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)

	if filter != "" {
		for _, routeEntry := range RouteEntryList {
			var filteredMap *client.Route_entry_region

			err := mapstructure.Decode(routeEntry, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))
			if diags != nil {
				return diags
			}
			if flag {
				returnRouteEntryList = append(returnRouteEntryList, *routeEntry)
			}
		}
		if len(returnRouteEntryList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, routeEntry := range RouteEntryList {
			returnRouteEntryList = append(returnRouteEntryList, *routeEntry)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	routeEntryItems := flattenRouteEntryRegionItemsData(returnRouteEntryList)
	d.Set("regional_prefix_list", routeEntryItems)
	d.Set("regional_prefix_count", len(returnRouteEntryList))
	return diags
}

func flattenRouteEntryRegionItemsData(RouteEntryItems []client.Route_entry_region) []interface{} {
	if RouteEntryItems != nil {
		ois := make([]interface{}, len(RouteEntryItems), len(RouteEntryItems))

		for i, RouteEntryItem := range RouteEntryItems {
			oi := make(map[string]interface{})
			oi["id"] = RouteEntryItem.ID
			oi["status"] = RouteEntryItem.Status
			oi["prefixes"] = RouteEntryItem.Prefixes

			RouteEntryRegions := make([]map[string]interface{}, 0)
			RouteEntryRegionsTF := make(map[string]interface{})
			RouteEntryRegionsTF["all"] = RouteEntryItem.Regions.All
			RouteEntryRegionsTF["selected"] = flattenRouteEntrySelectedRegionItemsData(RouteEntryItem.Regions.Selected)
			RouteEntryRegions = append(RouteEntryRegions, RouteEntryRegionsTF)
			oi["regions"] = RouteEntryRegions
			oi["type"] = RouteEntryItem.Type
			oi["enabled"] = RouteEntryItem.Enabled
			oi["overwriteroute"] = RouteEntryItem.OverWriteRoute
			oi["teamid"] = RouteEntryItem.TeamID
			oi["createdtime"] = RouteEntryItem.CreatedTime
			oi["updatedtime"] = RouteEntryItem.UpdatedTime

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenRouteEntrySelectedRegionItemsData(SelectedRegionItems []client.Selected_Reg) []interface{} {
	if SelectedRegionItems != nil {
		ois := make([]interface{}, len(SelectedRegionItems), len(SelectedRegionItems))

		for i, SelectedRegionItem := range SelectedRegionItems {
			oi := make(map[string]interface{})
			oi["csp"] = SelectedRegionItem.CSP
			oi["names"] = SelectedRegionItem.Names

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
