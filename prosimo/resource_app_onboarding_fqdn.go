package prosimo

import (
	"context"
	"log"
	"strings"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppOnboarding_FQDN() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to onboard TCP/UDP apps.",
		CreateContext: resourceAppOnboarding_FQDN_Create,
		UpdateContext: resourceAppOnboarding_FQDN_Update,
		DeleteContext: resourceAppOnboarding_FQDN_Delete,
		ReadContext:   resourceAppOnboarding_FQDN_Read,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"idp_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IDP provider name.",
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name for the application",
			},
			"app_access_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "e.g: Agent or Agentless",
			},
			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "e.g: type of app onboarded, e.g: citrix, web, fqdn, jumpbox",
			},
			"app_urls": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetAppDomainTypes(), false),
							Description:  "Type of Domain: e.g custom or prosimo",
						},
						"app_fqdn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fqdn of the app that user would access after onboarding ",
						},
						"subdomain_included": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set True to onboard subdomains of the application else False",
						},
						"protocols": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Protocol that prosimo edge uses to connect to App",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetAppProtocolsLFQDN(), false),
										Description:  "Protocol type, e.g: “http”, “https”, “ssh”, “vnc”, or “rdp",
									},
									"port_list": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "target port number",
									},
									"web_socket_enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Set to true if tou would like prosimo edges to communicate with app via websocket",
									},
									"is_valid_protocol_port": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"paths": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Customized websocket paths",
									},
								},
							},
						},
						"health_check_info": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Description: "Application health check config from edge",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"endpoint": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "HealthCheck Endpoints",
									},
								},
							},
						},
						"cloud_config": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_hosted_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "PUBLIC",
										ValidateFunc: validation.StringInSlice(client.AppHostedOptn(), false),
										Description:  "Wheather app is hosted in Public cloud like AWS/AZURE/GCP or private DC. Available options PRIVATE/PUBLIC",
									},
									"connection_option": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetCloudConnectionOptions(), false),
										Description:  "Public, if the app domain has a public IP address / DNS A record on the internet currently, and the Prosimo Edge should connect to the application using a public connection.Private, if the application only has a private IP address, and Edge should connect to it over a private connection.",
									},
									"cloud_creds_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "cloud account under which application is hosted",
									},
									"dc_app_ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Applicable only if  app_hosted_type is PRIVATE, IP of the app hosted in PRIVATE DC",
									},
									"is_show_connection_options": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"has_private_connection_options": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"edge_regions": {
										Type: schema.TypeList,
										// MinItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the region where app is available",
												},
												"conn_option": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(client.AppOnboardConnOptn(), false),
													Description:  "Connection option for private connection: e.g: peering/transitGateway/awsPrivateLink/azurePrivateLink/azureTransitVnet/vwanHub",
												},
												"region_type": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(client.AppOnboardRegionType(), false),
													Description:  "Type of region: e.g:active, backup etc",
												},
												"tgw_app_routetable": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(client.GetTgwAppRoutetableType(), false),
												},
												"app_network_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "App network id details",
												},
												"attach_point_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Attach Point id details",
												},
												"backend_ip_address_discover": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "if Set to true, auto discovers available endpoints",
												},
												"backend_ip_address_manual": {
													Type:        schema.TypeList,
													Optional:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Pass endpoints manually.",
												},
												"backend_ip_address_dns": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"dns_custom": {
													Type:        schema.TypeSet,
													MaxItems:    1,
													Optional:    true,
													Description: "Custom DNS setup",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"dns_app": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "DNS App name",
															},
															"dns_server": {
																Type:        schema.TypeList,
																Optional:    true,
																Elem:        &schema.Schema{Type: schema.TypeString},
																Description: "DNS Server List",
															},
															"is_healthcheck_enabled": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Health check to ensure application domains being resolved by dns servers ",
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
						"cache_rule": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cache Rules for your App Domains",
						},
						"waf_policy_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "WAF Policies for your App Domains, applicable when the Edge to App Protocol is either HTTP or HTTPS.",
						},
					},
				},
			},

			"saml_rewrite": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "App authentication option while selecting prosimo domain",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Required while selecting SAML based authentication",
						},
						"metadata_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Required while selecting SAML based authentication",
						},
						"selected_auth_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(client.GetSelectedAuthTypes(), false),
							Description:  "Type of authentication: e.g. SAML, OIDC, Others",
						},
					},
				},
			},
			"optimization_option": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(client.GetAppOnboardingOptimization(), false),
				Description:  "Optimization option for app: e.g: CostSaving, PerformanceEnhanced, FastLane",
			},
			"enable_multi_cloud_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Setting this to true would leverage multi clouds to optimize the app performance ",
			},
			"policy_name": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Select policy name.e.g: ALLOW-ALL-USERS, DENY-ALL-USERS or CUSTOMIZE.Conditional access policies and Web Application Firewall policies for the application.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"customize_policy": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "Choose any custom policy created from the policy library or create one.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"onboard_app": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like app to be onboarded to fabric",
			},
			"decommission_app": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like app to be offboarded from fabric",
			},
			"wait_for_rollout": {
				Type:        schema.TypeBool,
				Description: "Wait for the rollout of the task to complete. Defaults to true.",
				Default:     true,
				Optional:    true,
			},
			"force_offboard": {
				Type:        schema.TypeBool,
				Description: "Force app offboarding incase of normal offboarding failure.",
				Default:     true,
				Optional:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceAppOnboarding_FQDN_Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	appOffboardFlag := d.Get("decommission_app").(bool)
	if appOffboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid  decommission_app flag.",
			Detail:   "decommission_app can't be set to true while creating apponboarding resource.",
		})
		return diags
	}

	//Validate IDP
	idpName := d.Get("idp_name").(string)
	if idpName != "" {
		diags = validate_primaryIDP(ctx, idpName, meta)
		if diags != nil {
			return diags
		}
	}

	appOnboardObjOpts, diags := getAppOnboardConfigObj_FQDN(d)
	if diags != nil {
		return diags
	}
	// Step 1: create settings config
	diags = createAppOnboardSettings(ctx, d, meta, appOnboardObjOpts)
	if diags != nil {
		return diags
	}

	appOnboardSettingsDbObj, err := prosimoClient.GetAppOnboardSettings(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	appOnboardObjOpts.ID = appOnboardSettingsDbObj.ID
	for _, appURLDB := range appOnboardSettingsDbObj.AppURLs {

		for _, appURLOpts := range appOnboardObjOpts.AppURLsOpts {

			if appURLOpts.InternalDomain == appURLDB.InternalDomain {
				appURLOpts.ID = appURLDB.ID
			}

		}

	}

	// Step 2: create cloud config
	diags = createAppOnboardCloudConfigs(ctx, d, meta, appOnboardObjOpts)
	if diags != nil {
		return diags
	}

	// Step 3: create optimization option
	diags = createAppOnboardOptOption(ctx, d, meta, appOnboardObjOpts)
	if diags != nil {
		return diags
	}

	// Step 4: create waf and policy
	diags = createAppOnboardSecurity(ctx, d, meta, appOnboardObjOpts)
	if diags != nil {
		return diags
	}

	// do summary endpoint before app onboarding
	_, err = prosimoClient.GetAppOnboardSummary(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// resourceAppOnboardingRead(ctx, d, meta)

	// Step 5: onboard app
	diags = onboardApp(ctx, d, meta, appOnboardObjOpts)
	if diags != nil {
		return diags
	}

	resourceAppOnboarding_FQDN_Read(ctx, d, meta)

	return diags
}

func resourceAppOnboarding_FQDN_Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	appOnboardFlag := d.Get("onboard_app").(bool)
	appOffboardFlag := d.Get("decommission_app").(bool)
	if appOnboardFlag && appOffboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid onboard_app and decommission_app flag combination.",
			Detail:   "Both onboard_app and decommission_app have been set to true.",
		})
		return diags
	}

	appOnboardObjOpts, diags, flag := getAppOnboardConfigObj_FQDN_Update(d)
	if diags != nil {
		return diags
	}
	if flag {
		appOnboardSettingsDbObj, err := prosimoClient.GetAppOnboardSettings(ctx, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		appOnboardObjOpts.ID = appOnboardSettingsDbObj.ID
		for _, appURLDB := range appOnboardSettingsDbObj.AppURLs {
			for _, appURLOpts := range appOnboardObjOpts.AppURLsOpts {
				if appURLOpts.InternalDomain == appURLDB.InternalDomain {
					appURLOpts.ID = appURLDB.ID
					for _, appURLDBRegion := range appURLDB.Regions {
						for _, appURLOptsRegion := range appURLOpts.CloudConfigOpts.Regions {
							if appURLDBRegion.Name == appURLOptsRegion.Name {
								appURLOptsRegion.ID = appURLDBRegion.ID
							}
						}
					}
				}
			}
		}

		//Offboard app
		offBoardApp := false
		if d.HasChange("decommission_app") && !d.IsNewResource() {
			isDecommission := d.Get("decommission_app").(bool)
			if isDecommission {
				offBoardApp = true
				diags = offboardApp(ctx, d, meta, appOnboardObjOpts)
				if diags != nil {
					return diags
				}
			}
		}

		if !offBoardApp {
			// Step 1: create settings config
			diags = updateAppOnboardSettings(ctx, d, meta, appOnboardObjOpts, d.Id())
			if diags != nil {
				return diags
			}

			// updating again if any new apps are added
			appOnboardSettingsDbObj, err := prosimoClient.GetAppOnboardSettings(ctx, d.Id())
			if err != nil {
				return diag.FromErr(err)
			}
			appOnboardObjOpts.ID = appOnboardSettingsDbObj.ID
			for _, appURLDB := range appOnboardSettingsDbObj.AppURLs {

				for _, appURLOpts := range appOnboardObjOpts.AppURLsOpts {

					if appURLOpts.InternalDomain == appURLDB.InternalDomain {
						appURLOpts.ID = appURLDB.ID
					}

				}

			}

			// Step 2: create cloud config (Validate if it's a reboarding, If so clould config call would be skipped.)
			appSummary, err := prosimoClient.GetAppOnboardSummary(ctx, d.Id())
			if err != nil {
				return diag.FromErr(err)
			}
			if !appSummary.Deployed {
				diags = createAppOnboardCloudConfigs(ctx, d, meta, appOnboardObjOpts)
				if diags != nil {
					return diags
				}
			} else {
				log.Println("[DEBUG] Skipping Cloud config changes.Can't modify cloud config for a deployed app.")
			}

			// Step 3: create optimization option
			diags = createAppOnboardOptOption(ctx, d, meta, appOnboardObjOpts)
			if diags != nil {
				return diags
			}

			// Step 4: create waf and policy
			diags = createAppOnboardSecurity(ctx, d, meta, appOnboardObjOpts)
			if diags != nil {
				return diags
			}
		}

		//App reboard
		isOnboard := d.Get("onboard_app").(bool)
		if isOnboard {

			// do summary endpoint before app onboarding
			_, err = prosimoClient.GetAppOnboardSummary(ctx, d.Id())
			if err != nil {
				return diag.FromErr(err)
			}

			// Step 5: onboard app
			diags = onboardApp(ctx, d, meta, appOnboardObjOpts)
			if diags != nil {
				return diags
			}
		}
	}
	resourceAppOnboarding_FQDN_Read(ctx, d, meta)

	return diags
}

func resourceAppOnboarding_FQDN_Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	appOnboardSettingsDbObj, err := prosimoClient.GetAppOnboardSettings(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(appOnboardSettingsDbObj.ID)
	d.Set("app_name", appOnboardSettingsDbObj.App_Name)
	d.Set("app_access_type", appOnboardSettingsDbObj.App_Access_Type)
	d.Set("app_type", appOnboardSettingsDbObj.AppType)
	d.Set("optimization_option", appOnboardSettingsDbObj.OptOption)
	d.Set("enable_multi_cloud_access", appOnboardSettingsDbObj.EnableMultiCloud)
	d.Set("onboard_app", appOnboardSettingsDbObj.Deployed)
	d.Set("decommission_app", d.Get("decommission_app").(bool))

	return diags

}

func resourceAppOnboarding_FQDN_Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOffBoardSettingsID := d.Id()

	appSummary, err := prosimoClient.GetAppOnboardSummary(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if appSummary.Status == "DEPLOYED" {
		log.Printf("[INFO] App is in Onboarded State, Initiating Offboard ")
		appOffboardResData, err := prosimoClient.OffboardAppDeployment(ctx, appOffBoardSettingsID)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[DEBUG] Waiting for task id %s to complete", appOffboardResData.ResourceData.ID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete(ctx, d, meta, appOffboardResData.ResourceData.ID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] task %s is successful", appOffboardResData.ResourceData.ID)
		}
	}
	del_err := prosimoClient.DeleteApp(ctx, appOffBoardSettingsID)
	if del_err != nil {
		return diag.FromErr(del_err)
	}
	return diags

}

func getAppOnboardConfigObj_FQDN(d *schema.ResourceData) (*client.AppOnboardSettingsOpts, diag.Diagnostics) {
	var diags diag.Diagnostics
	appURLsOptsConfig := d.Get("app_urls").(*schema.Set).List()
	appURLsOptsList := []*client.AppURLOpts{}
	dnsCustomOpts := &client.AppOnboardDnsCustom{}
	var ipAddressDns bool

	for _, appURLOptsList := range appURLsOptsConfig {
		appURLOptsValues := appURLOptsList.(map[string]interface{})

		//Edge Region Config
		cloudConfig := appURLOptsValues["cloud_config"].(*schema.Set).List()[0].(map[string]interface{})
		edgeRegionsConfig := cloudConfig["edge_regions"].([]interface{})
		regionOptsList := []*client.AppOnboardCloudConfigRegionOpts{}
		for _, edgeRegion := range edgeRegionsConfig {
			edgeRegionValues := edgeRegion.(map[string]interface{})

			cloudConfigRegionOpts := &client.AppOnboardCloudConfigRegionOpts{
				Name: edgeRegionValues["region_name"].(string),
				// ConnOption: edgeRegionValues["conn_option"].(string),
				RegionType: edgeRegionValues["region_type"].(string),
			}
			ConnOptns := edgeRegionValues["conn_option"].(string)
			cloudConfigRegionOpts.ConnOption = ConnOptns
			if ConnOptns == client.Optnpeering {
				cloudConfigRegionOpts.ConnOption = client.OptnpeeringInput
			} else if ConnOptns == client.OptntransitGateway || ConnOptns == client.OptnvwanHub {
				cloudConfigRegionOpts.ConnOption = ConnOptns
				if ConnOptns == client.OptntransitGateway {
					cloudConfigRegionOpts.ModifyTgwAppRouteTable = edgeRegionValues["tgw_app_routetable"].(string)
				}
				cloudConfigRegionOpts.AppnetworkID = edgeRegionValues["app_network_id"].(string)
				cloudConfigRegionOpts.AttachPointID = edgeRegionValues["attach_point_id"].(string)
			} else if ConnOptns == client.OptnawsPrivateLink {
				cloudConfigRegionOpts.ConnOption = ConnOptns
				cloudConfigRegionOpts.AppnetworkID = edgeRegionValues["app_network_id"].(string)
			} else {
				cloudConfigRegionOpts.ConnOption = ConnOptns
			}

			ipAddressDiscover := edgeRegionValues["backend_ip_address_discover"].(bool)
			ipAddressDns = edgeRegionValues["backend_ip_address_dns"].(bool)
			if v, ok := edgeRegionValues["dns_custom"].(*schema.Set); ok && v.Len() > 0 {
				dnsCustom := v.List()[0].(map[string]interface{})
				// dnsCustomOpts := &client.AppOnboardDnsCustom{}
				dnsCustomOpts.IsHealthCheckEnabled = dnsCustom["is_healthcheck_enabled"].(bool)
				if v, ok := dnsCustom["dns_app"]; ok {
					dnsApp := v.(string)
					dnsCustomOpts.DnsAppName = dnsApp
				}
				if v, ok := dnsCustom["dns_server"]; ok {
					dnsServer := v.([]interface{})
					if len(dnsServer) > 0 {
						dnsCustomOpts.DnsServers = expandStringList(dnsServer)
					}
				}
			}
			if ipAddressDiscover {
				cloudConfigRegionOpts.BackendIPAddressDiscover = ipAddressDiscover
			} else {
				if v, ok := edgeRegionValues["backend_ip_address_manual"]; ok {
					ipAddressList := v.([]interface{})
					if len(ipAddressList) > 0 {
						cloudConfigRegionOpts.BackendIPAddressEntry = expandStringList(v.([]interface{}))
					}
				}
			}
			regionOptsList = append(regionOptsList, cloudConfigRegionOpts)
		}

		//Cloud Config
		appOnboardCloudConfigOpts := &client.AppOnboardCloudConfigOpts{
			AppHOstedType:              cloudConfig["app_hosted_type"].(string),
			ConnectionOption:           cloudConfig["connection_option"].(string),
			CloudCredsName:             cloudConfig["cloud_creds_name"].(string),
			IsShowConnectionOptions:    cloudConfig["is_show_connection_options"].(bool),
			HasPrivateConnectionOption: cloudConfig["has_private_connection_options"].(bool),
			Regions:                    regionOptsList,
		}
		// Read DC app IP if APP is hosted in private DC
		if appOnboardCloudConfigOpts.AppHOstedType == client.HostedPrivate {
			appOnboardCloudConfigOpts.DCAappIP = cloudConfig["dc_app_ip"].(string)
		}

		// app protocols
		protocols := []*client.AppProtocol{}
		inprotocolsValues := appURLOptsValues["protocols"].(*schema.Set)
		for _, protocolsValues := range inprotocolsValues.List() {
			formatedrotocolsValues := protocolsValues.(map[string]interface{})

			appProtocol := &client.AppProtocol{
				Protocol:            formatedrotocolsValues["protocol"].(string),
				PortList:            expandStringList(formatedrotocolsValues["port_list"].([]interface{})),
				WebSocketEnabled:    formatedrotocolsValues["web_socket_enabled"].(bool),
				IsValidProtocolPort: formatedrotocolsValues["web_socket_enabled"].(bool),
			}
			if appProtocol.WebSocketEnabled {
				if v, ok := formatedrotocolsValues["paths"]; ok {
					pathList := v.([]interface{})
					if len(pathList) > 0 {
						appProtocol.Paths = expandStringList(v.([]interface{}))
					}
				}
			}
			protocols = append(protocols, appProtocol)
		}

		// app healthcheck
		healthCheckInfo := appURLOptsValues["health_check_info"].(*schema.Set)
		enabled := healthCheckInfo.List()[0].(map[string]interface{})["enabled"].(bool)
		endpoint := healthCheckInfo.List()[0].(map[string]interface{})["endpoint"].(string)
		if !strings.HasPrefix(endpoint, "/") {
			endpoint = "/" + endpoint
		}

		appHealthCheckInfo := &client.AppHealthCheckInfo{
			Enabled:  enabled,
			Endpoint: endpoint,
		}

		//App Url Config
		appURLOpts := &client.AppURLOpts{
			// InternalDomain:    appURLOptsValues["internal_domain"].(string),
			DomainType:        appURLOptsValues["domain_type"].(string),
			AppFqdn:           appURLOptsValues["app_fqdn"].(string),
			SubdomainIncluded: appURLOptsValues["subdomain_included"].(bool),
			Protocols:         protocols,
			HealthCheckInfo:   appHealthCheckInfo,
			CloudConfigOpts:   appOnboardCloudConfigOpts,
			// DnsCustom:         dnsCustomOpts,
			CacheRuleName: appURLOptsValues["cache_rule"].(string),
			WafPolicyName: appURLOptsValues["waf_policy_name"].(string),
		}
		if appURLOpts.SubdomainIncluded || ipAddressDns {
			appURLOpts.DnsCustom = dnsCustomOpts
		}
		appURLsOptsList = append(appURLsOptsList, appURLOpts)

	}

	appSamlRewrite := &client.AppSamlRewrite{}
	if v, ok := d.Get("saml_rewrite").(*schema.Set); ok && v.Len() > 0 {
		samlCOnfig := v.List()[0].(map[string]interface{})
		if v, ok := samlCOnfig["selected_auth_type"]; ok {
			selectedAuthType := v.(string)
			appSamlRewrite.SelectedAuthType = selectedAuthType
			if selectedAuthType == client.SamlAuth {
				if v, ok := samlCOnfig["metadata"]; ok {
					appSamlRewrite.Metadata = v.(string)
				} else {
					appSamlRewrite.MetadataURL = d.Get("metadata_url").(string)
				}
			}
		}
	}

	//App Onboard Config
	appOnboardSettingsOpts := &client.AppOnboardSettingsOpts{
		App_Name: d.Get("app_name").(string),
		// Dns_Discovery:     d.Get("dns_discovery").(bool),
		AppOnboardType:   client.TypeFQDN,
		AppURLsOpts:      appURLsOptsList,
		OptOption:        d.Get("optimization_option").(string),
		AppSamlRewrite:   appSamlRewrite,
		EnableMultiCloud: d.Get("enable_multi_cloud_access").(bool),
		PolicyName:       expandStringList(d.Get("policy_name").([]interface{})),
		Dns_Discovery:    ipAddressDns,
		OnboardApp:       d.Get("onboard_app").(bool),
		DecommissionApp:  d.Get("decommission_app").(bool),
	}
	return appOnboardSettingsOpts, diags
}

func getAppOnboardConfigObj_FQDN_Update(d *schema.ResourceData) (*client.AppOnboardSettingsOpts, diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	appOnboardSettingsOpts := &client.AppOnboardSettingsOpts{}
	updateReq := false
	if d.HasChange("app_name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify App Name",
			Detail:   "App Name can't be modified",
		})
		return nil, diags, updateReq
	}
	if v, ok := d.Get("saml_rewrite").(*schema.Set); ok && v.Len() > 0 {
		if d.HasChange("saml_rewrite") {
			updateReq = true
		}
	}
	if d.HasChange("enable_multi_cloud_access") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("onboard_app") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("decommission_app") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("optimization_option") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("policy_name") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("app_urls") {
		updateReq = true
	}

	if updateReq {
		appOnboardObjOpts, diags := getAppOnboardConfigObj_FQDN(d)
		if diags != nil {
			return nil, diags, false
		}
		appOnboardSettingsOpts = appOnboardObjOpts
	}

	return appOnboardSettingsOpts, diags, updateReq
}

func getAppURL_FQDN(i int, appURL *client.AppURL, prosimoClient *client.ProsimoClient, ctx context.Context, d *schema.ResourceData) interface{} {

	var appURLTF = make(map[string]interface{})

	appURLTF["id"] = appURL.ID
	appURLTF["internal_domain"] = appURL.InternalDomain
	appURLTF["domain_type"] = appURL.DomainType
	appURLTF["app_fqdn"] = appURL.AppFqdn
	appURLTF["subdomain_included"] = appURL.SubdomainIncluded

	protocols := make([]map[string]interface{}, 0)
	for _, appProtocol := range appURL.Protocols {
		protocolTF := make(map[string]interface{})
		protocolTF["protocol"] = appProtocol.Protocol
		protocolTF["port_list"] = flattenStringList(appProtocol.PortList)
		if appProtocol.WebSocketEnabled {
			protocolTF["web_socket_enabled"] = appProtocol.WebSocketEnabled
			protocolTF["paths"] = flattenStringList(appProtocol.Paths)
		}
		protocols = append(protocols, protocolTF)

	}
	appURLTF["protocols"] = protocols

	healthCheckInfo := make([]map[string]interface{}, 0)
	appHealthCheckInfo := appURL.HealthCheckInfo
	healthCheckInfoTF := make(map[string]interface{})
	healthCheckInfoTF["enabled"] = appHealthCheckInfo.Enabled
	healthCheckInfoTF["endpoint"] = appHealthCheckInfo.Endpoint
	healthCheckInfo = append(healthCheckInfo, healthCheckInfoTF)
	appURLTF["health_check_info"] = healthCheckInfo

	cloudConfig := make([]map[string]interface{}, 0)
	cloudConfigTF := make(map[string]interface{})
	cloudConfigTF["connection_option"] = appURL.ConnectionOption

	edgeRegions := make([]map[string]interface{}, 0)
	for _, appRegions := range appURL.Regions {
		edgeRegionsTF := make(map[string]interface{})
		edgeRegionsTF["region_name"] = appRegions.Name
		edgeRegionsTF["region_type"] = appRegions.RegionType
		if appRegions.InputType != "entry" {
			edgeRegionsTF["backend_ip_address_discover"] = true
		} else {
			edgeRegionsTF["backend_ip_address_discover"] = false
			ips := make([]string, 0)
			for _, endpoints := range appRegions.Endpoints {
				ips = append(ips, endpoints.AppIP)
			}
			edgeRegionsTF["backend_ip_address_manual"] = ips
		}

		edgeRegions = append(edgeRegions, edgeRegionsTF)
	}
	cloudConfigTF["edge_regions"] = edgeRegions

	// get cloud name for cloud key id
	if appURL.CloudKeyID != "" {
		// log.Println("appURL.CloudKeyID", appURL.CloudKeyID)
		cloudCreds, err := prosimoClient.GetCloudCredsById(ctx, appURL.CloudKeyID)
		if err != nil {
			return diag.FromErr(err)
		}
		cloudConfigTF["cloud_creds_name"] = cloudCreds.Nickname
		cloudConfig = append(cloudConfig, cloudConfigTF)
		appURLTF["cloud_config"] = cloudConfig
	}

	if appURL.CacheRuleID != "" {
		cacheRule, err := prosimoClient.GetCacheRuleByID(ctx, appURL.CacheRuleID)
		if err != nil {
			return diag.FromErr(err)
		}
		appURLTF["cache_rule"] = cacheRule.Name
	}

	if appURL.WafHTTP != "" {
		wafRule, err := prosimoClient.GetWafByID(ctx, appURL.WafHTTP)
		if err != nil {
			return diag.FromErr(err)
		}
		appURLTF["waf_policy_name"] = wafRule.Name
	}

	return appURLTF
}
