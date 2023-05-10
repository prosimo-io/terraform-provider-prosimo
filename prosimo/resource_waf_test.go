package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceWaf_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWaf_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafExists("prosimo_waf.test"),
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

func TestAccResourceWaf_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWafPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafExists("prosimo_waf.test_update"),
				),
			},
			{
				Config: testAccResourceWafPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafExists("prosimo_waf.test_update"),
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

func testAccResourceWaf_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_waf" "test" {
		waf_name        = "terraform-waf"
		mode = "enforce"
		threshold = 30
		
		rulesets {
			basic {
				rule_groups = ["11000_whitelist"]
			}
	  
			owasp_crs_v32 {
				rule_groups = ["REQUEST-903.9001-DRUPAL-EXCLUSION-RULES"]
			}
		}
	  
		app_domains = ["google.com"]
	  
	  }
	  
`)
}

func testAccResourceWafPre() string {
	return fmt.Sprintf(`
	resource "prosimo_waf" "test_update" {
		waf_name        = "terraform-waf"
		mode = "enforce"
		threshold = 30
		
		rulesets {
			basic {
				rule_groups = ["11000_whitelist"]
			}
	  
			owasp_crs_v32 {
				rule_groups = ["REQUEST-903.9001-DRUPAL-EXCLUSION-RULES"]
			}
		}
	  
		app_domains = ["yahoo.com"]
	  
	  }
`)
}

func testAccResourceWafPost() string {
	return fmt.Sprintf(`
	resource "prosimo_waf" "test_update" {
		waf_name        = "terraform-waf"
		mode = "enforce"
		threshold = 30
		
		rulesets {
			basic {
				rule_groups = ["11000_whitelist"]
			}
	  
			owasp_crs_v32 {
				rule_groups = ["REQUEST-903.9001-DRUPAL-EXCLUSION-RULES"]
			}
		}
	  
		app_domains = ["google.com"]
	  
	  }
`)
}

func testAccCheckWafExists(resource string) resource.TestCheckFunc {
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

func testAccCheckWafDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_waf" {
			continue
		}
		_, flag := client.GetWaffiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
