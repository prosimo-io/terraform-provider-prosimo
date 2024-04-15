package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConnPlacement() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify connector placement.",
		CreateContext: resourceConnPlacementCreate,
		UpdateContext: resourceConnPlacementCreate,
		ReadContext:   resourceConnPlacementRead,
		DeleteContext: resourceConnPlacementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"connector_placement_app_vpc": {
				Type:     schema.TypeBool,
				Required: true,
				// ForceNew: true,
			},
		},
	}
}

func resourceConnPlacementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	conPlacementinput := &client.Connector_Placement{}

	if v, ok := d.GetOk("connector_placement_app_vpc"); ok {
		// appconvpc := v.(bool)
		conPlacementinput.ConnectorPlacementAppVpc = v.(bool)
	}

	log.Printf("[DEBUG] Updating connector placement block: %v", conPlacementinput.ConnectorPlacementAppVpc)
	err := prosimoClient.PutConnectorPlacement(ctx, conPlacementinput)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Updated connector placement block")
	d.SetId("Connector Placement")

	return resourceConnPlacementRead(ctx, d, meta)

}

func resourceConnPlacementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	log.Printf("[DEBUG] Reading connector placement block")
	conplacementres, err := prosimoClient.GetConnectorPlacement(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("connector_placement_app_vpc", conplacementres.ConnectorPlacementAppVpcStatus.ConnectorPlacementAppVpc)

	return diags
}

func resourceConnPlacementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	conPlacementinput := &client.Connector_Placement{
		ConnectorPlacementAppVpc: true,
	}
	err := prosimoClient.PutConnectorPlacement(ctx, conPlacementinput)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] reset connector placement")
	d.SetId("")
	return resourceConnPlacementRead(ctx, d, meta)

}
