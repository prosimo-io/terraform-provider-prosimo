package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceInsertion() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify service insertion policy.",
		CreateContext: resourceSICreate,
		UpdateContext: resourceSIUpdate,
		DeleteContext: resourceSIDelete,
		ReadContext:   resourceSIRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of Service Insertion",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Policy Namespace, Defaults to default",
			},
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Shared Service",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "FWPolicy_IP",
				Description: "Service Insertion Type",
			},
			"prosimo_managed_routing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TRUE if you would like Prosimo to update Firewal VNET Roue Table",
			},
			"route_tables": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of Route Table ID",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service Insertion Deployment Status",
			},
			"source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"networks": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source Network Name",
									},
								},
							},
						},
					},
				},
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"networks": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Target Network Name",
									},
								},
							},
						},
						"apps": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Target App Name",
									},
								},
							},
						},
					},
				},
			},
			"ip_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Insertion Policy Rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_addresses": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Source Ip Address",
						},
						"source_ports": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Source Port",
						},
						"destination_addresses": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Target Ip Address",
						},
						"destination_ports": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Destination Port",
						},
						"protocols": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Protocols",
						},
					},
				},
			},
			// "wait_for_rollout": {
			// 	Type:        schema.TypeBool,
			// 	Description: "Wait for the rollout of the task to complete. Defaults to true.",
			// 	Default:     true,
			// 	Optional:    true,
			// },
		},
	}
}

func resourceSICreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	// serviceInsertionInput := client.Service_Insertion{}

	sourceNetwork := &client.Source{}
	sourceNetworkList := []client.Service_Input{}
	if v, ok := d.GetOk("source"); ok {
		for i, _ := range v.([]interface{}) {
			sourceConfig := v.(*schema.Set).List()[i].(map[string]interface{})
			if v, ok := sourceConfig["networks"].(*schema.Set); ok && v.Len() > 0 {

				for i, val := range v.List() {
					_ = val

					networkConfig := v.List()[i].(map[string]interface{})
					selectnetworkName := networkConfig["name"].(string)
					selectnetworkid, err := prosimoClient.GetNetworkID(ctx, selectnetworkName)
					if err != nil {
						return diag.FromErr(err)
					}
					serviceInput := client.Service_Input{
						Name: selectnetworkName,
						ID:   selectnetworkid,
					}
					sourceNetworkList = append(sourceNetworkList, serviceInput)
				}
				sourceNetwork.Networks = sourceNetworkList
			}
		}
	}

	target := &client.Target{}
	targetNetworkList := []client.Service_Input{}
	targetAppList := []client.Service_Input{}
	if v, ok := d.GetOk("target"); ok {
		for i, _ := range v.([]interface{}) {
			targetConfig := v.(*schema.Set).List()[i].(map[string]interface{})
			if v, ok := targetConfig["networks"].(*schema.Set); ok && v.Len() > 0 {

				for i, val := range v.List() {
					_ = val

					networkConfig := v.List()[i].(map[string]interface{})
					selectnetworkName := networkConfig["name"].(string)
					selectnetworkid, err := prosimoClient.GetNetworkID(ctx, selectnetworkName)
					if err != nil {
						return diag.FromErr(err)
					}
					serviceInput := client.Service_Input{
						Name: selectnetworkName,
						ID:   selectnetworkid,
					}
					targetNetworkList = append(targetNetworkList, serviceInput)
				}
				target.Networks = targetNetworkList
			}

			if v, ok := targetConfig["apps"].(*schema.Set); ok && v.Len() > 0 {

				for i, _ := range v.List() {
					appConfig := v.List()[i].(map[string]interface{})
					selectappName := appConfig["name"].(string)
					selectappid, err := prosimoClient.GetAppID(ctx, selectappName)
					if err != nil {
						return diag.FromErr(err)
					}
					serviceInput := client.Service_Input{
						Name: selectappName,
						ID:   selectappid,
					}
					targetAppList = append(targetAppList, serviceInput)
				}
				target.Apps = targetAppList
			}
		}
	}
	ipRulesConfigInputList := []client.IpRule{}
	if v, ok := d.GetOk("ip_rules"); ok {
		for i, _ := range v.([]interface{}) {

			ipRulesConfig := v.([]interface{})[i].(map[string]interface{})
			ipRulesConfigInput := &client.IpRule{
				SrcAddr:  expandStringList(ipRulesConfig["source_addresses"].([]interface{})),
				SrcPort:  expandStringList(ipRulesConfig["source_ports"].([]interface{})),
				DestAddr: expandStringList(ipRulesConfig["destination_addresses"].([]interface{})),
				DestPort: expandStringList(ipRulesConfig["destination_ports"].([]interface{})),
				Protocol: expandStringList(ipRulesConfig["protocols"].([]interface{})),
			}
			ipRulesConfigInputList = append(ipRulesConfigInputList, *ipRulesConfigInput)
		}
	}

	ss, err := prosimoClient.GetSharedServiceByName(ctx, d.Get("service_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	serviceInsertionInput := &client.Service_Insertion{
		Name:         d.Get("name").(string),
		Service_name: d.Get("service_name").(string),
		RegionID:     ss.Region.ID,
		Type:         d.Get("type").(string),
		Source:       sourceNetwork,
		Target:       target,
		IpRules:      &ipRulesConfigInputList,
	}
	nameSpaceDetails, err := prosimoClient.GetNamespaceByName(ctx, d.Get("namespace").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	serviceInsertionInput.NameSpaceID = nameSpaceDetails.ID

	if ss.Region.CloudType == client.AzureCloudType {
		if v, ok := d.GetOk("prosimo_managed_routing"); ok {
			serviceInsertionInput.ProsimoManagedRouting = v.(bool)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing Prosimo managed Routing options",
				Detail:   "Missing Prosimo managed Routing options",
			})
			return diags
		}
		if serviceInsertionInput.ProsimoManagedRouting {
			if v, ok := d.GetOk("route_tables"); ok {
				serviceInsertionInput.RouteTable = expandStringList(v.([]interface{}))
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Missing Route Table ID details",
					Detail:   "Missing Route Table ID details",
				})
				return diags
			}

		}
	}
	siRes, err := prosimoClient.CreateServiceInsertion(ctx, serviceInsertionInput)
	if err != nil {
		return diag.FromErr(err)
	}
	// if d.Get("wait_for_rollout").(bool) {
	// 	log.Printf("[INFO] Waiting for task id %s to complete", siRes.Service_Insertion_Response.ID)
	// 	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
	// 		retryUntilTaskComplete(ctx, d, meta, siRes.Service_Insertion_Response.ID))
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// 	log.Printf("[INFO] task %s is successful", siRes.Service_Insertion_Response.ID)
	// }
	// si, err := prosimoClient.GetServiceInsertionByName(ctx, serviceInsertionInput.Name)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// log.Println("si.ID", si.ID)
	log.Printf("[INFO]Service Insertion  with id  %s is updated", siRes.Service_Insertion_Response.ID)
	d.SetId(siRes.Service_Insertion_Response.ID)
	resourceSIRead(ctx, d, meta)

	return diags
}

func resourceSIUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)

	updateReq := false
	if d.HasChange("name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Name",
			Detail:   "Name can't be modified",
		})
		return diags
	}

	if d.HasChange("type") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Service Insertion Type",
			Detail:   "Service Insertion Type can't be modified",
		})
		return diags
	}

	if d.HasChange("service_name") {
		updateReq = true
	}
	if d.HasChange("source") && !d.IsNewResource() {
		updateReq = true
	}
	if d.HasChange("target") && !d.IsNewResource() {
		updateReq = true
	}
	if d.HasChange("ip_rules") && !d.IsNewResource() {
		updateReq = true
	}

	if updateReq {
		sourceNetwork := &client.Source{}
		sourceNetworkList := []client.Service_Input{}
		if v, ok := d.GetOk("source"); ok {
			for i, _ := range v.([]interface{}) {
				sourceConfig := v.(*schema.Set).List()[i].(map[string]interface{})
				if v, ok := sourceConfig["networks"].(*schema.Set); ok && v.Len() > 0 {

					for i, val := range v.List() {
						_ = val

						networkConfig := v.List()[i].(map[string]interface{})
						selectnetworkName := networkConfig["name"].(string)
						selectnetworkid, err := prosimoClient.GetNetworkID(ctx, selectnetworkName)
						if err != nil {
							return diag.FromErr(err)
						}
						serviceInput := client.Service_Input{
							Name: selectnetworkName,
							ID:   selectnetworkid,
						}
						sourceNetworkList = append(sourceNetworkList, serviceInput)
					}
					sourceNetwork.Networks = sourceNetworkList
				}
			}
		}

		target := &client.Target{}
		targetNetworkList := []client.Service_Input{}
		targetAppList := []client.Service_Input{}
		if v, ok := d.GetOk("target"); ok {
			for i, _ := range v.([]interface{}) {
				targetConfig := v.(*schema.Set).List()[i].(map[string]interface{})
				if v, ok := targetConfig["networks"].(*schema.Set); ok && v.Len() > 0 {

					for i, val := range v.List() {
						_ = val

						networkConfig := v.List()[i].(map[string]interface{})
						selectnetworkName := networkConfig["name"].(string)
						selectnetworkid, err := prosimoClient.GetNetworkID(ctx, selectnetworkName)
						if err != nil {
							return diag.FromErr(err)
						}
						serviceInput := client.Service_Input{
							Name: selectnetworkName,
							ID:   selectnetworkid,
						}
						targetNetworkList = append(targetNetworkList, serviceInput)
					}
					target.Networks = targetNetworkList
				}

				if v, ok := targetConfig["apps"].(*schema.Set); ok && v.Len() > 0 {

					for i, _ := range v.List() {
						appConfig := v.List()[i].(map[string]interface{})
						selectappName := appConfig["name"].(string)
						selectappid, err := prosimoClient.GetAppID(ctx, selectappName)
						if err != nil {
							return diag.FromErr(err)
						}
						serviceInput := client.Service_Input{
							Name: selectappName,
							ID:   selectappid,
						}
						targetAppList = append(targetAppList, serviceInput)
					}
					target.Apps = targetAppList
				}
			}
		}
		ipRulesConfigInputList := []client.IpRule{}
		if v, ok := d.GetOk("ip_rules"); ok {
			for i, _ := range v.([]interface{}) {

				ipRulesConfig := v.([]interface{})[i].(map[string]interface{})
				ipRulesConfigInput := &client.IpRule{
					SrcAddr:  expandStringList(ipRulesConfig["source_addresses"].([]interface{})),
					SrcPort:  expandStringList(ipRulesConfig["source_ports"].([]interface{})),
					DestAddr: expandStringList(ipRulesConfig["destination_addresses"].([]interface{})),
					DestPort: expandStringList(ipRulesConfig["destination_ports"].([]interface{})),
					Protocol: expandStringList(ipRulesConfig["protocols"].([]interface{})),
				}
				ipRulesConfigInputList = append(ipRulesConfigInputList, *ipRulesConfigInput)
			}
		}

		ss, err := prosimoClient.GetSharedServiceByName(ctx, d.Get("service_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		serviceInsertionInput := &client.Service_Insertion{
			Name:               d.Get("name").(string),
			ID:                 d.Id(),
			Service_name:       d.Get("service_name").(string),
			RegionID:           ss.Region.ID,
			Type:               d.Get("type").(string),
			Source:             sourceNetwork,
			Target:             target,
			IpRules:            &ipRulesConfigInputList,
			CloudType:          ss.Region.CloudType,
			CloudRegion:        ss.Region.CloudRegion,
			SharedServiceCreds: ss.Region.CloudKeyID,
			GwLoadbalancerID:   ss.Region.GwLoadBalancerID,
			ServiceID:          ss.ID,
		}
		nameSpaceDetails, err := prosimoClient.GetNamespaceByName(ctx, d.Get("namespace").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		serviceInsertionInput.NameSpaceID = nameSpaceDetails.ID
		if ss.Region.CloudType == client.AzureCloudType {
			if v, ok := d.GetOk("prosimo_managed_routing"); ok {
				serviceInsertionInput.ProsimoManagedRouting = v.(bool)
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Missing Prosimo managed Routing options",
					Detail:   "Missing Prosimo managed Routing options",
				})
				return diags
			}
			if serviceInsertionInput.ProsimoManagedRouting {
				if v, ok := d.GetOk("route_tables"); ok {
					serviceInsertionInput.RouteTable = expandStringList(v.([]interface{}))
				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Missing Route Table ID details",
						Detail:   "Missing Route Table ID details",
					})
					return diags
				}

			}
		}
		putRes, err := prosimoClient.UpdateServiceInsertion(ctx, serviceInsertionInput)
		if err != nil {
			return diag.FromErr(err)
		}
		// if d.Get("wait_for_rollout").(bool) {
		// 	log.Printf("[INFO] Waiting for task id %s to complete", siRes.Service_Insertion_Response.ID)
		// 	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
		// 		retryUntilTaskComplete(ctx, d, meta, siRes.Service_Insertion_Response.ID))
		// 	if err != nil {
		// 		return diag.FromErr(err)
		// 	}
		// 	log.Printf("[INFO] task %s is successful", siRes.Service_Insertion_Response.ID)
		// }
		log.Printf("[INFO]Service Insertion  with id  %s is updated", putRes.Service_Insertion_Response.ID)
	}

	resourceSSRead(ctx, d, meta)

	return diags
}

func resourceSIRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	siID := d.Id()

	log.Printf("[DEBUG] Get Service Insertion Policy for %s", siID)

	si, err := prosimoClient.GetServiceInsertionByID(ctx, siID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", si.ID)
	d.Set("name", si.Name)
	d.Set("type", si.Type)
	d.Set("status", si.Status)
	d.Set("service_name", si.Service_name)

	return diags
}

func resourceSIDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	siID := d.Id()

	res, err := prosimoClient.DeleteServiceInsertion(ctx, siID)
	if err != nil {
		return diag.FromErr(err)
	}
	// if d.Get("wait_for_rollout").(bool) {
	// 	log.Printf("[INFO] Waiting for task id %s to complete", res.Service_Insertion_Response.ID)
	// 	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
	// 		retryUntilTaskComplete(ctx, d, meta, res.Service_Insertion_Response.ID))
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// 	log.Printf("[INFO] task %s is successful", res.Service_Insertion_Response.ID)
	// }
	log.Printf("[DEBUG] Deleted Service Insertion with - id - %s", res.Service_Insertion_Response.ID)
	d.SetId("")

	return diags
}
