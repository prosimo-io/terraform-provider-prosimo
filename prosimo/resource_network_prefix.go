package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkPrefix() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Network Prefixes.",
		CreateContext: resourceNPCreate,
		UpdateContext: resourceNPUpdate,
		DeleteContext: resourceNPDelete,
		ReadContext:   resourceNPRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"cloud_account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cloud account name",
			},
			"cloud_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of cloud region",
			},
			"cloud_network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of cloud network",
			},
			"prefix_route_tables": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Selected Regions where the entries would be added.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CIDR range",
						},
						"route_tables": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of Cloud Regions: e.g: us-east-1, eastus",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"route_table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of the route table",
									},
								},
							},
						},
					},
				},
			},
			"enable": {
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

func resourceNPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	prefixRTList := []client.Prefix_RouteID{}
	routeTableList := []client.Route_Tables{}
	if v, ok := d.GetOk("prefix_route_tables"); ok {
		for _, value := range v.([]interface{}) {
			prefixRT := value.(map[string]interface{})
			if v, ok := prefixRT["route_tables"]; ok {
				for _, value := range v.([]interface{}) {
					routeTable := value.(map[string]interface{})
					inputRT := client.Route_Tables{
						ID: routeTable["route_table"].(string),
					}
					routeTableList = append(routeTableList, inputRT)
				}
			}
			inputPrefixRT := client.Prefix_RouteID{
				Prefix:        prefixRT["ip_prefix"].(string),
				RouteTableIDS: routeTableList,
			}
			prefixRTList = append(prefixRTList, inputPrefixRT)
		}
	}
	cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, d.Get("cloud_account").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	networkRouteInput := &client.Route_entry_network{
		CloudKeyID:     cloudCreds.ID,
		CSP:            cloudCreds.CloudType,
		CloudRegion:    d.Get("cloud_region").(string),
		CloudNetworkID: d.Get("cloud_network").(string),
		PrefixesRT:     &prefixRTList,
	}

	log.Printf("[DEBUG]Creating network prefix entry : %v", networkRouteInput)
	createNWRouteEntry, err := prosimoClient.CreateNetworkRouteEntry(ctx, networkRouteInput)
	if err != nil {
		log.Printf("[ERROR] Error in network prefix config")
		return diag.FromErr(err)
	}
	log.Println("id", createNWRouteEntry.Data.ID)
	d.SetId(createNWRouteEntry.Data.ID)
	if d.Get("enable").(bool) {
		_, err := prosimoClient.EnableNetworkRouteEntry(ctx, createNWRouteEntry.Data.ID)
		if err != nil {
			log.Printf("[ERROR] Error in enabling network prefix config")
			return diag.FromErr(err)
		}
	}

	return resourceNPRead(ctx, d, meta)
}

func resourceNPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics
	update_req := false
	if d.HasChange("cloud_account") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify cloud account",
			Detail:   "cloud account can't be modified",
		})
		return diags
	}
	if d.HasChange("cloud_region") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify cloud region",
			Detail:   "cloud region can't be modified",
		})
		return diags
	}
	if d.HasChange("cloud_network") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify cloud network",
			Detail:   "cloud network can't be modified",
		})
		return diags
	}
	if d.HasChange("prefix_route_tables") && !d.IsNewResource() {
		update_req = true
	}
	// if d.HasChange("enable") && !d.IsNewResource() {
	// 	update_req = true
	// }
	if update_req {
		prefixRTList := []client.Prefix_RouteID{}
		routeTableList := []client.Route_Tables{}
		ns, err := prosimoClient.GetNetworkRouteEntryByID(ctx, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		if ns.Enabled {
			_, err := prosimoClient.DisableNetworkRouteEntry(ctx, d.Id())
			if err != nil {
				log.Printf("[ERROR] Error in disabling network prefix config")
				return diag.FromErr(err)
			}
		}
		if v, ok := d.GetOk("prefix_route_tables"); ok {
			for _, value := range v.([]interface{}) {
				prefixRT := value.(map[string]interface{})
				if v, ok := prefixRT["route_tables"]; ok {
					for _, value := range v.([]interface{}) {
						routeTable := value.(map[string]interface{})
						inputRT := client.Route_Tables{
							ID: routeTable["route_table"].(string),
						}
						routeTableList = append(routeTableList, inputRT)
					}
				}
				inputPrefixRT := client.Prefix_RouteID{
					Prefix:        prefixRT["ip_prefix"].(string),
					RouteTableIDS: routeTableList,
				}
				prefixRTList = append(prefixRTList, inputPrefixRT)
			}
		}
		cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, d.Get("cloud_account").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		networkRouteInput := &client.Route_entry_network{
			ID:             d.Id(),
			CloudKeyID:     cloudCreds.ID,
			CSP:            cloudCreds.CloudType,
			CloudRegion:    d.Get("cloud_region").(string),
			CloudNetworkID: d.Get("cloud_network").(string),
			PrefixesRT:     &prefixRTList,
		}

		log.Printf("[DEBUG]Creating network prefix entry : %v", networkRouteInput)
		update_err := prosimoClient.UpdateNetworkRouteEntry(ctx, networkRouteInput)
		if err != nil {
			log.Printf("[ERROR] Error in network prefix config")
			return diag.FromErr(update_err)
		}
		if d.Get("enable").(bool) {
			_, err := prosimoClient.EnableNetworkRouteEntry(ctx, d.Id())
			if err != nil {
				log.Printf("[ERROR] Error in enabling network prefix config")
				return diag.FromErr(err)
			}
		}
	} else {
		if d.HasChange("enable") && !d.IsNewResource() {
			if d.Get("enable").(bool) {
				_, err := prosimoClient.EnableNetworkRouteEntry(ctx, d.Id())
				if err != nil {
					log.Printf("[ERROR] Error in enabling network prefix config")
					return diag.FromErr(err)
				}
			} else {
				_, err := prosimoClient.DisableNetworkRouteEntry(ctx, d.Id())
				if err != nil {
					log.Printf("[ERROR] Error in disabling network prefix config")
					return diag.FromErr(err)
				}
			}

		}
	}
	return resourceNPRead(ctx, d, meta)
}

func resourceNPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	npID := d.Id()

	log.Printf("[DEBUG] Get Network Prefix with id  %s", npID)

	ns, err := prosimoClient.GetNetworkRouteEntryByID(ctx, npID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", ns.ID)
	// d.Set("prefix_route_tables", ns.PrefixesRT)
	d.Set("cloud_region", ns.CloudRegion)
	d.Set("cloud_network", ns.CloudNetworkID)
	d.Set("status", ns.Status)
	return diags
}

func resourceNPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	npID := d.Id()
	ns, err := prosimoClient.GetNetworkRouteEntryByID(ctx, npID)
	if err != nil {
		return diag.FromErr(err)
	}
	if ns.Enabled {
		_, err := prosimoClient.DisableNetworkRouteEntry(ctx, npID)
		if err != nil {
			log.Printf("[ERROR] Error in disabling network prefix config")
			return diag.FromErr(err)
		}
	}

	err_res := prosimoClient.DeleteNetworkRouteEntry(ctx, npID)
	if err != nil {
		return diag.FromErr(err_res)
	}
	log.Printf("[DEBUG] Deleted Route Prefix with - id - %s", npID)

	return diags
}
