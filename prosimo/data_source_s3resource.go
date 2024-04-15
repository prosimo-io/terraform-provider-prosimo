package prosimo

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceS3bucket() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on existing s3 buckets.",
		ReadContext: dataSourceS3bucketRead,
		Schema: map[string]*schema.Schema{
			"input_cloud_cred_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Filter based upon bucket name",
			},
			"input_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Filter based upon cloud region",
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceS3bucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	nickName := d.Get("input_cloud_cred_name").(string)
	region := d.Get("input_region").(string)

	// validate cloud credentials by name
	cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, nickName)

	if err != nil {
		return diag.FromErr(err)
	}

	if cloudCreds == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Cloud Credentials",
			Detail:   fmt.Sprintf("Unable to find Cloud Credentials for Nickname %s", nickName),
		})

		return diags
	}

	// fetch the cloud ID for AWS
	cloudCredsID := cloudCreds.ID
	s3InputData := &client.S3Input{
		ID:     cloudCredsID,
		Region: region,
	}
	log.Println("cloudCredsID", cloudCreds.ID)

	S3Data, err := prosimoClient.ReadS3Credentionals(ctx, s3InputData)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC850))

	d.Set("data", S3Data)

	return diags

}
