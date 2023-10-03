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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEdge() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify edges.",
		CreateContext: resourceEdgeCreate,
		UpdateContext: resourceEdgeUpdate,
		DeleteContext: resourceEdgeDelete,
		ReadContext:   resourceEdgeRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"ip_range": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet Range",
			},
			"vpc_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "VPC Source: Available options: Prosimo/Existing, applicable only for AWS",
				ValidateFunc: validation.StringInSlice(client.GetVPCSourceOptions(), false),
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of source vpc, applicable when vpc_source is existing vpc",
			},
			"node_size_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  " Available Options: <1 Gbps, 1-5 Gbps, 5-10 Gbps, >10 Gbps",
							ValidateFunc: validation.StringInSlice(client.GetConnectorBandwidthOptions(), false),
						},
						// "bandwidth_name": {
						// 	Type:        schema.TypeString,
						// 	Required:    true,
						// 	Description: "EX: <1 Gbps, >1 Gbps",
						// },
						"instance_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(client.GetConnectorInstanceOptions(), false),
							Description: "Available Options wrt cloud and bandwidth :" +
								"Cloud_Provider: AWS:" +
								"Bandwidth:  <1 Gbps, Available Options: t3.medium/t3a.medium/c5.large" +
								"Bandwidth:  1-5 Gbps, Available Options: c5a.large/c5.xlarge/c5a.xlarge/c5n.xlarge" +
								"Bandwidth: 5-10 Gbps, Available Options: c5a.8xlarge/c5.9xlarge" +
								"Bandwidth: >10 Gbps, Available Options: c5n.9xlarge/c5a.16xlarge/c5.18xlarge/c5n.18xlarge",
						},
					},
				},
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
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deploy_edge": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like to deploy the edge ",
			},
			"decommission_edge": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like the edge to be decommissioned.",
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

	appOnboardFlag := d.Get("deploy_edge").(bool)
	appOffboardFlag := d.Get("decommission_edge").(bool)

	if appOffboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid  decommission_edge flag.",
			Detail:   "decommission_edge can't be set to true while creating edge resource.",
		})
		return diags
	}
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
		CloudType:   cloudCreds.CloudType,
		CloudRegion: region,
		Subnet:      d.Get("ip_range").(string),
	}
	//NodeSize Settings
	if cloudCreds.CloudType == client.AWSCloudType {
		if v, ok := d.GetOk("vpc_source"); ok {
			vpcSource := v.(string)
			if vpcSource == client.ExistingVPCSource {
				byoresource := &client.ByoResource{
					VpcID: d.Get("vpc_name").(string),
				}
				edge.Byoresource = byoresource
			}
		}
		if v, ok := d.GetOk("node_size_settings"); ok {
			nodesizesettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
			nodesizesettingInput := &client.ConnectorSettings{
				// Bandwidth:     nodesizesettingConfig["bandwidth"].(string),
				BandwidthName: nodesizesettingConfig["bandwidth"].(string),
				InstanceType:  nodesizesettingConfig["instance_type"].(string),
			}
			switch nodesizesettingInput.BandwidthName {
			case client.LessThan1GBPS:
				nodesizesettingInput.Bandwidth = client.ConnectorSizeSmall
			case client.OneToFiveGBPS:
				nodesizesettingInput.Bandwidth = client.ConnectorSizeMedium
			case client.FiveToTenGBPS:
				nodesizesettingInput.Bandwidth = client.ConnectorSizeLarge
			case client.MoreThanTenGBPS:
				nodesizesettingInput.Bandwidth = client.ConnectorSizeExtraLarge
			}
			edge.NodeSizesettings = nodesizesettingInput
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing Node Size options",
				Detail:   "Node Size options are required for edges to be spun up in aws..",
			})
			return diags
		}
	}
	edgeResponseData, err := prosimoClient.CreateEdge(ctx, edge)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Created Edge for cloud name - %s, region - (%s)", cloudName, region)
	d.SetId(edgeResponseData.ResourceData.ID)

	// deploy edge
	if appOnboardFlag {
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
	}

	resourceCloudCredentialsRead(ctx, d, meta)

	return diags
}

func resourceEdgeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	edgeID := d.Id()

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardFlag := d.Get("deploy_edge").(bool)
	appOffboardFlag := d.Get("decommission_edge").(bool)
	if appOnboardFlag && appOffboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid deploy_edge and decommission_edge flag combination.",
			Detail:   "Both deploy_edge and decommission_edge have been set to true.",
		})
		return diags
	}

	if d.HasChange("cloud_name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify Cloud Name",
			Detail:   "Cloud Name can't be modified",
		})
		return diags
	}

	if d.HasChange("cloud_region") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify Cloud Region",
			Detail:   "Cloud Region can't be modified",
		})
		return diags
	}

	if d.HasChange("vpc_source") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify VPC Source",
			Detail:   "VPC SOURCE can't be modified",
		})
		return diags
	}

	//Offboard app
	// decomEdge := false
	if d.HasChange("decommission_edge") && !d.IsNewResource() {
		isDecommission := d.Get("decommission_edge").(bool)
		if isDecommission {
			// decomEdge = true
			// deployAppEdge := &client.Edge{
			// 	ID: edgeResponseData.ResourceData.ID,
			// }
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
			log.Printf("[DEBUG] Successfully Decommissioned Edge")
		}
	}

	// Get Cloud Type

	cloudName := d.Get("cloud_name").(string)
	cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, cloudName)
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := false

	if d.HasChange("ip_range") && !d.IsNewResource() {
		updateReq = true
		patchSubnet := client.Edge{
			Subnet:    d.Get("ip_range").(string),
			CloudType: cloudCreds.CloudType,
		}
		err := prosimoClient.PatchSubnetRange(ctx, edgeID, &patchSubnet)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if cloudCreds.CloudType == client.AWSCloudType {
		if d.HasChange("node_size_settings") {
			updateReq = true
		}
	}

	if updateReq {
		// update edge
		edge := &client.Edge{
			CloudKeyID:  cloudCreds.ID,
			CloudRegion: d.Get("cloud_region").(string),
			Subnet:      d.Get("ip_range").(string),
		}
		//NodeSize Settings
		if cloudCreds.CloudType == client.AWSCloudType {
			if v, ok := d.GetOk("node_size_settings"); ok {
				nodesizesettingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
				nodesizesettingInput := &client.ConnectorSettings{
					// Bandwidth:     nodesizesettingConfig["bandwidth"].(string),
					BandwidthName: nodesizesettingConfig["bandwidth"].(string),
					InstanceType:  nodesizesettingConfig["instance_type"].(string),
				}
				switch nodesizesettingInput.BandwidthName {
				case client.LessThan1GBPS:
					nodesizesettingInput.Bandwidth = client.ConnectorSizeSmall
				case client.OneToFiveGBPS:
					nodesizesettingInput.Bandwidth = client.ConnectorSizeMedium
				case client.FiveToTenGBPS:
					nodesizesettingInput.Bandwidth = client.ConnectorSizeLarge
				case client.MoreThanTenGBPS:
					nodesizesettingInput.Bandwidth = client.ConnectorSizeExtraLarge
				}
				edge.NodeSizesettings = nodesizesettingInput
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Missing Node Size options",
					Detail:   "Node Size options are required for edges to be spun up in aws..",
				})
				return diags
			}
		}
		edgeResponseData, err := prosimoClient.UpdateEdge(ctx, edgeID, edge)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(edgeResponseData.ResourceData.ID)
	}

	// deploy edge
	if appOnboardFlag {
		deployAppEdge := &client.Edge{
			ID: edgeID,
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
	}

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
	d.Set("team_id", edge.TeamID)

	return diags
}

func resourceEdgeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	edgeID := d.Id()

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

	if edge.Status == "DEPLOYED" {

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
