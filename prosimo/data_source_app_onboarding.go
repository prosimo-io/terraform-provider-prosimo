package prosimo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourceAppOnboarding() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on onboarded applications.",
		ReadContext: dataSourceAppOnboardingRead,
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
			"filter_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter based upon protocol used, e.g: http, https, tcp",
			},
			"filter_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter based on port number, eg: 80, 22 etc",
			},
			"filter_internal_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter based on internal domain name, eg: abc.com, 10.5.0.4/32 etc",
			},
			"filter_papp_fqdn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter based upon papp fqdn",
			},
			"filter_app_fqdn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter based upon onboarded app fqdn",
			},
			"app_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of onboarded apps",
			},
			"onboarded_apps": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"optimize_app_experience": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"optoption": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enablemulticloud": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"saml_rewrite": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apptype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"onboardtype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interactiontype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"addresstype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_urls": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"teamid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"internaldomain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Actual Fqdn of the app",
									},
									"domaintype": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of Domain: e.g Custom or Prosimo domian",
									},
									"appfqdn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Fqdn of the app that used would access after onboarding ",
									},
									"pappfqdn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Fqdn of the app that used would access after onboarding ",
									},
									"subdomainincluded": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Set True if Subdomians need to be included else False",
									},
									"cloud_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cacheruleid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocols": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Protocol that prosimo edge uses to connect to App",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ptotocol type, e.g: tcp, udp, http etc",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "target port number",
												},
												"portlist": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"websocketenabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Set to true if tou would like prosimo edges to communicate with app via websocket",
												},
												"isvalidprotocolport": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"paths": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Customized websocket paths",
												},
											},
										},
									},
									"ext_protocols": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"health_check_info": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"endpoint": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HealthCheck Endpoints",
												},
											},
										},
									},
									"connectionoption": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of connection: e.g- Private,Public.",
									},
									"edge_regions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the region where app is available",
												},
												"connoption": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Connection option for private connection: e.g Peering, Transitgateway",
												},
												"regiontype": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of region: e.g: Active, Backup etc",
												},
												"endpoints": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"appnetworkid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"appip": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"dnsservice": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"configured": {
													Type:     schema.TypeBool,
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

func dataSourceAppOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnAPPList []client.AppOnboardSettings
	count := 0
	// flag := false

	apppnboardSearchops := &client.AppOnboardSearch{}
	apppnboardSearchops.Category = "app"
	onboardAppList, err := prosimoClient.SearchAppOnboardApps(ctx, apppnboardSearchops)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)
	cloudType := d.Get("filter_cloud_type").([]interface{})
	cloudRegions := d.Get("filter_cloud_region").([]interface{})
	protocol := d.Get("filter_protocol").(string)
	port := d.Get("filter_port").(int)
	internalDomain := d.Get("filter_internal_domain").(string)
	pappFqdn := d.Get("filter_papp_fqdn").(string)
	appFqdn := d.Get("filter_app_fqdn").(string)

	if filter != "" {
		count += 1
		for _, onboardApp := range onboardAppList.Data.Records {
			filteredMap := map[string]interface{}{}

			err := mapstructure.Decode(onboardApp, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, filteredMap)
			if diags != nil {
				return diags
			}
			if flag {
				returnAPPList = append(returnAPPList, *onboardApp)
			}
		}
		if len(returnAPPList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})

			return diags
		}
	}

	if pappFqdn != "" {
		count += 1
		flag := false
		for _, onboardApp := range onboardAppList.Data.Records {
			for _, appUrl := range onboardApp.AppURLs {
				if appUrl.PappFqdn == pappFqdn {
					returnAPPList = append(returnAPPList, *onboardApp)
					flag = true
					break
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
	if appFqdn != "" {
		count += 1
		flag := false
		for _, onboardApp := range onboardAppList.Data.Records {
			for _, appUrl := range onboardApp.AppURLs {
				if appUrl.AppFqdn == appFqdn {
					returnAPPList = append(returnAPPList, *onboardApp)
					flag = true
					break
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
	if internalDomain != "" {
		count += 1
		flag := false
		for _, onboardApp := range onboardAppList.Data.Records {
			for _, appUrl := range onboardApp.AppURLs {
				if appUrl.InternalDomain == internalDomain {
					returnAPPList = append(returnAPPList, *onboardApp)
					flag = true
					break
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
	if protocol != "" {
		count += 1
		flag := false
		for _, onboardApp := range onboardAppList.Data.Records {
			for _, appUrl := range onboardApp.AppURLs {
				for _, protocoL := range appUrl.Protocols {
					if protocoL.Protocol == protocol {
						returnAPPList = append(returnAPPList, *onboardApp)
						flag = true
					}
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
	if port != 0 {
		count += 1
		flag := false
		for _, onboardApp := range onboardAppList.Data.Records {
			for _, appUrl := range onboardApp.AppURLs {
				for _, protocoL := range appUrl.Protocols {
					if protocoL.Port == port {
						returnAPPList = append(returnAPPList, *onboardApp)
						flag = true
					} else if len(protocoL.PortList) > 0 {
						for _, porT := range protocoL.PortList {
							i, _ := strconv.Atoi(porT)
							if i == port {
								returnAPPList = append(returnAPPList, *onboardApp)
								flag = true
							}
						}
					}
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
				for _, onboardApp := range onboardAppList.Data.Records {
					for _, appurl := range onboardApp.AppURLs {
						if appurl.CloudKeyID == cloudkey {
							returnAPPList = append(returnAPPList, *onboardApp)
							flag = true
						}
					}
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
			for _, onboardApp := range onboardAppList.Data.Records {
				for _, appUrl := range onboardApp.AppURLs {
					for _, region := range appUrl.Regions {
						if region.Name == regionName {
							returnAPPList = append(returnAPPList, *onboardApp)
							flag = true
						}
					}
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
	// }
	if count == 0 {
		for _, onboardapp := range onboardAppList.Data.Records {
			returnAPPList = append(returnAPPList, *onboardapp)
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
	appItems := flattenAppItemsData(returnAPPList)
	d.Set("onboarded_apps", appItems)
	d.Set("app_count", len(returnAPPList))
	return diags
}

func flattenAppItemsData(AppItems []client.AppOnboardSettings) []interface{} {
	if AppItems != nil {
		ois := make([]interface{}, len(AppItems), len(AppItems))

		for i, AppItem := range AppItems {
			oi := make(map[string]interface{})

			oi["app_name"] = AppItem.App_Name
			oi["id"] = AppItem.ID
			oi["team_id"] = AppItem.Team_ID
			oi["idp_id"] = AppItem.IDP_ID
			oi["app_access_type"] = AppItem.App_Access_Type
			oi["policy_group_id"] = AppItem.PolicyGroupID
			oi["optimize_app_experience"] = AppItem.Optimize_App_Experience
			oi["optoption"] = AppItem.OptOption
			oi["enablemulticloud"] = AppItem.EnableMultiCloud
			oi["status"] = AppItem.Status
			oi["apptype"] = AppItem.AppType
			oi["onboardtype"] = AppItem.OnboardType
			oi["interactiontype"] = AppItem.InterActionType
			oi["addresstype"] = AppItem.AddressType
			appUrlItems := flattenAppUrlItemsData(AppItem.AppURLs)
			oi["app_urls"] = appUrlItems

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenAppUrlItemsData(AppUrlItems []*client.AppURL) []interface{} {
	if AppUrlItems != nil {
		ois := make([]interface{}, len(AppUrlItems), len(AppUrlItems))

		for i, AppUrlItem := range AppUrlItems {
			oi := make(map[string]interface{})

			oi["id"] = AppUrlItem.ID
			oi["teamid"] = AppUrlItem.TeamID
			oi["internaldomain"] = AppUrlItem.InternalDomain
			oi["domaintype"] = AppUrlItem.DomainType
			oi["appfqdn"] = AppUrlItem.AppFqdn
			oi["pappfqdn"] = AppUrlItem.PappFqdn
			oi["subdomainincluded"] = AppUrlItem.SubdomainIncluded
			oi["cloud_key_id"] = AppUrlItem.CloudKeyID
			oi["certid"] = AppUrlItem.CertID
			oi["cacheruleid"] = AppUrlItem.CacheRuleID
			protocolItems := flattenProtocolItemsData(AppUrlItem.Protocols)
			oi["protocols"] = protocolItems
			extprotocolItems := flattenextProtocolItemsData(AppUrlItem.ExtProtocols)
			oi["ext_protocols"] = extprotocolItems
			healthCheckInfo := make([]map[string]interface{}, 0)
			appHealthCheckInfo := AppUrlItem.HealthCheckInfo
			healthCheckInfoTF := make(map[string]interface{})
			healthCheckInfoTF["enabled"] = appHealthCheckInfo.Enabled
			healthCheckInfoTF["endpoint"] = appHealthCheckInfo.Endpoint
			healthCheckInfo = append(healthCheckInfo, healthCheckInfoTF)
			oi["health_check_info"] = healthCheckInfo
			oi["connectionoption"] = AppUrlItem.ConnectionOption
			cloudConfigItems := flattenCloudConfiglItemsData(AppUrlItem.Regions)
			oi["edge_regions"] = cloudConfigItems
			dnsInfo := make([]map[string]interface{}, 0)
			appDnsInfo := AppUrlItem.DNSService
			dnsInfoTF := make(map[string]interface{})
			dnsInfoTF["type"] = appDnsInfo.Type
			dnsInfoTF["id"] = appDnsInfo.ID
			dnsInfo = append(dnsInfo, dnsInfoTF)
			oi["dnsservice"] = dnsInfo

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
func flattenProtocolItemsData(ProtocolItems []*client.AppProtocol) []interface{} {
	if ProtocolItems != nil {
		ois := make([]interface{}, len(ProtocolItems), len(ProtocolItems))

		for i, ProtocolItem := range ProtocolItems {
			oi := make(map[string]interface{})

			oi["protocol"] = ProtocolItem.Protocol
			oi["port"] = ProtocolItem.Port
			oi["portlist"] = ProtocolItem.PortList
			oi["websocketenabled"] = ProtocolItem.WebSocketEnabled
			oi["isvalidprotocolport"] = ProtocolItem.IsValidProtocolPort
			oi["paths"] = ProtocolItem.Paths

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenextProtocolItemsData(extProtocolItems []*client.AppProtocol) []interface{} {
	if extProtocolItems != nil {
		ois := make([]interface{}, len(extProtocolItems), len(extProtocolItems))

		for i, extProtocolItem := range extProtocolItems {
			oi := make(map[string]interface{})

			oi["protocol"] = extProtocolItem.Protocol
			oi["port"] = extProtocolItem.Port
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
func flattenCloudConfiglItemsData(cloudConfigItems []*client.AppOnboardCloudConfigRegions) []interface{} {
	if cloudConfigItems != nil {
		ois := make([]interface{}, len(cloudConfigItems), len(cloudConfigItems))

		for i, cloudConfigItem := range cloudConfigItems {
			oi := make(map[string]interface{})

			oi["id"] = cloudConfigItem.ID
			oi["name"] = cloudConfigItem.Name
			oi["connoption"] = cloudConfigItem.ConnOption
			oi["regiontype"] = cloudConfigItem.RegionType
			cloudRegionEPItems := flattenCloudRegionEPItemsData(cloudConfigItem.Endpoints)
			oi["endpoints"] = cloudRegionEPItems
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
func flattenCloudRegionEPItemsData(EndpointItems []*client.AppOnboardCloudRegionEndpoints) []interface{} {
	if EndpointItems != nil {
		ois := make([]interface{}, len(EndpointItems), len(EndpointItems))

		for i, EndpointItem := range EndpointItems {
			oi := make(map[string]interface{})

			oi["appnetworkid"] = EndpointItem.AppNetworkID
			oi["appip"] = EndpointItem.AppIP
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
