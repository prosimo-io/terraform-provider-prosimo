package prosimo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudtype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudkeyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudregion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clustername": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clustertype": {
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
						"publicip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"privateip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodesizesettings": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bandwidthname": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instancetype": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bandwidthrange": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"max": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"appnames": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"appusedcount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"networknames": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"networkusedcount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"edgeconnectivitycount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"privatelinksourcenames": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"privatelinkusedcount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sharedservicenames": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sharedservicecount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"createdtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locationid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nickname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ranchertoken": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tokenactivated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"updatedtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wgexternalendpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wginternalendpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"networkinfo": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpcid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ilbdns": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"byoresourcedetails": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpcid": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
								},
							},
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fabricconnectinfo": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"teamid": {
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
									"haspublic": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"attachments": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"conntype": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"pappfqdn": {
													Type:     schema.TypeString,
													Computed: true,
													Optional: true,
												},
												"attachtype": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"weight": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
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
			fmt.Println("returnEdge", returnEdge)
			diags, flag := checkMainOperand(filter, reflect.ValueOf(returnEdge))
			if diags != nil {
				return diags
			}
			if flag {
				returnEdgeList = append(returnEdgeList, returnEdge)
			}
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

			oi["id"] = EdgeItem.ID
			oi["cloudtype"] = EdgeItem.CloudType
			oi["cloudkeyid"] = EdgeItem.CloudKeyID
			oi["cloudregion"] = EdgeItem.CloudRegion
			oi["clustername"] = EdgeItem.ClusterName
			oi["clustertype"] = EdgeItem.ClusterType
			oi["pappfqdn"] = EdgeItem.PappFqdn
			oi["regstatus"] = EdgeItem.RegStatus
			oi["status"] = EdgeItem.Status
			oi["subnet"] = EdgeItem.Subnet
			oi["teamid"] = EdgeItem.TeamID
			oi["publicip"] = EdgeItem.PublicIP
			oi["privateip"] = EdgeItem.PrivateIP

			connectorSettings := make([]map[string]interface{}, 0)
			connectorSettingTF := make(map[string]interface{})
			bandwithRange := make([]map[string]interface{}, 0)
			bandwithRangeTF := make(map[string]interface{})
			// nodeSizesettings := EdgeItem.NodeSizesettings
			if EdgeItem.NodeSizesettings != nil {
				connectorSettingTF["bandwidth"] = EdgeItem.NodeSizesettings.Bandwidth
				connectorSettingTF["bandwidthname"] = EdgeItem.NodeSizesettings.BandwidthName
				connectorSettingTF["instancetype"] = EdgeItem.NodeSizesettings.InstanceType
				bandwidthRange := EdgeItem.NodeSizesettings.BandwidthRange
				bandwithRangeTF["min"] = bandwidthRange.Min
				bandwithRangeTF["max"] = bandwidthRange.Max
				bandwithRange = append(bandwithRange, bandwithRangeTF)
				connectorSettingTF["bandwidthrange"] = bandwithRange
				connectorSettings = append(connectorSettings, connectorSettingTF)
				oi["nodesizesettings"] = connectorSettings
			}
			oi["appnames"] = EdgeItem.AppNames
			oi["appusedcount"] = EdgeItem.AppUsedCount
			oi["networknames"] = EdgeItem.NetworkNames
			oi["networkusedcount"] = EdgeItem.NetworkUsedCount
			oi["edgeconnectivitycount"] = EdgeItem.EdgeConnectivityCount
			oi["privatelinksourcenames"] = EdgeItem.PrivateLinkSourceNames
			oi["privatelinkusedcount"] = EdgeItem.PrivateLinkUsedCount
			oi["sharedservicenames"] = EdgeItem.SharedServiceNames
			oi["sharedservicecount"] = EdgeItem.SharedServiceCount
			oi["city"] = EdgeItem.City
			oi["country"] = EdgeItem.Country
			oi["createdtime"] = EdgeItem.CreatedTime
			oi["flavor"] = EdgeItem.Flavor
			oi["locationid"] = EdgeItem.LocationId
			oi["nickname"] = EdgeItem.NickName

			oi["ranchertoken"] = EdgeItem.RancherToken
			oi["token"] = EdgeItem.Token
			oi["tokenactivated"] = EdgeItem.TokenActivated
			oi["updatedtime"] = EdgeItem.UpdatedTime
			oi["wgexternalendpoint"] = EdgeItem.WgExternalEndpoint
			oi["wginternalendpoint"] = EdgeItem.WgInternalEndpoint

			networkInfo := make([]map[string]interface{}, 0)
			networkInfoTF := make(map[string]interface{})
			networkInfoTF["vpcid"] = EdgeItem.NetworkInfo.VpcId
			networkInfoTF["ilbdns"] = EdgeItem.NetworkInfo.IlbDns
			networkInfo = append(networkInfo, networkInfoTF)
			oi["networkinfo"] = networkInfo

			networkbyoResource := make([]map[string]interface{}, 0)
			networkbyoResourceTF := make(map[string]interface{})
			networkbyoResourceTF["vpcid"] = EdgeItem.Byoresource.VpcID
			networkbyoResource = append(networkbyoResource, networkbyoResourceTF)
			oi["byoresourcedetails"] = networkbyoResource
			oi["state"] = EdgeItem.State
			fabricConnectInfo := make([]map[string]interface{}, 0)
			fabricConnectInfoTF := make(map[string]interface{})
			fabricConnectInfoTF["attachments"] = flattenAttachmentItemsData(EdgeItem.FabricConnectInfo.Attachments)
			fabricConnectInfoTF["cloudtype"] = EdgeItem.FabricConnectInfo.CloudType
			fabricConnectInfoTF["haspublic"] = EdgeItem.FabricConnectInfo.HasPublic
			fabricConnectInfoTF["id"] = EdgeItem.FabricConnectInfo.ID
			fabricConnectInfoTF["name"] = EdgeItem.FabricConnectInfo.Name
			fabricConnectInfoTF["teamid"] = EdgeItem.FabricConnectInfo.TeamID
			fabricConnectInfo = append(fabricConnectInfo, fabricConnectInfoTF)
			oi["fabricconnectinfo"] = fabricConnectInfo

			ois[i] = oi
			// fmt.Printf("Anil ois state %d : %v", i, ois)
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenAttachmentItemsData(AttachmentItems []*client.Attachments) []interface{} {
	if AttachmentItems != nil {
		ois := make([]interface{}, len(AttachmentItems), len(AttachmentItems))

		for i, attachmentItem := range AttachmentItems {
			oi := make(map[interface{}]interface{})
			oi["id"] = attachmentItem.ID
			oi["conntype"] = attachmentItem.ConnType
			oi["attachtype"] = attachmentItem.AttachType
			oi["pappfqdn"] = attachmentItem.PappFqdn
			oi["weight"] = attachmentItem.Weight
			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}
