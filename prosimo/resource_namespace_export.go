package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNamespaceExport() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to export/withdraw networks.",
		CreateContext: resourceNSECreate,
		UpdateContext: resourceNSECreate,
		DeleteContext: resourceNSEDelete,
		ReadContext:   resourceNSERead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the namespace",
			},
			"wait_for_rollout": {
				Type:        schema.TypeBool,
				Description: "Wait for the rollout of the task to complete. Defaults to true.",
				Default:     true,
				Optional:    true,
			},
			"export": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Export local networks to other namespace",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_network": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the network to be exported to other namespace",
						},
						"namespaces": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of namespaces where network would be exported",
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

func resourceNSECreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	nameSpace := d.Get("name").(string)
	namespaceRes, err := prosimoClient.GetNamespaceByName(ctx, nameSpace)
	if err != nil {
		return diag.FromErr(err)
	}
	nameSpaceID := namespaceRes.ID

	if v, ok := d.GetOk("export"); ok {
		inexportList := []client.NetActNamespace{}
		for _, value := range v.([]interface{}) {
			exportinput := value.(map[string]interface{})
			inExport := client.NetActNamespace{}
			if v, ok := exportinput["source_network"]; ok {
				sourceNetwork := v.(string)
				onboardNetworkList, err := prosimoClient.SearchOnboardNetworks(ctx)
				if err != nil {
					return diag.FromErr(err)
				}
				for _, network := range onboardNetworkList.Data.Records {

					if network.Name == sourceNetwork {
						inExport.NetworkID = network.ID
						break
					}
				}
			}
			if v, ok := exportinput["namespaces"]; ok {
				namespaceList := expandStringList(v.([]interface{}))
				innsList := []string{}
				for _, namespace := range namespaceList {
					namespaceRes, err := prosimoClient.GetNamespaceByName(ctx, namespace)
					if err != nil {
						return diag.FromErr(err)
					}
					innsList = append(innsList, namespaceRes.ID)
				}
				inExport.Namespaces = innsList
			}

			inexportList = append(inexportList, inExport)
		}
		expres, err := prosimoClient.ExportNetworkToNamespace(ctx, &inexportList, nameSpaceID)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[DEBUG] Waiting for task id %s to complete", expres.NamespaceResponse.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete(ctx, d, meta, expres.NamespaceResponse.TaskID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", expres.NamespaceResponse.TaskID)
		}
	}
	d.SetId(nameSpaceID)

	resourceNSERead(ctx, d, meta)
	return diags
}

func resourceNSERead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	nsID := d.Id()

	log.Printf("[DEBUG] Get namespace with id  %s", nsID)

	ns, err := prosimoClient.GetNamespaceByID(ctx, nsID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", ns.ID)

	return diags
}

func resourceNSEDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	nsID := d.Id()
	inexportList := []client.NetActNamespace{}
	expres, err := prosimoClient.ExportNetworkToNamespace(ctx, &inexportList, nsID)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[INFO] Waiting for task id %s to complete", expres.NamespaceResponse.TaskID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, expres.NamespaceResponse.TaskID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] task %s is successful", expres.NamespaceResponse.TaskID)
	}
	log.Printf("[DEBUG] Deleted namespace Mapping with - id - %s", nsID)
	d.SetId("")

	return diags
}
