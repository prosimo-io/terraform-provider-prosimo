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

func resourceInternetEgress() *schema.Resource {
	return &schema.Resource{
		Description:   "The Internet Egress in Prosimo helps control internet access rules for applications. Use this resource to create/modify egress rules.",
		CreateContext: resourceInternetEgressCreate,
		UpdateContext: resourceInternetEgressUpdate,
		ReadContext:   resourceInternetEgressRead,
		DeleteContext: resourceInternetEgressDelete,
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
				Description: "Name of Internet Egress Policy",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(client.GetPolicyActionTypes(), false),
				Description:  "Policy action, e.g: allow, deny",
			},
			"namespaces": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Policy Namespace where the policy can be in the action",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace_entries": {
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
			"networks": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Network details to attach to the policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_entries": {
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
			"network_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Network group details to attach to the policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_group_entries": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the network-group",
									},
								},
							},
						},
					},
				},
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
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Select policy match condition type i.e. - fqdn",
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
										Description:  "Operation of the selected property, available options are Is, Is NOT",
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
															"id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Input domain/fqdn value",
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

func inputIEDataops(ctx context.Context, d *schema.ResourceData, meta interface{}) (diag.Diagnostics, client.InternetEgress) {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	newpolicy := &client.InternetEgress{}
	// newdetails := client.Detaxils{}

	newmatchList := client.IE_Matches{}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		newpolicy.Name = name
	}
	if v, ok := d.GetOk("action"); ok {
		action := v.(string)
		newpolicy.Action = action
	}
	if v, ok := d.GetOk("matches"); ok {
		matchdetails := v.(*schema.Set).List()[0].(map[string]interface{})
		fqdnlist := []client.FqdnDetails{}
		if v, ok := matchdetails["match_entries"].(*schema.Set); ok && v.Len() > 0 {
			for i, _ := range v.List() {
				matchEntries := v.List()[i].(map[string]interface{})
				if v, ok := matchEntries["type"]; ok {
					types := v.(string)
					if types == "fqdn" {
						fqdnInput, diags := matchFqdnEntry(meta, matchEntries)
						if diags != nil {
							return diags, *newpolicy
						}
						fqdnlist = append(fqdnlist, fqdnInput)
						newmatchList.Fqdn = &fqdnlist
					}
				}
			}
		}
		newpolicy.Matches = &newmatchList

	}
	if v, ok := d.GetOk("networks"); ok {
		networkdetails := v.(*schema.Set).List()[0].(map[string]interface{})
		slectitemlist := []client.NetworkList{}
		if v, ok := networkdetails["network_entries"].(*schema.Set); ok && v.Len() > 0 {
			for i, _ := range v.List() {
				selectnetwork := v.List()[i].(map[string]interface{})
				selectnetworkName := selectnetwork["name"].(string)
				selectnetworkid, _ := prosimoClient.GetNetworkID(ctx, selectnetworkName)
				newinputitems1 := client.NetworkList{}
				newinputitems1.ID = selectnetworkid
				slectitemlist = append(slectitemlist, newinputitems1)
			}
		}
		newpolicy.Networks = &slectitemlist
	}
	if v, ok := d.GetOk("network_groups"); ok {
		networkdetails := v.(*schema.Set).List()[0].(map[string]interface{})
		slectitemlist := []client.NetworkGroupList{}
		if v, ok := networkdetails["network_group_entries"].(*schema.Set); ok && v.Len() > 0 {
			for i, _ := range v.List() {
				selectnetworkGroup := v.List()[i].(map[string]interface{})
				selectnetworkGroupName := selectnetworkGroup["name"].(string)
				selectnetworkgroupconfig, _ := prosimoClient.GetGrpConfByName(ctx, "NETWORK", selectnetworkGroupName)
				newinputitems1 := client.NetworkGroupList{}
				newinputitems1.ID = selectnetworkgroupconfig.Id
				slectitemlist = append(slectitemlist, newinputitems1)
			}
		}
		newpolicy.NetworkGroups = &slectitemlist
	}
	if v, ok := d.GetOk("namespaces"); ok {
		namespacedetails := v.(*schema.Set).List()[0].(map[string]interface{})
		slectitemlist := []client.NamespaceList{}
		if v, ok := namespacedetails["namespace_entries"].(*schema.Set); ok && v.Len() > 0 {
			for i, _ := range v.List() {
				selectnamespace := v.List()[i].(map[string]interface{})
				selectnamespaceName := selectnamespace["name"].(string)
				selectnamespace1, _ := prosimoClient.GetNamespaceByName(ctx, selectnamespaceName)
				newinputitems1 := client.NamespaceList{}
				newinputitems1.ID = selectnamespace1.ID
				slectitemlist = append(slectitemlist, newinputitems1)
			}
		}
		newpolicy.Namespaces = &slectitemlist
	}
	return diags, *newpolicy
}

func matchFqdnEntry(meta interface{}, matchEntries map[string]interface{}) (client.FqdnDetails, diag.Diagnostics) {
	prosimoClient := meta.(*client.ProsimoClient)
	newfqdnmatchdetails := client.FqdnDetails{}
	var diags diag.Diagnostics
	var FinalOperation string
	var property string
	flag := false
	matchDetails := prosimoClient.ReadInternetEgressJson()
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
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid entry in type Application fileds",
			Detail:   fmt.Sprintf("There seems to be a mismatch with the defined json payload for internet egress policy\nFinalOperation : %s , property : %s", FinalOperation, property),
		})
	}
	if FinalOperation == "" || property == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Either FinalOperation or property didn't match the expected schema fileds",
			Detail:   fmt.Sprintf("FinalOperation : %s , property : %s", FinalOperation, property),
		})
	}
	if v, ok := matchEntries["values"].(*schema.Set); ok && v.Len() > 0 {
		inputval := v.List()[0].(map[string]interface{})
		newuservalues := getIEMatchValues(inputval)
		newfqdnmatchdetails.Operation = FinalOperation
		newfqdnmatchdetails.Property = property
		newfqdnmatchdetails.Values = &newuservalues
	}
	return newfqdnmatchdetails, diags
}

func getIEMatchValues(inputval map[string]interface{}) client.Values {
	newmatchvalues := client.Values{}
	newinputitems := client.InputItems{}
	// newselecteditems := client.InputItems{}
	inputitemlist := []client.InputItems{}
	if v, ok := inputval["inputitems"].(*schema.Set); ok && v.Len() > 0 {
		for _, val := range v.List() {
			inputitem := val.(map[string]interface{})
			// newinputitems.ItemName = inputitem["name"].(string)
			newinputitems.ItemID = inputitem["id"].(string)
			inputitemlist = append(inputitemlist, newinputitems)
		}
	}
	newmatchvalues.InputItems = inputitemlist
	return newmatchvalues
}

func resourceInternetEgressCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)
	diags, newpolicy := inputIEDataops(ctx, d, meta)
	if diags != nil {
		return diags
	}

	policyListData, err := prosimoClient.CreateInternetEgressPolicy(ctx, &newpolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policyListData.Data.ID)

	return resourceInternetEgressRead(ctx, d, meta)
}

func resourceInternetEgressUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//resourcePolicyCreate(ctx, d, meta)
	prosimoClient := meta.(*client.ProsimoClient)
	policyID := d.Id()
	_, newpolicy := inputIEDataops(ctx, d, meta)
	newpolicy.ID = policyID
	policyListData, err := prosimoClient.UpdateInternetEgressPolicy(ctx, &newpolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policyListData.Data.ID)

	return resourceInternetEgressRead(ctx, d, meta)

}

func resourceInternetEgressRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	policyID := d.Id()
	res, err := prosimoClient.GetInternetEgressPolicyByID(ctx, policyID)

	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Internet Egress Policy Details",
			Detail:   fmt.Sprintf("Unable to find Policy for ID %s", policyID),
		})

		return diags
	}
	d.Set("id", res.ID)
	d.Set("name", res.Name)
	d.Set("action", res.Action)
	// d.Set("networks", res.Networks)
	// d.Set("networkgroups", res.NetworkGroups)
	// d.Set("namespaces", res.Namespaces)
	// d.Set("matches", res.Matches)

	return diags
}

func resourceInternetEgressDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policyId := d.Id()

	res_err := prosimoClient.DeleteInternetEgressPolicy(ctx, policyId)
	if res_err != nil {
		return diag.FromErr(res_err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags

}
