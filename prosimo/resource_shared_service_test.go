package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceSharedService_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSharedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSharedService_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSharedServiceExists("prosimo_shared_services.test"),
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

func TestAccResourceSharedService_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSharedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSharedServicePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSharedServiceExists("prosimo_shared_services.test_update"),
				),
			},
			{
				Config: testAccResourceSharedServicePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSharedServiceExists("prosimo_shared_services.test_update"),
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

func testAccResourceSharedService_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_shared_services" "test" {
		name = "firewall_svc"
		region {
			cloud_region = "us-west-2"
			gateway_lb = "com.amazonaws.vpce.us-west-2.vpce-svc-0fdda8395ea088814"
			cloud_creds_name = "prosimo-aws-app-iam"
		}
		onboard = true
		decommission = false
	}
`)
}

func testAccResourceSharedServicePre() string {
	return fmt.Sprintf(`
	resource "prosimo_shared_services" "test_update" {
		name = "firewall_svc"
		region {
			cloud_region = "us-west-2"
			gateway_lb = "com.amazonaws.vpce.us-west-2.vpce-svc-0fdda8395ea088814"
			cloud_creds_name = "prosimo-aws-app-iam"
		}
		onboard = true
		decommission = false
	}
`)
}

func testAccResourceSharedServicePost() string {
	return fmt.Sprintf(`
	resource "prosimo_shared_services" "test_update" {
		name = "firewall_svc"
		region {
			cloud_region = "us-west-2"
			gateway_lb = "com.amazonaws.vpce.us-west-2.vpce-svc-0fdda8395ea088814"
			cloud_creds_name = "prosimo-aws-app-iam"
		}
		onboard = false
		decommission = true
	}
`)
}

func testAccCheckSharedServiceExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		err, flag := client.GetSharedServicefiltered(rs.Primary.ID)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %p", resource, err)
		}
		return nil
	}
}

func testAccCheckSharedServiceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_shared_services" {
			continue
		}
		_, flag := client.GetSharedServicefiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
