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

func resourceNamespace() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify/Delete namespace policies and assignment of networks.",
		CreateContext: resourceNSCreate,
		UpdateContext: resourceNSUpdate,
		DeleteContext: resourceNSDelete,
		ReadContext:   resourceNSRead,
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
			"assign": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Assign the network to the namespace",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_networks": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							ForceNew:    true,
							Description: "Name of the networks to be assigned to the namespace",
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

func resourceNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	inns := &client.Namespace{
		Name: d.Get("name").(string),
	}
	postres, err := prosimoClient.CreateNamespace(ctx, inns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(postres.NamespaceResponse.ID)
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[DEBUG] Waiting for task id %s to complete", postres.NamespaceResponse.TaskID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, postres.NamespaceResponse.TaskID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] task %s is successful", postres.NamespaceResponse.TaskID)
	}

	log.Printf("[INFO] New Namespace with id  %s is created", postres.NamespaceResponse.ID)
	if v, ok := d.GetOk("assign"); ok {
		assigninput := v.(*schema.Set).List()[0].(map[string]interface{})
		inassignList := []client.NetActNamespace{}
		if v, ok := assigninput["source_networks"]; ok {
			sourceNetworkList := expandStringList(v.([]interface{}))
			for _, sourceNW := range sourceNetworkList {
				onboardNetworkList, err := prosimoClient.SearchOnboardNetworks(ctx)
				if err != nil {
					return diag.FromErr(err)
				}
				for _, network := range onboardNetworkList.Data.Records {
					if network.Name == sourceNW {
						inassign := client.NetActNamespace{
							ID: network.ID,
						}
						inassignList = append(inassignList, inassign)
						break
					}
				}
			}
			res, err := prosimoClient.AssignNetworkToNamespace(ctx, &inassignList, d.Id())
			if err != nil {
				return diag.FromErr(err)
			}
			if d.Get("wait_for_rollout").(bool) {
				log.Printf("[DEBUG] Waiting for task id %s to complete", res.NamespaceResponse.TaskID)
				err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
					retryUntilTaskComplete(ctx, d, meta, res.NamespaceResponse.TaskID))
				if err != nil {
					return diag.FromErr(err)
				}
				log.Printf("[INFO] task %s is successful", res.NamespaceResponse.TaskID)
			}
		}

	}
	resourceNSRead(ctx, d, meta)
	return diags
}

func resourceNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	if d.HasChange("name") && !d.IsNewResource() {
		inns := &client.Namespace{
			Name: d.Get("name").(string),
		}
		postres, err := prosimoClient.UpdateNamespace(ctx, inns)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[DEBUG] Waiting for task id %s to complete", postres.NamespaceResponse.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete(ctx, d, meta, postres.NamespaceResponse.TaskID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", postres.NamespaceResponse.TaskID)
		}
		log.Printf("[INFO] New Namespace with id  %s is updated", postres.NamespaceResponse.ID)
	}

	if d.HasChange("assign") && !d.IsNewResource() {
		if v, ok := d.GetOk("assign"); ok {
			assigninput := v.(*schema.Set).List()[0].(map[string]interface{})
			inassignList := []client.NetActNamespace{}
			if v, ok := assigninput["source_networks"]; ok {
				sourceNetworkList := expandStringList(v.([]interface{}))
				for _, sourceNW := range sourceNetworkList {
					onboardNetworkList, err := prosimoClient.SearchOnboardNetworks(ctx)
					if err != nil {
						return diag.FromErr(err)
					}
					for _, network := range onboardNetworkList.Data.Records {
						if network.Name == sourceNW {
							inassign := client.NetActNamespace{
								ID: network.ID,
							}
							inassignList = append(inassignList, inassign)
							break
						}
					}
				}
				res, err := prosimoClient.AssignNetworkToNamespace(ctx, &inassignList, d.Id())
				if err != nil {
					return diag.FromErr(err)
				}
				if d.Get("wait_for_rollout").(bool) {
					log.Printf("[DEBUG] Waiting for task id %s to complete", res.NamespaceResponse.TaskID)
					err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
						retryUntilTaskComplete(ctx, d, meta, res.NamespaceResponse.TaskID))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[INFO] task %s is successful", res.NamespaceResponse.TaskID)
				}
			}
		}
	}
	resourcePVSRead(ctx, d, meta)
	return diags
}

func resourceNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	nsID := d.Id()

	res, err := prosimoClient.DeleteNamespace(ctx, nsID)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[INFO] Waiting for task id %s to complete", res.NamespaceResponse.TaskID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, res.NamespaceResponse.TaskID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] task %s is successful", res.NamespaceResponse.TaskID)
	}
	log.Printf("[DEBUG] Deleted namespace Mapping with - id - %s", nsID)
	d.SetId("")

	return diags
}
