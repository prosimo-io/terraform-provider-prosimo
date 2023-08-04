package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePrivateLinkSource() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Private Link Sources.",
		CreateContext: resourcePVSCreate,
		UpdateContext: resourcePVSUpdate,
		DeleteContext: resourcePVSDelete,
		ReadContext:   resourcePVSRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of Private Link Source",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"cloud_creds_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cloud account under which application is hosted",
			},
			"cloud_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EX: us-west-2, eu-east-1",
			},
			"cloud_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_network": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of source VPC/VNET.",
									},
								},
							},
						},
						"subnets": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Subnet Details",
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourcePVSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	cloudSourceList := []client.Cloud_Source{}
	cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, d.Get("cloud_creds_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if v, ok := d.GetOk("cloud_sources"); ok {
		for i, _ := range v.([]interface{}) {
			cloudSourceConfig := v.([]interface{})[i].(map[string]interface{})
			cloudSource := &client.Cloud_Source{}
			if v, ok := cloudSourceConfig["cloud_network"].(*schema.Set); ok && v.Len() > 0 {
				cloudNetworkConfig := v.List()[0].(map[string]interface{})
				inVPC := cloudNetworkConfig["name"].(string)
				searchInput := &client.Search_Input{
					CloudCredsID: cloudCreds.ID,
					Region:       d.Get("cloud_region").(string),
				}
				res, err := prosimoClient.DiscoverPVSNetworks(ctx, searchInput)
				if err != nil {
					return diag.FromErr(err)
				}
				for _, cloudNet := range res.CloudNetworks {
					if cloudNet.Name == inVPC {
						cloudSource.CloudNetwork = &cloudNet
						break
					}
				}
			}
			if v, ok := cloudSourceConfig["subnets"].([]interface{}); ok && len(v) > 0 {
				inSubList := []client.Subnet{}
				for i, _ := range v {
					subnetConfig := v[i].(map[string]interface{})
					insub := client.Subnet{
						Cidr: subnetConfig["cidr"].(string),
					}
					inSubList = append(inSubList, insub)
					log.Println("insub", inSubList, insub)
				}
				cloudSource.Subnets = &inSubList
			}
			cloudSourceList = append(cloudSourceList, *cloudSource)
		}
	}
	inpvs := &client.PL_Source{
		Name:         d.Get("name").(string),
		CloudCredsID: cloudCreds.ID,
		Region:       d.Get("cloud_region").(string),
		CloudSources: &cloudSourceList,
	}
	postres, err := prosimoClient.CreatePrivateLinkSource(ctx, inpvs)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] New Private Link Source with id  %s is deployed", postres.PLSource_ResID.ID)
	d.SetId(postres.PLSource_ResID.ID)
	resourcePVSRead(ctx, d, meta)
	return diags
}

func resourcePVSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	updateReq := false
	if d.HasChange("name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Name",
			Detail:   "Name can't be modified",
		})
		return diags
	}
	if d.HasChange("cloud_creds_name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Cloud Credentials",
			Detail:   "Cloud Account details can't be modified",
		})
		return diags
	}

	if d.HasChange("cloud_region") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Cloud Regions",
			Detail:   "Cloud Region can't be modified",
		})
		return diags
	}

	if d.HasChange("cloud_sources") && !d.IsNewResource() {
		updateReq = true
	}
	if updateReq {
		cloudSourceList := []client.Cloud_Source{}
		cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, d.Get("cloud_creds_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if v, ok := d.GetOk("cloud_sources"); ok {
			for i, _ := range v.([]interface{}) {
				cloudSourceConfig := v.([]interface{})[i].(map[string]interface{})
				cloudSource := &client.Cloud_Source{}
				if v, ok := cloudSourceConfig["cloud_network"].(*schema.Set); ok && v.Len() > 0 {
					cloudNetworkConfig := v.List()[0].(map[string]interface{})
					inVPC := cloudNetworkConfig["name"].(string)
					searchInput := &client.Search_Input{
						CloudCredsID: cloudCreds.ID,
						Region:       d.Get("cloud_region").(string),
					}
					res, err := prosimoClient.DiscoverPVSNetworks(ctx, searchInput)
					if err != nil {
						return diag.FromErr(err)
					}
					for _, cloudNet := range res.CloudNetworks {
						if cloudNet.Name == inVPC {
							cloudSource.CloudNetwork = &cloudNet
							break
						}
					}
				}
				if v, ok := cloudSourceConfig["subnets"].([]interface{}); ok && len(v) > 0 {
					inSubList := []client.Subnet{}
					for i, _ := range v {
						subnetConfig := v[i].(map[string]interface{})
						insub := client.Subnet{
							Cidr: subnetConfig["cidr"].(string),
						}
						inSubList = append(inSubList, insub)
						log.Println("insub", inSubList, insub)
					}
					cloudSource.Subnets = &inSubList
				}
				cloudSourceList = append(cloudSourceList, *cloudSource)
			}
		}
		inpvs := &client.PL_Source{
			ID:           d.Id(),
			Name:         d.Get("name").(string),
			CloudCredsID: cloudCreds.ID,
			Region:       d.Get("cloud_region").(string),
			CloudSources: &cloudSourceList,
		}
		postres, err := prosimoClient.UpdatePrivateLinkSource(ctx, inpvs)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO]  Private Link Source with id  %s is updated", postres.PLSource_ResID.ID)
	}
	resourcePVSRead(ctx, d, meta)
	return diags
}

func resourcePVSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	pvsID := d.Id()

	log.Printf("[DEBUG] Get Private Link Source with id  %s", pvsID)

	pvs, err := prosimoClient.GetPrivateLinkSourceByID(ctx, pvsID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", pvs.ID)
	d.Set("name", pvs.Name)
	d.Set("cloud_region", pvs.Region)

	return diags
}

func resourcePVSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pvsID := d.Id()

	_, err := prosimoClient.DeletePrivateLinkSource(ctx, pvsID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Deleted Private Link Source with - id - %s", pvsID)
	d.SetId("")

	return diags
}
