package prosimo

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSharedServices() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify shared services.",
		CreateContext: resourceSSCreate,
		UpdateContext: resourceSSUpdate,
		DeleteContext: resourceSSDelete,
		ReadContext:   resourceSSRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Shared Service",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"teamid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource TEAM ID",
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Firewall",
				Description: "Type of Shared Service",
			},
			"region": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "EX: us-west-2, eu-east-1",
						},
						"gateway_lb": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Gateway Load Balance Service Name",
						},
						"cloud_creds_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cloud account under which application is hosted",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Shared Service Deployment Status",
			},
			"onboard": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like to onboard  a saved Shared Service with out any config changes",
			},
			"decommission": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like to decommission an already onboarded Shared Service",
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

func resourceSSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	offboardFlag := d.Get("decommission").(bool)
	if offboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid  decommission flag.",
			Detail:   "decommission can't be set to true while creating  resource.",
		})
		return diags
	}

	regionInput := &client.Region{}
	if v, ok := d.GetOk("region"); ok {
		regionConfig := v.(*schema.Set).List()[0].(map[string]interface{})
		cloudCredName := regionConfig["cloud_creds_name"].(string)
		cloudRregion := regionConfig["cloud_region"].(string)

		// validate cloud name
		cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, cloudCredName)
		if err != nil {
			return diag.FromErr(err)
		}

		if cloudCreds == nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get Cloud Credentials",
				Detail:   fmt.Sprintf("Unable to find Cloud Credentials for name %s", cloudCredName),
			})

			return diags
		}

		// validate cloud region
		regionExists, err := prosimoClient.CheckIfCloudRegionExists(ctx, cloudCreds.ID, cloudRregion)
		if err != nil {
			return diag.FromErr(err)
		}

		if !regionExists {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get Cloud Region",
				Detail:   fmt.Sprintf("Unable to find region %s for Cloud Name %s", cloudRregion, cloudCredName),
			})

			return diags
		}
		regionInput = &client.Region{
			CloudRegion:      cloudRregion,
			CloudKeyID:       cloudCreds.ID,
			CloudType:        cloudCreds.CloudType,
			GwLoadBalancerID: regionConfig["gateway_lb"].(string),
			CloudZones:       "",
		}
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing Region details",
			Detail:   fmt.Sprintln("Unable to find region config in tf input"),
		})

		return diags
	}
	ssInput := &client.Shared_Service{
		Name:   d.Get("name").(string),
		Type:   d.Get("service_type").(string),
		Region: regionInput,
	}

	// create SharedService
	ssResponseData, err := prosimoClient.CreateSharedService(ctx, ssInput)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("onboard").(bool) {
		onboardresponse, err := prosimoClient.OnboardSharedService(ctx, ssResponseData.Shared_Service_Response.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[INFO] Waiting for task id %s to complete", onboardresponse.Shared_Service_Response.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskCompleteSharedService(ctx, d, meta, onboardresponse.Shared_Service_Response.TaskID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", onboardresponse.Shared_Service_Response.TaskID)
		}
	}

	log.Printf("[DEBUG] Deployed Shared Service - name - %s", ssInput.Name)
	d.SetId(ssResponseData.Shared_Service_Response.ID)

	resourceSSRead(ctx, d, meta)

	return diags
}

func resourceSSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)

	ssOnboardFlag := d.Get("onboard").(bool)
	ssOffboardFlag := d.Get("decommission").(bool)
	if ssOnboardFlag && ssOffboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid onboard and decommission flag combination.",
			Detail:   "Both onboard and decommission have been set to true.",
		})
		return diags
	}

	updateReq := false
	if d.HasChange("name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Name",
			Detail:   "Name can't be modified",
		})
		return diags
	}

	if d.HasChange("service_type") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Service Type",
			Detail:   "Service Type can't be modified",
		})
		return diags
	}

	if d.HasChange("region") {
		updateReq = true
	}

	if d.HasChange("onboard") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("decommission") && !d.IsNewResource() {
		updateReq = true
	}

	//Offboard Shared Service
	if updateReq {
		offBoardApp := false
		if d.HasChange("decommission") && !d.IsNewResource() {
			isDecommission := d.Get("decommission").(bool)
			if isDecommission {
				offBoardApp = true
				offboardresponse, err := prosimoClient.DecomSharedService(ctx, d.Id())
				if err != nil {
					return diag.FromErr(err)
				}
				if d.Get("wait_for_rollout").(bool) {
					log.Printf("[INFO] Waiting for task id %s to complete", offboardresponse.Shared_Service_Response.TaskID)
					err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
						retryUntilTaskComplete(ctx, d, meta, offboardresponse.Shared_Service_Response.TaskID))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[INFO] task %s is successful", offboardresponse.Shared_Service_Response.TaskID)
				}
			}
		}

		if !offBoardApp {
			regionInput := &client.Region{}
			if v, ok := d.GetOk("region"); ok {
				regionConfig := v.(*schema.Set).List()[0].(map[string]interface{})
				cloudCredName := regionConfig["cloud_creds_name"].(string)
				cloudRregion := regionConfig["cloud_region"].(string)

				// validate cloud name
				cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, cloudCredName)
				if err != nil {
					return diag.FromErr(err)
				}

				if cloudCreds == nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Unable to get Cloud Credentials",
						Detail:   fmt.Sprintf("Unable to find Cloud Credentials for name %s", cloudCredName),
					})

					return diags
				}

				// validate cloud region
				regionExists, err := prosimoClient.CheckIfCloudRegionExists(ctx, cloudCreds.ID, cloudRregion)
				if err != nil {
					return diag.FromErr(err)
				}

				if !regionExists {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Unable to get Cloud Region",
						Detail:   fmt.Sprintf("Unable to find region %s for Cloud Name %s", cloudRregion, cloudCredName),
					})

					return diags
				}
				regionInput = &client.Region{
					CloudRegion:      cloudRregion,
					CloudKeyID:       cloudCreds.ID,
					CloudType:        cloudCreds.CloudType,
					GwLoadBalancerID: regionConfig["gateway_lb"].(string),
					CloudZones:       "",
				}
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Missing Region details",
					Detail:   fmt.Sprintln("Unable to find region config in tf input"),
				})

				return diags
			}
			ssInput := &client.Shared_Service{
				ID:     d.Id(),
				Name:   d.Get("name").(string),
				Type:   d.Get("service_type").(string),
				Region: regionInput,
			}

			updateres, err := prosimoClient.UpdateSharedService(ctx, ssInput)
			if err != nil {
				return diag.FromErr(err)
			}
			if d.Get("onboard").(bool) {
				onboardresponse, err := prosimoClient.OnboardSharedService(ctx, updateres.Shared_Service_Response.ID)
				if err != nil {
					return diag.FromErr(err)
				}
				if d.Get("wait_for_rollout").(bool) {
					log.Printf("[INFO] Waiting for task id %s to complete", onboardresponse.Shared_Service_Response.TaskID)
					err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
						retryUntilTaskCompleteSharedService(ctx, d, meta, onboardresponse.Shared_Service_Response.TaskID))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[INFO] task %s is successful", onboardresponse.Shared_Service_Response.TaskID)
				}
			}
			log.Printf("[DEBUG] Updated Shared Service - id - %s", updateres.Shared_Service_Response.ID)
		}
	}

	resourceSSRead(ctx, d, meta)

	return diags
}

func resourceSSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	ssID := d.Id()

	log.Printf("[DEBUG] Get Shared Service for %s", ssID)

	ss, err := prosimoClient.GetSharedServiceByID(ctx, ssID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", ss.ID)
	d.Set("teamid", ss.TeamID)
	d.Set("name", ss.Name)
	d.Set("service_type", ss.Type)
	d.Set("status", ss.Status)
	d.Set("onboard", d.Get("onboard").(bool))
	d.Set("decommission", d.Get("decommission").(bool))

	return diags
}

func resourceSSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	ssID := d.Id()
	ss, err := prosimoClient.GetSharedServiceByID(ctx, ssID)
	if err != nil {
		return diag.FromErr(err)
	}
	if ss.Status == "DEPLOYED" {
		offboardresponse, err := prosimoClient.DecomSharedService(ctx, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[INFO] Waiting for task id %s to complete", offboardresponse.Shared_Service_Response.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete(ctx, d, meta, offboardresponse.Shared_Service_Response.TaskID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", offboardresponse.Shared_Service_Response.TaskID)
		}
	}
	delete_err := prosimoClient.DeleteSharedService(ctx, ssID)
	if delete_err != nil {
		return diag.FromErr(delete_err)
	}
	log.Printf("[DEBUG] Deleted Shared Service with - id - %s", ssID)
	d.SetId("")

	return diags
}
