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

func datasourcePolicyAccess() *schema.Resource {

	return &schema.Resource{
		Description: "Use this data source to get information on existing access policies.",
		ReadContext: datasourcePolicyAccessRead,
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
												"networks": {
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
																		"selecteditems": {
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
												"devices": {
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
																		"inputitems": {
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
												"users": {
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
																		"selecteditems": {
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
												"idp": {
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
																		"selecteditems": {
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
																		"selecteditems": {
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
																		"selecteditems": {
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
																		"selecteditems": {
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
																		"inputitems": {
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
												"devicepostureprofiles": {
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
																		"inputitems": {
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

func datasourcePolicyAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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
			if filteredList.App_Access_Type == "access" {
				filteredMap := map[string]interface{}{}
				err := mapstructure.Decode(filteredList, &filteredMap)
				if err != nil {
					panic(err)
				}
				diags, flag := checkMainOperand(filter, filteredMap)
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
			if policyDetails.App_Access_Type == "access" {
				returnedPolicy = append(returnedPolicy, policyDetails)
			}

		}
	}
	d.SetId(time.Now().Format(time.RFC850))
	policyItems := flattenAccessPolicyItemsData(returnedPolicy)
	d.Set("policy", policyItems)
	return diags
}

func flattenAccessPolicyItemsData(PolicyItems []*client.Policy) []interface{} {
	if PolicyItems != nil {
		ois := make([]interface{}, len(PolicyItems), len(PolicyItems))

		for i, PolicyItem := range PolicyItems {
			oi := make(map[interface{}]interface{})
			oi["id"] = PolicyItem.ID
			oi["name"] = PolicyItem.DisplayName
			oi["type"] = PolicyItem.Type
			oi["teamid"] = PolicyItem.TeamID
			oi["app_access_type"] = PolicyItem.App_Access_Type
			detailItems := flattenDetailItemsDataAccess(&PolicyItem.Details)
			oi["details"] = detailItems
			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}
func flattenDetailItemsDataAccess(DetailItems *client.Details) interface{} {
	if DetailItems != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		oi["actions"] = DetailItems.Actions
		if len(DetailItems.Matches.Advanced)+len(DetailItems.Matches.Device_Posture_Profile)+len(DetailItems.Matches.Devices)+len(DetailItems.Matches.FQDN)+len(DetailItems.Matches.IDP)+len(DetailItems.Matches.UserList)+len(DetailItems.Matches.Time)+len(DetailItems.Matches.Networks)+len(DetailItems.Matches.URL) > 0 {
			matchesItems := flattenMatchesItemsDataAccess(&DetailItems.Matches)
			oi["matches"] = matchesItems
		}
		if len(DetailItems.Apps.SelectedItems) > 0 {
			appsItems := flattenAppItemsAccess(DetailItems.Apps)
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
func flattenMatchesItemsDataAccess(MatchDetails *client.MatchDetailList) interface{} {
	if MatchDetails != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		if MatchDetails.Networks != nil {
			matchentriesItems := flattenMatchEntriesItemsData(MatchDetails.Networks)
			oi["networks"] = matchentriesItems
		}
		if MatchDetails.UserList != nil {
			userlistItems := flattenMatchEntriesItemsData(MatchDetails.UserList)
			oi["users"] = userlistItems
		}
		if MatchDetails.Device_Posture_Profile != nil {
			devicepostureprofilesItems := flattenMatchEntriesItemsData(MatchDetails.Device_Posture_Profile)
			oi["devicepostureprofiles"] = devicepostureprofilesItems
		}
		if MatchDetails.Devices != nil {
			devicesItems := flattenMatchEntriesItemsData(MatchDetails.Devices)
			oi["devices"] = devicesItems
		}
		if MatchDetails.IDP != nil {
			idpItems := flattenMatchEntriesItemsData(MatchDetails.IDP)
			oi["idp"] = idpItems
		}
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
		ois = append(ois, oi)
		return ois
	}
	return make([]interface{}, 0)
}

func flattenAppItemsAccess(AppItems client.Values) interface{} {
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
