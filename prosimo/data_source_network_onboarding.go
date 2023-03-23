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
			"filter_cloud_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based upon cloud type, e.g: AWS, AZURE, GCP",
			},
			"filter_cloud_region": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter based upon cloud region, e.g: europe-central2, us-east-2",
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
						"public_cloud": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cloud_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"connection_option": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_networks": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cloud_network_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"connector_group_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"edge_connectivity_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hub_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"connectivity_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"subnets": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"connector_placement": {
													Type:     schema.TypeString,
													Computed: true,
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
	count := 0

	onboardNetworkList, err := prosimoClient.SearchOnboardNetworks(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)
	cloudType := d.Get("filter_cloud_type").([]interface{})
	cloudRegions := d.Get("filter_cloud_region").([]interface{})

	if filter != "" {
		count += 1
		for _, onboardNetwork := range onboardNetworkList.Data.Records {
			filteredMap := map[string]interface{}{}

			err := mapstructure.Decode(onboardNetwork, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, filteredMap)
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
	}
	if len(cloudType) > 0 {
		count += 1
		flag := false
		CloudNameList := expandStringList(cloudType)
		getCloud, err := prosimoClient.GetCloudCreds(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, cloudName := range CloudNameList {
			var idList []string
			for _, cloudCred := range getCloud.CloudCreds {
				if cloudCred.CloudType == cloudName {
					idList = append(idList, cloudCred.ID)
				}
			}
			for _, cloudkey := range idList {
				for _, onboardNetwork := range onboardNetworkList.Data.Records {
					if onboardNetwork.PublicCloud.CloudKeyID == cloudkey {
						returnNetworkList = append(returnNetworkList, *onboardNetwork)
						flag = true
					}
					// }
				}
			}
		}
		if !flag {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})

			return diags
		}
	}
	if len(cloudRegions) > 0 {
		count += 1
		flag := false
		CloudRegionList := expandStringList(cloudRegions)
		for _, regionName := range CloudRegionList {
			for _, onboardNetwork := range onboardNetworkList.Data.Records {
				if onboardNetwork.PublicCloud.CloudRegion == regionName {
					returnNetworkList = append(returnNetworkList, *onboardNetwork)
					flag = true
				}
			}
		}
		if !flag {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})

			return diags
		}
	}

	if count == 0 {
		for _, onboardNetwork := range onboardNetworkList.Data.Records {
			returnNetworkList = append(returnNetworkList, *onboardNetwork)
		}
	}
	if count > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Input, more than one filter condition is not supported",
			Detail:   fmt.Sprintln("As of now Prosimo datasources support a single filtering entity."),
		})

		return diags
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
			publicCloudItems := flattenPublicCloudItemsData(NetworkItem.PublicCloud)
			oi["public_cloud"] = publicCloudItems
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

		oi["cloud_id"] = PublicCloudItems.Id
		oi["cloud"] = PublicCloudItems.Cloud
		oi["cloud_type"] = PublicCloudItems.CloudType
		oi["connection_option"] = PublicCloudItems.ConnectionOption
		oi["cloud_key_id"] = PublicCloudItems.CloudKeyID
		oi["cloud_region"] = PublicCloudItems.CloudRegion
		cloudNetworkItems := flattenCloudNetworksItemsData(PublicCloudItems.CloudNetworks)
		oi["cloud_networks"] = cloudNetworkItems

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
			oi["cloud_network_id"] = CloudNetworkItem.CloudNetworkID
			oi["connector_group_id"] = CloudNetworkItem.ConnectorGrpID
			oi["edge_connectivity_id"] = CloudNetworkItem.EdgeConnectivityID
			oi["hub_id"] = CloudNetworkItem.HubID
			oi["connectivity_type"] = CloudNetworkItem.ConnectivityType
			oi["subnets"] = CloudNetworkItem.Subnets

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
