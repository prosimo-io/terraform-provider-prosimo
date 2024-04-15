package prosimo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// var filtered map[string]string

func datasourcePolicyTransit() *schema.Resource {

	return &schema.Resource{
		Description: "Use this data source to get information on existing Transit Policies.",
		ReadContext: datasourcePolicyTransitRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = app_access_type==agent",
			},
			"policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"teamid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_posture_configured": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"details": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"actions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"matches": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"property": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"operations": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"values": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"inputitems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"selecteditems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
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
												"url": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"property": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"operations": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"values": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"inputitems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"selecteditems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
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
												"fqdn": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"property": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"operations": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"values": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"inputitems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"selecteditems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
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
												"advanced": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"property": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"operations": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"values": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selecteditems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"inputitems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
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
												"prosimonetworks": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"property": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"operations": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"values": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selecteditems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"inputitems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"itemid": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"itemname": {
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
												"networkacl": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"property": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"operations": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"values": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"inputitems": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keyvalues": {
																						Type:     schema.TypeSet,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"protocol": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem:     &schema.Schema{Type: schema.TypeString},
																								},
																								"sourceip": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem:     &schema.Schema{Type: schema.TypeString},
																								},
																								"sourceport": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem:     &schema.Schema{Type: schema.TypeString},
																								},
																								"targetip": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem:     &schema.Schema{Type: schema.TypeString},
																								},
																								"targetport": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem:     &schema.Schema{Type: schema.TypeString},
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
													},
												},
											},
										},
									},
									"apps": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"selecteditems": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"itemid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"itemname": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"networks": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"selecteditems": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"itemid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"itemname": {
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

func datasourcePolicyTransitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics
	var returnedPolicy []*client.Policy

	policyList, err := prosimoClient.GetPolicy(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	filter := d.Get("filter").(string)
	if filter != "" {
		for _, filteredList := range policyList {
			if filteredList.App_Access_Type == "transit" {
				diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredList))
				if diags != nil {
					return diags
				}
				if flag {
					returnedPolicy = append(returnedPolicy, filteredList)
				}
			}
		}
		if len(returnedPolicy) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, policyDetails := range policyList {
			if policyDetails.App_Access_Type == "transit" {
				returnedPolicy = append(returnedPolicy, policyDetails)
			}

		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	policyItems := flattenPolicyItemsData(returnedPolicy)
	d.Set("policy", policyItems)
	return diags

}

func flattenPolicyItemsData(PolicyItems []*client.Policy) []interface{} {
	if PolicyItems != nil {
		ois := make([]interface{}, len(PolicyItems), len(PolicyItems))

		for i, PolicyItem := range PolicyItems {
			oi := make(map[interface{}]interface{})
			oi["id"] = PolicyItem.ID
			oi["name"] = PolicyItem.Name
			oi["type"] = PolicyItem.Type
			oi["teamid"] = PolicyItem.TeamID
			oi["device_posture_configured"] = PolicyItem.Device_Posture_Configured
			oi["app_access_type"] = PolicyItem.App_Access_Type
			detailItems := flattenDetailItemsData(&PolicyItem.Details)
			oi["details"] = detailItems
			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}

func flattenDetailItemsData(DetailItems *client.Details) interface{} {
	if DetailItems != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		oi["actions"] = DetailItems.Actions
		if len(DetailItems.Matches.Advanced)+len(DetailItems.Matches.Device_Posture_Profile)+len(DetailItems.Matches.Devices)+len(DetailItems.Matches.FQDN)+len(DetailItems.Matches.IDP)+len(DetailItems.Matches.UserList)+len(DetailItems.Matches.Time)+len(DetailItems.Matches.Networks)+len(DetailItems.Matches.URL) > 0 {
			matchesItems := flattenMatchesItemsData(&DetailItems.Matches)
			oi["matches"] = matchesItems
		}
		if len(DetailItems.Apps.SelectedItems) > 0 {
			appsItems := flattenAppItems(DetailItems.Apps)
			oi["apps"] = appsItems
		}
		if len(DetailItems.Networks.SelectedItems) > 0 {
			networkItems := flattenAppItems(DetailItems.Networks)
			oi["networks"] = networkItems
		}
		ois = append(ois, oi)
		return ois
	}
	return make([]interface{}, 0)
}

func flattenMatchesItemsData(MatchDetails *client.MatchDetailList) interface{} {
	if MatchDetails != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		if MatchDetails.Time != nil {
			timeItems := flattenMatchEntriesItemsData(MatchDetails.Time)
			oi["time"] = timeItems
		}
		if MatchDetails.URL != nil {
			urlItems := flattenMatchEntriesItemsData(MatchDetails.URL)
			oi["url"] = urlItems
		}
		if MatchDetails.FQDN != nil {
			fqdnItems := flattenMatchEntriesItemsData(MatchDetails.FQDN)
			oi["fqdn"] = fqdnItems
		}
		if MatchDetails.Advanced != nil {
			advancedItems := flattenMatchEntriesItemsData(MatchDetails.Advanced)
			oi["advanced"] = advancedItems
		}
		if MatchDetails.NetworkACL != nil {
			networkACLItes := flattenMatchEntriesItemsData(MatchDetails.NetworkACL)
			oi["networkacl"] = networkACLItes
		}
		if MatchDetails.ProsimoNetworks != nil {
			prosimoNetworks := flattenMatchEntriesItemsData(MatchDetails.ProsimoNetworks)
			oi["prosimonetworks"] = prosimoNetworks
		}
		ois = append(ois, oi)
		return ois
	}
	return make([]interface{}, 0)
}

func flattenMatchEntriesItemsData(MatchEnteriesItems []client.MatchDetails) []interface{} {
	if MatchEnteriesItems != nil {
		ois := make([]interface{}, len(MatchEnteriesItems), len(MatchEnteriesItems))

		for i, MatchEnteriesItem := range MatchEnteriesItems {
			oi := make(map[interface{}]interface{})
			oi["property"] = MatchEnteriesItem.Property
			oi["operations"] = MatchEnteriesItem.Operations
			if len(MatchEnteriesItem.Values.InputItems)+len(MatchEnteriesItem.Values.SelectedItems)+len(MatchEnteriesItem.Values.SelectedItems) > 0 {
				valuesItem := flattenValuesItems(&MatchEnteriesItem.Values)
				oi["values"] = valuesItem
			}

			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}

func flattenValuesItems(ValuesItems *client.Values) interface{} {
	if ValuesItems != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		if ValuesItems.InputItems != nil {
			inputItems := flattenValuesItemData(ValuesItems.InputItems)
			oi["inputitems"] = inputItems
		}
		if ValuesItems.SelectedItems != nil {
			selectedItems := flattenValuesItemData(ValuesItems.SelectedItems)
			oi["selecteditems"] = selectedItems
		}
		if ValuesItems.SelectedGroups != nil {
			selectedGroups := flattenValuesItemData(ValuesItems.SelectedGroups)
			oi["selectedgroups"] = selectedGroups
		}
		ois = append(ois, oi)
		return ois
	}
	return make([]interface{}, 0)
}

func flattenAppItems(AppItems client.Values) interface{} {
	if &AppItems != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		selectedItems := flattenValuesItemData(AppItems.SelectedItems)
		oi["selecteditems"] = selectedItems
		ois = append(ois, oi)
		return ois
	}
	return make([]interface{}, 0)
}

func flattenValuesItemData(Items []client.InputItems) []interface{} {
	if Items != nil {
		ois := make([]interface{}, len(Items), len(Items))

		for i, Item := range Items {
			oi := make(map[interface{}]interface{})
			if Item.KeyValues != nil {
				keyvaluesItem := flattenkeyValuesItemData(Item.KeyValues)
				oi["keyvalues"] = keyvaluesItem
			} else {
				oi["itemid"] = Item.ItemID
				oi["	"] = Item.ItemName
			}

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenkeyValuesItemData(KeyValues *client.KeyValues) interface{} {
	if KeyValues != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		oi["protocol"] = KeyValues.Protocol
		oi["sourceip"] = KeyValues.SourceIp
		oi["sourceport"] = KeyValues.SourcePort
		oi["targetip"] = KeyValues.TargetIp
		oi["targetport"] = KeyValues.TargetPort

		ois = append(ois, oi)
		return ois
	}

	return make([]interface{}, 0)
}
