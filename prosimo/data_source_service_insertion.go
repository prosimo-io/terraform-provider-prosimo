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

func dataSourceServiceInsertion() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on service insertion.",
		ReadContext: dataSourceServiceInsertionRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = app_access_type==agent",
			},
			"service_insertion_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of configured/onboarded service insertion policies",
			},
			"service_insertions": {
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
						"regionid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"servicename": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serviceid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudtype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudregion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gwloadbalancerid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sharedservicecreds": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routetable": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"routingmanagedby": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resourcegroup": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prosimomanagedrouting": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vnetforpeering": {
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
						"iprules": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"srcaddr": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"srcport": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"destaddr": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"destport": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"protocol": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"createdtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updatedtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServiceInsertionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnServiceInsertionList []client.Service_Insertion

	ServiceInsertionList, err := prosimoClient.GetServiceInsertion(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)
	if filter != "" {
		for _, serviceInsertion := range ServiceInsertionList {
			var filteredMap *client.Service_Insertion

			err := mapstructure.Decode(serviceInsertion, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))

			if diags != nil {
				return diags
			}
			if flag {
				returnServiceInsertionList = append(returnServiceInsertionList, *serviceInsertion)
			}
		}

		if len(returnServiceInsertionList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, serviceInsertion := range ServiceInsertionList {
			returnServiceInsertionList = append(returnServiceInsertionList, *serviceInsertion)
		}
	}
	fmt.Println("Service Insertion List :", returnServiceInsertionList)
	d.SetId(time.Now().Format(time.RFC850))
	serviceInsertionItems := flattenServiceInsertionItemsData(returnServiceInsertionList)
	d.Set("service_insertions", serviceInsertionItems)
	d.Set("service_insertion_count", len(returnServiceInsertionList))
	return diags
}

func flattenServiceInsertionItemsData(ServiceInsertionItems []client.Service_Insertion) []interface{} {
	if ServiceInsertionItems != nil {
		ois := make([]interface{}, len(ServiceInsertionItems), len(ServiceInsertionItems))

		for i, ServiceInsertionItem := range ServiceInsertionItems {
			oi := make(map[string]interface{})
			oi["id"] = ServiceInsertionItem.ID
			oi["name"] = ServiceInsertionItem.Name
			oi["regionid"] = ServiceInsertionItem.RegionID
			oi["servicename"] = ServiceInsertionItem.Service_name
			oi["serviceid"] = ServiceInsertionItem.ServiceID
			oi["type"] = ServiceInsertionItem.Type
			oi["status"] = ServiceInsertionItem.Status
			oi["cloudtype"] = ServiceInsertionItem.CloudType
			oi["cloudregion"] = ServiceInsertionItem.CloudRegion
			oi["gwloadbalancerid"] = ServiceInsertionItem.GwLoadbalancerID
			oi["sharedservicecreds"] = ServiceInsertionItem.SharedServiceCreds
			oi["routetable"] = ServiceInsertionItem.RouteTable
			oi["routingmanagedby"] = ServiceInsertionItem.RoutingManagedBy
			oi["resourcegroup"] = ServiceInsertionItem.ResourceGroup
			oi["prosimomanagedrouting"] = ServiceInsertionItem.ProsimoManagedRouting
			oi["vnetforpeering"] = ServiceInsertionItem.VnetForPeering

			ServiceInsertionSource := make([]map[string]interface{}, 0)
			ServiceInsertionSourceTF := make(map[string]interface{})
			ServiceInsertionSourceTF["networks"] = flattenServiceInsertionServiceInputItemsData(ServiceInsertionItem.Source.Networks)
			ServiceInsertionSource = append(ServiceInsertionSource, ServiceInsertionSourceTF)
			oi["source"] = ServiceInsertionSource
			ServiceInsertionTarget := make([]map[string]interface{}, 0)
			ServiceInsertionTargetTF := make(map[string]interface{})
			if ServiceInsertionItem.Target.Networks != nil {
				ServiceInsertionTargetTF["networks"] = flattenServiceInsertionServiceInputItemsData(ServiceInsertionItem.Target.Networks)
			}
			if ServiceInsertionItem.Target.Apps != nil {
				ServiceInsertionTargetTF["apps"] = flattenServiceInsertionServiceInputItemsData(ServiceInsertionItem.Target.Apps)
			}
			ServiceInsertionTarget = append(ServiceInsertionTarget, ServiceInsertionTargetTF)
			oi["target"] = ServiceInsertionTarget
			if ServiceInsertionItem.IpRules != nil {
				oi["iprules"] = flattenServiceInsertionIpRuleItemsData(*ServiceInsertionItem.IpRules)
			} else {
				oi["iprules"] = make([]interface{}, 0)
			}

			oi["createdtime"] = ServiceInsertionItem.CreatedTime
			oi["updatedtime"] = ServiceInsertionItem.UpdatedTime

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenServiceInsertionServiceInputItemsData(ServiceInputItems []client.Service_Input) []interface{} {
	if ServiceInputItems != nil {
		ois := make([]interface{}, len(ServiceInputItems), len(ServiceInputItems))

		for i, ServiceInputItem := range ServiceInputItems {
			oi := make(map[string]interface{})
			oi["id"] = ServiceInputItem.ID
			oi["name"] = ServiceInputItem.Name
			oi["cidrs"] = ServiceInputItem.Cidrs

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenServiceInsertionIpRuleItemsData(IpRuleItems []client.IpRule) []interface{} {
	if IpRuleItems != nil {
		ois := make([]interface{}, len(IpRuleItems), len(IpRuleItems))

		for i, IpRuleItem := range IpRuleItems {
			oi := make(map[string]interface{})
			oi["id"] = IpRuleItem.ID
			oi["srcaddr"] = IpRuleItem.SrcAddr
			oi["srcport"] = IpRuleItem.SrcPort
			oi["destaddr"] = IpRuleItem.DestAddr
			oi["destport"] = IpRuleItem.DestPort
			oi["protocol"] = IpRuleItem.Protocol

			ois[i] = oi
			fmt.Println("IpRules oi 2:", oi)
		}

		return ois
	}
	return make([]interface{}, 0)
}
