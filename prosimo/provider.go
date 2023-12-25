package prosimo

import (
	"context"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROSIMO_TOKEN", nil),
				Description: descriptions["token"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROSIMO_BASE_URL", nil),
				Description: descriptions["base_url"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROSIMO_INSECURE", true),
				Description: descriptions["insecure"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"prosimo_cloud_creds":              resourceCloudCreds(),
			"prosimo_edge":                     resourceEdge(),
			"prosimo_idp":                      resourceIDP(),
			"prosimo_app_onboarding_web":       resourceAppOnboarding_Web(),
			"prosimo_app_onboarding_fqdn":      resourceAppOnboarding_FQDN(),
			"prosimo_app_onboarding_dns":       resourceAppOnboarding_DNS(),
			"prosimo_app_onboarding_cloudsvc":  resourceAppOnboarding_CloudSVC(),
			"prosimo_app_onboarding_jumpbox":   resourceAppOnboarding_JumpBox(),
			"prosimo_app_onboarding_citrixvdi": resourceAppOnboarding_VDI(),
			"prosimo_network_onboarding":       resourceNetworkOnboarding(),
			"prosimo_waf":                      resourceWAF(),
			"prosimo_ip_reputation":            resourceIPReputation(),
			"prosimo_dynamic_risk":             resourceDynamicRisk(),
			"prosimo_geo_location":             resourceGeoLocation(),
			"prosimo_user_settings":            resourceUserSettings(),
			"prosimo_certificates":             resourceCertficate(),
			"prosimo_policy":                   resourcePolicy(),
			"prosimo_cache_rules":              resourceCache(),
			"prosimo_edr_profile":              resourceEdrProfile(),
			"prosimo_edr_integration":          resourceEdrIntegration(),
			"prosimo_dp_profile":               resourceDPProfile(),
			"prosimo_dp_settings":              resourceDPSettings(),
			"prosimo_log_exporter":             resourceLogConfig(),
			"prosimo_grouping":                 resourceGrpConfig(),
			"prosimo_connector_placement":      resourceConnPlacement(),
			"prosimo_shared_services":          resourceSharedServices(),
			"prosimo_service_insertion":        resourceServiceInsertion(),
			"prosimo_private_link_source":      resourcePrivateLinkSource(),
			"prosimo_private_link_mapping":     resourcePrivateLinkMapping(),
			"prosimo_namespace":                resourceNamespace(),
			"prosimo_cloud_gateway":            resourceCloudGateway(),
			"prosimo_visual_transit":           resourceVisualTransit(),
			"prosimo_regional_prefix":          resourceRegionalPrefix(),
			"prosimo_network_prefix":           resourceNetworkPrefix(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"prosimo_cloud_creds":         dataSourceCloudCreds(),
			"prosimo_edge":                dataSourceEdge(),
			"prosimo_idp":                 dataSourceIDP(),
			"prosimo_s3bucket":            dataSourceS3bucket(),
			"prosimo_app_onboarding":      dataSourceAppOnboarding(),
			"prosimo_certificates":        datasourceCertficate(),
			"prosimo_policy_access":       datasourcePolicyAccess(),
			"prosimo_policy_transit":      datasourcePolicyTransit(),
			"prosimo_network_onboarding":  dataSourceNetworkOnboarding(),
			"prosimo_discovered_networks": dataSourceNetworkDiscovered(),
			// "prosimo_waf_policy":          dataSourceWAFPolicy(),
		},

		// DataSourcesList: []*schema.Resouece{
		// 	"prosimo_s3bucket": dataSourceS3bucket(),
		// },
		ConfigureContextFunc: providerConfigure,
	}
	return p
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The API token used to connect to Prosimo. ",

		"base_url": "The Prosimo Base API URL",

		"insecure": "Defaults to False.Enable `insecure` mode only for testing purposes",
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiToken := d.Get("token").(string)
	baseURL := d.Get("base_url").(string)
	insecure := d.Get("insecure").(bool)

	insecure = true

	prosimoClient, err := client.NewProsimoClient(baseURL, apiToken, insecure)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Prosimo client",
			Detail:   "Unable to authenticate token for Prosimo client",
		})

		return nil, diags
	}

	return prosimoClient, diags
}
