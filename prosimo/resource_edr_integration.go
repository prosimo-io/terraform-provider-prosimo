package prosimo

import (
	"context"
	"fmt"
	"log"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEdrIntegration() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify end point security integrations.",
		CreateContext: resourceEdrProfileUpdate,
		UpdateContext: resourceEdrProfileUpdate,
		ReadContext:   resourceEdrProfileRead,
		DeleteContext: resourceEdrProfileDelete,
		// UpdateContext: resourceEdrProfileUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"crowdstrike": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Endpoint Security Integration name",
						},
						"vendor": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetEDRvendorTypes(), false),
							Description:  "Select EDR Vendor, for now only CrowdStrike is supported.",
						},
						"criteria": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sensor_active": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetEDRProfileInputTypes(), false),
										Description:  "Activate sensor, e.g: enabled, disabled",
									},
									"status": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetEDRProfileInputTypes(), false),
										Description:  "Status, e.g: enabled, disabled",
									},
									"zta_score": {
										Type:        schema.TypeSet,
										Required:    true,
										Description: "Zero Trust Access Score",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"from": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"to": {
													Type:     schema.TypeInt,
													Required: true,
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
		},
	}
}

func resourceEdrProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	edrProfileList := []client.EDR_Profile{}
	edrProfile := client.EDR_Profile{}
	criteria := client.CRITERIA{}
	ztaScore := client.ZtaScore{}

	if v, ok := d.GetOk("crowdstrike"); ok {
		profileList := v.([]interface{})
		if len(profileList) > 0 {
			for _, inprofile := range profileList {
				profile := inprofile.(map[string]interface{})
				// edrProfile.Id = ""
				if v, ok := profile["name"]; ok {
					name := v.(string)
					edrProfile.Name = name
				}

				if v, ok := profile["vendor"]; ok {
					vendor := v.(string)
					flag, err := prosimoClient.GetEDRVendor(ctx, vendor)
					if err != nil {
						return diag.FromErr(err)
					}
					if flag {
						edrProfile.Vendor = vendor
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "EDR config doesn't exist",
							Detail:   fmt.Sprintf("Unable to find edr config for vendor %s", vendor),
						})
						return diags
					}
				}

				if v, ok := profile["criteria"]; ok {
					criteriadetails := v.(*schema.Set).List()[0].(map[string]interface{})
					if v, ok := criteriadetails["sensor_active"]; ok {
						sensorActive := v.(string)
						criteria.SensorActive = sensorActive
					}
					if v, ok := criteriadetails["status"]; ok {
						status := v.(string)
						criteria.Status = status
					}
					if v, ok := criteriadetails["zta_score"].(*schema.Set); ok && v.Len() > 0 {
						ztaScores := v.List()[0].(map[string]interface{})
						if v, ok := ztaScores["from"]; ok {
							from := v.(int)
							ztaScore.From = from
						}
						if v, ok := ztaScores["to"]; ok {
							to := v.(int)
							ztaScore.To = to
						}
						criteria.Ztascore = ztaScore
					}
					edrProfile.Criteria = criteria
				}
				edrProfileList = append(edrProfileList, edrProfile)
			}
		}
	}

	existingval, err := prosimoClient.GetEDRProfile(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	existingProfileList := existingval.EdrProfileRes.CrowdStrike
	if len(existingProfileList) > 0 {
		for _, profile := range edrProfileList {
			isAddList := true
			for _, profile1 := range existingProfileList {
				if profile.Name == profile1.Name {
					isAddList = false
				} else {
					continue
				}
			}
			if isAddList == true {
				existingProfileList = append(existingProfileList, profile)
			}
		}
	} else {
		existingProfileList = edrProfileList
	}
	// log.Println("len(edrProfileList)", len(edrProfileList))
	if len(edrProfileList) > 0 {
		for i, profile := range existingProfileList {
			isDeleteList := true
			for _, profile1 := range edrProfileList {
				if profile.Name == profile1.Name {
					isDeleteList = false
				} else {
					continue
				}
			}
			if isDeleteList == true {
				existingProfileList = append(existingProfileList[:i], existingProfileList[i+1:]...)
			}
		}
	} else {
		// existingProfileList = existingProfileList[:0]
		existingProfileList = nil
	}
	log.Printf("[DEBUG] Creating EDR Profile : %v", existingProfileList)
	res, err := prosimoClient.UpdateEDRProfile(ctx, existingProfileList)
	if err != nil {
		log.Printf("[DEBUG] Error in creating EDRprofile")
		return diag.FromErr(err)
	}
	d.SetId(res.EdrProfileUpdateRes.AuditID)
	return resourceEdrProfileRead(ctx, d, meta)
}

func resourceEdrProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	res, err := prosimoClient.GetEDRProfile(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	crowdstrikeitems := flattenCrowdStrikeItemsData(&res.EdrProfileRes.CrowdStrike)
	d.Set("crowdstrike", crowdstrikeitems)

	return diags
}

func flattenCrowdStrikeItemsData(Profiles *[]client.EDR_Profile) []interface{} {
	if Profiles != nil {
		ois := make([]interface{}, len(*Profiles), len(*Profiles))

		for i, Profile := range *Profiles {
			oi := make(map[string]interface{})

			oi["name"] = Profile.Name
			oi["id"] = Profile.Id
			oi["vendor"] = Profile.Vendor
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func resourceEdrProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	existingProfileList := []client.EDR_Profile{}
	existingProfileList = nil
	_, err := prosimoClient.UpdateEDRProfile(ctx, existingProfileList)
	if err != nil {
		log.Printf("[DEBUG] Error in creating EDRprofile")
		return diag.FromErr(err)
	}
	return resourceEdrProfileRead(ctx, d, meta)

}
