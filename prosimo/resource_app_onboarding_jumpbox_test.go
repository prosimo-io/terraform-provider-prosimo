package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAppOnboardingJumpBox_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingJumpBoxDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingJumpBox_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingJumpBoxExists("prosimo_app_onboarding_jumpbox.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceAppOnboardingJumpBox_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingJumpBoxDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingJumpBoxPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingJumpBoxExists("prosimo_app_onboarding_jumpbox.test_update"),
				),
			},
			{
				Config: testAccResourceAppOnboardingJumpBoxPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingJumpBoxExists("prosimo_app_onboarding_jumpbox.test_update"),
				),
			},
		},
	})
}

func testAccResourceAppOnboardingJumpBox_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_app_onboarding_jumpbox" "test" {

	app_name = "common-app-jumpbox"
	idp_name = "azure_ad"
	app_urls {
		internal_domain = "tf-jumpbox.psonar.us"
		app_fqdn = "tf-jumpbox.psonar.us"
		cloud_config {
			connection_option = "private"
			cloud_creds_name = "prosimo-aws-app-iam"
			edge_regions {
				region_name = "us-east-2"
				region_type = "active"
				conn_option = "peering"
				backend_ip_address_discover = false
				backend_ip_address_manual = ["10.100.0.142"]
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
	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingJumpBoxPre() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_jumpbox" "test_update" {

	app_name = "common-app-jumpbox-new"
	idp_name = "azure_ad"
	app_urls {
		internal_domain = "tf-jumpbox.psonar.us"
		app_fqdn = "tf-jumpbox.psonar.us"
		cloud_config {
			connection_option = "private"
			cloud_creds_name = "prosimo-aws-app-iam"
			edge_regions {
				region_name = "us-east-2"
				region_type = "active"
				conn_option = "peering"
				backend_ip_address_discover = false
				backend_ip_address_manual = ["10.100.0.142"]
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
	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingJumpBoxPost() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_jumpbox" "test_update" {

	app_name = "common-app-jumpbox-new"
	idp_name = "azure_ad"
	app_urls {
		internal_domain = "tf-jumpbox.psonar.us"
		app_fqdn = "tf-jumpbox.psonar.us"
		cloud_config {
			connection_option = "private"
			cloud_creds_name = "prosimo-aws-app-iam"
			edge_regions {
				region_name = "us-east-2"
				region_type = "active"
				conn_option = "peering"
				backend_ip_address_discover = false
				backend_ip_address_manual = ["10.100.0.142"]
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
	policy_name = ["DENY-ALL-USERS"]
	onboard_app = false
	decommission_app = false
}
`)
}

func testAccCheckAppOnboardingJumpBoxExists(resource string) resource.TestCheckFunc {
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

func testAccCheckAppOnboardingJumpBoxDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_app_onboarding_jumpbox" {
			continue
		}
		_, flag := client.GetAppOnboardingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
