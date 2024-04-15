package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePrivateLinkMapping() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Private Link Mappings.",
		CreateContext: resourcePVMCreate,
		UpdateContext: resourcePVMUpdate,
		DeleteContext: resourcePVMDelete,
		ReadContext:   resourcePVMRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Private Link Source Name",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target App/Network name",
			},
			"hosted_zones": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the source VPC",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target Domain name",
						},
						"private_hosted_zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Private Hosted Zones",
						},
					},
				},
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

func resourcePVMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	hostedZonesList := []client.Hosted_Zones{}
	inpvm := &client.PL_Map{}
	source := d.Get("source").(string)
	target := d.Get("target").(string)
	if v, ok := d.GetOk("hosted_zones"); ok {
		for i, _ := range v.([]interface{}) {
			hostedZoneConfig := v.([]interface{})[i].(map[string]interface{})
			vpcName := hostedZoneConfig["vpc_name"]
			domainName := hostedZoneConfig["domain_name"]
			hostedZones := &client.Hosted_Zones{
				Name: hostedZoneConfig["private_hosted_zone"].(string),
			}
			pvsList, err := prosimoClient.GetPrivateLinkSource(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, pvs := range pvsList {
				if pvs.Name == source {
					serviceInputSource := &client.Service_Input{
						Name: source,
						ID:   pvs.ID,
					}
					inpvm.Source = serviceInputSource
					for _, cloudSource := range *pvs.CloudSources {
						if cloudSource.CloudNetwork.Name == vpcName {
							hostedZones.SourceID = cloudSource.ID
							break
						}
					}
					break
				}
			}
			appSearchInput := &client.AppOnboardSearch{}
			appList, err := prosimoClient.SearchAppOnboardApps(ctx, appSearchInput)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, app := range appList.Data.Records {
				if app.App_Name == target {
					serviceInputTarget := &client.Service_Input{
						Name: target,
						ID:   app.ID,
					}
					inpvm.Target = serviceInputTarget
					for _, appUrl := range app.AppURLs {
						if appUrl.AppFqdn == domainName {
							hostedZones.DomainID = appUrl.ID
							break
						}
					}
					break
				}
			}

			hostedZonesList = append(hostedZonesList, *hostedZones)
		}
		inpvm.HostedZones = &hostedZonesList
	}

	postres, err := prosimoClient.CreatePrivateLinkMapping(ctx, inpvm)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[INFO] Waiting for task id %s to complete", postres.PlMapRes.TaskID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, postres.PlMapRes.TaskID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] task %s is successful", postres.PlMapRes.TaskID)
	}
	log.Printf("[INFO] New Private Link Mapping with id  %s is deployed", postres.PlMapRes.ID)
	d.SetId(postres.PlMapRes.ID)
	resourcePVMRead(ctx, d, meta)
	return diags
}

func resourcePVMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	updateReq := false
	if d.HasChange("source") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Private Link Source",
			Detail:   "Private Link Source can't be modified",
		})
		return diags
	}
	if d.HasChange("target") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Target App",
			Detail:   "Target App details can't be modified",
		})
		return diags
	}

	if d.HasChange("cloud_region") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Cloud Regions",
			Detail:   "Cloud Region can't be modified",
		})
		return diags
	}

	if d.HasChange("hosted_zones") && !d.IsNewResource() {
		updateReq = true
	}
	if updateReq {
		hostedZonesList := []client.Hosted_Zones{}
		inpvm := &client.PL_Map{
			ID: d.Id(),
		}
		source := d.Get("source").(string)
		target := d.Get("target").(string)
		if v, ok := d.GetOk("hosted_zones"); ok {
			for i, _ := range v.([]interface{}) {
				hostedZoneConfig := v.([]interface{})[i].(map[string]interface{})
				vpcName := hostedZoneConfig["vpc_name"]
				domainName := hostedZoneConfig["domain_name"]
				hostedZones := &client.Hosted_Zones{
					Name: hostedZoneConfig["private_hosted_zone"].(string),
				}
				pvsList, err := prosimoClient.GetPrivateLinkSource(ctx)
				if err != nil {
					return diag.FromErr(err)
				}
				for _, pvs := range pvsList {
					if pvs.Name == source {
						serviceInputSource := &client.Service_Input{
							Name: source,
							ID:   pvs.ID,
						}
						inpvm.Source = serviceInputSource
						for _, cloudSource := range *pvs.CloudSources {
							if cloudSource.CloudNetwork.Name == vpcName {
								hostedZones.SourceID = cloudSource.ID
								break
							}
						}
						break
					}
				}
				appSearchInput := &client.AppOnboardSearch{}
				appList, err := prosimoClient.SearchAppOnboardApps(ctx, appSearchInput)
				if err != nil {
					return diag.FromErr(err)
				}
				for _, app := range appList.Data.Records {
					if app.App_Name == target {
						serviceInputTarget := &client.Service_Input{
							Name: target,
							ID:   app.ID,
						}
						inpvm.Target = serviceInputTarget
						for _, appUrl := range app.AppURLs {
							if appUrl.AppFqdn == domainName {
								hostedZones.DomainID = appUrl.ID
								break
							}
						}
						break
					}
				}

				hostedZonesList = append(hostedZonesList, *hostedZones)
			}
			inpvm.HostedZones = &hostedZonesList
		}
		postres, err := prosimoClient.UpdatePrivateLinkMapping(ctx, inpvm)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[INFO] Waiting for task id %s to complete", postres.PlMapRes.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete(ctx, d, meta, postres.PlMapRes.TaskID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", postres.PlMapRes.TaskID)
		}
		log.Printf("[INFO]  Private Link Source with id  %s is updated", postres.PlMapRes.ID)
	}
	resourcePVSRead(ctx, d, meta)
	return diags
}

func resourcePVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	pvmID := d.Id()

	log.Printf("[DEBUG] Get Private Link Mapping with id  %s", pvmID)

	pvm, err := prosimoClient.GetPrivateLinkMappingByID(ctx, pvmID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", pvm.ID)

	return diags
}

func resourcePVMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pvmID := d.Id()

	res, err := prosimoClient.DeletePrivateLinkMapping(ctx, pvmID)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[INFO] Waiting for task id %s to complete", res.PlMapRes.TaskID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, res.PlMapRes.TaskID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] task %s is successful", res.PlMapRes.TaskID)
	}
	log.Printf("[DEBUG] Deleted Private Link Mapping with - id - %s", pvmID)
	d.SetId("")

	return diags
}
