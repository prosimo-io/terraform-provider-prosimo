package prosimo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIpReputation_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIpReputation_basic(name),
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckIpReputationExists("prosimo_ip_reputation.test"),
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

func TestAccResourceIpReputation_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIpReputationPre(),
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckIpReputationExists("prosimo_ip_reputation.test_update"),
				),
			},
			{
				Config: testAccResourceIpReputationPost(),
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckIpReputationExists("prosimo_ip_reputation.test_update"),
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

func testAccResourceIpReputation_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_ip_reputation" "test" {
		enabled        = true
		allowlist = ["1.1.1.3/16", "1.1.1.2/16", "1.1.1.4/16"]
	  }
`)
}

func testAccResourceIpReputationPre() string {
	return fmt.Sprintf(`
	resource "prosimo_ip_reputation" "test_update" {
		enabled        = true
		allowlist = ["1.1.1.3/16", "1.1.1.2/16"]
	  }
`)
}

func testAccResourceIpReputationPost() string {
	return fmt.Sprintf(`
	resource "prosimo_ip_reputation" "test_update" {
		enabled        = true
		allowlist = ["1.1.1.3/16", "1.1.1.2/16", "1.1.1.4/16"]
	  }
`)
}
