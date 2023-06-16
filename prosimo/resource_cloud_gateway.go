package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudGateway() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Cloud Gateway configuration.",
		CreateContext: resourceCloudGatewayCreate,
		UpdateContext: resourceCloudGatewayUpdate,
		DeleteContext: resourceCloudGatewayDelete,
		ReadContext:   resourceCloudGatewayRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud region",
			},
			"attach_point": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Details of attach points, (VPC, Transit Gateway etc)",
			},
			"attachment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Transit Gateway Attachment, applicable when attach_point is a TGW.",
			},
			"route_table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Transit Gateway routing table, applicable when attach_point is a TGW",
			},
		},
	}
}

func resourceCloudGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	cloudGWInput := &client.CloudGateway{}

	region := d.Get("region").(string)

	attachPoint := d.Get("attach_point").(string)

	connectivityOptions, err := prosimoClient.GetCnnectivityOptions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, connectivityOption := range connectivityOptions {
		if connectivityOption.Region == region && connectivityOption.AttachPoint == attachPoint {
			cloudGWInput.EdgeConnectivityID = connectivityOption.ID
			break
		}
	}

	if v, ok := d.GetOk("attachment"); ok {
		cloudGWInput.Attachment = v.(string)
	}

	if v, ok := d.GetOk("route_table"); ok {
		cloudGWInput.RouteTable = v.(string)
	}

	createCloudGateway, err := prosimoClient.CreateCloudGateway(ctx, cloudGWInput)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createCloudGateway.Data.ID)
	return resourceCloudGatewayRead(ctx, d, meta)
}

func resourceCloudGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)
	updateReq := false

	if d.HasChange("region") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("attach_point") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("attachment") && !d.IsNewResource() {
		updateReq = true
	}

	if d.HasChange("route_table") && !d.IsNewResource() {
		updateReq = true
	}
	if updateReq {
		cloudGWInput := &client.CloudGateway{ID: d.Id()}

		region := d.Get("region").(string)

		attachPoint := d.Get("attach_point").(string)

		connectivityOptions, err := prosimoClient.GetCnnectivityOptions(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, connectivityOption := range connectivityOptions {
			if connectivityOption.Region == region && connectivityOption.AttachPoint == attachPoint {
				cloudGWInput.EdgeConnectivityID = connectivityOption.ID
				break
			}
		}

		if v, ok := d.GetOk("attachment"); ok {
			cloudGWInput.Attachment = v.(string)
		}

		if v, ok := d.GetOk("route_table"); ok {
			cloudGWInput.RouteTable = v.(string)
		}

		err1 := prosimoClient.UpdateCloudGateway(ctx, cloudGWInput)
		if err1 != nil {
			return diag.FromErr(err1)
		}
		log.Printf("Cloud gateway with id %s updated", d.Id())
	}
	return resourceCloudGatewayRead(ctx, d, meta)

}

func resourceCloudGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	gwID := d.Id()

	gwRes, err := prosimoClient.GetCloudGatewayByID(ctx, gwID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("region", gwRes.Region)
	d.Set("attach_point", gwRes.AttachPoint)
	d.Set("attachment", gwRes.Attachment)
	d.Set("route_table", gwRes.RouteTable)
	return diags
}

func resourceCloudGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	cgwID := d.Id()

	err := prosimoClient.DeleteCloudGateway(ctx, cgwID)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
