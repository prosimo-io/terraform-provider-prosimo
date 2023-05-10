package prosimo

import (
	"context"
	"reflect"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIPAddresses() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on available ip ranges.",
		ReadContext: dataSourceIPAddressesRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_pools": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudtype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"totalsubnets": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"subnetsinuse": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIPAddressesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnIPPoolList []*client.IPPool

	ipPoolList, err := prosimoClient.GetIPPool(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	filter := d.Get("filter").(string)
	if filter != "" {
		for _, filteredIPPoolList := range ipPoolList.IPPools {
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredIPPoolList))
			if diags != nil {
				return diags
			}
			if flag {
				returnIPPoolList = append(returnIPPoolList, filteredIPPoolList)
			}
		}
	} else {
		for _, filteredIPPoolList := range ipPoolList.IPPools {
			returnIPPoolList = append(returnIPPoolList, filteredIPPoolList)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	ipPoolItems := flattenIpPoolsItemsData(returnIPPoolList)
	d.Set("ip_pools", ipPoolItems)

	return diags
}

func flattenIpPoolsItemsData(IpPoolItems []*client.IPPool) []interface{} {
	if IpPoolItems != nil {
		ois := make([]interface{}, len(IpPoolItems), len(IpPoolItems))

		for i, IpPoolItem := range IpPoolItems {
			oi := make(map[string]interface{})

			oi["name"] = IpPoolItem.Name
			oi["cidr"] = IpPoolItem.Cidr
			oi["cloudtype"] = IpPoolItem.CloudType
			oi["subnetsinuse"] = IpPoolItem.SubnetsInUse
			oi["totalsubnets"] = IpPoolItem.TotalSubnets

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
