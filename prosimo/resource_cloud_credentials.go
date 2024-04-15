package prosimo

import (
	"context"
	"fmt"
	"log"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloudCreds() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify cloud credentials.",
		CreateContext: resourceCloudCredentialsCreate,
		ReadContext:   resourceCloudCredentialsRead,
		UpdateContext: resourceCloudCredentialsUpdate,
		DeleteContext: resourceCloudCredentialsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(client.GetCloudTypes(), false),
				Description:  "Select Cloud Service Provider, e.g: AWS, AZURE, GCP",
			},
			"nickname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Nickname of the cloud credential",
			},
			"bulk_onboard": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag for bulk upload, set it to true if want to bulk onboard accounts: defaults to false ",
			},

			"aws": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "AWS options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preferred_auth": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(client.GetAWSAuthTypes(), false),
							Description:  "Select preferred Authorization option, e.g: IAM Role, Access Keys",
						},
						"iam_role": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "IAM Role options",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_arn": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Role ARN, ref: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html",
									},
									"external_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: " External ID, ref: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html",
									},
								},
							},
						},
						"access_keys": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Access Key options",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Acces Key ID, ref: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html",
									},
									"secret_key_id": {
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
										Description: "Secret Key ID",
									},
								},
							},
						},
					},
				},
			},
			"azure": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "AZURE options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subscription ID, ref: https://learn.microsoft.com/en-us/azure/azure-portal/get-subscription-tenant-id",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tenant ID",
						},
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secret_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"gcp": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "GCP options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_path": {
							Type:         schema.TypeString,
							Required:     true,
							DefaultFunc:  schema.EnvDefaultFunc("HTTPFILEUPLOAD_FILE_PATH", nil),
							Description:  "Path of GCP credential file to upload.",
							ValidateFunc: validateFilePath,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of Credential, e.g: Service Account",
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
							// Default:  "Project ID",
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
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
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
			"bulk": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Bulk Account Onboarding options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							DefaultFunc:  schema.EnvDefaultFunc("HTTPFILEUPLOAD_FILE_PATH", nil),
							Description:  "Path of GCP credential file to upload.",
							ValidateFunc: validateFilePath,
						},
						"key_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tenant ID",
						},
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"external_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	cloudCreds := &client.CloudCreds{}
	cloudCredsDetails := &client.CloudCredsDetails{}
	var filePath string

	if v, ok := d.GetOk("cloud_type"); ok {
		cloudType := v.(string)
		cloudCreds.CloudType = cloudType
	}

	if v, ok := d.GetOk("nickname"); ok {
		nickname := v.(string)
		cloudCreds.Nickname = nickname
	}

	if d.Get("bulk_onboard").(bool) {
		if cloudCreds.CloudType != client.AWSCloudType {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Bulk Onboarding is only supported for AWS at the moment",
				Detail:   "Bulk Onboarding is only supported for AWS at the moment",
			})

			return diags
		}
		cloudCreds.KeyType = client.AWSBulkKeyType
		if v, ok := d.GetOk("bulk"); ok {
			bulkConfig := v.([]interface{})
			bulk := bulkConfig[0].(map[string]interface{})
			cloudCreds.AccountID = bulk["account_id"].(string)
			cloudCreds.ExternalID = bulk["external_id"].(string)
			filePath = bulk["file_path"].(string)
		}
		log.Printf("[DEBUG] Creating bulk Cloud Credentials for %v", cloudCreds)
		createdCloudCreds, err := prosimoClient.CreateBulkAcct(ctx, cloudCreds, filePath)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(createdCloudCreds.CloudCreds.ID)
	} else {
		switch cloudCreds.CloudType {
		case client.AWSCloudType:

			if v, ok := d.GetOk("aws"); ok {
				// awsConfig := v.(*schema.Set).List()
				awsConfig := v.([]interface{})
				aws := awsConfig[0].(map[string]interface{})

				prefferedAuth := aws["preferred_auth"].(string)
				if prefferedAuth == client.AWSIAMRoleAuth {

					cloudCreds.KeyType = client.AWSIAMRoleAuth

					if v, ok := aws["iam_role"].([]interface{}); ok {
						// iamRoleConfig := aws["iam_role"].(*schema.Set).List()
						// iamRoleConfig := aws["iam_role"].([]interface{})
						iamRole := v[0].(map[string]interface{})
						cloudCredsDetails.IAMRoleArn = iamRole["role_arn"].(string)
						cloudCredsDetails.IAMExternalID = iamRole["external_id"].(string)

					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "AWS IAM Role preferred auth requires iam_role block",
							Detail:   "AWS IAM Role preferred auth requires iam_role block",
						})

						return diags
					}

				} else if prefferedAuth == client.AWSAccessKeyAuth {

					cloudCreds.KeyType = client.AWSAccessKeyAuth

					if v, ok := aws["access_keys"].([]interface{}); ok {
						// accessKeyConfig := aws["access_keys"].(*schema.Set).List()
						// accessKeyConfig := aws["access_keys"].([]interface{})
						accessKey := v[0].(map[string]interface{})
						cloudCredsDetails.AccessKeyID = accessKey["access_key_id"].(string)
						cloudCredsDetails.SecretKeyID = accessKey["secret_key_id"].(string)
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "AWS Access Keys preferred auth requires access_keys block",
							Detail:   "AWS Access Keys  preferred auth requires access_keys block",
						})

						return diags
					}

				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Not a type of preferred auth for AWS",
						Detail:   "Not a type of preferred auth for AWS",
					})

					return diags
				}

			}

		case client.AzureCloudType:
			cloudCreds.KeyType = client.AzureKeyType
			if v, ok := d.GetOk("azure"); ok {
				// azureConfig := v.(*schema.Set).List()
				azureConfig := v.([]interface{})
				azure := azureConfig[0].(map[string]interface{})

				cloudCredsDetails.SubscriptionID = azure["subscription_id"].(string)
				cloudCredsDetails.TenantID = azure["tenant_id"].(string)
				cloudCredsDetails.ClientID = azure["client_id"].(string)
				cloudCredsDetails.SecretID = azure["secret_id"].(string)

			}

		case client.GCPCloudType:
			cloudCreds.KeyType = client.GCPKeyType
			if v, ok := d.GetOk("gcp"); ok {
				gcpConfig := v.([]interface{})
				gcp := gcpConfig[0].(map[string]interface{})

				filePath = gcp["file_path"].(string)

			}

		}
		if filePath == "" {
			cloudCreds.CloudCredsDetails = cloudCredsDetails
			log.Printf("[DEBUG] Creating Cloud Credentials for %v", cloudCreds)
			createdCloudCreds, err := prosimoClient.CreateCloudCreds(ctx, cloudCreds)
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] Created Cloud Credentials for cloud type - %s, nickname - (%s), id - (%s)", cloudCreds.CloudType, cloudCreds.Nickname, createdCloudCreds.CloudCreds.ID)
			d.SetId(createdCloudCreds.CloudCreds.ID)
		} else {
			log.Printf("[DEBUG] Creating Cloud Credentials for %v", cloudCreds)
			createdCloudCreds, err := prosimoClient.UploadGcpCloudCreds(ctx, cloudCreds, filePath)
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] Created Cloud Credentials for cloud type - %s, nickname - (%s), id - (%s)", cloudCreds.CloudType, cloudCreds.Nickname, createdCloudCreds.CloudCreds.ID)
			d.SetId(createdCloudCreds.CloudCreds.ID)
		}
	}

	resourceCloudCredentialsRead(ctx, d, meta)

	return diags
}

func resourceCloudCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	cloudCredsID := d.Id()

	prosimoClient := meta.(*client.ProsimoClient)

	cloudType := d.Get("cloud_type").(string)

	cloudCreds := &client.CloudCreds{
		ID:        cloudCredsID,
		CloudType: cloudType,
	}
	cloudCredsDetails := &client.CloudCredsDetails{}
	var filePath string

	updateReq := false

	if d.HasChange("nickname") && !d.IsNewResource() {
		updateReq = true
	}
	nickname := d.Get("nickname").(string)
	cloudCreds.Nickname = nickname

	switch cloudType {
	case client.AWSCloudType:
		if v, ok := d.GetOk("aws"); ok {
			awsConfig := v.([]interface{})
			aws := awsConfig[0].(map[string]interface{})

			if d.HasChange("aws.0.preferred_auth") && !d.IsNewResource() {
				updateReq = true
			}
			prefferedAuth := aws["preferred_auth"].(string)
			if prefferedAuth == client.AWSIAMRoleAuth {
				cloudCreds.KeyType = client.AWSIAMRoleAuth

				if v, ok := aws["iam_role"]; ok {
					iamRoleConfig := v.([]interface{})
					iamRole := iamRoleConfig[0].(map[string]interface{})

					if d.HasChange("aws.0.iam_role.0.role_arn") && !d.IsNewResource() {
						updateReq = true
					}

					if d.HasChange("aws.0.iam_role.0.external_id") && !d.IsNewResource() {
						updateReq = true
					}

					cloudCredsDetails.IAMRoleArn = iamRole["role_arn"].(string)
					cloudCredsDetails.IAMExternalID = iamRole["external_id"].(string)
				}

			} else if prefferedAuth == client.AWSAccessKeyAuth {
				cloudCreds.KeyType = client.AWSAccessKeyAuth

				if v, ok := aws["access_keys"]; ok {
					accessKeyConfig := v.([]interface{})
					accessKey := accessKeyConfig[0].(map[string]interface{})

					if d.HasChange("aws.0.access_keys.0.access_key_id") && !d.IsNewResource() {
						updateReq = true
					}

					if d.HasChange("aws.0.access_keys.0.secret_key_id") && !d.IsNewResource() {
						updateReq = true
					}

					cloudCredsDetails.AccessKeyID = accessKey["access_key_id"].(string)
					cloudCredsDetails.SecretKeyID = accessKey["secret_key_id"].(string)
				}

			}

		}

	case client.AzureCloudType:
		cloudCreds.KeyType = client.AzureKeyType
		if v, ok := d.GetOk("azure"); ok {
			azureConfig := v.([]interface{})
			azure := azureConfig[0].(map[string]interface{})

			if d.HasChange("azure.0.subscription_id") && !d.IsNewResource() {
				updateReq = true
			}

			if d.HasChange("azure.0.tenant_id") && !d.IsNewResource() {
				updateReq = true
			}
			if d.HasChange("azure.0.client_id") && !d.IsNewResource() {
				updateReq = true
			}
			if d.HasChange("azure.0.secret_id") && !d.IsNewResource() {
				updateReq = true
			}

			cloudCredsDetails.SubscriptionID = azure["subscription_id"].(string)
			cloudCredsDetails.TenantID = azure["tenant_id"].(string)
			cloudCredsDetails.ClientID = azure["client_id"].(string)
			cloudCredsDetails.SecretID = azure["secret_id"].(string)

		}

	case client.GCPCloudType:
		cloudCreds.KeyType = client.GCPKeyType
		if v, ok := d.GetOk("gcp"); ok {
			gcpConfig := v.([]interface{})
			gcp := gcpConfig[0].(map[string]interface{})
			//if d.HasChange("gcp.0.file_path") && !d.IsNewResource() {
			updateReq = true
			//}
			filePath = gcp["file_path"].(string)
		}

	}

	log.Printf("[DEBUG] Updating Cloud Credentials for cloud type - %s, nickname - (%s)", cloudType, cloudCreds.Nickname)
	cloudCreds.CloudCredsDetails = cloudCredsDetails

	if updateReq {
		if cloudType == "PRIVATE" {
			getCloud, err := prosimoClient.GetCloudCreds(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, cloud := range getCloud.CloudCreds {
				if cloud.CloudType == "PRIVATE" {
					cloudCreds.ID = cloud.ID
					// log.Println("cloudCreds", cloudCreds)
					updateCloudCreds, err := prosimoClient.UpdateCloudCreds(ctx, cloudCreds)
					if err != nil {
						return diag.FromErr(err)
					}
					d.SetId(updateCloudCreds.CloudCreds.ID)
					log.Printf("[DEBUG] Updated Cloud Credentials for cloud type - %s, nickname - (%s), id - (%s)", cloudType, cloudCreds.Nickname, updateCloudCreds.CloudCreds.ID)
				}
			}
		} else if filePath != "" {
			log.Printf("[DEBUG] Updating Cloud Credentials for %v", cloudCreds)
			createdCloudCreds, err := prosimoClient.UpdateGcpCloudCreds(ctx, cloudCreds, filePath)
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[DEBUG] Created Cloud Credentials for cloud type - %s, nickname - (%s), id - (%s)", cloudCreds.CloudType, cloudCreds.Nickname, createdCloudCreds.CloudCreds.ID)
			d.SetId(createdCloudCreds.CloudCreds.ID)
		} else {
			createdCloudCreds, err := prosimoClient.CreateCloudCreds(ctx, cloudCreds)
			if err != nil {
				return diag.FromErr(err)
			}
			d.SetId(createdCloudCreds.CloudCreds.ID)
			log.Printf("[DEBUG] Updated Cloud Credentials for cloud type - %s, nickname - (%s), id - (%s)", cloudType, cloudCreds.Nickname, createdCloudCreds.CloudCreds.ID)
		}
	}

	return resourceCloudCredentialsRead(ctx, d, meta)
}

func resourceCloudCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	var filePath string

	var diags diag.Diagnostics

	cloudCredsID := d.Id()

	log.Printf("[DEBUG] Get Cloudcreds for %s", cloudCredsID)

	cloudCreds, err := prosimoClient.GetCloudCredsById(ctx, cloudCredsID)
	if err != nil {
		return diag.FromErr(err)
	}
	if cloudCreds == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Cloud credentials",
			Detail:   fmt.Sprintf("Unable to find Cloud credentials for ID %s", cloudCredsID),
		})

		return diags
	}

	d.Set("cloud_type", cloudCreds.CloudType)
	d.Set("nickname", cloudCreds.Nickname)

	switch cloudCreds.CloudType {
	case client.AWSCloudType:
		aws := getAWSCreds(cloudCreds)
		d.Set("aws", aws)
	case client.AzureCloudType:
		azure := getAzureCreds(cloudCreds)
		d.Set("azure", azure)
	case client.GCPCloudType:
		cloudCreds.KeyType = client.GCPKeyType
		if v, ok := d.GetOk("gcp"); ok {
			gcpConfig := v.([]interface{})
			gcp := gcpConfig[0].(map[string]interface{})

			filePath = gcp["file_path"].(string)
		}

		log.Println("[DEBUG] file_path", filePath)
		gcp := getGCPCreds(cloudCreds, filePath)

		d.Set("gcp", gcp)
	}

	log.Printf("[DEBUG] Cloud credentials for %s - %+v", cloudCredsID, cloudCreds)

	return diags
}

func resourceCloudCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	log.Printf("[IMP]Please make sure the credential is not used by any existing edges before deleting")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cloudCredsID := d.Id()

	err := prosimoClient.DeleteCloudCreds(ctx, cloudCredsID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func getAWSCreds(cloudCreds *client.CloudCreds) interface{} {

	var aws []interface{}

	cloudCredsDetails := cloudCreds.CloudCredsDetails
	if cloudCredsDetails.IAMRoleArn != "" && cloudCredsDetails.IAMExternalID != "" {
		iamRole := []interface{}{
			map[string]interface{}{
				"role_arn":    cloudCredsDetails.IAMRoleArn,
				"external_id": cloudCredsDetails.IAMExternalID,
			},
		}

		aws = []interface{}{
			map[string]interface{}{
				"preferred_auth": client.AWSIAMRoleAuth,
				"iam_role":       iamRole,
			},
		}
	} else if cloudCredsDetails.AccessKeyID != "" && cloudCredsDetails.SecretKeyID != "" {
		accessKeys := []interface{}{
			map[string]interface{}{
				"access_key_id": cloudCredsDetails.AccessKeyID,
				"secret_key_id": cloudCredsDetails.SecretKeyID,
			},
		}

		aws = []interface{}{
			map[string]interface{}{
				"preferred_auth": client.AWSAccessKeyAuth,
				"access_keys":    accessKeys,
			},
		}
	}

	return aws
}

func getAzureCreds(cloudCreds *client.CloudCreds) interface{} {

	var azure []interface{}

	cloudCredsDetails := cloudCreds.CloudCredsDetails

	azure = []interface{}{
		map[string]interface{}{
			"subscription_id": cloudCredsDetails.SubscriptionID,
			"tenant_id":       cloudCredsDetails.TenantID,
			"client_id":       cloudCredsDetails.ClientID,
			"secret_id":       cloudCredsDetails.SecretID,
		},
	}

	return azure
}

func getGCPCreds(cloudCreds *client.CloudCreds, filePath string) interface{} {

	var gcp []interface{}

	cloudCredsDetails := cloudCreds.CloudCredsDetails

	gcp = []interface{}{
		map[string]interface{}{
			"type":                        cloudCredsDetails.GcpType,
			"project_id":                  cloudCredsDetails.ProjectID,
			"client_email":                cloudCredsDetails.ClientEmail,
			"client_id":                   cloudCredsDetails.GcpClientID,
			"auth_uri":                    cloudCredsDetails.AuthURI,
			"token_uri":                   cloudCredsDetails.TokenURI,
			"auth_provider_x509_cert_url": cloudCredsDetails.AuthCertURL,
			"client_x509_cert_url":        cloudCredsDetails.ClientCertURL,
			"file_path":                   filePath,
		},
	}

	return gcp
}
