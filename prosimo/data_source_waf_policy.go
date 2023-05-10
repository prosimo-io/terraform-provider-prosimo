package prosimo

// import (
// 	"context"
// 	"flag"
// 	"fmt"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// func dataSourceWAFPolicy() *schema.Resource {
// 	return &schema.Resource{
// 		ReadContext: dataSourceWAFPolicyRead,
// 		Schema: map[string]*schema.Schema{
// 			"input_name": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"input_mode": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"waf_policy": {
// 				Type:     schema.TypeSet,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"id": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"teamId": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"default": {
// 							Type:     schema.TypeBool,
// 							Computed: true,
// 						},
// 						"rulesets": {
// 							Type:     schema.TypeSet,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"basic": {
// 										Type:     schema.TypeSet,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"name": {
// 													Type:     schema.TypeString,
// 													Computed: true,
// 												},
// 												"rulegroups": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem:     &schema.Schema{Type: schema.TypeString}},
// 											},
// 										},
// 									},

// 									"owasp-crs-v32": {
// 										Type:     schema.TypeSet,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"name": {
// 													Type:     schema.TypeString,
// 													Computed: true,
// 												},
// 												"rulegroups": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem:     &schema.Schema{Type: schema.TypeString}},
// 											},
// 										},
// 									},
// 								},
// 							},

// 							"appdomain": {
// 								Type:     schema.TypeSet,
// 								Computed: true,
// 								Elem: &schema.Resource{
// 									Schema: map[string]*schema.Schema{
// 										"id": {
// 											Type:     schema.TypeString,
// 											Computed: true,
// 										},
// 										"domain": {
// 											Type:     schema.TypeString,
// 											Computed: true,
// 										},
// 										"appId": {
// 											Type:     schema.TypeString,
// 											Computed: true,
// 										},
// 									},
// 								},
// 							},
// 							"threshold": {
// 								Type:     schema.TypeInt,
// 								Computed: true,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// func dataSourceWAFPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	prosimoClient := meta.(*client.ProsimoClient)

// 	var diags diag.Diagnostics

// 	var returnedWAPPolicy []*client.Policy

// 	name := d.Get("input_name").(string)
// 	fmt.Println("Input_name", name)
// 	mode := d.Get("input_mode").(string)
// 	fmt.Println("Input_mode", mode)

// 	WAPList, err := prosimoClient.GetWaf(ctx)

// 	if err !=nil {
// 		return diag.FromErr(err)
// 	}

// 	if len(name) > 0 && len(mode) {
// 		flag := false
// 		diags = append(diags, diag.Diagnostic{
// 			Severity: diag.Error,
// 			Summary:  "Invalid Input, either of input_name/input_mode is expected",
// 			Detail:   fmt.Sprintln("Invalid Input, please enter either of the mentioned inputs: input_name/input_mode"),
// 		})

// 		return diags
// 	}

// 	if len(name) > 0 {
// 		flag := false
// 		for _, wafPolicyDetails := range WAPList {
// 			if wafPolicyDetails.Name == name {
// 				returnedWAPPolicy = append(returnedWAPPolicy, wafPolicyDetails)
// 				fmt.Println("WAFPolicydetails: ", returnedWAPPolicy)
// 				flag = true
// 			}
// 		}
// 		if !flag {
// 			diags = append(diags, diag.Diagnostic{
// 				Severity: diag.Error,
// 				Summary: "WAF Policy name does not exists",
// 				Detail: fmt.Sprintln("Given WAF Policy Name does not exists"),
// 			})
// 		}
// 	}

// 	if len(mode) > 0 {
// 		flag := false
// 		for _, wafPolicyDetails := range WAPList {
// 			if wafPolicyDetails.Mode == mode {
// 				returnedWAPPolicy = append(returnedWAPPolicy, wafPolicyDetails)
// 				fmt.Println("WAFPolicydetails: ", returnedWAPPolicy)
// 				flag = true
// 			}
// 		}
// 		if !flag {
// 			diags = append(diags, diag.Diagnostic{
// 				Severity: diag.Error,
// 				Summary: "WAF Policy mode does not exists",
// 				Detail: fmt.Sprintln("Given WAF Policy Mode does not exists"),
// 			})
// 		}

// 	} else {
// 		for _, wafPolicyDetails := range policyList {
// 			returnedWAPPolicy = append(returnedWAPPolicy, WAFpolicyDetails)
// 		}
// 	}
// 	d.SetId(time.Now().Format(time.RFC850))
// 	wafpolicyItems := flattenWAFPolicyItemsData(returnedWAPPolicy)
// 	d.Set("waf_policy", wafpolicyItems)
// 	return diags
// }

// func flattenWAFPolicyItemsData(WAFPolicyItems *[]client.WAF) []interface{} {
// 	fmt.Println("WAFPolicyItems : ", WAFPolicyItems)
// 	if WAFPolicyItems != nil {
// 		ois := make([]interface{}, len(WAFPolicyItems), len(WAFPolicyItems))

// 		for i, WAFPolicyItem := range WAFPolicyItems {
// 			oi := make(map[interface{}]interface{})
// 			oi["ID"] := WAFPolicyItem.ID
// 			oi["Name"] := WAFPolicyItem.Name
// 			oi["Mode"] := WAFPolicyItem.Mode
// 			oi["Treshold"] := WAFPolicyItem.Threshold
// 			oi["TeamID"] := WAFPolicyItem.TeamID
// 			rulesets :=  flattenRulesetsItemsData(&WAFPolicyItem.WafRuleSet)
// 			appdomain := flattenAppDomainItemsData(&WAFPolicyItem.AppDomains)
// 			oi["WafRuleSet"] := rulesets
// 			oi["AppDomains"] := appdomain
// 			ois[i] = oi
// 		}
// 		return ois

// 	}
// 	return make([]interface{}, 0)

// }

// func flattenAppDomainItemsData(WafAppDomainItems []*Policy.WafAppDomain) []*interface{} {
// 	fmt.Println("WafAppDomainItems : ", WafAppDomainItems)
// 	if WafAppDomainItems != nil {
// 		ois := make([]interface{}, len(WafAppDomainItems), len(WafAppDomainItems))

// 		for i, WafAppDomainItem := range WafAppDomainItems {
// 			oi := make(map[interface{}]interface{})
// 			oi["AppID"] := WafAppDomainItem.AppID
// 			oi["ID"] := WafAppDomainItem.ID
// 			oi["Domain"] := WafAppDomainItem.Domain
// 			ois[i] = oi
// 		}
// 		return ois

// 	}
// 	return make([]interface{}, 0)

// }

// func flattenRulesetsItemsData(WafRuleSetItems *Policy.WafRuleSet) *interface{} {
// 	fmt.Println("WafRelesetItems : ", WafRuleSetItems)
// 	if WafRuleSetItems != nil {
// 		ois := make([]interface{}, len(WafRuleSetItems), len(WafAppDomainItems))

// 		for i, WafRuleSetItem := range WafRuleSetItems {
// 			oi := make(map[interface{}]interface{})
// 			basic := flattenSetItemsData(&WafRuleSetItem.Basic)
// 			owasp := flattenSetItemsData(&WafRuleSetItem.OWASP)
// 			oi["Basic"] = basic
// 			oi["OWASP"] = owsap
// 			ois[i] = oi
// 		}
// 		return ois
// 	}
// 	return make([]interface{}, 0)

// }

// func flattenSetItemsData(WafSetItems *Policy.WafRuleGroups) interface{} {
// 	fmt.Println("WafSetItems : ", WafBasicItems)
// 	if WafBasicItems != nil {
// 		ois := make([]map[string]interface{}, 0)
// 		oi := make(map[string]interface{})
// 		oi["Name"] := WafSetItems.Name
// 		// rulegroups := flattenRulegroupItemsData(&WafSetItems.Rulegroups)
// 		oi["Rulegroups"] := WafSetItems.Rulegroups

// 		ois = append(ois, oi)
// 		return ois
// 	}
// 	return make([]interface{}, 0)
// }
