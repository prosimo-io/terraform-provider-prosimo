package prosimo

import (
	"context"
	"log"
	"time"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProsimoManagedFirewalls() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create/modify Prosimo Managed Firewalls.",
		CreateContext: resourcePMFCreate,
		UpdateContext: resourcePMFUpdate,
		DeleteContext: resourcePMFDelete,
		ReadContext:   resourcePMFRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Resource",
			},
			"firewall_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of Firewall, e.g: vmseries",
			},
			"cloud_creds_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Account Name",
			},
			"cloud_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Region",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIDR range",
			},
			"instance_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance size to be filled in Instance details Section",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version to be filled in Instance details Section",
			},
			"auth_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Auth Key",
			},
			"auth_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Auth Key",
			},
			"bootstrap": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region Level IP Prefixes",
			},
			"scaling_settings": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Scaling Settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"desired": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Default Capacity",
						},
						"min": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Minimum Capacity",
						},
						"max": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum Capacity",
						},
					},
				},
			},
			"assignments": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Assignment Config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of Template",
						},
						"device_group": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Device Group Name",
						},
					},
				},
			},
			"access_details": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Access Details",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "UserName",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Password",
						},
						"select_option_for_ssh": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Should be one of Generate new key pair/Provide public key",
						},
						"key_pair_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of Keypair, applicable when `Existing key pair` option has been selected",
						},
						"public_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public Key details when selected option is `Provide public key`",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Satatus of deployment",
			},
			"onboard": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like to onboard Managed Firewall",
			},
			"decommission": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Set this to true if you would like to Decommission Managed Firewall",
			},
			"wait_for_rollout": {
				Type:        schema.TypeBool,
				Description: "Wait for the rollout of the task to complete. Defaults to true.",
				Default:     true,
				Optional:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourcePMFCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	scalingConfigInput := client.ScalingConfig{}
	assignmentsConfigInput := client.DeviceConfig{}
	accessConfigInput := client.AccessConfig{}
	prosimoClient := meta.(*client.ProsimoClient)
	offboardFlag := d.Get("decommission").(bool)
	if offboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid  decommission flag.",
			Detail:   "decommission can't be set to true while creating  resource.",
		})
		return diags
	}
	cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, d.Get("cloud_creds_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	activationConfig := &client.ActivationConfig{
		AuthKey:  d.Get("auth_key").(string),
		AuthCode: d.Get("auth_code").(string),
	}

	fwConfig := &client.FWConfig{
		Type:             d.Get("firewall_type").(string),
		Name:             d.Get("name").(string),
		CloudKeyID:       cloudCreds.ID,
		CloudRegion:      d.Get("cloud_region").(string),
		CIDR:             d.Get("cidr").(string),
		InstanceSize:     d.Get("instance_size").(string),
		Version:          d.Get("version").(string),
		ActivationConfig: activationConfig,
		Bootstrap:        d.Get("bootstrap").(string),
		ScalingConfig:    &scalingConfigInput,
		DeviceConfig:     &assignmentsConfigInput,
		AccessConfig:     &accessConfigInput,
	}
	if v, ok := d.GetOk("scaling_settings"); ok {
		scalingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
		scalingConfigInput = client.ScalingConfig{
			DefaultCapacity: scalingConfig["desired"].(int),
			MinCapacity:     scalingConfig["min"].(int),
			MaxCapacity:     scalingConfig["max"].(int),
		}
		log.Println("scalingConfigInput", scalingConfigInput)
	}
	if v, ok := d.GetOk("assignments"); ok {
		assignmentsConfig := v.(*schema.Set).List()[0].(map[string]interface{})
		assignmentsConfigInput = client.DeviceConfig{
			DGName:       assignmentsConfig["template_name"].(string),
			TPLStackname: assignmentsConfig["device_group"].(string),
		}
		log.Println("assignmentsConfig", assignmentsConfigInput)
	}
	if v, ok := d.GetOk("access_details"); ok {
		accessConfig := v.(*schema.Set).List()[0].(map[string]interface{})
		accessCreds := &client.AccessCred{
			UserName: accessConfig["username"].(string),
			PassWord: accessConfig["password"].(string),
		}
		accessConfigInput.AccessCreds = accessCreds
		pemDetailsInput := &client.PEMDetails{}
		if accessConfig["select_option_for_ssh"].(string) == "Generate new key pair" {
			pemDetailsInput.KeyGenerate = true
		} else if accessConfig["select_option_for_ssh"].(string) == "Provide public key" {
			if v, ok := accessConfig["public_key"].(string); ok {
				pemDetailsInput.PublicKey = v
			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Missing public key details",
					Detail:   "Missing public key details",
				})
				return diags
			}
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid Access details input",
				Detail:   "Invalid Access details input",
			})
			return diags
		}
		log.Println("pemDetailsInput", pemDetailsInput)
		accessConfigInput.PEMDetails = pemDetailsInput
	}

	log.Printf("[DEBUG]Creating FireWall : %v", fwConfig)
	createFirewall, err := prosimoClient.CreateFirewall(ctx, fwConfig)
	if err != nil {
		log.Printf("[ERROR] Error in creating firewall")
		return diag.FromErr(err)
	}
	d.SetId(createFirewall.PlMapRes.ID)
	if d.Get("onboard").(bool) {
		onboardresponse, err := prosimoClient.DeployFirewall(ctx, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[INFO] Waiting for task id %s to complete", onboardresponse.FWConfig.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskCompleteSharedService(ctx, d, meta, onboardresponse.FWConfig.TaskID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", onboardresponse.FWConfig.TaskID)
		}
	}

	return resourcePMFRead(ctx, d, meta)
}

func resourcePMFUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)

	fwOnboardFlag := d.Get("onboard").(bool)
	fwOffboardFlag := d.Get("decommission").(bool)
	if fwOnboardFlag && fwOffboardFlag {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid onboard and decommission flag combination.",
			Detail:   "Both onboard and decommission have been set to true.",
		})
		return diags
	}

	update_req := false
	if d.HasChange("cloud_creds_name") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Can't Modify  Cloud Account Name",
			Detail:   "Cloud Account Name can't be modified",
		})
		return diags
	}
	if d.HasChange("name") {
		update_req = true
	}
	if d.HasChange("firewall_type") {
		update_req = true
	}
	if d.HasChange("cloud_region") {
		update_req = true
	}
	if d.HasChange("cidr") {
		update_req = true
	}
	if d.HasChange("instance_size") {
		update_req = true
	}
	if d.HasChange("version") {
		update_req = true
	}
	if d.HasChange("auth_key") {
		update_req = true
	}
	if d.HasChange("auth_code") {
		update_req = true
	}
	if d.HasChange("bootstrap") {
		update_req = true
	}
	if d.HasChange("scaling_settings") {
		update_req = true
	}
	if d.HasChange("assignments") {
		update_req = true
	}
	if d.HasChange("access_details") {
		update_req = true
	}
	if d.HasChange("onboard") && !d.IsNewResource() {
		update_req = true
	}

	if d.HasChange("decommission") && !d.IsNewResource() {
		update_req = true
	}
	//Offboard Firewall
	if update_req {
		offBoardApp := false
		if d.HasChange("decommission") && !d.IsNewResource() {
			isDecommission := d.Get("decommission").(bool)
			if isDecommission {
				offBoardApp = true
				offboardresponse, err := prosimoClient.DecomFirewall(ctx, d.Id())
				if err != nil {
					return diag.FromErr(err)
				}
				if d.Get("wait_for_rollout").(bool) {
					log.Printf("[INFO] Waiting for task id %s to complete", offboardresponse.FWConfig.TaskID)
					err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
						retryUntilTaskComplete(ctx, d, meta, offboardresponse.FWConfig.TaskID))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[INFO] task %s is successful", offboardresponse.FWConfig.TaskID)
				}
			}
		}

		if !offBoardApp {
			scalingConfigInput := client.ScalingConfig{}
			assignmentsConfigInput := client.DeviceConfig{}
			accessConfigInput := client.AccessConfig{}
			prosimoClient := meta.(*client.ProsimoClient)
			offboardFlag := d.Get("decommission").(bool)
			if offboardFlag {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid  decommission flag.",
					Detail:   "decommission can't be set to true while creating  resource.",
				})
				return diags
			}
			cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, d.Get("cloud_creds_name").(string))
			if err != nil {
				return diag.FromErr(err)
			}
			activationConfig := &client.ActivationConfig{
				AuthKey:  d.Get("auth_key").(string),
				AuthCode: d.Get("auth_code").(string),
			}

			fwConfig := &client.FWConfig{
				ID:               d.Id(),
				Type:             d.Get("firewall_type").(string),
				Name:             d.Get("name").(string),
				CloudKeyID:       cloudCreds.ID,
				CloudRegion:      d.Get("cloud_region").(string),
				CIDR:             d.Get("cidr").(string),
				InstanceSize:     d.Get("instance_size").(string),
				Version:          d.Get("version").(string),
				ActivationConfig: activationConfig,
				Bootstrap:        d.Get("bootstrap").(string),
				ScalingConfig:    &scalingConfigInput,
				DeviceConfig:     &assignmentsConfigInput,
				AccessConfig:     &accessConfigInput,
			}
			if v, ok := d.GetOk("scaling_settings"); ok {
				scalingConfig := v.(*schema.Set).List()[0].(map[string]interface{})
				scalingConfigInput = client.ScalingConfig{
					DefaultCapacity: scalingConfig["desired"].(int),
					MinCapacity:     scalingConfig["min"].(int),
					MaxCapacity:     scalingConfig["max"].(int),
				}
				log.Println("scalingConfigInput", scalingConfigInput)
			}
			if v, ok := d.GetOk("assignments"); ok {
				assignmentsConfig := v.(*schema.Set).List()[0].(map[string]interface{})
				assignmentsConfigInput = client.DeviceConfig{
					DGName:       assignmentsConfig["template_name"].(string),
					TPLStackname: assignmentsConfig["device_group"].(string),
				}
				log.Println("assignmentsConfig", assignmentsConfigInput)
			}
			if v, ok := d.GetOk("access_details"); ok {
				accessConfig := v.(*schema.Set).List()[0].(map[string]interface{})

				accessCreds := &client.AccessCred{
					UserName: accessConfig["username"].(string),
					PassWord: accessConfig["password"].(string),
				}
				accessConfigInput.AccessCreds = accessCreds
				pemDetailsInput := &client.PEMDetails{}
				if accessConfig["select_option_for_ssh"].(string) == "Generate new key pair" {
					pemDetailsInput.KeyGenerate = true
				} else if accessConfig["select_option_for_ssh"].(string) == "Provide public key" {
					if v, ok := accessConfig["public_key"].(string); ok {
						pemDetailsInput.PublicKey = v
					} else {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Missing public key details",
							Detail:   "Missing public key details",
						})
						return diags
					}
				} else {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid Access details input",
						Detail:   "Invalid Access details input",
					})
					return diags
				}
				log.Println("pemDetailsInput", pemDetailsInput)
				accessConfigInput.PEMDetails = pemDetailsInput
			}

			log.Printf("[DEBUG]Updating FireWall : %v", fwConfig)
			updateFirewall, err := prosimoClient.UpdateFirewall(ctx, fwConfig)
			if err != nil {
				log.Printf("[ERROR] Error in updating firewall")
				return diag.FromErr(err)
			}
			if d.Get("onboard").(bool) {
				onboardresponse, err := prosimoClient.DeployFirewall(ctx, d.Id())
				if err != nil {
					return diag.FromErr(err)
				}
				if d.Get("wait_for_rollout").(bool) {
					log.Printf("[INFO] Waiting for task id %s to complete", onboardresponse.FWConfig.TaskID)
					err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
						retryUntilTaskCompleteManagedFirewall(ctx, d, meta, onboardresponse.FWConfig.TaskID))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("[INFO] task %s is successful", onboardresponse.FWConfig.TaskID)
				}
			}
			log.Printf("[DEBUG] Updated firewall - id - %s", updateFirewall.PlMapRes.ID)
		}
	}
	resourceRPRead(ctx, d, meta)
	return diags
}

func resourcePMFRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	var diags diag.Diagnostics

	nsID := d.Id()

	log.Printf("[DEBUG] Get firewall with id  %s", nsID)

	ns, err := prosimoClient.GetFirewallByID(ctx, nsID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", ns.ID)
	d.Set("cidr", ns.CIDR)
	d.Set("firewall_type", ns.Type)
	d.Set("cloud_region", ns.CloudRegion)
	d.Set("instance_size", ns.InstanceSize)
	d.Set("version", ns.Version)
	d.Set("status", ns.Status)
	return diags
}

func resourcePMFDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	prosimoClient := meta.(*client.ProsimoClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	rpID := d.Id()

	err := prosimoClient.DeleteFirewall(ctx, rpID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Deleted firewall with - id - %s", rpID)

	return diags
}
