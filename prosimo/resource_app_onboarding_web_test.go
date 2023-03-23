package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAppOnboardingWeb_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingWebDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingWeb_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingWebExists("prosimo_app_onboarding_web.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceAppOnboardingWeb_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingWebDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingWebPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingWebExists("prosimo_app_onboarding_web.test_update"),
				),
			},
			{
				Config: testAccResourceAppOnboardingWebPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingWebExists("prosimo_app_onboarding_web.test_update"),
				),
			},
		},
	})
}

func testAccResourceAppOnboardingWeb_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_app_onboarding_web" "test" {

	app_name = "agentless-multi-VMs-tf"
	idp_name = "azure_ad"
	app_urls {
		internal_domain = "10.100.0.142"
		domain_type = "custom"
		app_fqdn = "alex-app-102.abc.com"
		subdomain_included = false

		protocols {
			protocol = "ssh"
			port = 22
		}

		health_check_info {
			enabled = true
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-aws-iam"
			edge_regions {
				region_type = "active"
				region_name = "us-west-1"
				conn_option = "public"
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

	onboard_app =false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingWebPre() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_web" "test_update" {

	app_name = "agentless-multi-VMs-tf-new"
	idp_name = "azure_ad"
	app_urls {
		internal_domain = "10.100.0.142"
		domain_type = "custom"
		app_fqdn = "alex-app-102.abc.com"
		subdomain_included = false

		protocols {
			protocol = "ssh"
			port = 22
		}

		health_check_info {
			enabled = true
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-aws-iam"
			edge_regions {
				region_type = "active"
				region_name = "us-west-1"
				conn_option = "public"
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

	onboard_app =false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingWebPost() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_web" "test_update" {

	app_name = "agentless-multi-VMs-tf-new"
	idp_name = "azure_ad"
	app_urls {
		internal_domain = "10.100.0.142"
		domain_type = "custom"
		app_fqdn = "alex-app-102.abc.com"
		subdomain_included = false

		protocols {
			protocol = "ssh"
			port = 22
		}

		health_check_info {
			enabled = true
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-aws-iam"
			edge_regions {
				region_type = "active"
				region_name = "us-west-1"
				conn_option = "public"
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

	onboard_app =false
	decommission_app = false
}
`)
}

func testAccCheckAppOnboardingWebExists(resource string) resource.TestCheckFunc {
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

func testAccCheckAppOnboardingWebDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_app_onboarding_web" {
			continue
		}
		_, flag := client.GetAppOnboardingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
