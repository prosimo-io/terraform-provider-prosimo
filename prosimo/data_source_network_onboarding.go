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

func dataSourceNetworkOnboarding() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on onboarded networks.",
		ReadContext: dataSourceNetworkOnboardingRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = app_access_type==agent",
			},
			"network_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of Onboarded Networks",
			},
			"onboarded_networks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
						"status": {
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
						"policy_updated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_nid": {
							Type:     schema.TypeInt,
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
	}
}

func dataSourceNetworkOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnNetworkList []client.NetworkOnboardoptns

	onboardNetworkList, err := prosimoClient.SearchOnboardNetworks(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)

	if filter != "" {
		for _, onboardNetwork := range onboardNetworkList.Data.Records {
			var filteredMap *client.NetworkOnboardoptns

			err := mapstructure.Decode(onboardNetwork, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))
			if diags != nil {
				return diags
			}
			if flag {
				returnNetworkList = append(returnNetworkList, *onboardNetwork)
			}
		}
		if len(returnNetworkList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, onboardNetwork := range onboardNetworkList.Data.Records {
			returnNetworkList = append(returnNetworkList, *onboardNetwork)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	networkItems := flattenNetworkItemsData(returnNetworkList)
	d.Set("onboarded_networks", networkItems)
	d.Set("network_count", len(returnNetworkList))
	return diags
}

func flattenNetworkItemsData(NetworkItems []client.NetworkOnboardoptns) []interface{} {
	if NetworkItems != nil {
		ois := make([]interface{}, len(NetworkItems), len(NetworkItems))

		for i, NetworkItem := range NetworkItems {
			oi := make(map[string]interface{})
			oi["name"] = NetworkItem.Name
			oi["id"] = NetworkItem.ID
			oi["teamid"] = NetworkItem.TeamID
			oi["status"] = NetworkItem.Status
			oi["pamcname"] = NetworkItem.PamCname
			oi["createdtime"] = NetworkItem.CreatedTime
			oi["updatedtime"] = NetworkItem.UpdatedTime
			oi["policy_updated"] = NetworkItem.PolicyUpdated
			oi["namespace_name"] = NetworkItem.NamespaceName
			oi["namespace_id"] = NetworkItem.NamespaceID
			oi["namespace_nid"] = NetworkItem.NamespaceNID
			publicCloudItems := flattenPublicCloudItemsData(NetworkItem.PublicCloud)
			oi["publiccloud"] = publicCloudItems
			securityItems := flattenSecurityItemsData(NetworkItem.Security)
			oi["security"] = securityItems

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenPublicCloudItemsData(PublicCloudItems *client.PublicCloud) interface{} {
	if PublicCloudItems != nil {
		ois := make([]map[string]interface{}, 0)
		// for i, PublicCloudItem := range PublicCloudItems {
		oi := make(map[string]interface{})

		oi["id"] = PublicCloudItems.Id
		oi["cloud"] = PublicCloudItems.Cloud
		oi["cloudtype"] = PublicCloudItems.CloudType
		oi["connectionoption"] = PublicCloudItems.ConnectionOption
		oi["cloudkeyid"] = PublicCloudItems.CloudKeyID
		oi["cloudregion"] = PublicCloudItems.CloudRegion
		cloudNetworkItems := flattenCloudNetworksItemsData(PublicCloudItems.CloudNetworks)
		oi["cloudnetworks"] = cloudNetworkItems

		ois = append(ois, oi)
		return ois
	}
	return make([]interface{}, 0)
}

func flattenCloudNetworksItemsData(CloudNetworkItems []client.CloudNetworkops) []interface{} {
	if CloudNetworkItems != nil {
		ois := make([]interface{}, len(CloudNetworkItems), len(CloudNetworkItems))

		for i, CloudNetworkItem := range CloudNetworkItems {
			oi := make(map[string]interface{})

			oi["id"] = CloudNetworkItem.Id
			oi["cloudnetworkid"] = CloudNetworkItem.CloudNetworkID
			oi["connectorgrpid"] = CloudNetworkItem.ConnectorGrpID
			oi["edgeconnectivityid"] = CloudNetworkItem.EdgeConnectivityID
			oi["hubid"] = CloudNetworkItem.HubID
			oi["connectivitytype"] = CloudNetworkItem.ConnectivityType
			subnets := flattenSubnetItemsData(CloudNetworkItem.Subnets)
			oi["subnets"] = subnets
			oi["connectorplacement"] = CloudNetworkItem.ConnectorPlacement

			connectorSettings := make([]map[string]interface{}, 0)
			connectorSettingTF := make(map[string]interface{})
			bandwithRange := make([]map[string]interface{}, 0)
			bandwithRangeTF := make(map[string]interface{})
			if CloudNetworkItem.Connectorsettings != nil {
				connectorSettingTF["bandwidth"] = CloudNetworkItem.Connectorsettings.Bandwidth
				connectorSettingTF["bandwidthname"] = CloudNetworkItem.Connectorsettings.BandwidthName
				connectorSettingTF["instancetype"] = CloudNetworkItem.Connectorsettings.InstanceType
				connectorSettingTF["cloudnetworkid"] = CloudNetworkItem.Connectorsettings.CloudNetworkID
				connectorSettingTF["updatestatus"] = CloudNetworkItem.Connectorsettings.UpdateStatus
				connectorSettingTF["subnets"] = CloudNetworkItem.Connectorsettings.Subnets
				bandwidthRange := CloudNetworkItem.Connectorsettings.BandwidthRange
				if bandwidthRange != nil {
					bandwithRangeTF["min"] = bandwidthRange.Min
					bandwithRangeTF["max"] = bandwidthRange.Max
					bandwithRange = append(bandwithRange, bandwithRangeTF)
				}
				connectorSettingTF["bandwidthrange"] = bandwithRange
				connectorSettings = append(connectorSettings, connectorSettingTF)
			}
			oi["connectorsettings"] = connectorSettings

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenSubnetItemsData(SubnetItems []client.InputSubnet) interface{} {
	if SubnetItems != nil {
		ois := make([]interface{}, len(SubnetItems), len(SubnetItems))

		for i, subnetItem := range SubnetItems {
			oi := make(map[string]interface{})
			oi["subnet"] = subnetItem.Subnet
			oi["virtual_subnet"] = subnetItem.VirtualSubnet
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenSecurityItemsData(SecurityItems *client.Security) interface{} {
	if SecurityItems != nil {
		ois := make([]map[string]interface{}, 0)
		oi := make(map[string]interface{})
		policyItems := flattenPolicyItemData(SecurityItems.Policies)
		oi["policies"] = policyItems

		ois = append(ois, oi)
		return ois
	}
	return make([]interface{}, 0)
}

func flattenPolicyItemData(PolicyItems []client.Policyops) []interface{} {
	if PolicyItems != nil {
		ois := make([]interface{}, len(PolicyItems), len(PolicyItems))

		for i, PolicyItem := range PolicyItems {
			oi := make(map[string]interface{})

			oi["id"] = PolicyItem.ID
			oi["name"] = PolicyItem.Name
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
