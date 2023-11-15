package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRegionalPrefix() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Private Link Sources.",
		CreateContext: resourceRPCreate,
		UpdateContext: resourceRPUpdate,
		DeleteContext: resourceRPDelete,
		ReadContext:   resourceRPRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"cidr": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Region Level IP Prefixes",
			},
			"all_regions": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "By default summary entries are added to all VPCs and VNETs onboarded to Prosimo.",
			},
			"selected_regions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Selected Regions where the entries would be added.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CSP Details: E.G: AZURE/GCP/AWS",
						},
						"cloud_region": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of Cloud Regions: e.g: us-east-1, eastus",
						},
					},
				},
			},
			"overwrite_route": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceRPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	selectedRegList := []client.Selected_Reg{}
	if v, ok := d.GetOk("selected_regions"); ok {
		for _, value := range v.([]interface{}) {
			selectReg := value.(map[string]interface{})
			selectRegInput := client.Selected_Reg{
				CSP:   selectReg["cloud_type"].(string),
				Names: expandStringList(selectReg["cloud_region"].([]interface{})),
			}
			selectedRegList = append(selectedRegList, selectRegInput)
		}
	}
	regRouteInput := &client.Regions_route{
		All:      d.Get("all_regions").(bool),
		Selected: selectedRegList,
	}

	inputRouteEntry := client.Route_entry_region{
		Prefixes:       expandStringList(d.Get("cidr").([]interface{})),
		Regions:        regRouteInput,
		OverWriteRoute: d.Get("overwrite_route").(bool),
	}

	log.Printf("[DEBUG]Creating route entry : %v", inputRouteEntry)
	createRouteEntry, err := prosimoClient.CreateRouteEntry(ctx, &inputRouteEntry)
	if err != nil {
		log.Printf("[ERROR] Error in creating Group config")
		return diag.FromErr(err)
	}
	log.Println("id", createRouteEntry.Data.ID)
	d.SetId(createRouteEntry.Data.ID)
	return resourceRPRead(ctx, d, meta)
}

func resourceRPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	update_req := false
	if d.HasChange("cidr") && !d.IsNewResource() {
		update_req = true
	}

	if d.HasChange("all_regions") && !d.IsNewResource() {
		update_req = true
	}

	if d.HasChange("selected_regions") {
		update_req = true
	}

	if update_req {
		selectedRegList := []client.Selected_Reg{}
		if v, ok := d.GetOk("selected_regions"); ok {
			for _, value := range v.([]interface{}) {
				selectReg := value.(map[string]interface{})
				selectRegInput := client.Selected_Reg{
					CSP:   selectReg["cloud_type"].(string),
					Names: expandStringList(selectReg["cloud_region"].([]interface{})),
				}
				selectedRegList = append(selectedRegList, selectRegInput)
			}
		}
		regRouteInput := &client.Regions_route{
			All:      d.Get("all_regions").(bool),
			Selected: selectedRegList,
		}

		inputRouteEntry := client.Route_entry_region{
			ID:             d.Id(),
			Prefixes:       expandStringList(d.Get("cidr").([]interface{})),
			Regions:        regRouteInput,
			OverWriteRoute: d.Get("overwrite_route").(bool),
		}

		log.Printf("[DEBUG]Updating route entry : %v", inputRouteEntry)
		err := prosimoClient.UpdateRouteEntry(ctx, &inputRouteEntry)
		if err != nil {
			log.Printf("[ERROR] Error in creating Group config")
			return diag.FromErr(err)
		}
	}
	return resourceRPRead(ctx, d, meta)
}

func resourceRPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	nsID := d.Id()

	log.Printf("[DEBUG] Get Route Prefix with id  %s", nsID)

	ns, err := prosimoClient.GetRouteEntryByID(ctx, nsID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", ns.ID)
	d.Set("cidr", ns.Prefixes)
	d.Set("all_regions", ns.Regions.All)
	d.Set("selected_regions", ns.Regions.Selected)
	d.Set("overwrite_route", ns.OverWriteRoute)
	d.Set("status", ns.Status)
	return diags
}

func resourceRPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	rpID := d.Id()

	err := prosimoClient.DeleteRouteEntry(ctx, rpID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Deleted Route Prefix with - id - %s", rpID)

	return diags
}
