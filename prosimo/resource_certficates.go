package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCertficate() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify certficates.",
		CreateContext: resourceCertficateCreate,
		ReadContext:   resourceCertficateRead,
		DeleteContext: resourceCertficateDelete,
		UpdateContext: resourceCertficateUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certficate Option: e.g Domain, Client and CA",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of certficate, e.g Valid, Invalid",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Common Name of the certificate",
			},
			"upload_domain_cert": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Certficate option for uploading Custom Domain Certificate",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path to the certficate",
						},
						"private_key_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Path to the private key",
						},
					},
				},
			},
			"upload_client_cert": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Certficate option for uploading Source Certificate",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path to the certificate",
						},
						"private_key_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Path to the private key",
						},
					},
				},
			},
			"upload_ca_cert": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Certficate option for uploading Certificate Authority (CA)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path to the  CA certificate",
						},
					},
				},
			},
		},
	}
}

func resourceCertficateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	// uploadCert := &client.UploadCert{}

	if v, ok := d.GetOk("upload_domain_cert"); ok {
		uploadCert := &client.UploadCert{}
		uploadCertDetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := uploadCertDetails["cert_path"]; ok {
			Cert_path := v.(string)
			uploadCert.CertPath = Cert_path
		}
		if v, ok := uploadCertDetails["private_key_path"]; ok {
			Key_path := v.(string)
			uploadCert.KeyPath = Key_path
		}
		log.Printf("[DEBUG] uploading domian certificate")
		uploadCertResponseData, err := prosimoClient.UploadCert(ctx, uploadCert.CertPath, uploadCert.KeyPath)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(uploadCertResponseData.ResourceData.ID)
	}

	if v, ok := d.GetOk("upload_client_cert"); ok {
		uploadCert := &client.UploadCert{}
		uploadCertDetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := uploadCertDetails["cert_path"]; ok {
			Cert_path := v.(string)
			uploadCert.CertPath = Cert_path
		}
		if v, ok := uploadCertDetails["private_key_path"]; ok {
			Key_path := v.(string)
			uploadCert.KeyPath = Key_path
		}
		log.Printf("[DEBUG] uploading client certificate")
		uploadCertResponseData, err := prosimoClient.UploadCertClient(ctx, uploadCert.CertPath, uploadCert.KeyPath)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(uploadCertResponseData.ResourceData.ID)
	}
	if v, ok := d.GetOk("upload_ca_cert"); ok {
		uploadCert := &client.UploadCert{}
		uploadCertDetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := uploadCertDetails["cert_path"]; ok {
			Cert_path := v.(string)
			uploadCert.CertPath = Cert_path
		}
		log.Printf("[DEBUG] uploading ca certificate")
		uploadCertResponseData, err := prosimoClient.UploadCertCA(ctx, uploadCert.CertPath)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(uploadCertResponseData.ResourceData.ID)
	}

	return resourceCertficateRead(ctx, d, meta)

}

func resourceCertficateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	certID := d.Id()
	var diags diag.Diagnostics

	res, err := prosimoClient.GetCertDetails(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, certDetails := range res {
		if certDetails.ID == certID {
			d.Set("team_id", certDetails.TeamID)
			d.Set("type", certDetails.Type)
			d.Set("status", certDetails.Status)
			d.Set("url", certDetails.URL)
		}
	}

	return diags
}

func resourceCertficateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)
	certID := d.Id()

	if v, ok := d.GetOk("upload_domain_cert"); ok {
		uploadCert := &client.UploadCert{}
		uploadCertDetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := uploadCertDetails["cert_path"]; ok {
			Cert_path := v.(string)
			uploadCert.CertPath = Cert_path
		}
		if v, ok := uploadCertDetails["private_key_path"]; ok {
			Key_path := v.(string)
			uploadCert.KeyPath = Key_path
		}
		log.Printf("[DEBUG] uploading domain certificate")
		_, err := prosimoClient.UploadCertUpdate(ctx, uploadCert.CertPath, uploadCert.KeyPath, certID)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if v, ok := d.GetOk("upload_client_cert"); ok {
		uploadCert := &client.UploadCert{}
		uploadCertDetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := uploadCertDetails["cert_path"]; ok {
			Cert_path := v.(string)
			uploadCert.CertPath = Cert_path
		}
		if v, ok := uploadCertDetails["private_key_path"]; ok {
			Key_path := v.(string)
			uploadCert.KeyPath = Key_path
		}
		log.Printf("[DEBUG] uploading client certificate")
		_, err := prosimoClient.UploadCertClientUpdate(ctx, uploadCert.CertPath, uploadCert.KeyPath, certID)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if v, ok := d.GetOk("upload_ca_cert"); ok {
		uploadCert := &client.UploadCert{}
		uploadCertDetails := v.(*schema.Set).List()[0].(map[string]interface{})
		if v, ok := uploadCertDetails["cert_path"]; ok {
			Cert_path := v.(string)
			uploadCert.CertPath = Cert_path
		}
		log.Printf("[DEBUG] uploading ca certificate")
		_, err := prosimoClient.UploadCertCAUpdate(ctx, uploadCert.CertPath, certID)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	return resourceCertficateRead(ctx, d, meta)
}

func resourceCertficateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	certID := d.Id()

	err := prosimoClient.DeleteCert(ctx, certID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Deleted Certficate")
	d.SetId("")
	return diags
}
