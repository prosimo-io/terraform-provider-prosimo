package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNetworkOnboarding_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkOnboardingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetworkOnboarding_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkOnboardingExists("prosimo_network_onboarding.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
			// {
			// 	ResourceName: "prosimo_app_onboarding_citrixvdi.test",
			// 	ImportState:  true,
			// 	// ImportStateVerify: true,
			// },
		},
	})
}

func TestAccResourceNetworkOnboarding_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkOnboardingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetworkOnboardingPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkOnboardingExists("prosimo_network_onboarding.test_update"),
				),
			},
			{
				Config: testAccResourceNetworkOnboardingPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkOnboardingExists("prosimo_network_onboarding.test_update"),
				),
			},
			// {
			// 	ResourceName: "prosimo_app_onboarding_citrixvdi.test_update",
			// 	ImportState:  true,
			// 	// ImportStateVerify: true,
			// },
		},
	})
}

func testAccResourceNetworkOnboarding_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_network_onboarding" "test" {

	name = "demo_network_new"
	public_cloud {
		cloud_type = "public"
		connection_option = "private"
		cloud_creds_name = "prosimo-aws-iam"
		region_name = "us-west-2"
		cloud_networks {
			vpc = "vpc-0748dd68349f5ada2"
			# hub_id = "tgw-04db5eac6fe3de45e"
			connector_placement = "none"
			connectivity_type = "vpc-peering"
			subnets = ["10.0.0.1/24"]
		}
		connect_type = "connector"

	}
	policies = ["DENY-ALL-NETWORKS"]
	onboard_app = false
	decommission_app = false
	}
`)
}

func testAccResourceNetworkOnboardingPre() string {
	return fmt.Sprintf(`
resource "prosimo_network_onboarding" "test_update" {

	name = "demo_network_new"
	public_cloud {
		cloud_type = "public"
		connection_option = "private"
		cloud_creds_name = "prosimo-aws-iam"
		region_name = "us-west-2"
		cloud_networks {
			vpc = "vpc-0748dd68349f5ada2"
			# hub_id = "tgw-04db5eac6fe3de45e"
			connector_placement = "none"
			connectivity_type = "vpc-peering"
			subnets = ["10.0.0.1/24"]
		}
		connect_type = "connector"

	}
	policies = ["DENY-ALL-NETWORKS"]
	onboard_app = false
	decommission_app = false
	}
`)
}

func testAccResourceNetworkOnboardingPost() string {
	return fmt.Sprintf(`
resource "prosimo_network_onboarding" "test_update" {

	name = "demo_network_new"
	public_cloud {
		cloud_type = "public"
		connection_option = "private"
		cloud_creds_name = "prosimo-aws-iam"
		region_name = "us-west-2"
		cloud_networks {
			vpc = "vpc-0748dd68349f5ada2"
			# hub_id = "tgw-04db5eac6fe3de45e"
			connector_placement = "none"
			connectivity_type = "vpc-peering"
			subnets = ["10.0.0.1/24"]
		}
		connect_type = "connector"

	}
	policies = ["ALLOW-ALL-NETWORKS"]
	onboard_app = false
	decommission_app = false
	}
`)
}

func testAccCheckNetworkOnboardingExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		err, flag := client.GetNetworkOnboardingfiltered(rs.Primary.ID)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %p", resource, err)
		}
		return nil
	}
}

func testAccCheckNetworkOnboardingDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_network_onboarding" {
			continue
		}
		_, flag := client.GetNetworkOnboardingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
