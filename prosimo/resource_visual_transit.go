package prosimo

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVisualTransit() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify transit resources.",
		CreateContext: resourceVisualTransitUpdate,
		ReadContext:   resourceVisualTransitRead,
		DeleteContext: resourceVisualTransitDelete,
		UpdateContext: resourceVisualTransitUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"transit_input": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Transit setup input",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID, computed post resource creation",
						},
						"cloud_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud Type, e.g: AWS, AZURE, GCP",
						},
						"cloud_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud Region, e.g: us-east-2, westus ",
						},
						"transit_deployment": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Transit Deployment Config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tgws": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "TWG details",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of TGW",
												},
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "TGW ID",
												},
												"action": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(client.GetTransitTgwActionTypes(), false),
													Description:  "Action on TGW, e.g: ADD, MOD, DEL",
												},
												"account": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "AWS account Details: Applicable while creating a new TGW  ",
												},
												"connection": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "TGW connection Details",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringInSlice(client.GetTransitTgwConnectionTypes(), false),
																Description:  "Type of connection, e.g: EDGE, VPC",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of VPC if connection type is VPC",
															},
															"action": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringInSlice(client.GetTransitVpcActionTypes(), false),
																Description:  "Connection Action, e.g: ADD, DEL",
															},
														},
													},
												},
											},
										},
									},
									"vpcs": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "VPC Details",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Name of VPC",
												},
												"action": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(client.GetTransitVpcActionTypes(), false),
													Description:  "Action on VPC, e.g: ADD, DEL",
												},
											},
										},
									},
									"vhubs": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "VHUB details",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Name of VHUB",
												},
												"action": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(client.GetTransitTgwActionTypes(), false),
													Description:  "Action on VHUB, e.g: ADD, MOD, DEL. ADD action would create a new vHUB ",
												},
												"account": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Azure account Details: Applicable while creating a new vHUB  ",
												},
												"vwan": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "vWAN Details, Applicable while creating a new vHUB",
												},
												"address_space": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Address space of vHUB, Applicable while creating a new vHUB",
												},
												"connection": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "VHUB connection Details",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringInSlice(client.GetTransitVhubConnectionTypes(), false),
																Description:  "Type of connection, e.g: EDGE, VNET",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of VNET if connection type is VNET",
															},
															"action": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringInSlice(client.GetTransitVpcActionTypes(), false),
																Description:  "Connection Action, e.g: ADD, DEL",
															},
														},
													},
												},
											},
										},
									},
									"vnets": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "VNET Details",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Name of VNET",
												},
												"action": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(client.GetTransitVpcActionTypes(), false),
													Description:  "Action on VNET, e.g: ADD, DEL",
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
			"deploy_transit_setup": {
				Type:        schema.TypeBool,
				Description: "Flag to deploy transit setup, if set to true the setup is deployed else would be in config state, defaults to true.",
				Default:     true,
				Optional:    true,
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

func resourceVisualTransitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	inputVisualTransitSetup := []client.Visual_Transit_Setup{}
	edgeInput := client.EdgeInput{}
	edgeInputList := []client.EdgeInput{}
	transitSearchInput := client.TransitSearchInput{}
	if v, ok := d.GetOk("transit_input"); ok {
		for i, _ := range v.([]interface{}) {
			transitInputConfig := v.([]interface{})[i].(map[string]interface{})
			edgeInput = client.EdgeInput{
				CloudType:   transitInputConfig["cloud_type"].(string),
				CloudRegion: transitInputConfig["cloud_region"].(string),
			}
			edgeInputList = append(edgeInputList, edgeInput)
			// 	}
			// }
			transitSearchInput.Edges = edgeInputList
			transitops, err := prosimoClient.TransitSetupSearch(ctx, transitSearchInput)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, inops := range transitops {
				for _, inEdge := range edgeInputList {
					if inops.CloudType == inEdge.CloudType && inops.CloudRegion == inEdge.CloudRegion {
						deploymentInput := &client.Deployment{}
						if v, ok := transitInputConfig["transit_deployment"].(*schema.Set); ok && v.Len() > 0 {
							deploymentConfig := v.List()[0].(map[string]interface{})

							switch edgeInput.CloudType { //Switch Case for Cloud Type
							case client.AWSCloudType: // Cloud Type AWS

								if v, ok := deploymentConfig["vhubs"].(*schema.Set); ok && v.Len() > 0 {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Invalid Input",
										Detail:   fmt.Sprintln("vHUB is not a deployment option for cloudType AWS"),
									})

									return diags
								}

								if _, ok := deploymentConfig["vnets"].(*schema.Set); ok && v.Len() > 0 {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Invalid Input",
										Detail:   fmt.Sprintln("vNET is not a deployment option for cloudType AWS"),
									})

									return diags
								}
								if v, ok := deploymentConfig["tgws"]; ok {
									tgwinputList := []client.Constructs{}
									for i, _ := range v.([]interface{}) {
										tgwconfig := v.([]interface{})[i].(map[string]interface{})
										tgwID := tgwconfig["id"].(string)
										tgwname := tgwconfig["name"].(string)
										tgwAction := tgwconfig["action"].(string)
										connectionInputList := []client.Connection{}
										if v, ok := tgwconfig["connection"]; ok {
											for i, _ := range v.([]interface{}) {
												connectionInput := v.([]interface{})[i].(map[string]interface{})
												if connectionInput["type"].(string) == "EDGE" {
													connectionInputConfig := client.Connection{
														Name:      connectionInput["name"].(string),
														Type:      connectionInput["type"].(string),
														Action:    connectionInput["action"].(string),
														ID:        inops.Edge.ID,
														AccountID: inops.Edge.AccountID,
													}
													connectionInputList = append(connectionInputList, connectionInputConfig)
												} else {
													netresList, err := prosimoClient.CloudNetworkSearch(ctx, edgeInput)
													if err != nil {
														return diag.FromErr(err)
													}
													for _, netres := range netresList {
														for _, vpc := range netres.Account.VPCS {
															if vpc.Name == connectionInput["name"].(string) {
																connectionInputConfig := client.Connection{
																	Name:      connectionInput["name"].(string),
																	Type:      connectionInput["type"].(string),
																	Action:    connectionInput["action"].(string),
																	ID:        vpc.ID,
																	AccountID: netres.Account.AccountID,
																	Subnets:   vpc.Subnets,
																}
																connectionInputList = append(connectionInputList, connectionInputConfig)
																break
															}
														}
													}
												}

											}
										}
										if tgwAction == client.AddTgwAction {
											tgwAccount := tgwconfig["account"].(string)
											cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, tgwAccount)
											if err != nil {
												return diag.FromErr(err)
											}
											if cloudCreds == nil {
												diags = append(diags, diag.Diagnostic{
													Severity: diag.Error,
													Summary:  "Unable to get Cloud Credentials",
													Detail:   fmt.Sprintf("Unable to find Cloud Credentials for Account %s", tgwAccount),
												})

												return diags
											}
											tgw := &client.Constructs{
												ID:          tgwID,
												Action:      tgwAction,
												Name:        tgwname,
												AccountID:   cloudCreds.AccountID,
												Connections: connectionInputList,
											}
											tgwinputList = append(tgwinputList, *tgw)
										} else {
											flag := false
											for _, tgw := range inops.Operation.TGWS {
												if tgw.ID == tgwID {
													tgw.Action = tgwAction
													tgw.Connections = connectionInputList
													tgwinputList = append(tgwinputList, *tgw)
													flag = true
													break
												}
											}
											if !flag {
												diags = append(diags, diag.Diagnostic{
													Severity: diag.Error,
													Summary:  "Invalid TGW details, pls validate",
													Detail:   "Invalid TGW details, pls validate",
												})
												return diags
											}
										}
									}
									deploymentInput.TGWS = tgwinputList
								}
								if v, ok := deploymentConfig["vpcs"]; ok {
									vpcinputList := []client.Constructs{}
									for i, _ := range v.([]interface{}) {
										vpcconfig := v.([]interface{})[i].(map[string]interface{})
										vpcname := vpcconfig["name"].(string)
										vpcAction := vpcconfig["action"].(string)
										netresList, err := prosimoClient.CloudNetworkSearch(ctx, edgeInput)
										if err != nil {
											return diag.FromErr(err)
										}
										for _, netres := range netresList {
											for _, vpc := range netres.Account.VPCS {
												if vpc.Name == vpcname {
													invpc := client.Constructs{
														Name:      vpcname,
														Action:    vpcAction,
														ID:        vpc.ID,
														AccountID: netres.Account.AccountID,
														Subnets:   vpc.Subnets,
													}
													vpcinputList = append(vpcinputList, invpc)
													break
												}
											}
										}
									}
									deploymentInput.VPCS = vpcinputList
								}

							case client.AzureCloudType: //Cloud Type AZURE
								if _, ok := deploymentConfig["tgws"].(*schema.Set); ok && v.Len() > 0 {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Invalid Input",
										Detail:   fmt.Sprintln("TGW is not a deployment option for cloudType AZURE"),
									})

									return diags
								}

								if _, ok := deploymentConfig["vpcs"].(*schema.Set); ok && v.Len() > 0 {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Invalid Input",
										Detail:   fmt.Sprintln("VPC is not a deployment option for cloudType AZURE"),
									})

									return diags
								}

								if v, ok := deploymentConfig["vhubs"]; ok {
									vhubinputList := []client.Constructs{}
									for i, _ := range v.([]interface{}) {
										vhubconfig := v.([]interface{})[i].(map[string]interface{})
										connectionInputList := []client.Connection{}
										if v, ok := vhubconfig["connection"]; ok {
											for i, _ := range v.([]interface{}) {
												connectionInput := v.([]interface{})[i].(map[string]interface{})
												if connectionInput["type"].(string) == "EDGE" {
													connectionInputConfig := client.Connection{
														Name:      connectionInput["name"].(string),
														Type:      connectionInput["type"].(string),
														Action:    connectionInput["action"].(string),
														ID:        inops.Edge.ID,
														AccountID: inops.Edge.AccountID,
													}
													connectionInputList = append(connectionInputList, connectionInputConfig)
												} else {
													netresList, err := prosimoClient.CloudNetworkSearch(ctx, edgeInput)
													if err != nil {
														return diag.FromErr(err)
													}
													for _, netres := range netresList {
														for _, vnet := range netres.Account.VNETS {
															if vnet.Name == connectionInput["name"].(string) {
																connectionInputConfig := client.Connection{
																	Name:      connectionInput["name"].(string),
																	Type:      connectionInput["type"].(string),
																	Action:    connectionInput["action"].(string),
																	ID:        vnet.ID,
																	AccountID: netres.Account.AccountID,
																	Subnets:   vnet.Subnets,
																}
																connectionInputList = append(connectionInputList, connectionInputConfig)
																break
															}
														}
													}
												}

											}
										}
										vhubname := vhubconfig["name"].(string)
										vhubAction := vhubconfig["action"].(string)
										if vhubAction == client.AddTgwAction {
											vhubAccount := vhubconfig["account"].(string)
											cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, vhubAccount)
											if err != nil {
												return diag.FromErr(err)
											}
											if cloudCreds == nil {
												diags = append(diags, diag.Diagnostic{
													Severity: diag.Error,
													Summary:  "Unable to get Cloud Credentials",
													Detail:   fmt.Sprintf("Unable to find Cloud Credentials for Account %s", vhubAccount),
												})

												return diags
											}
											vhub := &client.Constructs{
												Action:       vhubAction,
												Name:         vhubconfig["name"].(string),
												AccountID:    cloudCreds.AccountID,
												Connections:  connectionInputList,
												VwanID:       vhubconfig["vwan"].(string),
												AddressSpace: vhubconfig["address_space"].(string),
											}
											vhubinputList = append(vhubinputList, *vhub)
										} else {
											flag := false
											for _, vhub := range inops.Operation.VHUBS {
												if vhub.Name == vhubname {
													vhub.Action = vhubAction
													vhub.Connections = connectionInputList
													vhubinputList = append(vhubinputList, *vhub)
													flag = true
													break
												}

											}
											if !flag {
												diags = append(diags, diag.Diagnostic{
													Severity: diag.Error,
													Summary:  "Invalid VHUB details, pls validate",
													Detail:   "Invalid VHUB details, pls validate",
												})
												return diags
											}
										}
									}
									deploymentInput.VHUBS = vhubinputList
								}
								if v, ok := deploymentConfig["vnets"]; ok {
									vnetinputList := []client.Constructs{}
									for i, _ := range v.([]interface{}) {
										vnetconfig := v.([]interface{})[i].(map[string]interface{})
										vnetname := vnetconfig["name"].(string)
										vnetAction := vnetconfig["action"].(string)
										netresList, err := prosimoClient.CloudNetworkSearch(ctx, edgeInput)
										if err != nil {
											return diag.FromErr(err)
										}
										for _, netres := range netresList {
											for _, vnet := range netres.Account.VNETS {
												if vnet.Name == vnetname {
													invnet := client.Constructs{
														Name:      vnetname,
														Action:    vnetAction,
														ID:        vnet.ID,
														AccountID: netres.Account.AccountID,
														Subnets:   vnet.Subnets,
													}
													vnetinputList = append(vnetinputList, invnet)
													break
												}
											}
										}
									}
									deploymentInput.VNETS = vnetinputList
								}
							}

						}
						inputTransit := client.Visual_Transit_Setup{
							CloudType:   inops.CloudType,
							CloudRegion: inops.CloudRegion,
							TeamID:      inops.TeamID,
							Edge:        inops.Edge,
							Operation:   inops.Operation,
							Deployment:  deploymentInput,
						}
						inputVisualTransitSetup = append(inputVisualTransitSetup, inputTransit)
					}
				}
			}
		}
	}

	errres := prosimoClient.TransitSetup(ctx, inputVisualTransitSetup)
	if errres != nil {
		return diag.FromErr(errres)
	}
	transitopsres, err := prosimoClient.TransitSetupSearch(ctx, transitSearchInput)
	if err != nil {
		return diag.FromErr(err)
	}
	setupIDList := []string{}
	for _, inops := range transitopsres {
		for _, inEdge := range edgeInputList {
			if inops.CloudType == inEdge.CloudType && inops.CloudRegion == inEdge.CloudRegion {
				setupIDList = append(setupIDList, inops.ID)
			}
		}
	}
	if d.Get("deploy_transit_setup").(bool) {
		deployRes, err_res_deploy := prosimoClient.CreateTransitDeploy(ctx, setupIDList)
		if err_res_deploy != nil {
			return diag.FromErr(err_res_deploy)
		}
		if d.Get("wait_for_rollout").(bool) {
			for _, res := range deployRes {
				// if res.CloudType == edgeInput.CloudType && res.CloudRegion == edgeInput.CloudRegion {
				log.Printf("[INFO] Waiting for task id %s to complete", res.TaskID)
				err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
					retryUntilTaskComplete(ctx, d, meta, res.TaskID))
				if err != nil {
					return diag.FromErr(err)
				}
				log.Printf("[INFO] task %s is successful", res.TaskID)
				// }
			}
		}
	}

	return resourceVisualTransitRead(ctx, d, meta)
}

func resourceVisualTransitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	edgeInputList := []client.EdgeInput{}
	transitSearchInput := client.TransitSearchInput{}
	inputVisualTransitSetup := []client.Visual_Transit_Setup{}

	var diags diag.Diagnostics

	if v, ok := d.GetOk("transit_input"); ok {
		for i, _ := range v.([]interface{}) {
			transitInputConfig := v.([]interface{})[i].(map[string]interface{})
			edgeInput := client.EdgeInput{
				CloudType:   transitInputConfig["cloud_type"].(string),
				CloudRegion: transitInputConfig["cloud_region"].(string),
			}
			edgeInputList = append(edgeInputList, edgeInput)
		}
	}
	transitSearchInput.Edges = edgeInputList
	transitops, err := prosimoClient.TransitSetupSummary(ctx, transitSearchInput)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, inops := range transitops {
		for _, inEdge := range edgeInputList {
			if inops.CloudType == inEdge.CloudType && inops.CloudRegion == inEdge.CloudRegion {
				inputVisualTransitSetup = append(inputVisualTransitSetup, *inops)
				break
			}
		}
	}
	d.SetId("transit setup")
	flatteninputVisualTransitSetup := flattenTransitInputData(inputVisualTransitSetup)
	d.Set("transit_input", flatteninputVisualTransitSetup)
	return diags
}

func flattenTransitInputData(inputVisualTransitSetup []client.Visual_Transit_Setup) []interface{} {
	if inputVisualTransitSetup != nil {
		ois := make([]interface{}, len(inputVisualTransitSetup), len(inputVisualTransitSetup))

		for i, inputVisualTransitSetupItem := range inputVisualTransitSetup {
			oi := make(map[interface{}]interface{})
			oi["id"] = inputVisualTransitSetupItem.ID
			oi["cloud_region"] = inputVisualTransitSetupItem.CloudRegion
			oi["cloud_type"] = inputVisualTransitSetupItem.CloudType
			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}

func resourceVisualTransitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	edgeInputList := []client.EdgeInput{}
	transitSearchInput := client.TransitSearchInput{}

	var diags diag.Diagnostics

	if v, ok := d.GetOk("transit_input"); ok {
		for i, _ := range v.([]interface{}) {
			transitInputConfig := v.([]interface{})[i].(map[string]interface{})
			edgeInput := client.EdgeInput{
				CloudType:   transitInputConfig["cloud_type"].(string),
				CloudRegion: transitInputConfig["cloud_region"].(string),
			}
			edgeInputList = append(edgeInputList, edgeInput)
		}
	}
	transitSearchInput.Edges = edgeInputList
	transitopsres, err := prosimoClient.TransitSetupSummary(ctx, transitSearchInput)
	if err != nil {
		return diag.FromErr(err)
	}
	setupIDList := []string{}
	for _, inops := range transitopsres {
		for _, inEdge := range edgeInputList {
			if inops.CloudType == inEdge.CloudType && inops.CloudRegion == inEdge.CloudRegion {
				setupIDList = append(setupIDList, inops.ID)
			}
		}
	}
	err_res_delete := prosimoClient.DeleteTransitSetup(ctx, setupIDList)
	if err_res_delete != nil {
		return diag.FromErr(err_res_delete)
	}

	log.Println("[DEBUG] Deleted Transit setup")
	d.SetId("")
	return diags
}
