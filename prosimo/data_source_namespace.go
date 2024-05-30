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

func dataSourceNamespace() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on namespaces.",
		ReadContext: dataSourceNamespaceRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = region==us-west-1",
			},
			"namespace_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of configured namespaces",
			},
			"namespace_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
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
						"teamid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignednetworks": {
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
												"teamid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"pamcname": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"policyupdated": {
													Type:     schema.TypeBool,
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
												"progress": {
													Type:     schema.TypeInt,
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
												"namespaceid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"namespacename": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"namespacenid": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"exportable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"publiccloud": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloud": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudtype": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"connectionoption": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudkeyid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudregion": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudnetworks": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"cloudnetworkid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectorgrpid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"edgeconnectivityid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"hubid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectivitytype": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"subnets": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"subnet": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"virtual_subnet": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"connectorplacement": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectorsettings": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			// Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bandwidth": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"bandwidthname": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"instancetype": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"cloudnetworkid": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"updatestatus": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"subnets": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Optional: true,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"bandwidthrange": {
																						Type:     schema.TypeSet,
																						Computed: true,
																						Optional: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"min": {
																									Type:     schema.TypeInt,
																									Computed: true,
																								},
																								"max": {
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
													},
												},
												"security": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"policies": {
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
						"exportednetworks": {
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
												"teamid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"pamcname": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"policyupdated": {
													Type:     schema.TypeBool,
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
												"progress": {
													Type:     schema.TypeInt,
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
												"namespaceid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"namespacename": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"namespacenid": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"exportable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"publiccloud": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloud": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudtype": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"connectionoption": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudkeyid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudregion": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudnetworks": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"cloudnetworkid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectorgrpid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"edgeconnectivityid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"hubid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectivitytype": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"subnets": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"subnet": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"virtual_subnet": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"connectorplacement": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectorsettings": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			// Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bandwidth": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"bandwidthname": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"instancetype": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"cloudnetworkid": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"updatestatus": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"subnets": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Optional: true,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"bandwidthrange": {
																						Type:     schema.TypeSet,
																						Computed: true,
																						Optional: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"min": {
																									Type:     schema.TypeInt,
																									Computed: true,
																								},
																								"max": {
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
													},
												},
												"security": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"policies": {
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
						"importednetworks": {
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
												"teamid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"pamcname": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"policyupdated": {
													Type:     schema.TypeBool,
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
												"progress": {
													Type:     schema.TypeInt,
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
												"namespaceid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"namespacename": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"namespacenid": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"exportable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"publiccloud": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloud": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudtype": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"connectionoption": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudkeyid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudregion": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"cloudnetworks": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"cloudnetworkid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectorgrpid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"edgeconnectivityid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"hubid": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectivitytype": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"subnets": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"subnet": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"virtual_subnet": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"connectorplacement": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"connectorsettings": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			// Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bandwidth": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"bandwidthname": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"instancetype": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"cloudnetworkid": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"updatestatus": {
																						Type:     schema.TypeString,
																						Computed: true,
																						Optional: true,
																					},
																					"subnets": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Optional: true,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"bandwidthrange": {
																						Type:     schema.TypeSet,
																						Computed: true,
																						Optional: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"min": {
																									Type:     schema.TypeInt,
																									Computed: true,
																								},
																								"max": {
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
													},
												},
												"security": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"policies": {
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
					},
				},
			},
		},
	}
}

func dataSourceNamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnNamespaceList []client.Namespace

	NamespaceList, err := prosimoClient.GetNamespace(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)

	if filter != "" {
		for _, namespace := range NamespaceList {
			var filteredMap *client.Namespace

			err := mapstructure.Decode(namespace, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))
			if diags != nil {
				return diags
			}
			if flag {
				returnNamespaceList = append(returnNamespaceList, *namespace)
			}
		}
		if len(returnNamespaceList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, namespace := range NamespaceList {
			returnNamespaceList = append(returnNamespaceList, *namespace)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	namespaceItems := flattenNamespaceItemsData(returnNamespaceList)
	d.Set("namespace_list", namespaceItems)
	d.Set("namespace_count", len(returnNamespaceList))
	return diags
}

func flattenNamespaceItemsData(NamespaceItems []client.Namespace) []interface{} {
	if NamespaceItems != nil {
		ois := make([]interface{}, len(NamespaceItems), len(NamespaceItems))

		for i, NamespaceItem := range NamespaceItems {
			oi := make(map[string]interface{})
			oi["id"] = NamespaceItem.ID
			oi["nid"] = NamespaceItem.NID
			oi["name"] = NamespaceItem.Name
			oi["createdtime"] = NamespaceItem.CreatedTime
			oi["updatedtime"] = NamespaceItem.UpdatedTime
			oi["teamid"] = NamespaceItem.TeamID

			AssignedNetworks := make([]map[string]interface{}, 0)
			AssignedNetworksTF := make(map[string]interface{})
			AssignedNetworksTF["networks"] = flattenNsNetworkItemsData(*NamespaceItem.AssignedNetworks.Networks)
			AssignedNetworks = append(AssignedNetworks, AssignedNetworksTF)
			oi["assignednetworks"] = AssignedNetworks

			ExportedNetworks := make([]map[string]interface{}, 0)
			ExportedNetworksTF := make(map[string]interface{})
			ExportedNetworksTF["networks"] = flattenNsNetworkItemsData(*NamespaceItem.ExportedNetworks.Networks)
			ExportedNetworks = append(ExportedNetworks, ExportedNetworksTF)
			oi["exportednetworks"] = ExportedNetworks

			ImportedNetworks := make([]map[string]interface{}, 0)
			ImportedNetworksTF := make(map[string]interface{})
			ImportedNetworksTF["networks"] = flattenNsNetworkItemsData(*NamespaceItem.ImportedNetworks.Networks)
			ImportedNetworks = append(ImportedNetworks, ExportedNetworksTF)
			oi["importednetworks"] = ImportedNetworks

			oi["status"] = NamespaceItem.Status

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenNsNetworkItemsData(NsNetworkItems []client.NS_Networks) []interface{} {
	if NsNetworkItems != nil {
		ois := make([]interface{}, len(NsNetworkItems), len(NsNetworkItems))

		for i, NsNetworkItem := range NsNetworkItems {
			oi := make(map[string]interface{})
			oi["id"] = NsNetworkItem.ID
			oi["name"] = NsNetworkItem.Name
			oi["teamid"] = NsNetworkItem.TeamID
			oi["pamcname"] = NsNetworkItem.PamCname
			oi["policyupdated"] = NsNetworkItem.PolicyUpdated
			oi["deployed"] = NsNetworkItem.Deployed
			oi["status"] = NsNetworkItem.Status
			oi["progress"] = NsNetworkItem.Progress
			oi["createdtime"] = NsNetworkItem.CreatedTime
			oi["updatedtime"] = NsNetworkItem.UpdatedTime
			oi["namespaceid"] = NsNetworkItem.NamespaceID
			oi["namespacename"] = NsNetworkItem.NamespaceName
			oi["namespacenid"] = NsNetworkItem.NamespaceNID
			oi["exportable"] = NsNetworkItem.Exportable
			publicCloudItems := flattenPublicCloudItemsData(NsNetworkItem.PublicCloud)
			oi["publiccloud"] = publicCloudItems
			securityItems := flattenSecurityItemsData(NsNetworkItem.Security)
			oi["security"] = securityItems
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
