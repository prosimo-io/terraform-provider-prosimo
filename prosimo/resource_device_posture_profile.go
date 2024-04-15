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

func resourceDPProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify device posture profiles.",
		CreateContext: resourceDPProfileUpdate,
		UpdateContext: resourceDPProfileUpdate,
		ReadContext:   resourceDPProfileRead,
		DeleteContext: resourceDPProfileDelete,
		// UpdateContext: resourceEdrProfileUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"inprofile_list": {
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
							Description: "Profile name",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Set TRUE if you want to enable the profile, else false",
						},
						"risk_level": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetDPRiskLevelTypes(), false),
							Description:  "Risk level of the profile, e.g: high, medium, low",
						},
						"edr_profiles": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Select end point security profiles from the list. If the list is empty, profile needs to created first",
						},
						"criteria": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"os": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetOsTypes(), false),
										Description:  "OS type, e.g: windows, mac",
									},
									"firewall_status": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "na",
										ValidateFunc: validation.StringInSlice(client.GetEDRProfileInputTypes(), false),
										Description:  "Set firewall status, e.g: enabled, disabled",
									},
									"disk_encryption_status": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "na",
										ValidateFunc: validation.StringInSlice(client.GetEDRProfileInputTypes(), false),
										Description:  "Set disk encryption status, e.g: enabled, disabled",
									},
									"domain_of_interest": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The user's device should match one of the domains selected from the dropdown",
									},
									"running_process": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Select the processes that should be running. If more than one process is added, all of them should be running",
									},
									"os_operator": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(client.GetOsOperatorTypes(), false),
										Description:  "select one of is, is-notm is-atleast",
									},
									"os_version": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Os version details",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"build": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "OS Build Number, e.g: 11(21H2) 22000, 10(1511) 10586 etc",
												},
												"patch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Patch details, e.g: 1540, 1358 etc",
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

func resourceDPProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	dpProfileListHigh := []client.DevicePosture_Profile{}
	dpProfileListMedium := []client.DevicePosture_Profile{}
	dpProfileListLow := []client.DevicePosture_Profile{}
	dpProfile := client.DevicePosture_Profile{}
	criteria := client.CRITERIADP{}

	if v, ok := d.GetOk("inprofile_list"); ok {
		profileList := v.([]interface{})
		if len(profileList) > 0 {
			for _, inprofile := range profileList {
				profile := inprofile.(map[string]interface{})

				if v, ok := profile["name"]; ok {
					name := v.(string)
					dpProfile.Name = name
				}
				if v, ok := profile["enabled"]; ok {
					profile := v.(bool)
					dpProfile.Enabled = profile
				}
				if v, ok := profile["risk_level"]; ok {
					riskLevel := v.(string)
					dpProfile.RiskLevel = riskLevel
				}

				if v, ok := profile["criteria"]; ok {
					criteriadetails := v.(*schema.Set).List()[0].(map[string]interface{})
					if v, ok := criteriadetails["os"]; ok {
						os := v.(string)
						criteria.Os = os
					}
					if v, ok := criteriadetails["firewall_status"]; ok {
						firewallStatus := v.(string)
						criteria.FirewallStatus = firewallStatus
					}

					if v, ok := criteriadetails["disk_encryption_status"]; ok {
						diskEncryptionStatus := v.(string)
						criteria.DiskEncryptionStatus = diskEncryptionStatus
					}

					if v, ok := criteriadetails["domain_of_interest"]; ok {
						domainList := v.([]interface{})

						if len(domainList) > 0 {
							criteria.DomainOfInterest = expandStringList(v.([]interface{}))
						}
					}
					if v, ok := criteriadetails["running_process"]; ok {
						processes := v.([]interface{})

						if len(processes) > 0 {
							criteria.RunningProcess = expandStringList(v.([]interface{}))
						}
					}
					if v, ok := criteriadetails["os_operator"]; ok {
						osOperator := v.(string)
						criteria.OsOperator = osOperator
					}

					if v, ok := criteriadetails["os_version"]; ok {
						osVersonDetails := v.([]interface{})
						osVersion := client.Os_Version{}
						osVersionList := []client.Os_Version{}
						if len(osVersonDetails) > 0 {
							for _, inosVersion := range osVersonDetails {
								osversion := inosVersion.(map[string]interface{})
								if v, ok := osversion["build"]; ok {
									build := v.(string)
									osVersion.Build = build
								}
								if v, ok := osversion["patch"]; ok {
									patch := v.(string)
									osVersion.Patch = patch
								}
								osVersionList = append(osVersionList, osVersion)
							}
						}
						criteria.OsVersions = osVersionList
					}
					dpProfile.Criteria = criteria
				}

				if v, ok := profile["edr_profiles"]; ok {
					edrprofileList := v.([]interface{})
					edrProfileList := []client.EDR_Profile{}
					edrProfile := client.EDR_Profile{}
					if len(edrprofileList) > 0 {
						edrprofileListFormated := expandStringList(v.([]interface{}))
						for _, edrprofilename := range edrprofileListFormated {
							status, matchingedr, err := prosimoClient.GetEDRProfileByName(ctx, edrprofilename)
							// log.Println("matchingedr", matchingedr)
							if err != nil {
								return diag.FromErr(err)
							}
							if status {
								// log.Println("matchingedr", *matchingedr)
								edrProfile = *matchingedr
								edrProfileList = append(edrProfileList, edrProfile)
							} else {
								diags = append(diags, diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "EDR profile doesn't exist",
									Detail:   fmt.Sprintf("Unable to find input edr profile."),
								})
								return diags
							}
						}
					}
					dpProfile.EdrProfiles = edrProfileList
				}
				if dpProfile.RiskLevel == client.RiskLevelHigh {
					dpProfileListHigh = append(dpProfileListHigh, dpProfile)
				} else if dpProfile.RiskLevel == client.RiskLevelMedium {
					dpProfileListMedium = append(dpProfileListMedium, dpProfile)
				} else {
					dpProfileListLow = append(dpProfileListLow, dpProfile)
				}
			}
			if len(dpProfileListHigh) >= 0 {
				existingdplist, err := prosimoClient.GetDPProfileBYRiskLevel(ctx, client.RiskLevelHigh)
				if err != nil {
					return diag.FromErr(err)
				}
				inputProfileListHigh, err := dp_profile_match(dpProfileListHigh, existingdplist.DPProfileRes)
				if err != nil {
					return diag.FromErr(err)
				}
				res, err := prosimoClient.UpdateDPProfileHigh(ctx, inputProfileListHigh)
				if err != nil {
					log.Printf("[DEBUG] Error in creating DPprofile")
					return diag.FromErr(err)
				}
				_ = res
			}

			if len(dpProfileListMedium) >= 0 {
				existingdplist, err := prosimoClient.GetDPProfileBYRiskLevel(ctx, client.RiskLevelMedium)
				if err != nil {
					return diag.FromErr(err)
				}
				inputProfileListMedium, err := dp_profile_match(dpProfileListMedium, existingdplist.DPProfileRes)
				if err != nil {
					return diag.FromErr(err)
				}
				res, err := prosimoClient.UpdateDPProfileMedium(ctx, inputProfileListMedium)
				if err != nil {
					log.Printf("[DEBUG] Error in creating DPprofile")
					return diag.FromErr(err)
				}
				_ = res
			}

			if len(dpProfileListLow) >= 0 {
				existingdplist, err := prosimoClient.GetDPProfileBYRiskLevel(ctx, client.RiskLevelLow)
				if err != nil {
					return diag.FromErr(err)
				}
				inputProfileListLow, err := dp_profile_match(dpProfileListLow, existingdplist.DPProfileRes)
				if err != nil {
					return diag.FromErr(err)
				}
				res, err := prosimoClient.UpdateDPProfilelow(ctx, inputProfileListLow)
				if err != nil {
					log.Printf("[DEBUG] Error in creating DPprofile")
					return diag.FromErr(err)
				}
				_ = res
			}
		}
	}
	d.SetId("DP_Profile_List")
	return resourceDPProfileRead(ctx, d, meta)
}

func dp_profile_match(dpProfileList []client.DevicePosture_Profile, existingdplist []client.DevicePosture_Profile) ([]client.DevicePosture_Profile, error) {
	if len(existingdplist) > 0 {
		for _, profile := range dpProfileList {
			isAddList := true
			for _, profile1 := range existingdplist {
				if profile.Name == profile1.Name {
					isAddList = false
				} else {
					continue
				}
			}
			if isAddList {
				existingdplist = append(existingdplist, profile)
			}
		}
	} else {
		existingdplist = dpProfileList
	}
	if len(dpProfileList) > 0 {
		for i, profile := range existingdplist {
			isDeleteList := true
			for _, profile1 := range dpProfileList {
				if profile.Name == profile1.Name {
					isDeleteList = false
				} else {
					continue
				}
			}
			if isDeleteList {
				existingdplist = append(existingdplist[:i], existingdplist[i+1:]...)
			}
		}
	} else {
		existingdplist = nil
	}
	return existingdplist, nil
}

func resourceDPProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	res, err := prosimoClient.GetDPProfile(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	dprisklevelhigh := flattenDPItemsData(&res.High)
	d.Set("high", dprisklevelhigh)
	dprisklevelmedium := flattenDPItemsData(&res.Medium)
	d.Set("medium", dprisklevelmedium)
	dprisklevellow := flattenDPItemsData(&res.Low)
	d.Set("low", dprisklevellow)
	return diags
}

func flattenDPItemsData(Profiles *[]client.DevicePosture_Profile) []interface{} {
	if Profiles != nil {
		ois := make([]interface{}, len(*Profiles), len(*Profiles))

		for i, Profile := range *Profiles {
			oi := make(map[string]interface{})

			oi["name"] = Profile.Name
			oi["id"] = Profile.Id
			oi["enabled"] = Profile.Enabled
			oi["risk_level"] = Profile.RiskLevel
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func resourceDPProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	inDPProfileList := []client.DevicePosture_Profile{}
	inDPProfileList = nil

	res0, err := prosimoClient.UpdateDPProfileHigh(ctx, inDPProfileList)
	if err != nil {
		log.Printf("[DEBUG] Error in updating EDRprofile")
		return diag.FromErr(err)
	}
	_ = res0
	res1, err := prosimoClient.UpdateDPProfileMedium(ctx, inDPProfileList)
	if err != nil {
		log.Printf("[DEBUG] Error in updating EDRprofile")
		return diag.FromErr(err)
	}
	_ = res1
	res2, err := prosimoClient.UpdateDPProfilelow(ctx, inDPProfileList)
	if err != nil {
		log.Printf("[DEBUG] Error in updating EDRprofile")
		return diag.FromErr(err)
	}
	_ = res2
	return resourceDPProfileRead(ctx, d, meta)

}
