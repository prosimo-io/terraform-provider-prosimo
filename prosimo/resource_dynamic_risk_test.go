package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDynamicRisk_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynamicRiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDynamicRisk_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynamicRiskExists("prosimo_dynamic_risk.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      "prosimo_dynamic_risk.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDynamicRisk_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynamicRiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDynamicRisk_pre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynamicRiskExists("prosimo_dynamic_risk.test_update"),
				),
			},
			{
				Config: testAccResourceDynamicRisk_post(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynamicRiskExists("prosimo_dynamic_risk.test_update"),
				),
			},
			{
				ResourceName:      "prosimo_dynamic_risk.test_update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceDynamicRisk_basic(name string) string {
	return fmt.Sprintf(` 

resource "prosimo_dynamic_risk" "test" {
	threshold  {
		name  =  "alert" 
		enabled  =  true
		value  =  80
	}
		threshold  {
		name  =  "mfa" 
		enabled  =  true
		value  =  80
	}
		threshold  {
		name  =  "lockUser" 
		enabled  =  false
		value  =  100
	}
	} 
`)
}

func testAccResourceDynamicRisk_pre() string {
	return fmt.Sprintf(` 

resource "prosimo_dynamic_risk" "test_update" {
	threshold  {
		name  =  "alert" 
		enabled  =  true
		value  =  81
	}
		threshold  {
		name  =  "mfa" 
		enabled  =  true
		value  =  81
	}
		threshold  {
		name  =  "lockUser" 
		enabled  =  false
		value  =  100
	}
	} 
`)
}

func testAccResourceDynamicRisk_post() string {
	return fmt.Sprintf(` 

resource "prosimo_dynamic_risk" "test_update" {
	threshold  {
		name  =  "alert" 
		enabled  =  true
		value  =  82
	}
		threshold  {
		name  =  "mfa" 
		enabled  =  true
		value  =  82
	}
		threshold  {
		name  =  "lockUser" 
		enabled  =  false
		value  =  100
	}
	} 
`)
}

func testAccCheckDynamicRiskExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID
		err, flag := client.GetDynamicRiskfiltered(name)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckDynamicRiskDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_dynamic_risk" {
			continue
		}
		_, flag := client.GetDynamicRiskfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
