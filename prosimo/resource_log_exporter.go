package prosimo

import (
	"context"
	"log"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLogConfig() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify log exporter settings.",
		CreateContext: resourceLogConfigCreate,
		ReadContext:   resourceLogConfigRead,
		DeleteContext: resourceLogConfigDelete,
		UpdateContext: resourceLogConfigUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of log receiver endpoint",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address of log receiver endpoint",
			},
			"tcp_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "port of log receiver endpoint",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Defaults to false, set it true to enable tls verification",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description about log receiver",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Authentication token from receiver endpoint",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLogConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	logConfig := &client.Log_Config{}
	// auth := client.AUTH{}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		logConfig.Name = name
	}
	if v, ok := d.GetOk("ip"); ok {
		ip := v.(string)
		logConfig.IP = ip
	}
	if v, ok := d.GetOk("tcp_port"); ok {
		port := v.(int)
		logConfig.TcpPort = port
	}
	if v, ok := d.GetOk("tls_enabled"); ok {
		tlsStatus := v.(bool)
		logConfig.TlsEnabled = tlsStatus
	}
	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		logConfig.Description = description
	}
	if v, ok := d.GetOk("auth_token"); ok {
		authToken := v.(string)
		logConfig.AuthenticationToken = authToken
	}

	log.Printf("[DEBUG] Creating Log Config : %v", logConfig)
	createLogConfig, err := prosimoClient.CreateLogConf(ctx, logConfig)
	if err != nil {
		log.Printf("[ERROR] Error in creating Log config")
		return diag.FromErr(err)
	}
	d.SetId(createLogConfig.LogConfig.Id)
	return resourceLogConfigRead(ctx, d, meta)
}

func resourceLogConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	logID := d.Id()

	log.Printf("[DEBUG] Get LOG profile for %s", logID)

	res, err := prosimoClient.GetLogConf(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	var logConfig *client.Log_Config
	for _, returnedLogConfig := range res.LogConfigList {
		if returnedLogConfig.Id == logID {
			logConfig = returnedLogConfig
			break
		}
	}
	if logConfig == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DEBUG] Unable to get LOG config",
		})
		return diags
	}

	d.Set("team_id", logConfig.TeamID)
	d.Set("name", logConfig.Name)
	d.Set("ip", logConfig.IP)
	d.Set("tcp_port", logConfig.TcpPort)
	d.Set("description", logConfig.Description)
	d.Set("status", logConfig.Status)
	d.Set("tls_enabled", logConfig.TlsEnabled)

	return diags
}

func resourceLogConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	logConfig := &client.Log_Config{}
	// auth := client.AUTH{}

	updateReq := false
	logID := d.Id()
	logConfig.Id = logID

	if d.HasChange("name") && !d.IsNewResource() {
		updateReq = true
	}
	logConfig.Name = d.Get("name").(string)

	if d.HasChange("ip") && !d.IsNewResource() {
		updateReq = true
	}
	logConfig.IP = d.Get("ip").(string)

	if d.HasChange("tcp_port") && !d.IsNewResource() {
		updateReq = true
	}
	logConfig.TcpPort = d.Get("tcp_port").(int)

	if d.HasChange("description") && !d.IsNewResource() {
		updateReq = true
	}
	logConfig.Description = d.Get("description").(string)

	if d.HasChange("tls_enabled") && !d.IsNewResource() {
		updateReq = true
	}
	logConfig.TlsEnabled = d.Get("tls_enabled").(bool)

	if d.HasChange("auth_token") && !d.IsNewResource() {
		updateReq = true
	}
	logConfig.AuthenticationToken = d.Get("auth_token").(string)

	if len(diags) > 0 {
		return diags
	}

	if updateReq {
		log.Printf("[DEBUG] Updating LOG Config : %v", logConfig)
		_, err := prosimoClient.UpdateLogConf(ctx, logConfig)
		if err != nil {
			log.Printf("[ERROR] Error in updating LOGconfig")
			return diag.FromErr(err)
		}
		// d.SetId(updateEDR.EdrConfig.Id)
	}
	return resourceLogConfigRead(ctx, d, meta)
}

func resourceLogConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)
	var diags diag.Diagnostics

	logID := d.Id()
	err := prosimoClient.DeleteLogConf(ctx, logID)

	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
