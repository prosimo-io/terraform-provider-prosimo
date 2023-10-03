package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIPReputation() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify ip reputation settngs.",
		CreateContext: resourceIPReputationCreate,
		ReadContext:   resourceIPReputationRead,
		DeleteContext: resourceIPReputationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Defaults to false, set it true to enable IP Reputation",
			},
			"allowlist": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of CIDRs to allow",
			},
		},
	}
}

func resourceIPReputationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	ipReP := &client.IP_Reputation{}

	if v, ok := d.GetOk("enabled"); ok {
		enabled := v.(bool)
		ipReP.Enabled = enabled
	}

	if v, ok := d.GetOk("allowlist"); ok {
		CidrList := v.([]interface{})

		if len(CidrList) > 0 {
			ipReP.AllowList = expandStringList(v.([]interface{}))
		}
		// d.SetId("Ipreputation")
	}

	log.Printf("[DEBUG] Updating IP Reputation block: %t %v", ipReP.Enabled, ipReP.AllowList)
	_, err := prosimoClient.PutIPREP(ctx, ipReP)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Updated IP Reputation block: %t %v ", ipReP.Enabled, ipReP.AllowList)
	d.SetId("Ipreputation")

	return resourceIPReputationRead(ctx, d, meta)

}

func resourceIPReputationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	res, err := prosimoClient.GetIPREP(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("enabled", res.IpReps.Enabled)
	d.Set("allowlist", res.IpReps.AllowList)
	return diags
}

func resourceIPReputationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	ipRep := &client.IP_Reputation{
		Enabled:   false,
		AllowList: []string{},
	}
	_, err := prosimoClient.PutIPREP(ctx, ipRep)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Disabled ip_reputation")
	d.SetId("")
	return resourceIPReputationRead(ctx, d, meta)

}
