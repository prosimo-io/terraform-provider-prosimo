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

func dataSourceSharedService() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on shared services.",
		ReadContext: dataSourceSharedServiceRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = status==DEPLOYED",
			},
			"shared_service_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of configured/deployed Shared Services",
			},
			"shared_service_list": {
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
						"deployed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"serviceinsert": {
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
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"regionid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"networks": {
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
															"cidrs": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: false,
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
									"target": {
										Type:     schema.TypeSet,
										Computed: true,
										// Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"networks": {
													Type:     schema.TypeSet,
													Computed: false,
													Optional: true,
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
															"cidrs": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"apps": {
													Type:     schema.TypeSet,
													Computed: false,
													Optional: true,
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
															"cidrs": {
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
									"namespacenid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"namespaceid": {
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
	}
}

func dataSourceSharedServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnSharedServiceList []client.Shared_Service

	SharedServiceList, err := prosimoClient.GetSharedService(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)

	if filter != "" {
		for _, sharedService := range SharedServiceList {
			var filteredMap *client.Shared_Service

			err := mapstructure.Decode(sharedService, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))
			if diags != nil {
				return diags
			}
			if flag {
				returnSharedServiceList = append(returnSharedServiceList, *sharedService)
			}
		}
		if len(returnSharedServiceList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, sharedService := range SharedServiceList {
			returnSharedServiceList = append(returnSharedServiceList, *sharedService)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	sharedServiceItems := flattenSharedServiceItemsData(returnSharedServiceList)
	d.Set("shared_service_list", sharedServiceItems)
	d.Set("shared_service_count", len(returnSharedServiceList))
	return diags
}

func flattenSharedServiceItemsData(SharedServiceItems []client.Shared_Service) []interface{} {
	if SharedServiceItems != nil {
		ois := make([]interface{}, len(SharedServiceItems), len(SharedServiceItems))

		for i, SharedServiceItem := range SharedServiceItems {
			oi := make(map[string]interface{})
			oi["id"] = SharedServiceItem.ID
			oi["name"] = SharedServiceItem.Name
			oi["deployed"] = SharedServiceItem.Deployed
			oi["status"] = SharedServiceItem.Status
			oi["type"] = SharedServiceItem.Type
			if SharedServiceItem.ServiceInsert != nil {
				oi["serviceinsert"] = flattenServiceInsertItemsData(*SharedServiceItem.ServiceInsert)
			}
			oi["createdtime"] = SharedServiceItem.CreatedTime
			oi["updatedtime"] = SharedServiceItem.UpdatedTime
			oi["progress"] = SharedServiceItem.Progress

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenServiceInsertItemsData(ServiceInsertItems []client.ServiceInsert) []interface{} {
	if ServiceInsertItems != nil {
		ois := make([]interface{}, len(ServiceInsertItems), len(ServiceInsertItems))

		for i, ServiceInsertItem := range ServiceInsertItems {
			oi := make(map[string]interface{})
			oi["id"] = ServiceInsertItem.ID
			oi["name"] = ServiceInsertItem.Name
			oi["type"] = ServiceInsertItem.Type
			oi["regionid"] = ServiceInsertItem.RegionID

			ServiceInsertionSource := make([]map[string]interface{}, 0)
			ServiceInsertionSourceTF := make(map[string]interface{})
			ServiceInsertionSourceTF["networks"] = flattenServiceInsertionServiceInputItemsData(ServiceInsertItem.Source.Networks)
			ServiceInsertionSource = append(ServiceInsertionSource, ServiceInsertionSourceTF)
			oi["source"] = ServiceInsertionSource
			ServiceInsertionTarget := make([]map[string]interface{}, 0)
			ServiceInsertionTargetTF := make(map[string]interface{})
			if ServiceInsertItem.Target.Networks != nil {
				ServiceInsertionTargetTF["networks"] = flattenServiceInsertionServiceInputItemsData(ServiceInsertItem.Target.Networks)
			}
			if ServiceInsertItem.Target.Apps != nil {
				ServiceInsertionTargetTF["apps"] = flattenServiceInsertionServiceInputItemsData(ServiceInsertItem.Target.Apps)
			}
			if len(ServiceInsertionTargetTF) > 0 {
				ServiceInsertionTarget = append(ServiceInsertionTarget, ServiceInsertionTargetTF)
			}
			oi["target"] = ServiceInsertionTarget

			oi["namespacenid"] = ServiceInsertItem.NamespaceNID
			oi["namespaceid"] = ServiceInsertItem.NamespaceID

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
