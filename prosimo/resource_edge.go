package prosimo

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEdge() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify edges.",
		CreateContext: resourceEdgeCreate,
		UpdateContext: resourceEdgeCreate,
		DeleteContext: resourceEdgeDelete,
		ReadContext:   resourceEdgeRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cloud_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Cloud Account",
			},
			"cloud_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud Region",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the Cloud Service Provider, e.g: AWS, AZURE, GCP",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the edge cluster",
			},
			"papp_fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "pappFqdn URL",
			},
			"reg_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Deployment Status",
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet Range",
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wait_for_rollout": {
				Type:        schema.TypeBool,
				Description: "Wait for the rollout of the task to complete. Defaults to true.",
				Default:     true,
				Optional:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceEdgeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	cloudName := d.Get("cloud_name").(string)
	region := d.Get("cloud_region").(string)

	// validate cloud name
	cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, cloudName)
	if err != nil {
		return diag.FromErr(err)
	}

	if cloudCreds == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Cloud Credentials",
			Detail:   fmt.Sprintf("Unable to find Cloud Credentials for name %s", cloudName),
		})

		return diags
	}

	// validate cloud region
	regionExists, err := prosimoClient.CheckIfCloudRegionExists(ctx, cloudCreds.ID, region)
	if err != nil {
		return diag.FromErr(err)
	}

	if !regionExists {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Cloud Region",
			Detail:   fmt.Sprintf("Unable to find region %s for Cloud Name %s", region, cloudName),
		})

		return diags
	}

	// create edge
	edge := &client.Edge{
		CloudKeyID:  cloudCreds.ID,
		CloudRegion: region,
	}
	edgeResponseData, err := prosimoClient.CreateEdge(ctx, edge)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Created Edge for cloud name - %s, region - (%s)", cloudName, region)
	d.SetId(edgeResponseData.ResourceData.ID)

	// deploy app
	deployAppEdge := &client.Edge{
		ID: edgeResponseData.ResourceData.ID,
	}
	appResponseData, err := prosimoClient.DeployApp(ctx, deployAppEdge)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[INFO] Waiting for task id %s to complete", appResponseData.ResourceData.ID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, appResponseData.ResourceData.ID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] task %s is successful", appResponseData.ResourceData.ID)
	}
	log.Printf("[DEBUG] Deployed App for Edge - cloud name - %s, region - (%s)", cloudName, region)

	resourceCloudCredentialsRead(ctx, d, meta)

	return diags
}

func resourceEdgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	edgeID := d.Id()

	log.Printf("[DEBUG] Get Edge for %s", edgeID)

	edgeList, err := prosimoClient.GetEdge(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var edge *client.Edge
	for _, returnedEdge := range edgeList.Edges {
		if returnedEdge.ID == edgeID {
			edge = returnedEdge
			break
		}
	}
	if edge == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Edge",
			Detail:   fmt.Sprintf("Unable to find Edge for ID %s", edgeID),
		})

		return diags
	}

	// get cloud name for cloud key id
	cloudCreds, err := prosimoClient.GetCloudCredsById(ctx, edge.CloudKeyID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("cloud_name", cloudCreds.Nickname)
	d.Set("cloud_region", edge.CloudRegion)
	d.Set("id", edge.ID)
	d.Set("cloud_type", edge.CloudType)
	d.Set("cluster_name", edge.ClusterName)
	d.Set("papp_fqdn", edge.PappFqdn)
	d.Set("reg_status", edge.RegStatus)
	d.Set("status", edge.Status)
	d.Set("subnet", edge.Subnet)
	d.Set("team_id", edge.TeamID)

	return diags
}

func resourceEdgeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	edgeID := d.Id()

	appResponseData, err := prosimoClient.DeleteAppDeployment(ctx, edgeID)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[INFO] Waiting for task id %s to complete", appResponseData.ResourceData.ID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, appResponseData.ResourceData.ID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] task %s is successful", appResponseData.ResourceData.ID)
	}

	ret_err := prosimoClient.DeleteEdge(ctx, edgeID)
	if ret_err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
