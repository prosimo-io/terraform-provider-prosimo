package prosimo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func dataSourcePrivateLink() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on private links.",
		ReadContext: dataSourcePrivateLinkRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = region==us-west-1",
			},
			"pvt_link_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total Number of configured/onboarded Private Links",
			},
			"private_links": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"inuse": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"totalcount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloudcredsid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"credentials": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"credentials": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cloudsources": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"deleted": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cloudnetwork": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"subnets": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cidr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"deleted": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"endpoints": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"endpoint": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"domainid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"domain": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"policyid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"appid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"appname": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"status": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"protoports": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"protocol": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"port": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"websocketenabled": {
																			Type:     schema.TypeBool,
																			Computed: true,
																		},
																		"paths": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"portlist": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"createdtime": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"updatedtime": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"hostedzones": {
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
												"cloudcredentialsid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hostedzoneid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"vpcid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"prosimomanaged": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"createdtime": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"updatedtime": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"records": {
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
												"record": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hostedzoneid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"target": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ttl": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"createdtime": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"updatedtime": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"createdtime": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"updatedtime": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"createdtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updatedtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ports": {
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
									"edgeid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"plsid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domainid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"appid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"createdtime": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"updatedtime": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"edge": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloudcredsid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
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
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"policies": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourcePrivateLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	var returnPvtLinkList []client.PL_Source

	PrivateLinkList, err := prosimoClient.GetPrivateLinkSource(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)

	if filter != "" {
		for _, privateLink := range PrivateLinkList {
			var filteredMap *client.PL_Source

			err := mapstructure.Decode(privateLink, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))
			if diags != nil {
				return diags
			}
			if flag {
				returnPvtLinkList = append(returnPvtLinkList, *privateLink)
			}
		}
		if len(returnPvtLinkList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})
			return diags
		}
	} else {
		for _, privateLink := range PrivateLinkList {
			returnPvtLinkList = append(returnPvtLinkList, *privateLink)
		}
	}

	d.SetId(time.Now().Format(time.RFC850))
	pvtLinkItems := flattenPrivateLinkItemsData(returnPvtLinkList)
	d.Set("private_links", pvtLinkItems)
	d.Set("pvt_link_count", len(returnPvtLinkList))
	return diags
}

func flattenPrivateLinkItemsData(PrivateLinkItems []client.PL_Source) []interface{} {
	if PrivateLinkItems != nil {
		ois := make([]interface{}, len(PrivateLinkItems), len(PrivateLinkItems))

		for i, PvtLinkItem := range PrivateLinkItems {
			oi := make(map[string]interface{})
			oi["id"] = PvtLinkItem.ID
			oi["name"] = PvtLinkItem.Name
			oi["region"] = PvtLinkItem.Region
			oi["cloudcredsid"] = PvtLinkItem.CloudCredsID
			oi["deleted"] = PvtLinkItem.Deleted
			oi["inuse"] = PvtLinkItem.InUse
			oi["totalcount"] = PvtLinkItem.TotalCount

			PvtLinkCredentials := make([]map[string]interface{}, 0)
			PvtLinkCredentialsTF := make(map[string]interface{})
			PvtLinkCredentialsTF["id"] = PvtLinkItem.Credentials.ID
			PvtLinkCredentialsTF["cloud"] = PvtLinkItem.Credentials.Cloud
			PvtLinkCredentialsTF["type"] = PvtLinkItem.Credentials.Type
			PvtLinkCredentialsTF["credentials"] = PvtLinkItem.Credentials.Credentials
			PvtLinkCredentials = append(PvtLinkCredentials, PvtLinkCredentialsTF)
			oi["credentials"] = PvtLinkCredentials

			oi["cloudsources"] = flattenCloudSourceItemsData(*PvtLinkItem.CloudSources)
			oi["ports"] = flattenPlPortsItemsData(*PvtLinkItem.Ports)

			PvtLinkEdge := make([]map[string]interface{}, 0)
			PvtLinkEdgeTF := make(map[string]interface{})
			PvtLinkEdgeTF["id"] = PvtLinkItem.Edge.ID
			PvtLinkEdgeTF["domain"] = PvtLinkItem.Edge.Domain
			PvtLinkEdgeTF["cloud"] = PvtLinkItem.Edge.Cloud
			PvtLinkEdgeTF["cloudcredsid"] = PvtLinkItem.Edge.CloudCredsID
			PvtLinkEdgeTF["name"] = PvtLinkItem.Edge.Name

			EdgeNetworkInfo := make([]map[string]interface{}, 0)
			EdgeNetworkInfoTF := make(map[string]interface{})
			EdgeNetworkInfoTF["vpcid"] = PvtLinkItem.Edge.NetworkInfo.VpcID
			EdgeNetworkInfoTF["ilbdns"] = PvtLinkItem.Edge.NetworkInfo.IlbDns
			EdgeNetworkInfo = append(EdgeNetworkInfo, EdgeNetworkInfoTF)
			PvtLinkEdgeTF["networkinfo"] = EdgeNetworkInfo

			PvtLinkEdgeTF["region"] = PvtLinkItem.Edge.Region
			PvtLinkEdge = append(PvtLinkEdge, PvtLinkEdgeTF)
			oi["edge"] = PvtLinkEdge

			oi["policies"] = flattenPlPoliciesItemsData(*PvtLinkItem.Policies)

			oi["createdtime"] = PvtLinkItem.CreatedTime
			oi["updatedtime"] = PvtLinkItem.UpdatedTime

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenCloudSourceItemsData(CloudSourceItems []client.Cloud_Source) []interface{} {
	if CloudSourceItems != nil {
		ois := make([]interface{}, len(CloudSourceItems), len(CloudSourceItems))

		for i, CloudSourceItem := range CloudSourceItems {
			oi := make(map[string]interface{})
			oi["id"] = CloudSourceItem.ID
			oi["deleted"] = CloudSourceItem.Deleted

			CloudNetwork := make([]map[string]interface{}, 0)
			CloudNetworkTF := make(map[string]interface{})
			CloudNetworkTF["id"] = CloudSourceItem.CloudNetwork.ID
			CloudNetworkTF["name"] = CloudSourceItem.CloudNetwork.Name
			CloudNetwork = append(CloudNetwork, CloudNetworkTF)
			oi["cloudnetwork"] = CloudNetwork

			oi["hostedzones"] = flattenHostedZoneItemsData(*CloudSourceItem.HostedZones)
			oi["records"] = flattenCloudSourceRecordItemsData(*CloudSourceItem.Records)
			oi["subnets"] = flattenCloudSourceSubnetItemsData(*CloudSourceItem.Subnets)

			oi["createdtime"] = CloudSourceItem.CreatedTime
			oi["updatedtime"] = CloudSourceItem.UpdatedTime

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenHostedZoneItemsData(HostedZoneItems []client.HostedZone) []interface{} {
	if HostedZoneItems != nil {
		ois := make([]interface{}, len(HostedZoneItems), len(HostedZoneItems))

		for i, HostedZoneItem := range HostedZoneItems {
			oi := make(map[string]interface{})
			oi["id"] = HostedZoneItem.ID
			oi["teamid"] = HostedZoneItem.TeamID
			oi["cloudcredentialsid"] = HostedZoneItem.CloudCredentialsID
			oi["hostedzoneid"] = HostedZoneItem.HostedZoneID
			oi["name"] = HostedZoneItem.Name
			oi["vpcid"] = HostedZoneItem.VpcID
			oi["region"] = HostedZoneItem.Region
			oi["prosimomanaged"] = HostedZoneItem.ProsimoManaged
			oi["status"] = HostedZoneItem.Status
			oi["createdtime"] = HostedZoneItem.CreatedTime
			oi["updatedtime"] = HostedZoneItem.UpdatedTime
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenCloudSourceRecordItemsData(CloudSourceRecordItems []client.Cloud_Source_Record) []interface{} {
	if CloudSourceRecordItems != nil {
		ois := make([]interface{}, len(CloudSourceRecordItems), len(CloudSourceRecordItems))

		for i, CloudSourceRecordItem := range CloudSourceRecordItems {
			oi := make(map[string]interface{})
			oi["id"] = CloudSourceRecordItem.ID
			oi["teamid"] = CloudSourceRecordItem.TeamID
			oi["record"] = CloudSourceRecordItem.Record
			oi["hostedzoneid"] = CloudSourceRecordItem.HostedZoneID
			oi["type"] = CloudSourceRecordItem.Type
			oi["target"] = CloudSourceRecordItem.Target
			oi["ttl"] = CloudSourceRecordItem.Ttl
			oi["status"] = CloudSourceRecordItem.Status
			oi["createdtime"] = CloudSourceRecordItem.CreatedTime
			oi["updatedtime"] = CloudSourceRecordItem.UpdatedTime
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenCloudSourceSubnetItemsData(CloudSourceSubnetItems []client.Subnet) []interface{} {
	if CloudSourceSubnetItems != nil {
		ois := make([]interface{}, len(CloudSourceSubnetItems), len(CloudSourceSubnetItems))

		for i, CloudSourceSubnetItem := range CloudSourceSubnetItems {
			oi := make(map[string]interface{})
			oi["id"] = CloudSourceSubnetItem.ID
			oi["cidr"] = CloudSourceSubnetItem.Cidr
			oi["deleted"] = CloudSourceSubnetItem.Deleted
			oi["endpoints"] = flattenPlEndpointItemsData(*CloudSourceSubnetItem.Endpoints)
			oi["createdtime"] = CloudSourceSubnetItem.CreatedTime
			oi["updatedtime"] = CloudSourceSubnetItem.UpdatedTime
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenPlEndpointItemsData(PlEndpointItems []client.PL_Endpoint) []interface{} {
	if PlEndpointItems != nil {
		ois := make([]interface{}, len(PlEndpointItems), len(PlEndpointItems))

		for i, PlEndpointItem := range PlEndpointItems {
			oi := make(map[string]interface{})
			oi["id"] = PlEndpointItem.ID
			oi["endpoint"] = PlEndpointItem.Endpoint
			oi["name"] = PlEndpointItem.Name
			oi["domainid"] = PlEndpointItem.DomainID
			oi["domain"] = PlEndpointItem.Domain
			oi["policyid"] = PlEndpointItem.PolicyID
			oi["appid"] = PlEndpointItem.AppID
			oi["appname"] = PlEndpointItem.AppName
			oi["status"] = PlEndpointItem.Status
			oi["protoports"] = flattenPlEndpointProtoPortItemsData(*PlEndpointItem.ProtoPorts)

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenPlEndpointProtoPortItemsData(PlEndpointProtoPortItems []client.PL_Endpoint_ProtoPort) []interface{} {
	if PlEndpointProtoPortItems != nil {
		ois := make([]interface{}, len(PlEndpointProtoPortItems), len(PlEndpointProtoPortItems))

		for i, PlEndpointProtoPortItem := range PlEndpointProtoPortItems {
			oi := make(map[string]interface{})
			oi["protocol"] = PlEndpointProtoPortItem.Protocol
			oi["port"] = PlEndpointProtoPortItem.Port
			oi["websocketenabled"] = PlEndpointProtoPortItem.WebSocketEnabled
			oi["paths"] = PlEndpointProtoPortItem.Paths
			oi["portlist"] = PlEndpointProtoPortItem.PortList

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenPlPortsItemsData(PlPortsItems []client.PL_Port) []interface{} {
	if PlPortsItems != nil {
		ois := make([]interface{}, len(PlPortsItems), len(PlPortsItems))

		for i, PlPortsItem := range PlPortsItems {
			oi := make(map[string]interface{})
			oi["id"] = PlPortsItem.ID
			oi["teamid"] = PlPortsItem.TeamID
			oi["edgeid"] = PlPortsItem.EdgeID
			oi["port"] = PlPortsItem.Port
			oi["status"] = PlPortsItem.Status
			oi["plsid"] = PlPortsItem.PlsID
			oi["domainid"] = PlPortsItem.DomainID
			oi["appid"] = PlPortsItem.AppID
			oi["createdtime"] = PlPortsItem.CretedTime
			oi["updatedtime"] = PlPortsItem.UpdatedTime

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}

func flattenPlPoliciesItemsData(PlPoliciesItems []client.PL_Policy) []interface{} {
	if PlPoliciesItems != nil {
		ois := make([]interface{}, len(PlPoliciesItems), len(PlPoliciesItems))

		for i, PlPoliciesItem := range PlPoliciesItems {
			oi := make(map[string]interface{})
			oi["id"] = PlPoliciesItem.ID
			oi["status"] = PlPoliciesItem.Status
			oi["target"] = PlPoliciesItem.Target

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
