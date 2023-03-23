package prosimo

import (
	"context"
	"fmt"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIPAddresses() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify ip ranges.",
		CreateContext: resourceIPAddressesCreate,
		ReadContext:   resourceIPAddressesRead,
		DeleteContext: resourceIPAddressesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CIDR range",
			},
			"cloud_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(client.GetCloudTypes(), false),
				Description:  "Cloud Service Provider, e.g: AWS, AZURE, GCP ",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of cloud account",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_subnets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of Generated Subnets",
			},
			"subnets_in_use": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Subnets in use",
			},
		},
	}
}

func resourceIPAddressesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	ipPool := &client.IPPool{}

	if v, ok := d.GetOk("cidr"); ok {
		cidr := v.(string)
		ipPool.Cidr = cidr
	}

	if v, ok := d.GetOk("cloud_type"); ok {
		cloudType := v.(string)
		ipPool.CloudType = cloudType
	}
	log.Printf("[DEBUG] Creating IP Addresses block: %s (%s)", ipPool.Cidr, ipPool.CloudType)

	createdIPPool, err := prosimoClient.CreateIPPool(ctx, ipPool)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Created IP Addresses block: %s (%s) with ID (%s)", ipPool.Cidr, ipPool.CloudType, createdIPPool.IPPool.ID)
	d.SetId(createdIPPool.IPPool.ID)

	return resourceIPAddressesRead(ctx, d, meta)
}

func resourceIPAddressesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	ipPoolID := d.Id()

	log.Printf("[DEBUG] Get IP pool for %s", ipPoolID)

	ipPoolList, err := prosimoClient.GetIPPool(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipPool *client.IPPool
	for _, returnedIPPool := range ipPoolList.IPPools {
		if returnedIPPool.ID == ipPoolID {
			ipPool = returnedIPPool
			break
		}
	}
	if ipPool == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get IP Address",
			Detail:   fmt.Sprintf("Unable to find IP Address for ID %s", ipPoolID),
		})

		return diags
	}

	d.Set("cidr", ipPool.Cidr)
	d.Set("cloud_type", ipPool.CloudType)
	d.Set("name", ipPool.Name)
	d.Set("id", ipPool.ID)
	d.Set("total_subnets", ipPool.TotalSubnets)
	d.Set("subnets_in_use", ipPool.SubnetsInUse)

	log.Printf("[DEBUG] IP pool for %s - %+v", ipPoolID, ipPool)

	return diags
}

func resourceIPAddressesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	log.Printf("[IMP]Please make sure the ippool is not used by any existing edges before deleting")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	ipPoolID := d.Id()

	err := prosimoClient.DeleteIPPool(ctx, ipPoolID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
