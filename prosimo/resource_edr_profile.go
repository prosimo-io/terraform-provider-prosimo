package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEdrProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Endpoint Security Profiles.",
		CreateContext: resourceEdrConfigCreate,
		ReadContext:   resourceEdrConfigRead,
		DeleteContext: resourceEdrConfigDelete,
		UpdateContext: resourceEdrConfigUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the EDR Profile",
			},
			"vendor": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(client.GetEDRvendorTypes(), false),
				Description:  "Name of the EDR vendor, For now only CrowdStrike is supported",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Integration Status",
			},
			"auth": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Vendor Auth inputs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetEDRbaseurlTypes(), false),
						},

						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"customer_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mssp": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceEdrConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	edrConfig := &client.EDR_Config{}
	auth := client.AUTH{}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		edrConfig.Name = name
	}
	if v, ok := d.GetOk("vendor"); ok {
		name := v.(string)
		edrConfig.Vendor = name
	}

	if v, ok := d.GetOk("auth"); ok {
		authdetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := authdetails["base_url"]; ok {
			baseUrl := v.(string)
			auth.BaseURL = baseUrl
		}
		if v, ok := authdetails["client_id"]; ok {
			clientID := v.(string)
			auth.ClientID = clientID
		}
		if v, ok := authdetails["client_secret"]; ok {
			clientSecret := v.(string)
			auth.ClientSecret = clientSecret
		}
		if v, ok := authdetails["customer_id"]; ok {
			customerID := v.(string)
			auth.CustomerID = customerID
		}
		if v, ok := authdetails["mssp"]; ok {
			mssp := v.(bool)
			auth.MSSP = mssp
		}
		edrConfig.Auth = auth
	}
	log.Printf("[DEBUG] Creating EDR Config : %v", edrConfig)
	createEDR, err := prosimoClient.CreateEDRConf(ctx, edrConfig)
	if err != nil {
		log.Printf("[DEBUG] Error in creating EDRconfig")
		return diag.FromErr(err)
	}
	d.SetId(createEDR.EdrConfig.Id)
	return resourceEdrConfigRead(ctx, d, meta)
}

func resourceEdrConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	edrID := d.Id()

	log.Printf("[DEBUG] Get EDR profile for %s", edrID)

	res, err := prosimoClient.GetEDRConf(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	var edr *client.EDR_Config
	for _, returnedEdr := range res.EdrConfigList {
		if returnedEdr.Id == edrID {
			edr = returnedEdr
			break
		}
	}
	if edr == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get EDR config",
		})
		return diags
	}

	d.Set("name", edr.Name)
	d.Set("vendor", edr.Vendor)
	d.Set("status", edr.Status)

	return diags
}

func resourceEdrConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	edrConfig := &client.EDR_Config{}
	auth := client.AUTH{}

	updateReq := false
	edrID := d.Id()
	edrConfig.Id = edrID

	if d.HasChange("name") && !d.IsNewResource() {
		updateReq = true
	}
	edrConfig.Name = d.Get("name").(string)

	if d.HasChange("vendor") && !d.IsNewResource() {
		updateReq = true
	}
	edrConfig.Vendor = d.Get("vendor").(string)

	if v, ok := d.GetOk("auth"); ok {
		authdetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if d.HasChange("authdetails.base_url") && !d.IsNewResource() {
			updateReq = true
		}
		auth.BaseURL = authdetails["base_url"].(string)

		if d.HasChange("authdetails.client_id") && !d.IsNewResource() {
			updateReq = true
		}
		auth.ClientID = authdetails["client_id"].(string)

		if d.HasChange("authdetails.client_secret") && !d.IsNewResource() {
			updateReq = true
		}
		auth.ClientSecret = authdetails["client_secret"].(string)
		if d.HasChange("authdetails.customer_id") && !d.IsNewResource() {
			updateReq = true
		}
		auth.CustomerID = authdetails["customer_id"].(string)

		if d.HasChange("authdetails.mssp") && !d.IsNewResource() {
			updateReq = true
		}
		auth.MSSP = authdetails["mssp"].(bool)
		if v, ok := authdetails["mssp"]; ok {
			mssp := v.(bool)
			auth.MSSP = mssp
		}
		edrConfig.Auth = auth
	}

	if len(diags) > 0 {
		return diags
	}

	if updateReq {
		log.Printf("[DEBUG] Updating EDR Config : %v", edrConfig)
		_, err := prosimoClient.UpdateEDRConf(ctx, edrConfig)
		if err != nil {
			log.Printf("[DEBUG] Error in updating EDRconfig")
			return diag.FromErr(err)
		}
		// d.SetId(updateEDR.EdrConfig.Id)
	}
	return resourceEdrConfigRead(ctx, d, meta)
}

func resourceEdrConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics

	edrID := d.Id()
	err := prosimoClient.DeleteEDRConf(ctx, edrID)

	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
