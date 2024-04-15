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

func resourceIDP() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify identity providers.",
		CreateContext: resourceIDPCreate,
		UpdateContext: resourceIDPUpdate,
		DeleteContext: resourceIDPDelete,
		ReadContext:   resourceIDPRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"idp_account": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(client.GetIDPAccountTypes(), false),
				Description:  "Identity provider, choose among okta, azure_ad, one_login, ping-one, ping-federate, google, other",
			},
			"auth_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(client.GetIDPAuthTypes(), false),
				Description:  "Authentication type, e.g: oidc, saml",
			},
			"account_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IDP account url",
			},
			"api_cred_provided": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(client.GetAPICredProvided(), false),
			},
			"oidc": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Auth type OIDC options",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"api_client_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_secret_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"region": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(client.GetIDPPRegionTypes(), false),
							Description:  "Choose between eu, asia, us, default",
						},
						"env_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"admin_email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"customer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_file": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							DefaultFunc:  schema.EnvDefaultFunc("HTTPFILEUPLOAD_FILE_PATH", nil),
							Description:  descriptions["file_path"],
							ValidateFunc: validateFilePath,
						},
					},
				},
			},

			"saml": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metadata_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(client.GetIDPPRegionTypes(), false),
							Description:  "Choose between eu, asia, us, default",
						},
						"env_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_client_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_secret_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"admin_email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"customer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_file": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							DefaultFunc:  schema.EnvDefaultFunc("HTTPFILEUPLOAD_FILE_PATH", nil),
							Description:  descriptions["file_path"],
							ValidateFunc: validateFilePath,
						},
					},
				},
			},
			"select_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(client.GetIDPTypes(), false),
				Description:  "IDP type, e.g: primary, partner",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IDP status",
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IDP subnet range ",
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_file_updated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"api_creds_file_name": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"partner": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_domain": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of partner domain urls",
						},
						"apps": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of apps",
						},
					},
				},
			},
		},
	}
}

func resourceIDPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	idp := &client.IDP{}

	idpDetails := &client.IDPDetails{}

	if v, ok := d.GetOk("idp_account"); ok {
		idpAccount := v.(string)
		idp.IDPName = idpAccount
	}
	if v, ok := d.GetOk("auth_type"); ok {
		authType := v.(string)
		idp.Auth_Type = authType
	}
	if v, ok := d.GetOk("account_url"); ok {
		accountURL := v.(string)
		idp.AccountURL = accountURL
	}
	if v, ok := d.GetOk("api_cred_provided"); ok {
		apiCredProvided := v.(string)
		idp.Api_Creds_Provided = apiCredProvided
	}
	if v, ok := d.GetOk("select_type"); ok {
		selectType := v.(string)
		idp.Select_Type = selectType
		if idp.Select_Type == client.IDPPartnerType {
			if v, ok := d.GetOk("partner"); ok {
				partnerConfig := v.([]interface{})
				partner := partnerConfig[0].(map[string]interface{})
				if v, ok := partner["user_domain"]; ok {
					userDomainList := v.([]interface{})

					if len(userDomainList) > 0 {
						idp.DomainURL = expandStringList(v.([]interface{}))
					}
				}
				if v, ok := partner["apps"]; ok {
					appList := v.([]interface{})

					if len(appList) > 0 {
						idp.AppIDs = expandStringList(v.([]interface{}))
					}
				}
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid Partner Config.",
					Detail:   "Invalid Partner Config.",
				})
			}
		}
	}

	switch idp.IDPName {

	// okta account
	case client.IDPOktaAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // okta + oidc + apicredno

			if v, ok := d.GetOk("oidc"); ok {
				// oidcConfig := v.(*schema.Set).List()
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})
				if idp.Api_Creds_Provided == client.APICredYes {
					apiToken := oidc["api_token"].(string)

					if apiToken != "" {
						idpDetails.APIToken = apiToken
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for Okta IDP account, missing apitoken.",
							Detail:   "Invalid OIDC config for Okta IDP account, missing apitoken.",
						})
					}
				} else {
					secretID := oidc["secret_id"].(string)
					clientID := oidc["client_id"].(string)
					if secretID != "" && clientID != "" {
						idpDetails.ClientSecret = secretID
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for Okta IDP account, missing id/secret fields",
							Detail:   "Invalid OIDC config for Okta IDP account, missing id/secret fields",
						})
					}

				}
			}
		} else if idp.Auth_Type == client.IDPSAMLAuth { // okta + saml
			if v, ok := d.GetOk("saml"); ok {
				// samlConfig := v.(*schema.Set).List()
				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					apiToken := saml["api_token"].(string)
					if apiToken != "" {
						idpDetails.APIToken = apiToken
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for Okta IDP account, missing apitoken.",
							Detail:   "Invalid SAML config for Okta IDP account, missing apitoken.",
						})
					}

				} else {
					metadataURL := saml["metadata_url"].(string)
					metadata := saml["metadata"].(string)
					if metadataURL != "" && metadata != "" {
						idpDetails.MetadataURL = metadataURL
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for Okta IDP account, missing metadata.",
							Detail:   "Invalid SAML config for Okta IDP account, missing metadata.",
						})
					}
				}
			}
		}

	// azure_ad account
	case client.IDPAzureADAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // azure_ad + oidc

			// azure_ad + oidc - client_id, secret_id

			if v, ok := d.GetOk("oidc"); ok {
				// oidcConfig := v.(*schema.Set).List()
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})
				if idp.Api_Creds_Provided == client.APICredYes {
					apisecretID := oidc["api_secret_id"].(string)
					apiclientID := oidc["api_client_id"].(string)

					if apisecretID != "" && apiclientID != "" {
						idpDetails.APIClientSecret = apisecretID
						idpDetails.APIClientID = apiclientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for Azure AD IDP account, missing api client/secret field.",
							Detail:   "Invalid OIDC config for Azure AD IDP account, missing api client/secret field.",
						})
					}

				} else {
					secretID := oidc["secret_id"].(string)
					clientID := oidc["client_id"].(string)

					if secretID != "" && clientID != "" {
						idpDetails.ClientSecret = secretID
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for Azure AD IDP account, missing client/secret field.",
							Detail:   "Invalid OIDC config for Azure AD IDP account, missing client/secret field.",
						})
					}
				}
			}

		} else if idp.Auth_Type == client.IDPSAMLAuth { // azure_ad + saml
			// azure_ad + saml - metadata_url, metadata

			if v, ok := d.GetOk("saml"); ok {
				// samlConfig := v.(*schema.Set).List()
				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					apisecretID := saml["api_secret_id"].(string)
					apiclientID := saml["api_client_id"].(string)
					metadataURL := saml["metadata_url"].(string)
					metadata := saml["metadata"].(string)

					if apisecretID != "" && apiclientID != "" && metadataURL != "" && metadata != "" {
						idpDetails.APIClientSecret = apisecretID
						idpDetails.APIClientID = apiclientID
						idpDetails.MetadataURL = metadataURL
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for Azure AD IDP account",
							Detail:   "Invalid SAML config for Azure AD IDP account",
						})
					}
				} else {
					metadataURL := saml["metadata_url"].(string)
					metadata := saml["metadata"].(string)
					if metadataURL != "" && metadata != "" {
						idpDetails.MetadataURL = metadataURL
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for Azure AD IDP account",
							Detail:   "Invalid SAML config for Azure AD IDP account",
						})
					}

				}

			}

		}

	// one_login account
	case client.IDPOneLoginAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // one_login + oidc

			// one_login + oidc - client_id, secret_id, region

			if v, ok := d.GetOk("oidc"); ok {
				// oidcConfig := v.(*schema.Set).List()
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					apisecretID := oidc["api_secret_id"].(string)
					apiclientID := oidc["api_client_id"].(string)
					region := oidc["region"].(string)

					if apisecretID != "" && apiclientID != "" && region != "" {
						idpDetails.APIClientSecret = apisecretID
						idpDetails.APIClientID = apiclientID
						idpDetails.Region = region
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for OneLogin IDP account",
							Detail:   "Invalid OIDC config for OneLogin IDP account",
						})
					}
				} else {
					secretID := oidc["secret_id"].(string)
					clientID := oidc["client_id"].(string)

					if secretID != "" && clientID != "" {
						idpDetails.ClientSecret = secretID
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for OneLogin IDP account",
							Detail:   "Invalid OIDC config for OneLogin IDP account",
						})

					}

				}

			}

		} else if idp.Auth_Type == client.IDPSAMLAuth { // okta + saml
			// okta + saml - metadata_url, metadata

			if v, ok := d.GetOk("saml"); ok {
				// samlConfig := v.(*schema.Set).List()
				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					apisecretID := saml["api_secret_id"].(string)
					apiclientID := saml["api_client_id"].(string)
					region := saml["region"].(string)

					if apisecretID != "" && apiclientID != "" && region != "" {
						idpDetails.APIClientSecret = apisecretID
						idpDetails.APIClientID = apiclientID
						idpDetails.Region = region
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for OneLogin IDP account",
							Detail:   "Invalid OIDC config for OneLogin IDP account",
						})
					}
				} else {
					metadataURL := saml["metadata_url"].(string)
					metadata := saml["metadata"].(string)

					if metadataURL != "" && metadata != "" {
						idpDetails.MetadataURL = metadataURL
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for OneLogin IDP account",
							Detail:   "Invalid SAML config for OneLogin IDP account",
						})
					}

				}

			}

		}

	// other account
	case client.IDPOtherAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // other + oidc

			// other + oidc - client_id, secret_id, region

			if v, ok := d.GetOk("oidc"); ok {
				// oidcConfig := v.(*schema.Set).List()
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				secretID := oidc["secret_id"].(string)
				clientID := oidc["client_id"].(string)

				if secretID != "" && clientID != "" {
					idpDetails.ClientSecret = secretID
					idpDetails.ClientID = clientID
				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid OIDC config for Other IDP account",
						Detail:   "Invalid OIDC config for Other IDP account",
					})

				}

			}
		} else if idp.Auth_Type == client.IDPSAMLAuth { // other + saml
			// other + saml - metadata_url, metadata

			if v, ok := d.GetOk("saml"); ok {
				// samlConfig := v.(*schema.Set).List()
				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				metadataURL := saml["metadata_url"].(string)
				metadata := saml["metadata"].(string)

				if metadataURL != "" && metadata != "" {
					idpDetails.MetadataURL = metadataURL
					idpDetails.Metadata = metadata
				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid SAML config for other IDP account",
						Detail:   "Invalid SAML config for other IDP account",
					})

				}

			}

		}
	case client.IDPPingOneAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // pingone + oidc

			// pingone + oidc - client_id, secret_id, region, envid

			if v, ok := d.GetOk("oidc"); ok {
				// oidcConfig := v.(*schema.Set).List()
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					apisecretID := oidc["api_secret_id"].(string)
					apiclientID := oidc["api_client_id"].(string)
					region := oidc["region"].(string)
					envID := oidc["env_id"].(string)

					if apisecretID != "" && apiclientID != "" && region != "" && envID != "" {
						idpDetails.APIClientSecret = apisecretID
						idpDetails.APIClientID = apiclientID
						idpDetails.Region = region
						idpDetails.EnvID = envID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for PingOne IDP account",
							Detail:   "Invalid OIDC config for PingOne IDP account",
						})
					}
				} else {
					secretID := oidc["secret_id"].(string)
					clientID := oidc["client_id"].(string)
					if secretID != "" && clientID != "" {
						idpDetails.ClientSecret = secretID
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for PingOne IDP account",
							Detail:   "Invalid OIDC config for PingOne IDP account",
						})
					}

				}

			}
		} else if idp.Auth_Type == client.IDPSAMLAuth { // other + saml
			// other + saml - metadata_url, metadata

			if v, ok := d.GetOk("saml"); ok {
				// samlConfig := v.(*schema.Set).List()
				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					apisecretID := saml["api_secret_id"].(string)
					apiclientID := saml["api_client_id"].(string)
					region := saml["region"].(string)
					envID := saml["env_id"].(string)

					if apisecretID != "" && apiclientID != "" && region != "" && envID != "" {
						idpDetails.APIClientSecret = apisecretID
						idpDetails.APIClientID = apiclientID
						idpDetails.Region = region
						idpDetails.EnvID = envID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for PingOne IDP account",
							Detail:   "Invalid SAML config for PingOne IDP account",
						})
					}
				} else {
					metadataURL := saml["metadata_url"].(string)
					metadata := saml["metadata"].(string)

					if metadataURL != "" && metadata != "" {
						idpDetails.MetadataURL = metadataURL
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for PingOne IDP account",
							Detail:   "Invalid SAML config for PingOne IDP account",
						})
					}

				}

			}

		}
	case client.IDPGoogleWSAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // google + oidc

			// google + oidc - client_id, secret_id, region, envid

			if v, ok := d.GetOk("oidc"); ok {
				// oidcConfig := v.(*schema.Set).List()
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					secretID := oidc["secret_id"].(string)
					clientID := oidc["client_id"].(string)
					apiEmail := oidc["admin_email"].(string)
					customerID := oidc["customer_id"].(string)
					domain := oidc["domain"].(string)
					filePath := oidc["file_path"].(string)

					if secretID != "" && clientID != "" && apiEmail != "" && customerID != "" && domain != "" && filePath != "" {
						idpDetails.ClientSecret = secretID
						idpDetails.ClientID = clientID
						idpDetails.APIEmail = apiEmail
						idpDetails.CustomerID = customerID
						idpDetails.Domain = domain
						idpDetails.FilePath = filePath
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for Google IDP account",
							Detail:   "Invalid OIDC config for Google IDP account",
						})
					}
				} else {
					secretID := oidc["secret_id"].(string)
					clientID := oidc["client_id"].(string)

					if secretID != "" && clientID != "" {
						idpDetails.ClientSecret = secretID
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for Google IDP account",
							Detail:   "Invalid OIDC config for Google IDP account",
						})
					}
				}
			}
		} else if idp.Auth_Type == client.IDPSAMLAuth { // other + saml
			// other + saml - metadata_url, metadata

			if v, ok := d.GetOk("saml"); ok {
				// samlConfig := v.(*schema.Set).List()
				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					metadataURL := saml["metadata_url"].(string)
					metadata := saml["metadata"].(string)
					apiEmail := saml["admin_email"].(string)
					customerID := saml["customer_id"].(string)
					domain := saml["domain"].(string)
					filePath := saml["file_path"].(string)

					if metadataURL != "" && metadata != "" && apiEmail != "" && customerID != "" && domain != "" && filePath != "" {
						idpDetails.ClientSecret = metadataURL
						idpDetails.ClientID = metadata
						idpDetails.APIEmail = apiEmail
						idpDetails.CustomerID = customerID
						idpDetails.Domain = domain
						idpDetails.FilePath = filePath
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid SAML config for Google IDP account",
							Detail:   "Invalid SAML config for Google IDP account",
						})
					}
				} else {
					metadataURL := saml["metadata_url"].(string)
					metadata := saml["metadata"].(string)

					if metadataURL != "" && metadata != "" {
						idpDetails.ClientSecret = metadataURL
						idpDetails.ClientID = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid OIDC config for Google IDP account",
							Detail:   "Invalid OIDC config for Google IDP account",
						})
					}
				}
			}
		}

	}

	if len(diags) > 0 {
		return diags
	}
	idp.Details = *idpDetails
	if idp.IDPName != client.IDPGoogleWSAccount {
		log.Printf("[DEBUG] Creating IDP : %v", idp)
		createdIDP, err := prosimoClient.CreateIDP(ctx, idp)
		if err != nil {
			log.Printf("[ERROR] Error in creating IDP : %s, accountURL: (%s), idpDetails : (%v) - error %s", idp.IDPName, idp.AccountURL, idp.Details, err.Error())
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] Created IDP : %s, accountURL: (%s) with ID (%s)", idp.IDPName, idp.AccountURL, createdIDP.ResourceData.ID)
		d.SetId(createdIDP.ResourceData.ID)
	} else {
		log.Printf("[DEBUG] Creating IDP : %v", idp)
		createdIDP, err := prosimoClient.CreateGoogleIDP(ctx, idp)
		if err != nil {
			log.Printf("[ERROR] Error in creating IDP : %s, accountURL: (%s), idpDetails : (%v) - error %s", idp.IDPName, idp.AccountURL, idp.Details, err.Error())
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] Created IDP : %s, accountURL: (%s) with ID (%s)", idp.IDPName, idp.AccountURL, createdIDP.ResourceData.ID)
		d.SetId(createdIDP.ResourceData.ID)
	}
	return resourceIDPRead(ctx, d, meta)
}

func resourceIDPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	idpAccount := d.Get("idp_account").(string)

	idp := &client.IDP{
		IDPName: idpAccount,
	}

	updateReq := false

	idpDetails := &client.IDPDetails{}

	idpID := d.Id()
	idp.ID = idpID

	if d.HasChange("auth_type") && !d.IsNewResource() {
		updateReq = true
	}
	idp.Auth_Type = d.Get("auth_type").(string)

	if d.HasChange("account_url") && !d.IsNewResource() {
		updateReq = true
	}
	idp.AccountURL = d.Get("account_url").(string)

	if d.HasChange("api_cred_provided") && !d.IsNewResource() {
		updateReq = true
	}
	idp.Api_Creds_Provided = d.Get("api_cred_provided").(string)

	if d.HasChange("select_type") && !d.IsNewResource() {
		updateReq = true
	}
	idp.Select_Type = d.Get("select_type").(string)
	if idp.Select_Type == client.IDPPartnerType {
		if v, ok := d.GetOk("partner"); ok {
			partnerConfig := v.([]interface{})
			partner := partnerConfig[0].(map[string]interface{})

			// z := partner.(map[string]interface)
			// partner := partnerConfig[0].(map[string][]string)

			if v, ok := partner["user_domain"]; ok {
				userDomainList := v.([]interface{})
				if len(userDomainList) > 0 {
					for _, i := range userDomainList {
						if d.HasChange(i.(string)) && !d.IsNewResource() {
							updateReq = true
						}
					}
					// if len(userDomainList) > 0 {
					idp.DomainURL = expandStringList(v.([]interface{}))
				}
			}
			if v, ok := partner["apps"]; ok {
				appList := v.([]interface{})
				if len(appList) > 0 {
					for _, i := range appList {
						if d.HasChange(i.(string)) && !d.IsNewResource() {
							updateReq = true
						}
					}
					// if len(appList) > 0 {
					idp.AppIDs = expandStringList(v.([]interface{}))
				}
			}
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid Partner Config.",
				Detail:   "Invalid Partner Config.",
			})
		}
	}
	switch idpAccount {

	// okta account
	case client.IDPOktaAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // okta + oidc

			if v, ok := d.GetOk("oidc"); ok {
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					if apiToken, ok := oidc["api_token"].(string); ok { // okta + oidc + api creds - api_token
						if d.HasChange("oidc.0.api_token") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIToken = apiToken
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIToken is required for Okta IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APIToken is required for Okta IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}
				} else {
					// okta + oidc + no api creds - client_id, secret_id

					if secretID, ok := oidc["secret_id"].(string); ok {
						if d.HasChange("oidc.0.secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "SecretID is required for Okta IDP account with OIDC auth type",
							Detail:   "SecretID is required for Okta IDP account with OIDC auth type",
						})

					}

					if clientID, ok := oidc["client_id"].(string); ok {
						if d.HasChange("oidc.0.client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "ClientID is required for Okta IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "ClientID is required for Okta IDP account with OIDC auth type and API Creds Provided as NO",
						})

					}

				}

			}

		} else if idp.Auth_Type == client.IDPSAMLAuth { // okta + saml

			if v, ok := d.GetOk("saml"); ok {

				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					if apiToken, ok := saml["api_token"].(string); ok { // okta + saml + api creds - api_token
						if d.HasChange("saml.0.api_token") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIToken = apiToken
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIToken is required for Okta IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APIToken is required for Okta IDP account with SAML auth type and API Creds Provided as YES",
						})

					}
				} else {
					// okta + saml + no api creds - metadata_url, metadata

					if metadataURL, ok := saml["metadata_url"].(string); ok {
						if d.HasChange("saml.0.metadata_url") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.MetadataURL = metadataURL
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "MetadataURL is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "MetadataURL is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
						})

					}

					if metadata, ok := saml["metadata"].(string); ok {
						if d.HasChange("saml.0.metadata") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Metadata is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "Metadata is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
						})

					}

				}

			}

		}

	// azure_ad account
	case client.IDPAzureADAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // azure_ad + oidc

			if v, ok := d.GetOk("oidc"); ok {
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				// azure_ad + oidc - client_id, secret_id

				if idp.Api_Creds_Provided == client.APICredYes {
					if secretID, ok := oidc["api_secret_id"].(string); ok {
						if d.HasChange("oidc.0.api_secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APISecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APISecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}

					if clientID, ok := oidc["api_client_id"].(string); ok {
						if d.HasChange("oidc.0.api_client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APIClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}

				} else {
					if secretID, ok := oidc["secret_id"].(string); ok {
						if d.HasChange("oidc.0.secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "SecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "SecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
						})

					}

					if clientID, ok := oidc["client_id"].(string); ok {
						if d.HasChange("oidc.0.client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "ClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "ClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
						})

					}

				}
			}

		} else if idp.Auth_Type == client.IDPSAMLAuth { // azure_ad + saml
			// azure_ad + saml - metadata_url, metadata

			if v, ok := d.GetOk("saml"); ok {

				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					if metadataURL, ok := saml["metadata_url"].(string); ok {
						if d.HasChange("saml.0.metadata_url") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.MetadataURL = metadataURL
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "MetadataURL is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "MetadataURL is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
						})

					}

					if metadata, ok := saml["metadata"].(string); ok {
						if d.HasChange("saml.0.metadata") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Metadata is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "Metadata is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
						})

					}

					if secretID, ok := saml["api_secret_id"].(string); ok {
						if d.HasChange("saml.0.api_secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APISecretID is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APISecretID is required for Azure AD IDP account with SAMLauth type and API Creds Provided as YES",
						})

					}

					if clientID, ok := saml["api_client_id"].(string); ok {
						if d.HasChange("saml.0.api_client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIClientID is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APIClientID is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
						})

					}

				} else {
					if metadataURL, ok := saml["metadata_url"].(string); ok {
						if d.HasChange("saml.0.metadata_url") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.MetadataURL = metadataURL
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "MetadataURL is required for Azure AD IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "MetadataURL is required for Azure AD IDP account with SAML auth type and API Creds Provided as NO",
						})

					}

					if metadata, ok := saml["metadata"].(string); ok {
						if d.HasChange("saml.0.metadata") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Metadata is required for Azure AD IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "Metadata is required for Azure AD IDP account with SAML auth type and API Creds Provided as NO",
						})

					}
				}

			}
		}

	// one_login account
	case client.IDPOneLoginAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // one_login + oidc

			if v, ok := d.GetOk("oidc"); ok {
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				// one_login + oidc - client_id, secret_id, region

				if idp.Api_Creds_Provided == client.APICredYes {
					if secretID, ok := oidc["api_secret_id"].(string); ok {
						if d.HasChange("oidc.0.api_secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APISecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APISecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}

					if clientID, ok := oidc["api_client_id"].(string); ok {
						if d.HasChange("oidc.0.api_client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APIClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}

					if region, ok := oidc["region"].(string); ok {
						if d.HasChange("oidc.0.region") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = region
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Region is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "Region is required for Azure AD IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}

				} else {
					if secretID, ok := oidc["secret_id"].(string); ok {
						if d.HasChange("oidc.0.secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "SecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "SecretID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
						})
					}

					if clientID, ok := oidc["client_id"].(string); ok {
						if d.HasChange("oidc.0.client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "ClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "ClientID is required for Azure AD IDP account with OIDC auth type and API Creds Provided as NO",
						})
					}
				}

			}

		} else if idp.Auth_Type == client.IDPSAMLAuth { // one_login + saml

			if v, ok := d.GetOk("saml"); ok {

				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				// one_login + saml - metadata_url, metadata

				if idp.Api_Creds_Provided == client.APICredYes {
					if secretID, ok := saml["api_secret_id"].(string); ok {
						if d.HasChange("saml.0.api_secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APISecretID is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APISecretID is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
						})
					}

					if clientID, ok := saml["api_client_id"].(string); ok {
						if d.HasChange("saml.0.api_client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIClientID is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APIClientID is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
						})
					}

					if region, ok := saml["region"].(string); ok {
						if d.HasChange("saml.0.region") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = region
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Region is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "Region is required for Azure AD IDP account with SAML auth type and API Creds Provided as YES",
						})
					}
				} else {
					if metadataURL, ok := saml["metadata_url"].(string); ok {
						if d.HasChange("saml.0.metadata_url") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.MetadataURL = metadataURL
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "MetadataURL is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "MetadataURL is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
						})

					}

					if metadata, ok := saml["metadata"].(string); ok {
						if d.HasChange("saml.0.metadata") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Metadata is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "Metadata is required for Okta IDP account with SAML auth type and API Creds Provided as NO",
						})

					}

				}
			}

		}
	case client.IDPPingOneAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // pingone + oidc

			if v, ok := d.GetOk("oidc"); ok {
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				// one_login + oidc - client_id, secret_id, region

				if idp.Api_Creds_Provided == client.APICredYes {
					if secretID, ok := oidc["api_secret_id"].(string); ok {
						if d.HasChange("oidc.0.api_secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APISecretID is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APISecretID is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}

					if clientID, ok := oidc["api_client_id"].(string); ok {
						if d.HasChange("oidc.0.api_client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIClientID is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APIClientID is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}

					if region, ok := oidc["region"].(string); ok {
						if d.HasChange("oidc.0.region") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Region = region
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Region is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "Region is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}

					if envID, ok := oidc["env_id"].(string); ok {
						if d.HasChange("oidc.0.env_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.EnvID = envID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "EnvID is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "EnvID is required for Ping One IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}
				} else {
					if secretID, ok := oidc["secret_id"].(string); ok {
						if d.HasChange("oidc.0.secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "SecretID is required for Ping One IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "SecretID is required for Ping One IDP account with OIDC auth type and API Creds Provided as NO",
						})

					}

					if clientID, ok := oidc["client_id"].(string); ok {
						if d.HasChange("oidc.0.client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "ClientID is required for Ping One IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "ClientID is required for Ping One IDP account with OIDC auth type and API Creds Provided as NO",
						})

					}
				}

			}

		} else if idp.Auth_Type == client.IDPSAMLAuth { // one_login + saml

			if v, ok := d.GetOk("saml"); ok {

				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				// one_login + saml - metadata_url, metadata

				if idp.Api_Creds_Provided == client.APICredYes {
					if secretID, ok := saml["api_secret_id"].(string); ok {
						if d.HasChange("saml.0.api_secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APISecretID is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APISecretID is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
						})

					}

					if clientID, ok := saml["api_client_id"].(string); ok {
						if d.HasChange("saml.0.api_client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIClientID is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APIClientID is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
						})

					}

					if region, ok := saml["region"].(string); ok {
						if d.HasChange("saml.0.region") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Region = region
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Region is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "Region is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
						})

					}

					if envID, ok := saml["env_id"].(string); ok {
						if d.HasChange("saml.0.env_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.EnvID = envID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "EnvID is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "EnvID is required for Ping One IDP account with SAML auth type and API Creds Provided as YES",
						})

					}
				} else {
					if metadataURL, ok := saml["metadata_url"].(string); ok {
						if d.HasChange("saml.0.metadata_url") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.MetadataURL = metadataURL
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "MetadataURL is required for Ping One IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "MetadataURL is required for Ping One IDP account with SAML auth type and API Creds Provided as NO",
						})

					}

					if metadata, ok := saml["metadata"].(string); ok {
						if d.HasChange("saml.0.metadata") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Metadata is required for Ping One IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "Metadata is required for Ping One IDP account with SAML auth type and API Creds Provided as NO",
						})

					}
				}
			}

		}
	case client.IDPGoogleWSAccount:
		if idp.Auth_Type == client.IDPOIDCAuth { // pingone + oidc

			if v, ok := d.GetOk("oidc"); ok {
				oidcConfig := v.([]interface{})
				oidc := oidcConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					if filePath, ok := oidc["file_path"].(string); ok {
						if d.HasChange("oidc.0.file_path") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.FilePath = filePath
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "FilePath is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "FilePath is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
						})

					}
					if secretID, ok := oidc["secret_id"].(string); ok {
						if d.HasChange("oidc.0.secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "SecretID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "SecretID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}

					if clientID, ok := oidc["client_id"].(string); ok {
						if d.HasChange("oidc.0.client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "ClientID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "ClientID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}
					if apiEmail, ok := oidc["admin_email"].(string); ok {
						if d.HasChange("oidc.0.admin_email") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIEmail = apiEmail
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIEMAIL is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "APIEMAIL is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}
					if customerID, ok := oidc["customer_id"].(string); ok {
						if d.HasChange("oidc.0.customer_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.CustomerID = customerID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "CustomerID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "CustomerID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}
					if domain, ok := oidc["domain"].(string); ok {
						if d.HasChange("oidc.0.domain") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Domain = domain
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Domain is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
							Detail:   "Domain is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as YES",
						})
					}
				} else {
					if secretID, ok := oidc["secret_id"].(string); ok {
						if d.HasChange("oidc.0.secret_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientSecret = secretID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "SecretID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "SecretID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as NO",
						})
					}

					if clientID, ok := oidc["client_id"].(string); ok {
						if d.HasChange("oidc.0.client_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.ClientID = clientID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "ClientID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as NO",
							Detail:   "ClientID is required for Google Workspace IDP account with OIDC auth type and API Creds Provided as NO",
						})
					}
				}
			}
		} else if idp.Auth_Type == client.IDPSAMLAuth {
			// google + saml - metadata_url, metadata
			if v, ok := d.GetOk("saml"); ok {
				// samlConfig := v.(*schema.Set).List()
				samlConfig := v.([]interface{})
				saml := samlConfig[0].(map[string]interface{})

				if idp.Api_Creds_Provided == client.APICredYes {
					if filePath, ok := saml["file_path"].(string); ok {
						if d.HasChange("saml.0.file_path") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.FilePath = filePath
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "FilePath is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "FilePath is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
						})

					}
					if apiEmail, ok := saml["admin_email"].(string); ok {
						if d.HasChange("saml.0.admin_email") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.APIEmail = apiEmail
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "APIEMAIL is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "APIEMAIL is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
						})
					}
					if customerID, ok := saml["customer_id"].(string); ok {
						if d.HasChange("saml.0.customer_id") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.CustomerID = customerID
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "CustomerID is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "CustomerID is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
						})
					}
					if domain, ok := saml["domain"].(string); ok {
						if d.HasChange("saml.0.domain") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Domain = domain
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Domain is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "Domain is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
						})
					}
					if metadataURL, ok := saml["metadata_url"].(string); ok {
						if d.HasChange("saml.0.metadata_url") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.MetadataURL = metadataURL
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "MetadataURL is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "MetadataURL is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
						})

					}

					if metadata, ok := saml["metadata"].(string); ok {
						if d.HasChange("saml.0.metadata") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Metadata is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
							Detail:   "Metadata is required for Google Workspace IDP account with SAML auth type and API Creds Provided as YES",
						})

					}
				} else {
					if metadataURL, ok := saml["metadata_url"].(string); ok {
						if d.HasChange("saml.0.metadata_url") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.MetadataURL = metadataURL
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "MetadataURL is required for Google Workspace IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "MetadataURL is required for Google Workspace IDP account with SAML auth type and API Creds Provided as NO",
						})

					}

					if metadata, ok := saml["metadata"].(string); ok {
						if d.HasChange("saml.0.metadata") && !d.IsNewResource() {
							updateReq = true
						}
						idpDetails.Metadata = metadata
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Metadata is required for Google Workspace IDP account with SAML auth type and API Creds Provided as NO",
							Detail:   "Metadata is required for Google Workspace IDP account with SAML auth type and API Creds Provided as NO",
						})

					}
				}
			}
		}
	}

	if len(diags) > 0 {
		return diags
	}

	log.Printf("[DEBUG] Creating IDP : %s, accountURL: (%s)", idp.IDPName, idp.AccountURL)
	idp.Details = *idpDetails

	if updateReq {
		if idp.IDPName != client.IDPGoogleWSAccount {
			_, err := prosimoClient.UpdateIDP(ctx, idp)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] Created IDP : %s, accountURL: (%s) with ID (%s)", idp.IDPName, idp.AccountURL, d.Id())
			d.SetId(d.Id())
		} else {
			_, err := prosimoClient.UpdateGoogleIDP(ctx, idp)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] Created IDP : %s, accountURL: (%s) with ID (%s)", idp.IDPName, idp.AccountURL, d.Id())
			d.SetId(d.Id())
		}
	}

	return resourceIDPRead(ctx, d, meta)
}

func resourceIDPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	idpID := d.Id()

	log.Printf("Get IDP for %s", idpID)

	idpList, err := prosimoClient.GetIDP(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var idp *client.IDP
	for _, returnedIDP := range idpList.IDPs {
		if returnedIDP.ID == idpID {
			idp = returnedIDP
			break
		}
	}
	if idp == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get IDP",
			Detail:   fmt.Sprintf("Unable to find IDP for ID %s", idpID),
		})

		return diags
	}

	idpAccount := idp.IDPName
	authType := idp.Auth_Type
	accountURL := idp.AccountURL
	selectType := idp.Select_Type
	apiCredProvided := idp.Api_Creds_Provided

	d.Set("idp_account", idpAccount)
	d.Set("auth_type", authType)
	d.Set("account_url", accountURL)
	d.Set("select_type", selectType)
	d.Set("api_cred_provided", apiCredProvided)

	idpDetails := idp.Details

	switch idpAccount {

	// okta account
	case client.IDPOktaAccount:
		if authType == client.IDPOIDCAuth { // okta + oidc

			oktaOIDCConfig := map[string]interface{}{}

			if idpDetails.APIToken != "" { // okta + oidc + api creds - api_token
				oktaOIDCConfig["api_token"] = idpDetails.APIToken
			} else {
				// okta + oidc + no api creds - client_id, secret_id

				oktaOIDCConfig["secret_id"] = idpDetails.ClientSecret
				oktaOIDCConfig["client_id"] = idpDetails.ClientID
			}

			var oktaOIDCList []map[string]interface{}
			oktaOIDCList = append(oktaOIDCList, oktaOIDCConfig)
			d.Set("oidc", oktaOIDCList)
		} else if authType == client.IDPSAMLAuth { // okta + saml

			oktaSAMLConfig := map[string]interface{}{}

			if idpDetails.APIToken != "" { // okta + oidc + api creds - api_token
				oktaSAMLConfig["api_token"] = idpDetails.APIToken
			} else {
				// okta + oidc + no api creds - client_id, secret_id

				oktaSAMLConfig["metadata_url"] = idpDetails.MetadataURL
				oktaSAMLConfig["metadata"] = idpDetails.Metadata
			}

			var oktaSAMLList []map[string]interface{}
			oktaSAMLList = append(oktaSAMLList, oktaSAMLConfig)
			d.Set("saml", oktaSAMLList)

		}

	// azure_ad account
	case client.IDPAzureADAccount:
		if authType == client.IDPOIDCAuth { // azure_ad + oidc

			// azure_ad + oidc - client_id, secret_id

			azureADOIDCConfig := map[string]interface{}{}
			azureADOIDCConfig["secret_id"] = idpDetails.ClientSecret
			azureADOIDCConfig["client_id"] = idpDetails.ClientID

			var azureOIDCList []map[string]interface{}
			azureOIDCList = append(azureOIDCList, azureADOIDCConfig)
			d.Set("oidc", azureOIDCList)
		} else if authType == client.IDPSAMLAuth { // azure_ad + saml
			// azure_ad + saml - metadata_url, metadata

			azureADSAMLConfig := map[string]interface{}{}
			azureADSAMLConfig["metadata_url"] = idpDetails.MetadataURL
			azureADSAMLConfig["metadata"] = idpDetails.Metadata

			var azureSAMLList []map[string]interface{}
			azureSAMLList = append(azureSAMLList, azureADSAMLConfig)
			d.Set("saml", azureSAMLList)

		}

	// one_login account
	case client.IDPOneLoginAccount:
		if authType == client.IDPOIDCAuth { // one_login + oidc

			// one_login + oidc - client_id, secret_id, region

			oneLoginOIDCConfig := map[string]interface{}{}
			oneLoginOIDCConfig["secret_id"] = idpDetails.ClientSecret
			oneLoginOIDCConfig["client_id"] = idpDetails.ClientID
			oneLoginOIDCConfig["region"] = idpDetails.Region

			var oneLoginOIDCList []map[string]interface{}
			oneLoginOIDCList = append(oneLoginOIDCList, oneLoginOIDCConfig)
			d.Set("oidc", oneLoginOIDCList)

		} else if authType == client.IDPSAMLAuth { // okta + saml
			// one_login + saml - metadata_url, metadata

			oneLoginSAMLConfig := map[string]interface{}{}
			oneLoginSAMLConfig["metadata_url"] = idpDetails.MetadataURL
			oneLoginSAMLConfig["metadata"] = idpDetails.Metadata

			var oneLoginSAMLList []map[string]interface{}
			oneLoginSAMLList = append(oneLoginSAMLList, oneLoginSAMLConfig)
			d.Set("saml", oneLoginSAMLList)

		}

		// one_login account
	case client.IDPOtherAccount:
		if authType == client.IDPOIDCAuth { // other + oidc

			// other + oidc - client_id, secret_id, region

			otherOIDCConfig := map[string]interface{}{}
			otherOIDCConfig["secret_id"] = idpDetails.ClientSecret
			otherOIDCConfig["client_id"] = idpDetails.ClientID

			var otherOIDCList []map[string]interface{}
			otherOIDCList = append(otherOIDCList, otherOIDCConfig)
			d.Set("oidc", otherOIDCList)

		} else if authType == client.IDPSAMLAuth { // okta + saml
			// other + saml - metadata_url, metadata

			otherSAMLConfig := map[string]interface{}{}
			otherSAMLConfig["metadata_url"] = idpDetails.MetadataURL
			otherSAMLConfig["metadata"] = idpDetails.Metadata

			var otherSAMLList []map[string]interface{}
			otherSAMLList = append(otherSAMLList, otherSAMLConfig)
			d.Set("saml", otherSAMLList)

		}
	case client.IDPPingOneAccount:
		if authType == client.IDPOIDCAuth { // pingone + oidc

			// pingone + oidc - client_id, secret_id, region

			pingOneOIDCConfig := map[string]interface{}{}
			pingOneOIDCConfig["secret_id"] = idpDetails.ClientSecret
			pingOneOIDCConfig["client_id"] = idpDetails.ClientID
			pingOneOIDCConfig["region"] = idpDetails.Region
			pingOneOIDCConfig["env_id"] = idpDetails.EnvID

			var pingOneOIDCList []map[string]interface{}
			pingOneOIDCList = append(pingOneOIDCList, pingOneOIDCConfig)
			d.Set("oidc", pingOneOIDCList)

		} else if authType == client.IDPSAMLAuth { // pingone + saml
			// pingone + saml - metadata_url, metadata

			pingOneSAMLConfig := map[string]interface{}{}
			pingOneSAMLConfig["secret_id"] = idpDetails.ClientSecret
			pingOneSAMLConfig["client_id"] = idpDetails.ClientID
			pingOneSAMLConfig["region"] = idpDetails.Region
			pingOneSAMLConfig["env_id"] = idpDetails.EnvID

			var pingOneSAMLList []map[string]interface{}
			pingOneSAMLList = append(pingOneSAMLList, pingOneSAMLConfig)
			d.Set("saml", pingOneSAMLList)

		}
	case client.IDPGoogleWSAccount:
		if authType == client.IDPOIDCAuth { // pingone + oidc

			// google + oidc - client_id, secret_id, region

			googleOIDCConfig := map[string]interface{}{}
			googleOIDCConfig["secret_id"] = idpDetails.ClientSecret
			googleOIDCConfig["client_id"] = idpDetails.ClientID
			googleOIDCConfig["domain"] = idpDetails.Domain
			googleOIDCConfig["customer_id"] = idpDetails.CustomerID
			googleOIDCConfig["admin_email"] = idpDetails.APIEmail
			googleOIDCConfig["api_file"] = idpDetails.APIFile

			var googleOIDCList []map[string]interface{}
			googleOIDCList = append(googleOIDCList, googleOIDCConfig)
			d.Set("oidc", googleOIDCList)

		} else if authType == client.IDPSAMLAuth { // pingone + saml
			// google + saml - metadata_url, metadata

			googleSAMLConfig := map[string]interface{}{}
			googleSAMLConfig["secret_id"] = idpDetails.ClientSecret
			googleSAMLConfig["client_id"] = idpDetails.ClientID
			googleSAMLConfig["domain"] = idpDetails.Domain
			googleSAMLConfig["customer_id"] = idpDetails.CustomerID
			googleSAMLConfig["admin_email"] = idpDetails.APIEmail
			googleSAMLConfig["api_file"] = idpDetails.APIFile

			var googleSAMLList []map[string]interface{}
			googleSAMLList = append(googleSAMLList, googleSAMLConfig)
			d.Set("saml", googleSAMLList)

		}
	}

	return diags
}

func resourceIDPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	// idpList := client.IDPList{}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	idpID := d.Id()

	idpList, err := prosimoClient.GetIDP(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var idp *client.IDP
	for _, returnedIDP := range idpList.IDPs {
		if returnedIDP.ID == idpID {
			idp = returnedIDP
			break
		}
	}

	if idp == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get IDP",
			Detail:   fmt.Sprintf("Unable to find IDP for ID %s", idpID),
		})

		return diags
	}

	if idp.Select_Type == client.IDPPartnerType && len(idp.AppIDs) > 0 {
		idp.AppIDs = nil
		idp.Details.APIClientID = ""
		idp.Details.APIClientSecret = ""
		idp.Details.Region = ""
		idp.Details.EnvID = ""
		idp.Details.ClientID = ""
		idp.Details.ClientSecret = ""
		idp.Details.APIToken = ""
		idp.Details.APIFile = ""
		idp.Details.APIEmail = ""
		idp.Details.Domain = ""
		idp.Details.Metadata = ""
		idp.Details.FilePath = ""
		idp.Details.CustomerID = ""

		// idp.Status = ""
		log.Println("[DEBUG] update idp", idp)
		if idp.IDPName != client.IDPGoogleWSAccount {
			_, err := prosimoClient.UpdateIDP(ctx, idp)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, err := prosimoClient.UpdateGoogleIDP(ctx, idp)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		err1 := prosimoClient.DeleteIDP(ctx, idpID)
		if err1 != nil {
			return diag.FromErr(err)
		}
	} else {
		err := prosimoClient.DeleteIDP(ctx, idpID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
