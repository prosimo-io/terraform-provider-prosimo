package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDynamicRisk() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to modify dynamic risk settings.",
		CreateContext: resourceDynamicRiskCreate,
		ReadContext:   resourceDynamicRiskRead,
		DeleteContext: resourceDynamicRiskDelete,
		UpdateContext: resourceDynamicRiskCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"threshold": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    3,
				Description: "Threshold settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetDyRiskNameTypes(), false),
							Description:  "Name of the risk settings, e.g: alert, mfa, lockUser",
						},

						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Set the value to true to enable the risk profile",
						},
						"value": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Threshold value",
						},
					},
				},
			},
		},
	}
}

func resourceDynamicRiskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	dyRiSk := &client.Dynamic_Risk{}

	res, err := prosimoClient.GetDYRisk(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, returneddyRisk := range res.DyRisk {
		dyRiSk.Id = returneddyRisk.Id
	}

	thresholds := []client.Threshold{}

	if v, ok := d.GetOk("threshold"); ok {
		thresholdList := v.([]interface{})
		log.Printf("[DEBUG] Updating threshold block")
		for _, thresHold := range thresholdList {
			val := thresHold.(map[string]interface{})
			thresholdData := client.Threshold{}
			if name, ok := val["name"].(string); ok {
				thresholdData.Name = name
			}
			if enabled, ok := val["enabled"].(bool); ok {
				thresholdData.Enabled = enabled
			}
			if value, ok := val["value"].(int); ok {
				thresholdData.Value = value
			}

			thresholds = append(thresholds, thresholdData)
		}

	}
	alert, mfa, lockUser := false, false, false
	for i := range thresholds {
		// log.Println("name=", thresholds[i].Name)
		if thresholds[i].Name == "alert" && thresholds[i].Value <= 100 {
			alert = true
		} else if thresholds[i].Name == "mfa" && thresholds[i].Value <= 100 {
			mfa = true
		} else if thresholds[i].Name == "lockUser" && thresholds[i].Value == 100 {
			lockUser = true
		} else {
			log.Printf("[DEBUG] invalid threshold name input")
		}
	}
	if alert == true && mfa == true && lockUser == true {
		dyRiSk.Thresholds = thresholds
		_, err := prosimoClient.PutDYRisk(ctx, dyRiSk)
		if err != nil {
			return diag.FromErr(err)
		}

	} else {
		log.Println("[DEBUG] Invalid threshold Inputs")
		return diag.Errorf("Invalid threshold Inputs")
	}
	d.SetId("Dynamic_Risk")
	return resourceDynamicRiskRead(ctx, d, meta)
}

func resourceDynamicRiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	res, err := prosimoClient.GetDYRisk(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	var dyRisk *client.Dynamic_Risk
	for _, returneddyRisk := range res.DyRisk {
		dyRisk = returneddyRisk
		break
	}
	if dyRisk == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get dyrisk id",
		})
	}

	thresholditems := flattenThresholdItemsData(&dyRisk.Thresholds)
	d.Set("id", dyRisk.Id)
	d.Set("threshold", thresholditems)

	d.SetId(dyRisk.Id)
	return diags
}

func flattenThresholdItemsData(ThresholdItems *[]client.Threshold) []interface{} {
	if ThresholdItems != nil {
		ois := make([]interface{}, len(*ThresholdItems), len(*ThresholdItems))

		for i, ThresholdItem := range *ThresholdItems {
			oi := make(map[string]interface{})

			oi["name"] = ThresholdItem.Name
			oi["enabled"] = ThresholdItem.Enabled
			oi["value"] = ThresholdItem.Value

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func resourceDynamicRiskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics
	var dyRiskid string
	res, err := prosimoClient.GetDYRisk(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, returneddyRisk := range res.DyRisk {
		dyRiskid = returneddyRisk.Id
	}
	// log.Println("printing dyriskid", dyRiskid)

	thresholds := []client.Threshold{}
	threshold1 := client.Threshold{
		Name:    "alert",
		Enabled: true,
		Value:   45,
	}
	thresholds = append(thresholds, threshold1)
	threshold2 := client.Threshold{
		Name:    "mfa",
		Enabled: true,
		Value:   45,
	}
	thresholds = append(thresholds, threshold2)
	threshold3 := client.Threshold{
		Name:    "lockUser",
		Enabled: false,
		Value:   100,
	}
	thresholds = append(thresholds, threshold3)
	dyrisk_input := &client.Dynamic_Risk{
		Id:         dyRiskid,
		Thresholds: thresholds,
	}
	res1, err := prosimoClient.PutDYRisk(ctx, dyrisk_input)
	_ = res1
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] set to default")
	d.SetId("")

	return diags
}
