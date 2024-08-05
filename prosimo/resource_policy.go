package prosimo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "The policy engine in Prosimo helps control access rules between users, applications, and networking. Use this resource to create/modify policies.",
		CreateContext: resourcePolicyCreate,
		UpdateContext: resourcePolicyUpdate,
		ReadContext:   resourcePolicyRead,
		DeleteContext: resourcePolicynDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of policy",
			},
			"types": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "type of policy, e.g: default, managed",
			},
			"teamid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"device_posture_configured": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "only applicable for access app access type, set it to true to enable device posture",
			},

			"app_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(client.GetPolicyappAccessType(), false),
				Description:  "app access type, e.g: access, transit",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy Namespace, only applicable for transit app_access_type",
			},
			"details": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeString,
							Required: true,
							// Elem:     &schema.Schema{Type: schema.TypeString},
							ValidateFunc: validation.StringInSlice(client.GetPolicyActionTypes(), false),
							Description:  "policy action, e.g: allow, deny",
						},
						"lock_users": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "set this to true to lock the user defined in policy",
						},
						"alert": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "set this to true to trigger the alert as per policy config",
						},
						"mfa": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "set this to true to trigger",
						},
						"skipwaf": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "set this to true to skip waf",
						},
						"bypass": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "set this to true to bypass policy",
						},
						"matches": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_entries": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(client.GetPolicyResourceTypes(), false),
													Description:  "Select policy match condition type, for access policy options are users, location, idp, devices, time, url, device-posture, fqdn and advanced. For transit type options are time, url, networkacl, fqdn, egressfqdns, prosimonetworks, networks and advanced",
												},
												"property": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Select property of selected type",
												},
												"operation": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(client.GetPolicyResourceOperation(), false),
													Description:  "Operation of the selected property, available options are Id, Is NOT, Contains, Does NOT contain, Starts with, Ends with, In, NOT in, Is at least, Between ",
												},
												"values": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"inputitems": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Input value name",
																		},
																		"ip_details": {
																			Type:        schema.TypeSet,
																			Optional:    true,
																			Description: "Only applicable for type networkacl",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"source_ip": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "Source IP list",
																					},
																					"target_ip": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "Target IP list",
																					},
																					"protocol": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "List of protocols",
																					},
																					"source_port": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "Source port list",
																					},
																					"target_port": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "Target port list",
																					},
																				},
																			},
																		},
																		"egress_fqdn_details": {
																			Type:        schema.TypeSet,
																			Optional:    true,
																			Description: "Only applicable for type egressfqdn",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"fqdn_inverse_match": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "FQDNs which need to be excluded",
																					},
																					"fqdn_match": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "FQDNs which need to be included",
																					},
																					"protocol": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "List of protocols",
																					},
																					"source_ip": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "Source IP list",
																					},
																					"target_port": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						Elem:        &schema.Schema{Type: schema.TypeString},
																						Description: "Target port list",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"selecteditems": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Selected value name",
																		},
																		"country_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Country name, only applicable for type location",
																		},
																		"state_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "State name, only applicable for type location",
																		},
																		"city_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "City name, only applicable for type location",
																		},
																	},
																},
															},
															"selectedgroups": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Input value name",
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
							},
						},
						"apps": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "App details to attach to the policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selecteditems": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the app",
												},
											},
										},
									},
								},
							},
						},
						"networks": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Network details to attach to the policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selecteditems": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the network",
												},
											},
										},
									},
								},
							},
						},
						"internet_traffic_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "set it to true to enable internet access",
						},
					},
				},
			},
			"createdtime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updatedtime": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func inputDataops(ctx context.Context, d *schema.ResourceData, meta interface{}) (diag.Diagnostics, client.Policy) {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	newpolicy := &client.Policy{}
	newdetails := client.Details{}

	newmatchList := client.MatchDetailList{}

	newvalues1 := client.Values{}
	var appaccestype string
	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		newpolicy.Name = name
	}
	if v, ok := d.GetOk("type"); ok {
		types := v.(string)
		newpolicy.PolicyType = types
	}
	if v, ok := d.GetOk("teamid"); ok {
		teamid := v.(string)
		newpolicy.TeamID = teamid
	}
	if v, ok := d.GetOk("app_access_type"); ok {
		appaccestype = v.(string)
		newpolicy.App_Access_Type = appaccestype
		if newpolicy.App_Access_Type == client.PolicyTransit {
			inNameSpace := d.Get("namespace").(string)
			if inNameSpace != "" {
				nameSpaceDetails, _ := prosimoClient.GetNamespaceByName(ctx, inNameSpace)
				newpolicy.NamespaceID = nameSpaceDetails.ID
			} else {
				log.Println("[ERROR]: Missing Namespace details")
			}
		}

		// if appaccestype == "Agentless" {
		// 	newpolicy.AppAccessType = "agentless"
		// } else if appaccestype == "Agent Based Access" {
		// 	newpolicy.AppAccessType = "agent"
		// }
	}
	if v, ok := d.GetOk("device_posture_configured"); ok {
		// teamid := v.(string)
		newpolicy.Device_Posture_Configured = v.(bool)
	}
	if v, ok := d.GetOk("details"); ok {
		detailsinput := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := detailsinput["actions"]; ok {
			action := v.(string)
			var actionList []string
			actionList = append(actionList, action)
			if action == "allow" {
				if v, ok := detailsinput["mfa"]; ok {
					mfa := v.(bool)
					if mfa {
						actionList = append(actionList, "mfa")
					}
				}
				if v, ok := detailsinput["alert"]; ok {
					alert := v.(bool)
					if alert {
						actionList = append(actionList, "alert")
					}
				}
				if v, ok := detailsinput["skipwaf"]; ok {
					alert := v.(bool)
					if alert {
						actionList = append(actionList, "skipwaf")
					}
				}
				if v, ok := detailsinput["bypass"]; ok {
					alert := v.(bool)
					if alert {
						actionList = append(actionList, "bypass")
					}
				}
			} else if action == "deny" {
				if v, ok := detailsinput["lock_users"]; ok {
					lockuser := v.(bool)
					if lockuser {
						actionList = append(actionList, "lockUser")
					}
				}
				if v, ok := detailsinput["alert"]; ok {
					alert := v.(bool)
					if alert {
						actionList = append(actionList, "alert")
					}
				}
			}
			newdetails.Actions = actionList
		}
		nwAclFlag := false
		egressFqdnFlag := false
		if v, ok := detailsinput["matches"].(*schema.Set); ok && v.Len() > 0 {
			matchdetails := v.List()[0].(map[string]interface{})
			userlist := []client.MatchDetails{}
			devicelist := []client.MatchDetails{}
			networkslist := []client.MatchDetails{}
			egressfqdnlist := []client.MatchDetails{}
			timelist := []client.MatchDetails{}
			urllist := []client.MatchDetails{}
			applist := []client.MatchDetails{}
			advancelist := []client.MatchDetails{}
			idplist := []client.MatchDetails{}
			dplist := []client.MatchDetails{}
			nwlist := []client.MatchDetails{}
			nwACLlist := []client.MatchDetails{}
			if v, ok := matchdetails["match_entries"].(*schema.Set); ok && v.Len() > 0 {
				for i, val := range v.List() {
					_ = val
					matchEntries := v.List()[i].(map[string]interface{})
					if v, ok := matchEntries["type"]; ok {
						types := v.(string)
						if types == "users" {
							userInput := matchUserEntry(meta, matchEntries)
							userlist = append(userlist, userInput)
							newmatchList.UserList = userlist
						} else if types == "devices" {
							deviceInput := matchDeviceEntry(meta, matchEntries)
							devicelist = append(devicelist, deviceInput)
							newmatchList.Devices = devicelist
						} else if types == "location" {
							_, networksInput := matchNetworkEntry(ctx, d, meta, matchEntries)
							networkslist = append(networkslist, networksInput)
							newmatchList.Networks = networkslist
						} else if types == "time" {
							timeInput := matchTimeEntry(meta, matchEntries)
							timelist = append(timelist, timeInput)
							newmatchList.Time = timelist
						} else if types == "url" {
							urlInput := matchURLEntry(meta, matchEntries)
							urllist = append(urllist, urlInput)
							newmatchList.URL = urllist
						} else if types == "fqdn" {
							applicationInput := matchAppEntry(meta, matchEntries)
							applist = append(applist, applicationInput)
							newmatchList.FQDN = applist
						} else if types == "advanced" {
							advancedInput := matchAdvanceEntry(meta, matchEntries)
							advancelist = append(advancelist, advancedInput)
							newmatchList.Advanced = advancelist
						} else if types == "idp" {
							idpInput := matchIDPEntry(meta, matchEntries)
							idplist = append(idplist, idpInput)
							newmatchList.IDP = idplist
						} else if types == "device-posture" {
							dpInput := matchDPEntry(meta, matchEntries)
							dplist = append(dplist, dpInput)
							newmatchList.Device_Posture_Profile = dplist
						} else if types == "networks" {
							nwInput := matchProsimoNWEntry(ctx, meta, matchEntries)
							nwlist = append(dplist, nwInput)
							newmatchList.ProsimoNetworks = nwlist
						} else if types == "networkacl" {
							nwAclFlag = true
							nwaclInput := matchNWACLEntry(meta, matchEntries)
							nwACLlist = append(dplist, nwaclInput)
							newmatchList.NetworkACL = nwACLlist
						} else if types == "egressfqdn" {
							egressFqdnFlag = true
							nwegressfqdnInput := matchEgressFqdnEntry(meta, matchEntries)
							egressfqdnlist = append(dplist, nwegressfqdnInput)
							newmatchList.EgressFqdns = egressfqdnlist
						}
						newdetails.Matches = newmatchList
					}
				}
			}
			if nwAclFlag && egressFqdnFlag {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Network ACL and Egress FQDN match conditions configured together",
					Detail:   "Network ACL and Egress FQDN match conditions are mutually exclusive, configure any one of them",
				})
				return diags, *newpolicy
			}
		}
		appsAttachmentFlag := false
		if v, ok := detailsinput["apps"].(*schema.Set); ok && v.Len() > 0 {
			appsAttachmentFlag = true
			appdetails := v.List()[0].(map[string]interface{})
			slectitemlist := []client.InputItems{}
			if v, ok := appdetails["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
				for i, val := range v.List() {
					_ = val
					selectapp := v.List()[i].(map[string]interface{})
					selectappName := selectapp["name"].(string)
					selectappid, _ := prosimoClient.GetAppID(ctx, selectappName)
					//selectappid := selectapp["id"].(string)
					newinputitems1 := client.InputItems{}
					newinputitems1.ItemName = selectappName
					newinputitems1.ItemID = selectappid
					slectitemlist = append(slectitemlist, newinputitems1)
				}
				newvalues1.SelectedItems = slectitemlist
			}
			newdetails.Apps = newvalues1

		}
		nwAttachmentFlag := false
		if v, ok := detailsinput["networks"].(*schema.Set); ok && v.Len() > 0 {
			nwAttachmentFlag = true
			networkdetails := v.List()[0].(map[string]interface{})
			slectitemlist := []client.InputItems{}
			if v, ok := networkdetails["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
				for i, val := range v.List() {
					_ = val
					selectnetwork := v.List()[i].(map[string]interface{})
					selectnetworkName := selectnetwork["name"].(string)
					selectnetworkid, _ := prosimoClient.GetNetworkID(ctx, selectnetworkName)
					//selectappid := selectapp["id"].(string)
					newinputitems1 := client.InputItems{}
					newinputitems1.ItemName = selectnetworkName
					newinputitems1.ItemID = selectnetworkid
					slectitemlist = append(slectitemlist, newinputitems1)
				}
				newvalues1.SelectedItems = slectitemlist
			}
			newdetails.Networks = newvalues1
		}
		internetFlag := false
		if v, ok := detailsinput["internet_traffic_enabled"]; ok {
			newdetails.Internet_Traffic_Enabled = v.(bool)
			internetFlag = newdetails.Internet_Traffic_Enabled
		}
		if (appsAttachmentFlag || nwAttachmentFlag) && internetFlag {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Internet Access is configured along with networks/apps attachments",
				Detail:   "Internet Access and apps/networks attachments are mutually exclusive, use any one of them",
			})
			return diags, *newpolicy
		}
		if (appsAttachmentFlag || nwAttachmentFlag) && egressFqdnFlag {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Egress FQDN match condition is configured along with networks/apps attachments",
				Detail:   "Egress FQDN match condition is only supported along with internet access",
			})
			return diags, *newpolicy
		}
		// } else {
		// 	newdetails.Networks = nil
		// }
		newpolicy.Details = newdetails

	}
	return nil, *newpolicy
}
func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	diags, newpolicy := inputDataops(ctx, d, meta)
	// log.Println("newpolicy", newpolicy)
	if diags != nil {
		return diags
	}
	policyListData, err := prosimoClient.CreatePolicy(ctx, &newpolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policyListData.Data.ID)

	return resourcePolicyRead(ctx, d, meta)
}

func matchUserEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newusermatchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	matchDetails := prosimoClient.ReadJson()
	flag := false
	for _, val := range matchDetails.Users.Property {
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
		if !flag {
			log.Println("[ERROR] Invalid entry for type User")
		}
		newusermatchdetails.Operations = FinalOperation
		newusermatchdetails.Property = property
	}
	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newusermatchdetails.Values = newuservalues

	}
	return newusermatchdetails
}

func matchTimeEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newtimematchdetails := client.MatchDetails{}
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	for _, val := range matchDetails.Time.Property {
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
	}
	if !flag {
		log.Println("[ERROR] Invalid entry for type Time")
	}
	Operations := "Between"

	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newtimematchdetails.Operations = Operations
		newtimematchdetails.Property = property
		newtimematchdetails.Values = newuservalues

	}
	return newtimematchdetails
}

func matchURLEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newurlmatchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	for _, val := range matchDetails.URL.Property {
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
		newurlmatchdetails.Operations = FinalOperation
		newurlmatchdetails.Property = property
	}
	if !flag {
		log.Println("[ERROR] Invalid value in type URL  fields")
	}
	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newurlmatchdetails.Operations = FinalOperation
		newurlmatchdetails.Property = property
		newurlmatchdetails.Values = newuservalues

	}
	return newurlmatchdetails
}

func matchIDPEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newIDPmatchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	property = matchEntries["property"].(string)
	for _, val := range matchDetails.IDP.Property {
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
	}
	if !flag {
		log.Println("[ERROR] Invalid entry in IDP operation field")
	}
	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newIDPmatchdetails.Operations = FinalOperation
		newIDPmatchdetails.Property = property
		newIDPmatchdetails.Values = newuservalues

	}
	return newIDPmatchdetails
}

func matchDPEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newDPmatchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	// log.Println("matchDetails", matchDetails)
	for _, val := range matchDetails.Device_Posture_Profile.Property {
		log.Println("dp val", val)
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
	}
	if !flag {
		log.Println("[ERROR] Invalid entry in Device Posture operation field")
	}
	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newDPmatchdetails.Operations = FinalOperation
		newDPmatchdetails.Property = property
		newDPmatchdetails.Values = newuservalues

	}
	return newDPmatchdetails
}

func matchAppEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newappmatchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	log.Println("matchDetails", matchDetails)
	for _, val := range matchDetails.FQDN.Property {
		log.Println("matchEntries[].(string)", matchEntries["property"].(string))
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
	}
	if !flag {
		log.Println("[ERROR] Invalid entry in type Application fileds")
	}

	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newappmatchdetails.Operations = FinalOperation
		newappmatchdetails.Property = property
		newappmatchdetails.Values = newuservalues

	}
	return newappmatchdetails
}

func matchNetworkEntry(ctx context.Context, d *schema.ResourceData, meta interface{}, matchEntries map[string]interface{}) (diag.Diagnostics, client.MatchDetails) {
	prosimoClient := meta.(*client.ProsimoClient)
	// var diags diag.Diagnostics
	newnetworkmatchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	newselecteditems := client.InputItems{}
	slectitemlist := []client.InputItems{}
	newmatchvalues := client.Values{}
	locSearch := &client.LocationSearchPayload{}
	flag := false
	matchDetails := prosimoClient.ReadJson()
	for _, val := range matchDetails.Location.Property {
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
	}
	if !flag {
		log.Println("[ERROR] Invalid value in type Network fields")
	}
	//--------Validate CITY/STATE/COUNTRY combination----------
	if property == "geoip_country_code" {
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			if v, ok := inputval["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
				for _, val := range v.List() {
					selecteditem := val.(map[string]interface{})
					if property == "geoip_country_code" {
						retLoc, _ := prosimoClient.FetchGeoLocation(ctx, locSearch)
						flag := false
						for _, countryDetails := range retLoc {
							// log.Println("country name", selecteditem["country_name"].(string))
							if countryDetails.CountryName != selecteditem["country_name"].(string) {
								continue
							} else {
								flag = true
								if v, ok := selecteditem["state_name"].(string); ok {
									if v != "" {
										locSearch.CountryCodeISO2 = countryDetails.CountryCodeISO2
										locSearch.Value = v
										geoDetails, _ := prosimoClient.FetchGeoLocation(ctx, locSearch)
										// for _, geoDetail := range geoDetails {
										// if geoDetail.StateName == v {
										if v1, ok := selecteditem["city_name"].(string); ok {
											if v1 != "" {
												// log.Println("valuev1", v1)
												flag1 := false
												for _, geoDetail := range geoDetails {
													// locSearch.Value = v1
													if geoDetail.CityName == v1 {
														flag1 = true
														// log.Println("geoDetail", geoDetail)
														// geoCityDetails, _ := prosimoClient.FetchGeoLocation(ctx, locSearch)
														// for _, countryDetails := range geoDetails {
														newselecteditems.CityCode = geoDetail.CityCode
														newselecteditems.CountryName = geoDetail.CountryName
														newselecteditems.ItemID = strconv.Itoa(geoDetail.CityCode)
														newselecteditems.StateName = geoDetail.StateName
														newselecteditems.CityName = geoDetail.CityName
														newselecteditems.ItemName = strings.Join([]string{geoDetail.CityName, geoDetail.StateName, geoDetail.CountryName}, ",")
														newselecteditems.Region = geoDetail.Region
														newselecteditems.CountryCodeISO2 = geoDetail.CountryCodeISO2
													}
												}
												if !flag1 {
													log.Println("[ERROR]: City doesnot exist")
												}

											} else {
												for _, geoDetail := range geoDetails {
													if geoDetail.CityName == "?" {
														newselecteditems.CityCode = geoDetail.CityCode
														newselecteditems.CountryName = geoDetail.CountryName
														newselecteditems.ItemID = strconv.Itoa(geoDetail.CityCode)
														newselecteditems.StateName = geoDetail.StateName
														// newselecteditems.CityName = geoDetail.CityName
														newselecteditems.ItemName = strings.Join([]string{geoDetail.StateName, geoDetail.CountryName}, ",")
														newselecteditems.Region = geoDetail.Region
														newselecteditems.CountryCodeISO2 = geoDetail.CountryCodeISO2
													}
												}
											}
										}
									} else {
										newselecteditems.CityCode = countryDetails.CityCode
										newselecteditems.CountryCodeISO2 = countryDetails.CountryCodeISO2
										newselecteditems.ItemID = strconv.Itoa(countryDetails.CityCode)
										newselecteditems.ItemName = countryDetails.CountryName
										newselecteditems.CountryName = countryDetails.CountryName

									}
								}
							}
						}
						if !flag {
							log.Println("[ERROR]: Invalid Country Name")
						}
					}
					slectitemlist = append(slectitemlist, newselecteditems)
				}
				newmatchvalues.SelectedItems = slectitemlist
			}
			newnetworkmatchdetails.Operations = FinalOperation
			newnetworkmatchdetails.Property = property
			newnetworkmatchdetails.Values = newmatchvalues
		}
	} else {
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			newmatchvalues := getMatchValues(inputval)
			newnetworkmatchdetails.Operations = FinalOperation
			newnetworkmatchdetails.Property = property
			newnetworkmatchdetails.Values = newmatchvalues
		}
	}

	return nil, newnetworkmatchdetails
}

// func locSearch()
func matchDeviceEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newdevicematchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	if matchEntries["property"].(string) == "Device OS" {
		for _, val := range matchDetails.Devices.Property {
			if matchEntries["property"].(string) == val.User_Property {
				property = val.Server_Property
				flag = true
			}
			for _, val := range val.Operations {
				if matchEntries["operation"].(string) == val.User_Operation_Name {
					FinalOperation = val.Server_Operation_Name
					flag = true
				}
			}
		}
		if !flag {
			log.Println("[ERROR]:Invalid value in type Device fields")
		}
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			device_os_list := [6]string{"Windows", "Linux", "Mac OSX", "iOS", "Android", "Chrome"}
			newmatchvalues := client.Values{}
			newselecteditems := client.InputItems{}
			slectitemlist := []client.InputItems{}
			if v, ok := inputval["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
				flag := true
				for i, val := range v.List() {
					_ = val
					selecteditem := v.List()[i].(map[string]interface{})
					selecteditemName := selecteditem["name"].(string)
					for _, item := range device_os_list {
						if item == selecteditemName {
							flag = false
						} else {
							continue
						}
					}
					if !flag {
						selecteditemId := selecteditem["name"].(string)
						newselecteditems.ItemName = selecteditemName
						newselecteditems.ItemID = selecteditemId
					} else {
						log.Println("[ERROR]:Invalid Selected item in Device")
					}
					slectitemlist = append(slectitemlist, newselecteditems)
				}
			}
			newmatchvalues.SelectedItems = slectitemlist
			newdevicematchdetails.Operations = FinalOperation
			newdevicematchdetails.Property = property
			newdevicematchdetails.Values = newmatchvalues
		}
	} else if matchEntries["property"].(string) == "Device OS Version" {
		for _, val := range matchDetails.Devices.Property {
			if matchEntries["property"].(string) == val.User_Property {
				property = val.Server_Property
				flag = true
			}
			for _, val := range val.Operations {
				if matchEntries["operation"].(string) == val.User_Operation_Name {
					FinalOperation = val.Server_Operation_Name
					flag = true
				}
			}
		}
		if !flag {
			log.Println("[ERROR]:Invalid value in type Device fields")
		}
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			newuservalues := getMatchValues(inputval)
			newdevicematchdetails.Operations = FinalOperation
			newdevicematchdetails.Property = property
			newdevicematchdetails.Values = newuservalues
		}

	} else if matchEntries["property"].(string) == "Device Category" {
		for _, val := range matchDetails.Devices.Property {
			if matchEntries["property"].(string) == val.User_Property {
				property = val.Server_Property
				flag = true
			}
			for _, val := range val.Operations {
				if matchEntries["operation"].(string) == val.User_Operation_Name {
					FinalOperation = val.Server_Operation_Name
					flag = true
				}
			}
		}
		if !flag {
			log.Println("[ERROR]:Invalid value in type Device fields")
		}
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			newmatchvalues := client.Values{}
			newselecteditems := client.InputItems{}
			slectitemlist := []client.InputItems{}
			var selecteditemId string
			if v, ok := inputval["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
				for i, val := range v.List() {
					_ = val
					selecteditem := v.List()[i].(map[string]interface{})
					selecteditemName := selecteditem["name"].(string)
					if selecteditemName == "Crawler" {
						selecteditemId = "crawler"
					} else if selecteditemName == "Desktop" {
						selecteditemId = "pc"
					} else if selecteditemName == "Mobile" {
						selecteditemId = "smartphone"
					} else {
						log.Println("[ERROR]:Invalid Selected item in Device")
					}
					newselecteditems.ItemName = selecteditemName
					newselecteditems.ItemID = selecteditemId
					slectitemlist = append(slectitemlist, newselecteditems)
				}
			}
			newmatchvalues.SelectedItems = slectitemlist
			newdevicematchdetails.Operations = FinalOperation
			newdevicematchdetails.Property = property
			newdevicematchdetails.Values = newmatchvalues
		}
	} else if matchEntries["property"].(string) == "Browser" {
		for _, val := range matchDetails.Devices.Property {
			if matchEntries["property"].(string) == val.User_Property {
				property = val.Server_Property
				flag = true
			}
			for _, val := range val.Operations {
				if matchEntries["operation"].(string) == val.User_Operation_Name {
					FinalOperation = val.Server_Operation_Name
					flag = true
				}
			}
		}
		if !flag {
			log.Println("[ERROR]:Invalid value in type Device fields")
		}
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			device_os_list := [6]string{"Chrome", "Firefox", "Internet Explorer", "Edge", "Safari", "Opera"}
			newmatchvalues := client.Values{}
			newselecteditems := client.InputItems{}
			slectitemlist := []client.InputItems{}
			if v, ok := inputval["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
				flag := true
				for i, val := range v.List() {
					_ = val
					selecteditem := v.List()[i].(map[string]interface{})
					selecteditemName := selecteditem["name"].(string)
					for _, item := range device_os_list {
						if item == selecteditemName {
							flag = false
						} else {
							continue
						}
					}
					if !flag {
						selecteditemId := selecteditem["name"].(string)
						newselecteditems.ItemName = selecteditemName
						newselecteditems.ItemID = selecteditemId
					} else {
						log.Println("[ERROR]:Invalid Selected item in Device")
					}

					slectitemlist = append(slectitemlist, newselecteditems)
				}
			}
			newmatchvalues.SelectedItems = slectitemlist
			newdevicematchdetails.Operations = FinalOperation
			newdevicematchdetails.Property = property
			newdevicematchdetails.Values = newmatchvalues
		}
	} else if matchEntries["property"].(string) == "Browser Version" {
		for _, val := range matchDetails.Devices.Property {
			if matchEntries["property"].(string) == val.User_Property {
				property = val.Server_Property
				flag = true
			}
			for _, val := range val.Operations {
				if matchEntries["operation"].(string) == val.User_Operation_Name {
					FinalOperation = val.Server_Operation_Name
					flag = true
				}
			}
		}
		if !flag {
			log.Println("[ERROR]:Invalid value in type Device fields")
		}
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			newuservalues := getMatchValues(inputval)
			newdevicematchdetails.Operations = FinalOperation
			newdevicematchdetails.Property = property
			newdevicematchdetails.Values = newuservalues
		}
	} else if matchEntries["property"].(string) == "Trusted Device Certificate" {
		for _, val := range matchDetails.Devices.Property {
			if matchEntries["property"].(string) == val.User_Property {
				property = val.Server_Property
				flag = true
			}
			for _, val := range val.Operations {
				if matchEntries["operation"].(string) == val.User_Operation_Name {
					FinalOperation = val.Server_Operation_Name
					flag = true
				}
			}
		}
		if !flag {
			log.Println("[ERROR]:Invalid value in type Device fields")
		}
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			newuservalues := getMatchValues(inputval)
			newdevicematchdetails.Operations = FinalOperation
			newdevicematchdetails.Property = property
			newdevicematchdetails.Values = newuservalues
		}
	} else {
		log.Println("[ERROR]:Invalid value in Network property field")
	}
	return newdevicematchdetails
}

func matchAdvanceEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newadvancematchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	for _, val := range matchDetails.Advanced.Property {
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
		if !flag {
			log.Println("[ERROR]:Invalid entry in type Advanced fileds")
		}
		if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
			inputval := v.List()[0].(map[string]interface{})
			device_os_list := [9]string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
			newmatchvalues := client.Values{}
			newselecteditems := client.InputItems{}
			slectitemlist := []client.InputItems{}
			if v, ok := inputval["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
				flag := true
				for i, val := range v.List() {
					_ = val
					selecteditem := v.List()[i].(map[string]interface{})
					selecteditemName := selecteditem["name"].(string)
					for _, item := range device_os_list {
						if item == selecteditemName {
							flag = false
						} else {
							continue
						}
					}
					if !flag {
						selecteditemId := selecteditem["name"].(string)
						newselecteditems.ItemName = selecteditemName
						newselecteditems.ItemID = selecteditemId
					} else {
						log.Println("[ERROR]:Invalid Selected item in Device")
					}
					slectitemlist = append(slectitemlist, newselecteditems)
				}
			}
			newmatchvalues.SelectedItems = slectitemlist
			newadvancematchdetails.Operations = FinalOperation
			newadvancematchdetails.Property = property
			newadvancematchdetails.Values = newmatchvalues
		} else {
			log.Println("[ERROR]:Invalid value in Advance property field")
		}
	}
	return newadvancematchdetails
}

func matchNWACLEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newtimematchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	for _, val := range matchDetails.NetworkACL.Property {
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
	}
	if !flag {
		log.Println("[ERROR] Invalid entry for type Networks")
	}

	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newtimematchdetails.Operations = FinalOperation
		newtimematchdetails.Property = property
		newtimematchdetails.Values = newuservalues

	}
	return newtimematchdetails
}

func matchEgressFqdnEntry(meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newtimematchdetails := client.MatchDetails{}
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadJson()
	for _, val := range matchDetails.EgressFqdn.Property {
		if matchEntries["property"].(string) == val.User_Property {
			property = val.Server_Property
			flag = true
		}
		for _, val := range val.Operations {
			if matchEntries["operation"].(string) == val.User_Operation_Name {
				FinalOperation = val.Server_Operation_Name
				flag = true
			}
		}
	}
	if !flag {
		log.Println("[ERROR] Invalid entry for type EgressFqdns")
	}

	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getMatchValues(inputval)
		newtimematchdetails.Operations = FinalOperation
		newtimematchdetails.Property = property
		newtimematchdetails.Values = newuservalues

	}
	return newtimematchdetails
}

func matchProsimoNWEntry(ctx context.Context, meta interface{}, matchEntries map[string]interface{}) client.MatchDetails {
	prosimoClient := meta.(*client.ProsimoClient)
	newmatchdetails := client.MatchDetails{}
	var property string
	var Operations string
	matchDetails := prosimoClient.ReadJson()
	for _, val := range matchDetails.Networks.Property {
		property = val.Server_Property
		Operations = val.Operations[0].Server_Operation_Name
	}
	// if flag == false {
	// 	log.Println("[ERROR] Invalid entry for type Time")
	// }
	// Operations := "Between"

	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newmatchvalues := client.Values{}
		// newuservalues := getMatchValues(inputval)
		newmatchdetails.Operations = Operations
		newmatchdetails.Property = property

		if v, ok := inputval["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
			slectitemlist := []client.InputItems{}
			for _, val := range v.List() {
				selecteditem := val.(map[string]interface{})
				newselecteditems := client.InputItems{}
				newselecteditems.ItemName = selecteditem["name"].(string)
				selectnetworkid, _ := prosimoClient.GetNetworkID(ctx, newselecteditems.ItemName)
				if selectnetworkid != "" {
					// if err != nil {
					// 	return diag.FromErr(err)
					// }
					newselecteditems.ItemID = selectnetworkid
					slectitemlist = append(slectitemlist, newselecteditems)
				} else {
					log.Println("Invalid network details")
				}
			}
			newmatchvalues.SelectedItems = slectitemlist
		}
		newmatchdetails.Values = newmatchvalues
	}
	return newmatchdetails
}

func getMatchValues(inputval map[string]interface{}) client.Values {
	newmatchvalues := client.Values{}
	newinputitems := client.InputItems{}
	newselecteditems := client.InputItems{}
	newselectedgroups := client.InputItems{}
	inputitemlist := []client.InputItems{}
	slectitemlist := []client.InputItems{}
	selectitemgroup := []client.InputItems{}
	if v, ok := inputval["inputitems"].(*schema.Set); ok && v.Len() > 0 {
		for _, val := range v.List() {
			inputitem := val.(map[string]interface{})
			newinputitems.ItemName = inputitem["name"].(string)
			newinputitems.ItemID = inputitem["name"].(string)
			if v, ok := inputitem["ip_details"].(*schema.Set); ok && v.Len() > 0 {
				ipDetails := v.List()[0].(map[string]interface{})
				ipdetailsInput := &client.KeyValues{
					SourceIp:   expandStringList(ipDetails["source_ip"].([]interface{})),
					TargetIp:   expandStringList(ipDetails["target_ip"].([]interface{})),
					Protocol:   expandStringList(ipDetails["protocol"].([]interface{})),
					SourcePort: expandStringList(ipDetails["source_port"].([]interface{})),
					TargetPort: expandStringList(ipDetails["target_port"].([]interface{})),
				}
				newinputitems.KeyValues = ipdetailsInput
			} else if v, ok := inputitem["egress_fqdn_details"].(*schema.Set); ok && v.Len() > 0 {
				ipDetails := v.List()[0].(map[string]interface{})
				ipdetailsInput := &client.KeyValues{
					FqdnInverseMatch: expandStringList(ipDetails["fqdn_inverse_match"].([]interface{})),
					FqdnMatch:        expandStringList(ipDetails["fqdn_match"].([]interface{})),
					Protocol:         expandStringList(ipDetails["protocol"].([]interface{})),
					SourceIp:         expandStringList(ipDetails["source_ip"].([]interface{})),
					TargetPort:       expandStringList(ipDetails["target_port"].([]interface{})),
				}
				newinputitems.KeyValues = ipdetailsInput
			}
			inputitemlist = append(inputitemlist, newinputitems)
		}
	}
	if v, ok := inputval["selecteditems"].(*schema.Set); ok && v.Len() > 0 {
		for _, val := range v.List() {
			selecteditem := val.(map[string]interface{})
			newselecteditems.ItemName = selecteditem["name"].(string)
			newselecteditems.ItemID = selecteditem["name"].(string)
			slectitemlist = append(slectitemlist, newselecteditems)
		}
	}
	if v, ok := inputval["selectedgroups"].(*schema.Set); ok && v.Len() > 0 {
		for _, val := range v.List() {
			selectedgroup := val.(map[string]interface{})
			newselectedgroups.ItemName = selectedgroup["name"].(string)
			newselectedgroups.ItemID = selectedgroup["name"].(string)
			selectitemgroup = append(slectitemlist, newselectedgroups)
		}
	}
	newmatchvalues.InputItems = inputitemlist
	newmatchvalues.SelectedItems = slectitemlist
	newmatchvalues.SelectedGroups = selectitemgroup
	return newmatchvalues
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//resourcePolicyCreate(ctx, d, meta)
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	policyID := d.Id()
	// diags, newpolicy := inputDataops(ctx, d, meta, false)
	diags, newpolicy := inputDataops(ctx, d, meta)
	if diags != nil {
		return diags
	}
	newpolicy.ID = policyID
	policyListData, err := prosimoClient.UpdatePolicy(ctx, &newpolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policyListData.Data.ID)

	return resourcePolicyRead(ctx, d, meta)

}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	policyID := d.Id()
	res, err := prosimoClient.GetPolicyByID(ctx, policyID)

	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Policy Details",
			Detail:   fmt.Sprintf("Unable to find Policy for ID %s", policyID),
		})

		return diags
	}
	d.Set("id", res.ID)
	d.Set("name", res.Name)
	d.Set("teamid", res.TeamID)
	d.Set("types", res.PolicyType)
	d.Set("app_access_type", res.App_Access_Type)
	if res.App_Access_Type == client.PolicyTransit {
		d.Set("namespace", d.Get("namespace").(string))
	}

	return diags
}

func resourcePolicynDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policyId := d.Id()
	res, err := prosimoClient.GetPolicyByID(ctx, policyId)
	if err != nil {
		return diag.FromErr(err)
	}
	// Detach apps/networks or internet access from the policy before the destroy operation
	res.Details.Apps.SelectedItems = res.Details.Apps.SelectedItems[:0]
	res.Details.Networks.SelectedItems = res.Details.Networks.SelectedItems[:0]
	res.Details.Internet_Traffic_Enabled = false
	updateData, err := prosimoClient.UpdatePolicy(ctx, res)
	_ = updateData
	if err != nil {
		return diag.FromErr(err)
	}
	// Delete the policy
	res_err := prosimoClient.DeletePolicy(ctx, policyId)
	if res_err != nil {
		return diag.FromErr(res_err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags

}
