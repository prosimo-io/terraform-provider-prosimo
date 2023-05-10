package prosimo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func datasourceCertficate() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on existing certificates.",
		ReadContext: dataSourceCertificateRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom filters to scope specific results. Usage: filter = app_access_type==agent",
			},
			"certs": {
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
						"ca": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generatedby": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isteamcert": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"notified": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"issuetime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expirytime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"san": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificatehash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signingalgorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"publickeyalgorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keysize": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pstatus": {
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
		},
	}
}

func dataSourceCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	var returnedCerts []client.ReadCertDetails

	getCerts, err := prosimoClient.GetCertDetails(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	filter := d.Get("filter").(string)
	if filter != "" {
		for _, certDetails := range getCerts {
			// filteredMap := map[string]interface{}{}
			var filteredMap *client.AppOnboardSettings

			err := mapstructure.Decode(onboardApp, &filteredMap)
			if err != nil {
				panic(err)
			}
			diags, flag := checkMainOperand(filter, reflect.ValueOf(filteredMap))
			if diags != nil {
				return diags
			}
			if flag {
				returnedCerts = append(returnedCerts, *&certDetails)
			}
		}
		if len(returnedCerts) == 0 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No match for input attribute",
				Detail:   fmt.Sprintln("No match for input attribute"),
			})

			return diags
		}
	} else {
		for _, certDetails := range getCerts {
			returnedCerts = append(returnedCerts, certDetails)
		}
	}
	d.SetId(time.Now().Format(time.RFC850))
	certItems := flattenCertItemsData(returnedCerts)
	d.Set("certs", certItems)
	return diags
}

func flattenCertItemsData(CertItems []client.ReadCertDetails) []interface{} {
	if CertItems != nil {
		ois := make([]interface{}, len(CertItems), len(CertItems))

		for i, CertItem := range CertItems {
			oi := make(map[string]interface{})

			oi["id"] = CertItem.ID
			oi["teamid"] = CertItem.TeamID
			oi["ca"] = CertItem.CA
			oi["dn"] = CertItem.DN
			oi["generatedby"] = CertItem.Generatedby
			oi["isteamcert"] = CertItem.ISTeamCert
			oi["notified"] = CertItem.Notified
			oi["issuetime"] = CertItem.Issuetime
			oi["expirytime"] = CertItem.Expirytime
			oi["san"] = CertItem.SAN
			oi["certificatehash"] = CertItem.Certificatehash
			oi["signingalgorithm"] = CertItem.Signingalgorithm
			oi["publickeyalgorithm"] = CertItem.Publickeyalgorithm
			oi["keysize"] = CertItem.Keysize
			oi["certificate"] = CertItem.Certificate
			oi["createdtime"] = CertItem.Createdtime
			oi["updatedtime"] = CertItem.Updatedtime

			ois[i] = oi
		}
		return ois
	}
	return make([]interface{}, 0)
}
