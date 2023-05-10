package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppOnboarding_IP() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to onboard IP Endpoint apps.",
		CreateContext: resourceAppOnboarding_IP_Create,
		UpdateContext: resourceAppOnboarding_IP_Update,
		DeleteContext: resourceAppOnboarding_IP_Delete,
		ReadContext:   resourceAppOnboarding_IP_Read,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
						"app_fqdn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fqdn of the app that user would access after onboarding ",
						},
						"protocols": {
							Type:        schema.TypeSet,
							MaxItems:    1,
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
										Optional:    true,
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
										Type:     schema.TypeList,
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
												// "backend_ip_address_discover": {
												// 	Type:     schema.TypeBool,
												// 	Optional: true,
												// },
											},
										},
									},
								},
							},
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
				Description: " Select policy name.e.g: ALLOW-ALL-USERS, DENY-ALL-USERS or CUSTOMIZE.Conditional access policies and Web Application Firewall policies for the application",
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
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceAppOnboarding_IP_Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

	appOnboardObjOpts, diags := getAppOnboardConfigObj_IP(d)
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

	resourceAppOnboarding_IP_Read(ctx, d, meta)

	return diags
}

func resourceAppOnboarding_IP_Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

	appOnboardObjOpts, diags, flag := getAppOnboardConfigObj_IP_Update(d)
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
	resourceAppOnboarding_IP_Read(ctx, d, meta)

	return diags
}

func resourceAppOnboarding_IP_Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

func resourceAppOnboarding_IP_Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func getAppOnboardConfigObj_IP(d *schema.ResourceData) (*client.AppOnboardSettingsOpts, diag.Diagnostics) {
	var diags diag.Diagnostics
	appURLsOptsConfig := d.Get("app_urls").(*schema.Set).List()
	appURLsOptsList := []*client.AppURLOpts{}
	// dnsCustomOpts := &client.AppOnboardDnsCustom{}

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
			// cloudConfigRegionOpts.ConnOption = ConnOptns
			if ConnOptns == client.Optnpeering {
				cloudConfigRegionOpts.ConnOption = client.OptnpeeringInput
			} else if ConnOptns == client.OptntransitGateway || ConnOptns == client.OptnvwanHub || ConnOptns == client.OptnazureTransitVnet {
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
			cloudConfigRegionOpts.BackendIPAddressDiscover = true
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
				Protocol: formatedrotocolsValues["protocol"].(string),
				// PortList:            expandStringList(formatedrotocolsValues["port_list"].([]interface{})),
				WebSocketEnabled:    formatedrotocolsValues["web_socket_enabled"].(bool),
				IsValidProtocolPort: formatedrotocolsValues["web_socket_enabled"].(bool),
			}
			if appProtocol.Protocol == client.DNSfqdnAppProtocol {
				appProtocol.PortList = []string{"53"}
			} else {
				appProtocol.PortList = expandStringList(formatedrotocolsValues["port_list"].([]interface{}))
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

		//App Url Config
		appURLOpts := &client.AppURLOpts{
			AppFqdn:         appURLOptsValues["app_fqdn"].(string) + "/32",
			Protocols:       protocols,
			CloudConfigOpts: appOnboardCloudConfigOpts,
		}
		appURLsOptsList = append(appURLsOptsList, appURLOpts)

	}

	//App Onboard Config
	appOnboardSettingsOpts := &client.AppOnboardSettingsOpts{
		App_Name:       d.Get("app_name").(string),
		AppOnboardType: client.TypeIP,
		AppURLsOpts:    appURLsOptsList,
		OptOption:      d.Get("optimization_option").(string),
		// AppSamlRewrite:   appSamlRewrite,
		EnableMultiCloud: d.Get("enable_multi_cloud_access").(bool),
		PolicyName:       expandStringList(d.Get("policy_name").([]interface{})),
		OnboardApp:       d.Get("onboard_app").(bool),
		DecommissionApp:  d.Get("decommission_app").(bool),
	}

	return appOnboardSettingsOpts, diags
}

func getAppOnboardConfigObj_IP_Update(d *schema.ResourceData) (*client.AppOnboardSettingsOpts, diag.Diagnostics, bool) {
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
		appOnboardObjOpts, diags := getAppOnboardConfigObj_IP(d)
		if diags != nil {
			return nil, diags, false
		}
		appOnboardSettingsOpts = appOnboardObjOpts
	}

	return appOnboardSettingsOpts, diags, updateReq
}

func getAppURL_IP(i int, appURL *client.AppURL, prosimoClient *client.ProsimoClient, ctx context.Context, d *schema.ResourceData) interface{} {

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
