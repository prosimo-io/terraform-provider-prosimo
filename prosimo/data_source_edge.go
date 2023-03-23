package prosimo

import (
	"context"
	"fmt"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourceEdge() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on existing edges.",
		ReadContext: dataSourceEdgeRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edges": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloudregion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudtype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clustername": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pappfqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"regstatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"teamid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEdgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnEdgeList []*client.Edge

	edgeList, err := prosimoClient.GetEdge(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	filter := d.Get("filter").(string)
	fmt.Println("filter:", filter)
	if filter != "" {
		for _, returnEdge := range edgeList.Edges {
			filteredMap := map[string]interface{}{}
			err := mapstructure.Decode(returnEdge, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, filteredMap)
			if diags != nil {
				return diags
			}
			if flag {
				returnEdgeList = append(returnEdgeList, returnEdge)
			}
		}
		if len(returnEdgeList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})

			return diags
		}
	} else {
		for _, returnEdge := range edgeList.Edges {
			returnEdgeList = append(returnEdgeList, returnEdge)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	edgeItems := flattenEdgeItemsData(returnEdgeList)
	d.Set("edges", edgeItems)
	return diags
}

func flattenEdgeItemsData(EdgeItems []*client.Edge) []interface{} {
	if EdgeItems != nil {
		ois := make([]interface{}, len(EdgeItems), len(EdgeItems))

		for i, EdgeItem := range EdgeItems {
			oi := make(map[string]interface{})

			oi["cloudtype"] = EdgeItem.CloudType
			oi["cloudregion"] = EdgeItem.CloudRegion
			oi["clustername"] = EdgeItem.ClusterName
			oi["pappfqdn"] = EdgeItem.PappFqdn
			oi["regstatus"] = EdgeItem.RegStatus
			oi["status"] = EdgeItem.Status
			oi["subnet"] = EdgeItem.Subnet
			oi["teamid"] = EdgeItem.TeamID

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
