package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
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
	name = "aws-u2n-euspoke3-tf"
	public_cloud {
		cloud_type        = "public"
		connection_option = "private"
		cloud_creds_name  = "prosimo-aws-app-iam"
		region_name       = "eu-west-1"
		cloud_networks {
		hub_id              = "tgw-06d2d8db5d344a1ed"
		vpc                 = "vpc-02669f01859cd3545"
		connectivity_type   = "transit-gateway"
		connector_placement = "none"
		subnets             = ["10.24.3.0/28","10.24.3.32/28"]
		}
	}
	policies         = ["ALLOW-ALL-NETWORKS"]
	onboard_app      = true
	decommission_app = false
	}
`)
}

func testAccResourceNetworkOnboardingPre() string {
	return fmt.Sprintf(`
resource "prosimo_network_onboarding" "test_update" {
	name = "aws-u2n-euspoke3-tf-new"
	public_cloud {
		cloud_type        = "public"
		connection_option = "private"
		cloud_creds_name  = "prosimo-aws-app-iam"
		region_name       = "eu-west-1"
		cloud_networks {
		hub_id              = "tgw-06d2d8db5d344a1ed"
		vpc                 = "vpc-02669f01859cd3545"
		connectivity_type   = "transit-gateway"
		connector_placement = "none"
		subnets             = ["10.24.3.0/28","10.24.3.32/28"]
		}
	}
	policies         = ["ALLOW-ALL-NETWORKS"]
	onboard_app      = true
	decommission_app = false
	}
`)
}

func testAccResourceNetworkOnboardingPost() string {
	return fmt.Sprintf(`
resource "prosimo_network_onboarding" "test_update" {
	name = "aws-u2n-euspoke3-tf-new"
	public_cloud {
		cloud_type        = "public"
		connection_option = "private"
		cloud_creds_name  = "prosimo-aws-app-iam"
		region_name       = "eu-west-1"
		cloud_networks {
		hub_id              = "tgw-06d2d8db5d344a1ed"
		vpc                 = "vpc-02669f01859cd3545"
		connectivity_type   = "transit-gateway"
		connector_placement = "none"
		subnets             = ["10.24.3.0/28","10.24.3.32/28"]
		}
	}
	policies         = ["DENY-ALL-NETWORKS"]
	onboard_app      = true
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
