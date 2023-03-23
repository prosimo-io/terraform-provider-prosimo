package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAppOnboardingFQDN_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingFQDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingFQDN_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingFQDNExists("prosimo_app_onboarding_fqdn.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceAppOnboardingFQDN_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingFQDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingFQDNPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingFQDNExists("prosimo_app_onboarding_fqdn.test_update"),
				),
			},
			{
				Config: testAccResourceAppOnboardingFQDNPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingFQDNExists("prosimo_app_onboarding_fqdn.test_update"),
				),
			},
		},
	})
}

func testAccResourceAppOnboardingFQDN_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_app_onboarding_fqdn" "test" {

	app_name = "common-app-agent-fqdn"
	idp_name = "azure_ad"
	app_urls {
		domain_type = "custom"
		app_fqdn = "alex-app-101.abc.com"
		subdomain_included = false

		protocols {
			protocol = "tcp"
			port_list = ["80", "90"]
		}

		health_check_info {
			enabled = false
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-gcp-infra"
			edge_regions {
				region_type = "active"
				region_name = "us-west2"
				conn_option = "public"
				backend_ip_address_discover = false
				backend_ip_address_manual = ["23.99.84.98"]
				# dns_custom {                                   
				#     dns_app = "agent-DNS-Server-tf"
				#     is_healthcheck_enabled = true
				#   }
			}
		}
	}
	saml_rewrite{
		selected_auth_type = "oidc"
	}
	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingFQDNPre() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_fqdn" "test_update" {

	app_name = "common-app-agent-fqdn-new"
	idp_name = "azure_ad"
	app_urls {
		domain_type = "custom"
		app_fqdn = "alex-app-101.abc.com"
		subdomain_included = false

		protocols {
			protocol = "tcp"
			port_list = ["80", "90"]
		}

		health_check_info {
			enabled = false
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-gcp-infra"
			edge_regions {
				region_type = "active"
				region_name = "us-west2"
				conn_option = "public"
				backend_ip_address_discover = false
				backend_ip_address_manual = ["23.99.84.98"]
				# dns_custom {                                   
				#     dns_app = "agent-DNS-Server-tf"
				#     is_healthcheck_enabled = true
				#   }
			}
		}
	}
	saml_rewrite{
		selected_auth_type = "oidc"
	}
	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingFQDNPost() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_fqdn" "test_update" {

	app_name = "common-app-agent-fqdn-new"
	idp_name = "azure_ad"
	app_urls {
		domain_type = "custom"
		app_fqdn = "alex-app-101.abc.com"
		subdomain_included = false

		protocols {
			protocol = "tcp"
			port_list = ["80", "90"]
		}

		health_check_info {
			enabled = false
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-gcp-infra"
			edge_regions {
				region_type = "active"
				region_name = "us-west2"
				conn_option = "public"
				backend_ip_address_discover = false
				backend_ip_address_manual = ["23.99.84.98"]
				# dns_custom {                                   
				#     dns_app = "agent-DNS-Server-tf"
				#     is_healthcheck_enabled = true
				#   }
			}
		}
	}
	saml_rewrite{
		selected_auth_type = "oidc"
	}
	optimization_option = "PerformanceEnhanced"

	policy_name = ["DENY-ALL-USERS"]

	onboard_app = false
	decommission_app = false
}
`)
}

func testAccCheckAppOnboardingFQDNExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID
		err, flag := client.GetAppOnboardingfiltered(name)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckAppOnboardingFQDNDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_app_onboarding_fqdn" {
			continue
		}
		_, flag := client.GetAppOnboardingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
