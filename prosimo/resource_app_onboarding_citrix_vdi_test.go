package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAppOnboardingCitrixVDI_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingCitrixVDIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingCitrixVDI_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingCitrixVDIExists("prosimo_app_onboarding_citrixvdi.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceAppOnboardingCitrixVDI_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingCitrixVDIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingCitrixVDIPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingCitrixVDIExists("prosimo_app_onboarding_citrixvdi.test_update"),
				),
			},
			{
				Config: testAccResourceAppOnboardingCitrixVDIPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingCitrixVDIExists("prosimo_app_onboarding_citrixvdi.test_update"),
				),
			},
		},
	})
}

func testAccResourceAppOnboardingCitrixVDI_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_app_onboarding_citrixvdi" "test" {

	app_name = "common-app"
	// idp_name = "azure_ad"
	citrix_ip = ["10.1.1.1"]
	app_urls {
		internal_domain = "google.com"
		domain_type = "custom"
		app_fqdn = "google.com"
		subdomain_included = false

		protocols {
			protocol = "http"
			port = 80
		}

		health_check_info {
			enabled = false
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-aws-iam"
			edge_regions {
				region_type = "active"
				region_name = "us-west-2"
				conn_option = "public"
				backend_ip_address_discover = true
			}
		}
		dns_service {
			type = "manual"
		}

		ssl_cert {
			generate_cert = true
		}
	}

	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app =false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingCitrixVDIPre() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_citrixvdi" "test_update" {

	app_name = "common-app-new"
	// idp_name = "azure_ad"
	citrix_ip = ["10.1.1.2"]
	app_urls {
		internal_domain = "amazon.com"
		domain_type = "custom"
		app_fqdn = "amazon.com"
		subdomain_included = false

		protocols {
			protocol = "http"
			port = 80
		}

		health_check_info {
			enabled = false
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-aws-iam"
			edge_regions {
				region_type = "active"
				region_name = "us-west-2"
				conn_option = "public"
				backend_ip_address_discover = true
			}
		}
		dns_service {
			type = "manual"
		}

		ssl_cert {
			generate_cert = true
		}
	}

	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app =false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingCitrixVDIPost() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_citrixvdi" "test_update" {

	app_name = "common-app-new"
	// idp_name = "azure_ad"
	citrix_ip = ["10.1.1.1"]
	app_urls {
		internal_domain = "amazon.com"
		domain_type = "custom"
		app_fqdn = "amazon.com"
		subdomain_included = false

		protocols {
			protocol = "http"
			port = 80
		}

		health_check_info {
			enabled = false
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-aws-iam"
			edge_regions {
				region_type = "active"
				region_name = "us-west-2"
				conn_option = "public"
				backend_ip_address_discover = true
			}
		}
		dns_service {
			type = "manual"
		}

		ssl_cert {
			generate_cert = true
		}
	}

	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app =false
	decommission_app = false
}
`)
}

func testAccCheckAppOnboardingCitrixVDIExists(resource string) resource.TestCheckFunc {
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

func testAccCheckAppOnboardingCitrixVDIDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_app_onboarding_citrixvdi" {
			continue
		}
		_, flag := client.GetAppOnboardingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
