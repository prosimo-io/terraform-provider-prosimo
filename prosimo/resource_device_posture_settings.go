package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDPSettings() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify device posture settings.",
		CreateContext: resourceDPSettingsUpdate,
		UpdateContext: resourceDPSettingsUpdate,
		ReadContext:   resourceDPSettingsRead,
		DeleteContext: resourceDPSettingsDelete,
		// UpdateContext: resourceEdrProfileUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enable_dp_feature": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this to TRUE to enable device posture",
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
			// Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceDPSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	dpFetaureSettings := client.DevicePosture_Settings{}

	if v, ok := d.GetOk("enable_dp_feature"); ok {
		enableDPFeature := v.(bool)
		dpFetaureSettings.Enabled = enableDPFeature
	}

	getInputRes, err := prosimoClient.UpdateDPSettings(ctx, dpFetaureSettings)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[DEBUG] d Waiting for task id %s to complete", getInputRes.DPProfileUpdateRes.TaskID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, getInputRes.DPProfileUpdateRes.TaskID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] task %s is successful", getInputRes.DPProfileUpdateRes.TaskID)
	}

	d.SetId("DP_Settings")
	return resourceDPSettingsRead(ctx, d, meta)
}

func resourceDPSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	_, err := prosimoClient.GetDPSettings(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceDPSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	dpFetaureSettings := client.DevicePosture_Settings{}
	dpFetaureSettings.Enabled = false

	dpSettings_res, err := prosimoClient.UpdateDPSettings(ctx, dpFetaureSettings)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[DEBUG] Waiting for task id %s to complete", dpSettings_res.DPProfileUpdateRes.TaskID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, dpSettings_res.DPProfileUpdateRes.TaskID))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDPProfileRead(ctx, d, meta)

}
