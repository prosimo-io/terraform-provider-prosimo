package prosimo

import (
	"context"
	"fmt"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCache() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify cache rules.",
		CreateContext: resourceCacheCreate,
		UpdateContext: resourceCacheUpdate,
		ReadContext:   resourceCacheRead,
		DeleteContext: resourceCacheDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of cache rule",
			},
			"bypass_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false, set it to true if you want to bypass cache.",
			},
			"teamid": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cache_control_ignored": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false, set it to true if you want to skip cache control.",
			},
			"share_static_content": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false, set it to true if you want to share static content.",
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"editable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"path_patterns": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Path pattern list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to store cache",
						},
						"is_default": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"bypass_uri": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"is_new_path": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"status": {
							Type:     schema.TypeString,
							Required: true,
						},
						"settings": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id_ignored": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(client.GetCacheType(), false),
									},
									"cache_control_ignored": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"cookie_ignored": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"query_parameter_ignored": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"ttl": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Required: true,
												},
												"time": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"time_unit": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(client.GetCacheTimeUnit(), false),
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
			"bypass_info": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resp_hdrs": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"x_jenkins_session": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"content_type": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"app_domains": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"is_new": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func inputDataops_cache(ctx context.Context, d *schema.ResourceData, meta interface{}) (diag.Diagnostics, client.CacheRule) {
	prosimoClient := meta.(*client.ProsimoClient)
	//prosimoClient := meta.(*client.ProsimoClient)
	//var diags diag.Diagnostics
	incache := &client.CacheRule{}
	pathList := []client.PATH{}

	inappdoainList := []client.AppDomaiN{}
	inbypassInfo := client.ByPassInfo{}
	inresphdrs := client.RespHdrs{}

	insettings := client.Settings{}
	inttl := client.Ttl{}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		incache.Name = name
	}
	if v, ok := d.GetOk("default"); ok {
		Default := v.(bool)
		incache.DefaultCache = Default
	}
	if v, ok := d.GetOk("editable"); ok {
		Editable := v.(bool)
		incache.Editable = Editable
	}
	if v, ok := d.GetOk("share_static_content"); ok {
		ShareStaticContent := v.(bool)
		incache.ShareStaticContent = ShareStaticContent
	}
	if v, ok := d.GetOk("cache_control_ignored"); ok {
		CacheControlIgnored := v.(bool)
		incache.CacheControlIgnored = CacheControlIgnored
	}
	if v, ok := d.GetOk("bypass_cache"); ok {
		BypassCache := v.(bool)
		incache.BypassCache = BypassCache
	}

	if v, ok := d.GetOk("bypassInfo"); ok {
		byPassInfo := v.(map[string]interface{})
		if v, ok := byPassInfo["respHdrs"].(*schema.Set); ok {
			respHdrDetails := v.List()[0].(map[string]interface{})
			if v, ok := respHdrDetails["x_jenkins_session"]; ok {
				sessionList := v.([]interface{})
				if len(sessionList) > 0 {
					inresphdrs.X_Jenkins_Session = expandStringList(v.([]interface{}))
				}
			}
			if v, ok := respHdrDetails["content-type"]; ok {
				sessionList := v.([]interface{})
				if len(sessionList) > 0 {
					inresphdrs.ContentType = expandStringList(v.([]interface{}))
				}
			}
			inbypassInfo.RespHrdrs = inresphdrs
		}
		incache.ByPassInfo = inbypassInfo
	}

	if v, ok := d.GetOk("app_domains"); ok {
		inappDomainList := v.([]interface{})
		for _, appdomain := range inappDomainList {
			inappdomain := client.AppDomaiN{}
			val := appdomain.(map[string]interface{})
			if domain, ok := val["domain"].(string); ok {
				inappdomain.AppDomain = domain
				existingAppDomainList, err := prosimoClient.GetAppDomains(ctx)
				if err != nil {
					return diag.FromErr(err), *incache
				}
				for _, exisingAppDomain := range existingAppDomainList {
					if domain == exisingAppDomain.Domain {
						inappdomain.AppDomainID = exisingAppDomain.ID
					}
				}
			}
			inappdoainList = append(inappdoainList, inappdomain)
			// log.Println("app domain list", inappdoainList)
		}
		incache.AppDomains = inappdoainList
	}

	if v, ok := d.GetOk("path_patterns"); ok {
		//inpathList := v.(*schema.Set).List()[0].(map[string]interface{})
		inpathList := v.([]interface{})
		for _, path := range inpathList {
			inPath := client.PATH{}
			val := path.(map[string]interface{})
			if path, ok := val["path"].(string); ok {
				inPath.Path = path
			}
			if bypassURI, ok := val["bypass_uri"].(bool); ok {
				inPath.ByPassURI = bypassURI
			}
			if isDefault, ok := val["is_default"].(bool); ok {
				inPath.IsDefault = isDefault
			}
			if status, ok := val["status"].(string); ok {
				inPath.Status = status
			}
			if isNewPath, ok := val["is_new_path"].(bool); ok {
				inPath.IsNewPath = isNewPath
			}
			if v, ok := val["settings"].(*schema.Set); ok {
				settingDetails := v.List()[0].(map[string]interface{})
				if types, ok := settingDetails["type"].(string); ok {
					if types == client.CacheTypeDynamic {
						types = client.CacheAPIInputDynamic
					} else if types == client.CacheTypeStaticLongLived {
						types = client.CacheAPIInputStaticLongLived
					} else if types == client.CacheTypeStaticShortLived {
						types = client.CacheAPIInputStaticShortLived
					}
					insettings.Type = types
				}
				if userIDIgnored, ok := settingDetails["user_id_gnored"].(bool); ok {
					insettings.UserIDIgnored = userIDIgnored
				}
				if cacheControlIgnored, ok := settingDetails["cache_control_ignored"].(bool); ok {
					insettings.CacheControlIgnored = cacheControlIgnored
				}
				if cookieIgnored, ok := settingDetails["cookie_ignored"].(bool); ok {
					insettings.CookieIgnored = cookieIgnored
				}
				if queryParameterIgnored, ok := settingDetails["uery_parameter_ignored"].(bool); ok {
					insettings.QueryParamaterIgnored = queryParameterIgnored
				}
				if v, ok := settingDetails["ttl"].(*schema.Set); ok {
					ttlDetails := v.List()[0].(map[string]interface{})
					if enabled, ok := ttlDetails["enabled"].(bool); ok {
						inttl.Enabled = enabled
					}
					if intime, ok := ttlDetails["time"].(int); ok {
						inttl.Time = intime
					}
					if timeUnit, ok := ttlDetails["time_unit"].(string); ok {
						if timeUnit == client.CacheTimeUnitHours {
							timeUnit = client.CacheAPITimeUnitHours
						} else if timeUnit == client.CacheTimeUnitMinutes {
							timeUnit = client.CacheAPITimeUnitMinutes
						} else if timeUnit == client.CacheTimeUnitSeconds {
							timeUnit = client.CacheAPITimeUnitSeconds
						} else if timeUnit == client.CacheTimeUnitDays {
							timeUnit = client.CacheAPITimeUnitDays
						}
						inttl.TimeUnit = timeUnit
					}
					insettings.TTL = inttl
				}
				inPath.Settings = insettings

			}
			pathList = append(pathList, inPath)
		}
		incache.PathPatterns = pathList

	}

	return nil, *incache
}

func resourceCacheCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	prosimoClient := meta.(*client.ProsimoClient)
	_, newcache := inputDataops_cache(ctx, d, meta)
	// log.Println("newcache", newcache)
	cacheData, err := prosimoClient.CreateCacheRule(ctx, &newcache)
	if err != nil {
		return diag.FromErr(err)
	}
	// log.Println("cache  data", cacheData)
	d.SetId(cacheData.Data.ID)
	//log.Println("policy id", policyListData.ID)

	return resourceCacheRead(ctx, d, meta)
	//return nil
}

func resourceCacheUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//resourcePolicyCreate(ctx, d, meta)
	prosimoClient := meta.(*client.ProsimoClient)
	cacheID := d.Id()
	_, newcache := inputDataops_cache(ctx, d, meta)
	newcache.ID = cacheID
	// log.Println("newcache", newcache)
	_, err := prosimoClient.UpdateCacheRule(ctx, &newcache)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cacheID)
	//log.Println("policy id", policyListData.ID)

	return resourceCacheRead(ctx, d, meta)
}

func resourceCacheRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics
	cacheID := d.Id()
	// log.Println("cacheId", cacheID)
	res, err := prosimoClient.GetCacheRuleByID(ctx, cacheID)

	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Cache Details",
			Detail:   fmt.Sprintf("Unable to find cache details for ID %s", cacheID),
		})

		return diags
	}
	d.Set("id", res.ID)
	d.Set("name", res.Name)
	d.Set("teamid", res.TeamID)

	return diags
}

func resourceCacheDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cacheID := d.Id()
	res_err := prosimoClient.DeleteCacheRule(ctx, cacheID)
	if res_err != nil {
		return diag.FromErr(res_err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
