package prosimo

import (
	"context"
	"fmt"
	"reflect"
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
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"connectorplacement": {
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
			oi["subnets"] = CloudNetworkItem.Subnets
			oi["connectorplacement"] = CloudNetworkItem.ConnectorPlacement

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
