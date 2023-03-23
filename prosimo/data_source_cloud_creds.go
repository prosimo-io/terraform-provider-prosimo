package prosimo

import (
	"context"
	"fmt"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/mapstructure"
)

func dataSourceCloudCreds() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on existing cloud credentials.",
		ReadContext: dataSourceCloudCredsRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = app_access_type==agent",
			},
			"cloud_creds": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloudtype": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetCloudTypes(), false),
						},
						"nickname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"aws": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"preferred_auth": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetAWSAuthTypes(), false),
									},
									"iam_role": {
										Type:     schema.TypeSet,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"role_arn": {
													Type:     schema.TypeString,
													Required: true,
												},
												"external_id": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"access_keys": {
										Type:     schema.TypeSet,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"access_key_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"secret_key_id": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"azure": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subscription_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tenant_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secret_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"gcp": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_account": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_email": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_uri": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"token_uri": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_provider_x509_cert_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_x509_cert_url": {
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

func dataSourceCloudCredsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	var returnCloudCredList []*client.CloudCreds

	cloudCredsList, err := prosimoClient.GetCloudCreds(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)
	if filter != "" {
		for _, returnedCloudCreds := range cloudCredsList.CloudCreds {
			filteredMap := map[string]interface{}{}

			err := mapstructure.Decode(returnedCloudCreds, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, filteredMap)
			if diags != nil {
				return diags
			}
			if flag {
				returnCloudCredList = append(returnCloudCredList, returnedCloudCreds)
			}
		}
		if len(returnCloudCredList) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})

			return diags
		}
	} else {
		for _, filteredCreds := range cloudCredsList.CloudCreds {
			returnCloudCredList = append(returnCloudCredList, filteredCreds)
		}
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if len(returnCloudCredList) > 0 {
		cloudCredItems := flattenCloudCredItemsData(returnCloudCredList)
		d.SetId(time.Now().Format(time.RFC850))
		d.Set("cloud_creds", cloudCredItems)
	}
	return diags

}

func flattenCloudCredItemsData(CloudCredItems []*client.CloudCreds) []interface{} {
	if CloudCredItems != nil {
		ois := make([]interface{}, len(CloudCredItems), len(CloudCredItems))

		for i, cloudCredItem := range CloudCredItems {
			oi := make(map[string]interface{})

			oi["nickname"] = cloudCredItem.Nickname
			oi["cloudtype"] = cloudCredItem.CloudType
			switch cloudCredItem.CloudType {
			case client.AWSCloudType:
				aws := getAWSCreds(cloudCredItem)
				oi["aws"] = aws
			case client.AzureCloudType:
				azure := getAzureCreds(cloudCredItem)
				oi["azure"] = azure
			case client.GCPCloudType:
				if cloudCredItem.CloudCredsDetails != nil {
					gcp := getGCPCreds(cloudCredItem, "")
					oi["gcp"] = gcp
				}
			}

			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
