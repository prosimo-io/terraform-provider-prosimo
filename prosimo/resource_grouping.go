package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceGrpConfig() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify resource group settings.",
		CreateContext: resourceGrpConfigCreate,
		ReadContext:   resourceGrpConfigRead,
		DeleteContext: resourceGrpConfigDelete,
		UpdateContext: resourceGrpConfigUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the grouping",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(client.GetGroupingTypes(), false),
				Description:  "Grouping type, e.g: USER, APP, DEVICE, TIME, IP_RANGE, GEO",
			},
			"sub_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": {
				Type:     schema.TypeList,
				MinItems: 1,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"details": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps": {
							Type:     schema.TypeList,
							MinItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"names": {
							Type:     schema.TypeList,
							MinItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"time": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from": {
										Type:     schema.TypeString,
										Required: true,
									},
									"to": {
										Type:     schema.TypeString,
										Required: true,
									},
									"timezone": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"ranges": {
							Type:     schema.TypeList,
							MinItems: 1,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func inputDataopsGrouping(ctx context.Context, d *schema.ResourceData, meta interface{}) (diag.Diagnostics, client.Grp_Config) {
	prosimoClient := meta.(*client.ProsimoClient)
	grpConfig := &client.Grp_Config{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}
	deTails := client.DeTails{}

	if grpConfig.Type == client.TypeUSER {
		if v, ok := d.GetOk("members"); ok {
			ranges := v.([]interface{})

			if len(ranges) > 0 {
				grpConfig.Members = expandStringList(v.([]interface{}))
			}
		}
	}

	if grpConfig.Type == client.TypeDEVICE {
		flag := false
		// var FinalOperation string
		if v, ok := d.GetOk("sub_type"); ok {
			subType := v.(string)
			// log.Println("subtype", subType)
			// grpConfig.SubType = subType
			matchDetails := prosimoClient.ReadJson()

			for _, val := range matchDetails.Devices.Property {
				if subType == val.User_Property {
					grpConfig.SubType = val.Server_Property
					flag = true
				}
			}

		}
		if !flag {
			log.Println("[ERROR]:Invalid/Missing value in  Device sub_type field")
		}
	}

	if v, ok := d.GetOk("details"); ok {
		details := v.(*schema.Set).List()[0].(map[string]interface{})
		apps := []client.App{}
		if grpConfig.Type == client.TypeAPP {
			if v, ok := details["apps"]; ok {
				appList := v.([]interface{})
				for _, app := range appList {
					val := app.(map[string]interface{})
					appInput := client.App{}
					if v, ok := val["name"].(string); ok {
						appInput.Name = v
						appid, _ := prosimoClient.GetAppID(ctx, appInput.Name)
						appInput.Id = appid
					}
					apps = append(apps, appInput)

				}
				deTails.Apps = apps
			}
		}
		if grpConfig.Type == client.TypeDEVICE {
			if v, ok := details["names"]; ok {
				nameList := v.([]interface{})
				names := []client.Name{}
				for _, name := range nameList {
					val := name.(map[string]interface{})
					nameInput := client.Name{}
					// if v, ok := val["id"].(string); ok {
					// 	nameInput.Id = v
					// }
					if v, ok := val["name"].(string); ok {
						nameInput.Name = v
						nameInput.Id = nameInput.Name
					}

					names = append(names, nameInput)

				}
				deTails.Names = names
				// regionName := v.(string)
				// cloudconfigInput.RegionName = regionName
			}
		}

		if grpConfig.Type == client.TypeTIME {
			if v, ok := details["time"]; ok {
				timeInput := client.Time{}
				timeInputList := []client.Time{}
				timeDetails := v.(*schema.Set).List()[0].(map[string]interface{})

				from := timeDetails["from"].(string)
				timeInput.From = from

				to := timeDetails["to"].(string)
				timeInput.To = to

				timeZone := timeDetails["timezone"].(string)
				timeInput.TimeZone = timeZone

				timeInputList = append(timeInputList, timeInput)
				deTails.Time = timeInputList
			}
		}

		if grpConfig.Type == client.TypeIP_RANGE {
			if v, ok := details["ranges"]; ok {
				ranges := v.([]interface{})

				if len(ranges) > 0 {
					deTails.Ranges = expandStringList(v.([]interface{}))
				}
				// d.SetId("Ipreputation")
			}
		}

		grpConfig.Details = deTails

	}
	return nil, *grpConfig
}

func resourceGrpConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	_, grpConfig := inputDataopsGrouping(ctx, d, meta)

	log.Printf("[DEBUG]Creating grp Config : %v", grpConfig)
	createGrpConfig, err := prosimoClient.CreateGrpConf(ctx, &grpConfig)
	if err != nil {
		log.Printf("[ERROR] Error in creating Group config")
		return diag.FromErr(err)
	}
	log.Println("grpid", createGrpConfig.Grp_Config_Res.Id)
	d.SetId(createGrpConfig.Grp_Config_Res.Id)
	return resourceGrpConfigRead(ctx, d, meta)
}

func resourceGrpConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	grpID := d.Id()
	_, grpconfig := inputDataopsGrouping(ctx, d, meta)

	log.Printf("[DEBUG] Get Grouping profile for %s", grpID)

	res, err := prosimoClient.GetGrpConf(ctx, grpconfig.Type)
	if err != nil {
		return diag.FromErr(err)
	}
	var grpConfig *client.Grp_Config
	for _, returnedGrpConfig := range res.Grp_Config_Res.Records {
		if returnedGrpConfig.Id == grpID {
			grpConfig = returnedGrpConfig
			break
		}
	}
	if grpConfig == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[ERROR] Unable to get Group config",
		})
		return diags
	}

	d.Set("name", grpConfig.Name)
	d.Set("type", grpConfig.Type)

	return diags
}

func resourceGrpConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)
	grpID := d.Id()

	_, grpConfig := inputDataopsGrouping(ctx, d, meta)
	grpConfig.Id = grpID

	log.Printf("[DEBUG] Creating grp Config : %v", grpConfig)
	updateGrpConfig, err := prosimoClient.UpdateGrpConf(ctx, &grpConfig)
	if err != nil {
		log.Printf("[ERROR] Error in creating Group config")
		return diag.FromErr(err)
	}
	d.SetId(updateGrpConfig.Grp_Config_Res.Id)
	return resourceGrpConfigRead(ctx, d, meta)
}

func resourceGrpConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics

	grpID := d.Id()
	err := prosimoClient.DeleteGrpConf(ctx, grpID)

	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
