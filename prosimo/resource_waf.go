package prosimo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWAF() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify waf settings.",
		CreateContext: resourceWafCreate,
		UpdateContext: resourceWafUpdate,
		ReadContext:   resourceWafRead,
		DeleteContext: resourceWafDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"waf_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the waf policy set",
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(client.GetWafModeTypes(), false),
				Description:  "waf detect mode, e.g: enforce, detect ",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "When the Anomaly Score exceeds the Anomaly Threshold, a notification is generated by prosimo fabric in Detect mode or along with the notification the request is blocked in Enforce mode",
			},
			"rulesets": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_groups": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "List of basic rules sets waf would apply",
									},
								},
							},
						},
						"owasp_crs_v32": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_groups": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "OWASP Modsecurity Core Ruleset v3.2  waf would apply",
									},
								},
							},
						},
					},
				},
			},
			"app_domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of app domain to which waf would apply",
			},
		},
	}
}

func resourceWafCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	waf := &client.Waf{}

	if v, ok := d.GetOk("waf_name"); ok {
		wafName := v.(string)
		waf.Name = wafName
	}

	if v, ok := d.GetOk("mode"); ok {
		mode := v.(string)
		waf.Mode = mode
	}

	if v, ok := d.GetOk("threshold"); ok {
		threshold := v.(int)
		waf.Threshold = threshold
	}

	if v, ok := d.GetOk("rulesets"); ok {
		rulesets := v.([]interface{})
		if len(rulesets) > 0 && rulesets[0] != nil {

			wafRulesetsTF := rulesets[0].(map[string]interface{})

			wafRuleSet := &client.WafRuleSet{}
			waf.WafRuleSet = wafRuleSet

			if v, ok := wafRulesetsTF["basic"].(*schema.Set); ok && v.Len() > 0 {

				basicRuleSet := v.List()[0].(map[string]interface{})

				if v, ok := basicRuleSet["rule_groups"]; ok {
					basicRuleGroups := v.([]interface{})
					if len(basicRuleGroups) > 0 {
						wafBasicRuleGroups := &client.WafRuleGroups{}
						wafBasicRuleGroups.Name = client.WafRuleSetBasic
						wafBasicRuleGroups.Rulegroups = expandStringList(v.([]interface{}))
						wafRuleSet.Basic = wafBasicRuleGroups
					}
				}

			}

			if v, ok := wafRulesetsTF["owasp_crs_v32"].(*schema.Set); ok && v.Len() > 0 {

				owaspRuleSet := v.List()[0].(map[string]interface{})

				if v, ok := owaspRuleSet["rule_groups"]; ok {
					owaspRuleGroups := v.([]interface{})
					if len(owaspRuleGroups) > 0 {
						wafOWASPRuleGroups := &client.WafRuleGroups{}
						wafOWASPRuleGroups.Name = client.WafRuleSetOWASP
						wafOWASPRuleGroups.Rulegroups = expandStringList(v.([]interface{}))
						wafRuleSet.OWASP = wafOWASPRuleGroups
					}
				}
			}
		}

	}

	if waf.WafRuleSet != nil {
		err := validateWafRuleSet(ctx, d, meta, waf.WafRuleSet)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var wafAppDomainIds *client.WafAppDomainIds
	if v, ok := d.GetOk("app_domains"); ok {

		appDomainsList := v.([]interface{})
		if len(appDomainsList) > 0 {
			wafAppDomainIds = &client.WafAppDomainIds{}

			appDomainNames := expandStringList(v.([]interface{}))

			existingAppDomainList, err := prosimoClient.GetAppDomains(ctx)
			if err != nil {
				return diag.FromErr(err)
			}

			var addDomainIDs []string
			for _, appDomain := range appDomainNames {
				appDomainExists := false
				for _, exisingAppDomain := range existingAppDomainList {
					if appDomain == exisingAppDomain.Domain {
						appDomainExists = true
						addDomainIDs = append(addDomainIDs, exisingAppDomain.ID)
					}
				}

				if !appDomainExists {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid App Domain",
						Detail:   fmt.Sprintf("Invalid App Domain %s", appDomain),
					})
				}

			}

			wafAppDomainIds.AddDomainIDs = addDomainIDs

		}

	}
	log.Printf("[DEBUG] Creating Waf: %s", waf.Name)

	waf.DefaultWaf = false
	createdWAF, err := prosimoClient.CreateWaf(ctx, waf)
	if err != nil {
		return diag.FromErr(err)
	}

	if wafAppDomainIds != nil {
		log.Printf("[DEBUG] Creating Waf App Domains: %v", wafAppDomainIds)
		_, err := prosimoClient.UpdateWafAppDomains(ctx, wafAppDomainIds, createdWAF.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] Created Waf App Domains: %v for waf %s", wafAppDomainIds, waf.Name)
	}

	log.Printf("[DEBUG] Created WAF: %v", waf)
	d.SetId(createdWAF.ID)

	return resourceWafRead(ctx, d, meta)
}

func resourceWafUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	wafID := d.Id()

	prosimoClient := meta.(*client.ProsimoClient)

	waf := &client.Waf{}
	waf.ID = wafID

	if v, ok := d.GetOk("waf_name"); ok {
		wafName := v.(string)
		waf.Name = wafName
	}

	updateReq := false

	if d.HasChange("mode") && !d.IsNewResource() {
		updateReq = true
	}
	if v, ok := d.GetOk("mode"); ok {
		mode := v.(string)
		waf.Mode = mode
	}

	if d.HasChange("threshold") && !d.IsNewResource() {
		updateReq = true
	}
	if v, ok := d.GetOk("threshold"); ok {
		threshold := v.(int)
		waf.Threshold = threshold
	}

	if v, ok := d.GetOk("rulesets"); ok {
		rulesets := v.([]interface{})

		if len(rulesets) > 0 {
			updateReq = true
			wafRulesetsTF := rulesets[0].(map[string]interface{})

			wafRuleSet := &client.WafRuleSet{}
			waf.WafRuleSet = wafRuleSet

			if v, ok := wafRulesetsTF["basic"].(*schema.Set); ok && v.Len() > 0 {
				basicRuleSet := v.List()[0].(map[string]interface{})

				if v, ok := basicRuleSet["rule_groups"]; ok {
					basicRuleGroups := v.([]interface{})
					wafBasicRuleGroups := &client.WafRuleGroups{}
					wafRuleSet.Basic = wafBasicRuleGroups
					if len(basicRuleGroups) > 0 {
						wafBasicRuleGroups.Name = client.WafRuleSetBasic
						wafBasicRuleGroups.Rulegroups = expandStringList(v.([]interface{}))
					}
				}

			}

			// if d.HasChange("rulesets.0.owasp_crs_v32") && !d.IsNewResource() {
			// 	updateReq = true
			// }

			if v, ok := wafRulesetsTF["owasp_crs_v32"].(*schema.Set); ok && v.Len() > 0 {

				owaspRuleSet := v.List()[0].(map[string]interface{})

				if v, ok := owaspRuleSet["rule_groups"]; ok {
					owaspRuleGroups := v.([]interface{})
					wafOWASPRuleGroups := &client.WafRuleGroups{}
					wafRuleSet.OWASP = wafOWASPRuleGroups
					if len(owaspRuleGroups) > 0 {
						wafOWASPRuleGroups.Name = client.WafRuleSetOWASP
						wafOWASPRuleGroups.Rulegroups = expandStringList(v.([]interface{}))
					}
				}
			}
		}

	}

	if waf.WafRuleSet != nil {
		err := validateWafRuleSet(ctx, d, meta, waf.WafRuleSet)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	createdWaf, err := prosimoClient.GetWafByID(ctx, wafID)
	if err != nil {
		return diag.FromErr(err)
	}
	createdAppDomains := make(map[string]bool)
	if len(createdWaf.AppDomains) > 0 {
		for _, appDomain := range createdWaf.AppDomains {
			log.Printf("[DEBUG] createdWaf.AppDomains for - %v", appDomain)
			createdAppDomains[appDomain.Domain] = false
		}
	}

	var wafAppDomainIds *client.WafAppDomainIds
	if v, ok := d.GetOk("app_domains"); ok {

		appDomainsList := v.([]interface{})

		if len(appDomainsList) > 0 {

			if d.HasChange("app_domains") && !d.IsNewResource() {
				updateReq = true
			}

			wafAppDomainIds = &client.WafAppDomainIds{}

			appDomainNames := expandStringList(v.([]interface{}))

			existingAppDomainList, err := prosimoClient.GetAppDomains(ctx)
			if err != nil {
				return diag.FromErr(err)
			}

			var addDomainIDs []string
			for _, appDomain := range appDomainNames {
				appDomainExists := false
				for _, exisingAppDomain := range existingAppDomainList {
					if appDomain == exisingAppDomain.Domain {
						appDomainExists = true
						addDomainIDs = append(addDomainIDs, exisingAppDomain.ID)
						createdAppDomains[exisingAppDomain.Domain] = true
					}
				}

				if !appDomainExists {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid App Domain",
						Detail:   fmt.Sprintf("Invalid App Domain %s", appDomain),
					})
				}

			}

			var deleteDomainIDs []string
			for domain, isAvailable := range createdAppDomains {

				if !isAvailable {

					for _, exisingAppDomain := range existingAppDomainList {
						if domain == exisingAppDomain.Domain {
							deleteDomainIDs = append(deleteDomainIDs, exisingAppDomain.ID)
						}
					}

				}
			}
			wafAppDomainIds.AddDomainIDs = addDomainIDs
			wafAppDomainIds.DeleteDomainIDs = deleteDomainIDs

		}

	} else if len(createdWaf.AppDomains) > 0 {
		log.Printf("[DEBUG] delete appDomainsList")
		wafAppDomainIds = &client.WafAppDomainIds{}
		var deleteDomainIDs []string

		for _, appDomain := range createdWaf.AppDomains {
			deleteDomainIDs = append(deleteDomainIDs, appDomain.ID)
		}
		wafAppDomainIds.DeleteDomainIDs = deleteDomainIDs
		updateReq = true
	}

	log.Printf("[DEBUG] Updating waf - %s", waf.Name)
	if updateReq {
		updatedWAF, err := prosimoClient.UpdateWaf(ctx, waf)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(updatedWAF.ID)
		log.Printf("[DEBUG] Updated WAF for - %v", waf)

		if wafAppDomainIds != nil {
			log.Printf("[DEBUG] Creating Waf App Domains: %v", wafAppDomainIds)
			_, err := prosimoClient.UpdateWafAppDomains(ctx, wafAppDomainIds, updatedWAF.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] Created Waf App Domains: %v for waf %s", wafAppDomainIds, waf.Name)
		}
	}

	return resourceWafRead(ctx, d, meta)
}

func resourceWafRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	wafID := d.Id()

	log.Printf("Get WAF for %s", wafID)

	waf, err := prosimoClient.GetWafByID(ctx, wafID)
	if err != nil {
		return diag.FromErr(err)
	}

	if waf == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get WAF",
			Detail:   fmt.Sprintf("[ERROR] Unable to find WAF for ID %s", wafID),
		})

		return diags
	}

	d.Set("waf_name", waf.Name)
	d.Set("mode", waf.Mode)
	d.Set("id", waf.ID)
	d.Set("threshold", waf.Threshold)

	if waf.WafRuleSet != nil {
		rulesets := map[string]interface{}{}

		if waf.WafRuleSet.Basic != nil {
			basicRulesets := map[string]interface{}{}
			basicRulesets["rule_groups"] = flattenStringList(waf.WafRuleSet.Basic.Rulegroups)
			rulesets["basic"] = basicRulesets
		}

		if waf.WafRuleSet.OWASP != nil {
			owaspRulesets := map[string]interface{}{}
			owaspRulesets["rule_groups"] = flattenStringList(waf.WafRuleSet.OWASP.Rulegroups)
			rulesets["owasp_crs_v32"] = owaspRulesets
		}

		d.Set("rulesets", rulesets)

	}

	log.Printf("[DEBUG] Waf for %s - %+v", waf.Name, waf)

	return diags
}

func resourceWafDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	ipPoolID := d.Id()

	err := prosimoClient.DeleteWaf(ctx, ipPoolID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func validateWafRuleSet(ctx context.Context, d *schema.ResourceData, meta interface{}, wafeRuleSet *client.WafRuleSet) error {

	prosimoClient := meta.(*client.ProsimoClient)

	prosimoWafRuleSet, err := prosimoClient.GetWafRuleSet(ctx)
	if err != nil {
		return err
	}

	basicWafRuleGroups := prosimoWafRuleSet.Basic
	basicRuleGroupsList := make(map[string]bool)
	for _, ruleGroup := range basicWafRuleGroups.Rulegroups {
		basicRuleGroupsList[ruleGroup] = true
	}

	owaspWafRuleGroups := prosimoWafRuleSet.OWASP
	owaspRuleGroupsList := make(map[string]bool)
	for _, ruleGroup := range owaspWafRuleGroups.Rulegroups {
		owaspRuleGroupsList[ruleGroup] = true
	}

	var invalidRuleGroups []string
	if wafeRuleSet.Basic != nil {
		userBasicWafRuleGroups := wafeRuleSet.Basic.Rulegroups
		for _, ruleGroup := range userBasicWafRuleGroups {
			if _, ok := basicRuleGroupsList[ruleGroup]; !ok {
				invalidRuleGroups = append(invalidRuleGroups, ruleGroup)
			}
		}
	}

	if wafeRuleSet.OWASP != nil {
		userOWASPWafRuleGroups := wafeRuleSet.OWASP.Rulegroups
		for _, ruleGroup := range userOWASPWafRuleGroups {
			if _, ok := owaspRuleGroupsList[ruleGroup]; !ok {
				invalidRuleGroups = append(invalidRuleGroups, ruleGroup)
			}
		}
	}

	if len(invalidRuleGroups) > 0 {
		errorText := fmt.Sprintf("invalid rule groups - %s", strings.Join(invalidRuleGroups, ", "))
		return errors.New(errorText)
	}

	return nil

}
