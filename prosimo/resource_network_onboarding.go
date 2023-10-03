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

func resourceNetworkOnboarding() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to onboard networks.",
		CreateContext: resourceNetworkOnboardingCreate,
		ReadContext:   resourceNetworkOnboardingRead,
		DeleteContext: resourceNetworkOnboardingDelete,
		UpdateContext: resourceNetworkOnboardingUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name for the application",
			},
			"network_exportable_policy": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Mark Network Exportable in Policy",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Assigned Namespace",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pam_cname": {
				Type:     schema.TypeString,
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
			"private_cloud": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_creds_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cloud application account name.",
						},
						"subnets": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "subnet cider list",
						},
					},
				},
			},
			"public_cloud": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "public",
							ValidateFunc: validation.StringInSlice(client.GetCloudTypeOptions(), false),
							Description:  "public or private cloud",
						},
						"connection_option": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetCloudConnectionOptions(), false),
							Description:  "public or private cloud",
						},
						"cloud_creds_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cloud application account name.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of cloud region",
						},
						"cloud_networks": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPC ID",
									},
									"vnet": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VNET ID",
									},
									"connector_placement": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetConnectorPlacementOptions(), false),
										Description:  "Infra VPC, Workload VPC or none.",
									},
									"hub_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "(Required if transit-gateway is selected) tgw-id",
									},
									"connectivity_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "transit-gateway, vpc-peering",
									},
									"subnets": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "subnet cider list",
									},
									"service_insertion_endpoint_subnets": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(client.GetServiceInsertionOptions(), false),
										Description:  "Service Insertion Endpoint, applicable when connector is placed in Workload VPC",
									},
									"connector_settings": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bandwidth": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  " Available Options: <1 Gbps, 1-5 Gbps, 5-10 Gbps, >10 Gbps",
													ValidateFunc: validation.StringInSlice(client.GetConnectorBandwidthOptions(), false),
												},
												"instance_type": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(client.GetConnectorInstanceOptions(), false),
													Description: "Available Options wrt cloud and bandwidth :" +
														"Cloud_Provider: AWS:" +
														"Bandwidth:  <1 Gbps, Available Options: t3.medium/t3a.medium/c5.large" +
														"Bandwidth:  1-5 Gbps, Available Options: c5a.large/c5.xlarge/c5a.xlarge/c5n.xlarge" +
														"Bandwidth: 5-10 Gbps, Available Options: c5a.8xlarge/c5.9xlarge" +
														"Bandwidth: >10 Gbps, Available Options: c5n.9xlarge/c5a.16xlarge/c5.18xlarge/c5n.18xlarge" +
														"Cloud_Provider: AZURE:" +
														"For AZURE Default Connector settings are used,hence user does not have to specify is explicitly" +
														"Provided values: Bandwidth: <1 Gbps, Instance Type: Standard_A2_v2" +
														"Cloud_Provider: GCP:" +
														"Bandwidth:  <1 Gbps, Available Options: e2-standard-2" +
														"Bandwidth:  1-5 Gbps, Available Options: e2-standard-4" +
														"Bandwidth: 5-10 Gbps, Available Options: e2-standard-8/e2-standard-16" +
														"Bandwidth: >10 Gbps, Available Options: c2-standard-16",
												},
												"connector_subnets": {
													Type:        schema.TypeList,
													Optional:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "connector subnet cider list, Applicable when connector placement is in workload VPC/VNET ",
												},
											},
										},
									},
								},
							},
						},
						"connect_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "connector",
						},
					},
				},
			},
			"policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Select policy name.e.g: ALLOW-ALL-NETWORKS, DENY-ALL-NETWORKS or Custom Policies",
			},
			"onboard_app": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like the network to be onboarded to fabric",
			},
			"decommission_app": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like the network  to be offboarded from fabric",
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

func resourceNetworkOnboardingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	appOffboardFlag := d.Get("decommission_app").(bool)
	if appOffboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid  decommission_app flag.",
			Detail:   "decommission_app can't be set to true while creating network onboarding resource.",
		})
		return diags
	}

	// CloudCredName := d.Get("app_name").(string)
	nameSpace, err := prosimoClient.GetNamespaceByName(ctx, d.Get("namespace").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	networkOnboardops := &client.NetworkOnboardoptns{
		Name:        d.Get("name").(string),
		Exportable:  d.Get("network_exportable_policy").(bool),
		NamespaceID: nameSpace.ID,
	}

	onboardresponse, err := prosimoClient.NetworkOnboard(ctx, networkOnboardops)
	if err != nil {
		return diag.FromErr(err)
	}
	networkOnboardops.ID = onboardresponse.NetworkOnboardResponseData.ID

	diags = resourceNetworkOnboardingSettings(ctx, d, meta, networkOnboardops)
	if diags != nil {
		return diags
	}
	d.SetId(networkOnboardops.ID)
	if d.Get("onboard_app").(bool) {
		res, err2 := prosimoClient.OnboardNetworkDeployment(ctx, networkOnboardops.ID)
		if err2 != nil {
			return diag.FromErr(err2)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[DEBUG] Waiting for task id %s to complete", res.NetworkDeploymentResops.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskCompleteNetwork(ctx, d, meta, res.NetworkDeploymentResops.TaskID, networkOnboardops))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", res.NetworkDeploymentResops.TaskID)
		}
	}

	resourceNetworkOnboardingRead(ctx, d, meta)

	return diags
}

func resourceNetworkOnboardingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

	updateReq := false
	if d.HasChange("name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify App Name",
			Detail:   "App Name can't be modified",
		})
		return diags
	}
	if d.HasChange("network_exportable_policy") {
		updateReq = true
	}
	if d.HasChange("namespace") {
		updateReq = true
	}
	if d.HasChange("public_cloud") {
		updateReq = true
	}
	if d.HasChange("private_cloud") {
		updateReq = true
	}
	if d.HasChange("policies") {
		updateReq = true
	}
	if d.HasChange("onboard_app") && !d.IsNewResource() {
		updateReq = true
	}
	if d.HasChange("decommission_app") && !d.IsNewResource() {
		updateReq = true
	}

	if updateReq {
		networkOnboardops := &client.NetworkOnboardoptns{
			Name:       d.Get("name").(string),
			Exportable: d.Get("network_exportable_policy").(bool),
			ID:         d.Id(),
		}
		nameSpace, _ := prosimoClient.GetNamespaceByName(ctx, d.Get("namespace").(string))
		networkOnboardops.NamespaceID = nameSpace.ID
		diags, networkOnboardops = resourceNetworkOnboardingSettingsUpdate(ctx, d, meta, networkOnboardops)
		if diags != nil {
			return diags
		}
		if d.Get("decommission_app").(bool) {
			onboardresponse, err := prosimoClient.OffboardNetworkDeployment(ctx, networkOnboardops.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			if d.Get("wait_for_rollout").(bool) {
				log.Printf("[DEBUG] Waiting for task id %s to complete", onboardresponse.NetworkDeploymentResops.TaskID)
				err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
					retryUntilTaskComplete(ctx, d, meta, onboardresponse.NetworkDeploymentResops.TaskID))
				if err != nil {
					return diag.FromErr(err)
				}
				log.Printf("[DEBUG] task %s is successful", onboardresponse.NetworkDeploymentResops.TaskID)
			}
		} else if d.Get("onboard_app").(bool) {
			networkOnboardSettingsDbObj, err := prosimoClient.GetNetworkSettings(ctx, d.Id())
			if err != nil {
				return diag.FromErr(err)
			}
			if networkOnboardSettingsDbObj.Deployed {
				res, err := prosimoClient.OnboardNetworkDeploymentPost(ctx, networkOnboardops)
				if err != nil {
					return diag.FromErr(err)
				}

				if d.Get("wait_for_rollout").(bool) {
					log.Printf("[DEBUG] Waiting for task id %s to complete", res.NetworkDeploymentResops.TaskID)
					err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
						retryUntilTaskCompleteNetwork(ctx, d, meta, res.NetworkDeploymentResops.TaskID, networkOnboardops))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[INFO] task %s is successful", res.NetworkDeploymentResops.TaskID)
				}
			} else {
				res, err2 := prosimoClient.OnboardNetworkDeployment(ctx, networkOnboardops.ID)
				if err2 != nil {
					return diag.FromErr(err2)
				}
				if d.Get("wait_for_rollout").(bool) {
					log.Printf("[DEBUG] Waiting for task id %s to complete", res.NetworkDeploymentResops.TaskID)
					err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
						retryUntilTaskCompleteNetwork(ctx, d, meta, res.NetworkDeploymentResops.TaskID, networkOnboardops))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[INFO] task %s is successful", res.NetworkDeploymentResops.TaskID)
				}
			}
		}

	}
	resourceNetworkOnboardingRead(ctx, d, meta)

	return diags
}
func resourceNetworkOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	// log.Printf("resourceAppOnboardingRead %s", d.Id())
	networkOnboardSettingsDbObj, err := prosimoClient.GetNetworkSettings(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(networkOnboardSettingsDbObj.ID)
	d.Set("name", networkOnboardSettingsDbObj.Name)
	d.Set("network_exportable_policy", networkOnboardSettingsDbObj.Exportable)
	d.Set("namespace", d.Get("namespace").(string))
	d.Set("pam_cname", networkOnboardSettingsDbObj.PamCname)
	d.Set("deployed", networkOnboardSettingsDbObj.Deployed)
	d.Set("status", networkOnboardSettingsDbObj.Status)
	d.Set("policies", d.Get("policies").([]interface{}))
	d.Set("onboard_app", networkOnboardSettingsDbObj.Deployed)
	d.Set("decommission_app", d.Get("decommission_app").(bool))

	return diags

}

func resourceNetworkOnboardingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	networkOffBoardSettingsID := d.Id()

	appSummary, err := prosimoClient.GetNetworkSettings(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if appSummary.Status == "DEPLOYED" {
		log.Printf("[INFO] Network is in Onboarded State, Initiating Offboard ")
		appOffboardResData, err := prosimoClient.OffboardNetworkDeployment(ctx, networkOffBoardSettingsID)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[DEBUG] Waiting for task id %s to complete", appOffboardResData.NetworkDeploymentResops.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete(ctx, d, meta, appOffboardResData.NetworkDeploymentResops.TaskID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] task %s is successful", appOffboardResData.NetworkDeploymentResops.TaskID)
		}
	}
	del_err := prosimoClient.DeleteNetworkDeployment(ctx, networkOffBoardSettingsID)
	if del_err != nil {
		return diag.FromErr(del_err)
	}
	return diags
}

func resourceNetworkOnboardingSettings(ctx context.Context, d *schema.ResourceData, meta interface{}, networkOnboardops *client.NetworkOnboardoptns) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics

	// Public Cloud configuration.
	if v, ok := d.GetOk("public_cloud"); ok && v.(*schema.Set).Len() > 0 {
		publiccloudOptsConfig := v.(*schema.Set).List()[0].(map[string]interface{})
		cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, publiccloudOptsConfig["cloud_creds_name"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
		cloudNetworkInputList := []client.CloudNetworkops{}
		if v, ok := publiccloudOptsConfig["cloud_networks"]; ok && v.(*schema.Set).Len() > 0 {
			cloudNetworkListConfig := v.(*schema.Set).List()
			for _, cloudNetwork := range cloudNetworkListConfig {
				cloudNetworkConfig := cloudNetwork.(map[string]interface{})
				cloudNetworkInput := &client.CloudNetworkops{
					ConnectivityType: cloudNetworkConfig["connectivity_type"].(string),
					HubID:            cloudNetworkConfig["hub_id"].(string),
					Subnets:          expandStringList(cloudNetworkConfig["subnets"].([]interface{})),
				}
				connectorPlacement := cloudNetworkConfig["connector_placement"].(string)
				if connectorPlacement == client.WorkloadVpcConnectorPlacementOptions {
					cloudNetworkInput.ConnectorPlacement = client.AppConnectorPlacementOptions
				} else if connectorPlacement == client.InfraVPCConnectorPlacementOptions {
					cloudNetworkInput.ConnectorPlacement = client.InfraConnectorPlacementOptions
				} else {
					cloudNetworkInput.ConnectorPlacement = connectorPlacement
				}
				if cloudCreds.CloudType == client.AzureCloudType {
					cloudNetworkInput.CloudNetworkID = cloudNetworkConfig["vnet"].(string)
				} else {
					cloudNetworkInput.CloudNetworkID = cloudNetworkConfig["vpc"].(string)
				}
				if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions && cloudCreds.CloudType == client.AWSCloudType {
					if v, ok := cloudNetworkConfig["service_insertion_endpoint_subnets"].(string); ok {
						serviceSubnet := &client.ServiceSubnets{
							Mode: v,
						}
						cloudNetworkInput.Servicesubnets = serviceSubnet
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Missing Service Endpoint details",
							Detail:   "Service Endpoint details are required if connector placement is in Infra or workload vpc and cloud type is AWS.",
						})

						return diags
					}
				}
				switch cloudCreds.CloudType {
				case client.AWSCloudType:
					// if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions || cloudNetworkInput.ConnectorPlacement == client.InfraConnectorPlacementOptions {
					if v, ok := cloudNetworkConfig["connector_settings"]; ok && v.(*schema.Set).Len() > 0 {
						connectorsettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
						connectorsettingInput := &client.ConnectorSettings{
							// Bandwidth:     connectorsettingConfig["bandwidth"].(string),
							BandwidthName: connectorsettingConfig["bandwidth"].(string),
							InstanceType:  connectorsettingConfig["instance_type"].(string),
						}
						switch connectorsettingInput.BandwidthName {
						case client.LessThan1GBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeSmall
						case client.OneToFiveGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeMedium
						case client.FiveToTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeLarge
						case client.MoreThanTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeExtraLarge
						}

						if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions {
							if v, ok := connectorsettingConfig["connector_subnets"]; ok {
								if len(expandStringList(v.([]interface{}))) > 0 {
									connectorsettingInput.Subnets = expandStringList(v.([]interface{}))
								} else {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Missing Connector group Subnets",
										Detail:   "Connector group Subnets are required if connector placement is in Workload VPC.",
									})

									return diags
								}
							}
						}
						cloudNetworkInput.Connectorsettings = connectorsettingInput
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Missing Connector Active setting options",
							Detail:   "Active setting options are required if Cloud Type is AWS.",
						})

						return diags
					}
				case client.AzureCloudType:
					log.Println("entering Azure block")
					connectorsettingInput := &client.ConnectorSettings{
						Bandwidth:     client.AzureBandwidth,
						BandwidthName: client.AzureBandwidthName,
						InstanceType:  client.AzureInstanceType,
					}
					if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions {
						if v, ok := cloudNetworkConfig["connector_settings"]; ok && v.(*schema.Set).Len() > 0 {
							connectorsettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
							connectorsettingInput.Subnets = expandStringList(connectorsettingConfig["connector_subnets"].([]interface{}))
						} else {
							diags = append(diags, diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "Missing Connector group Subnets",
								Detail:   "Connector group Subnets are required if connector placement is in Workload VPC.",
							})

							return diags
						}
					}
					log.Println("connectorsettingInput", connectorsettingInput)
					cloudNetworkInput.Connectorsettings = connectorsettingInput

				case client.GCPCloudType:
					if v, ok := cloudNetworkConfig["connector_settings"]; ok && v.(*schema.Set).Len() > 0 {
						connectorsettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
						connectorsettingInput := &client.ConnectorSettings{
							// Bandwidth:     connectorsettingConfig["bandwidth"].(string),
							BandwidthName: connectorsettingConfig["bandwidth"].(string),
							InstanceType:  connectorsettingConfig["instance_type"].(string),
							Subnets:       expandStringList(connectorsettingConfig["connector_subnets"].([]interface{})),
						}
						switch connectorsettingInput.BandwidthName {
						case client.LessThan1GBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeSmall
						case client.OneToFiveGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeMedium
						case client.FiveToTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeLarge
						case client.MoreThanTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeExtraLarge
						}
						cloudNetworkInput.Connectorsettings = connectorsettingInput
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Missing Connector Active setting options",
							Detail:   "Active setting options are required if Cloud Type is GCP.",
						})

						return diags
					}

				}
				cloudNetworkInputList = append(cloudNetworkInputList, *cloudNetworkInput)
			}
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing Cloud Network settings.",
				Detail:   "Please provide cloud network details like vpc, vnet etc ",
			})
			return diags
		}
		publicCloudoptn := &client.PublicCloud{
			CloudType:        publiccloudOptsConfig["cloud_type"].(string),
			ConnectionOption: publiccloudOptsConfig["connection_option"].(string),
			CloudKeyID:       cloudCreds.ID,
			CloudRegion:      publiccloudOptsConfig["region_name"].(string),
			CloudNetworks:    cloudNetworkInputList,
			ConnectType:      publiccloudOptsConfig["connect_type"].(string),
		}
		networkOnboardops.PublicCloud = publicCloudoptn

	}

	// Private Cloud configuration.
	if v, ok := d.GetOk("private_cloud"); ok && v.(*schema.Set).Len() > 0 {
		privatecloudOptsConfig := v.(*schema.Set).List()[0].(map[string]interface{})

		privateCloudoptn := &client.PrivateCloud{
			CloudType:        "private",
			ConnectionOption: "private",
			Subnets:          expandStringList(privatecloudOptsConfig["subnets"].([]interface{})),
		}
		cloudCredName := privatecloudOptsConfig["cloud_creds_name"].(string)
		cloudCreds, err := prosimoClient.GetCloudCredsPrivate(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, cloudCred := range cloudCreds.CloudCreds {
			if cloudCred.Nickname == cloudCredName {
				privateCloudoptn.PrivateCloudID = cloudCred.ID
			}
		}
		networkOnboardops.PrivateCloud = privateCloudoptn
	}
	err := prosimoClient.NetworkOnboardCloud(ctx, networkOnboardops)
	if err != nil {
		return diag.FromErr(err)
	}
	// Securirty policy configuration.
	networkPolicyList := []client.Policyops{}
	if v, ok := d.GetOk("policies"); ok {
		inputPolicies := expandStringList(v.([]interface{}))
		for _, inputpolicy := range inputPolicies {
			networkPolicy := client.Policyops{}
			policyDbObj, err := prosimoClient.GetPolicyByName(ctx, inputpolicy)
			if err != nil {
				return diag.FromErr(err)
			}
			networkPolicy.ID = policyDbObj.ID
			networkPolicyList = append(networkPolicyList, networkPolicy)
		}
		policyList := &client.Security{
			Policies: networkPolicyList,
		}
		securityInput := &client.NetworkSecurityInput{
			Security: policyList,
		}
		err := prosimoClient.NetworkOnboardSecurity(ctx, securityInput, networkOnboardops.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceNetworkOnboardingSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}, networkOnboardops *client.NetworkOnboardoptns) (diag.Diagnostics, *client.NetworkOnboardoptns) {
	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics

	// Public Cloud configuration.
	if v, ok := d.GetOk("public_cloud"); ok {
		publiccloudOptsConfig := v.(*schema.Set).List()[0].(map[string]interface{})
		cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, publiccloudOptsConfig["cloud_creds_name"].(string))
		if err != nil {
			return diag.FromErr(err), nil
		}
		cloudNetworkInputList := []client.CloudNetworkops{}
		if v, ok := publiccloudOptsConfig["cloud_networks"]; ok && v.(*schema.Set).Len() > 0 {
			cloudNetworkListConfig := v.(*schema.Set).List()
			for _, cloudNetwork := range cloudNetworkListConfig {
				cloudNetworkConfig := cloudNetwork.(map[string]interface{})
				cloudNetworkInput := &client.CloudNetworkops{
					ConnectivityType:   cloudNetworkConfig["connectivity_type"].(string),
					HubID:              cloudNetworkConfig["hub_id"].(string),
					ConnectorPlacement: cloudNetworkConfig["connector_placement"].(string),
					Subnets:            expandStringList(cloudNetworkConfig["subnets"].([]interface{})),
				}
				connectorPlacement := cloudNetworkConfig["connector_placement"].(string)
				if connectorPlacement == client.WorkloadVpcConnectorPlacementOptions {
					cloudNetworkInput.ConnectorPlacement = client.AppConnectorPlacementOptions
				} else if connectorPlacement == client.InfraVPCConnectorPlacementOptions {
					cloudNetworkInput.ConnectorPlacement = client.InfraConnectorPlacementOptions
				} else {
					cloudNetworkInput.ConnectorPlacement = connectorPlacement
				}
				if cloudCreds.CloudType == client.AzureCloudType {
					cloudNetworkInput.CloudNetworkID = cloudNetworkConfig["vnet"].(string)
				} else {
					cloudNetworkInput.CloudNetworkID = cloudNetworkConfig["vpc"].(string)
				}
				if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions && cloudCreds.CloudType == client.AWSCloudType {
					if v, ok := cloudNetworkConfig["service_insertion_endpoint_subnets"].(string); ok {
						serviceSubnet := &client.ServiceSubnets{
							Mode: v,
						}
						cloudNetworkInput.Servicesubnets = serviceSubnet
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Missing Service Endpoint details",
							Detail:   "Service Endpoint details are required if connector placement is in Infra or workload vpc and cloud type is AWS.",
						})

						return diags, nil
					}
				}
				switch cloudCreds.CloudType {
				case client.AWSCloudType:
					// if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions || cloudNetworkInput.ConnectorPlacement == client.InfraConnectorPlacementOptions {
					if v, ok := cloudNetworkConfig["connector_settings"]; ok && v.(*schema.Set).Len() > 0 {
						connectorsettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
						connectorsettingInput := &client.ConnectorSettings{
							// Bandwidth:     connectorsettingConfig["bandwidth"].(string),
							BandwidthName: connectorsettingConfig["bandwidth"].(string),
							InstanceType:  connectorsettingConfig["instance_type"].(string),
						}
						switch connectorsettingInput.BandwidthName {
						case client.LessThan1GBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeSmall
						case client.OneToFiveGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeMedium
						case client.FiveToTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeLarge
						case client.MoreThanTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeExtraLarge
						}

						if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions {
							if v, ok := connectorsettingConfig["connector_subnets"]; ok {
								if len(expandStringList(v.([]interface{}))) > 0 {
									connectorsettingInput.Subnets = expandStringList(v.([]interface{}))
								} else {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Missing Connector group Subnets",
										Detail:   "Connector group Subnets are required if connector placement is in Workload VPC.",
									})

									return diags, nil
								}
							}
						}
						cloudNetworkInput.Connectorsettings = connectorsettingInput
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Missing Connector Active setting options",
							Detail:   "Active setting options are required if Cloud Type is AWS.",
						})

						return diags, nil
					}
				case client.AzureCloudType:
					connectorsettingInput := &client.ConnectorSettings{
						Bandwidth:     client.AzureBandwidth,
						BandwidthName: client.AzureBandwidthName,
						InstanceType:  client.AzureInstanceType,
					}
					if cloudNetworkInput.ConnectorPlacement == client.AppConnectorPlacementOptions {
						if v, ok := cloudNetworkConfig["connector_settings"]; ok && v.(*schema.Set).Len() > 0 {
							connectorsettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
							connectorsettingInput.Subnets = expandStringList(connectorsettingConfig["connector_subnets"].([]interface{}))
						} else {
							diags = append(diags, diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "Missing Connector group Subnets",
								Detail:   "Connector group Subnets are required if connector placement is in Workload VPC.",
							})

							return diags, nil
						}
					}
					cloudNetworkInput.Connectorsettings = connectorsettingInput

				case client.GCPCloudType:
					if v, ok := cloudNetworkConfig["connector_settings"]; ok && v.(*schema.Set).Len() > 0 {
						connectorsettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
						connectorsettingInput := &client.ConnectorSettings{
							// Bandwidth:     connectorsettingConfig["bandwidth"].(string),
							BandwidthName: connectorsettingConfig["bandwidth"].(string),
							InstanceType:  connectorsettingConfig["instance_type"].(string),
							Subnets:       expandStringList(connectorsettingConfig["connector_subnets"].([]interface{})),
						}
						switch connectorsettingInput.BandwidthName {
						case client.LessThan1GBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeSmall
						case client.OneToFiveGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeMedium
						case client.FiveToTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeLarge
						case client.MoreThanTenGBPS:
							connectorsettingInput.Bandwidth = client.ConnectorSizeExtraLarge
						}
						cloudNetworkInput.Connectorsettings = connectorsettingInput
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Missing Connector Active setting options",
							Detail:   "Active setting options are required if Cloud Type is GCP.",
						})

						return diags, nil
					}
				}
				cloudNetworkInputList = append(cloudNetworkInputList, *cloudNetworkInput)
			}
		}

		publicCloudoptn := &client.PublicCloud{
			CloudType:        publiccloudOptsConfig["cloud_type"].(string),
			ConnectionOption: publiccloudOptsConfig["connection_option"].(string),
			CloudKeyID:       cloudCreds.ID,
			CloudRegion:      publiccloudOptsConfig["region_name"].(string),
			CloudNetworks:    cloudNetworkInputList,
			ConnectType:      publiccloudOptsConfig["connect_type"].(string),
		}
		networkOnboardops.PublicCloud = publicCloudoptn

	}

	// Private Cloud configuration.
	if v, ok := d.GetOk("private_cloud"); ok && v.(*schema.Set).Len() > 0 {
		privatecloudOptsConfig := v.(*schema.Set).List()[0].(map[string]interface{})

		privateCloudoptn := &client.PrivateCloud{
			CloudType:        "private",
			ConnectionOption: "private",
			Subnets:          expandStringList(privatecloudOptsConfig["subnets"].([]interface{})),
		}
		cloudCredName := privatecloudOptsConfig["cloud_creds_name"].(string)
		cloudCreds, err := prosimoClient.GetCloudCredsPrivate(ctx)
		if err != nil {
			return diag.FromErr(err), nil
		}
		for _, cloudCred := range cloudCreds.CloudCreds {
			if cloudCred.Nickname == cloudCredName {
				privateCloudoptn.PrivateCloudID = cloudCred.ID
			}
		}
		networkOnboardops.PrivateCloud = privateCloudoptn
	}

	networkOnboardSettingsDbObj, err := prosimoClient.GetNetworkSettings(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err), nil
	}
	if !networkOnboardSettingsDbObj.Deployed {
		err := prosimoClient.NetworkOnboardCloud(ctx, networkOnboardops)
		if err != nil {
			return diag.FromErr(err), nil
		}
	}

	// Securirty policy configuration.
	networkPolicyList := []client.Policyops{}
	if v, ok := d.GetOk("policies"); ok {
		inputPolicies := expandStringList(v.([]interface{}))
		for _, inputpolicy := range inputPolicies {
			networkPolicy := client.Policyops{}
			policyDbObj, err := prosimoClient.GetPolicyByName(ctx, inputpolicy)
			if err != nil {
				return diag.FromErr(err), nil
			}
			networkPolicy.ID = policyDbObj.ID
			networkPolicyList = append(networkPolicyList, networkPolicy)
		}
	}
	policyList := &client.Security{
		Policies: networkPolicyList,
	}

	networkOnboardops.Security = policyList
	securityInput := &client.NetworkSecurityInput{
		Security: policyList,
	}

	if !networkOnboardSettingsDbObj.Deployed {
		err := prosimoClient.NetworkOnboardSecurity(ctx, securityInput, networkOnboardops.ID)
		if err != nil {
			return diag.FromErr(err), nil
		}
	}

	return diags, networkOnboardops
}
