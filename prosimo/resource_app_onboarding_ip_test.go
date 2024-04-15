package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAppOnboardingIP_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingIP_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingIPExists("prosimo_app_onboarding_ip.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceAppOnboardingIP_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingIPPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingIPExists("prosimo_app_onboarding_ip.test_update"),
				),
			},
			{
				Config: testAccResourceAppOnboardingIPPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingIPExists("prosimo_app_onboarding_ip.test_update"),
				),
			},
		},
	})
}

func testAccResourceAppOnboardingIP_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_app_onboarding_ip" "test" {

	app_name = "agent-test"
	// idp_name = "azure_ad"
	app_urls {
		app_fqdn = "10.10.0.36"
		protocols {
			protocol = "tcp"
			port_list = ["443", "22"]
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-infra"
			edge_regions {
				region_type = "active"
				region_name = "us-east-2"
				conn_option = "public"				
			}
		}
	}
	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingIPPre() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_ip" "test_update" {

	app_name = "agent-test_update"
	// idp_name = "azure_ad"
	app_urls {
		app_fqdn = "10.10.0.36"
		protocols {
			protocol = "tcp"
			port_list = ["443", "22"]
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-infra"
			edge_regions {
				region_type = "active"
				region_name = "us-east-2"
				conn_option = "public"				
			}
		}
	}
	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingIPPost() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_ip" "test_update" {

	app_name = "agent-test_update"
	// idp_name = "azure_ad"
	app_urls {
		app_fqdn = "10.10.0.36"
		protocols {
			protocol = "tcp"
			port_list = ["443"]
		}

		cloud_config {
			connection_option = "public"
			cloud_creds_name = "prosimo-infra"
			edge_regions {
				region_type = "active"
				region_name = "us-east-2"
				conn_option = "public"				
			}
		}
	}
	optimization_option = "PerformanceEnhanced"

	policy_name = ["ALLOW-ALL-USERS"]

	onboard_app = false
	decommission_app = false
}
`)
}

func testAccCheckAppOnboardingIPExists(resource string) resource.TestCheckFunc {
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

func testAccCheckAppOnboardingIPDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_app_onboarding_ip" {
			continue
		}
		_, flag := client.GetAppOnboardingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
