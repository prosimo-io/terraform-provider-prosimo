package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewallManger() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Firewall Manager.",
		CreateContext: resourceFMCreate,
		UpdateContext: resourceFMUpdate,
		DeleteContext: resourceFMDelete,
		ReadContext:   resourceFMRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"integration_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of Integration, e.g: panorama",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target IP Address",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Target API Key for authentication",
			},
			"license_settings": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "License Settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"license_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Mode of license, e.g: Bring your own license (BYOL) Pay as you go (PAYG)",
						},
						"firewall_family": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Firewall Family",
						},
						"instance_family": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance Family",
						},
						"license_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "License Type, e.g: Bundle1, Bundle2",
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceFMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics
	// fwmConfig := client.FWMConfig{}

	prosimoClient := meta.(*client.ProsimoClient)
	fwmConfig := client.FWMConfig{
		IntegrationType: d.Get("integration_type").(string),
		IPAddress:       d.Get("ip_address").(string),
		APIKey:          d.Get("api_key").(string),
	}
	if v, ok := d.GetOk("license_settings"); ok {
		licenseSettings := v.(*schema.Set).List()[0].(map[string]interface{})
		licenseSettingsInput := &client.LicenseDetails{
			LicenseMode:    licenseSettings["license_mode"].(string),
			FirewallFamily: licenseSettings["firewall_family"].(string),
			InstanceFamily: licenseSettings["instance_family"].(string),
			LicenseType:    licenseSettings["license_type"].(string),
		}
		fwmConfig.LicenseDetails = licenseSettingsInput
		log.Println("fwmConfig.LicenseDetails", fwmConfig.LicenseDetails)
	}

	log.Printf("[DEBUG]Creating FireWallManager : %v", fwmConfig)
	createFirewallManager, err := prosimoClient.CreateFirewallManager(ctx, &fwmConfig)
	if err != nil {
		log.Printf("[ERROR] Error in creating firewall")
		return diag.FromErr(err)
	}
	d.SetId(createFirewallManager.PlMapRes.ID)
	return resourceFMRead(ctx, d, meta)
}

func resourceFMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	if d.HasChange("integration_type") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify Integration Type",
			Detail:   "Integration Type can't be modified",
		})
		return diags
	}

	if d.HasChange("license_settings") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify License Settings",
			Detail:   "License Settings can't be modified",
		})
		return diags
	}
	update_req := false
	if d.HasChange("ip_address") && !d.IsNewResource() {
		update_req = true
	}

	if d.HasChange("api_key") && !d.IsNewResource() {
		update_req = true
	}

	if update_req {
		inputFM := &client.FWMConfig{
			ID:        d.Id(),
			IPAddress: d.Get("ip_address").(string),
			APIKey:    d.Get("api_key").(string),
		}

		log.Printf("[DEBUG]Updating Firewall Manager : %v", inputFM)
		err := prosimoClient.UpdateFirewallManager(ctx, inputFM)
		if err != nil {
			log.Printf("[ERROR] Error Updating Firewall Manager")
			return diag.FromErr(err)
		}
	}
	return resourceFMRead(ctx, d, meta)
}

func resourceFMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	nsID := d.Id()

	log.Printf("[DEBUG] Get firewall Manager with id  %s", nsID)

	fm, err := prosimoClient.GetFirewallManagerByID(ctx, nsID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", fm.ID)
	d.Set("integration_type", fm.IntegrationType)
	d.Set("ip_address", fm.IPAddress)
	d.Set("api_key", fm.APIKey)
	return diags
}

func resourceFMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	fmID := d.Id()

	err := prosimoClient.DeleteFirewallManager(ctx, fmID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Deleted firewall Manager with - id - %s", fmID)

	return diags
}
